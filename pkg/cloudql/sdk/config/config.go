package config

import (
	"context"
	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type ClientConfig struct {
	Addresses    []string `cty:"addresses"`
	Username     *string  `cty:"username"`
	Password     *string  `cty:"password"`
	IsOpenSearch *bool    `cty:"is_open_search"`
	AwsRegion    *string  `cty:"aws_region"`

	PgHost     *string `cty:"pg_host"`
	PgPort     *string `cty:"pg_port"`
	PgUser     *string `cty:"pg_user"`
	PgPassword *string `cty:"pg_password"`
	PgDatabase *string `cty:"pg_database"`
	PgSslMode  *string `cty:"pg_ssl_mode"`

	ComplianceServiceBaseURL *string `cty:"compliance_service_baseurl"`
}

func Schema() map[string]*schema.Attribute {
	return map[string]*schema.Attribute{
		"addresses": {
			Type: schema.TypeList,
			Elem: &schema.Attribute{Type: schema.TypeString},
		},
		"username": {
			Type: schema.TypeString,
		},
		"password": {
			Type: schema.TypeString,
		},
		"is_open_search": {
			Type:     schema.TypeBool,
			Required: false,
		},
		"aws_region": {
			Type:     schema.TypeString,
			Required: false,
		},
		"pg_host": {
			Type: schema.TypeString,
		},
		"pg_port": {
			Type: schema.TypeString,
		},
		"pg_user": {
			Type: schema.TypeString,
		},
		"pg_password": {
			Type: schema.TypeString,
		},
		"pg_database": {
			Type: schema.TypeString,
		},
		"pg_ssl_mode": {
			Type: schema.TypeString,
		},
		"compliance_service_baseurl": {
			Type:     schema.TypeString,
			Required: false,
		},
	}
}

func Instance() any {
	return &ClientConfig{}
}

func GetConfig(connection *plugin.Connection) ClientConfig {
	if connection == nil || connection.Config == nil {
		return ClientConfig{}
	}
	config, _ := connection.Config.(ClientConfig)
	return config
}

func NewClientCached(c ClientConfig, cache *connection.ConnectionCache, ctx context.Context) (essdk.Client, error) {
	value, ok := cache.Get(ctx, "opengovernance-es-client")
	if ok {
		return value.(essdk.Client), nil
	}

	plugin.Logger(ctx).Warn("client is not cached, creating a new one")

	client, err := NewClient(ctx, c)
	if err != nil {
		return essdk.Client{}, err
	}

	cache.Set(ctx, "opengovernance-es-client", client)

	return client, nil
}

func NewClient(ctx context.Context, c ClientConfig) (essdk.Client, error) {
	client, err := essdk.NewClient(essdk.ClientConfig{
		Addresses:    c.Addresses,
		Username:     c.Username,
		Password:     c.Password,
		IsOpenSearch: c.IsOpenSearch,
		AwsRegion:    c.AwsRegion,
	})
	if err != nil {
		return essdk.Client{}, err
	}

	return client, nil
}
