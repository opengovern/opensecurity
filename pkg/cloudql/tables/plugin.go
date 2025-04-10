package opengovernance

import (
	"context"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/extra/create_functions"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/extra/utils"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/extra/view-sync"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "cloudql",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: config.Instance,
			Schema:      config.Schema(),
		},
		TableMap: map[string]*plugin.Table{
			"platform_findings":                 tablePlatformFindings(ctx),
			"platform_resources":                tablePlatformResources(ctx),
			"platform_lookup":                   tablePlatformLookup(ctx),
			"platform_integrations":             tablePlatformConnections(ctx),
			"tasks":                             tablePlatformTasks(ctx),
			"platform_integration_groups":       tablePlatformIntegrationGroups(ctx),
			"platform_api_benchmark_summary":    tablePlatformApiBenchmarkSummary(ctx),
			"platform_api_benchmark_controls":   tablePlatformApiBenchmarkControls(ctx),
			"platform_artifact_vulnerabilities": tablePlatformArtifactVulnerabilities(ctx),
			"platform_artifact_sbom":            tablePlatformArtifactSboms(ctx),
			"packages_with_vulnerabilities":     tablePlatformPackageVulnerabilities(ctx),
			"vulnerability_details":             tablePlatformOsvVulnerabilityDetails(ctx),
			"software_packages":                 tableArtifactPackageList(ctx),
			"cve_details":                       tableCveDetails(ctx),
		},
	}

	extraLogger, _ := utils.NewZapLogger()

	viewSync := view_sync.NewViewSync(extraLogger)
	go viewSync.Start(ctx)

	create_functions.CreatePostgresFunctions(ctx, extraLogger)

	return p
}
