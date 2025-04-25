package auth

import (
	"context"
	"crypto/rsa"
	"crypto/sha512"
	_ "embed" // Keep if needed for email templates etc.
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"
	api2 "github.com/opengovern/og-util/pkg/api" // Assuming this is the correct path for Role type
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/db" // Local DB package
	"github.com/opengovern/opensecurity/services/auth/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	// Use v5 if that's standard in your project, otherwise v4
	"github.com/golang-jwt/jwt" // Or jwt "github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/auth/api" // Local API definitions
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DexClaims struct defines the expected claims from a Dex-issued token.
// Keep defined here or import if moved elsewhere
type DexClaims struct {
	Email           string                 `json:"email"`
	EmailVerified   bool                   `json:"email_verified"`
	Groups          []string               `json:"groups"`
	Name            string                 `json:"name"`
	FederatedClaims map[string]interface{} `json:"federated_claims,omitempty"` // Use omitempty if not always present
	jwt.StandardClaims
}

// httpRoutes holds dependencies for HTTP handlers.
type httpRoutes struct {
	logger *zap.Logger

	platformPrivateKey *rsa.PrivateKey
	platformKeyID      string // Stores the calculated Key ID (JWK Thumbprint)
	db                 db.Database
	authServer         *Server // Reference to the main auth server logic
}

// Register registers the HTTP routes with the Echo server.
func (r *httpRoutes) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	// Public / Semi-public endpoints
	v1.GET("/check", r.Check)  // For Envoy auth checks
	v1.POST("/token", r.Token) // For OAuth/OIDC code exchange

	// User Management (protected by roles)
	v1.GET("/users", httpserver.AuthorizeHandler(r.GetUsers, api2.EditorRole))
	v1.GET("/user/:id", httpserver.AuthorizeHandler(r.GetUserDetails, api2.EditorRole)) // ID likely DB uint ID
	v1.GET("/me", httpserver.AuthorizeHandler(r.GetMe, api2.ViewerRole))                // Viewer should be enough
	v1.POST("/user", httpserver.AuthorizeHandler(r.CreateUser, api2.EditorRole))
	v1.PUT("/user", httpserver.AuthorizeHandler(r.UpdateUser, api2.EditorRole))
	v1.GET("/user/password/check", httpserver.AuthorizeHandler(r.CheckUserPasswordChangeRequired, api2.ViewerRole))
	v1.POST("/user/password/reset", httpserver.AuthorizeHandler(r.ResetUserPassword, api2.ViewerRole))
	v1.DELETE("/user/:id", httpserver.AuthorizeHandler(r.DeleteUser, api2.AdminRole)) // ID likely DB uint ID

	// API Key Management (protected by Admin role)
	v1.POST("/keys", httpserver.AuthorizeHandler(r.CreateAPIKey, api2.AdminRole))
	v1.GET("/keys", httpserver.AuthorizeHandler(r.ListAPIKeys, api2.AdminRole))
	v1.DELETE("/key/:id", httpserver.AuthorizeHandler(r.DeleteAPIKey, api2.AdminRole)) // ID likely DB uint ID
	v1.PUT("/key/:id", httpserver.AuthorizeHandler(r.EditAPIKey, api2.AdminRole))      // ID likely DB uint ID

	// Connector Management (protected by Admin role)
	v1.GET("/connectors", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole))
	v1.GET("/connectors/supported-connector-types", httpserver.AuthorizeHandler(r.GetSupportedType, api2.AdminRole))
	// Note: Route path fixed from original code '/connector/:type' assumes type is the unique ID
	v1.GET("/connector/:type", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole)) // Reuses GetConnectors with type filter in handler logic
	v1.POST("/connector", httpserver.AuthorizeHandler(r.CreateConnector, api2.AdminRole))
	v1.POST("/connector/auth0", httpserver.AuthorizeHandler(r.CreateAuth0Connector, api2.AdminRole))
	v1.PUT("/connector", httpserver.AuthorizeHandler(r.UpdateConnector, api2.AdminRole))
	v1.DELETE("/connector/:id", httpserver.AuthorizeHandler(r.DeleteConnector, api2.AdminRole)) // ID is Dex Connector ID string
}

// bindValidate binds and validates the request body.
func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		// Return a more specific error if possible, or log details
		return fmt.Errorf("failed to bind request: %w", err)
	}
	// Use Echo's validator (ensure a validator is registered with Echo instance)
	if err := ctx.Validate(i); err != nil {
		// Return validation errors directly
		return fmt.Errorf("validation failed: %w", err) // Consider returning a 400 Bad Request here
	}
	return nil
}

// Check handles Envoy's Check request.
func (r *httpRoutes) Check(ctx echo.Context) error {
	// Reconstruct CheckRequest from HTTP headers
	checkRequest := envoyauth.CheckRequest{
		Attributes: &envoyauth.AttributeContext{
			Request: &envoyauth.AttributeContext_Request{
				Http: &envoyauth.AttributeContext_HttpRequest{
					Headers: make(map[string]string),
				},
			},
		},
	}

	// Populate headers
	for k, v := range ctx.Request().Header {
		headerKey := strings.ToLower(k) // Envoy usually sends lowercase
		if len(v) > 0 {
			checkRequest.Attributes.Request.Http.Headers[headerKey] = v[0]
		} else {
			checkRequest.Attributes.Request.Http.Headers[headerKey] = ""
		}
	}

	// Get original URI and Method from proxy headers (adjust header names if needed)
	originalURIStr := ctx.Request().Header.Get("x-original-uri")    // Common header name
	originalMethod := ctx.Request().Header.Get("x-original-method") // Common header name

	if originalURIStr != "" {
		originalUri, err := url.Parse(originalURIStr)
		if err != nil {
			r.logger.Warn("Failed to parse X-Original-URI", zap.String("uri", originalURIStr), zap.Error(err))
			checkRequest.Attributes.Request.Http.Path = "/" // Default or error?
		} else {
			checkRequest.Attributes.Request.Http.Path = originalUri.Path
		}
	} else {
		checkRequest.Attributes.Request.Http.Path = ctx.Request().URL.Path
	}

	if originalMethod != "" {
		checkRequest.Attributes.Request.Http.Method = originalMethod
	} else {
		checkRequest.Attributes.Request.Http.Method = ctx.Request().Method
	}
	checkRequest.Attributes.Request.Http.Id = ctx.Request().Header.Get("x-request-id") // Pass request ID if available

	// Call the core Check logic in the server component
	res, err := r.authServer.Check(ctx.Request().Context(), &checkRequest)
	if err != nil {
		r.logger.Error("Auth server Check failed", zap.String("path", checkRequest.Attributes.Request.Http.Path), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal authorization error")
	}

	// Process the CheckResponse from the server
	if res.Status.Code != int32(codes.OK) {
		httpStatusCode := http.StatusUnauthorized // Default
		respBody := "Access Denied"
		if deniedResp := res.GetDeniedResponse(); deniedResp != nil {
			if deniedResp.Status != nil {
				httpStatusCode = int(deniedResp.Status.Code)
			}
			if deniedResp.Body != "" {
				respBody = deniedResp.Body
			}
		}
		r.logger.Info("Access explicitly denied by auth server",
			zap.String("path", checkRequest.Attributes.Request.Http.Path),
			zap.Int32("grpc_code", res.Status.Code),
			zap.String("message", res.Status.Message),
			zap.Int("http_status", httpStatusCode),
		)
		return echo.NewHTTPError(httpStatusCode, respBody)
	}

	// Allowed access
	okResp := res.GetOkResponse()
	if okResp == nil {
		r.logger.Error("Auth server returned OK status but nil OkResponse")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal authorization configuration error")
	}

	// Append headers from OkResponse to the downstream response
	for _, headerOpt := range okResp.GetHeaders() {
		if header := headerOpt.GetHeader(); header != nil {
			ctx.Response().Header().Set(header.Key, header.Value)
		}
	}

	r.logger.Debug("Check request approved, returning OK", zap.String("path", checkRequest.Attributes.Request.Http.Path))
	return ctx.NoContent(http.StatusOK)
}

// Token handles the OAuth/OIDC code exchange.
func (r *httpRoutes) Token(ctx echo.Context) error {
	var req api.GetTokenRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// --- Exchange code with Dex ---
	domain := os.Getenv("DEX_AUTH_DOMAIN")
	if domain == "" {
		r.logger.Error("DEX_AUTH_DOMAIN environment variable not set")
		return echo.NewHTTPError(http.StatusInternalServerError, "Identity provider configuration error")
	}
	dexTokenURL := fmt.Sprintf("%s/token", domain)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", req.Code)
	data.Set("redirect_uri", req.CallBackUrl)
	data.Set("client_id", "public-client") // Assuming public client for code flow
	// data.Set("client_secret", "YOUR_SECRET") // If using confidential client

	r.logger.Info("Exchanging code with Dex", zap.String("url", dexTokenURL), zap.String("clientId", "public-client"))

	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequestWithContext(ctx.Request().Context(), "POST", dexTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		r.logger.Error("Failed to create Dex token request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token request")
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		r.logger.Error("Failed to make Dex token request", zap.String("url", dexTokenURL), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to communicate with identity provider")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		r.logger.Error("Dex token exchange failed",
			zap.String("url", dexTokenURL),
			zap.Int("status", httpResp.StatusCode),
			zap.String("body", string(bodyBytes)))
		return echo.NewHTTPError(http.StatusBadGateway, "Token exchange with identity provider failed")
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&tokenResponse); err != nil {
		r.logger.Error("Failed to decode Dex token response", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process identity provider response")
	}

	// Prefer ID token for OIDC claims
	tokenToVerify := ""
	if idToken, ok := tokenResponse["id_token"].(string); ok && idToken != "" {
		tokenToVerify = idToken
		r.logger.Debug("Using id_token from Dex response for verification")
	} else if accessToken, ok := tokenResponse["access_token"].(string); ok && accessToken != "" {
		// Fallback only if ID token isn't available AND verifier supports access tokens
		tokenToVerify = accessToken
		r.logger.Debug("Using access_token from Dex response for verification (fallback)")
	} else {
		r.logger.Error("Neither id_token nor access_token found in Dex response", zap.Any("responseKeys", mapsKeys(tokenResponse)))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid token response from identity provider")
	}

	// --- Verify Dex Token & Enrich Claims ---
	dv, err := r.authServer.dexVerifier.Verify(ctx.Request().Context(), tokenToVerify)
	if err != nil {
		r.logger.Warn("Failed to verify token from Dex", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to verify token")
	}

	var claims json.RawMessage
	if err := dv.Claims(&claims); err != nil {
		r.logger.Error("Failed to get claims from verified Dex token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process token claims")
	}

	var claimsMap DexClaims // Use the specific struct
	if err = json.Unmarshal(claims, &claimsMap); err != nil {
		r.logger.Error("Failed to unmarshal Dex claims", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process token claims")
	}

	// Find the user in the local DB based on email from verified claims
	user, err := r.db.GetUserByEmail(claimsMap.Email)
	if err != nil || user == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			r.logger.Warn("User from Dex token not found in local DB", zap.String("email", claimsMap.Email))
			return echo.NewHTTPError(http.StatusForbidden, "User not registered in this application")
		}
		r.logger.Error("Failed to get user from local DB during token exchange", zap.String("email", claimsMap.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}

	if !user.IsActive {
		r.logger.Warn("Login attempt via token exchange by inactive user", zap.String("email", user.Email), zap.Uint("userID", user.ID))
		return echo.NewHTTPError(http.StatusForbidden, "User account is inactive")
	}

	// Enrich claims with local data
	enrichedClaims := claimsMap                                              // Start with Dex claims
	enrichedClaims.Groups = append(enrichedClaims.Groups, string(user.Role)) // Add local role
	enrichedClaims.Name = user.Username                                      // Override name with local username
	enrichedClaims.Subject = user.ExternalId                                 // IMPORTANT: Override 'sub' with the canonical ExternalId from DB
	enrichedClaims.Id = uuid.NewString()                                     // Generate a new JTI for this token instance
	enrichedClaims.Issuer = "platform-auth-service"                          // Identify this service as the issuer
	enrichedClaims.Audience = "platform-client"                              // Set appropriate audience(s)
	// enrichedClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() // Set desired expiry

	// --- Create and Sign Enriched Token with KID ---
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, enrichedClaims)

	if r.platformKeyID != "" {
		token.Header["kid"] = r.platformKeyID
		r.logger.Debug("Adding kid to enriched token JWT header", zap.String("kid", r.platformKeyID))
	} else {
		r.logger.Warn("Platform Key ID (kid) is not configured. Enriched token JWT header will not contain kid.")
	}

	signedToken, err := token.SignedString(r.platformPrivateKey)
	if err != nil {
		r.logger.Error("Failed to sign enriched token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to sign token")
	}

	// Replace the original access_token with the new signed token
	tokenResponse["access_token"] = signedToken
	// delete(tokenResponse, "id_token") // Optionally remove original id_token

	r.logger.Info("Successfully exchanged code and issued enriched token", zap.String("email", claimsMap.Email), zap.String("externalId", user.ExternalId))
	return ctx.JSON(http.StatusOK, tokenResponse)
}

// Helper to get map keys for logging without exposing values
func mapsKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// --- GetUsers ---
func (r *httpRoutes) GetUsers(ctx echo.Context) error {
	users, err := r.db.GetUsers()
	if err != nil {
		r.logger.Error("Failed to get users from DB", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	resp := make([]api.GetUsersResponse, 0, len(users))
	for _, u := range users {
		tempResp := api.GetUsersResponse{
			ID:            u.ID,
			UserName:      u.Username,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			ExternalId:    u.ExternalId,
			RoleName:      u.Role,
			CreatedAt:     u.CreatedAt,
			IsActive:      u.IsActive,
			ConnectorId:   u.ConnectorId,
			FullName:      u.FullName,
		}
		if !u.LastLogin.IsZero() {
			tempResp.LastActivity = &u.LastLogin
		}
		resp = append(resp, tempResp)
	}
	return ctx.JSON(http.StatusOK, resp)
}

// --- GetUserDetails ---
func (r *httpRoutes) GetUserDetails(ctx echo.Context) error {
	userIDParam := ctx.Param("id")
	// Assuming ID is the database uint ID based on route and DeleteUser pattern
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		r.logger.Warn("Invalid user ID format in GetUserDetails", zap.String("idParam", userIDParam), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	// DB GetUser expects stringified ID
	user, err := r.db.GetUser(strconv.FormatUint(userID, 10))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user details from DB", zap.Uint64("userID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user details")
	}
	if user == nil { // Should be caught by ErrRecordNotFound
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	resp := api.GetUserResponse{
		ID:            user.ID,
		UserName:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		Blocked:       !user.IsActive, // Blocked is inverse of IsActive
		RoleName:      user.Role,
	}
	if !user.LastLogin.IsZero() {
		resp.LastActivity = &user.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)
}

// --- GetMe ---
func (r *httpRoutes) GetMe(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx) // External ID from token
	if userID == "" {
		r.logger.Error("UserID missing from context in GetMe handler")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in request context")
	}

	dbUser, err := r.db.GetUserByExternalID(userID) // Lookup by External ID
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("User from token context not found in DB", zap.String("externalID", userID))
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get current user from DB", zap.String("externalID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}
	if dbUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Convert db.User to api.GetMeResponse
	resp := api.GetMeResponse{
		ID:            dbUser.ID,
		UserName:      dbUser.Username,
		Email:         dbUser.Email,
		EmailVerified: dbUser.EmailVerified,
		CreatedAt:     dbUser.CreatedAt,
		Blocked:       !dbUser.IsActive,
		Role:          string(dbUser.Role),
		MemberSince:   dbUser.CreatedAt,
		ConnectorId:   dbUser.ConnectorId,
	}
	if !dbUser.LastLogin.IsZero() {
		resp.LastLogin = &dbUser.LastLogin
		resp.LastActivity = &dbUser.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)
}

// --- CreateAPIKey ---
func (r *httpRoutes) CreateAPIKey(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx) // ExternalID of the creator
	var req api.CreateAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Fetch creator details (optional, good for context)
	usr, err := utils.GetUser(userID, r.db)
	if err != nil || usr == nil {
		r.logger.Error("Failed to get creator user details for API key", zap.String("userID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user details")
	}

	// Check key limit (Make limit configurable)
	keyLimit := int64(5) // Example limit
	currentKeyCount, err := r.db.CountApiKeysForUser(userID)
	if err != nil {
		r.logger.Error("Failed to count user API keys", zap.String("userID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to count API keys")
	}
	if currentKeyCount >= keyLimit {
		r.logger.Warn("API key limit reached for user", zap.String("userID", userID), zap.Int64("limit", keyLimit))
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Maximum number of %d API keys for user reached", keyLimit))
	}

	jti := uuid.NewString()
	apiKeyClaims := &jwt.StandardClaims{
		Issuer:    "platform-auth-service",
		Subject:   userID, // Subject is the ExternalID of the creator
		Audience:  "platform-api",
		ExpiresAt: 0, // No expiry, or use time.Now().AddDate(1, 0, 0).Unix() for 1 year
		IssuedAt:  jwt.TimeFunc().Unix(),
		Id:        jti,
	}
	// Add role as a custom claim if services expect it directly in the token:
	// type ApiKeyClaims struct {
	// 	Role api.Role `json:"role"`
	// 	jwt.StandardClaims
	// }
	// claims := ApiKeyClaims{ Role: req.Role, StandardClaims: *apiKeyClaims }
	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	if r.platformPrivateKey == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Platform API key signing is disabled")
	}

	// Create Token Object and Add KID
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, apiKeyClaims) // Use StandardClaims for now
	if r.platformKeyID != "" {
		token.Header["kid"] = r.platformKeyID
		r.logger.Debug("Adding kid to API Key JWT header", zap.String("kid", r.platformKeyID))
	} else {
		r.logger.Warn("Platform Key ID (kid) is not configured. API Key JWT header will not contain kid.")
	}

	// Sign the API Key token
	signedToken, err := token.SignedString(r.platformPrivateKey)
	if err != nil {
		r.logger.Error("Failed to sign API key token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create API key token")
	}

	// Hashing and DB Storage
	masked := fmt.Sprintf("%s...%s", signedToken[:min(10, len(signedToken))], signedToken[max(0, len(signedToken)-10):])
	hash := sha512.New()
	_, err = hash.Write([]byte(signedToken))
	if err != nil {
		r.logger.Error("Failed to hash API key token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash API key")
	}
	keyHash := hex.EncodeToString(hash.Sum(nil))

	r.logger.Info("Creating API Key DB entry", zap.String("name", req.Name), zap.String("role", string(req.Role)))
	apikey := db.ApiKey{
		Name:          req.Name,
		Role:          req.Role, // Store the intended role from the request
		CreatorUserID: userID,   // Store the ExternalID of the creator
		IsActive:      true,     // New keys are active by default
		MaskedKey:     masked,
		KeyHash:       keyHash,
	}

	err = r.db.AddApiKey(&apikey)
	if err != nil {
		r.logger.Error("Failed to add API Key to db", zap.Error(err))
		// TODO: Check for potential duplicate name errors if there's a unique constraint
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save API key")
	}
	r.logger.Info("Successfully created and stored API Key", zap.Uint("apiKeyID", apikey.ID), zap.String("name", apikey.Name))

	// Return the response including the actual token (ONLY ON CREATION)
	return ctx.JSON(http.StatusCreated, api.CreateAPIKeyResponse{
		ID:        apikey.ID,
		Name:      apikey.Name,
		Active:    apikey.IsActive,
		CreatedAt: apikey.CreatedAt,
		RoleName:  apikey.Role,
		Token:     signedToken, // Return the full token
	})
}

// DeleteAPIKey deletes an API key by its database ID.
func (r *httpRoutes) DeleteAPIKey(ctx echo.Context) error {
	idParam := ctx.Param("id")
	apiKeyID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		r.logger.Warn("Invalid API Key ID format in DeleteAPIKey", zap.String("idParam", idParam), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid API Key ID format")
	}

	err = r.db.DeleteAPIKey(apiKeyID)
	if err != nil {
		// GORM might not return ErrRecordNotFound on Delete if nothing matched.
		// Check affected rows if needed, or just log and return error.
		r.logger.Error("Failed to delete API Key", zap.Uint64("apiKeyID", apiKeyID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete API Key")
	}

	r.logger.Info("Deleted API Key", zap.Uint64("apiKeyID", apiKeyID))
	return ctx.NoContent(http.StatusAccepted) // 202 or 204
}

// EditAPIKey updates the role and/or active status of an API key.
func (r *httpRoutes) EditAPIKey(ctx echo.Context) error {
	idParam := ctx.Param("id")
	// DB UpdateAPIKey expects string ID
	// apiKeyID, err := strconv.ParseUint(idParam, 10, 64) ... if conversion needed

	var req api.EditAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := r.db.UpdateAPIKey(idParam, req.IsActive, req.Role)
	if err != nil {
		// GORM Update might return ErrRecordNotFound if key doesn't exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "API Key not found")
		}
		r.logger.Error("Failed to update API Key", zap.String("apiKeyID", idParam), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update API Key")
	}

	r.logger.Info("Updated API Key", zap.String("apiKeyID", idParam), zap.Bool("isActive", req.IsActive), zap.String("role", string(req.Role)))
	return ctx.NoContent(http.StatusAccepted) // 202 or 200 OK
}

// ListAPIKeys lists API keys. If called by Admin, consider listing all keys?
// Currently lists keys created *by* the requesting user (who must be Admin based on route).
func (r *httpRoutes) ListAPIKeys(ctx echo.Context) error {
	requestingUserID := httpserver.GetUserID(ctx) // ExternalID of the admin

	// Decide whether to list all keys or just those created by this admin
	var keys []db.ApiKey
	var err error
	// Example: If you want Admins to see *all* keys, use ListApiKeys() instead
	// keys, err = r.db.ListApiKeys()
	keys, err = r.db.ListApiKeysForUser(requestingUserID) // Lists keys created by this user

	if err != nil {
		r.logger.Error("Failed to list API keys for user", zap.String("userID", requestingUserID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve API keys")
	}

	resp := make([]api.APIKeyResponse, 0, len(keys))
	for _, key := range keys {
		resp = append(resp, api.APIKeyResponse{
			ID:            key.ID,
			CreatedAt:     key.CreatedAt,
			UpdatedAt:     key.UpdatedAt,
			Name:          key.Name,
			RoleName:      key.Role,
			CreatorUserID: key.CreatorUserID,
			Active:        key.IsActive,
			MaskedKey:     key.MaskedKey,
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

// --- CreateUser ---
func (r *httpRoutes) CreateUser(ctx echo.Context) error {
	var req api.CreateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Call internal function to handle user creation logic
	err := r.DoCreateUser(req)
	if err != nil {
		// DoCreateUser should return appropriate echo.HTTPError
		return err
	}

	r.logger.Info("User created successfully", zap.String("email", req.EmailAddress))
	return ctx.NoContent(http.StatusCreated)
}

// DoCreateUser contains the core logic for creating a user.
func (r *httpRoutes) DoCreateUser(req api.CreateUserRequest) error {
	if req.EmailAddress == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is required")
	}
	email := strings.ToLower(strings.TrimSpace(req.EmailAddress)) // Normalize email

	// Check if user already exists
	existingUser, err := r.db.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Handle actual DB errors
		r.logger.Error("Failed to check existing user by email", zap.String("email", email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check user existence")
	}
	if existingUser != nil {
		r.logger.Warn("Attempt to create user with existing email", zap.String("email", email))
		return echo.NewHTTPError(http.StatusConflict, "Email address already in use") // 409 Conflict
	}

	// Determine if this is the first user (potential admin bootstrap)
	count, err := r.db.GetUsersCount()
	if err != nil {
		r.logger.Error("Failed to get users count", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get users count")
	}

	isAdminBootstrap := count == 0
	role := api2.ViewerRole // Default role
	if req.Role != nil {
		role = *req.Role
	}

	if isAdminBootstrap {
		r.logger.Info("Creating first user, assigning Admin role", zap.String("email", email))
		role = api2.AdminRole // Force first user to be admin
	} else if req.Role == nil {
		r.logger.Info("No role specified for new user, defaulting to Viewer", zap.String("email", email))
	}

	connectorType := req.ConnectorId // Use provided connector ID/type
	externalID := ""
	requirePasswordChange := true // Default for non-admin, non-bootstrap

	if req.Password != nil && *req.Password != "" {
		// --- Handle Local User with Password ---
		connectorType = "local" // Override connector type for password users
		externalID = fmt.Sprintf("local|%s", email)
		r.logger.Info("Creating local user with password", zap.String("email", email), zap.String("externalId", externalID))

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("Failed to hash user password", zap.String("email", email), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process password")
		}

		// Create Dex password entry
		dexReq := &dexApi.CreatePasswordReq{
			Password: &dexApi.Password{
				UserId:   externalID, // Use the generated external ID
				Email:    email,
				Hash:     hashedPassword,
				Username: email, // Use email as username for Dex password DB
			},
		}
		resp, err := r.authServer.dexClient.CreatePassword(context.TODO(), dexReq) // Use request context if available
		if err != nil {
			// Check if gRPC error indicates specific Dex issue
			r.logger.Error("Failed to create dex password", zap.String("email", email), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadGateway, "Failed to create user identity with provider")
		}
		if resp.AlreadyExists {
			// This shouldn't happen if email check passed, but handle defensively
			r.logger.Warn("Dex password entry already exists for new user, attempting update", zap.String("email", email))
			updateReq := &dexApi.UpdatePasswordReq{
				Email:       email,
				NewHash:     hashedPassword,
				NewUsername: email,
			}
			_, err = r.authServer.dexClient.UpdatePassword(context.TODO(), updateReq)
			if err != nil {
				r.logger.Error("Failed to update potentially existing dex password", zap.String("email", email), zap.Error(err))
				return echo.NewHTTPError(http.StatusBadGateway, "Failed to update user identity with provider")
			}
		}
		// Password was just set, don't require change immediately unless specifically requested
		requirePasswordChange = false

	} else {
		// --- Handle User from External Connector (or local without initial password) ---
		if connectorType == "" || connectorType == "local" {
			connectorType = "local" // Default to local if not specified or password not set
			externalID = fmt.Sprintf("local|%s", email)
			r.logger.Info("Creating local user without initial password", zap.String("email", email))
			// They might need a password reset flow later
		} else {
			// Assume external connector
			externalID = fmt.Sprintf("%s|%s", connectorType, email) // Construct external ID
			r.logger.Info("Creating user linked to external connector", zap.String("email", email), zap.String("connector", connectorType))
			requirePasswordChange = false // Password managed by external IDP
		}
	}

	if isAdminBootstrap {
		requirePasswordChange = false // First admin doesn't need to change password
	}

	// Create user in local DB
	newUser := &db.User{
		Email:                 email,
		Username:              email, // Default username to email
		FullName:              email, // Default full name to email
		Role:                  role,
		EmailVerified:         false, // Assume not verified initially
		ConnectorId:           connectorType,
		ExternalId:            externalID,
		RequirePasswordChange: requirePasswordChange,
		IsActive:              true, // New users active by default (req.IsActive was unused?)
	}
	err = r.db.CreateUser(newUser)
	if err != nil {
		r.logger.Error("Failed to create user in local database", zap.String("email", email), zap.Error(err))
		// Attempt to clean up Dex password entry if DB insert failed? Complex rollback.
		// dexClient.DeletePassword(...)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user data")
	}

	return nil // Success
}

// --- UpdateUser ---
func (r *httpRoutes) UpdateUser(ctx echo.Context) error {
	var req api.UpdateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.EmailAddress == "" { // Email address is used to identify the user to update
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is required to identify user")
	}
	email := strings.ToLower(strings.TrimSpace(req.EmailAddress))

	// Get the existing user
	user, err := r.db.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for update", zap.String("email", email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user")
	}
	if user == nil { // Safeguard
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// --- Update Password (if provided and user is local) ---
	if req.Password != nil && *req.Password != "" {
		if user.ConnectorId != "local" {
			r.logger.Warn("Attempt to set password for non-local user", zap.String("email", email), zap.String("connector", user.ConnectorId))
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot set password for user linked to external connector")
		}

		r.logger.Info("Updating password for local user", zap.String("email", email))
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("Failed to hash updated password", zap.String("email", email), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process password")
		}

		// Update Dex password entry
		dexUpdateReq := &dexApi.UpdatePasswordReq{
			Email:       email,
			NewHash:     hashedPassword,
			NewUsername: email, // Keep username consistent
		}
		resp, err := r.authServer.dexClient.UpdatePassword(context.TODO(), dexUpdateReq)
		if err != nil {
			r.logger.Error("Failed to update dex password", zap.String("email", email), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadGateway, "Failed to update user identity with provider")
		}
		if resp.NotFound {
			// If Dex entry was missing, try creating it
			r.logger.Warn("Dex password entry not found during update, attempting creation", zap.String("email", email))
			dexCreateReq := &dexApi.CreatePasswordReq{
				Password: &dexApi.Password{
					UserId:   user.ExternalId, // Use existing external ID
					Email:    email,
					Hash:     hashedPassword,
					Username: email,
				},
			}
			_, err = r.authServer.dexClient.CreatePassword(context.TODO(), dexCreateReq)
			if err != nil {
				r.logger.Error("Failed to create dex password during update fallback", zap.String("email", email), zap.Error(err))
				return echo.NewHTTPError(http.StatusBadGateway, "Failed to create user identity with provider")
			}
		}
		// Mark password change requirement as false since admin/update set it
		err = r.db.UserPasswordUpdate(user.ID)
		if err != nil {
			r.logger.Error("Failed to mark user password as updated in db", zap.Uint("userID", user.ID), zap.Error(err))
			// Non-fatal? Log and continue with other updates.
		}
	}

	// --- Update Other User Fields ---
	// Build update struct with only fields that are present in the request
	// Using GORM's Updates method handles partial updates well if fields are pointers
	// But the request struct uses value types for IsActive, UserName etc.
	// So, we update the user object and then save it.

	updateNeeded := false
	if req.Role != nil && user.Role != *req.Role {
		user.Role = *req.Role
		updateNeeded = true
	}
	if user.IsActive != req.IsActive {
		user.IsActive = req.IsActive
		updateNeeded = true
	}
	if req.UserName != "" && user.Username != req.UserName {
		user.Username = req.UserName
		updateNeeded = true
	}
	if req.FullName != "" && user.FullName != req.FullName {
		user.FullName = req.FullName
		updateNeeded = true
	}
	// Update ConnectorId and ExternalId if connector changes
	// Be careful: Changing ConnectorId for an existing user can break login if not handled properly.
	if req.ConnectorId != "" && user.ConnectorId != req.ConnectorId {
		user.ConnectorId = req.ConnectorId
		user.ExternalId = fmt.Sprintf("%s|%s", req.ConnectorId, user.Email) // Recalculate ExternalId
		// If changing to 'local', should we attempt to delete Dex password?
		// If changing from 'local', should we attempt to delete Dex password?
		r.logger.Warn("User connector changed", zap.String("email", email), zap.String("oldConnector", user.ConnectorId), zap.String("newConnector", req.ConnectorId))
		updateNeeded = true
	}

	if updateNeeded {
		r.logger.Info("Updating user details in database", zap.String("email", email))
		err = r.db.UpdateUser(user) // UpdateUser should save all fields passed in the user struct
		if err != nil {
			r.logger.Error("Failed to update user in database", zap.String("email", email), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user data")
		}
	} else {
		r.logger.Info("No user details needed updating", zap.String("email", email))
	}

	return ctx.NoContent(http.StatusOK)
}

// --- DeleteUser ---
func (r *httpRoutes) DeleteUser(ctx echo.Context) error {
	idParam := ctx.Param("id")
	// Assuming ID is the database uint ID
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		r.logger.Warn("Invalid user ID format in DeleteUser", zap.String("idParam", idParam), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	err = r.DoDeleteUser(uint(userID)) // Pass uint ID to internal logic
	if err != nil {
		return err // DoDeleteUser should return appropriate echo.HTTPError
	}

	r.logger.Info("User deleted successfully", zap.Uint("userID", uint(userID)))
	return ctx.NoContent(http.StatusAccepted) // 202 or 204
}

// DoDeleteUser contains the core logic for deleting a user.
func (r *httpRoutes) DoDeleteUser(userID uint) error {
	// Get user details first to find email for Dex password deletion
	user, err := r.db.GetUser(strconv.FormatUint(uint64(userID), 10)) // GetUser takes string ID
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for deletion", zap.Uint("userID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user before deletion")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Prevent deletion of the very first user (ID=1, common convention)
	if user.ID == 1 {
		r.logger.Warn("Attempt to delete the first user (ID 1)", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot delete the initial administrator user")
	}

	// If the user is a local user, attempt to delete their Dex password entry
	if user.ConnectorId == "local" {
		r.logger.Info("Deleting Dex password entry for local user", zap.String("email", user.Email))
		dexReq := &dexApi.DeletePasswordReq{
			Email: user.Email,
		}
		resp, err := r.authServer.dexClient.DeletePassword(context.TODO(), dexReq)
		if err != nil {
			// Log error but proceed with local deletion? Or fail? Fail seems safer.
			r.logger.Error("Failed to remove dex password during user deletion", zap.String("email", user.Email), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadGateway, "Failed to remove user identity from provider")
		}
		if resp.NotFound {
			r.logger.Warn("Dex password entry not found during deletion, proceeding", zap.String("email", user.Email))
		}
	}

	// Delete user from local database
	err = r.db.DeleteUser(user.ID) // DeleteUser takes uint ID
	if err != nil {
		// Should we attempt rollback on Dex if this fails? Complex.
		r.logger.Error("Failed to delete user from local database", zap.Uint("userID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user data")
	}

	return nil // Success
}

// --- CheckUserPasswordChangeRequired ---
func (r *httpRoutes) CheckUserPasswordChangeRequired(ctx echo.Context) error {
	userId := httpserver.GetUserID(ctx) // External ID from token
	if userId == "" {
		r.logger.Error("UserID missing from context in CheckUserPasswordChangeRequired")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in request context")
	}

	user, err := r.db.GetUserByExternalID(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for password check", zap.String("externalID", userId), zap.Error(err))
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

// --- ResetUserPassword ---
func (r *httpRoutes) ResetUserPassword(ctx echo.Context) error {
	userId := httpserver.GetUserID(ctx) // External ID of the user resetting their own password
	if userId == "" {
		r.logger.Error("UserID missing from context in ResetUserPassword")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in request context")
	}

	// Get user details
	user, err := r.db.GetUserByExternalID(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		r.logger.Error("Failed to get user for password reset", zap.String("externalID", userId), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Check if user is local
	if user.ConnectorId != "local" {
		r.logger.Warn("Password reset attempt for non-local user", zap.String("externalID", userId), zap.String("connector", user.ConnectorId))
		return echo.NewHTTPError(http.StatusBadRequest, "Password reset only available for local accounts")
	}

	// Bind and validate request body
	var req api.ResetUserPasswordRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Current and new passwords are required")
	}
	if req.CurrentPassword == req.NewPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password must be different from the current password")
	}
	// Add password complexity rules if needed

	// Verify current password with Dex
	dexVerifyReq := &dexApi.VerifyPasswordReq{
		Email:    user.Email,
		Password: req.CurrentPassword,
	}
	resp, err := r.authServer.dexClient.VerifyPassword(context.TODO(), dexVerifyReq)
	if err != nil {
		r.logger.Error("Failed to verify current password with Dex", zap.String("email", user.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to verify current password")
	}
	if resp.NotFound {
		// Should not happen if user exists locally and is local type
		r.logger.Error("Dex password entry not found for existing local user during verification", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusInternalServerError, "User identity inconsistency")
	}
	if !resp.Verified {
		r.logger.Info("Incorrect current password provided during reset", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect current password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash new password", zap.String("email", user.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process new password")
	}

	// Update password in Dex
	passwordUpdateReq := &dexApi.UpdatePasswordReq{
		Email:       user.Email,
		NewHash:     hashedPassword,
		NewUsername: user.Username, // Keep username consistent
	}
	passwordUpdateResp, err := r.authServer.dexClient.UpdatePassword(context.TODO(), passwordUpdateReq)
	if err != nil {
		r.logger.Error("Failed to update dex password during reset", zap.String("email", user.Email), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to update password with identity provider")
	}
	if passwordUpdateResp.NotFound {
		// Should not happen if verification succeeded, but handle defensively
		r.logger.Error("Dex password entry not found for existing local user during update", zap.String("email", user.Email))
		return echo.NewHTTPError(http.StatusInternalServerError, "User identity inconsistency")
	}

	// Mark password change as completed in local DB
	err = r.db.UserPasswordUpdate(user.ID)
	if err != nil {
		r.logger.Error("Failed to mark user password as updated in db after reset", zap.Uint("userID", user.ID), zap.Error(err))
		// Non-fatal? Log and still return success as password was updated in Dex.
	}

	r.logger.Info("User successfully reset password", zap.String("email", user.Email))
	return ctx.NoContent(http.StatusAccepted) // 202 Accepted or 200 OK
}

// --- GetConnectors ---
func (r *httpRoutes) GetConnectors(ctx echo.Context) error {
	// Check for optional type filter from path parameter
	connectorTypeFilter := ctx.Param("type")
	if connectorTypeFilter != "" {
		r.logger.Info("Filtering connectors by type", zap.String("type", connectorTypeFilter))
	}

	// List connectors from Dex
	req := &dexApi.ListConnectorReq{}
	respDex, err := r.authServer.dexClient.ListConnectors(context.TODO(), req)
	if err != nil {
		r.logger.Error("Failed to list connectors from Dex", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to retrieve connector list from provider")
	}

	connectorsFromDex := respDex.Connectors
	resp := make([]api.GetConnectorsResponse, 0, len(connectorsFromDex))

	for _, dexConnector := range connectorsFromDex {
		// Skip the built-in 'local' connector unless specifically requested?
		if dexConnector.Id == "local" {
			continue
		}

		// Apply type filter if provided
		if connectorTypeFilter != "" && !strings.EqualFold(connectorTypeFilter, dexConnector.Type) {
			continue
		}

		// Get corresponding local DB record for additional metadata
		localConnector, err := r.db.GetConnectorByConnectorID(dexConnector.Id)
		if err != nil {
			// Log error but potentially continue? Or fail?
			// If local record is missing, maybe still show Dex info?
			r.logger.Warn("Failed to get local DB record for Dex connector", zap.String("connectorID", dexConnector.Id), zap.Error(err))
			// Decide how to handle - skip, show partial, or error out? Skipping for now.
			continue
			// If DB error is critical:
			// return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve connector metadata")
		}
		if localConnector == nil {
			// Dex has connector, but we don't. Log inconsistency.
			r.logger.Warn("Connector exists in Dex but not in local DB", zap.String("connectorID", dexConnector.Id))
			continue // Skip inconsistent entries
		}

		// Build response item
		info := api.GetConnectorsResponse{
			ID:          localConnector.ID, // Local DB ID
			ConnectorID: dexConnector.Id,   // Dex ID
			Type:        dexConnector.Type,
			Name:        dexConnector.Name,
			SubType:     localConnector.ConnectorSubType, // From local DB
			UserCount:   localConnector.UserCount,        // From local DB (how is this updated?)
			CreatedAt:   localConnector.CreatedAt,        // From local DB
			LastUpdate:  localConnector.LastUpdate,       // From local DB
		}

		// Extract common OIDC config fields if applicable
		if strings.EqualFold(dexConnector.Type, "oidc") && len(dexConnector.Config) > 0 {
			var oidcConfig struct { // Unmarshal only needed fields
				Issuer   string `json:"issuer"`
				ClientID string `json:"clientID"`
				// TenantID string `json:"tenantID"` // Only for specific subtypes
			}
			err := json.Unmarshal(dexConnector.Config, &oidcConfig)
			if err != nil {
				r.logger.Warn("Failed to unmarshal OIDC config for connector", zap.String("connectorID", dexConnector.Id), zap.Error(err))
			} else {
				info.Issuer = oidcConfig.Issuer
				info.ClientID = oidcConfig.ClientID
				// Note: Omitting ClientSecret for security
			}
		}
		resp = append(resp, info)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// --- GetSupportedType ---
func (r *httpRoutes) GetSupportedType(ctx echo.Context) error {
	// Use maps defined in utils package
	supportedConnectors := utils.SupportedConnectors
	supportedNames := utils.SupportedConnectorsNames

	responseList := make([]api.GetSupportedConnectorTypeResponse, 0, len(supportedConnectors))

	// Currently only supports 'oidc', make dynamic if more types added
	if subTypes, ok := supportedConnectors["oidc"]; ok {
		subTypeNames := supportedNames["oidc"]
		if len(subTypes) != len(subTypeNames) {
			r.logger.Error("Mismatch between supported OIDC subtypes and names in utils config")
			// Handle error appropriately
		}

		apiSubTypes := make([]api.ConnectorSubTypes, 0, len(subTypes))
		for i, subTypeID := range subTypes {
			name := subTypeID // Default name to ID
			if i < len(subTypeNames) {
				name = subTypeNames[i]
			}
			apiSubTypes = append(apiSubTypes, api.ConnectorSubTypes{
				ID:   subTypeID,
				Name: name,
			})
		}
		responseList = append(responseList, api.GetSupportedConnectorTypeResponse{
			ConnectorType: "oidc",
			SubTypes:      apiSubTypes,
		})
	}
	// Add loops for other connector types (e.g., "saml") if supported in the future

	return ctx.JSON(http.StatusOK, responseList)
}

// --- CreateConnector ---
func (r *httpRoutes) CreateConnector(ctx echo.Context) error {
	var req api.CreateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Prepare request for utility function (validation happens there too)
	dexUtilReq := utils.CreateConnectorRequest{
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		ID:               req.ID, // Dex Connector ID
		Name:             req.Name,
	}

	// Use utility function to build Dex API request
	dexAPICreator := utils.GetConnectorCreator(strings.ToLower(dexUtilReq.ConnectorType))
	if dexAPICreator == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported connector type: %s", dexUtilReq.ConnectorType))
	}

	dexAPIReq, err := dexAPICreator(dexUtilReq)
	if err != nil {
		// Error from utility function (e.g., validation failure, unsupported subtype)
		r.logger.Warn("Failed to prepare Dex connector creation request", zap.Error(err), zap.Any("request", req))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to create connector request: %s", err.Error()))
	}

	// Call Dex gRPC API
	r.logger.Info("Creating Dex connector", zap.String("id", dexAPIReq.Connector.Id), zap.String("type", dexAPIReq.Connector.Type), zap.String("name", dexAPIReq.Connector.Name))
	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexAPIReq)
	if err != nil {
		r.logger.Error("Failed to create Dex connector via gRPC", zap.String("id", dexAPIReq.Connector.Id), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to create connector with identity provider")
	}
	if res.AlreadyExists {
		r.logger.Warn("Attempt to create Dex connector that already exists", zap.String("id", dexAPIReq.Connector.Id))
		return echo.NewHTTPError(http.StatusConflict, "Connector with this ID already exists")
	}

	// Create corresponding record in local DB
	// Use the ID from the Dex request (which incorporated defaults if needed)
	localConnector := &db.Connector{
		ConnectorID:      dexAPIReq.Connector.Id,
		ConnectorType:    dexAPIReq.Connector.Type,
		ConnectorSubType: req.ConnectorSubType, // Store the subtype provided
		LastUpdate:       time.Now(),
		// UserCount starts at 0
	}
	err = r.db.CreateConnector(localConnector)
	if err != nil {
		r.logger.Error("Failed to create local DB record for new connector", zap.String("connectorID", dexAPIReq.Connector.Id), zap.Error(err))
		// Critical inconsistency: Connector exists in Dex but not locally. Manual intervention needed?
		// Should we try to delete the Dex connector? Rollback is complex.
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save connector metadata locally after creation")
	}

	// Restart Dex Pod (handle with care)
	r.logger.Info("Restarting Dex pod to apply connector changes", zap.String("connectorID", dexAPIReq.Connector.Id))
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after connector creation", zap.String("connectorID", dexAPIReq.Connector.Id), zap.Error(err))
		// Non-fatal? Connector is created but might not be active in Dex yet.
		// Return success but maybe include a warning?
		// return echo.NewHTTPError(http.StatusInternalServerError, "Failed to restart identity provider service")
	}

	r.logger.Info("Successfully created connector", zap.String("id", dexAPIReq.Connector.Id))
	// Return success, maybe return the created connector details?
	// The Dex response 'res' is empty on success, so return based on localConnector.
	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":           localConnector.ID,
		"connector_id": localConnector.ConnectorID,
		"type":         localConnector.ConnectorType,
		"sub_type":     localConnector.ConnectorSubType,
		"created_at":   localConnector.CreatedAt,
	})
}

// --- CreateAuth0Connector --- (Specialized endpoint)
func (r *httpRoutes) CreateAuth0Connector(ctx echo.Context) error {
	var req api.CreateAuth0ConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Use utility function to build Dex request for Auth0 connector
	dexUtilReq := utils.CreateAuth0ConnectorRequest{
		Issuer:       req.Issuer, // Issuer might be derived from Domain in util if empty
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Domain:       req.Domain, // Used by util to construct RedirectURIs etc.
	}
	dexAPIReq, err := utils.CreateAuth0Connector(dexUtilReq)
	if err != nil {
		r.logger.Warn("Failed to prepare Dex Auth0 connector creation request", zap.Error(err), zap.Any("request", req))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to create Auth0 connector request: %s", err.Error()))
	}

	// --- Ensure Dex Clients are Updated/Created (moved from original CreateAuth0Connector util?) ---
	// This logic seems more appropriate here in the handler than inside the connector creation util.
	// It updates Dex's *OAuth Clients*, not the connector itself.
	publicUris := req.PublicURIS // Use URIs from the request
	if len(publicUris) > 0 {
		err = r.ensureDexClient("public-client", "Public Client", publicUris, true) // true for public
		if err != nil {
			return err
		} // ensureDexClient should return echo.HTTPError
	}
	privateUris := req.PrivateURIS // Use URIs from the request
	if len(privateUris) > 0 {
		// Assuming private client needs a fixed secret "secret" based on original code
		err = r.ensureDexClient("private-client", "Private Client", privateUris, false) // false for confidential
		if err != nil {
			return err
		}
	}
	// --- End Dex Client Update ---

	// Call Dex gRPC API to create the *connector*
	r.logger.Info("Creating Dex Auth0 connector", zap.String("id", dexAPIReq.Connector.Id), zap.String("name", dexAPIReq.Connector.Name))
	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexAPIReq)
	if err != nil {
		r.logger.Error("Failed to create Dex Auth0 connector via gRPC", zap.String("id", dexAPIReq.Connector.Id), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to create connector with identity provider")
	}
	if res.AlreadyExists {
		r.logger.Warn("Attempt to create Dex Auth0 connector that already exists", zap.String("id", dexAPIReq.Connector.Id))
		// Should we update instead? Or just return conflict?
		return echo.NewHTTPError(http.StatusConflict, "Connector with ID 'auth0' already exists")
	}

	// Create corresponding record in local DB
	localConnector := &db.Connector{
		ConnectorID:      dexAPIReq.Connector.Id, // Should be "auth0"
		ConnectorType:    "oidc",
		ConnectorSubType: "auth0", // Specific subtype
		LastUpdate:       time.Now(),
	}
	err = r.db.CreateConnector(localConnector)
	if err != nil {
		r.logger.Error("Failed to create local DB record for Auth0 connector", zap.String("connectorID", dexAPIReq.Connector.Id), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save connector metadata locally after creation")
	}

	// Restart Dex Pod
	r.logger.Info("Restarting Dex pod to apply Auth0 connector changes", zap.String("connectorID", dexAPIReq.Connector.Id))
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after Auth0 connector creation", zap.Error(err))
		// Log and continue, maybe return warning?
	}

	r.logger.Info("Successfully created Auth0 connector", zap.String("id", dexAPIReq.Connector.Id))
	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":           localConnector.ID,
		"connector_id": localConnector.ConnectorID,
		"type":         localConnector.ConnectorType,
		"sub_type":     localConnector.ConnectorSubType,
		"created_at":   localConnector.CreatedAt,
	})
}

// Helper for ensuring Dex OAuth client exists/is updated
func (r *httpRoutes) ensureDexClient(id, name string, redirectUris []string, isPublic bool) error {
	ctx := context.TODO() // Use request context if available
	clientResp, _ := r.authServer.dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: id})

	if clientResp != nil && clientResp.Client != nil {
		// Update existing client
		r.logger.Info("Updating Dex OAuth client", zap.String("id", id), zap.Strings("redirectUris", redirectUris))
		updateReq := dexApi.UpdateClientReq{
			Id:           id,
			Name:         name, // Ensure name is updated too if changed
			RedirectUris: redirectUris,
		}
		_, err := r.authServer.dexClient.UpdateClient(ctx, &updateReq)
		if err != nil {
			r.logger.Error("Failed to update Dex OAuth client", zap.String("id", id), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("Failed to update OAuth client '%s'", id))
		}
	} else {
		// Create new client
		r.logger.Info("Creating Dex OAuth client", zap.String("id", id), zap.Strings("redirectUris", redirectUris), zap.Bool("public", isPublic))
		createReq := dexApi.CreateClientReq{
			Client: &dexApi.Client{
				Id:           id,
				Name:         name,
				RedirectUris: redirectUris,
				Public:       isPublic,
				// Secret is only set for non-public clients
			},
		}
		if !isPublic {
			createReq.Client.Secret = "secret" // Use configurable secret
		}
		_, err := r.authServer.dexClient.CreateClient(ctx, &createReq)
		if err != nil {
			r.logger.Error("Failed to create Dex OAuth client", zap.String("id", id), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("Failed to create OAuth client '%s'", id))
		}
	}
	return nil
}

// --- UpdateConnector ---
func (r *httpRoutes) UpdateConnector(ctx echo.Context) error {
	var req api.UpdateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// ID in request is the local DB uint ID. We need Dex ConnectorID for Dex update.
	if req.ID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Local database connector ID is required for update")
	}
	if req.ConnectorID == "" {
		// We could fetch the local record using req.ID to get the ConnectorID,
		// but the request should ideally include it for clarity.
		return echo.NewHTTPError(http.StatusBadRequest, "Dex connector_id is required in the request body for update")
	}

	// Prepare request for utility function
	dexUtilReq := utils.UpdateConnectorRequest{
		ID:               req.ConnectorID, // Use ConnectorID from request for Dex target
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		// Name update is not directly supported by Dex UpdateConnector (only config)
	}

	// Use utility function to build Dex API update request
	// Assuming UpdateOIDCConnector exists and works similarly to Create
	dexAPIReq, err := utils.UpdateOIDCConnector(dexUtilReq)
	if err != nil {
		r.logger.Warn("Failed to prepare Dex connector update request", zap.Error(err), zap.Any("request", req))
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to create connector update request: %s", err.Error()))
	}

	// Call Dex gRPC API
	r.logger.Info("Updating Dex connector", zap.String("id", dexAPIReq.Id))
	res, err := r.authServer.dexClient.UpdateConnector(context.TODO(), dexAPIReq)
	if err != nil {
		r.logger.Error("Failed to update Dex connector via gRPC", zap.String("id", dexAPIReq.Id), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to update connector with identity provider")
	}
	if res.NotFound {
		r.logger.Warn("Attempt to update Dex connector that does not exist", zap.String("id", dexAPIReq.Id))
		return echo.NewHTTPError(http.StatusNotFound, "Connector not found in identity provider")
	}

	// Update corresponding record in local DB
	localConnectorUpdate := &db.Connector{
		Model: gorm.Model{ID: req.ID}, // Identify record by local DB ID
		// Update fields based on request - be careful about overwriting existing data
		ConnectorID:      req.ConnectorID, // Ensure this matches Dex ID
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		LastUpdate:       time.Now(),
		// Name is not updated in Dex, should we update it locally? Based on req.Name?
		// UserCount should not be reset here.
	}
	err = r.db.UpdateConnector(localConnectorUpdate) // Assumes UpdateConnector updates non-zero fields
	if err != nil {
		r.logger.Error("Failed to update local DB record for connector", zap.Uint("localID", req.ID), zap.String("connectorID", req.ConnectorID), zap.Error(err))
		// Inconsistency: Dex updated, local DB failed.
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save connector metadata locally after update")
	}

	// Restart Dex Pod (handle with care)
	r.logger.Info("Restarting Dex pod to apply connector changes", zap.String("connectorID", req.ConnectorID))
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after connector update", zap.String("connectorID", req.ConnectorID), zap.Error(err))
		// Non-fatal?
	}

	r.logger.Info("Successfully updated connector", zap.String("id", req.ConnectorID))
	return ctx.NoContent(http.StatusAccepted) // 202 or 200
}

// --- DeleteConnector ---
func (r *httpRoutes) DeleteConnector(ctx echo.Context) error {
	// ID in path parameter is the Dex Connector ID string (e.g., "oidc-google", "auth0")
	connectorID := ctx.Param("id")
	if connectorID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Connector ID is required in path")
	}
	if connectorID == "local" {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot delete the built-in 'local' connector")
	}

	// Call Dex gRPC API to delete connector
	r.logger.Info("Deleting Dex connector", zap.String("id", connectorID))
	dexReq := &dexApi.DeleteConnectorReq{
		Id: connectorID,
	}
	resp, err := r.authServer.dexClient.DeleteConnector(context.TODO(), dexReq)
	if err != nil {
		r.logger.Error("Failed to delete Dex connector via gRPC", zap.String("id", connectorID), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to delete connector with identity provider")
	}
	if resp.NotFound {
		r.logger.Warn("Attempt to delete Dex connector that does not exist", zap.String("id", connectorID))
		// Connector doesn't exist in Dex, should we still delete local record? Yes.
	}

	// Delete corresponding record from local DB
	err = r.db.DeleteConnector(connectorID) // DeleteConnector uses the Dex Connector ID string
	if err != nil {
		// If Dex deletion succeeded but local failed, log inconsistency
		r.logger.Error("Failed to delete local DB record for connector after Dex deletion", zap.String("connectorID", connectorID), zap.Error(err))
		// Return error, but Dex connector is already gone. Manual cleanup might be needed.
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete connector metadata locally")
	}

	// Restart Dex Pod (handle with care)
	r.logger.Info("Restarting Dex pod to apply connector changes", zap.String("connectorID", connectorID))
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("Failed to restart Dex pod after connector deletion", zap.String("connectorID", connectorID), zap.Error(err))
		// Non-fatal?
	}

	r.logger.Info("Successfully deleted connector", zap.String("id", connectorID))
	return ctx.NoContent(http.StatusAccepted) // 202 or 204
}

// Helper min/max for masking - keep at end or move to utils
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
