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
	CveDetailsIndex = "cve_details"
)

type CveDetails struct {
	ID                    string           `json:"id"` // UPPERCASE
	SourceIdentifier      string           `json:"source_identifier"`
	Published             string           `json:"published"`
	LastModified          string           `json:"last_modified"`
	VulnStatus            string           `json:"vuln_status"`
	Description           string           `json:"description"`          // Single string
	Metrics               []interface{}    `json:"metrics,omitempty"`    // Flat array
	Weaknesses            []TargetWeakness `json:"weaknesses,omitempty"` // Array, might be empty
	CisaExploitAdd        *string          `json:"cisa_exploit_add,omitempty"`
	CisaActionDue         *string          `json:"cisa_action_due,omitempty"`
	CisaRequiredAction    *string          `json:"cisa_required_action,omitempty"`
	CisaVulnerabilityName *string          `json:"cisa_vulnerability_name,omitempty"`
}

type TargetWeakness struct {
	Source      string             `json:"source"`
	Type        string             `json:"type"`
	Description []InputDescription `json:"description"` // Re-use InputDescription
}

type InputDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type CveDetailsResult struct {
	PlatformID   string            `json:"platform_id"`
	ResourceID   string            `json:"resource_id"`
	ResourceName string            `json:"resource_name"`
	Description  CveDetails        `json:"Description"`
	TaskType     string            `json:"task_type"`
	ResultType   string            `json:"result_type"`
	Metadata     map[string]string `json:"metadata"`
	DescribedBy  string            `json:"described_by"`
	DescribedAt  int64             `json:"described_at"`
}

type CveDetailsHit struct {
	ID      string           `json:"_id"`
	Score   float64          `json:"_score"`
	Index   string           `json:"_index"`
	Type    string           `json:"_type"`
	Version int64            `json:"_version,omitempty"`
	Source  CveDetailsResult `json:"_source"`
	Sort    []any            `json:"sort"`
}

type CveDetailsHits struct {
	Total es.SearchTotal  `json:"total"`
	Hits  []CveDetailsHit `json:"hits"`
}

type CveDetailsResponse struct {
	PitID string         `json:"pit_id"`
	Hits  CveDetailsHits `json:"hits"`
}

type CveDetailsPaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewCveDetailsPaginator(filters []es.BoolFilter, limit *int64) (CveDetailsPaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), CveDetailsIndex, filters, limit)
	if err != nil {
		return CveDetailsPaginator{}, err
	}

	p := CveDetailsPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p CveDetailsPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p CveDetailsPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p CveDetailsPaginator) NextPage(ctx context.Context) ([]CveDetailsResult, error) {
	var response CveDetailsResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
	if err != nil {
		return nil, err
	}

	var values []CveDetailsResult
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

var cveDetailsMapping = map[string]string{
	"id":                      "Description.ID",
	"source_identifier":       "Description.SourceIdentifier",
	"published":               "Description.Published",
	"last_modified":           "Description.LastModified",
	"vuln_status":             "Description.VulnStatus",
	"description":             "Description.Description",
	"cisa_exploit_add":        "Description.CisaExploitAdd",
	"cisa_action_due":         "Description.CisaActionDue",
	"cisa_required_action":    "Description.CisaRequiredAction",
	"cisa_vulnerability_name": "Description.CisaVulnerabilityName",
}

func ListCveDetails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListCveDetails", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListCveDetails NewClientCached", "error", err)
		return nil, err
	}
	k := Client{ES: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListCveDetails NewSelfClientCached", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListCveDetails GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListCveDetails GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.FetchType)
	paginator, err := k.NewCveDetailsPaginator(
		es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, cveDetailsMapping,
			nil, encodedResourceCollectionFilters, clientType, true),
		d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListCveDetails NewArtifactSbomPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListCveDetails NextPage", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Trace("ListCveDetails", "next page")

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
