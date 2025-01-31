package job

import (
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/auth"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/compliance"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/core"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/elasticsearch"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/integration"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/inventory"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/manifest"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/resource_collection"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/resource_info"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/types"
)

var migrations = map[string]types.Migration{
	"elasticsearch": elasticsearch.Migration{},
}
var Order = []string{
	"elasticsearch",
}

var manualMigrations = map[string]types.Migration{
	"manifest":            manifest.Migration{},
	"core":                core.Migration{},
	"integration":         integration.Migration{},
	"inventory":           inventory.Migration{},
	"resource_collection": resource_collection.Migration{},
	"elasticsearch":       elasticsearch.Migration{},
	"compliance":          compliance.Migration{},
	"resource_info":       resource_info.Migration{},
	"auth":                auth.Migration{},
}

// Ordered keys slice
var ManualOrder = []string{
	"manifest",
	"core",
	"integration",
	"inventory",
	"resource_collection",
	"elasticsearch",
	"compliance",
	"resource_info",
	"auth",
}