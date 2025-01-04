package compliance

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/opengovern/og-util/pkg/model"
	opengovernanceTypes "github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/compliance/api"
	"github.com/opengovern/opencomply/services/compliance/es"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"go.uber.org/zap"
)

func (h *HttpHandler) getBenchmarkTree(ctx context.Context, benchmarkId string) (*api.NestedBenchmark, error) {
	b, err := h.db.GetBenchmark(ctx, benchmarkId)
	if err != nil {
		return nil, err
	}

	var children []api.NestedBenchmark
	for _, child := range b.Children {
		childNested, err := h.getBenchmarkTree(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		children = append(children, *childNested)
	}

	nb := api.NestedBenchmark{
		ID:                b.ID,
		Title:             b.Title,
		ReferenceCode:     b.DisplayCode,
		Description:       b.Description,
		LogoURI:           b.LogoURI,
		Category:          b.Category,
		DocumentURI:       b.DocumentURI,
		AutoAssign:        b.AutoAssign,
		TracksDriftEvents: b.TracksDriftEvents,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
		Tags:              b.GetTagsMap(),
		Children:          children,
	}
	if b.IntegrationType != nil {
		nb.IntegrationTypes = integration_type.ParseTypes(b.IntegrationType)
	}
	for _, control := range b.Controls {
		nb.Controls = append(nb.Controls, control.ID)
	}
	return &nb, err
}

func (h *HttpHandler) getBenchmarkPath(ctx context.Context, benchmarkId string) (string, error) {
	parent, err := h.db.GetBenchmarkParent(ctx, benchmarkId)
	if err != nil {
		return "", err
	}
	if parent == "" {
		return benchmarkId, nil
	}
	parentPath, err := h.getBenchmarkPath(ctx, parent)
	if err != nil {
		return "", err
	}
	if parentPath == "" {
		return parent, nil
	}
	return parentPath + "/" + benchmarkId, nil
}

func (h *HttpHandler) getBenchmarkComplianceResultSummary(ctx context.Context, benchmarkId string, complianceResultFilters *api.ComplianceResultSummaryFilters) (*api.GetBenchmarkDetailsComplianceResults, error) {
	complianceResults, evaluatedAt, err := es.BenchmarkIntegrationSummary(ctx, h.logger, h.client, benchmarkId)
	if err != nil {
		return nil, err
	}

	var complianceResultsResult api.GetBenchmarkDetailsComplianceResults
	complianceResultsResult.LastEvaluatedAt = time.Unix(evaluatedAt, 0)
	for connection, resultGroup := range complianceResults {
		if complianceResultFilters != nil && len(complianceResultFilters.IntegrationID) > 0 {
			if !listContains(complianceResultFilters.IntegrationID, connection) {
				continue
			}
		}
		if complianceResultFilters != nil && len(complianceResultFilters.ResourceTypeID) > 0 {
			complianceResultsResult.Results = make(map[opengovernanceTypes.ComplianceStatus]int)
			for resourceType, result := range resultGroup.ResourceTypes {
				if listContains(complianceResultFilters.ResourceTypeID, resourceType) {
					for k, v := range result.QueryResult {
						if _, ok := complianceResultsResult.Results[k]; ok {
							complianceResultsResult.Results[k] += v
						} else {
							complianceResultsResult.Results[k] = v
						}
					}
				}
			}
		} else {
			complianceResultsResult.Results = resultGroup.Result.QueryResult
		}
		complianceResultsResult.IntegrationIDs = append(complianceResultsResult.IntegrationIDs, connection)
	}
	return &complianceResultsResult, nil
}

type BenchmarkControlsCache struct {
	Controls map[string]bool
}

// getControlsUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (h *HttpHandler) getControlsUnderBenchmark(ctx context.Context, benchmarkId string, benchmarksCache map[string]BenchmarkControlsCache) (map[string]bool, error) {
	controls := make(map[string]bool)

	benchmark, err := h.db.GetBenchmarkWithControlQueries(ctx, benchmarkId)
	if err != nil {
		h.logger.Error("failed to fetch benchmarks", zap.Error(err))
		return nil, err
	}
	for _, c := range benchmark.Controls {
		controls[c.ID] = true
	}

	for _, child := range benchmark.Children {
		var childControls map[string]bool
		if cache, ok := benchmarksCache[child.ID]; ok {
			childControls = cache.Controls
		} else {
			childControls, err = h.getControlsUnderBenchmark(ctx, child.ID, benchmarksCache)
			if err != nil {
				return nil, err
			}
			benchmarksCache[child.ID] = BenchmarkControlsCache{Controls: childControls}
		}
		for k, _ := range childControls {
			controls[k] = true
		}
	}
	return controls, nil
}

type BenchmarkTablesCache struct {
	PrimaryTables map[string]bool
	ListTables    map[string]bool
}

// getTablesUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (h *HttpHandler) getTablesUnderBenchmark(ctx context.Context, benchmarkId string, benchmarkCache map[string]BenchmarkTablesCache) (map[string]bool, map[string]bool, error) {
	primaryTables := make(map[string]bool)
	listOfTables := make(map[string]bool)

	benchmark, err := h.db.GetBenchmarkWithControlQueries(ctx, benchmarkId)
	if err != nil {
		h.logger.Error("failed to fetch benchmarks", zap.Error(err))
		return nil, nil, err
	}
	for _, c := range benchmark.Controls {
		if c.Policy != nil {
			if c.Policy.PrimaryResource != nil && *c.Policy.PrimaryResource != "" {
				primaryTables[*c.Policy.PrimaryResource] = true
			}
			for _, t := range c.Policy.ListOfResources {
				if t == "" {
					continue
				}
				listOfTables[t] = true
			}
		}
	}

	for _, child := range benchmark.Children {
		var childPrimaryTables, childListOfTables map[string]bool
		if cache, ok := benchmarkCache[child.ID]; ok {
			childPrimaryTables = cache.PrimaryTables
			childListOfTables = cache.ListTables
		} else {
			childPrimaryTables, childListOfTables, err = h.getTablesUnderBenchmark(ctx, child.ID, benchmarkCache)
			if err != nil {
				return nil, nil, err
			}
			benchmarkCache[child.ID] = BenchmarkTablesCache{
				PrimaryTables: childPrimaryTables,
				ListTables:    childListOfTables,
			}
		}

		for k, _ := range childPrimaryTables {
			primaryTables[k] = true
		}
		for k, _ := range childListOfTables {
			childListOfTables[k] = true
		}
	}
	return primaryTables, listOfTables, nil
}

func (h *HttpHandler) getChildBenchmarksWithDetails(ctx context.Context, benchmarkId string, req api.GetBenchmarkDetailsRequest) ([]api.GetBenchmarkDetailsChildren, error) {
	var benchmarks []api.GetBenchmarkDetailsChildren
	benchmark, err := h.db.GetBenchmark(ctx, benchmarkId)
	if err != nil {
		h.logger.Error("failed to fetch benchmarks", zap.Error(err))
		return nil, err
	}
	for _, child := range benchmark.Children {
		var childChildren []api.GetBenchmarkDetailsChildren
		if req.BenchmarkChildren {
			childBenchmarks, err := h.getChildBenchmarksWithDetails(ctx, child.ID, req)
			if err != nil {
				return nil, err
			}
			childChildren = append(childChildren, childBenchmarks...)
		}
		var controlIDs []string
		for _, c := range child.Controls {
			controlIDs = append(controlIDs, c.ID)
		}

		complianceResults, evaluatedAt, err := es.BenchmarkIntegrationSummary(ctx, h.logger, h.client, benchmark.ID)
		if err != nil {
			return nil, err
		}

		var complianceResultsResult api.GetBenchmarkDetailsComplianceResults
		complianceResultsResult.LastEvaluatedAt = time.Unix(evaluatedAt, 0)
		for connection, resultGroup := range complianceResults {
			if req.ComplianceResultFilters != nil && len(req.ComplianceResultFilters.IntegrationID) > 0 {
				if !listContains(req.ComplianceResultFilters.IntegrationID, connection) {
					continue
				}
			}
			if req.ComplianceResultFilters != nil && len(req.ComplianceResultFilters.ResourceTypeID) > 0 {
				complianceResultsResult.Results = make(map[opengovernanceTypes.ComplianceStatus]int)
				for resourceType, result := range resultGroup.ResourceTypes {
					if listContains(req.ComplianceResultFilters.ResourceTypeID, resourceType) {
						for k, v := range result.QueryResult {
							if _, ok := complianceResultsResult.Results[k]; ok {
								complianceResultsResult.Results[k] += v
							} else {
								complianceResultsResult.Results[k] = v
							}
						}
					}
				}
			} else {
				complianceResultsResult.Results = resultGroup.Result.QueryResult
			}
			complianceResultsResult.IntegrationIDs = append(complianceResultsResult.IntegrationIDs, connection)
		}

		benchmarks = append(benchmarks, api.GetBenchmarkDetailsChildren{
			ID:                child.ID,
			Title:             child.Title,
			Tags:              filterTagsByRegex(req.TagsRegex, model.TrimPrivateTags(child.GetTagsMap())),
			ControlIDs:        controlIDs,
			ComplianceResults: complianceResultsResult,
			Children:          childChildren,
		})
	}
	return benchmarks, nil
}

func (h *HttpHandler) getChildBenchmarks(ctx context.Context, benchmarkId string) ([]string, error) {
	var benchmarks []string
	benchmark, err := h.db.GetBenchmark(ctx, benchmarkId)
	if err != nil {
		h.logger.Error("failed to fetch benchmarks", zap.Error(err))
		return nil, err
	}
	if benchmark == nil {
		return nil, fmt.Errorf("benchmark %s not found", benchmarkId)
	}
	for _, child := range benchmark.Children {
		childBenchmarks, err := h.getChildBenchmarks(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		benchmarks = append(benchmarks, childBenchmarks...)
	}
	benchmarks = append(benchmarks, benchmarkId)
	return benchmarks, nil
}

func listContains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// listContainsList list1 > list2
func listContainsList(list1 []string, list2 []string) bool {
	for _, v1 := range list2 {
		if !listContains(list1, v1) {
			return false
		}
	}
	return true
}

func mapToArray(input map[string]bool) []string {
	var result []string
	for k, _ := range input {
		result = append(result, k)
	}
	return result
}

func filterTagsByRegex(regexPattern *string, tags map[string][]string) map[string][]string {
	if regexPattern == nil {
		return tags
	}
	re := regexp.MustCompile(*regexPattern)

	resultsMap := make(map[string][]string)
	for k, v := range tags {
		if re.MatchString(k) {
			resultsMap[k] = v
		}
	}
	return resultsMap
}
