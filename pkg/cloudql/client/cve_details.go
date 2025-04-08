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
	CveID                  string              `json:"cve_id"`
	SourceIdentifier       string              `json:"source_identifier"`
	Published              string              `json:"published"`
	LastModified           string              `json:"last_modified"`
	VulnStatus             string              `json:"vuln_status"`
	Descriptions           []OutputDescription `json:"descriptions"`
	CvssVersion            string              `json:"cvss_version,omitempty"`
	CvssScore              float64             `json:"cvss_score,omitempty"`
	CvssSeverity           string              `json:"cvss_severity,omitempty"`
	CvssAttackVector       string              `json:"cvss_attack_vector,omitempty"`
	CvssAttackComplexity   string              `json:"cvss_attack_complexity,omitempty"`
	CvssPrivilegesRequired string              `json:"cvss_privileges_required,omitempty"`
	CvssUserInteraction    string              `json:"cvss_user_interaction,omitempty"`
	CvssConfImpact         string              `json:"cvss_conf_impact,omitempty"`
	CvssIntegImpact        string              `json:"cvss_integ_impact,omitempty"`
	CvssAvailImpact        string              `json:"cvss_avail_impact,omitempty"`
	Metrics                OutputMetrics       `json:"metrics,omitempty"`
	Weaknesses             []OutputWeakness    `json:"weaknesses"`
	CisaExploitAdd         string              `json:"cisa_exploit_add,omitempty"`
	CisaActionDue          string              `json:"cisa_action_due,omitempty"`
	CisaRequiredAction     string              `json:"cisa_required_action,omitempty"`
	CisaVulnerabilityName  string              `json:"cisa_vulnerability_name,omitempty"`
}

type OutputDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}
type OutputMetrics struct {
	CvssMetricV40 []InputCvssMetricV40 `json:"cvssMetricV40,omitempty"`
	CvssMetricV31 []InputCvssMetricV31 `json:"cvssMetricV31,omitempty"`
	CvssMetricV2  []InputCvssMetricV2  `json:"cvssMetricV2,omitempty"`
}

type InputCvssMetricV2 struct {
	Source                  string          `json:"source"`
	Type                    string          `json:"type"`
	CvssData                InputCvssDataV2 `json:"cvssData"`
	BaseSeverity            string          `json:"baseSeverity"`
	ExploitabilityScore     float64         `json:"exploitabilityScore"`
	ImpactScore             float64         `json:"impactScore"`
	AcInsufInfo             bool            `json:"acInsufInfo"`
	ObtainAllPrivilege      bool            `json:"obtainAllPrivilege"`
	ObtainUserPrivilege     bool            `json:"obtainUserPrivilege"`
	ObtainOtherPrivilege    bool            `json:"obtainOtherPrivilege"`
	UserInteractionRequired bool            `json:"userInteractionRequired"`
}
type InputCvssDataV2 struct {
	Version               string  `json:"version"`
	VectorString          string  `json:"vectorString"`
	AccessVector          string  `json:"accessVector"`
	AccessComplexity      string  `json:"accessComplexity"`
	Authentication        string  `json:"authentication"`
	ConfidentialityImpact string  `json:"confidentialityImpact"`
	IntegrityImpact       string  `json:"integrityImpact"`
	AvailabilityImpact    string  `json:"availabilityImpact"`
	BaseScore             float64 `json:"baseScore"`
}

type OutputWeakness struct {
	Source      string              `json:"source"`
	Type        string              `json:"type"`
	Description []OutputDescription `json:"description"`
}

type InputCvssMetricV40 struct {
	Source   string           `json:"source"`
	Type     string           `json:"type"`
	CvssData InputCvssDataV40 `json:"cvssData"`
}
type InputCvssDataV40 struct {
	Version                           string  `json:"version"`
	VectorString                      string  `json:"vectorString"`
	BaseScore                         float64 `json:"baseScore"`
	BaseSeverity                      string  `json:"baseSeverity"`
	AttackVector                      string  `json:"attackVector"`
	AttackComplexity                  string  `json:"attackComplexity"`
	AttackRequirements                string  `json:"attackRequirements"`
	PrivilegesRequired                string  `json:"privilegesRequired"`
	UserInteraction                   string  `json:"userInteraction"`
	VulnConfidentialityImpact         string  `json:"vulnConfidentialityImpact"`
	VulnIntegrityImpact               string  `json:"vulnIntegrityImpact"`
	VulnAvailabilityImpact            string  `json:"vulnAvailabilityImpact"`
	SubConfidentialityImpact          string  `json:"subConfidentialityImpact"`
	SubIntegrityImpact                string  `json:"subIntegrityImpact"`
	SubAvailabilityImpact             string  `json:"subAvailabilityImpact"`
	ExploitMaturity                   string  `json:"exploitMaturity"`
	ConfidentialityRequirement        string  `json:"confidentialityRequirement"`
	IntegrityRequirement              string  `json:"integrityRequirement"`
	AvailabilityRequirement           string  `json:"availabilityRequirement"`
	ModifiedAttackVector              string  `json:"modifiedAttackVector"`
	ModifiedAttackComplexity          string  `json:"modifiedAttackComplexity"`
	ModifiedAttackRequirements        string  `json:"modifiedAttackRequirements"`
	ModifiedPrivilegesRequired        string  `json:"modifiedPrivilegesRequired"`
	ModifiedUserInteraction           string  `json:"modifiedUserInteraction"`
	ModifiedVulnConfidentialityImpact string  `json:"modifiedVulnConfidentialityImpact"`
	ModifiedVulnIntegrityImpact       string  `json:"modifiedVulnIntegrityImpact"`
	ModifiedVulnAvailabilityImpact    string  `json:"modifiedVulnAvailabilityImpact"`
	ModifiedSubConfidentialityImpact  string  `json:"modifiedSubConfidentialityImpact"`
	ModifiedSubIntegrityImpact        string  `json:"modifiedSubIntegrityImpact"`
	ModifiedSubAvailabilityImpact     string  `json:"modifiedSubAvailabilityImpact"`
	Safety                            string  `json:"Safety"`
	Automatable                       string  `json:"Automatable"`
	Recovery                          string  `json:"Recovery"`
	ValueDensity                      string  `json:"valueDensity"`
	VulnerabilityResponseEffort       string  `json:"vulnerabilityResponseEffort"`
	ProviderUrgency                   string  `json:"providerUrgency"`
}
type InputCvssMetricV31 struct {
	Source              string           `json:"source"`
	Type                string           `json:"type"`
	CvssData            InputCvssDataV31 `json:"cvssData"`
	ExploitabilityScore float64          `json:"exploitabilityScore"`
	ImpactScore         float64          `json:"impactScore"`
}
type InputCvssDataV31 struct {
	Version               string  `json:"version"`
	VectorString          string  `json:"vectorString"`
	AttackVector          string  `json:"attackVector"`
	AttackComplexity      string  `json:"attackComplexity"`
	PrivilegesRequired    string  `json:"privilegesRequired"`
	UserInteraction       string  `json:"userInteraction"`
	Scope                 string  `json:"scope"`
	ConfidentialityImpact string  `json:"confidentialityImpact"`
	IntegrityImpact       string  `json:"integrityImpact"`
	AvailabilityImpact    string  `json:"availabilityImpact"`
	BaseScore             float64 `json:"baseScore"`
	BaseSeverity          string  `json:"baseSeverity"`
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
	"cve_id":                   "CveID",
	"source_identifier":        "SourceIdentifier",
	"published":                "Published",
	"last_modified":            "LastModified",
	"vuln_status":              "VulnStatus",
	"cvss_version":             "CvssVersion",
	"cvss_score":               "CvssScore",
	"cvss_severity":            "CvssSeverity",
	"cvss_attack_vector":       "CvssAttackVector",
	"cvss_attack_complexity":   "CvssAttackComplexity",
	"cvss_privileges_required": "CvssPrivilegesRequired",
	"cvss_user_interaction":    "CvssUserInteraction",
	"cvss_conf_impact":         "CvssConfImpact",
	"cvss_integ_impact":        "CvssIntegImpact",
	"cvss_avail_impact":        "CvssAvailImpact",
	"cisa_exploit_add":         "CisaExploitAdd",
	"cisa_action_due":          "CisaActionDue",
	"cisa_required_action":     "CisaRequiredAction",
	"cisa_vulnerability_name":  "CisaVulnerabilityName",
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
