package opengovernance_client

import (
	"context"
	"runtime"

	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	PackageVulnerabilitiesIndex = "package_vulnerability"
)

type PackageWithVulnIDs struct {
	PackageIdentifier string   `json:"package_identifier"`
	PackageName       string   `json:"package_name"`
	Ecosystem         string   `json:"ecosystem"`
	Version           string   `json:"version"`
	Vulnerabilities   []string `json:"vulnerabilities"`
}

type PackageWithVulnIDsResult struct {
	PlatformID   string             `json:"platform_id"`
	ResourceID   string             `json:"resource_id"`
	ResourceName string             `json:"resource_name"`
	Description  PackageWithVulnIDs `json:"Description"`
	TaskType     string             `json:"task_type"`
	ResultType   string             `json:"result_type"`
	Metadata     map[string]string  `json:"metadata"`
	DescribedBy  string             `json:"described_by"`
	DescribedAt  int64              `json:"described_at"`
}

type PackageWithVulnIDsHit struct {
	ID      string                   `json:"_id"`
	Score   float64                  `json:"_score"`
	Index   string                   `json:"_index"`
	Type    string                   `json:"_type"`
	Version int64                    `json:"_version,omitempty"`
	Source  PackageWithVulnIDsResult `json:"_source"`
	Sort    []any                    `json:"sort"`
}

type PackageWithVulnIDsHits struct {
	Total es.SearchTotal          `json:"total"`
	Hits  []PackageWithVulnIDsHit `json:"hits"`
}

type PackageWithVulnIDsSearchResponse struct {
	PitID string                 `json:"pit_id"`
	Hits  PackageWithVulnIDsHits `json:"hits"`
}

type PackageWithVulnIDsPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewPackageWithVulnIDsPaginator(filters []es.BoolFilter, limit *int64) (PackageWithVulnIDsPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), PackageVulnerabilitiesIndex, filters, limit)
	if err != nil {
		return PackageWithVulnIDsPaginator{}, err
	}

	p := PackageWithVulnIDsPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p PackageWithVulnIDsPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p PackageWithVulnIDsPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p PackageWithVulnIDsPaginator) NextPage(ctx context.Context) ([]PackageWithVulnIDsResult, error) {
	var response PackageWithVulnIDsSearchResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []PackageWithVulnIDsResult
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var packageWithVulnIDsMapping = map[string]string{
	"image_url":   "Description.ImageURL",
	"artifact_id": "Description.ArtifactID",
	"packages":    "Description.Packages",
}

func ListPackageWithVulnIDs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListPackageWithVulnIDs", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListPackageWithVulnIDs NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListPackageWithVulnIDs NewSelfClientCached", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListPackageWithVulnIDs GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListPackageWithVulnIDs GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.FetchType)
	paginator, err := k.NewPackageWithVulnIDsPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, packageWithVulnIDsMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListPackageWithVulnIDs NewArtifactSbomPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListPackageWithVulnIDs NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListPackageWithVulnIDs", "next page")

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
