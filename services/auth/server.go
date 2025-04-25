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
	host              string
	platformPublicKey *rsa.PublicKey
	platformKeyID     string // Stores the calculated Key ID (JWK Thumbprint)
	dexVerifier       *oidc.IDTokenVerifier
	dexClient         dexApi.DexClient
	logger            *zap.Logger
	// --- MODIFIED ---
	db db.DatabaseInterface // Use the database interface
	// --- END MODIFIED ---
	updateLogin chan User // Channel for queuing login updates
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

// Valid implements the jwt.Claims interface for the userClaim struct.
// This method is called by jwt.ParseWithClaims after signature verification
// and after default time-based claims (exp, nbf, iat) have been checked.
// Use this method to enforce application-specific payload validation rules.
func (u *userClaim) Valid() error {
	// 1. Standard time validations (exp, nbf, iat) are handled by the default
	//    jwt.ParseWithClaims logic before this method is called. We don't need
	//    to repeat them here unless we want extremely custom error reporting.
	//    Example using v5 helpers if needed (usually not necessary):
	//    now := time.Now()
	//    if !u.VerifyExpiresAt(now, true) { // true = required
	// 	   return jwt.ErrTokenExpired
	//    }
	//    if !u.VerifyIssuedAt(now, true) { // true = required
	// 	   return jwt.ErrTokenUsedBeforeIssued
	//    }
	//    if !u.VerifyNotBefore(now, true) { // true = required
	// 	   return jwt.ErrTokenNotValidYet
	//    }

	// 2. Validate presence of essential application-specific claims.
	//    The 'sub' (Subject) claim is critical for identifying the user.
	//    It's mapped to the embedded StandardClaims.Subject field.
	if u.Subject == "" {
		// Use standard validation error types from the jwt package
		return jwt.NewValidationError("token missing required 'sub' (subject) claim", jwt.ValidationErrorClaimsInvalid)
	}

	// 3. Validate Issuer if expected for platform tokens.
	//    OIDC tokens issuer is validated by the oidc library.
	//    Platform tokens (like API keys) should have a specific issuer.
	//    We check this during parsing in Verify() rather than here, but could add here too.
	// if u.Issuer != "platform-auth-service" { // Example check
	//  return jwt.NewValidationError(fmt.Sprintf("invalid 'iss' (issuer) claim: %s", u.Issuer), jwt.ValidationErrorClaimsInvalid)
	// }

	// 4. Validate Audience if expected.
	//    Audience checks are often done during ParseWithClaims using jwt.WithAudience option,
	//    but can be done here if complex logic is needed.
	// if !u.VerifyAudience("platform-api", true) { // Example check
	//     return jwt.NewValidationError(fmt.Sprintf("invalid 'aud' (audience) claim: %v", u.Audience), jwt.ValidationErrorClaimsInvalid)
	// }

	// Add any other custom checks relevant to your specific claims payload here.
	// For example, if you added a custom "scope" claim:
	// if u.Scope == "" {
	//     return jwt.NewValidationError("token missing required 'scope' claim", jwt.ValidationErrorClaimsInvalid)
	// }

	// If all application-specific checks pass, return nil.
	return nil
}

// UpdateLastLoginLoop processes queued users to update their last login timestamp.
func (s *Server) UpdateLastLoginLoop() {
	ticker := time.NewTicker(1 * time.Minute) // Check less frequently
	defer ticker.Stop()
	usersToUpdate := make(map[string]User) // Use map for efficient deduplication

	// Background context for DB calls originating from this loop
	loopCtx := context.Background()

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

			processedIDs := []string{} // Keep track of processed IDs in this batch
			for extId := range usersToUpdate {
				processedIDs = append(processedIDs, extId) // Mark for processing

				// --- MODIFIED: Pass context to DB calls ---
				dbUser, err := s.db.GetUserByExternalID(loopCtx, extId)
				if err != nil {
					// Check explicitly for not found (assuming db returns nil,nil or error)
					if errors.Is(err, gorm.ErrRecordNotFound) || dbUser == nil && err == nil {
						s.logger.Warn("User no longer exists in DB, cannot update last login", zap.String("externalId", extId))
					} else {
						s.logger.Error("Failed to get user from DB for last login update", zap.String("externalId", extId), zap.Error(err))
					}
					continue // Skip this user for now if error or not found
				}
				// Note: dbUser check for nil after error check might be redundant if GetUserByExternalID contract is clear

				updateInterval := 15 * time.Minute
				if time.Since(dbUser.LastLogin) > updateInterval {
					s.logger.Info("Updating last login timestamp", zap.String("externalId", extId), zap.Time("previousLogin", dbUser.LastLogin))
					updateTime := time.Now()
					// --- MODIFIED: Pass context to DB calls ---
					err = s.db.UpdateUserLastLoginWithExternalID(loopCtx, extId, updateTime)
					if err != nil {
						s.logger.Error("Failed to update user last login in DB", zap.String("externalId", extId), zap.Error(err))
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
			// Add case for context cancellation if the loop needs to stop gracefully
			// case <- loopCtx.Done():
			//     s.logger.Info("UpdateLastLoginLoop stopping due to context cancellation.")
			//     return
		}
	}
}

// UpdateLastLogin queues a user for potential last login timestamp update check.
func (s *Server) UpdateLastLogin(claim *userClaim) {
	if claim != nil && claim.ExternalUserID != "" && claim.Email != "" {
		select {
		case s.updateLogin <- User{ExternalId: claim.ExternalUserID, Email: claim.Email}:
			s.logger.Debug("Queued user for last login update check", zap.String("externalId", claim.ExternalUserID))
		default:
			s.logger.Warn("Update last login channel is full, dropping update request", zap.String("externalId", claim.ExternalUserID))
		}
	}
}

// Check implements the Envoy External Authorization gRPC Check method.
func (s *Server) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
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
		return unAuth, fmt.Errorf("missing http attributes in check request")
	}
	headers := httpRequest.GetHeaders()
	path := httpRequest.GetPath()
	method := httpRequest.GetMethod()

	authHeader := headers[echo.HeaderAuthorization]
	if authHeader == "" {
		authHeader = headers[strings.ToLower(echo.HeaderAuthorization)]
	}
	if authHeader == "" {
		s.logger.Debug("Authorization header missing", zap.String("path", path), zap.String("method", method))
		return unAuth, nil
	}

	verifiedClaim, err := s.Verify(ctx, authHeader) // Pass context
	if err != nil {
		s.logger.Info("Token verification failed", zap.String("path", path), zap.String("method", method), zap.Error(err))
		return unAuth, nil
	}

	verifiedClaim.Email = strings.ToLower(strings.TrimSpace(verifiedClaim.Email))
	if verifiedClaim.Email == "" {
		s.logger.Warn("Verified token missing email claim", zap.String("sub", verifiedClaim.Subject), zap.String("jti", verifiedClaim.Id))
		return unAuth, nil
	}

	// Get user from local DB using the interface and passing context
	theUser, err := s.db.GetUserByEmail(ctx, verifiedClaim.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("User from token not found in local database", zap.String("email", verifiedClaim.Email))
		} else {
			s.logger.Error("Failed to get user from database during Check", zap.String("email", verifiedClaim.Email), zap.Error(err))
		}
		return unAuth, nil
	}
	if theUser == nil {
		s.logger.Warn("GetUserByEmail returned nil user without error during Check", zap.String("email", verifiedClaim.Email))
		return unAuth, nil
	}

	if !theUser.IsActive {
		s.logger.Warn("Authentication attempt by inactive user", zap.String("email", theUser.Email), zap.Uint("dbID", theUser.ID))
		forbiddenResp := &envoyauth.CheckResponse{Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)}, HttpResponse: &envoyauth.CheckResponse_DeniedResponse{DeniedResponse: &envoyauth.DeniedHttpResponse{Status: &envoytype.HttpStatus{Code: http.StatusForbidden}, Body: "User account is inactive"}}}
		return forbiddenResp, nil
	}

	// Populate claim details from reliable DB source
	verifiedClaim.Role = theUser.Role
	verifiedClaim.ExternalUserID = theUser.ExternalId
	verifiedClaim.MemberSince = &theUser.CreatedAt
	verifiedClaim.UserLastLogin = &theUser.LastLogin

	go s.UpdateLastLogin(verifiedClaim) // Queue async update

	s.logger.Info("Access granted", zap.String("path", path), zap.String("method", method), zap.String("email", verifiedClaim.Email), zap.String("role", string(verifiedClaim.Role)), zap.String("externalId", verifiedClaim.ExternalUserID))
	headersToSend := []*envoycore.HeaderValueOption{
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserIDHeader, Value: verifiedClaim.ExternalUserID}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserRoleHeader, Value: string(verifiedClaim.Role)}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
	}
	return &envoyauth.CheckResponse{Status: &status.Status{Code: int32(rpc.OK)}, HttpResponse: &envoyauth.CheckResponse_OkResponse{OkResponse: &envoyauth.OkHttpResponse{Headers: headersToSend}}}, nil
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
	idToken, errDex := s.dexVerifier.Verify(ctx, tokenString) // Pass context
	if errDex == nil {
		s.logger.Debug("Token successfully verified by Dex OIDC Verifier")
		var dexClaims DexClaims
		if err := idToken.Claims(&dexClaims); err != nil {
			s.logger.Error("Failed to extract claims from Dex verified token", zap.Error(err))
			return nil, fmt.Errorf("failed to extract Dex token claims: %w", err)
		}
		claim := &userClaim{Email: dexClaims.Email, EmailVerified: dexClaims.EmailVerified, ExternalUserID: dexClaims.Subject, StandardClaims: dexClaims.StandardClaims}
		return claim, nil
	}
	s.logger.Debug("Dex OIDC verification failed, attempting platform key verification", zap.Error(errDex))

	// 2. Try verifying with Platform Public Key
	if s.platformPublicKey != nil {
		s.logger.Debug("Attempting token verification via Platform Public Key")
		var platformClaims userClaim
		token, errPlatform := jwt.ParseWithClaims(tokenString, &platformClaims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method for platform token: %v", token.Header["alg"])
			}
			if kid, ok := token.Header["kid"].(string); ok {
				if kid != s.platformKeyID {
					return nil, fmt.Errorf("token 'kid' [%s] does not match expected platform key ID [%s]", kid, s.platformKeyID)
				}
				s.logger.Debug("Verifying platform token using kid", zap.String("kid", kid))
			} else {
				s.logger.Warn("Platform token header missing 'kid', verification may be ambiguous if keys rotate.")
			}
			return s.platformPublicKey, nil
		})

		if errPlatform == nil && token.Valid {
			s.logger.Debug("Token successfully verified by Platform Public Key")
			if platformClaims.Subject == "" {
				return nil, fmt.Errorf("platform token missing 'sub' claim")
			}
			platformClaims.ExternalUserID = platformClaims.Subject
			return &platformClaims, nil
		}
		s.logger.Debug("Platform key verification failed", zap.String("expected_kid", s.platformKeyID), zap.Error(errPlatform))
	} else {
		s.logger.Debug("Platform public key not configured, skipping platform verification")
	}

	return nil, fmt.Errorf("token verification failed") // Generic failure
}

// (Helper functions like newDexOidcVerifier, newDexClient are defined in cmd.go)
