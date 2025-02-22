package opengovernance_client

import (
	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	pgsdk "github.com/opengovern/opensecurity/pkg/cloudql/sdk/pg"
)

type Client struct {
	ES essdk.Client
	PG pgsdk.Client
}
