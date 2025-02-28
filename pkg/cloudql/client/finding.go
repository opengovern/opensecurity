package opengovernance_client

import (
	"context"
	"runtime"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type FindingHit struct {
	ID      string                 `json:"_id"`
	Score   float64                `json:"_score"`
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Version int64                  `json:"_version,omitempty"`
	Source  types.ComplianceResult `json:"_source"`
	Sort    []any                  `json:"sort"`
}

type FindingHits struct {
	Total es.SearchTotal `json:"total"`
	Hits  []FindingHit   `json:"hits"`
}

type FindingSearchResponse struct {
	PitID string      `json:"pit_id"`
	Hits  FindingHits `json:"hits"`
}

type FindingPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewFindingPaginator(filters []es.BoolFilter, limit *int64) (FindingPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), types.ComplianceResultsIndex, filters, limit)
	if err != nil {
		return FindingPaginator{}, err
	}

	p := FindingPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p FindingPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p FindingPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p FindingPaginator) NextPage(ctx context.Context) ([]types.ComplianceResult, error) {
	var response FindingSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []types.ComplianceResult
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

var listFindingFilters = map[string]string{
	"id":                   "ID",
	"benchmark_id":         "benchmarkID",
	"policy_id":            "policyID",
	"integration_id":       "integrationID",
	"described_at":         "describedAt",
	"evaluated_at":         "evaluatedAt",
	"state_active":         "stateActive",
	"result":               "result",
	"severity":             "severity",
	"evaluator":            "evaluator",
	"integration_type":     "integrationType",
	"platform_resource_id": "platformResourceID",
	"resource_name":        "resourceName",
	"resource_type":        "resourceType",
	"reason":               "reason",
}

func ListFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListFindings")
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListFindings NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	paginator, err := k.NewFindingPaginator(es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, listFindingFilters,
		nil, nil, nil, true), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListFindings NewFindingPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListFindings NextPage", "error", err)
			return nil, err
		}

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
