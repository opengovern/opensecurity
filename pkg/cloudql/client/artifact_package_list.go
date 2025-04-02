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
	ArtifactPackageListIndex = "artifact_package_list"
)

type ArtifactPackageList struct {
	ImageURL   string    `json:"image_url"`
	ArtifactID string    `json:"artifact_id"`
	Packages   []Package `json:"packages"`
}
type Package struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Version   string `json:"version"`
}

type ArtifactPackageListResult struct {
	PlatformID   string              `json:"platform_id"`
	ResourceID   string              `json:"resource_id"`
	ResourceName string              `json:"resource_name"`
	Description  ArtifactPackageList `json:"Description"`
	TaskType     string              `json:"task_type"`
	ResultType   string              `json:"result_type"`
	Metadata     map[string]string   `json:"metadata"`
	DescribedBy  string              `json:"described_by"`
	DescribedAt  int64               `json:"described_at"`
}

type ArtifactPackageListHit struct {
	ID      string                    `json:"_id"`
	Score   float64                   `json:"_score"`
	Index   string                    `json:"_index"`
	Type    string                    `json:"_type"`
	Version int64                     `json:"_version,omitempty"`
	Source  ArtifactPackageListResult `json:"_source"`
	Sort    []any                     `json:"sort"`
}

type ArtifactPackageListHits struct {
	Total es.SearchTotal           `json:"total"`
	Hits  []ArtifactPackageListHit `json:"hits"`
}

type ArtifactPackageListResponse struct {
	PitID string                  `json:"pit_id"`
	Hits  ArtifactPackageListHits `json:"hits"`
}

type ArtifactPackageListPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewArtifactPackageListPaginator(filters []es.BoolFilter, limit *int64) (ArtifactPackageListPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), ArtifactPackageListIndex, filters, limit)
	if err != nil {
		return ArtifactPackageListPaginator{}, err
	}

	p := ArtifactPackageListPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ArtifactPackageListPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ArtifactPackageListPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ArtifactPackageListPaginator) NextPage(ctx context.Context) ([]ArtifactPackageListResult, error) {
	var response ArtifactPackageListResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []ArtifactPackageListResult
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

var artifactPackageListMapping = map[string]string{
	"image_url":   "Description.ImageURL",
	"artifact_id": "Description.ArtifactID",
}

func ListArtifactPackageList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
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
	paginator, err := k.NewArtifactPackageListPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, artifactPackageListMapping,
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
