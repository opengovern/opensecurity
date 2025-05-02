package job

import (
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/compliance"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/core"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/integration"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/inventory"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/manifest"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/plugins"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/tasks"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/types"
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
