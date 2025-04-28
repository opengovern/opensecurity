package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"

	"github.com/coreos/go-oidc/v3/oidc"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/auth/db"
	"github.com/opengovern/opensecurity/services/auth/utils"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type User struct {
	ID         string
	Email      string
	ExternalId string
	Role       api.Role
	LastLogin  time.Time
	CreatedAt  time.Time
}

type Server struct {
	host                string
	platformPublicKey   *rsa.PublicKey
	dexVerifier         *oidc.IDTokenVerifier
	dexClient           dexApi.DexClient
	logger              *zap.Logger
	db                  db.Database
	updateLoginUserList []User
	updateLogin         chan User
}

type DexClaims struct {
	Email           string                 `json:"email"`
	EmailVerified   bool                   `json:"email_verified"`
	Groups          []string               `json:"groups"`
	Name            string                 `json:"name"`
	FederatedClaims map[string]interface{} `json:"federated_claims"`
	jwt.StandardClaims
}

func (s *Server) UpdateLastLoginLoop() {
	for {
		finished := false
		for !finished {
			select {
			case userId := <-s.updateLogin:
				alreadyExists := false
				for _, user := range s.updateLoginUserList {
					if user.ExternalId == userId.ExternalId {
						alreadyExists = true
					}
				}

				if !alreadyExists {
					s.updateLoginUserList = append(s.updateLoginUserList, userId)
				}
			default:
				finished = true
			}
		}

		for i := 0; i < len(s.updateLoginUserList); i++ {
			user := s.updateLoginUserList[i]
			if user.ExternalId != "" {
				usr, err := utils.GetUserByEmail(user.Email, s.db)
				if err != nil {
					s.logger.Error("failed to get user metadata", zap.String(" External", user.ExternalId), zap.Error(err))
					continue
				}
				tim := time.Time{}
				if !usr.LastLogin.IsZero() {
					tim = usr.LastLogin
				}

				if time.Now().After(tim.Add(15 * time.Minute)) {
					s.logger.Info("updating metadata", zap.String("External Id", user.ExternalId))

					tim = time.Now()
					s.logger.Info("time is", zap.Time("time", tim))

					err = utils.UpdateUserLastLogin(user.ExternalId, tim, s.db)
					if err != nil {
						s.logger.Error("failed to update user metadata", zap.String("External Id", user.ExternalId), zap.Error(err))
					}
				}
			}

			s.updateLoginUserList = append(s.updateLoginUserList[:i], s.updateLoginUserList[i+1:]...)
			i--
		}
		time.Sleep(time.Second)
	}
}

func (s *Server) UpdateLastLogin(claim *userClaim) {
	timeNow := time.Now()
	doUpdate := false
	if claim.MemberSince == nil {
		claim.MemberSince = &timeNow
		doUpdate = true
	}
	if claim.UserLastLogin == nil {
		claim.UserLastLogin = &timeNow
		doUpdate = true
	} else {
		if time.Now().After(claim.UserLastLogin.Add(15 * time.Minute)) {
			claim.UserLastLogin = &timeNow
			doUpdate = true
		}
	}

	if doUpdate {
		s.updateLogin <- User{
			ExternalId: claim.ExternalUserID,
			LastLogin:  *claim.UserLastLogin,
			CreatedAt:  *claim.MemberSince,
			Email:      claim.Email,
		}
	}
}

func (s *Server) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	unAuth := &envoyauth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &envoyauth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoyauth.DeniedHttpResponse{
				Status: &envoytype.HttpStatus{Code: 401},
				Body:   http.StatusText(http.StatusUnauthorized),
			},
		},
	}

	httpRequest := req.GetAttributes().GetRequest().GetHttp()
	headers := httpRequest.GetHeaders()

	authHeader := headers[echo.HeaderAuthorization]
	if authHeader == "" {
		authHeader = headers[strings.ToLower(echo.HeaderAuthorization)]
	}

	user, err := s.Verify(ctx, authHeader)
	if err != nil {
		s.logger.Warn("denied access due to unsuccessful token verification",
			zap.String("reqId", httpRequest.Id),
			zap.String("path", httpRequest.Path),
			zap.String("method", httpRequest.Method),
			zap.Error(err))
		return unAuth, nil
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if user.Email == "" {
		s.logger.Warn("denied access due to failure to get email from token",
			zap.String("reqId", httpRequest.Id),
			zap.String("path", httpRequest.Path),
			zap.String("method", httpRequest.Method),
			zap.Error(err))
		return unAuth, nil
	}

	theUser, err := utils.GetUserByEmail(user.Email, s.db)
	if err != nil {
		s.logger.Warn("failed to get user",
			zap.String("userId", user.ExternalUserID),
			zap.String("email", user.Email),
			zap.Error(err))
		if errors.Is(err, errors.New("user disabled")) {
			return unAuth, nil
		}
		if errors.Is(err, errors.New("user not found")) {
			return unAuth, nil
		}

	}
	if theUser == nil {
		return unAuth, nil
	}
	user.Role = (api.Role)(theUser.Role)
	user.ExternalUserID = theUser.ExternalId
	user.MemberSince = &theUser.CreatedAt
	user.UserLastLogin = &theUser.LastLogin

	go s.UpdateLastLogin(user)

	return &envoyauth.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.OK),
		},

		HttpResponse: &envoyauth.CheckResponse_OkResponse{
			OkResponse: &envoyauth.OkHttpResponse{
				Headers: []*envoycore.HeaderValueOption{
					// {
					// 	Header: &envoycore.HeaderValue{
					// 		Key:   httpserver.XplatformWorkspaceIDHeader,
					// 		Value: rb.WorkspaceID,
					// 	},
					// },
					// {
					// 	Header: &envoycore.HeaderValue{
					// 		Key:   httpserver.XplatformWorkspaceNameHeader,
					// 		Value: rb.WorkspaceName,
					// 	},
					// },
					{
						Header: &envoycore.HeaderValue{
							Key:   httpserver.XPlatformUserIDHeader,
							Value: user.ExternalUserID,
						},
					},
					{
						Header: &envoycore.HeaderValue{
							Key:   httpserver.XPlatformUserRoleHeader,
							Value: string(user.Role),
						},
					},
					{
						Header: &envoycore.HeaderValue{
							Key:   httpserver.XPlatformUserConnectionsScope,
							Value: theUser.ExternalId,
						},
					},
				},
			},
		},
	}, nil
}

type userClaim struct {
	Role           api.Role
	Email          string
	MemberSince    *time.Time
	UserLastLogin  *time.Time
	ConnectionIDs  map[string][]string
	ExternalUserID string `json:"sub"`
	EmailVerified  bool
}

func (u userClaim) Valid() error {
	return nil
}

func (s *Server) Verify(ctx context.Context, authToken string) (*userClaim, error) {
	if !strings.HasPrefix(authToken, "Bearer ") {
		return nil, errors.New("invalid authorization token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authToken, "Bearer "))
	if token == "" {
		return nil, errors.New("missing authorization token")
	}

	var u userClaim

	s.logger.Info("dex verifier verifying")
	dv, err := s.dexVerifier.Verify(ctx, token)
	if err == nil {
		var claims json.RawMessage
		if err := dv.Claims(&claims); err != nil {
			s.logger.Error("dex verifier claim error", zap.Error(err))

			return nil, err
		}
		s.logger.Info("raw dex verifier claims", zap.Any("claims", string(claims)))
		var claimsMap DexClaims
		if err = json.Unmarshal(claims, &claimsMap); err != nil {
			s.logger.Error("dex verifier claim error", zap.Error(err))

			return nil, err
		}
		s.logger.Info("dex verifier claims", zap.Any("claims", claimsMap))

		return &userClaim{
			Email:         claimsMap.Email,
			EmailVerified: claimsMap.EmailVerified,
		}, nil
	} else {
		s.logger.Error("dex verifier verify error", zap.Error(err))
	}

	if s.platformPublicKey != nil {
		_, errk := jwt.ParseWithClaims(token, &u, func(token *jwt.Token) (interface{}, error) {
			return s.platformPublicKey, nil
		})
		if errk == nil {
			return &u, nil
		} else {
			fmt.Println("failed to auth with platform cred due to", errk)
		}
	}
	return nil, err
}

func newDexOidcVerifier(ctx context.Context, domain, clientId string) (*oidc.IDTokenVerifier, error) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 10,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	provider, err := oidc.NewProvider(
		oidc.InsecureIssuerURLContext(
			oidc.ClientContext(ctx, httpClient),
			domain,
		), domain,
	)
	if err != nil {
		return nil, err
	}

	return provider.Verifier(&oidc.Config{
		ClientID:          clientId,
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}), nil
}

func newDexClient(hostAndPort string) (dexApi.DexClient, error) {
	conn, err := grpc.NewClient(hostAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return dexApi.NewDexClient(conn), nil
}
