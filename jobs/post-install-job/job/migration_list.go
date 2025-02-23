package job

import (
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/auth"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/compliance"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/core"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/elasticsearch"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/integration"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/inventory"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/manifest"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/resource_collection"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/resource_info"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/types"
)

var migrations = map[string]types.Migration{
	"elasticsearch": elasticsearch.Migration{},
	"manifest":      manifest.Migration{},
}
var Order = []string{
	"elasticsearch",
	"manifest",
}

var manualMigrations = map[string]types.Migration{
	"elasticsearch":       elasticsearch.Migration{},
	"manifest":            manifest.Migration{},
	"core":                core.Migration{},
	"integration":         integration.Migration{},
	"inventory":           inventory.Migration{},
	"resource_collection": resource_collection.Migration{},
	"compliance":          compliance.Migration{},
	"resource_info":       resource_info.Migration{},
	"auth":                auth.Migration{},
}

// Ordered keys slice
var ManualOrder = []string{
	"elasticsearch",
	"manifest",
	"core",
	"integration",
	"inventory",
	"resource_collection",
	"compliance",
	"resource_info",
	"auth",
}
