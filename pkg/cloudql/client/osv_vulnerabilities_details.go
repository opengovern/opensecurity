package opengovernance_client

import (
	"context"
	"encoding/json"
	"runtime"

	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	OsvVulnerabilitiesDetailsIndex = "osv_vulnerability_detail"
)

type OsvVulnerabilityDetail struct {
	SchemaVersion    string          `json:"schema_version,omitempty"`
	ID               string          `json:"id"`
	Modified         string          `json:"modified"`
	Published        string          `json:"published,omitempty"`
	Withdrawn        string          `json:"withdrawn,omitempty"`
	Aliases          []string        `json:"aliases,omitempty"`
	Related          []string        `json:"related,omitempty"`
	Summary          string          `json:"summary,omitempty"`
	Details          string          `json:"details,omitempty"`
	Severity         []interface{}   `json:"severity,omitempty"`
	Affected         []interface{}   `json:"affected"`
	References       []interface{}   `json:"references,omitempty"`
	Credits          []interface{}   `json:"credits,omitempty"`
	DatabaseSpecific json.RawMessage `json:"database_specific,omitempty"`
}

type OsvVulnerabilityDetailResult struct {
	PlatformID   string                 `json:"platform_id"`
	ResourceID   string                 `json:"resource_id"`
	ResourceName string                 `json:"resource_name"`
	Description  OsvVulnerabilityDetail `json:"Description"`
	TaskType     string                 `json:"task_type"`
	ResultType   string                 `json:"result_type"`
	Metadata     map[string]string      `json:"metadata"`
	DescribedBy  string                 `json:"described_by"`
	DescribedAt  int64                  `json:"described_at"`
}

type OsvVulnerabilityDetailHit struct {
	ID      string                       `json:"_id"`
	Score   float64                      `json:"_score"`
	Index   string                       `json:"_index"`
	Type    string                       `json:"_type"`
	Version int64                        `json:"_version,omitempty"`
	Source  OsvVulnerabilityDetailResult `json:"_source"`
	Sort    []any                        `json:"sort"`
}

type OsvVulnerabilityDetailHits struct {
	Total es.SearchTotal              `json:"total"`
	Hits  []OsvVulnerabilityDetailHit `json:"hits"`
}

type OsvVulnerabilityDetailResponse struct {
	PitID string                     `json:"pit_id"`
	Hits  OsvVulnerabilityDetailHits `json:"hits"`
}

type OsvVulnerabilityDetailPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewOsvVulnerabilityDetailPaginator(filters []es.BoolFilter, limit *int64) (OsvVulnerabilityDetailPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), OsvVulnerabilitiesDetailsIndex, filters, limit)
	if err != nil {
		return OsvVulnerabilityDetailPaginator{}, err
	}

	p := OsvVulnerabilityDetailPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p OsvVulnerabilityDetailPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p OsvVulnerabilityDetailPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p OsvVulnerabilityDetailPaginator) NextPage(ctx context.Context) ([]OsvVulnerabilityDetailResult, error) {
	var response OsvVulnerabilityDetailResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []OsvVulnerabilityDetailResult
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

var osvVulnerabilityDetailMapping = map[string]string{
	"schema_version": "Description.SchemaVersion",
	"id":             "Description.ID",
	"modified":       "Description.Modified",
	"published":      "Description.Published",
	"withdrawn":      "Description.Withdrawn",
	"summary":        "Description.Summary",
	"details":        "Description.Details",
}

func ListOsvVulnerabilityDetail(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListOsvVulnerabilityDetail", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail NewSelfClientCached", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.FetchType)
	paginator, err := k.NewOsvVulnerabilityDetailPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, osvVulnerabilityDetailMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail NewOsvVulnerabilityDetailPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListOsvVulnerabilityDetail NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListOsvVulnerabilityDetail", "next page")

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
