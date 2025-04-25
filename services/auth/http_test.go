// http_test.go
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
	"github.com/opengovern/opensecurity/services/auth/db"
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

// MockDatabase is a mock type for the db.DatabaseInterface
type MockDatabase struct {
	mock.Mock
}

// Implement methods defined in db.DatabaseInterface that are used by handlers under test
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

// --- Fix #3: Correct Role type in UpdateAPIKey mock signature ---
func (m *MockDatabase) UpdateAPIKey(ctx context.Context, id string, isActive bool, role api2.Role) error {
	args := m.Called(ctx, id, isActive, role)
	return args.Error(0)
}
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

// --- Mock Dex Client ---

// MockDexClient mocks the dexApi.DexClient interface
type MockDexClient struct {
	mock.Mock
}

// --- Fix #1 & #2: Correct GetDiscovery return type ---
func (m *MockDexClient) GetDiscovery(ctx context.Context, in *dexApi.DiscoveryReq, opts ...grpc.CallOption) (*dexApi.DiscoveryResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	// Type assertion must match the corrected return type
	return respArg.(*dexApi.DiscoveryResp), args.Error(1)
}

// --- Implement other methods called by handlers or needed to satisfy interface ---
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
func (m *MockDexClient) GetVersion(ctx context.Context, in *dexApi.VersionReq, opts ...grpc.CallOption) (*dexApi.VersionResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.VersionResp), args.Error(1)
}
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
// Add ListKeys and GetKeys needed to fully satisfy the interface

// --- Corrected ListKeys Mock ---
// ListKeys mocks the Dex gRPC ListKeys method.
// The input type is ListKeysReq (plural).
// In http_test.go - MockDexClient
func (m *MockDexClient) ListKeys(ctx context.Context, in *dexApi.ListKeysReq, opts ...grpc.CallOption) (*dexApi.ListKeysResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	if respArg == nil {
		return nil, args.Error(1)
	}
	return respArg.(*dexApi.ListKeysResp), args.Error(1)
}

// --- Corrected GetKeys Mock ---
// GetKeys mocks the Dex gRPC GetKeys method.
// (The original signature provided was actually correct for this one)
func (m *MockDexClient) GetKeys(ctx context.Context, in *dexApi.GetKeysReq, opts ...grpc.CallOption) (*dexApi.GetKeysResp, error) {
	args := m.Called(ctx, in)
	respArg := args.Get(0)
	errArg := args.Error(1)

	if respArg == nil {
		return nil, errArg
	}
	return respArg.(*dexApi.GetKeysResp), errArg
}

// --- Helper Functions ---
func calculateTestKeyID(pub *rsa.PublicKey) (string, error) {
	jwk := jose.JSONWebKey{Key: pub}
	tb, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tb), nil
}
func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey, string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pub := &priv.PublicKey
	kid, err := calculateTestKeyID(pub)
	require.NoError(t, err)
	require.NotEmpty(t, kid)
	return priv, pub, kid
}
func now() int64 { return time.Now().Unix() }

// --- Setup Test Suite Helper ---
type httpTestDeps struct {
	routes         *httpRoutes
	mockDb         *MockDatabase
	mockDexClient  *MockDexClient
	testPrivateKey *rsa.PrivateKey
	testPublicKey  *rsa.PublicKey
	testKid        string
	echoInstance   *echo.Echo
}

func setupHttpTests(t *testing.T) httpTestDeps {
	testPrivateKey, testPublicKey, testKid := generateTestKeys(t)
	mockDb := new(MockDatabase)
	mockDex := new(MockDexClient) // MockDexClient should now implement DexClient fully
	testLogger := zap.NewNop()

	testAuthServer := &Server{
		logger:            testLogger.Named("testAuthServer"),
		dexClient:         mockDex, // This assignment should now work
		platformPublicKey: testPublicKey,
		platformKeyID:     testKid,
	}

	routes := &httpRoutes{
		logger:             testLogger.Named("httpRoutes"),
		platformPrivateKey: testPrivateKey,
		platformKeyID:      testKid,
		db:                 mockDb, // Assign mockDb (satisfies interface)
		authServer:         testAuthServer,
	}
	e := echo.New()
	return httpTestDeps{routes: routes, mockDb: mockDb, mockDexClient: mockDex, testPrivateKey: testPrivateKey, testPublicKey: testPublicKey, testKid: testKid, echoInstance: e}
}

// Helper to create a test Echo context
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
func TestCreateAPIKeyHandler(t *testing.T) {
	deps := setupHttpTests(t)
	creatorUserID := "local|creator@example.com"
	mockCreatorDbUser := &db.User{Model: gorm.Model{ID: 2}, Email: "creator@example.com", ExternalId: creatorUserID, Role: api2.AdminRole, Username: "creator", IsActive: true, EmailVerified: true, ConnectorId: "local"}
	apiKeyRole := api2.EditorRole
	apiKeyName := "My-Test-Key-1"
	requestBody := api.CreateAPIKeyRequest{Name: apiKeyName, Role: apiKeyRole}
	requestBodyBytes, _ := json.Marshal(requestBody)
	deps.mockDb.On("GetUserByExternalID", mock.Anything, creatorUserID).Return(mockCreatorDbUser, nil).Once()
	deps.mockDb.On("CountApiKeysForUser", mock.Anything, creatorUserID).Return(int64(0), nil).Once()
	deps.mockDb.On("AddApiKey", mock.Anything, mock.AnythingOfType("*db.ApiKey")).Return(nil).Once()
	headers := map[string]string{httpserver.XPlatformUserIDHeader: creatorUserID}
	c, rec := createTestContext(deps.echoInstance, http.MethodPost, "/api/v1/keys", bytes.NewReader(requestBodyBytes), headers)
	err := deps.routes.CreateAPIKey(c)
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

// TestGetUsersHandler tests listing users.
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

// TestDeleteUserHandler tests deleting a user.
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

// TestDeleteFirstUserHandler tests preventing deletion of user ID 1.
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
	deps.mockDexClient.AssertExpectations(t)
}

// TestResetUserPasswordHandler tests the password reset flow.
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

// TestResetUserPasswordHandler_WrongPassword tests failure on incorrect current password.
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

// TestTokenGenerationLogic verifies the JWT signing part of the Token exchange logic.
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
