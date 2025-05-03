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


// This map seems unchanged between branches in the non-conflicting part
var migrations = map[string]types.Migration{}

// This order seems unchanged between branches in the non-conflicting part
var Order = []string{}

// manualMigrations resolved: Keep removals from HEAD, add plugins from main
var manualMigrations = map[string]types.Migration{
	"manifest":    manifest.Migration{},
	"core":        core.Migration{},
	"integration": integration.Migration{},
	"inventory":   inventory.Migration{},
	"compliance":  compliance.Migration{},
	"tasks":       tasks.Migration{},
	"plugins":     plugins.Migration{},
}

// ManualOrder resolved: Keep removals from HEAD, ensure plugins is present
var ManualOrder = []string{
	"manifest",
	"core",
	"integration",
	"inventory",
	"compliance",
	"tasks",
	"plugins",
}
