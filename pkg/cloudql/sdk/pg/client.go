package pg

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClientCached(c config.ClientConfig, cache *connection.ConnectionCache, ctx context.Context) (Client, error) {
	value, ok := cache.Get(ctx, "opengovernance-pg-client")
	if ok {
		return value.(Client), nil
	}

	plugin.Logger(ctx).Warn("pg client is not cached, creating a new one")

	client, err := NewClient(ctx, c)
	if err != nil {
		return Client{}, err
	}

	cache.Set(ctx, "opengovernance-pg-client", client)

	return client, nil
}

func NewInventoryClientCached(c config.ClientConfig, cache *connection.ConnectionCache, ctx context.Context) (Client, error) {
	value, ok := cache.Get(ctx, "opengovernance-inventory-pg-client")
	if ok {
		return value.(Client), nil
	}

	plugin.Logger(ctx).Warn("pg client is not cached, creating a new one")

	c.PgDatabase = aws.String("inventory")
	client, err := NewClient(ctx, c)
	if err != nil {
		return Client{}, err
	}

	cache.Set(ctx, "opengovernance-inventory-pg-client", client)

	return client, nil
}

func NewCoreClient(c config.ClientConfig, ctx context.Context) (Client, error) {
	c.PgDatabase = aws.String("core")
	return NewClient(ctx, c)
}

func NewClient(ctx context.Context, c config.ClientConfig) (Client, error) {
	if c.PgHost == nil || len(*c.PgHost) == 0 {
		host := os.Getenv("PG_HOST")
		c.PgHost = &host
	}

	if c.PgPort == nil || len(*c.PgPort) == 0 {
		port := os.Getenv("PG_PORT")
		c.PgPort = &port
	}

	if c.PgUser == nil || len(*c.PgUser) == 0 {
		user := os.Getenv("PG_USER")
		c.PgUser = &user
	}

	if c.PgPassword == nil || len(*c.PgPassword) == 0 {
		password := os.Getenv("PG_PASSWORD")
		c.PgPassword = &password
	}

	if c.PgDatabase == nil || len(*c.PgDatabase) == 0 {
		database := os.Getenv("PG_DATABASE")
		c.PgDatabase = &database
	}

	if c.PgSslMode == nil || len(*c.PgSslMode) == 0 {
		sslMode := os.Getenv("PG_SSL_MODE")
		c.PgSslMode = &sslMode
	}

	cfg := postgres.Config{
		Host:    *c.PgHost,
		Port:    *c.PgPort,
		User:    *c.PgUser,
		Passwd:  *c.PgPassword,
		DB:      *c.PgDatabase,
		SSLMode: *c.PgSslMode,
	}

	orm, err := postgres.NewClient(&cfg, zap.NewNop())
	if err != nil {
		return Client{}, err
	}
	return Client{db: orm}, nil
}

func (c Client) DB() *gorm.DB {
	return c.db
}
