package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	// Keep grpc import if newDexClient is defined here
	"github.com/coreos/go-oidc/v3/oidc"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4" // Keep if needed elsewhere
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/authcache" // Import authcache
	"github.com/opengovern/opensecurity/services/auth/db"
	"github.com/opengovern/opensecurity/services/auth/utils"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
	// "google.golang.org/grpc/credentials" // Keep if newServerCredentials is defined here
)

// User struct used within the UpdateLastLoginLoop channel.
type User struct {
	ID         string // DB ID? Or ExternalID? Clarify usage. Assuming ExternalId based on loop logic.
	Email      string
	ExternalId string // Explicitly adding ExternalId
	Role       api.Role
	LastLogin  time.Time
	CreatedAt  time.Time
}

// Server holds the core components of the authentication service.
type Server struct {
	host                string
	platformPublicKey   *rsa.PublicKey
	dexVerifier         *oidc.IDTokenVerifier
	dexClient           dexApi.DexClient
	logger              *zap.Logger
	db                  db.Database
	authCache           *authcache.AuthCacheService // Injected cache service instance
	updateLoginUserList []User                      // State for UpdateLastLoginLoop - consider alternatives if scaling
	updateLogin         chan User                   // Channel for UpdateLastLoginLoop
}

// DexClaims represents the expected claims structure within a Dex ID token.
type DexClaims struct {
	Email           string                 `json:"email"`
	EmailVerified   bool                   `json:"email_verified"`
	Groups          []string               `json:"groups"`           // Optional groups claim
	Name            string                 `json:"name"`             // Optional name claim
	FederatedClaims map[string]interface{} `json:"federated_claims"` // Optional federated claims
	jwt.StandardClaims
}

// userClaim represents the unified claims extracted from either a Dex token or a platform API key token.
type userClaim struct {
	Role           api.Role   // Role derived from DB/Cache or platform token
	Email          string     // User's email
	MemberSince    *time.Time // User's creation time (from DB/Cache)
	UserLastLogin  *time.Time // User's last login time (from DB/Cache)
	ExternalUserID string     `json:"sub"` // Subject claim (External ID)
	EmailVerified  bool       // From Dex token if available
}

// Valid implements jwt.Claims interface (basic validation).
func (u userClaim) Valid() error {
	if u.ExternalUserID == "" || u.Email == "" {
		return errors.New("user claim missing external ID or email")
	}
	return nil
}

// UpdateLastLoginLoop periodically updates the last login time for users in the database.
func (s *Server) UpdateLastLoginLoop() {
	s.logger.Info("Starting UpdateLastLoginLoop...")
	ticker := time.NewTicker(30 * time.Second) // Check more frequently?
	defer ticker.Stop()

	for {
		select {
		case userUpdate := <-s.updateLogin:
			exists := false
			for _, u := range s.updateLoginUserList {
				if u.ExternalId == userUpdate.ExternalId {
					exists = true
					break
				}
			}
			if !exists {
				s.updateLoginUserList = append(s.updateLoginUserList, userUpdate)
				s.logger.Debug("Added user to update queue", zap.String("externalId", userUpdate.ExternalId))
			}

		case <-ticker.C:
			if len(s.updateLoginUserList) == 0 {
				continue
			}

			s.logger.Debug("Processing last login update queue", zap.Int("queueSize", len(s.updateLoginUserList)))
			processedList := s.updateLoginUserList
			s.updateLoginUserList = nil // Clear the list for next cycle

			for _, user := range processedList {
				if user.ExternalId == "" || user.Email == "" {
					s.logger.Warn("Skipping update for user with missing ExternalId or Email", zap.Any("user", user))
					continue
				}

				dbUser, err := utils.GetUserByEmail(user.Email, s.db)
				if err != nil {
					s.logger.Error("Failed to get user for last login check", zap.String("externalId", user.ExternalId), zap.Error(err))
					continue
				}
				if dbUser == nil {
					s.logger.Warn("User from update queue not found in DB", zap.String("externalId", user.ExternalId))
					continue
				}

				if time.Since(dbUser.LastLogin) > (15 * time.Minute) { // Use time.Since for clarity
					now := time.Now()
					s.logger.Info("Updating last login time in DB", zap.String("externalId", user.ExternalId), zap.Time("newLoginTime", now))
					err = utils.UpdateUserLastLogin(user.ExternalId, now, s.db)
					if err != nil {
						s.logger.Error("Failed to update last login time in DB", zap.String("externalId", user.ExternalId), zap.Error(err))
					}
				} else {
					s.logger.Debug("Skipping last login update, not enough time passed", zap.String("externalId", user.ExternalId), zap.Time("lastLoginDB", dbUser.LastLogin))
				}
			}
			// Add a way to gracefully exit the loop, e.g., listening on context cancellation
			// case <-ctx.Done():
			//  s.logger.Info("Stopping UpdateLastLoginLoop due to context cancellation.")
			//  return
		}
	}
}

// UpdateLastLogin sends user details to the UpdateLastLoginLoop if conditions are met.
func (s *Server) UpdateLastLogin(claim *userClaim) {
	if claim == nil || claim.ExternalUserID == "" {
		s.logger.Warn("UpdateLastLogin called with nil or invalid claim")
		return
	}

	s.logger.Debug("Queueing user for potential last login update", zap.String("externalId", claim.ExternalUserID))
	select {
	case s.updateLogin <- User{ExternalId: claim.ExternalUserID, Email: claim.Email}:
		// Successfully queued
	default:
		s.logger.Warn("UpdateLastLogin channel buffer full, dropping update request", zap.String("externalId", claim.ExternalUserID))
	}
}

// Check performs the authorization check for Envoy.
// It verifies the token, checks the cache, falls back to the database, and populates the cache.
func (s *Server) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	// Standard Unauthorized response structure
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
	headers := httpRequest.GetHeaders()

	// Extract Authorization header
	authHeader := headers[echo.HeaderAuthorization]
	if authHeader == "" {
		authHeader = headers[strings.ToLower(echo.HeaderAuthorization)]
	}

	// 1. Verify the token (JWT signature, expiry, platform key etc.)
	verifiedClaim, err := s.Verify(ctx, authHeader)
	if err != nil {
		s.logger.Warn("Token verification failed",
			zap.String("reqId", httpRequest.Id),
			zap.String("path", httpRequest.Path),
			zap.Error(err))
		return unAuth, nil // Return UNAUTHENTICATED
	}

	// Ensure email is present and clean after verification
	verifiedClaim.Email = strings.ToLower(strings.TrimSpace(verifiedClaim.Email))
	if verifiedClaim.Email == "" {
		s.logger.Warn("Token verified but email claim is missing or empty", zap.String("externalId", verifiedClaim.ExternalUserID))
		return unAuth, nil // Return UNAUTHENTICATED
	}
	logFields := []zap.Field{
		zap.String("email", verifiedClaim.Email),
		zap.String("externalId", verifiedClaim.ExternalUserID),
		zap.String("path", httpRequest.Path),
		zap.String("method", httpRequest.Method),
	}
	s.logger.Debug("Token verified successfully", logFields...)

	// --- Authorization Check with Cache ---
	// var userInfo *utils.User // No longer needed here, use specific types
	var userRole api.Role     // Final resolved role
	var userExternalId string // Final resolved external ID

	// 2. Check Cache
	cachedInfo, cacheErr := s.authCache.GetUser(ctx, verifiedClaim.Email)

	if cacheErr == nil {
		// Cache Hit!
		s.logger.Debug("Auth cache hit", logFields...)
		// Check if the cached user is active
		if !cachedInfo.IsActive {
			s.logger.Warn("Access denied: User found in cache but is inactive", logFields...)
			return unAuth, nil // Return UNAUTHENTICATED for inactive user
		}
		// Use cached data
		userRole = api.Role(cachedInfo.Role) // Convert string back to api.Role
		userExternalId = cachedInfo.ExternalID
		// Queue for last login update based on verified claim (which has necessary IDs)
		go s.UpdateLastLogin(verifiedClaim)

	} else if errors.Is(cacheErr, authcache.ErrUserInfoNotFound) {
		// Cache Miss - Proceed to DB lookup
		s.logger.Debug("Auth cache miss, checking database", logFields...)

		// 3. DB Lookup (Fallback)
		dbUser, dbErr := utils.GetUserByEmail(verifiedClaim.Email, s.db)
		if dbErr != nil {
			// Handle specific errors from GetUserByEmail (already logs internally)
			// e.g., "user disabled", "user not found"
			s.logger.Warn("Access denied: User lookup failed or user invalid", append(logFields, zap.Error(dbErr))...)
			return unAuth, nil // Return UNAUTHENTICATED
		}
		if dbUser == nil { // Defensive check
			s.logger.Error("GetUserByEmail returned nil user without error", logFields...)
			return unAuth, nil
		}

		// DB Hit! User is valid and active (checked within GetUserByEmail)
		userRole = api.Role(dbUser.Role)
		userExternalId = dbUser.ExternalId

		// 4. Populate Cache
		cacheUserInfo := &authcache.CachedUserInfo{
			ID:         dbUser.ID,
			Role:       dbUser.Role, // Store as string
			ExternalID: dbUser.ExternalId,
			IsActive:   dbUser.IsActive,
		}
		if cacheSetErr := s.authCache.AddUserToCache(ctx, verifiedClaim.Email, cacheUserInfo); cacheSetErr != nil {
			// Log cache population error but don't fail the request
			s.logger.Error("Failed to populate auth cache after DB lookup", append(logFields, zap.Error(cacheSetErr))...)
		}

		// Queue for last login update based on verified claim
		go s.UpdateLastLogin(verifiedClaim)

	} else {
		// Other Cache Error
		s.logger.Error("Auth cache error during GET, falling back to DB", append(logFields, zap.Error(cacheErr))...)

		// Fallback to DB as if it was a cache miss
		dbUser, dbErr := utils.GetUserByEmail(verifiedClaim.Email, s.db)
		if dbErr != nil {
			s.logger.Warn("Access denied: User lookup failed or user invalid (after cache error)", append(logFields, zap.Error(dbErr))...)
			return unAuth, nil
		}
		if dbUser == nil {
			s.logger.Error("GetUserByEmail returned nil user without error (after cache error)", logFields...)
			return unAuth, nil
		}
		// Use DB data, but don't attempt to cache again
		userRole = api.Role(dbUser.Role)
		userExternalId = dbUser.ExternalId
		go s.UpdateLastLogin(verifiedClaim)
	}

	// --- Authorization Decision ---
	// At this point, we have a valid, active user (either from cache or DB)
	// Add more complex authorization logic here if needed (e.g., checking roles against path/method)
	s.logger.Info("Authorization check successful", logFields...)

	// 5. Construct OK Response
	// Ensure userExternalId and userRole are populated correctly from either cache or DB path
	if userExternalId == "" || userRole == "" {
		s.logger.Error("Internal error: ExternalID or Role missing after successful check", logFields...)
		return unAuth, fmt.Errorf("internal check error: missing user details") // Should not happen
	}

	return &envoyauth.CheckResponse{
		Status: &status.Status{Code: int32(rpc.OK)},
		HttpResponse: &envoyauth.CheckResponse_OkResponse{
			OkResponse: &envoyauth.OkHttpResponse{
				Headers: []*envoycore.HeaderValueOption{
					{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserIDHeader, Value: userExternalId}},
					{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserRoleHeader, Value: string(userRole)}},
					// XPlatformUserConnectionsScope seems redundant if it's the same as ExternalId? Adjust if needed.
					{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserConnectionsScope, Value: userExternalId}},
				},
			},
		},
	}, nil
}

// Verify checks the validity of a Bearer token using either Dex OIDC or the platform's RSA key.
// It returns a userClaim containing essential info if verification succeeds.
func (s *Server) Verify(ctx context.Context, authToken string) (*userClaim, error) {
	if !strings.HasPrefix(authToken, "Bearer ") {
		return nil, errors.New("invalid authorization token format")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authToken, "Bearer "))
	if token == "" {
		return nil, errors.New("missing authorization token")
	}

	// --- Attempt 1: Verify using Dex OIDC Verifier ---
	s.logger.Debug("Attempting token verification via Dex OIDC Verifier...")
	idToken, err := s.dexVerifier.Verify(ctx, token)
	if err == nil {
		// Dex verification successful
		var claims DexClaims // Use the specific Dex claims struct
		if claimErr := idToken.Claims(&claims); claimErr != nil {
			s.logger.Error("Dex OIDC token verified, but failed to extract claims", zap.Error(claimErr))
			return nil, fmt.Errorf("failed to extract dex claims: %w", claimErr)
		}
		s.logger.Debug("Dex OIDC verification successful", zap.String("subject", claims.Subject), zap.String("email", claims.Email))
		// Construct userClaim from Dex claims
		// Role is NOT set here, will be determined later in Check()
		return &userClaim{
			Email:          claims.Email,
			EmailVerified:  claims.EmailVerified,
			ExternalUserID: claims.Subject, // Standard OIDC subject claim maps to ExternalUserID
		}, nil
	} else {
		// Log Dex verification error, but don't return yet, try platform key
		s.logger.Warn("Dex OIDC verification failed, proceeding to check platform key", zap.Error(err))
	}

	// --- Attempt 2: Verify using Platform RSA Public Key ---
	if s.platformPublicKey == nil {
		// Platform key verification is disabled or key not loaded
		s.logger.Debug("Platform public key not available, cannot verify platform token.")
		// Return the original Dex verification error as the final reason
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	s.logger.Debug("Attempting token verification via Platform Key...")
	var platformClaims userClaim // Use the userClaim struct directly for platform tokens
	parsedToken, parseErr := jwt.ParseWithClaims(token, &platformClaims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.platformPublicKey, nil
	})

	if parseErr == nil && parsedToken.Valid {
		// Platform key verification successful
		s.logger.Debug("Platform key verification successful", zap.String("subject", platformClaims.ExternalUserID), zap.String("email", platformClaims.Email), zap.String("role", string(platformClaims.Role)))
		// Role *is* expected in platform tokens, return the full claim
		return &platformClaims, nil
	} else {
		// Log platform key verification error
		s.logger.Warn("Platform key verification failed", zap.Error(parseErr))
		// Return the *original* Dex verification error as the primary failure reason,
		// potentially adding context about the platform key failure.
		return nil, fmt.Errorf("token verification failed (Dex error: %v / Platform key error: %v)", err, parseErr)
	}
}
