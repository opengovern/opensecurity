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
	ArtifactVulnerabilitiesIndex = "oci_container_vulnerabilities"
)

type OciArtifactVulnerabilities struct {
	ImageURL        string               `json:"imageUrl"`
	ArtifactDigest  string               `json:"artifactDigest"`
	Vulnerabilities []VulnerabilityMatch `json:"Vulnerabilities"`
}

type GrypeOutput struct {
	Matches []VulnerabilityMatch `json:"matches"`
}

type VulnerabilityMatch struct {
	Vulnerability          Vulnerability   `json:"vulnerability"`
	RelatedVulnerabilities []Vulnerability `json:"relatedVulnerabilities"`
	MatchDetail            interface{}     `json:"matchDetail"`
	Artifact               interface{}     `json:"artifact"`
}

type Vulnerability struct {
	ID          string             `json:"id"`
	DataSource  string             `json:"dataSource"`
	Namespace   string             `json:"namespace"`
	Severity    string             `json:"severity"`
	URLs        []string           `json:"urls"`
	Description string             `json:"Description"`
	CVSs        []VulnerabilityCVS `json:"cvss"`
	Fix         VulnerabilityFix   `json:"fix"`
	Advisories  interface{}        `json:"advisories"`
}

type VulnerabilityCVS struct {
	Source         string            `json:"source"`
	Type           string            `json:"type"`
	Version        string            `json:"version"`
	Vector         string            `json:"vector"`
	Metrics        map[string]string `json:"metrics"`
	VendorMetadata map[string]string `json:"vendorMetadata"`
}

type VulnerabilityFix struct {
	Versions []string `json:"versions"`
	State    string   `json:"state"`
}

type ArtifactVulnerabilitiesTaskResult struct {
	PlatformID   string                     `json:"platform_id"`
	ResourceID   string                     `json:"resource_id"`
	ResourceName string                     `json:"resource_name"`
	Description  OciArtifactVulnerabilities `json:"description"`
	TaskType     string                     `json:"task_type"`
	ResultType   string                     `json:"result_type"`
	Metadata     map[string]string          `json:"metadata"`
	DescribedBy  string                     `json:"described_by"`
	DescribedAt  int64                      `json:"described_at"`
}

type OciArtifactVulnerabilitiesHit struct {
	ID      string                            `json:"_id"`
	Score   float64                           `json:"_score"`
	Index   string                            `json:"_index"`
	Type    string                            `json:"_type"`
	Version int64                             `json:"_version,omitempty"`
	Source  ArtifactVulnerabilitiesTaskResult `json:"_source"`
	Sort    []any                             `json:"sort"`
}

type OciArtifactVulnerabilitiesHits struct {
	Total es.SearchTotal                  `json:"total"`
	Hits  []OciArtifactVulnerabilitiesHit `json:"hits"`
}

type OciArtifactVulnerabilitiesSearchResponse struct {
	PitID string                         `json:"pit_id"`
	Hits  OciArtifactVulnerabilitiesHits `json:"hits"`
}

type OciArtifactVulnerabilitiesPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewOciArtifactVulnerabilitiesPaginator(filters []es.BoolFilter, limit *int64) (OciArtifactVulnerabilitiesPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), ArtifactVulnerabilitiesIndex, filters, limit)
	if err != nil {
		return OciArtifactVulnerabilitiesPaginator{}, err
	}

	p := OciArtifactVulnerabilitiesPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p OciArtifactVulnerabilitiesPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p OciArtifactVulnerabilitiesPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p OciArtifactVulnerabilitiesPaginator) NextPage(ctx context.Context) ([]ArtifactVulnerabilitiesTaskResult, error) {
	var response OciArtifactVulnerabilitiesSearchResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []ArtifactVulnerabilitiesTaskResult
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

var artifactVulnerabilitiesMapping = map[string]string{
	"image_url":       "Description.imageUrl",
	"artifact_digest": "Description.artifactDigest",
}

func ListArtifactVulnerabilities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListArtifactVulnerabilities", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactVulnerabilities NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactVulnerabilities NewSelfClientCached", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.FetchType)
	paginator, err := k.NewOciArtifactVulnerabilitiesPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, artifactVulnerabilitiesMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactVulnerabilities NewOciArtifactVulnerabilitiesPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListArtifactVulnerabilities NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListArtifactVulnerabilities", "next page")

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
