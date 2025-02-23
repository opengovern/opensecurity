package opengovernance_client

import (
	"context"
	"runtime"

	"github.com/opengovern/og-util/pkg/integration"
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	InventorySummaryIndex = "inventory_summary"
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LookupResource struct {
	// PlatformID is the unique Global ID of the resource inside the platform
	PlatformID string `json:"platform_id"`
	// ResourceID is the globally unique ID of the resource.
	ResourceID string `json:"resource_id"`
	// ResourceName is the name of the resource.
	ResourceName string `json:"resource_name"`
	// IntegrationType is the type of the integration source of the resource, i.e. AWS Cloud, Azure Cloud.
	IntegrationType integration.Type `json:"integration_type"`
	// ResourceType is the type of the resource.
	ResourceType string `json:"resource_type"`
	// IntegrationID is aws account id or azure subscription id
	IntegrationID string `json:"integration_id"`
	// IsCommon
	IsCommon bool `json:"is_common"`
	// Tags
	Tags []Tag `json:"canonical_tags"`
	// DescribedBy is the resource describe job id
	DescribedBy string `json:"described_by"`
	// DescribedAt is when the DescribeSourceJob is created
	DescribedAt int64 `json:"described_at"`
}

type LookupResourceHit struct {
	ID      string         `json:"_id"`
	Score   float64        `json:"_score"`
	Index   string         `json:"_index"`
	Type    string         `json:"_type"`
	Version int64          `json:"_version,omitempty"`
	Source  LookupResource `json:"_source"`
	Sort    []any          `json:"sort"`
}

type LookupResourceHits struct {
	Total es.SearchTotal      `json:"total"`
	Hits  []LookupResourceHit `json:"hits"`
}

type LookupResourceSearchResponse struct {
	PitID string             `json:"pit_id"`
	Hits  LookupResourceHits `json:"hits"`
}

type LookupResourcePaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewLookupResourcePaginator(filters []es.BoolFilter, limit *int64) (LookupResourcePaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), InventorySummaryIndex, filters, limit)
	if err != nil {
		return LookupResourcePaginator{}, err
	}

	p := LookupResourcePaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p LookupResourcePaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p LookupResourcePaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p LookupResourcePaginator) NextPage(ctx context.Context) ([]LookupResource, error) {
	var response LookupResourceSearchResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []LookupResource
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

var lookupMapping = map[string]string{
	"integration_type": "integration_type",
	"region":           "location",
	"integration_id":   "integration_id",
}

func ListLookupResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListLookupResources", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources NewSelfClientCached", "error", err)
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
	paginator, err := k.NewLookupResourcePaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, lookupMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListLookupResources NewLookupResourcePaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListLookupResources NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListLookupResources", "next page")

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
