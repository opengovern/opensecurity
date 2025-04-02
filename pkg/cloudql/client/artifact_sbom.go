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
	ArtifactSbomsIndex = "artifact_sbom"
)

type ArtifactSbom struct {
	ImageURL          string      `json:"image_url"`
	ArtifactID        string      `json:"artifact_id"`
	Packages          []string    `json:"packages"`
	SbomSpdxJson      interface{} `json:"sbom_spdx_json"`
	SbomCyclonedxJson interface{} `json:"sbom_cyclonedx_json"`
}

type ArtifactSbomResult struct {
	PlatformID   string            `json:"platform_id"`
	ResourceID   string            `json:"resource_id"`
	ResourceName string            `json:"resource_name"`
	Description  ArtifactSbom      `json:"Description"`
	TaskType     string            `json:"task_type"`
	ResultType   string            `json:"result_type"`
	Metadata     map[string]string `json:"metadata"`
	DescribedBy  string            `json:"described_by"`
	DescribedAt  int64             `json:"described_at"`
}

type ArtifactSbomHit struct {
	ID      string             `json:"_id"`
	Score   float64            `json:"_score"`
	Index   string             `json:"_index"`
	Type    string             `json:"_type"`
	Version int64              `json:"_version,omitempty"`
	Source  ArtifactSbomResult `json:"_source"`
	Sort    []any              `json:"sort"`
}

type ArtifactSbomHits struct {
	Total es.SearchTotal    `json:"total"`
	Hits  []ArtifactSbomHit `json:"hits"`
}

type ArtifactSbomSearchResponse struct {
	PitID string           `json:"pit_id"`
	Hits  ArtifactSbomHits `json:"hits"`
}

type ArtifactSbomPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewArtifactSbomPaginator(filters []es.BoolFilter, limit *int64) (ArtifactSbomPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), ArtifactSbomsIndex, filters, limit)
	if err != nil {
		return ArtifactSbomPaginator{}, err
	}

	p := ArtifactSbomPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ArtifactSbomPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ArtifactSbomPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ArtifactSbomPaginator) NextPage(ctx context.Context) ([]ArtifactSbomResult, error) {
	var response ArtifactSbomSearchResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []ArtifactSbomResult
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

var artifactSbomsMapping = map[string]string{
	"image_url":   "Description.ImageURL",
	"artifact_id": "Description.ArtifactID",
	"packages":    "Description.Packages",
}

func ListArtifactSboms(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListArtifactSboms", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactSboms NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactSboms NewSelfClientCached", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactSboms GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.FetchType)
	paginator, err := k.NewArtifactSbomPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, artifactSbomsMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListArtifactSboms NewArtifactSbomPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListArtifactSboms NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListArtifactSboms", "next page")

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
