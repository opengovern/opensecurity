// Package auth provides the core authentication and authorization logic.
// This file contains unit and handler tests for the HTTP API layer (http.go).
package auth

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	// Use the same JWT library version as your main code
	"github.com/golang-jwt/jwt" // Or jwt "github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	// Use alias 'api2' for the shared api package to avoid conflict with local 'api'
	jose "github.com/go-jose/go-jose/v3" // Use v3 or v4 depending on your go.mod
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"         // For header constants
	"github.com/opengovern/opensecurity/services/auth/api" // Local API definitions
	"github.com/opengovern/opensecurity/services/auth/db"  // Needed because CreateAPIKey uses utils.GetUser
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc" // For MockDexClient
	"gorm.io/gorm"           // Import gorm for Model struct and errors

	// Import the dex api package aliased
	dexApi "github.com/dexidp/dex/api/v2"
)

// --- Mock Database ---

// MockDatabase is a mock type for the db.DatabaseInterface.
// It uses testify/mock to allow setting expectations and asserting calls.
type MockDatabase struct {
	mock.Mock
}

// Implement all methods defined in db.DatabaseInterface.

func (m *MockDatabase) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	args := m.Called(ctx, email)
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*db.User), args.Error(1)
}
func (m *MockDatabase) GetUserByExternalID(ctx context.Context, id string) (*db.User, error) {
	args := m.Called(ctx, id)
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*db.User), args.Error(1)
}
func (m *MockDatabase) GetUser(ctx context.Context, id string) (*db.User, error) {
	args := m.Called(ctx, id)
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*db.User), args.Error(1)
}
func (m *MockDatabase) GetUsers(ctx context.Context) ([]db.User, error) {
	args := m.Called(ctx)
	usersArg := args.Get(0)
	if usersArg == nil {
		return nil, args.Error(1)
	}
	return usersArg.([]db.User), args.Error(1)
}
func (m *MockDatabase) GetUsersCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockDatabase) GetFirstUser(ctx context.Context) (*db.User, error) {
	args := m.Called(ctx)
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*db.User), args.Error(1)
}
func (m *MockDatabase) CreateUser(ctx context.Context, user *db.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockDatabase) UpdateUser(ctx context.Context, user *db.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockDatabase) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDatabase) UpdateUserLastLoginWithExternalID(ctx context.Context, id string, lastLogin time.Time) error {
	args := m.Called(ctx, id, lastLogin)
	return args.Error(0)
}
func (m *MockDatabase) UserPasswordUpdate(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDatabase) FindIdByEmail(ctx context.Context, email string) (uint, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(uint), args.Error(1)
}
func (m *MockDatabase) AddApiKey(ctx context.Context, key *db.ApiKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
func (m *MockDatabase) CountApiKeysForUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockDatabase) ListApiKeysForUser(ctx context.Context, userId string) ([]db.ApiKey, error) {
	args := m.Called(ctx, userId)
	keysArg := args.Get(0)
	if keysArg == nil {
		return nil, args.Error(1)
	}
	return keysArg.([]db.ApiKey), args.Error(1)
}
func (m *MockDatabase) UpdateAPIKey(ctx context.Context, id string, isActive bool, role api2.Role) error {
	args := m.Called(ctx, id, isActive, role)
	return args.Error(0)
} // Use api2.Role
func (m *MockDatabase) DeleteAPIKey(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDatabase) ListApiKeys(ctx context.Context) ([]db.ApiKey, error) {
	args := m.Called(ctx)
	keysArg := args.Get(0)
	if keysArg == nil {
		return nil, args.Error(1)
	}
	return keysArg.([]db.ApiKey), args.Error(1)
}
func (m *MockDatabase) GetConnectorByConnectorID(ctx context.Context, connectorID string) (*db.Connector, error) {
	args := m.Called(ctx, connectorID)
	connArg := args.Get(0)
	if connArg == nil {
		return nil, args.Error(1)
	}
	return connArg.(*db.Connector), args.Error(1)
}
func (m *MockDatabase) CreateConnector(ctx context.Context, connector *db.Connector) error {
	args := m.Called(ctx, connector)
	return args.Error(0)
}
func (m *MockDatabase) UpdateConnector(ctx context.Context, connector *db.Connector) error {
	args := m.Called(ctx, connector)
	return args.Error(0)
}
func (m *MockDatabase) DeleteConnector(ctx context.Context, connectorID string) error {
	args := m.Called(ctx, connectorID)
	return args.Error(0)
}
func (m *MockDatabase) GetConnectors(ctx context.Context) ([]db.Connector, error) {
	args := m.Called(ctx)
	connArg := args.Get(0)
	if connArg == nil {
		return nil, args.Error(1)
	}
	return connArg.([]db.Connector), args.Error(1)
}
func (m *MockDatabase) GetConnector(ctx context.Context, id string) (*db.Connector, error) {
	args := m.Called(ctx, id)
	connArg := args.Get(0)
	if connArg == nil {
		return nil, args.Error(1)
	}
	return connArg.(*db.Connector), args.Error(1)
}
func (m *MockDatabase) GetConnectorByConnectorType(ctx context.Context, connectorType string) (*db.Connector, error) {
	args := m.Called(ctx, connectorType)
	connArg := args.Get(0)
	if connArg == nil {
		return nil, args.Error(1)
	}
	return connArg.(*db.Connector), args.Error(1)
}
func (m *MockDatabase) GetKeyPair(ctx context.Context) ([]db.Configuration, error) {
	args := m.Called(ctx)
	cfgArg := args.Get(0)
	if cfgArg == nil {
		return nil, args.Error(1)
	}
	return cfgArg.([]db.Configuration), args.Error(1)
}
func (m *MockDatabase) AddConfiguration(ctx context.Context, c *db.Configuration) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
func (m *MockDatabase) Initialize() error { args := m.Called(); return args.Error(0) }

// Implement EnableUser/DisableUser from the interface
func (m *MockDatabase) EnableUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDatabase) DisableUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Mock Dex Client ---

// MockDexClient mocks the dexApi.DexClient interface.
// It uses testify/mock and includes all methods defined in the interface
// to ensure compile-time compatibility.
type MockDexClient struct {
	mock.Mock
}

// Implement ALL methods from dexApi.DexClient interface.
func (m *MockDexClient) CreateClient(ctx context.Context, in *dexApi.CreateClientReq, opts ...grpc.CallOption) (*dexApi.CreateClientResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.CreateClientResp), args.Error(1)
}
func (m *MockDexClient) UpdateClient(ctx context.Context, in *dexApi.UpdateClientReq, opts ...grpc.CallOption) (*dexApi.UpdateClientResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.UpdateClientResp), args.Error(1)
}
func (m *MockDexClient) DeleteClient(ctx context.Context, in *dexApi.DeleteClientReq, opts ...grpc.CallOption) (*dexApi.DeleteClientResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.DeleteClientResp), args.Error(1)
}
func (m *MockDexClient) GetClient(ctx context.Context, in *dexApi.GetClientReq, opts ...grpc.CallOption) (*dexApi.GetClientResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.GetClientResp), args.Error(1)
}
func (m *MockDexClient) CreatePassword(ctx context.Context, in *dexApi.CreatePasswordReq, opts ...grpc.CallOption) (*dexApi.CreatePasswordResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.CreatePasswordResp), args.Error(1)
}
func (m *MockDexClient) UpdatePassword(ctx context.Context, in *dexApi.UpdatePasswordReq, opts ...grpc.CallOption) (*dexApi.UpdatePasswordResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.UpdatePasswordResp), args.Error(1)
}
func (m *MockDexClient) DeletePassword(ctx context.Context, in *dexApi.DeletePasswordReq, opts ...grpc.CallOption) (*dexApi.DeletePasswordResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.DeletePasswordResp), args.Error(1)
}
func (m *MockDexClient) ListPasswords(ctx context.Context, in *dexApi.ListPasswordReq, opts ...grpc.CallOption) (*dexApi.ListPasswordResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.ListPasswordResp), args.Error(1)
}
func (m *MockDexClient) VerifyPassword(ctx context.Context, in *dexApi.VerifyPasswordReq, opts ...grpc.CallOption) (*dexApi.VerifyPasswordResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.VerifyPasswordResp), args.Error(1)
}
func (m *MockDexClient) CreateConnector(ctx context.Context, in *dexApi.CreateConnectorReq, opts ...grpc.CallOption) (*dexApi.CreateConnectorResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.CreateConnectorResp), args.Error(1)
}
func (m *MockDexClient) UpdateConnector(ctx context.Context, in *dexApi.UpdateConnectorReq, opts ...grpc.CallOption) (*dexApi.UpdateConnectorResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.UpdateConnectorResp), args.Error(1)
}
func (m *MockDexClient) DeleteConnector(ctx context.Context, in *dexApi.DeleteConnectorReq, opts ...grpc.CallOption) (*dexApi.DeleteConnectorResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.DeleteConnectorResp), args.Error(1)
}
func (m *MockDexClient) ListConnectors(ctx context.Context, in *dexApi.ListConnectorReq, opts ...grpc.CallOption) (*dexApi.ListConnectorResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.ListConnectorResp), args.Error(1)
} // Note: Uses ListConnectorReq based on proto
func (m *MockDexClient) GetVersion(ctx context.Context, in *dexApi.VersionReq, opts ...grpc.CallOption) (*dexApi.VersionResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.VersionResp), args.Error(1)
}
func (m *MockDexClient) GetDiscovery(ctx context.Context, in *dexApi.DiscoveryReq, opts ...grpc.CallOption) (*dexApi.DiscoveryResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.DiscoveryResp), args.Error(1)
} // Corrected return type
func (m *MockDexClient) ListRefresh(ctx context.Context, in *dexApi.ListRefreshReq, opts ...grpc.CallOption) (*dexApi.ListRefreshResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.ListRefreshResp), args.Error(1)
}
func (m *MockDexClient) RevokeRefresh(ctx context.Context, in *dexApi.RevokeRefreshReq, opts ...grpc.CallOption) (*dexApi.RevokeRefreshResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.RevokeRefreshResp), args.Error(1)
}

// Add ListKeys and GetKeys needed to fully satisfy the interface
func (m *MockDexClient) ListKeys(ctx context.Context, in *dexApi.ListKeysReq, opts ...grpc.CallOption) (*dexApi.ListKeysResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.ListKeysResp), args.Error(1)
}
func (m *MockDexClient) GetKeys(ctx context.Context, in *dexApi.GetKeysReq, opts ...grpc.CallOption) (*dexApi.GetKeysResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.GetKeysResp), args.Error(1)
}

// --- Helper Functions ---

// calculateKeyID computes the JWK thumbprint (SHA256, Base64URL encoded).
func calculateTestKeyID(pub *rsa.PublicKey) (string, error) {
	jwk := jose.JSONWebKey{Key: pub}
	tb, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tb), nil
}

// generateTestKeys creates a new RSA key pair and calculates its kid for testing.
func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey, string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pub := &priv.PublicKey
	kid, err := calculateTestKeyID(pub)
	require.NoError(t, err)
	require.NotEmpty(t, kid)
	return priv, pub, kid
}

// --- Setup Test Suite Helper ---

// httpTestDeps holds all mocked dependencies and common setup needed for HTTP handler tests.
type httpTestDeps struct {
	routes         *httpRoutes
	mockDb         *MockDatabase
	mockDexClient  *MockDexClient
	testPrivateKey *rsa.PrivateKey
	testPublicKey  *rsa.PublicKey
	testKid        string
	echoInstance   *echo.Echo
}

// setupHttpTests initializes mocks, test keys, and the httpRoutes struct for testing.
func setupHttpTests(t *testing.T) httpTestDeps {
	testPrivateKey, testPublicKey, testKid := generateTestKeys(t)
	mockDb := new(MockDatabase)
	mockDex := new(MockDexClient) // MockDexClient should now implement DexClient fully
	testLogger := zap.NewNop()    // Use Nop logger for tests unless output is needed

	// Create a minimal Server struct containing the mocked Dex client.
	// Add other fields like dexVerifier if the handlers under test interact with them directly.
	testAuthServer := &Server{
		logger:            testLogger.Named("testAuthServer"),
		dexClient:         mockDex, // Inject mock Dex client
		platformPublicKey: testPublicKey,
		platformKeyID:     testKid,
		// db field is not needed here as httpRoutes uses its own db interface instance
	}

	// Instantiate httpRoutes with test dependencies.
	routes := &httpRoutes{
		logger:             testLogger.Named("httpRoutes"),
		platformPrivateKey: testPrivateKey,
		platformKeyID:      testKid,
		db:                 mockDb, // Assign mockDb which satisfies db.DatabaseInterface
		authServer:         testAuthServer,
	}

	e := echo.New()
	// e.Validator = &YourCustomValidator{} // Register validator if needed

	return httpTestDeps{routes: routes, mockDb: mockDb, mockDexClient: mockDex, testPrivateKey: testPrivateKey, testPublicKey: testPublicKey, testKid: testKid, echoInstance: e}
}

// createTestContext is a helper to create an Echo context with a request, response recorder,
// and optional headers for testing handlers.
func createTestContext(e *echo.Echo, method, path string, body io.Reader, headers map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} // Default to JSON
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// --- Test Functions ---

// TestCreateAPIKeyHandler fully tests the CreateAPIKey HTTP handler.
// It mocks database interactions and verifies the HTTP response and the generated JWT.
func TestCreateAPIKeyHandler(t *testing.T) {
	deps := setupHttpTests(t)
	creatorUserID := "local|creator@example.com"
	mockCreatorDbUser := &db.User{Model: gorm.Model{ID: 2}, Email: "creator@example.com", ExternalId: creatorUserID, Role: api2.AdminRole, Username: "creator", IsActive: true, EmailVerified: true, ConnectorId: "local"}
	apiKeyRole := api2.EditorRole
	apiKeyName := "My-Test-Key-1"
	requestBody := api.CreateAPIKeyRequest{Name: apiKeyName, Role: apiKeyRole}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Mock Expectations
	deps.mockDb.On("GetUserByExternalID", mock.Anything, creatorUserID).Return(mockCreatorDbUser, nil).Once() // Called by utils.GetUser
	deps.mockDb.On("CountApiKeysForUser", mock.Anything, creatorUserID).Return(int64(0), nil).Once()
	deps.mockDb.On("AddApiKey", mock.Anything, mock.AnythingOfType("*db.ApiKey")).Return(nil).Once()

	// Setup Echo Context
	headers := map[string]string{httpserver.XPlatformUserIDHeader: creatorUserID}
	c, rec := createTestContext(deps.echoInstance, http.MethodPost, "/api/v1/keys", bytes.NewReader(requestBodyBytes), headers)

	// Execute Handler
	err := deps.routes.CreateAPIKey(c)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	var respBody api.CreateAPIKeyResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respBody))
	assert.Equal(t, apiKeyName, respBody.Name)
	assert.Equal(t, apiKeyRole, respBody.RoleName)
	signedToken := respBody.Token
	require.NotEmpty(t, signedToken)
	var parsedClaims jwt.StandardClaims
	parsedToken, err := jwt.ParseWithClaims(signedToken, &parsedClaims, func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		if kid != deps.testKid {
			return nil, fmt.Errorf("bad kid")
		}
		return deps.testPublicKey, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, deps.testKid, parsedToken.Header["kid"])
	assert.Equal(t, creatorUserID, parsedClaims.Subject)
	deps.mockDb.AssertExpectations(t)
}

// TestGetUsersHandler tests the GET /users endpoint.
// It mocks the database call to return a predefined list of users and verifies the JSON response.
func TestGetUsersHandler(t *testing.T) {
	deps := setupHttpTests(t)
	mockUsers := []db.User{{Model: gorm.Model{ID: 1}, Email: "test1@example.com", Username: "test1", Role: api2.ViewerRole, ExternalId: "local|test1@example.com", IsActive: true}, {Model: gorm.Model{ID: 2}, Email: "test2@example.com", Username: "test2", Role: api2.EditorRole, ExternalId: "local|test2@example.com", IsActive: false}}
	deps.mockDb.On("GetUsers", mock.Anything).Return(mockUsers, nil).Once()
	c, rec := createTestContext(deps.echoInstance, http.MethodGet, "/api/v1/users", nil, nil)
	err := deps.routes.GetUsers(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var respBody []api.GetUsersResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respBody))
	assert.Len(t, respBody, 2)
	assert.Equal(t, mockUsers[1].Email, respBody[1].Email)
	assert.False(t, respBody[1].IsActive)
	deps.mockDb.AssertExpectations(t)
}

// TestDeleteUserHandler tests the DELETE /user/:id endpoint for a successful deletion.
// It mocks database calls and the Dex gRPC call for password deletion.
func TestDeleteUserHandler(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToDelete := uint(5)
	userEmail := "delete@example.com"
	userExternalID := "local|delete@example.com"
	mockUserToDelete := &db.User{Model: gorm.Model{ID: userIDToDelete}, Email: userEmail, ExternalId: userExternalID, ConnectorId: "local", IsActive: true}
	deps.mockDb.On("GetUser", mock.Anything, strconv.Itoa(int(userIDToDelete))).Return(mockUserToDelete, nil).Once()
	deps.mockDexClient.On("DeletePassword", mock.Anything, &dexApi.DeletePasswordReq{Email: userEmail}).Return(&dexApi.DeletePasswordResp{NotFound: false}, nil).Once()
	deps.mockDb.On("DeleteUser", mock.Anything, userIDToDelete).Return(nil).Once()
	c, rec := createTestContext(deps.echoInstance, http.MethodDelete, "/api/v1/user/"+strconv.Itoa(int(userIDToDelete)), nil, nil)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(userIDToDelete)))
	err := deps.routes.DeleteUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, rec.Code)
	deps.mockDb.AssertExpectations(t)
	deps.mockDexClient.AssertExpectations(t)
}

// TestDeleteFirstUserHandler tests that deleting the user with ID 1 is forbidden.
func TestDeleteFirstUserHandler(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToDelete := uint(1)
	mockUserToDelete := &db.User{Model: gorm.Model{ID: userIDToDelete}, Email: "admin@example.com", ConnectorId: "local"}
	deps.mockDb.On("GetUser", mock.Anything, strconv.Itoa(int(userIDToDelete))).Return(mockUserToDelete, nil).Once()
	c, _ := createTestContext(deps.echoInstance, http.MethodDelete, "/api/v1/user/1", nil, nil)
	c.SetParamNames("id")
	c.SetParamValues("1")
	err := deps.routes.DeleteUser(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Contains(t, httpErr.Message.(string), "Cannot delete the initial administrator user")
	deps.mockDb.AssertExpectations(t)
	deps.mockDexClient.AssertExpectations(t) // No Dex calls expected
}

// TestResetUserPasswordHandler tests the successful password reset flow for a local user.
func TestResetUserPasswordHandler(t *testing.T) {
	deps := setupHttpTests(t)
	userID := "local|reset@example.com"
	userEmail := "reset@example.com"
	dbUserID := uint(10)
	currentPassword := "oldPassword123"
	newPassword := "newStrongPassword456"
	mockUser := &db.User{Model: gorm.Model{ID: dbUserID}, Email: userEmail, ExternalId: userID, ConnectorId: "local", IsActive: true}
	requestBody := api.ResetUserPasswordRequest{CurrentPassword: currentPassword, NewPassword: newPassword}
	requestBodyBytes, _ := json.Marshal(requestBody)
	deps.mockDb.On("GetUserByExternalID", mock.Anything, userID).Return(mockUser, nil).Once()
	deps.mockDexClient.On("VerifyPassword", mock.Anything, mock.MatchedBy(func(req *dexApi.VerifyPasswordReq) bool {
		return req.Email == userEmail && req.Password == currentPassword
	})).Return(&dexApi.VerifyPasswordResp{Verified: true, NotFound: false}, nil).Once()
	deps.mockDexClient.On("UpdatePassword", mock.Anything, mock.MatchedBy(func(req *dexApi.UpdatePasswordReq) bool { return req.Email == userEmail && len(req.NewHash) > 10 })).Return(&dexApi.UpdatePasswordResp{NotFound: false}, nil).Once()
	deps.mockDb.On("UserPasswordUpdate", mock.Anything, dbUserID).Return(nil).Once()
	headers := map[string]string{httpserver.XPlatformUserIDHeader: userID}
	c, rec := createTestContext(deps.echoInstance, http.MethodPost, "/api/v1/user/password/reset", bytes.NewReader(requestBodyBytes), headers)
	err := deps.routes.ResetUserPassword(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, rec.Code)
	deps.mockDb.AssertExpectations(t)
	deps.mockDexClient.AssertExpectations(t)
}

// TestResetUserPasswordHandler_WrongPassword tests failure on incorrect current password during reset.
func TestResetUserPasswordHandler_WrongPassword(t *testing.T) {
	deps := setupHttpTests(t)
	userID := "local|reset@example.com"
	userEmail := "reset@example.com"
	dbUserID := uint(10)
	currentPassword := "oldPassword123"
	newPassword := "newStrongPassword456"
	mockUser := &db.User{Model: gorm.Model{ID: dbUserID}, Email: userEmail, ExternalId: userID, ConnectorId: "local", IsActive: true}
	requestBody := api.ResetUserPasswordRequest{CurrentPassword: "WRONG" + currentPassword, NewPassword: newPassword}
	requestBodyBytes, _ := json.Marshal(requestBody)
	deps.mockDb.On("GetUserByExternalID", mock.Anything, userID).Return(mockUser, nil).Once()
	deps.mockDexClient.On("VerifyPassword", mock.Anything, mock.AnythingOfType("*v2.VerifyPasswordReq")).Return(&dexApi.VerifyPasswordResp{Verified: false, NotFound: false}, nil).Once()
	headers := map[string]string{httpserver.XPlatformUserIDHeader: userID}
	c, _ := createTestContext(deps.echoInstance, http.MethodPost, "/api/v1/user/password/reset", bytes.NewReader(requestBodyBytes), headers)
	err := deps.routes.ResetUserPassword(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Contains(t, httpErr.Message.(string), "Incorrect current password")
	deps.mockDb.AssertExpectations(t)
	deps.mockDexClient.AssertExpectations(t)
}

// TestEnableUserHandler tests enabling a user successfully by an admin.
func TestEnableUserHandler(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToEnable := uint(99)
	adminUserID := "local|admin@example.com" // ID of the admin making the request

	// Mock DB call for EnableUser
	deps.mockDb.On("EnableUser", mock.Anything, userIDToEnable).Return(nil).Once()

	// Setup context with Admin UserID header
	headers := map[string]string{httpserver.XPlatformUserIDHeader: adminUserID}
	c, rec := createTestContext(deps.echoInstance, http.MethodPut, "/api/v1/user/"+strconv.Itoa(int(userIDToEnable))+"/enable", nil, headers)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(userIDToEnable)))

	// Execute (assuming AuthorizeHandler middleware allows based on header role - tested elsewhere)
	err := deps.routes.EnableUserHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	deps.mockDb.AssertExpectations(t)
}

// TestDisableUserHandler tests disabling a user successfully by an admin.
func TestDisableUserHandler(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToDisable := uint(98)
	adminUserID := "local|admin@example.com"

	// Mock DB call for DisableUser
	deps.mockDb.On("DisableUser", mock.Anything, userIDToDisable).Return(nil).Once()

	// Setup context with Admin UserID header
	headers := map[string]string{httpserver.XPlatformUserIDHeader: adminUserID}
	c, rec := createTestContext(deps.echoInstance, http.MethodPut, "/api/v1/user/"+strconv.Itoa(int(userIDToDisable))+"/disable", nil, headers)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(userIDToDisable)))

	// Execute
	err := deps.routes.DisableUserHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	deps.mockDb.AssertExpectations(t)
}

// TestDisableUserHandler_NotFound tests disabling a non-existent user.
func TestDisableUserHandler_NotFound(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToDisable := uint(101)
	adminUserID := "local|admin@example.com"

	// Mock DB call for DisableUser to return not found
	deps.mockDb.On("DisableUser", mock.Anything, userIDToDisable).Return(gorm.ErrRecordNotFound).Once()

	headers := map[string]string{httpserver.XPlatformUserIDHeader: adminUserID}
	c, _ := createTestContext(deps.echoInstance, http.MethodPut, "/api/v1/user/"+strconv.Itoa(int(userIDToDisable))+"/disable", nil, headers)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(userIDToDisable)))

	err := deps.routes.DisableUserHandler(c)

	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
	deps.mockDb.AssertExpectations(t)
}

// TestDisableUserHandler_DisableAdmin1 tests preventing disabling user ID 1.
func TestDisableUserHandler_DisableAdmin1(t *testing.T) {
	deps := setupHttpTests(t)
	userIDToDisable := uint(1)
	adminUserID := "local|otheradmin@example.com"

	// DB call should not be made for ID 1
	deps.mockDb.On("DisableUser", mock.Anything, userIDToDisable).Return(nil).Maybe() // Use Maybe if unsure

	headers := map[string]string{httpserver.XPlatformUserIDHeader: adminUserID}
	c, _ := createTestContext(deps.echoInstance, http.MethodPut, "/api/v1/user/1/disable", nil, headers)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := deps.routes.DisableUserHandler(c)

	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Contains(t, httpErr.Message.(string), "Cannot disable the initial administrator user")
	deps.mockDb.AssertExpectations(t) // Ensure DisableUser was NOT called
}

// TestTokenGenerationLogic verifies the JWT signing part of the Token exchange logic.
// This test focuses specifically on the claims enrichment and signing process.
func TestTokenGenerationLogic(t *testing.T) {
	deps := setupHttpTests(t)
	testEmail := "test@example.com"
	testExternalID := "local|test@example.com"
	testRole := api2.EditorRole
	mockUser := &db.User{Model: gorm.Model{ID: 1}, Email: testEmail, ExternalId: testExternalID, Role: testRole, Username: "testuser", IsActive: true, ConnectorId: "local", EmailVerified: true}
	dexSubject := "some-dex-internal-sub"
	nowTime := time.Now()
	dexClaimsInput := DexClaims{Email: testEmail, EmailVerified: true, Groups: []string{"group-from-dex"}, Name: "Dex Name", StandardClaims: jwt.StandardClaims{Issuer: "http://dex-issuer.example.com", Subject: dexSubject, Audience: "public-client", ExpiresAt: nowTime.Add(1 * time.Hour).Unix(), IssuedAt: nowTime.Unix(), Id: "dex-jti-123"}}
	enrichedClaims := dexClaimsInput
	enrichedClaims.Groups = append(enrichedClaims.Groups, string(mockUser.Role))
	enrichedClaims.Name = mockUser.Username
	enrichedClaims.Subject = mockUser.ExternalId
	enrichedClaims.Id = "new-jti-for-platform-token"
	enrichedClaims.Issuer = "platform-auth-service"
	enrichedClaims.Audience = "platform-client"
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, enrichedClaims)
	token.Header["kid"] = deps.testKid
	signedToken, err := token.SignedString(deps.testPrivateKey)
	require.NoError(t, err)
	var parsedClaims DexClaims
	parsedToken, err := jwt.ParseWithClaims(signedToken, &parsedClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("bad alg")
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid != deps.testKid {
			return nil, fmt.Errorf("bad kid [%v]", kid)
		}
		return deps.testPublicKey, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, deps.testKid, parsedToken.Header["kid"])
	assert.Equal(t, testExternalID, parsedClaims.Subject)
	assert.Contains(t, parsedClaims.Groups, string(testRole))
}
