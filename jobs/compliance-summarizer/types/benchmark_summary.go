package types

import (
	"fmt"
	"github.com/opengovern/opensecurity/pkg/utils"

	"github.com/axiomhq/hyperloglog"
	"github.com/opengovern/opensecurity/pkg/types"
)

type Result struct {
	QueryResult    map[types.ComplianceStatus]int
	SeverityResult map[types.ComplianceResultSeverity]int
	SecurityScore  float64
	CostImpact     *float64 `json:"CostImpact,omitempty"`
}

func (r Result) IsFullyPassed() bool {
	for status, count := range r.QueryResult {
		if status.IsPassed() {
			continue
		}
		if count > 0 {
			return false
		}
	}
	return true
}

type ResultGroup struct {
	Result        Result
	ResourceTypes map[string]Result
	Controls      map[string]ControlResult
}

type ControlResult struct {
	Passed bool

	FailedResourcesCount int
	TotalResourcesCount  int

	// these are not exported fields so they are not marshalled
	allResources    *hyperloglog.Sketch
	failedResources *hyperloglog.Sketch

	FailedIntegrationCount int
	TotalIntegrationCount  int

	// these are not exported fields so they are not marshalled
	allIntegrations    *hyperloglog.Sketch
	failedIntegrations *hyperloglog.Sketch

	CostImpact *float64 `json:"CostImpact,omitempty"`
}

type BenchmarkSummaryResult struct {
	BenchmarkResult ResultGroup
	Integrations    map[string]ResultGroup
}

type BenchmarkSummary struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	BenchmarkID      string
	JobID            uint
	EvaluatedAtEpoch int64

	Integrations        BenchmarkSummaryResult
	ResourceCollections map[string]BenchmarkSummaryResult
}

func (b BenchmarkSummary) KeysAndIndex() ([]string, string) {
	keys := []string{
		b.BenchmarkID,
		fmt.Sprintf("%d", b.JobID),
	}
	return keys, types.BenchmarkSummaryIndex
}

func (r *BenchmarkSummaryResult) addComplianceResult(complianceResult types.ComplianceResult) {
	if !complianceResult.ComplianceStatus.IsPassed() {
		r.BenchmarkResult.Result.SeverityResult[complianceResult.Severity]++
	}
	r.BenchmarkResult.Result.QueryResult[complianceResult.ComplianceStatus]++
	r.BenchmarkResult.Result.CostImpact = utils.PAdd(r.BenchmarkResult.Result.CostImpact, complianceResult.CostImpact)

	integration, ok := r.Integrations[complianceResult.IntegrationID]
	if !ok {
		integration = ResultGroup{
			Result: Result{
				QueryResult:    map[types.ComplianceStatus]int{},
				SeverityResult: map[types.ComplianceResultSeverity]int{},
				SecurityScore:  0,
			},
			ResourceTypes: map[string]Result{},
			Controls:      map[string]ControlResult{},
		}
	}
	if !complianceResult.ComplianceStatus.IsPassed() {
		integration.Result.SeverityResult[complianceResult.Severity]++
	}
	integration.Result.QueryResult[complianceResult.ComplianceStatus]++
	integration.Result.CostImpact = utils.PAdd(integration.Result.CostImpact, complianceResult.CostImpact)
	r.Integrations[complianceResult.IntegrationID] = integration

	resourceType, ok := r.BenchmarkResult.ResourceTypes[complianceResult.ResourceType]
	if !ok {
		resourceType = Result{
			QueryResult:    map[types.ComplianceStatus]int{},
			SeverityResult: map[types.ComplianceResultSeverity]int{},
			SecurityScore:  0,
		}
	}
	if !complianceResult.ComplianceStatus.IsPassed() {
		resourceType.SeverityResult[complianceResult.Severity]++
	}
	resourceType.QueryResult[complianceResult.ComplianceStatus]++
	resourceType.CostImpact = utils.PAdd(resourceType.CostImpact, complianceResult.CostImpact)
	r.BenchmarkResult.ResourceTypes[complianceResult.ResourceType] = resourceType

	integrationResourceType, ok := integration.ResourceTypes[complianceResult.ResourceType]
	if !ok {
		integrationResourceType = Result{
			QueryResult:    map[types.ComplianceStatus]int{},
			SeverityResult: map[types.ComplianceResultSeverity]int{},
			SecurityScore:  0,
		}
	}
	if !complianceResult.ComplianceStatus.IsPassed() {
		integrationResourceType.SeverityResult[complianceResult.Severity]++
	}
	integrationResourceType.QueryResult[complianceResult.ComplianceStatus]++
	integrationResourceType.CostImpact = utils.PAdd(integrationResourceType.CostImpact, complianceResult.CostImpact)
	integration.ResourceTypes[complianceResult.ResourceType] = integrationResourceType

	control, ok := r.BenchmarkResult.Controls[complianceResult.ControlID]
	if !ok {
		control = ControlResult{
			Passed:             true,
			allResources:       hyperloglog.New16(),
			failedResources:    hyperloglog.New16(),
			allIntegrations:    hyperloglog.New16(),
			failedIntegrations: hyperloglog.New16(),
		}
	}

	if !complianceResult.ComplianceStatus.IsPassed() {
		control.Passed = false

		control.failedResources.Insert([]byte(complianceResult.PlatformResourceID))
		control.failedIntegrations.Insert([]byte(complianceResult.IntegrationID))
	}
	control.allResources.Insert([]byte(complianceResult.PlatformResourceID))
	control.allIntegrations.Insert([]byte(complianceResult.IntegrationID))
	control.CostImpact = utils.PAdd(control.CostImpact, complianceResult.CostImpact)
	r.BenchmarkResult.Controls[complianceResult.ControlID] = control

	integrationControl, ok := integration.Controls[complianceResult.ControlID]
	if !ok {
		integrationControl = ControlResult{
			Passed:             true,
			allResources:       hyperloglog.New16(),
			failedResources:    hyperloglog.New16(),
			allIntegrations:    hyperloglog.New16(),
			failedIntegrations: hyperloglog.New16(),
		}
	}
	if !complianceResult.ComplianceStatus.IsPassed() {
		integrationControl.Passed = false
		integrationControl.failedResources.Insert([]byte(complianceResult.PlatformResourceID))
		integrationControl.failedIntegrations.Insert([]byte(complianceResult.IntegrationID))
	}
	integrationControl.allResources.Insert([]byte(complianceResult.PlatformResourceID))
	integrationControl.allIntegrations.Insert([]byte(complianceResult.IntegrationID))
	integrationControl.CostImpact = utils.PAdd(integrationControl.CostImpact, complianceResult.CostImpact)
	integration.Controls[complianceResult.ControlID] = integrationControl
}

func (r *BenchmarkSummaryResult) summarize() {
	// update security scores
	for controlID, summary := range r.BenchmarkResult.Controls {
		summary.FailedIntegrationCount = int(summary.failedIntegrations.Estimate())
		summary.TotalIntegrationCount = int(summary.allIntegrations.Estimate())

		summary.FailedResourcesCount = int(summary.failedResources.Estimate())
		summary.TotalResourcesCount = int(summary.allResources.Estimate())

		r.BenchmarkResult.Controls[controlID] = summary
	}

	for resourceType, summary := range r.BenchmarkResult.ResourceTypes {
		total := 0
		for _, count := range summary.QueryResult {
			total += count
		}

		if total > 0 {
			summary.SecurityScore = float64(summary.QueryResult[types.ComplianceStatusOK]) / float64(total) * 100.0
		}

		r.BenchmarkResult.ResourceTypes[resourceType] = summary
	}

	total := 0
	for _, count := range r.BenchmarkResult.Result.QueryResult {
		total += count
	}
	if total > 0 {
		r.BenchmarkResult.Result.SecurityScore = float64(r.BenchmarkResult.Result.QueryResult[types.ComplianceStatusOK]) / float64(total) * 100.0
	}

	for integrationID, summary := range r.Integrations {
		for controlID, controlSummary := range summary.Controls {
			controlSummary.FailedIntegrationCount = int(controlSummary.failedIntegrations.Estimate())
			controlSummary.TotalIntegrationCount = int(controlSummary.allIntegrations.Estimate())

			controlSummary.FailedResourcesCount = int(controlSummary.failedResources.Estimate())
			controlSummary.TotalResourcesCount = int(controlSummary.allResources.Estimate())

			summary.Controls[controlID] = controlSummary
		}

		for resourceType, resourceTypeSummary := range summary.ResourceTypes {
			total := 0
			for _, count := range resourceTypeSummary.QueryResult {
				total += count
			}

			if total > 0 {
				resourceTypeSummary.SecurityScore = float64(resourceTypeSummary.QueryResult[types.ComplianceStatusOK]) / float64(total) * 100.0
			}

			summary.ResourceTypes[resourceType] = resourceTypeSummary
		}

		total := 0
		for _, count := range summary.Result.QueryResult {
			total += count
		}

		if total > 0 {
			summary.Result.SecurityScore = float64(summary.Result.QueryResult[types.ComplianceStatusOK]) / float64(total) * 100.0
		}

		r.Integrations[integrationID] = summary
	}
}

func (b *BenchmarkSummary) summarize() {
	b.Integrations.summarize()
	for rcId, rc := range b.ResourceCollections {
		rc.summarize()
		b.ResourceCollections[rcId] = rc
	}
}
