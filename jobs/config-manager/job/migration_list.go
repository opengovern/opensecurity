package job

import (
	// Imports kept from both branches or specific to HEAD's intent
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/compliance"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/core"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/elasticsearch"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/integration"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/inventory"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/manifest"

	// Import kept from main as it seems unrelated to the removal
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/plugins"
	// Imports removed as per feat-removing-resource-info
	// "github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/resource_collection"
	// "github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/resource_info"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/migrations/tasks"
	"github.com/opengovern/opensecurity/jobs/config-manager/job/types"
)

// This order seems unchanged between branches in the non-conflicting part
var Order = []string{
	"elasticsearch",
	"manifest",
}

// This map seems unchanged between branches in the non-conflicting part
var migrations = map[string]types.Migration{
	"elasticsearch": elasticsearch.Migration{},
	"manifest":      manifest.Migration{},
}

// manualMigrations resolved: Keep removals from HEAD, add plugins from main
var manualMigrations = map[string]types.Migration{
	"elasticsearch": elasticsearch.Migration{},
	"manifest":      manifest.Migration{},
	"core":          core.Migration{},
	"integration":   integration.Migration{},
	"inventory":     inventory.Migration{},
	// "resource_collection": resource_collection.Migration{}, // Kept commented out / removed from HEAD
	"compliance": compliance.Migration{},
	// "resource_info":       resource_info.Migration{}, // Kept commented out / removed from HEAD
	"tasks":   tasks.Migration{},
	"plugins": plugins.Migration{}, // Added from main
}

// ManualOrder resolved: Keep removals from HEAD, ensure plugins is present
var ManualOrder = []string{
	"elasticsearch",
	"manifest",
	"core",
	"integration",
	"inventory",
	// "resource_collection", // Kept commented out / removed from HEAD
	"compliance",
	// "resource_info", // Kept commented out / removed from HEAD
	"tasks",
	"plugins", // Kept from HEAD/main
}
