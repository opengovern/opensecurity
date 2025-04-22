package auth

import (
	"context"
	"crypto/rsa"
	"crypto/sha512"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	"github.com/opengovern/opensecurity/services/auth/db"

	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/auth/api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// var (
// 	//go:embed email/invite.html
// 	inviteEmailTemplate string
// )

type httpRoutes struct {
	logger *zap.Logger

	platformPrivateKey *rsa.PrivateKey
	db                 db.Database
	authServer         *Server
}


func (r *httpRoutes) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	// VAlidate token
	v1.GET("/check", r.Check)
	v1.POST("/token", r.Token)

	// USERS
	v1.GET("/users", httpserver.AuthorizeHandler(r.GetUsers, api2.EditorRole))                                      //checked
	v1.GET("/user/:id", httpserver.AuthorizeHandler(r.GetUserDetails, api2.EditorRole))                             //checked
	v1.GET("/me", httpserver.AuthorizeHandler(r.GetMe, api2.EditorRole))                                            //checked
	v1.POST("/user", httpserver.AuthorizeHandler(r.CreateUser, api2.EditorRole))                                    //checked
	v1.PUT("/user", httpserver.AuthorizeHandler(r.UpdateUser, api2.EditorRole))                                     //checked
	v1.GET("/user/password/check", httpserver.AuthorizeHandler(r.CheckUserPasswordChangeRequired, api2.ViewerRole)) //checked
	v1.POST("/user/password/reset", httpserver.AuthorizeHandler(r.ResetUserPassword, api2.ViewerRole))              //checked
	v1.DELETE("/user/:id", httpserver.AuthorizeHandler(r.DeleteUser, api2.AdminRole))                               //checked
	// API KEYS
	v1.POST("/keys", httpserver.AuthorizeHandler(r.CreateAPIKey, api2.AdminRole)) //checked
	v1.GET("/keys", httpserver.AuthorizeHandler(r.ListAPIKeys, api2.AdminRole))   //checked
	v1.DELETE("/key/:id", httpserver.AuthorizeHandler(r.DeleteAPIKey, api2.AdminRole))
	v1.PUT("/key/:id", httpserver.AuthorizeHandler(r.EditAPIKey, api2.AdminRole))
	// connectors
	v1.GET("/connectors", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole))
	v1.GET("/connectors/supported-connector-types", httpserver.AuthorizeHandler(r.GetSupportedType, api2.AdminRole))
	v1.GET("/connector/:type", httpserver.AuthorizeHandler(r.GetConnectors, api2.AdminRole))
	v1.POST("/connector", httpserver.AuthorizeHandler(r.CreateConnector, api2.AdminRole))
	v1.POST("/connector/auth0", httpserver.AuthorizeHandler(r.CreateAuth0Connector, api2.AdminRole))

	v1.PUT("/connector", httpserver.AuthorizeHandler(r.UpdateConnector, api2.AdminRole))
	v1.DELETE("/connector/:id", httpserver.AuthorizeHandler(r.DeleteConnector, api2.AdminRole))

}

func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

func (r *httpRoutes) Check(ctx echo.Context) error {
	checkRequest := envoyauth.CheckRequest{
		Attributes: &envoyauth.AttributeContext{
			Request: &envoyauth.AttributeContext_Request{
				Http: &envoyauth.AttributeContext_HttpRequest{
					Headers: make(map[string]string),
				},
			},
		},
	}

	for k, v := range ctx.Request().Header {
		if len(v) == 0 {
			checkRequest.Attributes.Request.Http.Headers[k] = ""
		} else {
			checkRequest.Attributes.Request.Http.Headers[k] = v[0]
		}
	}
	originalUri, err := url.Parse(ctx.Request().Header.Get("X-Original-URI"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid original uri")
	}
	checkRequest.Attributes.Request.Http.Path = originalUri.Path
	checkRequest.Attributes.Request.Http.Method = ctx.Request().Header.Get("X-Original-Method")

	res, err := r.authServer.Check(ctx.Request().Context(), &checkRequest)
	if err != nil {
		return err
	}

	if res.Status.Code != int32(codes.OK) {
		return echo.NewHTTPError(http.StatusUnauthorized, res.Status.Message)
	}

	if res.GetOkResponse() == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "no ok response")
	}

	for _, header := range res.GetOkResponse().GetHeaders() {
		if header == nil || header.Header == nil {
			continue
		}
		ctx.Response().Header().Set(header.Header.Key, header.Header.Value)
	}

	return ctx.NoContent(http.StatusOK)
}
func (r *httpRoutes) Token(ctx echo.Context) error {
	var req api.GetTokenRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if req.Code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}
	if req.CallBackUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "callback url is required")
	}

	// make an http call to dex pod to get the token
	domain:=os.Getenv("DEX_DOMAIN")
	// http call to dex pod with code and callback url serach params
	// make body 
	body:= map[string]string{
		"code": req.Code,
		"redirect_uri": req.CallBackUrl,
		"grant_type": "authorization_code",
		"client_id": "public-client",
		"client_secret": "",
	}
	url := fmt.Sprintf("http://%s/token", domain)
	// set headers
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	// make http call
	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create request body")
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create request")
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, err := client.Do(request)
	if err != nil {
		r.logger.Error("failed to make request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to make request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		r.logger.Error("failed to get token", zap.String("status", response.Status))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get token")
	}
	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to decode response")
	}
	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get access token")
	}
	dv,err := r.authServer.dexVerifier.Verify(ctx.Request().Context(), accessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to verify token")
	}
	var claims json.RawMessage
	if err := dv.Claims(&claims); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get claims")
	}
	var claimsMap DexClaims
	if err = json.Unmarshal(claims, &claimsMap); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to unmarshal claims")
	}
	// check if the claimsMap is empty
	// find the user in the db
	user, err := r.db.GetUserByEmail(claimsMap.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get user")
	}
	// add user role to claims
	claimsMap.Groups = append(claimsMap.Groups, string(user.Role))
	// convert claims map to json
	
	// update dv claims

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,claimsMap )
	// sign the token with the private key
	signedToken, err := token.SignedString(r.platformPrivateKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to sign token")
	}
	// return the token
	tokenResponse["access_token"] = signedToken
	return ctx.JSON(http.StatusOK, tokenResponse)


}

// GetUsers godoc
//
//	@Summary		List Users
//	@Description	Retrieves a list of users who are members of the workspace.
//	@Security		BearerToken
//	@Tags			users
//	@Produce		json
//	@Param			request	body	api.GetUsersRequest	false	"Request Body"
//	@Success		200		{array}	api.GetUsersResponse
//	@Router			/auth/api/v1/users [get]
func (r *httpRoutes) GetUsers(ctx echo.Context) error {

	var req api.GetUsersRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	users, err := r.db.GetUsers()
	if err != nil {
		return err
	}
	var resp []api.GetUsersResponse
	for _, u := range users {
		temp_resp := api.GetUsersResponse{
			ID:            u.ID,
			UserName:      u.Username,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			ExternalId:    u.ExternalId,
			CreatedAt:     u.CreatedAt,
			RoleName:      u.Role,
			IsActive:      u.IsActive,
			ConnectorId:   u.ConnectorId,
		}
		if u.LastLogin.IsZero() {
			temp_resp.LastActivity = nil
		} else {
			temp_resp.LastActivity = &u.LastLogin
		}
		resp = append(resp, temp_resp)

	}
	return ctx.JSON(http.StatusOK, resp)
}

// GetUserDetails godoc
//
//	@Summary		Get User details
//	@Description	Returns user details by specified user id.
//	@Security		BearerToken
//	@Tags			users
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	api.GetUserResponse
//	@Router			/auth/api/v1/user/{userId} [get]
func (r *httpRoutes) GetUserDetails(ctx echo.Context) error {

	userID := ctx.Param("id")
	userID, err := url.QueryUnescape(userID)
	if err != nil {
		return err
	}
	user, err := r.db.GetUser(userID)
	if err != nil {
		return err
	}

	resp := api.GetUserResponse{
		ID:            user.ID,
		UserName:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		Blocked:       user.IsActive,
		RoleName:      user.Role,
	}
	// check if LastLogin is Default go time value remove it
	if user.LastLogin.IsZero() {
		resp.LastActivity = nil
	} else {
		resp.LastActivity = &user.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)

}

// GetMe godoc
//
//	@Summary		Get Me
//	@Description	Returns my user details
//	@Security		BearerToken
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	api.GetMeResponse
//	@Router			/auth/api/v1/me [get]
func (r *httpRoutes) GetMe(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)

	db_user, err := r.db.GetUserByExternalID(userID)
	user ,_:= utils.DbUserToApi(db_user)

	if err != nil {
		return err
	}

	resp := api.GetMeResponse{

		ID:            user.ID,
		UserName:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		Blocked:       user.IsActive,
		Role:          user.Role,
		MemberSince:   user.CreatedAt,
		ConnectorId:   user.ConnectorId,
	}
	if user.LastLogin.IsZero() {
		resp.LastLogin = nil
		resp.LastActivity = nil
	} else {
		resp.LastLogin = &user.LastLogin
		resp.LastActivity = &user.LastLogin
	}

	return ctx.JSON(http.StatusOK, resp)

}

// CreateAPIKey godoc
//
//	@Summary		Create Workspace Key
//	@Description	Creates workspace key for the defined role with the defined name in the workspace.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Param			request	body		api.CreateAPIKeyRequest	true	"Request Body"
//	@Success		200		{object}	api.CreateAPIKeyResponse
//	@Failure		406		{object}	echo.HTTPError
//	@Router			/auth/api/v1/key/create [post]
func (r *httpRoutes) CreateAPIKey(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)
	var req api.CreateAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	usr, err := utils.GetUser(userID, r.db)
	if err != nil {
		r.logger.Error("failed to get user", zap.Error(err))
		return err
	}

	if usr == nil {
		return errors.New("failed to find user in auth")
	}

	u := userClaim{
		Role:           req.Role,
		Email:          usr.Email,
		ExternalUserID: usr.ExternalId,
	}

	if r.platformPrivateKey == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "platform api key is disabled")
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, &u).SignedString(r.platformPrivateKey)
	if err != nil {
		r.logger.Error("failed to create token", zap.Error(err))
		return err
	}

	masked := fmt.Sprintf("%s...%s", token[:10], token[len(token)-10:])

	hash := sha512.New()
	_, err = hash.Write([]byte(token))
	if err != nil {
		r.logger.Error("failed to hash token", zap.Error(err))
		return err
	}
	keyHash := hex.EncodeToString(hash.Sum(nil))
	r.logger.Info("hashed token")

	currentKeyCount, err := r.db.CountApiKeysForUser(userID)
	if err != nil {
		r.logger.Error("failed to get user API Keys count", zap.Error(err))
		return err
	}
	if currentKeyCount > 5 {
		return echo.NewHTTPError(http.StatusNotAcceptable, "maximum number of keys for user reached")
	}
	r.logger.Info("creating API Key")
	apikey := db.ApiKey{
		Name:          req.Name,
		Role:          req.Role,
		CreatorUserID: userID,
		IsActive:      true,
		MaskedKey:     masked,
		KeyHash:       keyHash,
	}

	r.logger.Info("adding API Key")
	err = r.db.AddApiKey(&apikey)
	if err != nil {
		r.logger.Error("failed to add API Key", zap.Error(err))
		return err
	}

	return ctx.JSON(http.StatusCreated, api.CreateAPIKeyResponse{
		ID:        apikey.ID,
		Name:      apikey.Name,
		Active:    apikey.IsActive,
		CreatedAt: apikey.CreatedAt,
		RoleName:  apikey.Role,
		Token:     token,
	})
}

// DeleteAPIKey godoc
//
//	@Summary		Delete Workspace Key
//	@Description	Deletes the specified workspace key by ID.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Param			id	path		string	true	"Key ID"
//	@Success		200	{object}	nil
//	@Router			/auth/api/v1/key/{id}/delete [delete]
func (r *httpRoutes) DeleteAPIKey(ctx echo.Context) error {
	// TODO: Ask from ANIL what should i do
	// userId := httpserver.GetUserID(ctx)
	id := ctx.Param("id")

	integer_id, err := (strconv.ParseUint(id, 10, 32))

	err = r.db.DeleteAPIKey(integer_id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}
func (r *httpRoutes) EditAPIKey(ctx echo.Context) error {
	// TODO: Ask from ANIL what should i do
	// userId := httpserver.GetUserID(ctx)
	var req api.EditAPIKeyRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	err := r.db.UpdateAPIKey(id, req.IsActive, req.Role)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}

// ListAPIKeys godoc
//
//	@Summary		Get API keys List
//	@Description	Gets list of all keys.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Success		200	{object}	[]api.APIKeyResponse
//	@Router			/auth/api/v1/keys [get]
func (r *httpRoutes) ListAPIKeys(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)
	keys, err := r.db.ListApiKeysForUser(userID)
	if err != nil {
		return err
	}

	var resp []api.APIKeyResponse
	for _, key := range keys {
		resp = append(resp, api.APIKeyResponse{
			ID:            key.ID,
			CreatedAt:     key.CreatedAt,
			Name:          key.Name,
			RoleName:      key.Role,
			CreatorUserID: key.CreatorUserID,
			Active:        key.IsActive,
			MaskedKey:     key.MaskedKey,
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

// CreateUser godoc
//
//	@Summary		Create User
//	@Description	Creates User.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Param			request	body	api.CreateUserRequest	true	"Request Body"
//	@Success		200
//	@Router			/auth/api/v1/user [post]
func (r *httpRoutes) CreateUser(ctx echo.Context) error {

	var req api.CreateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := r.DoCreateUser(req)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (r *httpRoutes) DoCreateUser(req api.CreateUserRequest) error {

	if req.EmailAddress == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email address is required")
	}

	user, err := r.db.GetUserByEmail(req.EmailAddress)
	if err != nil {
		r.logger.Error("failed to get user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	if user != nil && user.Email != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email already used")
	}

	count, err := r.db.GetUsersCount()
	if err != nil {
		r.logger.Error("failed to get users count", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users count")
	}
	adminAccount := false
	var firstUser *db.User

	if count == 1 {
		firstUser, err = r.db.GetFirstUser()
		if err != nil {
			r.logger.Error("failed to get first user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get first user")
		}

	} else if count == 0 {
		adminAccount = true
	}
	if adminAccount && (req.Role == nil || *req.Role != api2.AdminRole) {
		return echo.NewHTTPError(http.StatusBadRequest, "You should define admin role")
	}

	if adminAccount && firstUser != nil {
		err = r.DoDeleteUser(firstUser.Email)
		if err != nil {
			r.logger.Error("failed to delete first user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete first user")
		}
	}

	connector := ""
	userId := fmt.Sprintf("%v|%s", req.ConnectorId, req.EmailAddress)
	if req.Password != nil {
		connector = "local"
		userId := fmt.Sprintf("local|%s", req.EmailAddress)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("failed to hash token", zap.Error(err))
			return err
		}

		dexReq := &dexApi.CreatePasswordReq{
			Password: &dexApi.Password{
				UserId: userId,
				Email:  req.EmailAddress,
				Hash:   hashedPassword,
			},
		}

		resp, err := r.authServer.dexClient.CreatePassword(context.TODO(), dexReq)
		if err != nil {
			r.logger.Error("failed to create dex password", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex password")
		}
		if resp.AlreadyExists {
			dexReq := &dexApi.UpdatePasswordReq{
				Email:   req.EmailAddress,
				NewHash: hashedPassword,
			}

			_, err = r.authServer.dexClient.UpdatePassword(context.TODO(), dexReq)
			if err != nil {
				r.logger.Error("failed to update dex password", zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex password")
			}
		}
	}

	role := api2.ViewerRole
	if req.Role != nil {
		role = *req.Role
	}

	requirePasswordChange := true
	if adminAccount {
		requirePasswordChange = false
	}

	newUser := &db.User{
		Email:                 req.EmailAddress,
		Username:              req.EmailAddress,
		FullName:              req.EmailAddress,
		Role:                  role,
		EmailVerified:         false,
		ConnectorId:           connector,
		ExternalId:            userId,
		RequirePasswordChange: requirePasswordChange,
		IsActive:              true,
	}
	err = r.db.CreateUser(newUser)
	if err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create user")
	}
	return nil
}

// UpdateUser godoc
//
//	@Summary		Update User
//	@Description	Updates User.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Param			request	body	api.UpdateUserRequest	true	"Request Body"
//	@Success		200
//	@Router			/auth/api/v1/user [put]
func (r *httpRoutes) UpdateUser(ctx echo.Context) error {
	var req api.UpdateUserRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.EmailAddress == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email address is required")
	}

	user, err := r.db.GetUserByEmail(req.EmailAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}

	if req.Password != nil && req.ConnectorId == "local" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.logger.Error("failed to hash token", zap.Error(err))
			return err
		}

		dexReq := &dexApi.UpdatePasswordReq{
			Email:   req.EmailAddress,
			NewHash: hashedPassword,
		}

		resp, err := r.authServer.dexClient.UpdatePassword(context.TODO(), dexReq)
		if err != nil {
			r.logger.Error("failed to update dex password", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex password")
		}
		if resp.NotFound {
			dexReq := &dexApi.CreatePasswordReq{
				Password: &dexApi.Password{
					UserId: fmt.Sprintf("local|%s", req.EmailAddress),
					Email:  req.EmailAddress,
					Hash:   hashedPassword,
				},
			}

			_, err = r.authServer.dexClient.CreatePassword(context.TODO(), dexReq)
			if err != nil {
				r.logger.Error("failed to create dex password", zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex password")
			}
		}

		err = r.db.UserPasswordUpdate(user.ID)
		if err != nil {
			r.logger.Error("failed to update user", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to update user")
		}
	}

	if req.Role != nil {
		update_user := &db.User{
			Model: gorm.Model{
				ID: user.ID,
			},
			Role:        *req.Role,
			IsActive:    req.IsActive,
			Username:    req.UserName,
			FullName:    req.FullName,
			Email:       user.Email,
			ExternalId:  fmt.Sprintf("%v|%s", req.ConnectorId, user.Email),
			ConnectorId: req.ConnectorId,
		}
		err = r.db.UpdateUser(update_user)
		if err != nil {
			r.logger.Error("failed to update user", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to update user")
		}
	}

	return ctx.NoContent(http.StatusOK)
}

// DeleteUser godoc
//
//	@Summary		Delete User
//	@Description	Delete User.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Param			email_address	path	string	true	"Request Body"
//	@Success		200
//	@Router			/auth/api/v3/user/{email_address}/delete [delete]
func (r *httpRoutes) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	err := r.DoDeleteUser(id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}

func (r *httpRoutes) DoDeleteUser(id string) error {
	user, err2 := r.db.GetUser(id)

	if err2 != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user does not exist")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user does not exist")
	}
	if user.ID == 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot delete the first user")
	}
	dexReq := &dexApi.DeletePasswordReq{
		Email: user.Email,
	}

	_, err := r.authServer.dexClient.DeletePassword(context.TODO(), dexReq)
	if err != nil {
		r.logger.Error("failed to remove dex password", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to remove dex password")
	}

	err = r.db.DeleteUser(user.ID)
	if err != nil {
		r.logger.Error("failed to delete user", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create user")
	}
	return nil
}

// CheckUserPasswordChangeRequired godoc
//
//	@Summary		Delete User
//	@Description	Delete User.
//	@Security		BearerToken
//	@Tags			keys
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v3/user/password/check [get]
func (r *httpRoutes) CheckUserPasswordChangeRequired(ctx echo.Context) error {
	userId := httpserver.GetUserID(ctx)

	user, err := r.db.GetUserByExternalID(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}

	if user.RequirePasswordChange {
		return ctx.String(http.StatusOK, "CHANGE_REQUIRED")
	} else {
		return ctx.String(http.StatusOK, "CHANGE_NOT_REQUIRED")
	}
}

// ResetUserPassword godoc
//
//	@Summary		Reset current user password
//	@Description	Reset current user password
//	@Security		BearerToken
//	@Tags			user
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v3/user/password/reset [post]
func (r *httpRoutes) ResetUserPassword(ctx echo.Context) error {

	userId := httpserver.GetUserID(ctx)

	user, err := r.db.GetUserByExternalID(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}

	var req api.ResetUserPasswordRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if user.ConnectorId != "local" {
		return echo.NewHTTPError(http.StatusBadRequest, "user connector should be local")
	}

	dexReq := &dexApi.VerifyPasswordReq{
		Email:    user.Email,
		Password: req.CurrentPassword,
	}

	resp, err := r.authServer.dexClient.VerifyPassword(context.TODO(), dexReq)
	if err != nil {
		r.logger.Error("failed to validate dex password", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate dex password")
	}
	if resp.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	if !resp.Verified {
		return echo.NewHTTPError(http.StatusUnauthorized, "current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("failed to hash token", zap.Error(err))
		return err
	}

	passwordUpdateReq := &dexApi.UpdatePasswordReq{
		Email:   user.Email,
		NewHash: hashedPassword,
	}

	passwordUpdateResp, err := r.authServer.dexClient.UpdatePassword(context.TODO(), passwordUpdateReq)
	if err != nil {
		r.logger.Error("failed to update dex password", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update dex password")
	}
	if passwordUpdateResp.NotFound {
		dexReq := &dexApi.CreatePasswordReq{
			Password: &dexApi.Password{
				UserId: fmt.Sprintf("local|%s", user.Email),
				Email:  user.Email,
				Hash:   hashedPassword,
			},
		}

		_, err = r.authServer.dexClient.CreatePassword(context.TODO(), dexReq)
		if err != nil {
			r.logger.Error("failed to create dex password", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex password")
		}
	}

	err = r.db.UserPasswordUpdate(user.ID)
	if err != nil {
		r.logger.Error("failed to update user", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update user")
	}

	return ctx.NoContent(http.StatusAccepted)
}

// GetConnector godoc
//
//	@Summary		Get  Connectors list
//	@Description	Returns a list  connectors. can have connector type param
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connectors [GET]

func (r *httpRoutes) GetConnectors(ctx echo.Context) error {
	req := &dexApi.ListConnectorReq{}
	connectorType := ctx.Param("type")
	// Execute the ListConnectors RPC.
	respDex, err := r.authServer.dexClient.ListConnectors(context.TODO(), req)
	if err != nil {
		r.logger.Error("failed to list connectors", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to list connectors")

	}

	connectors := respDex.Connectors

	var resp []api.GetConnectorsResponse
	for _, connector := range connectors {

		localConnector, err := r.db.GetConnectorByConnectorID(connector.Id)
		if err != nil {
			r.logger.Error("failed to get connector", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to get connector")
		}
		if connectorType != "" && strings.ToLower(connectorType) != strings.ToLower(connector.Type) {
			continue
		}
		if connector.Id == "local" {
			continue
		}
		info := api.GetConnectorsResponse{
			ID:          localConnector.ID,
			ConnectorID: connector.Id,
			Type:        connector.Type,
			Name:        connector.Name,
			SubType:     localConnector.ConnectorSubType,
			UserCount:   localConnector.UserCount,
			CreatedAt:   localConnector.CreatedAt,
			LastUpdate:  localConnector.LastUpdate,
		}

		// If the connector is of type "oidc", attempt to extract Issuer and ClientID
		if strings.ToLower(connector.Type) == "oidc" {
			var config api.OIDCConfig
			var data map[string]interface{}
			err := json.Unmarshal(connector.Config, &config)
			new_err := json.Unmarshal(connector.Config, &data)
			r.logger.Info("data", zap.Any("data", data))
			if new_err != nil {
				r.logger.Error("Failed to unmarshal OIDC config for connector", zap.Error(err))
			}
			if err != nil {
				r.logger.Error("Failed to unmarshal OIDC config for connector", zap.Error(err))
			} else {
				info.Issuer = config.Issuer
				info.ClientID = config.ClientID
				// Note: Omitting ClientSecret for security reasons
			}
		}

		resp = append(resp, info)
	}
	return ctx.JSON(http.StatusOK, resp)
}

// GetSupportedConnectors godoc
//
//	@Summary		Get Supported Connectors
//	@Description	Returns a list of supported connectors.
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connectors/ [GET]

func (r *httpRoutes) GetSupportedType(ctx echo.Context) error {
	var connectors []api.GetSupportedConnectorTypeResponse

	subTypes := utils.SupportedConnectors["oidc"]
	subTypesNames := utils.SupportedConnectorsNames["oidc"]

	var types []api.ConnectorSubTypes
	for i, key := range subTypes {
		types = append(types, api.ConnectorSubTypes{
			ID:   key,
			Name: subTypesNames[i],
		})
	}
	connectors = append(connectors, api.GetSupportedConnectorTypeResponse{
		ConnectorType: "oidc",
		SubTypes:      types,
	})

	return ctx.JSON(http.StatusOK, connectors)

}

// CreateConnector godoc
//
//	@Summary		Create Connector
//	@Description	Creates new OIDC connector.
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connector/supported-connector-types [post]
func (r *httpRoutes) CreateConnector(ctx echo.Context) error {
	var req api.CreateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.ConnectorType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "connector type is required")
	}
	connectorTypeLower := strings.ToLower(req.ConnectorType)
	creator := utils.GetConnectorCreator(connectorTypeLower)
	if creator == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "connector type is not supported")
	}

	// default
	connectorSubTypeLower := "general" // default
	if req.ConnectorSubType != "" {
		connectorSubTypeLower = strings.ToLower(req.ConnectorSubType)
	} else {
		r.logger.Info("No connector_sub_type specified. Defaulting to 'general'")
	}
	if !utils.IsSupportedSubType(connectorTypeLower, connectorSubTypeLower) {
		err := fmt.Sprintf("unsupported connector_sub_type '%s' for connector_type '%s'", connectorSubTypeLower, connectorTypeLower)
		r.logger.Info(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)

	}
	switch connectorSubTypeLower {
	case "general":
		// Required: issuer, client_id, client_secret
		if strings.TrimSpace(req.Issuer) == "" {
			r.logger.Warn("Missing 'issuer' for 'general' OIDC connector")
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "issuer is required for 'general' OIDC connector",
			})

		}

		// Set default id and name if not provided
		if strings.TrimSpace(req.ID) == "" {
			req.ID = "oidc-default"
		}
		if strings.TrimSpace(req.Name) == "" {
			req.Name = "General OIDC"
		}

	case "entraid":
		// Required: tenant_id, client_id, client_secret
		if strings.TrimSpace(req.TenantID) == "" {
			err := "Missing 'tenant_id' for 'entraid' OIDC connector"
			r.logger.Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)

		}
		// fetching issuer

		// Set default id and name if not provided
		if strings.TrimSpace(req.ID) == "" {
			req.ID = "entra-id"

		}
		if strings.TrimSpace(req.Name) == "" {
			req.Name = "AzureAD/EntraID"

		}

	case "google-workspace":
		// Required: client_id, client_secret

		// Set default id and name if not provided
		if strings.TrimSpace(req.ID) == "" {
			req.ID = "google-oidc"

		}
		if strings.TrimSpace(req.Name) == "" {
			req.Name = "Google Workspaces "

		}
	}
	dexRequest := utils.CreateConnectorRequest{
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		ID:               req.ID,
		Name:             req.Name,
	}
	dexreq, err := creator(dexRequest)
	if err != nil {
		r.logger.Error("Error on Creating dex request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexreq)
	if err != nil {
		r.logger.Error("failed to create dex connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex connector")
	}
	if res.AlreadyExists {
		return echo.NewHTTPError(http.StatusBadRequest, "connector already exists")
	}
	err = r.db.CreateConnector(&db.Connector{
		LastUpdate:       time.Now(),
		ConnectorID:      req.ID,
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
	})
	if err != nil {
		r.logger.Error("failed to create connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create connector")
	}
	// restart dex pod on connector creation
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("failed to restart dex pod", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to restart dex pod")
	}

	return ctx.JSON(http.StatusCreated, res)
}

// CreateConnector godoc
//
//	@Summary		Create Connector
//	@Description	Creates new OIDC connector.
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connector/supported-connector-types [post]
func (r *httpRoutes) CreateAuth0Connector(ctx echo.Context) error {
	var req api.CreateAuth0ConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	creator := utils.CreateAuth0Connector

	dexRequest := utils.CreateAuth0ConnectorRequest{
		Issuer:       req.Issuer,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Domain:       req.Domain,
	}
	dexreq, err := creator(dexRequest)
	if err != nil {
		r.logger.Error("Error on Creating dex request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	publicUris := req.PublicURIS
	publicClientResp, _ := r.authServer.dexClient.GetClient(context.TODO(), &dexApi.GetClientReq{
		Id: "public-client",
	})

	r.logger.Info("public URIS", zap.Any("uris", publicUris))

	if publicClientResp != nil && publicClientResp.Client != nil {
		publicClientReq := dexApi.UpdateClientReq{
			Id:           "public-client",
			Name:         "Public Client",
			RedirectUris: publicUris,
		}

		_, err = r.authServer.dexClient.UpdateClient(context.TODO(), &publicClientReq)
		if err != nil {
			r.logger.Error("Auth Migrator: failed to create dex public client", zap.Error(err))
			return err
		}
	} else {
		publicClientReq := dexApi.CreateClientReq{
			Client: &dexApi.Client{
				Id:           "public-client",
				Name:         "Public Client",
				RedirectUris: publicUris,
				Public:       true,
			},
		}

		_, err = r.authServer.dexClient.CreateClient(context.TODO(), &publicClientReq)
		if err != nil {
			r.logger.Error("Auth Migrator: failed to create dex public client", zap.Error(err))
			return err
		}
	}
	privateUris := req.PrivateURIS

	r.logger.Info("private URIS", zap.Any("uris", privateUris))

	privateClientResp, _ := r.authServer.dexClient.GetClient(context.TODO(), &dexApi.GetClientReq{
		Id: "private-client",
	})
	if privateClientResp != nil && privateClientResp.Client != nil {
		privateClientReq := dexApi.UpdateClientReq{
			Id:           "private-client",
			Name:         "Private Client",
			RedirectUris: privateUris,
		}

		_, err = r.authServer.dexClient.UpdateClient(context.TODO(), &privateClientReq)
		if err != nil {
			r.logger.Error("Auth Migrator: failed to create dex private client", zap.Error(err))
			return err
		}
	} else {
		privateClientReq := dexApi.CreateClientReq{
			Client: &dexApi.Client{
				Id:           "private-client",
				Name:         "Private Client",
				RedirectUris: privateUris,
				Secret:       "secret",
			},
		}

		_, err = r.authServer.dexClient.CreateClient(context.TODO(), &privateClientReq)
		if err != nil {
			r.logger.Error("Auth Migrator: failed to create dex private client", zap.Error(err))
			return err
		}
	}

	res, err := r.authServer.dexClient.CreateConnector(context.TODO(), dexreq)
	if err != nil {
		r.logger.Error("failed to create dex connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create dex connector")
	}
	if res.AlreadyExists {
		return echo.NewHTTPError(http.StatusBadRequest, "connector already exists")
	}
	err = r.db.CreateConnector(&db.Connector{
		LastUpdate:       time.Now(),
		ConnectorID:      "auth0",
		ConnectorType:    "oidc",
		ConnectorSubType: "auth0",
	})
	if err != nil {
		r.logger.Error("failed to create connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create connector")
	}

	// restart dex pod on connector creation
	err = utils.RestartDexPod()
	if err != nil {
		r.logger.Error("failed to restart dex pod", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to restart dex pod")
	}

	return ctx.JSON(http.StatusCreated, res)
}

// UpdateConnector godoc
//
//	@Summary		Update Connector
//	@Description	Update new OIDC connector.
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connector [put]

func (r *httpRoutes) UpdateConnector(ctx echo.Context) error {
	var req api.UpdateConnectorRequest
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.ID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID required")
	}

	if !utils.IsSupportedSubType(req.ConnectorType, req.ConnectorSubType) {
		err := fmt.Sprintf("unsupported connector_sub_type '%s' for connector_type '%s'", req.ConnectorType, req.ConnectorSubType)
		r.logger.Info(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)

	}
	switch req.ConnectorSubType {
	case "general":
		// Required: issuer, client_id, client_secret
		if strings.TrimSpace(req.Issuer) == "" {
			err := "Missing 'issuer' for 'general' OIDC connector update"
			r.logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)

		}
		// client_id and client_secret are already validated as required in the struct

	case "entraid":
		// Required: tenant_id, client_id, client_secret
		if strings.TrimSpace(req.TenantID) == "" {
			err := "Missing 'tenant_id' for 'entraid' OIDC connector update"
			r.logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)

		}
		// client_id and client_secret are already validated as required in the struct

	case "google-workspace":
		// Required: client_id, client_secret
		// No additional fields needed
		// client_id and client_secret are already validated as required in the struct

	default:
		err := fmt.Sprintf("unsupported connector_sub_type: %s", req.ConnectorSubType)
		r.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)

	}
	dexRequest := utils.UpdateConnectorRequest{
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
		Issuer:           req.Issuer,
		TenantID:         req.TenantID,
		ClientID:         req.ClientID,
		ClientSecret:     req.ClientSecret,
		ID:               req.ConnectorID,
	}

	dexreq, err := utils.UpdateOIDCConnector(dexRequest)
	if err != nil {
		r.logger.Error("Error on Creating dex request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := r.authServer.dexClient.UpdateConnector(context.TODO(), dexreq)
	if err != nil {
		r.logger.Error("failed to update dex connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update dex connector")

	}

	if res.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "connector not found")
	}
	err = r.db.UpdateConnector(&db.Connector{
		Model: gorm.Model{
			ID: req.ID,
		},
		LastUpdate:       time.Now(),
		ConnectorID:      req.ConnectorID,
		ConnectorType:    req.ConnectorType,
		ConnectorSubType: req.ConnectorSubType,
	})
	if err != nil {
		r.logger.Error("failed to update connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to update connector")
	}
	return ctx.JSON(http.StatusAccepted, res)

}

// DeleteConnector godoc
//
//	@Summary		Delete Connector
//	@Description	Delete  OIDC connector.
//	@Security		BearerToken
//	@Tags			connectors
//	@Produce		json
//	@Success		200
//	@Router			/auth/api/v1/connector/:id [Delete]

func (r *httpRoutes) DeleteConnector(ctx echo.Context) error {
	connectorID := ctx.Param("id")
	if connectorID == "" {
		r.logger.Error("Missing connector_id in DeleteConnectorByIDHandler request")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "connector_id is required in the URL path",
		})
	}
	req := &dexApi.DeleteConnectorReq{
		Id: connectorID,
	}

	resp, err := r.authServer.dexClient.DeleteConnector(context.TODO(), req)
	if err != nil {
		r.logger.Error("failed to delete connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to delete connector")
	}
	if resp.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "connector not found")
	}
	err = r.db.DeleteConnector(connectorID)
	if err != nil {
		r.logger.Error("failed to delete connector", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to delete connector")
	}

	return ctx.NoContent(http.StatusAccepted)
}
