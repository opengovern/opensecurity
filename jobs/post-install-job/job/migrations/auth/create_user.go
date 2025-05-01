package auth

import (
	"context"
	"fmt"

	dexApi "github.com/dexidp/dex/api/v2"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/services/auth/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type Migration struct{}

func (m Migration) IsGitBased() bool {
	return false
}

func (m Migration) AttachmentFolderPath() string {
	return ""
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "auth",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{Orm: orm, Logger: logger}

	err = orm.AutoMigrate(&db.User{})
	if err != nil {
		return fmt.Errorf("auto migrate user: %w", err)
	}

	dexClient, err := newDexClient(conf.DexGrpcAddress)
	if err != nil {
		logger.Error("Auth Migrator: failed to create dex client", zap.Error(err))
		return err
	}

	count, err := dbm.GetUsersCount()
	if err != nil {
		return err
	}
	if count > 0 {
		logger.Info("users already exist")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(conf.DefaultDexUserPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Auth Migrator: failed to generate password", zap.Error(err))
		return err
	}

	req := dexApi.CreatePasswordReq{
		Password: &dexApi.Password{
			Email:    conf.DefaultDexUserEmail,
			Username: conf.DefaultDexUserName,
			UserId:   fmt.Sprintf("local|%s", conf.DefaultDexUserEmail),
			Hash:     hashedPassword,
		},
	}

	_, err = dexClient.CreatePassword(ctx, &req)
	if err != nil {
		logger.Error("Auth Migrator: failed to create dex password", zap.Error(err))
		return err
	}

	role := api.AdminRole

	user := &db.User{

		Email:                 conf.DefaultDexUserEmail,
		Username:              conf.DefaultDexUserEmail,
		FullName:              conf.DefaultDexUserEmail,
		Role:                  role,
		ExternalId:            fmt.Sprintf("local|%s", conf.DefaultDexUserEmail),
		ConnectorId:           "local",
		IsActive:              true,
		RequirePasswordChange: true,
	}
	err = dbm.CreateUser(user)
	if err != nil {
		logger.Error("Auth Migrator: failed to create user in database", zap.Error(err))
		return err
	}

	return nil
}

func newDexClient(hostAndPort string) (dexApi.DexClient, error) {
	conn, err := grpc.NewClient(hostAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return dexApi.NewDexClient(conn), nil
}
