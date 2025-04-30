package auth

import (
	"context"
	"crypto/rsa"
	"crypto/sha512"
	_ "embed" // Keep if needed for email templates later
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/api"       // Renamed import alias
	"github.com/opengovern/opensecurity/services/auth/authcache" // Import authcache
	"github.com/opengovern/opensecurity/services/auth/db"
	"github.com/opengovern/opensecurity/services/auth/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// httpRoutes holds dependencies for HTTP handlers.
type httpRoutes struct {
	logger             *zap.Logger
	platformPrivateKey *rsa.PrivateKey
	db                 db.Database
	authCache          *authcache.AuthCacheService // Injected cache service
	authServer         *Server                     // Injected main server logic
}

// Register defines and registers all HTTP routes for the auth service.
func (r *httpRoutes) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1") // Base path for v1 API

	// Health/Check Endpoint (usually public or less restricted)
	v1.GET("/check", r.Check) // Envoy external auth check endpoint

	// User Management Endpoints (Authorization applied per route)
	v1.GET("/users", httpserver.AuthorizeHandler(r.GetUsers, api2.EditorRole))
	v1.GET("/user/:id", httpserver.AuthorizeHandler(r.GetUserDetails, api2.EditorRole))
	v1.GET("/me", httpserver.AuthorizeHandler(r.GetMe, api2.ViewerRole)) // Viewer role might be sufficient for /me
	v1.POST("/user", httpserver.AuthorizeHandler(r.CreateUser, api2.EditorRole))
	v1.PUT("/user", httpserver.AuthorizeHandler(r.UpdateUser, api2.EditorRole))
	v1.GET("/user/password/check", httpserver.AuthorizeHandler(r.CheckUserPasswordChangeRequired, api2.ViewerRole))
	v1.POST("/user/password/reset", httpserver.AuthorizeHandler(r.ResetUserPassword, api2.ViewerRole))
	v1.DELETE("/user/:id", httpserver.AuthorizeHandler(r.DeleteUser, api2.AdminRole))

	// API Key Management Endpoints
	v1.POST("/keys", httpserver.AuthorizeHandler(r.CreateAPIKey, api2.AdminRole))
	v1.GET("/keys", httpserver.AuthorizeHandler(r.ListAPIKeys, api2.AdminRole))
	v1.DELETE("/key/:id", httpserver.AuthorizeHandler(r.DeleteAPIKey, api2.AdminRole))
	v1.PUT("/key/:id", httpserver.AuthorizeHandler(r.EditAPIKey, api2.AdminRole))

	// Connector Management Endpoints
	v1.GET("/connectors", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole))
	v1.GET("/connectors/supported-connector-types", httpserver.AuthorizeHandler(r.GetSupportedType, api2.AdminRole))
	v1.GET("/connector/:type", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole)) // Reuse GetConnectors, filtering happens inside
	v1.POST("/connector", httpserver.AuthorizeHandler(r.CreateConnector, api2.AdminRole))
	v1.POST("/connector/auth0", httpserver.AuthorizeHandler(r.CreateAuth0Connector, api2.AdminRole)) // Specific endpoint for Auth0
	v1.PUT("/connector", httpserver.AuthorizeHandler(r.UpdateConnector, api2.AdminRole))
	v1.DELETE("/connector/:id", httpserver.AuthorizeHandler(r.DeleteConnector, api2.AdminRole)) // Assuming delete by ConnectorID (string)
}

// bindValidate is a helper to bind and validate request data.
func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		// Return a more specific error for binding issues
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
	}
	if err := ctx.Validate(i); err != nil {
		// Return validation errors
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
	}
	return nil
}

// Check implements the Envoy External Authorization check endpoint.
func (r *httpRoutes) Check(ctx echo.Context) error {
	// Construct the CheckRequest for the auth server logic
	checkRequest := envoyauth.CheckRequest{
		Attributes: &envoyauth.AttributeContext{
			Request: &envoyauth.AttributeContext_Request{
				Http: &envoyauth.AttributeContext_HttpRequest{
					Headers: make(map[string]string),
				},
			},
		},
	}

	// Copy headers from incoming request
	for k, v := range ctx.Request().Header {
		// Envoy expects single values for headers in this map
		if len(v) > 0 {
			checkRequest.Attributes.Request.Http.Headers[strings.ToLower(k)] = v[0] // Lowercase header names
		}
	}

	// Extract original URI and Method from proxy headers (adjust header names if needed)
	originalUriStr := ctx.Request().Header.Get("X-Original-URI")
	originalMethod := ctx.Request().Header.Get("X-Original-Method")
	if originalUriStr == "" {
		// Fallback or error if header is missing
		r.logger.Warn("Missing X-Original-URI header in /check request")
		originalUriStr = ctx.Request().URL.RequestURI() // Use request URI as fallback
	}
	if originalMethod == "" {
		r.logger.Warn("Missing X-Original-Method header in /check request")
		originalMethod = ctx.Request().Method // Use request method as fallback
	}

	originalUri, err := url.Parse(originalUriStr)
	if err != nil {
		r.logger.Warn("Failed to parse X-Original-URI", zap.String("uri", originalUriStr), zap.Error(err))
		// Decide how to handle parse error - maybe deny? For now, use path from request URL
		checkRequest.Attributes.Request.Http.Path = ctx.Request().URL.Path
	} else {
		checkRequest.Attributes.Request.Http.Path = originalUri.Path
	}
	checkRequest.Attributes.Request.Http.Method = originalMethod

	// Call the core Check logic (likely involving cache/DB lookups)
	res, err := r.authServer.Check(ctx.Request().Context(), &checkRequest)
	if err != nil {
		// Log internal server errors during check
		r.logger.Error("Error during auth check execution", zap.Error(err))
		// Return a generic internal error to the client (Envoy)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal authorization error")
	}

	// Process the response from the auth server
	if res.Status.Code != int32(codes.OK) {
		// Denied access
		// Consider logging denied reason if available and appropriate
		r.logger.Info("Authorization check denied",
			zap.Int32("code", res.Status.Code),
			zap.String("message", res.Status.Message),
			zap.String("path", checkRequest.Attributes.Request.Http.Path),
			zap.String("method", checkRequest.Attributes.Request.Http.Method))
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized") // Simple unauthorized
	}

	// Allowed access
	okResponse := res.GetOkResponse()
	if okResponse == nil {
		// This case should ideally not happen if status code is OK, but handle defensively
		r.logger.Error("Auth check returned OK status but missing OkResponse")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal authorization configuration error")
	}

	// Add headers from the OkResponse to the upstream request
	for _, headerValueOption := range okResponse.GetHeaders() {
		if headerValueOption != nil && headerValueOption.Header != nil {
			// Use Add to handle potential multiple headers with the same key if needed, Set overwrites
			ctx.Response().Header().Set(headerValueOption.Header.Key, headerValueOption.Header.Value)
		}
	}

	// Return 200 OK to Envoy, allowing the original request to proceed
	return ctx.NoContent(http.StatusOK)
}

// GetUsers retrieves a list of all users.
func (r *httpRoutes) GetUsers(ctx echo.Context) error {
	// Note: Binding GetUsersRequest is currently not used for filtering.
	// Add filtering logic here if needed based on req.
	users, err := r.db.GetUsers()
	if err != nil {
		r.logger.Error("Failed to get users from database", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	// Map database users to API response users
	resp := make([]api.GetUsersResponse, 0, len(users))
	for _, u := range users {
		apiUser := api.GetUsersResponse{
			ID:            u.ID,
			UserName:      u.Username,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			ExternalId:    u.ExternalId,
			CreatedAt:     u.CreatedAt,
			RoleName:      u.Role,
			IsActive:      u.IsActive,
			FullName:      u.FullName,
			ConnectorId:   u.ConnectorId,
		}
		// Handle zero time for LastActivity
		if !u.LastLogin.IsZero() {
			apiUser.LastActivity = &u.LastLogin
		}
		resp = append(resp, apiUser)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// GetUserDetails retrieves details for a specific user by ID.
func (r *httpRoutes) GetUserDetails(ctx echo.Context) error {
	userIDStr := ctx.Param("id")
	// ID is likely numeric based on DB schema, attempt conversion
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		r.logger.Warn("Invalid user ID format in request", zap.String("id", userIDStr), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	// Fetch user by numeric ID
	user, err := r.db.GetUser(userIDStr) // Assuming GetUser handles string->uint conversion or takes string
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("User not found by ID", zap.Uint64("id", userIDUint))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user details from database", zap.Uint64("id", userIDUint), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user details")
	}
	if user == nil { // Should be caught by ErrRecordNotFound, but double-check
		r.logger.Warn("GetUser returned nil user without ErrRecordNotFound", zap.Uint64("id", userIDUint))
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Map to response
	resp := api.GetUserResponse{
		ID:            user.ID,
		UserName:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		Blocked:       !user.IsActive, // Map IsActive to Blocked
		RoleName:      user.Role,
	}
	if !user.LastLogin.IsZero() {
		resp.LastActivity = &user.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)
}

// GetMe retrieves details for the currently authenticated user.
func (r *httpRoutes) GetMe(ctx echo.Context) error {
	// Get user ID (External ID) from the authenticated context
	externalUserID := httpserver.GetUserID(ctx)
	if externalUserID == "" {
		// This should ideally not happen if AuthorizeHandler is working correctly
		r.logger.Error("Unable to get user ID from context in /me handler")
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot identify authenticated user")
	}

	// Use the utility function which handles potential errors like user not found/disabled
	user, err := utils.GetUser(externalUserID, r.db) // GetUser uses External ID
	if err != nil {
		if errors.Is(err, errors.New("user not found")) || errors.Is(err, errors.New("user disabled")) {
			r.logger.Warn("Attempt to get /me for non-existent or disabled user", zap.String("externalId", externalUserID), zap.Error(err))
			// Return 404 or 403 depending on policy
			return echo.NewHTTPError(http.StatusNotFound, "User not found or is inactive")
		}
		r.logger.Error("Failed to get user details for /me", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user details")
	}

	// Map to response
	resp := api.GetMeResponse{
		ID:            user.ID,
		UserName:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		Blocked:       !user.IsActive,
		Role:          user.Role,
		MemberSince:   user.CreatedAt,
		ConnectorId:   user.ConnectorId,
		// ColorBlindMode: user.ColorBlindMode, // Add if this field exists in utils.User
	}
	if !user.LastLogin.IsZero() {
		resp.LastLogin = &user.LastLogin
		resp.LastActivity = &user.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)
}

// CreateAPIKey generates a new API key (JWT) for the authenticated user.
func (r *httpRoutes) CreateAPIKey(ctx echo.Context) error {
	externalUserID := httpserver.GetUserID(ctx) // Get External ID of the user creating the key
	var req api.CreateAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err // bindValidate returns appropriate HTTP errors
	}

	// Fetch the user creating the key to get their email (needed for claim)
	// Use the DB directly here as utils.GetUser might return a limited struct
	creatorUser, err := r.db.GetUserByExternalID(externalUserID)
	if err != nil || creatorUser == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || creatorUser == nil {
			r.logger.Error("Creator user not found in DB for CreateAPIKey", zap.String("externalId", externalUserID))
			return echo.NewHTTPError(http.StatusNotFound, "Authenticated user not found")
		}
		r.logger.Error("Failed to get creator user for CreateAPIKey", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}

	// Check key limit
	currentKeyCount, err := r.db.CountApiKeysForUser(externalUserID)
	if err != nil {
		r.logger.Error("Failed to count API keys for user", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check key limits")
	}
	// TODO: Make key limit configurable
	if currentKeyCount >= 5 {
		r.logger.Warn("API key limit reached for user", zap.String("externalId", externalUserID), zap.Int64("limit", 5))
		return echo.NewHTTPError(http.StatusNotAcceptable, "Maximum number of API keys (5) reached for user")
	}

	// Prepare claims for the new API key JWT
	// Note: This JWT represents the API key itself, not the user's session
	apiKeyClaims := userClaim{ // Assuming userClaim struct is defined in server.go or common place
		Role:           req.Role,          // Role assigned to the key
		Email:          creatorUser.Email, // Email of the user *creating* the key
		ExternalUserID: externalUserID,    // External ID of the user *creating* the key
		// Add IssuedAt, maybe Expiry if keys should expire?
		// StandardClaims: jwt.StandardClaims{ IssuedAt: time.Now().Unix() },
	}

	// Ensure platform private key is available for signing
	if r.platformPrivateKey == nil {
		r.logger.Error("Platform private key is not configured, cannot create API key")
		return echo.NewHTTPError(http.StatusInternalServerError, "API key generation is disabled")
	}

	// Sign the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &apiKeyClaims).SignedString(r.platformPrivateKey)
	if err != nil {
		r.logger.Error("Failed to sign API key token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate API key token")
	}

	// Generate masked key and hash for storage
	maskedKey := "invalid"
	if len(token) > 20 { // Basic check to avoid panic on very short tokens
		maskedKey = fmt.Sprintf("%s...%s", token[:10], token[len(token)-10:])
	}
	hash := sha512.New()
	if _, err = hash.Write([]byte(token)); err != nil { // Check hash write error
		r.logger.Error("Failed to hash API key token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to secure API key")
	}
	keyHash := hex.EncodeToString(hash.Sum(nil))

	// Create DB record for the key
	apiKeyRecord := db.ApiKey{
		Name:          req.Name,
		Role:          req.Role,
		CreatorUserID: externalUserID, // Store the creator's external ID
		IsActive:      true,
		MaskedKey:     maskedKey,
		KeyHash:       keyHash, // Store the hash of the full token
	}

	// Add to database
	if err = r.db.AddApiKey(&apiKeyRecord); err != nil {
		r.logger.Error("Failed to save API key record to database", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save API key")
	}
	r.logger.Info("Created new API key", zap.Uint("id", apiKeyRecord.ID), zap.String("name", req.Name), zap.String("creator", externalUserID))

	// Return response including the full token (only time it's shown)
	return ctx.JSON(http.StatusCreated, api.CreateAPIKeyResponse{
		ID:        apiKeyRecord.ID,
		Name:      apiKeyRecord.Name,
		Active:    apiKeyRecord.IsActive,
		CreatedAt: apiKeyRecord.CreatedAt,
		RoleName:  apiKeyRecord.Role,
		Token:     token, // Return the actual JWT token
	})
}

// DeleteAPIKey deletes an API key by its database ID.
func (r *httpRoutes) DeleteAPIKey(ctx echo.Context) error {
	idStr := ctx.Param("id")
	apiKeyID, err := strconv.ParseUint(idStr, 10, 64) // Use ParseUint for uint64 ID
	if err != nil {
		r.logger.Warn("Invalid API key ID format for deletion", zap.String("id", idStr), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid API key ID format")
	}

	// Perform deletion
	err = r.db.DeleteAPIKey(apiKeyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("API key not found for deletion", zap.Uint64("id", apiKeyID))
			return echo.NewHTTPError(http.StatusNotFound, "API key not found")
		}
		r.logger.Error("Failed to delete API key from database", zap.Uint64("id", apiKeyID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete API key")
	}

	r.logger.Info("Deleted API key", zap.Uint64("id", apiKeyID))
	return ctx.NoContent(http.StatusAccepted) // 202 Accepted or 204 No Content are appropriate
}

// EditAPIKey updates the role and active status of an API key.
func (r *httpRoutes) EditAPIKey(ctx echo.Context) error {
	idStr := ctx.Param("id") // Get ID from path
	var req api.EditAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// Validate ID format (UpdateAPIKey expects string, but let's ensure it's reasonable)
	if _, err := strconv.ParseUint(idStr, 10, 64); err != nil {
		r.logger.Warn("Invalid API key ID format for edit", zap.String("id", idStr), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid API key ID format")
	}

	// Perform update
	err := r.db.UpdateAPIKey(idStr, req.IsActive, req.Role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("API key not found for edit", zap.String("id", idStr))
			return echo.NewHTTPError(http.StatusNotFound, "API key not found")
		}
		r.logger.Error("Failed to update API key in database", zap.String("id", idStr), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update API key")
	}

	r.logger.Info("Edited API key", zap.String("id", idStr), zap.Bool("isActive", req.IsActive), zap.String("role", string(req.Role)))
	return ctx.NoContent(http.StatusAccepted) // 202 Accepted or 200 OK / 204 No Content
}

// ListAPIKeys lists API keys for the currently authenticated user.
func (r *httpRoutes) ListAPIKeys(ctx echo.Context) error {
	externalUserID := httpserver.GetUserID(ctx)
	if externalUserID == "" {
		r.logger.Error("Unable to get user ID from context in ListAPIKeys")
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot identify authenticated user")
	}

	keys, err := r.db.ListApiKeysForUser(externalUserID)
	if err != nil {
		r.logger.Error("Failed to list API keys for user", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve API keys")
	}

	// Map to response
	resp := make([]api.APIKeyResponse, 0, len(keys))
	for _, key := range keys {
		resp = append(resp, api.APIKeyResponse{
			ID:            key.ID,
			CreatedAt:     key.CreatedAt,
			UpdatedAt:     key.UpdatedAt, // Include UpdatedAt if needed
			Name:          key.Name,
			RoleName:      key.Role,
			CreatorUserID: key.CreatorUserID,
			Active:        key.IsActive,
			MaskedKey:     key.MaskedKey,
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

// CreateUser creates a new user account.
func (r *httpRoutes) CreateUser(ctx echo.Context) error {
	var req api.CreateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// Call internal logic, handle potential HTTP errors returned
	err := r.DoCreateUser(req)
	if err != nil {
		// DoCreateUser should return appropriate echo.HTTPError
		return err
	}

	// Invalidate cache if creation was successful and user might be cached immediately (though unlikely)
	// It's safer to invalidate on update/delete.
	// _ = r.authCache.RemoveUserFromCache(ctx.Request().Context(), req.EmailAddress)

	return ctx.NoContent(http.StatusCreated)
}

// DoCreateUser contains the core logic for creating a user.
func (r *httpRoutes) DoCreateUser(req api.CreateUserRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.EmailAddress))
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is required")
	}
	// --- Add this block for validation ---
	if _, err := mail.ParseAddress(email); err != nil {
		r.logger.Warn("Invalid email format provided", zap.String("email", req.EmailAddress), zap.Error(err)) // Log original input
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid email address format: %s", req.EmailAddress))
	}
	// --- End validation block ---

	// Check if user already exists
	existingUser, err := r.db.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Handle DB errors other than not found
		r.logger.Error("Database error checking for existing user", zap.String("email", email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check user existence")
	}
	if existingUser != nil {
		r.logger.Warn("Attempted to create user with existing email", zap.String("email", email))
		return echo.NewHTTPError(http.StatusConflict, "Email address already in use") // 409 Conflict is appropriate
	}

	// Determine if this is the very first user (admin creation logic)
	count, err := r.db.GetUsersCount()
	if err != nil {
		r.logger.Error("Failed to get users count", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check user count")
	}
	isFirstUser := count == 0

	// Assign role - default to Viewer unless specified, force Admin for first user
	role := api2.ViewerRole
	if isFirstUser {
		role = api2.AdminRole
		r.logger.Info("Creating first user, assigning admin role", zap.String("email", email))
	} else if req.Role != nil {
		// TODO: Validate req.Role value against allowed roles
		role = *req.Role
	}

	// Determine connector and external ID
	connectorID := req.ConnectorId // Use provided connector ID
	if connectorID == "" {
		// If no connector specified, assume local password auth if password is provided
		if req.Password != nil && *req.Password != "" {
			connectorID = "local"
		} else {
			// Cannot create user without a connector or a password for local auth
			return echo.NewHTTPError(http.StatusBadRequest, "Connector ID or password required for user creation")
		}
	}
	externalID := fmt.Sprintf("%s|%s", connectorID, email) // Construct external ID

	// Handle Dex password creation if it's a local user
	if connectorID == "local" {
		if req.Password == nil || *req.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Password is required for local user creation")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("Failed to hash password", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Password hashing failed")
		}

		// Create or update password in Dex
		dexReq := &dexApi.CreatePasswordReq{
			Password: &dexApi.Password{
				UserId: externalID, // Use the constructed external ID
				Email:  email,
				Hash:   hashedPassword,
				// Username: email, // Optional: Dex might use email as username by default
			},
		}
		resp, err := r.authServer.dexClient.CreatePassword(context.TODO(), dexReq) // Use request context?
		if err != nil {
			// Check for specific Dex errors if possible
			r.logger.Error("Failed to create dex password", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set up user password")
		}
		if resp.AlreadyExists {
			// This shouldn't happen if GetUserByEmail check passed, but handle defensively
			r.logger.Warn("Dex password already exists for new user, attempting update", zap.String("email", email))
			updateReq := &dexApi.UpdatePasswordReq{
				Email:   email,
				NewHash: hashedPassword,
			}
			_, err = r.authServer.dexClient.UpdatePassword(context.TODO(), updateReq)
			if err != nil {
				r.logger.Error("Failed to update potentially existing dex password", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set up user password")
			}
		}
	}

	// Prepare user record for DB insertion
	newUser := &db.User{
		Email:                 email,
		Username:              email, // Default username to email
		FullName:              email, // Default full name to email
		Role:                  role,
		EmailVerified:         false, // Require verification?
		ConnectorId:           connectorID,
		ExternalId:            externalID,
		RequirePasswordChange: !isFirstUser, // First user doesn't need to change password immediately
		IsActive:              req.IsActive, // Use IsActive from request
	}

	// Create user in the database
	err = r.db.CreateUser(newUser) // CreateUser now handles OnConflict based on external_id
	if err != nil {
		r.logger.Error("Failed to create user in database", zap.String("email", email), zap.Error(err))
		// Check for specific DB errors (e.g., unique constraint violation if OnConflict fails)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user data")
	}

	r.logger.Info("Successfully created user", zap.String("email", email), zap.Uint("id", newUser.ID), zap.String("role", string(role)))
	return nil
}

// UpdateUser updates an existing user's details.
func (r *httpRoutes) UpdateUser(ctx echo.Context) error {
	var req api.UpdateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	email := strings.ToLower(strings.TrimSpace(req.EmailAddress))
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is required")
	}

	// Get the existing user by email
	user, err := r.db.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("User not found for update", zap.String("email", email))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for update", zap.String("email", email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user")
	}
	if user == nil { // Should be caught by ErrRecordNotFound
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Track if any changes were made that require cache invalidation
	cacheNeedsInvalidation := false

	// Update password if provided and connector is local
	if req.Password != nil && *req.Password != "" {
		if user.ConnectorId != "local" {
			return echo.NewHTTPError(http.StatusBadRequest, "Password can only be set for local users")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("Failed to hash password during update", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Password hashing failed")
		}

		// Update Dex password
		updateReq := &dexApi.UpdatePasswordReq{
			Email:   email,
			NewHash: hashedPassword,
		}
		resp, err := r.authServer.dexClient.UpdatePassword(context.TODO(), updateReq) // Use request context?
		if err != nil {
			r.logger.Error("Failed to update dex password", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
		}
		if resp.NotFound {
			// Password didn't exist, try creating it (edge case?)
			r.logger.Warn("Dex password not found during update, attempting create", zap.String("email", email))
			createReq := &dexApi.CreatePasswordReq{
				Password: &dexApi.Password{
					UserId: user.ExternalId, // Use existing external ID
					Email:  email,
					Hash:   hashedPassword,
				},
			}
			_, err = r.authServer.dexClient.CreatePassword(context.TODO(), createReq)
			if err != nil {
				r.logger.Error("Failed to create dex password after update attempt failed", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set password")
			}
		}

		// Update password change flag in our DB
		if user.RequirePasswordChange { // Only update if it was true
			err = r.db.UserPasswordUpdate(user.ID)
			if err != nil {
				r.logger.Error("Failed to update user password change flag", zap.Uint("id", user.ID), zap.Error(err))
				// Log but don't necessarily fail the whole request
			}
		}
	}

	// Update other user fields
	// Use a map for targeted updates to avoid overwriting fields unintentionally
	updateData := make(map[string]interface{})
	if req.Role != nil && user.Role != *req.Role {
		// TODO: Add validation for allowed roles
		updateData["role"] = *req.Role
		cacheNeedsInvalidation = true // Role change invalidates cache
	}
	if user.IsActive != req.IsActive {
		updateData["is_active"] = req.IsActive
		cacheNeedsInvalidation = true // Active status change invalidates cache
	}
	if req.UserName != "" && user.Username != req.UserName {
		updateData["username"] = req.UserName
	}
	if req.FullName != "" && user.FullName != req.FullName {
		updateData["full_name"] = req.FullName
	}
	// Handle potential connector change - this might need more complex logic
	// if changing connector implies changing external ID format.
	if req.ConnectorId != "" && user.ConnectorId != req.ConnectorId {
		updateData["connector_id"] = req.ConnectorId
		// IMPORTANT: If external ID format depends on connector, update it too!
		newExternalId := fmt.Sprintf("%s|%s", req.ConnectorId, user.Email)
		if user.ExternalId != newExternalId {
			updateData["external_id"] = newExternalId
			cacheNeedsInvalidation = true // External ID change invalidates cache
		}
	}

	// Perform the update only if there are changes
	if len(updateData) > 0 {
		tx := r.db.Orm.Model(&db.User{}).Where("id = ?", user.ID).Updates(updateData)
		if tx.Error != nil {
			r.logger.Error("Failed to update user in database", zap.Uint("id", user.ID), zap.Error(tx.Error))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user data")
		}
		if tx.RowsAffected == 0 {
			// Should not happen if GetUserByEmail succeeded, but check defensively
			r.logger.Warn("User found but no rows affected by update", zap.Uint("id", user.ID))
			// Return NotFound or InternalError?
			return echo.NewHTTPError(http.StatusNotFound, "User found but update failed")
		}
		r.logger.Info("Updated user details", zap.Uint("id", user.ID), zap.Any("changes", updateData))
	} else {
		r.logger.Info("No changes detected for user update", zap.Uint("id", user.ID))
	}

	// Invalidate cache if relevant fields changed
	if cacheNeedsInvalidation {
		err = r.authCache.RemoveUserFromCache(ctx.Request().Context(), email)
		if err != nil {
			// Log cache error but don't fail the request
			r.logger.Error("Failed to invalidate user cache after update", zap.String("email", email), zap.Error(err))
		}
	}

	return ctx.NoContent(http.StatusOK) // 200 OK or 204 No Content
}

// DeleteUser deletes a user account.
func (r *httpRoutes) DeleteUser(ctx echo.Context) error {
	userIDStr := ctx.Param("id")
	if userIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required in path")
	}

	// Call internal logic
	err := r.DoDeleteUser(userIDStr) // Pass the ID string
	if err != nil {
		// DoDeleteUser returns appropriate echo.HTTPError
		return err
	}

	// No need to invalidate cache here, DoDeleteUser handles it

	return ctx.NoContent(http.StatusAccepted) // 202 or 204
}

// DoDeleteUser contains the core logic for deleting a user.
func (r *httpRoutes) DoDeleteUser(idStr string) error {
	// Get user by ID string (GetUser should handle conversion/lookup)
	user, err := r.db.GetUser(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("User not found for deletion", zap.String("id", idStr))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for deletion", zap.String("id", idStr), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user")
	}
	if user == nil { // Should be caught by ErrRecordNotFound
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Prevent deletion of the first user (ID 1)
	if user.ID == 1 {
		r.logger.Warn("Attempted to delete the first user (ID 1)", zap.Uint("id", user.ID))
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot delete the primary admin user")
	}

	userEmail := user.Email // Store email before deleting user record

	// Delete Dex password if it's a local user
	if user.ConnectorId == "local" {
		dexReq := &dexApi.DeletePasswordReq{
			Email: userEmail,
		}
		_, err := r.authServer.dexClient.DeletePassword(context.TODO(), dexReq) // Use request context?
		if err != nil {
			// Log error but proceed with DB deletion attempt
			r.logger.Error("Failed to remove dex password during user deletion", zap.String("email", userEmail), zap.Error(err))
			// Consider if this should be a fatal error for the deletion process
		} else {
			r.logger.Info("Removed dex password for deleted user", zap.String("email", userEmail))
		}
	}

	// Delete user from database
	err = r.db.DeleteUser(user.ID)
	if err != nil {
		// Check if it was already deleted?
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("User already deleted from DB?", zap.Uint("id", user.ID))
			// Proceed to cache invalidation anyway
		} else {
			r.logger.Error("Failed to delete user from database", zap.Uint("id", user.ID), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user data")
		}
	} else {
		r.logger.Info("Deleted user from database", zap.Uint("id", user.ID), zap.String("email", userEmail))
	}

	// Invalidate user cache
	// Use context from the original request if possible, otherwise context.Background()
	cacheCtx := context.Background() // Assuming original context isn't easily available here
	err = r.authCache.RemoveUserFromCache(cacheCtx, userEmail)
	if err != nil {
		// Log cache error but don't fail the overall deletion
		r.logger.Error("Failed to invalidate user cache after deletion", zap.String("email", userEmail), zap.Error(err))
	}

	return nil // Indicate success to the calling handler
}

// CheckUserPasswordChangeRequired checks if the current user must change their password.
func (r *httpRoutes) CheckUserPasswordChangeRequired(ctx echo.Context) error {
	externalUserID := httpserver.GetUserID(ctx)
	if externalUserID == "" {
		r.logger.Error("Unable to get user ID from context in CheckUserPasswordChangeRequired")
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot identify authenticated user")
	}

	user, err := r.db.GetUserByExternalID(externalUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("User not found for password change check", zap.String("externalId", externalUserID))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for password change check", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if user.RequirePasswordChange {
		return ctx.String(http.StatusOK, "CHANGE_REQUIRED")
	} else {
		return ctx.String(http.StatusOK, "CHANGE_NOT_REQUIRED")
	}
}

// ResetUserPassword allows the authenticated user to change their own password.
func (r *httpRoutes) ResetUserPassword(ctx echo.Context) error {
	externalUserID := httpserver.GetUserID(ctx)
	if externalUserID == "" {
		r.logger.Error("Unable to get user ID from context in ResetUserPassword")
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot identify authenticated user")
	}

	// Get user from DB
	user, err := r.db.GetUserByExternalID(externalUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("User not found for password reset", zap.String("externalId", externalUserID))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for password reset", zap.String("externalId", externalUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Bind and validate request payload
	var req api.ResetUserPasswordRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// Ensure it's a local user
	if user.ConnectorId != "local" {
		r.logger.Warn("Attempt to reset password for non-local user", zap.String("externalId", externalUserID), zap.String("connector", user.ConnectorId))
		return echo.NewHTTPError(http.StatusBadRequest, "Password reset only available for local users")
	}

	// Verify current password with Dex
	verifyReq := &dexApi.VerifyPasswordReq{
		Email:    user.Email,
		Password: req.CurrentPassword,
	}
	resp, err := r.authServer.dexClient.VerifyPassword(context.TODO(), verifyReq) // Use request context?
	if err != nil {
		r.logger.Error("Failed to call Dex VerifyPassword", zap.String("email", user.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify current password")
	}
	if resp.NotFound {
		// Should not happen if user exists in our DB
		r.logger.Error("User found in local DB but not in Dex passwords", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusNotFound, "Password record not found")
	}
	if !resp.Verified {
		r.logger.Info("Incorrect current password provided for password reset", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect current password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash new password", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Password hashing failed")
	}

	// Update the password in Dex
	passwordUpdateReq := &dexApi.UpdatePasswordReq{
		Email:   user.Email,
		NewHash: hashedPassword,
	}
	_, err = r.authServer.dexClient.UpdatePassword(context.TODO(), passwordUpdateReq) // Use request context?
	if err != nil {
		// Handle case where password might have been deleted between Verify and Update? Unlikely.
		r.logger.Error("Failed to update dex password after verification", zap.String("email", user.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	// Update the require_password_change flag in our database
	if user.RequirePasswordChange { // Only update if it was true
		err = r.db.UserPasswordUpdate(user.ID)
		if err != nil {
			r.logger.Error("Failed to update user password change flag after reset", zap.Uint("id", user.ID), zap.Error(err))
			// Log but don't fail the request
		}
	}

	// Invalidate cache as user state (password change flag) might be implicitly part of authz check later
	err = r.authCache.RemoveUserFromCache(ctx.Request().Context(), user.Email)
	if err != nil {
		r.logger.Error("Failed to invalidate user cache after password reset", zap.String("email", user.Email), zap.Error(err))
	}

	r.logger.Info("User successfully reset password", zap.String("email", user.Email), zap.Uint("id", user.ID))
	return ctx.NoContent(http.StatusAccepted)
}

// GetConnectors lists configured identity connectors.
func (r *httpRoutes) GetConnectors(ctx echo.Context) error {
	req := &dexApi.ListConnectorReq{}
	filterType := strings.ToLower(ctx.Param("type")) // Get optional type filter from path

	// Get connectors from Dex
	respDex, err := r.authServer.dexClient.ListConnectors(context.TODO(), req) // Use request context?
	if err != nil {
		r.logger.Error("Failed to list connectors from Dex", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve connectors")
	}

	connectors := respDex.Connectors
	resp := make([]api.GetConnectorsResponse, 0, len(connectors))

	// Enrich Dex info with data from our DB (like UserCount, SubType)
	for _, dexConnector := range connectors {
		// Skip the built-in 'local' connector unless explicitly requested (unlikely)
		if dexConnector.Id == "local" && filterType != "local" {
			continue
		}

		// Apply type filter if provided
		if filterType != "" && strings.ToLower(dexConnector.Type) != filterType {
			continue
		}

		// Get corresponding record from our DB
		localConnector, err := r.db.GetConnectorByConnectorID(dexConnector.Id)
		if err != nil {
			// If not found in our DB, log warning but potentially still show Dex info? Or skip?
			// Skipping for now, assuming our DB should be in sync for configured connectors.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Warn("Connector found in Dex but not in local DB, skipping", zap.String("connectorId", dexConnector.Id))
				continue
			}
			r.logger.Error("Failed to get local connector info from DB", zap.String("connectorId", dexConnector.Id), zap.Error(err))
			// Maybe return partial list or error out? Returning error for now.
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve full connector details")
		}
		if localConnector == nil { // Should be caught by ErrRecordNotFound
			r.logger.Warn("GetConnectorByConnectorID returned nil without error", zap.String("connectorId", dexConnector.Id))
			continue
		}

		// Build response object
		info := api.GetConnectorsResponse{
			ID:          localConnector.ID, // Our DB primary key
			ConnectorID: dexConnector.Id,   // Dex's ID
			Type:        dexConnector.Type,
			Name:        dexConnector.Name,
			SubType:     localConnector.ConnectorSubType,
			UserCount:   localConnector.UserCount,
			CreatedAt:   localConnector.CreatedAt,
			LastUpdate:  localConnector.LastUpdate,
		}

		// Attempt to extract OIDC details if applicable
		if strings.ToLower(dexConnector.Type) == "oidc" {
			var oidcConfig api.OIDCConfig // Use the struct from your API package
			if err := json.Unmarshal(dexConnector.Config, &oidcConfig); err == nil {
				info.Issuer = oidcConfig.Issuer
				info.ClientID = oidcConfig.ClientID
				info.TenantID = oidcConfig.TenantID // Include TenantID if present
			} else {
				r.logger.Warn("Failed to unmarshal OIDC config for connector", zap.String("connectorId", dexConnector.Id), zap.Error(err))
				// Don't fail the request, just omit the OIDC details
			}
		}

		resp = append(resp, info)
	}
	return ctx.JSON(http.StatusOK, resp)
}

// GetSupportedType lists the connector types supported by the service configuration.
func (r *httpRoutes) GetSupportedType(ctx echo.Context) error {
	// Use the definitions from the utils package
	connectors := make([]api.GetSupportedConnectorTypeResponse, 0, len(utils.SupportedConnectors))

	for connType, subTypes := range utils.SupportedConnectors {
		apiSubTypes := make([]api.ConnectorSubTypes, 0, len(subTypes))
		subTypeNames := utils.SupportedConnectorsNames[connType] // Get corresponding names

		for i, subTypeID := range subTypes {
			name := subTypeID // Default name to ID
			if i < len(subTypeNames) {
				name = subTypeNames[i] // Use pretty name if available
			}
			apiSubTypes = append(apiSubTypes, api.ConnectorSubTypes{
				ID:   subTypeID,
				Name: name,
			})
		}

		connectors = append(connectors, api.GetSupportedConnectorTypeResponse{
			ConnectorType: connType,
			SubTypes:      apiSubTypes,
		})
	}

	return ctx.JSON(http.StatusOK, connectors)
}

// CreateConnector creates a new identity connector in Dex and our DB.
func (r *httpRoutes) CreateConnector(ctx echo.Context) error {
	var req api.CreateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// Basic validation
	if req.ConnectorType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "connector_type is required")
	}
	connectorTypeLower := strings.ToLower(req.ConnectorType)
	connectorSubTypeLower := strings.ToLower(req.ConnectorSubType) // Allow empty subtype initially

	// Get the appropriate creator function based on type
	creator := utils.GetConnectorCreator(connectorTypeLower)
	if creator == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Connector type '%s' is not supported", req.ConnectorType))
	}

	// Handle subtype validation and default assignment
	if connectorSubTypeLower == "" {
		if connectorTypeLower == "oidc" { // Default OIDC subtype if not provided
			connectorSubTypeLower = "general"
			req.ConnectorSubType = "general" // Update request struct for consistency
			r.logger.Info("No connector_sub_type specified for OIDC, defaulting to 'general'")
		} else {
			// If other types require subtypes, add checks here
		}
	}

	if !utils.IsSupportedSubType(connectorTypeLower, connectorSubTypeLower) {
		err := fmt.Sprintf("Unsupported connector_sub_type '%s' for connector_type '%s'", req.ConnectorSubType, req.ConnectorType)
		r.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Assign default ID and Name if not provided by user, based on subtype
	// This logic might be better placed within the utils creator functions
	if strings.TrimSpace(req.ID) == "" {
		// Generate default ID based on type/subtype
		req.ID = fmt.Sprintf("%s-%s-default", connectorTypeLower, connectorSubTypeLower)
		// Handle potential conflicts if user creates multiple defaults? Maybe add random suffix?
		r.logger.Info("Assigning default connector ID", zap.String("id", req.ID))
	}
	if strings.TrimSpace(req.Name) == "" {
		// Generate default Name based on type/subtype
		req.Name = fmt.Sprintf("%s (%s)", strings.ToTitle(connectorTypeLower), strings.ToTitle(connectorSubTypeLower))
		r.logger.Info("Assigning default connector Name", zap.String("name", req.Name))
	}

	// Validate required fields based on subtype (could also be in creator)
	switch connectorSubTypeLower {
	case "general":
		if strings.TrimSpace(req.Issuer) == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "issuer is required for 'general' OIDC connector")
		}
	case "entraid":
		if strings.TrimSpace(req.TenantID) == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "tenant_id is required for 'entraid' OIDC connector")
		}
	case "google-workspace":
		// Requires ClientID and ClientSecret (validated by struct tags)
		break
	}

	// Prepare request for Dex utils function
	dexUtilRequest := utils.CreateConnectorRequest{
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		ID:               req.ID,
		Name:             req.Name,
	}

	// Create the Dex gRPC request structure
	dexGrpcReq, err := creator(dexUtilRequest)
	if err != nil {
		r.logger.Error("Failed to prepare Dex connector request structure", zap.Error(err))
		// Return specific error if possible (e.g., from fetchEntraIDIssuer)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to prepare connector config: %v", err))
	}

	// Call Dex API to create the connector
	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexGrpcReq) // Use request context?
	if err != nil {
		// Handle potential gRPC errors or specific Dex errors
		r.logger.Error("Failed to create connector in Dex", zap.Error(err))
		// Check for common Dex errors if possible (e.g., invalid config)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create connector in Dex: %v", err))
	}
	if res.AlreadyExists {
		r.logger.Warn("Attempted to create connector that already exists in Dex", zap.String("id", req.ID))
		return echo.NewHTTPError(http.StatusConflict, "Connector with this ID already exists")
	}

	// Create corresponding record in our local database
	err = r.db.CreateConnector(&db.Connector{
		// Gorm automatically handles CreatedAt, UpdatedAt
		ConnectorID:      req.ID, // Use the (potentially defaulted) ID
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		LastUpdate:       time.Now(), // Set initial LastUpdate time
		// UserCount defaults to 0
	})
	if err != nil {
		r.logger.Error("Failed to create connector record in local database", zap.String("id", req.ID), zap.Error(err))
		// Attempt to clean up by deleting the connector from Dex? Complex rollback logic.
		// For now, return error indicating DB failure.
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save connector configuration locally")
	}

	// Restart Dex pod (if applicable and configured)
	// Consider making this conditional or configurable
	r.logger.Info("Attempting to restart Dex pod to apply connector changes...")
	err = utils.RestartDexPod()
	if err != nil {
		// Log error but don't fail the request, as connector was created
		r.logger.Error("Failed to restart Dex pod after connector creation", zap.Error(err))
		// Maybe return a warning in the response?
	}

	r.logger.Info("Successfully created connector", zap.String("id", req.ID), zap.String("type", req.ConnectorType), zap.String("subtype", req.ConnectorSubType))
	// Return 201 Created with the Dex response (which might be empty)
	return ctx.JSON(http.StatusCreated, res) // Or return a custom success message/object
}

// CreateAuth0Connector handles the specific flow for creating an Auth0 OIDC connector.
func (r *httpRoutes) CreateAuth0Connector(ctx echo.Context) error {
	// This endpoint seems redundant if CreateConnector handles different subtypes.
	// If Auth0 requires special handling beyond standard OIDC, keep it, otherwise consider merging.
	// Assuming it might have special client update logic...

	var req api.CreateAuth0ConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// Prepare request for Dex utils function for Auth0
	dexUtilRequest := utils.CreateAuth0ConnectorRequest{
		Issuer:       req.Issuer, // Validate Issuer format?
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Domain:       req.Domain, // Validate Domain format?
	}
	dexGrpcReq, err := utils.CreateAuth0Connector(dexUtilRequest) // Use the specific Auth0 creator
	if err != nil {
		r.logger.Error("Failed to prepare Dex Auth0 connector request structure", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to prepare Auth0 connector config: %v", err))
	}

	// Update Dex public/private clients with URIs from request
	r.logger.Info("Updating Dex clients for Auth0 connector...")
	if err := r.updateDexClients(ctx.Request().Context(), req.PublicURIS, req.PrivateURIS); err != nil {
		// Log the error but decide if it's fatal for connector creation
		r.logger.Error("Failed to update Dex clients for Auth0", zap.Error(err))
		// Potentially return error here if client updates are mandatory
		// return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update Dex client configuration")
	}

	// Call Dex API to create the connector
	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexGrpcReq) // Use request context?
	if err != nil {
		r.logger.Error("Failed to create Auth0 connector in Dex", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create Auth0 connector in Dex: %v", err))
	}
	if res.AlreadyExists {
		// Auth0 connector ID is hardcoded to "auth0" in utils.CreateAuth0Connector
		r.logger.Warn("Auth0 connector already exists in Dex")
		return echo.NewHTTPError(http.StatusConflict, "Auth0 connector already exists")
	}

	// Create corresponding record in our local database
	err = r.db.CreateConnector(&db.Connector{
		ConnectorID:      "auth0", // Hardcoded ID from utils
		ConnectorType:    "oidc",
		ConnectorSubType: "auth0", // Specific subtype
		LastUpdate:       time.Now(),
	})
	if err != nil {
		r.logger.Error("Failed to create Auth0 connector record in local database", zap.Error(err))
		// Attempt rollback?
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save Auth0 connector configuration locally")
	}

	// Restart Dex pod
	r.logger.Info("Attempting to restart Dex pod for Auth0 connector...")
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after Auth0 connector creation", zap.Error(err))
	}

	r.logger.Info("Successfully created Auth0 connector")
	return ctx.JSON(http.StatusCreated, res)
}

// updateDexClients is a helper to update redirect URIs for public/private clients.
func (r *httpRoutes) updateDexClients(ctx context.Context, publicUris, privateUris []string) error {
	// Update Public Client
	if len(publicUris) > 0 {
		publicClientReq := dexApi.UpdateClientReq{
			Id:           "public-client",
			RedirectUris: publicUris,
			// Name: "Public Client", // Avoid resetting name if not intended
		}
		_, err := r.authServer.dexClient.UpdateClient(ctx, &publicClientReq)
		if err != nil {
			r.logger.Error("Failed to update Dex public client URIs", zap.Error(err))
			return fmt.Errorf("failed to update dex public client: %w", err)
		}
		r.logger.Info("Updated Dex public client redirect URIs")
	}

	// Update Private Client
	if len(privateUris) > 0 {
		privateClientReq := dexApi.UpdateClientReq{
			Id:           "private-client",
			RedirectUris: privateUris,
			// Name: "Private Client",
		}
		_, err := r.authServer.dexClient.UpdateClient(ctx, &privateClientReq)
		if err != nil {
			r.logger.Error("Failed to update Dex private client URIs", zap.Error(err))
			return fmt.Errorf("failed to update dex private client: %w", err)
		}
		r.logger.Info("Updated Dex private client redirect URIs")
	}
	return nil
}

// UpdateConnector updates an existing identity connector.
func (r *httpRoutes) UpdateConnector(ctx echo.Context) error {
	var req api.UpdateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	// ID in the request body refers to our *database* ID for the connector record.
	// ConnectorID refers to Dex's ID.
	if req.ID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Database connector ID (id field) is required for update")
	}
	if req.ConnectorID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Dex connector ID (connector_id field) is required for update")
	}

	// Validate connector type/subtype combination
	connectorTypeLower := strings.ToLower(req.ConnectorType)
	connectorSubTypeLower := strings.ToLower(req.ConnectorSubType)
	if !utils.IsSupportedSubType(connectorTypeLower, connectorSubTypeLower) {
		err := fmt.Sprintf("Unsupported connector_sub_type '%s' for connector_type '%s'", req.ConnectorSubType, req.ConnectorType)
		r.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate required fields based on subtype (similar to create)
	switch connectorSubTypeLower {
	case "general":
		if strings.TrimSpace(req.Issuer) == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "issuer is required for 'general' OIDC connector update")
		}
	case "entraid":
		if strings.TrimSpace(req.TenantID) == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "tenant_id is required for 'entraid' OIDC connector update")
		}
	case "google-workspace":
		break // ClientID/Secret validated by struct tags
	}

	// Prepare request for Dex utils update function
	dexUtilRequest := utils.UpdateConnectorRequest{
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		ID:               req.ConnectorID, // Pass Dex's ID to the util function
		Name:             req.Name,        // Pass Name if updating it is supported/desired
	}

	// Create the Dex gRPC update request structure
	dexGrpcReq, err := utils.UpdateOIDCConnector(dexUtilRequest)
	if err != nil {
		r.logger.Error("Failed to prepare Dex connector update request structure", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to prepare connector update config: %v", err))
	}

	// Call Dex API to update the connector
	res, err := r.authServer.dexClient.UpdateConnector(context.TODO(), dexGrpcReq) // Use request context?
	if err != nil {
		r.logger.Error("Failed to update connector in Dex", zap.String("id", req.ConnectorID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to update connector in Dex: %v", err))
	}
	if res.NotFound {
		r.logger.Warn("Connector not found in Dex during update attempt", zap.String("id", req.ConnectorID))
		// If not found in Dex, should we delete our local record or return error?
		// Returning 404 seems appropriate.
		return echo.NewHTTPError(http.StatusNotFound, "Connector not found in identity provider")
	}

	// Update our local database record
	dbConnector := &db.Connector{
		Model: gorm.Model{
			ID: req.ID, // Target our DB record using its primary key
		},
		LastUpdate:       time.Now(),
		ConnectorID:      req.ConnectorID, // Ensure these are updated if they can change
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		// UserCount is likely managed separately
	}
	err = r.db.UpdateConnector(dbConnector)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("Connector found in Dex but corresponding record missing in local DB during update", zap.Uint("dbId", req.ID), zap.String("dexId", req.ConnectorID))
			// Data inconsistency - potentially critical error
			return echo.NewHTTPError(http.StatusInternalServerError, "Data inconsistency: Connector missing locally")
		}
		r.logger.Error("Failed to update connector record in local database", zap.Uint("dbId", req.ID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update local connector configuration")
	}

	// Restart Dex pod (if applicable)
	r.logger.Info("Attempting to restart Dex pod after connector update...")
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after connector update", zap.Error(err))
	}

	r.logger.Info("Successfully updated connector", zap.String("id", req.ConnectorID), zap.Uint("dbId", req.ID))
	return ctx.JSON(http.StatusAccepted, res) // 202 Accepted or 200 OK
}

// DeleteConnector deletes an identity connector from Dex and our DB.
func (r *httpRoutes) DeleteConnector(ctx echo.Context) error {
	// The :id in the route likely refers to Dex's connector ID (string)
	connectorID := ctx.Param("id")
	if connectorID == "" {
		r.logger.Error("Missing connector_id in path for DeleteConnector request")
		return echo.NewHTTPError(http.StatusBadRequest, "Connector ID is required in the URL path")
	}

	// Get our local DB record first to ensure it exists before deleting from Dex
	localConnector, err := r.db.GetConnectorByConnectorID(connectorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("Attempted to delete connector not found in local DB", zap.String("connectorId", connectorID))
			// Connector doesn't exist locally, maybe it doesn't exist in Dex either?
			// Proceed to attempt Dex deletion, but maybe return 404 if Dex also fails?
		} else {
			r.logger.Error("Failed to query local DB before Dex connector deletion", zap.String("connectorId", connectorID), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify connector before deletion")
		}
	}
	// If localConnector is nil here even without error, treat as not found
	if localConnector == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Warn("Attempted to delete connector not found in local DB (nil record)", zap.String("connectorId", connectorID))
	}

	// Delete from Dex
	dexReq := &dexApi.DeleteConnectorReq{
		Id: connectorID,
	}
	resp, err := r.authServer.dexClient.DeleteConnector(context.TODO(), dexReq) // Use request context?
	if err != nil {
		r.logger.Error("Failed to delete connector from Dex", zap.String("id", connectorID), zap.Error(err))
		// Don't fail if Dex returns NotFound, as our goal is deletion anyway
		if !strings.Contains(err.Error(), "not found") { // Basic check, might need refinement based on actual Dex error types
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete connector from identity provider: %v", err))
		}
		r.logger.Warn("Connector not found in Dex during deletion attempt (might be okay)", zap.String("id", connectorID))
	}
	if resp != nil && resp.NotFound { // Check response field too
		r.logger.Warn("Connector reported as not found by Dex during deletion", zap.String("id", connectorID))
		// If it wasn't found locally either, return 404
		if localConnector == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Connector not found")
		}
	}

	// Delete from our local database if it existed
	if localConnector != nil {
		err = r.db.DeleteConnector(connectorID)                    // Delete by ConnectorID string
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Ignore not found error here too
			r.logger.Error("Failed to delete connector record from local database", zap.String("id", connectorID), zap.Error(err))
			// Return error as DB state might be inconsistent now
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete local connector configuration")
		}
	}

	// Invalidate any cache entries related to users from this connector? Complex.
	// For now, rely on user cache TTL or invalidation during user-specific updates.

	// Restart Dex pod (if applicable)
	r.logger.Info("Attempting to restart Dex pod after connector deletion...")
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after connector deletion", zap.Error(err))
	}

	r.logger.Info("Successfully deleted connector (or ensured it was deleted)", zap.String("id", connectorID))
	return ctx.NoContent(http.StatusAccepted) // 202 or 204
}
