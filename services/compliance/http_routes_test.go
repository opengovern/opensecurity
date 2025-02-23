package compliance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/opengovern/og-util/pkg/httpserver"

	"gorm.io/gorm"

	idocker "github.com/opengovern/og-util/pkg/dockertest"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/ory/dockertest/v3"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/compliance/db"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HttpServerSuite struct {
	suite.Suite

	handler *HttpHandler
	router  *echo.Echo
}

func (s *HttpServerSuite) SetupSuite() {
	require := s.Require()

	t := s.T()
	pool, err := dockertest.NewPool("")
	s.NoError(err, "pool constructed")
	err = pool.Client.Ping()
	s.NoError(err, "pinged pool")
	user, pass := "postgres", "123456"
	resource, err := pool.Run(user, "14.2-alpine", []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", pass)})
	s.NoError(err, "status postgres")
	t.Cleanup(func() {
		err := pool.Purge(resource)
		s.NoError(err, "purge resource %s", resource)
	})
	time.Sleep(5 * time.Second)

	var adb *gorm.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		cfg := &postgres.Config{
			Host:   idocker.GetDockerHost(),
			Port:   resource.GetPort("5432/tcp"),
			User:   user,
			Passwd: pass,
			DB:     "postgres",
		}

		logger, err := zap.NewProduction()
		s.NoError(err, "new zap logger")

		adb, err = postgres.NewClient(cfg, logger)
		s.NoError(err, "new postgres client")

		d, err := adb.DB()
		if err != nil {
			return err
		}

		return d.Ping()
	})
	s.handler = &HttpHandler{
		db: db.Database{Orm: adb},
	}

	// err = s.handler.db.Initialize()
	require.NoError(err, "db initialize")

	logger, err := zap.NewProduction()
	require.NoError(err, "new logger")
	s.router, _ = httpserver.Register(logger, s.handler)

}

func (s *HttpServerSuite) BeforeTest(suiteName, testName string) {
	//require := s.Require()
}

func (s *HttpServerSuite) AfterTest(suiteName, testName string) {
	//require := s.Require()
}

func doRequestJSONResponse(router *echo.Echo, method string, path string, request, response interface{}) (*httptest.ResponseRecorder, error) {
	var requestBytes []byte
	var err error

	if request != nil {
		requestBytes, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	rec, resp, err := sendRequest(router, method, path, requestBytes, "application/json")
	if err != nil {
		return nil, err
	}

	if response != nil {
		if err := json.Unmarshal(resp, response); err != nil {
			fmt.Println("Unmarshal error on response: " + string(resp))
			return nil, err
		}
	}

	if rec.Code != 200 {
		return nil, errors.New("not ok response code")
	}
	return rec, nil
}

func sendRequest(router *echo.Echo, method string, path string, request []byte, accept string) (*httptest.ResponseRecorder, []byte, error) {
	var r io.Reader
	if request != nil {
		r = bytes.NewReader(request)
	}

	req := httptest.NewRequest(method, path, r)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", accept)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Wrap in NopCloser in case the calling method wants to also read the body
	b, err := io.ReadAll(io.NopCloser(rec.Body))
	if err != nil {
		return nil, nil, err
	}
	return rec, b, nil
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, &HttpServerSuite{})
}

func (s *HttpServerSuite) TestDatabaseTableStructure() {
	time.Sleep(5 * time.Minute)
}
