// server.go
package auth

import (
	"context"
	"crypto/rsa" // Ensure rsa is imported
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"

	"github.com/coreos/go-oidc/v3/oidc"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"

	// Use v4 or v5 consistently based on your go.mod
	"github.com/golang-jwt/jwt" // Or jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/db" // Local DB package
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
	"gorm.io/gorm" // Import gorm if checking for ErrRecordNotFound
)

// User struct for internal use, e.g., in UpdateLastLoginLoop
type User struct {
	ID         string // Typically ExternalId
	Email      string
	ExternalId string
	Role       api.Role
	LastLogin  time.Time
	CreatedAt  time.Time
}

// Server holds dependencies for the auth server logic (gRPC Check, Verify)
type Server struct {
	host                string
	platformPublicKey   *rsa.PublicKey
	platformKeyID       string // Stores the calculated Key ID (JWK Thumbprint)
	dexVerifier         *oidc.IDTokenVerifier
	dexClient           dexApi.DexClient
	logger              *zap.Logger
	db                  db.Database // Use the specific db.Database type
	updateLoginUserList []User      // Deprecated if using map below
	updateLogin         chan User   // Channel for queuing login updates
}

// userClaim represents the claims verified or derived by this auth service.
type userClaim struct {
	Role           api.Role
	Email          string
	MemberSince    *time.Time // Populated after DB lookup
	UserLastLogin  *time.Time // Populated after DB lookup
	ExternalUserID string     // Usually corresponds to 'sub' claim
	EmailVerified  bool

	// Embed standard claims for JWTs that use them (e.g., platform tokens)
	jwt.StandardClaims
}

// Valid implements the jwt.Claims interface. Add validation if needed.
func (u *userClaim) Valid() error {
	// Example: Check expiry from standard claims if present
	// if u.ExpiresAt != 0 {
	// 	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
	// 		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	// 	}
	// }
	// Add other standard claim validations (nbf, iat) if necessary
	return nil
}

// UpdateLastLoginLoop processes queued users to update their last login timestamp.
func (s *Server) UpdateLastLoginLoop() {
	ticker := time.NewTicker(1 * time.Minute) // Check less frequently
	defer ticker.Stop()
	usersToUpdate := make(map[string]User) // Use map for efficient deduplication

	for {
		select {
		case user := <-s.updateLogin:
			if user.ExternalId != "" { // Ensure we have an ID to work with
				usersToUpdate[user.ExternalId] = user // Add/overwrite user in map
				s.logger.Debug("User added/updated in last login queue", zap.String("externalId", user.ExternalId))
			} else {
				s.logger.Warn("Received user for last login update with empty ExternalId", zap.String("email", user.Email))
			}
		case <-ticker.C:
			if len(usersToUpdate) == 0 {
				continue
			}
			s.logger.Debug("Processing user last login update batch", zap.Int("queue_size", len(usersToUpdate)))

			// Process users currently in the map
			processedIDs := []string{} // Keep track of processed IDs in this batch
			for extId := range usersToUpdate {
				processedIDs = append(processedIDs, extId) // Mark for processing

				// Get current user state from DB
				dbUser, err := s.db.GetUserByExternalID(extId) // Use the correct DB method
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						s.logger.Warn("User no longer exists in DB, cannot update last login", zap.String("externalId", extId))
					} else {
						s.logger.Error("Failed to get user from DB for last login update", zap.String("externalId", extId), zap.Error(err))
					}
					continue // Skip this user for now if error or not found
				}
				if dbUser == nil { // Safeguard
					s.logger.Warn("GetUserByExternalID returned nil user without error", zap.String("externalId", extId))
					continue
				}

				// Check if update is needed (e.g., > 15 mins since last recorded login)
				// Allow configurable interval?
				updateInterval := 15 * time.Minute
				if time.Since(dbUser.LastLogin) > updateInterval {
					s.logger.Info("Updating last login timestamp", zap.String("externalId", extId), zap.Time("previousLogin", dbUser.LastLogin))
					updateTime := time.Now()
					err = s.db.UpdateUserLastLoginWithExternalID(extId, updateTime)
					if err != nil {
						s.logger.Error("Failed to update user last login in DB", zap.String("externalId", extId), zap.Error(err))
						// Decide on retry: Keep in map or log and remove? Removing for now.
					} else {
						s.logger.Debug("Successfully updated last login", zap.String("externalId", extId))
					}
				} else {
					s.logger.Debug("Skipping last login update, too recent", zap.String("externalId", extId))
				}
			}
			// Remove processed users from the map
			for _, id := range processedIDs {
				delete(usersToUpdate, id)
			}
			s.logger.Debug("Finished processing user last login update batch")
		}
	}
}

// UpdateLastLogin queues a user for potential last login timestamp update check.
func (s *Server) UpdateLastLogin(claim *userClaim) {
	if claim != nil && claim.ExternalUserID != "" && claim.Email != "" {
		// Send minimal info needed for the loop to check DB
		select {
		case s.updateLogin <- User{ExternalId: claim.ExternalUserID, Email: claim.Email}:
			s.logger.Debug("Queued user for last login update check", zap.String("externalId", claim.ExternalUserID))
		default:
			// Channel buffer is full, update request is dropped.
			// This prevents blocking the Check request handler.
			s.logger.Warn("Update last login channel is full, dropping update request", zap.String("externalId", claim.ExternalUserID))
		}
	}
}

// Check implements the Envoy External Authorization gRPC Check method.
func (s *Server) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	// Standard unauthorized response
	unAuth := &envoyauth.CheckResponse{
		Status: &status.Status{Code: int32(rpc.UNAUTHENTICATED)},
		HttpResponse: &envoyauth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoyauth.DeniedHttpResponse{
				Status: &envoytype.HttpStatus{Code: http.StatusUnauthorized},
				Body:   http.StatusText(http.StatusUnauthorized),
			},
		},
	}

	httpRequest := req.GetAttributes().GetRequest().GetHttp()
	if httpRequest == nil {
		s.logger.Warn("Check request missing HTTP attributes")
		// Consider returning Internal error instead of Unauthorized
		return unAuth, fmt.Errorf("missing http attributes in check request")
	}
	headers := httpRequest.GetHeaders()
	path := httpRequest.GetPath()     // Use for logging
	method := httpRequest.GetMethod() // Use for logging

	// Extract Authorization header
	authHeader := headers[echo.HeaderAuthorization]
	if authHeader == "" {
		authHeader = headers[strings.ToLower(echo.HeaderAuthorization)]
	}
	if authHeader == "" {
		s.logger.Debug("Authorization header missing", zap.String("path", path), zap.String("method", method))
		return unAuth, nil // No error, just unauthorized
	}

	// Verify the token using the server's Verify method
	verifiedClaim, err := s.Verify(ctx, authHeader)
	if err != nil {
		s.logger.Info("Token verification failed", // Use Info for failed auth attempts
			zap.String("path", path),
			zap.String("method", method),
			zap.Error(err))
		// Return specific error from Verify if needed, otherwise unAuth
		// Check if error indicates expired token vs invalid signature etc.
		return unAuth, nil // No internal error, just failed verification
	}

	// Post-verification checks and DB lookup
	verifiedClaim.Email = strings.ToLower(strings.TrimSpace(verifiedClaim.Email))
	if verifiedClaim.Email == "" {
		s.logger.Warn("Verified token missing email claim", zap.String("sub", verifiedClaim.Subject), zap.String("jti", verifiedClaim.Id))
		return unAuth, nil
	}

	// Get user from local DB
	theUser, err := s.db.GetUserByEmail(verifiedClaim.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("User from token not found in local database", zap.String("email", verifiedClaim.Email))
		} else {
			s.logger.Error("Failed to get user from database during Check", zap.String("email", verifiedClaim.Email), zap.Error(err))
			// Return internal error? Or treat as unauthorized? For security, treat as UnAuth.
		}
		return unAuth, nil
	}
	if theUser == nil { // Safeguard
		s.logger.Warn("GetUserByEmail returned nil user without error during Check", zap.String("email", verifiedClaim.Email))
		return unAuth, nil
	}

	// Check if user is active locally
	if !theUser.IsActive {
		s.logger.Warn("Authentication attempt by inactive user", zap.String("email", theUser.Email), zap.Uint("dbID", theUser.ID))
		forbiddenResp := &envoyauth.CheckResponse{
			Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)},
			HttpResponse: &envoyauth.CheckResponse_DeniedResponse{
				DeniedResponse: &envoyauth.DeniedHttpResponse{
					Status: &envoytype.HttpStatus{Code: http.StatusForbidden}, // 403 Forbidden
					Body:   "User account is inactive",
				},
			},
		}
		return forbiddenResp, nil
	}

	// --- Populate claim details from reliable DB source ---
	// This ensures the headers sent downstream reflect the current DB state
	verifiedClaim.Role = theUser.Role
	verifiedClaim.ExternalUserID = theUser.ExternalId // Use DB ExternalId
	verifiedClaim.MemberSince = &theUser.CreatedAt
	verifiedClaim.UserLastLogin = &theUser.LastLogin

	// Queue async update for last login
	go s.UpdateLastLogin(verifiedClaim)

	// --- Build OK response ---
	s.logger.Info("Access granted",
		zap.String("path", path),
		zap.String("method", method),
		zap.String("email", verifiedClaim.Email),
		zap.String("role", string(verifiedClaim.Role)),
		zap.String("externalId", verifiedClaim.ExternalUserID),
	)
	// Headers to inject downstream
	headersToSend := []*envoycore.HeaderValueOption{
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserIDHeader, Value: verifiedClaim.ExternalUserID}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserRoleHeader, Value: string(verifiedClaim.Role)}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
		// Add other headers like Email if needed by downstream services
		// {Header: &envoycore.HeaderValue{Key: "X-Platform-User-Email", Value: verifiedClaim.Email}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
	}

	return &envoyauth.CheckResponse{
		Status: &status.Status{Code: int32(rpc.OK)},
		HttpResponse: &envoyauth.CheckResponse_OkResponse{
			OkResponse: &envoyauth.OkHttpResponse{
				Headers: headersToSend,
			},
		},
	}, nil
}

// Verify attempts to validate the token first via Dex OIDC, then via the platform key.
func (s *Server) Verify(ctx context.Context, authHeader string) (*userClaim, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return nil, fmt.Errorf("missing authorization token")
	}

	// 1. Try verifying with Dex OIDC Verifier
	s.logger.Debug("Attempting token verification via Dex OIDC Verifier")
	idToken, errDex := s.dexVerifier.Verify(ctx, tokenString)
	if errDex == nil {
		s.logger.Debug("Token successfully verified by Dex OIDC Verifier")
		var dexClaims DexClaims // Use the specific DexClaims struct
		if err := idToken.Claims(&dexClaims); err != nil {
			s.logger.Error("Failed to extract claims from Dex verified token", zap.Error(err))
			return nil, fmt.Errorf("failed to extract Dex token claims: %w", err)
		}

		// Map essential Dex claims to our internal userClaim structure
		claim := &userClaim{
			Email:          dexClaims.Email,
			EmailVerified:  dexClaims.EmailVerified,
			ExternalUserID: dexClaims.Subject, // Standard 'sub' from OIDC
			// Role is determined later from DB
			StandardClaims: dexClaims.StandardClaims, // Copy standard claims
		}
		return claim, nil
	}
	// Log Dex error only if it's not a typical validation error we expect to fall through
	// e.g., log network errors, config errors, but maybe Debug log standard token errors
	s.logger.Debug("Dex OIDC verification failed, attempting platform key verification", zap.Error(errDex))

	// 2. Try verifying with Platform Public Key
	if s.platformPublicKey != nil {
		s.logger.Debug("Attempting token verification via Platform Public Key")
		var platformClaims userClaim // Parse directly into our claim struct

		token, errPlatform := jwt.ParseWithClaims(tokenString, &platformClaims, func(token *jwt.Token) (interface{}, error) {
			// Check signing method matches what we use (RS256)
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method for platform token: %v", token.Header["alg"])
			}
			// Check kid header matches our platform key ID
			if kid, ok := token.Header["kid"].(string); ok {
				if kid != s.platformKeyID {
					return nil, fmt.Errorf("token 'kid' [%s] does not match expected platform key ID [%s]", kid, s.platformKeyID)
				}
				s.logger.Debug("Verifying platform token using kid", zap.String("kid", kid))
			} else {
				// If kid is mandatory for platform tokens, return error. If optional/fallback, log warning.
				s.logger.Warn("Platform token header missing 'kid', verification may be ambiguous if keys rotate.")
				// return nil, fmt.Errorf("platform token missing required 'kid' header") // Uncomment if kid is mandatory
			}
			return s.platformPublicKey, nil
		})

		// Check both error and token validity
		if errPlatform == nil && token.Valid {
			s.logger.Debug("Token successfully verified by Platform Public Key")
			// Ensure essential fields are present after parsing
			if platformClaims.Subject == "" {
				// Subject is crucial, usually maps to ExternalUserID
				return nil, fmt.Errorf("platform token missing 'sub' claim")
			}
			platformClaims.ExternalUserID = platformClaims.Subject // Standard mapping

			// Role might be embedded as a custom claim, parse if needed, otherwise set later from DB
			// if roleClaim, ok := platformClaims.PrivateClaims["role"].(string); ok {
			//     platformClaims.Role = api.Role(roleClaim)
			// }

			return &platformClaims, nil
		}
		// Log platform verification error, including validation errors like expiry
		s.logger.Debug("Platform key verification failed", zap.String("expected_kid", s.platformKeyID), zap.Error(errPlatform))

	} else {
		s.logger.Debug("Platform public key not configured, skipping platform verification")
	}

	// If both failed, return a consolidated error message
	// Prioritize returning Dex error if it exists, as it was the first attempt
	if errDex != nil {
		// Check if Dex error is a standard validation error vs. config/network error
		// Return a generic message for standard validation errors
		return nil, fmt.Errorf("token verification failed") // Don't leak Dex error details potentially
	}
	return nil, fmt.Errorf("token verification failed") // Generic failure if Dex didn't error and platform failed/skipped
}

// (newDexOidcVerifier and newDexClient functions are defined in cmd.go)
