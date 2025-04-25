// Package auth provides the core authentication and authorization logic,
// including OIDC integration with Dex, platform token issuance/verification,
// API key management, user management, and an Envoy external authorization check service.
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

// User struct is an internal representation used primarily for passing
// minimal user identifiers between the Check/Verify methods and the
// asynchronous UpdateLastLoginLoop.
type User struct {
	ID         string // Typically ExternalId, though name conflicts with ExternalId field below. Review needed.
	Email      string
	ExternalId string   // Canonical user identifier (e.g., connector|id)
	Role       api.Role // User's application role
	LastLogin  time.Time
	CreatedAt  time.Time
}

// Server holds dependencies and state for the core authentication server logic,
// primarily handling gRPC Check requests from Envoy and token verification.
type Server struct {
	host              string                // Hostname information (usage unclear in provided snippet)
	platformPublicKey *rsa.PublicKey        // Public key for verifying platform-issued JWTs/API Keys.
	platformKeyID     string                // Key ID ('kid') associated with the platform keys (JWK thumbprint).
	dexVerifier       *oidc.IDTokenVerifier // Verifier for OIDC ID Tokens issued by Dex.
	dexClient         dexApi.DexClient      // gRPC client for interacting with the Dex API (passwords, connectors).
	logger            *zap.Logger           // Structured logger.
	db                db.DatabaseInterface  // Interface for database operations.
	updateLogin       chan User             // Channel for queuing asynchronous last login updates.
}

// userClaim represents the custom claims structure used for platform-issued JWTs
// (both enriched user tokens and API Keys). It embeds standard JWT claims
// and includes fields populated after verification based on DB lookups.
type userClaim struct {
	Role           api.Role   // User's role within the application.
	Email          string     // User's email address.
	MemberSince    *time.Time // Timestamp when the user was created (populated post-verification).
	UserLastLogin  *time.Time // Timestamp of last verified login activity (populated post-verification).
	ExternalUserID string     // Canonical user identifier (maps to 'sub' claim).
	EmailVerified  bool       // Email verification status (usually from Dex).

	// Embed standard claims for JWTs that use them (e.g., platform tokens)
	jwt.StandardClaims
}

/*
Valid implements the jwt.Claims interface. It provides *custom* validation checks
on the token's payload *after* the standard validation steps within
jwt.ParseWithClaims have already passed.

The validation sequence performed by jwt.ParseWithClaims is roughly:
 1. Parse header and payload.
 2. Verify cryptographic signature using the key provided by the Keyfunc.
 3. Perform standard time-based claim validation (checking 'exp', 'nbf', 'iat'
    against the current time, considering any configured leeway).
 4. **If and only if** the signature and standard time claims are valid,
    call this `Valid()` method on the claims.
 5. If `Valid()` returns an error, `jwt.ParseWithClaims` returns that error.
 6. If `Valid()` returns `nil`, `jwt.ParseWithClaims` returns the parsed token
    (marked as valid) and a `nil` error.

Crucially, if the token signature is invalid OR if standard time validation fails
(e.g., the token is expired because 'exp' is in the past), jwt.ParseWithClaims
returns an error immediately, and **this `Valid()` method will not be called.**

Therefore, this method only needs to implement application-specific checks
that go beyond the standard ones performed by the library. Here, we ensure
the essential 'sub' (Subject) claim is present, as required by our application
logic to identify the user. We don't need to re-check expiration.
*/
func (u *userClaim) Valid() error {
	// Standard time checks (exp, nbf, iat) are implicitly handled by the
	// jwt.ParseWithClaims function *before* this method is invoked.

	// Application-specific check: Subject claim ('sub') MUST be present.
	// The 'sub' claim maps to the embedded StandardClaims.Subject field.
	if u.Subject == "" {
		// Return a standard validation error indicating which claim failed.
		return jwt.NewValidationError("token missing required 'sub' (subject) claim", jwt.ValidationErrorClaimsInvalid)
	}

	// Add other application-specific payload checks here if needed.

	// If all custom checks implemented here pass, return nil.
	return nil
}

// UpdateLastLoginLoop runs as a background goroutine, periodically processing
// users whose activity has been detected by the Check method. It updates the
// last_login timestamp in the database for users if sufficient time has passed
// since their last recorded login, using a background context for DB calls.
// It uses a map to efficiently deduplicate update requests within a time window.
func (s *Server) UpdateLastLoginLoop() {
	ticker := time.NewTicker(1 * time.Minute) // Frequency of update checks
	defer ticker.Stop()
	usersToUpdate := make(map[string]User) // Stores users needing potential update

	// Background context for DB operations initiated by this loop
	loopCtx := context.Background()

	for {
		select {
		case user := <-s.updateLogin: // Receive user info from Check handler
			if user.ExternalId != "" {
				usersToUpdate[user.ExternalId] = user
				s.logger.Debug("User added/updated in last login queue", zap.String("externalId", user.ExternalId))
			} else {
				s.logger.Warn("Received user for last login update with empty ExternalId", zap.String("email", user.Email))
			}
		case <-ticker.C: // Timer fires to process the queue
			if len(usersToUpdate) == 0 {
				continue // Skip if queue is empty
			}
			s.logger.Debug("Processing user last login update batch", zap.Int("queue_size", len(usersToUpdate)))

			processedIDs := make([]string, 0, len(usersToUpdate))
			for extId := range usersToUpdate {
				processedIDs = append(processedIDs, extId) // Track processed users for removal

				dbUser, err := s.db.GetUserByExternalID(loopCtx, extId)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) || (dbUser == nil && err == nil) {
						s.logger.Warn("User no longer exists in DB, cannot update last login", zap.String("externalId", extId))
					} else {
						s.logger.Error("Failed to get user from DB for last login update", zap.String("externalId", extId), zap.Error(err))
					}
					continue // Skip processing this user on error/not found
				}

				updateInterval := 15 * time.Minute // Interval after which to update login time
				if time.Since(dbUser.LastLogin) > updateInterval {
					s.logger.Info("Updating last login timestamp", zap.String("externalId", extId), zap.Time("previousLogin", dbUser.LastLogin))
					updateTime := time.Now()
					err = s.db.UpdateUserLastLoginWithExternalID(loopCtx, extId, updateTime)
					if err != nil {
						s.logger.Error("Failed to update user last login in DB", zap.String("externalId", extId), zap.Error(err))
						// Consider retry logic or permanent failure handling here
					} else {
						s.logger.Debug("Successfully updated last login", zap.String("externalId", extId))
					}
				} else {
					s.logger.Debug("Skipping last login update, too recent", zap.String("externalId", extId))
				}
			}
			// Remove processed users from the map for the next cycle
			for _, id := range processedIDs {
				delete(usersToUpdate, id)
			}
			s.logger.Debug("Finished processing user last login update batch")
			// Consider adding context cancellation check for graceful shutdown of the loop
			// case <- loopCtx.Done(): return
		}
	}
}

// UpdateLastLogin queues a user identified by their claims for a potential
// asynchronous update of their last_login timestamp via the UpdateLastLoginLoop.
// Uses a non-blocking send to avoid delaying the calling goroutine (e.g., Check handler).
func (s *Server) UpdateLastLogin(claim *userClaim) {
	if claim != nil && claim.ExternalUserID != "" && claim.Email != "" {
		// Send minimal required info to the update channel
		userUpdate := User{ExternalId: claim.ExternalUserID, Email: claim.Email}
		select {
		case s.updateLogin <- userUpdate:
			s.logger.Debug("Queued user for last login update check", zap.String("externalId", claim.ExternalUserID))
		default:
			// Log if the channel buffer is full and the update is dropped
			s.logger.Warn("Update last login channel is full, dropping update request", zap.String("externalId", claim.ExternalUserID))
		}
	}
}

// Check implements the Envoy External Authorization gRPC service `Check` method.
// It verifies the Authorization header token, checks the user's status and role
// against the local database, and returns an authorization decision (OK or Denied)
// with appropriate headers for downstream consumption.
func (s *Server) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	// Standard unauthorized response for failures
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

	// Extract Authorization Bearer token
	authHeader := headers[echo.HeaderAuthorization]
	if authHeader == "" {
		authHeader = headers[strings.ToLower(echo.HeaderAuthorization)]
	}
	if authHeader == "" {
		s.logger.Debug("Authorization header missing", zap.String("path", path), zap.String("method", method))
		return unAuth, nil
	}

	// Verify the token (handles both Dex OIDC and Platform JWTs/API Keys)
	verifiedClaim, err := s.Verify(ctx, authHeader) // Pass request context
	if err != nil {
		s.logger.Info("Token verification failed", zap.String("path", path), zap.String("method", method), zap.Error(err))
		return unAuth, nil
	}

	// Basic claim validation
	verifiedClaim.Email = strings.ToLower(strings.TrimSpace(verifiedClaim.Email))
	if verifiedClaim.Email == "" {
		s.logger.Warn("Verified token missing email claim", zap.String("sub", verifiedClaim.Subject), zap.String("jti", verifiedClaim.Id))
		return unAuth, nil
	}

	// Get authoritative user state from local database using email from token
	theUser, err := s.db.GetUserByEmail(ctx, verifiedClaim.Email) // Pass request context
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("User from token not found in local database", zap.String("email", verifiedClaim.Email))
		} else {
			s.logger.Error("Failed to get user from database during Check", zap.String("email", verifiedClaim.Email), zap.Error(err))
		}
		return unAuth, nil // Deny if user not found or DB error
	}
	if theUser == nil {
		s.logger.Warn("GetUserByEmail returned nil user without error during Check", zap.String("email", verifiedClaim.Email))
		return unAuth, nil
	}

	// Authorization Check: Ensure user is active in local system
	if !theUser.IsActive {
		s.logger.Warn("Authentication attempt by inactive user", zap.String("email", theUser.Email), zap.Uint("dbID", theUser.ID))
		forbiddenResp := &envoyauth.CheckResponse{Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)}, HttpResponse: &envoyauth.CheckResponse_DeniedResponse{DeniedResponse: &envoyauth.DeniedHttpResponse{Status: &envoytype.HttpStatus{Code: http.StatusForbidden}, Body: "User account is inactive"}}}
		return forbiddenResp, nil
	}

	// Populate claim details from reliable DB source for header injection
	verifiedClaim.Role = theUser.Role
	verifiedClaim.ExternalUserID = theUser.ExternalId // Use the canonical ID from DB
	// These aren't typically needed for Check decision but are populated if Verify was called
	// verifiedClaim.MemberSince = &theUser.CreatedAt
	// verifiedClaim.UserLastLogin = &theUser.LastLogin

	// Queue async update for last login timestamp
	go s.UpdateLastLogin(verifiedClaim)

	// Build successful authorization response
	s.logger.Info("Access granted", zap.String("path", path), zap.String("method", method), zap.String("email", verifiedClaim.Email), zap.String("role", string(verifiedClaim.Role)), zap.String("externalId", verifiedClaim.ExternalUserID))
	headersToSend := []*envoycore.HeaderValueOption{
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserIDHeader, Value: verifiedClaim.ExternalUserID}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
		{Header: &envoycore.HeaderValue{Key: httpserver.XPlatformUserRoleHeader, Value: string(verifiedClaim.Role)}, AppendAction: envoycore.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD},
		// Consider adding X-Platform-User-Email if needed downstream
	}
	return &envoyauth.CheckResponse{Status: &status.Status{Code: int32(rpc.OK)}, HttpResponse: &envoyauth.CheckResponse_OkResponse{OkResponse: &envoyauth.OkHttpResponse{Headers: headersToSend}}}, nil
}

// Verify attempts to validate the provided 'Authorization: Bearer <token>' header string.
// It first tries to validate it as an OIDC token using the Dex verifier. If that fails,
// it attempts to validate it as a platform-issued JWT (user token or API key) using
// the platform's public key and Key ID ('kid').
// Returns the validated claims in a userClaim struct or an error if validation fails.
func (s *Server) Verify(ctx context.Context, authHeader string) (*userClaim, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return nil, fmt.Errorf("missing authorization token")
	}

	// --- Attempt 1: Verify as Dex OIDC Token ---
	s.logger.Debug("Attempting token verification via Dex OIDC Verifier")
	idToken, errDex := s.dexVerifier.Verify(ctx, tokenString) // Pass context
	if errDex == nil {
		s.logger.Debug("Token successfully verified by Dex OIDC Verifier")
		var dexClaims DexClaims
		if err := idToken.Claims(&dexClaims); err != nil {
			s.logger.Error("Failed to extract claims from Dex verified token", zap.Error(err))
			return nil, fmt.Errorf("failed to extract Dex token claims: %w", err)
		}
		// Map Dex claims to userClaim (Role added later in Check)
		claim := &userClaim{Email: dexClaims.Email, EmailVerified: dexClaims.EmailVerified, ExternalUserID: dexClaims.Subject, StandardClaims: dexClaims.StandardClaims}
		return claim, nil
	}
	s.logger.Debug("Dex OIDC verification failed, attempting platform key verification", zap.Error(errDex)) // Log Dex error and continue

	// --- Attempt 2: Verify as Platform Token/API Key ---
	if s.platformPublicKey == nil {
		s.logger.Debug("Platform public key not configured, skipping platform verification")
		// If Dex verification also failed, return its error or a generic one
		if errDex != nil {
			return nil, fmt.Errorf("token verification failed (Dex verify error: %w)", errDex)
		}
		return nil, fmt.Errorf("platform key not available for token verification")
	}

	s.logger.Debug("Attempting token verification via Platform Public Key")
	var platformClaims userClaim // Parse directly into our claim struct
	token, errPlatform := jwt.ParseWithClaims(tokenString, &platformClaims, func(token *jwt.Token) (interface{}, error) {
		// Keyfunc: Validate alg and kid before returning key
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method for platform token: %v", token.Header["alg"])
		}
		if kid, ok := token.Header["kid"].(string); ok {
			if kid != s.platformKeyID {
				return nil, fmt.Errorf("token 'kid' [%s] does not match expected platform key ID [%s]", kid, s.platformKeyID)
			}
			s.logger.Debug("Verifying platform token using kid", zap.String("kid", kid))
		} else {
			s.logger.Warn("Platform token header missing 'kid', verification may be ambiguous if keys rotate.") /* return error if kid is mandatory */
		}
		return s.platformPublicKey, nil
	})

	// Check parsing/validation result
	if errPlatform == nil && token.Valid { // token.Valid confirms signature, time claims, and custom Valid() passed
		s.logger.Debug("Token successfully verified by Platform Public Key")
		if platformClaims.Subject == "" {
			// Should have been caught by platformClaims.Valid(), but double-check
			return nil, fmt.Errorf("verified platform token is missing 'sub' claim")
		}
		platformClaims.ExternalUserID = platformClaims.Subject // Ensure ExternalUserID is set from Subject
		return &platformClaims, nil
	}

	// Log specific platform verification error
	s.logger.Debug("Platform key verification failed", zap.String("expected_kid", s.platformKeyID), zap.Error(errPlatform))

	// If both Dex and Platform verification failed, return a consolidated error
	// Prioritize returning the platform error if it exists, otherwise Dex error
	finalErr := errors.New("token verification failed") // Generic fallback
	if errPlatform != nil {
		finalErr = fmt.Errorf("platform token validation failed: %w", errPlatform)
	} else if errDex != nil {
		// Only use Dex error if platform attempt didn't yield a more specific error
		finalErr = fmt.Errorf("OIDC token validation failed: %w", errDex)
	}
	return nil, finalErr
}
