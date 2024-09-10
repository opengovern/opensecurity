package compliance

import (
	"context"
	"github.com/kaytu-io/kaytu-engine/pkg/compliance/api"
	"github.com/kaytu-io/kaytu-engine/pkg/compliance/es"
	kaytuTypes "github.com/kaytu-io/kaytu-engine/pkg/types"
	"github.com/kaytu-io/kaytu-util/pkg/model"
	"go.uber.org/zap"
	"regexp"
	"time"
)

func (h *HttpHandler) getBenchmarkFindingSummary(ctx context.Context, benchmarkId string, findingFilters *api.FindingSummaryFilters) (*api.GetBenchmarkDetailsFindings, error) {
	findings, evaluatedAt, err := es.BenchmarkConnectionSummary(ctx, h.logger, h.client, benchmarkId)
	if err != nil {
		return nil, err
	}

	var findingsResult api.GetBenchmarkDetailsFindings
	findingsResult.LastEvaluatedAt = time.Unix(evaluatedAt, 0)
	for connection, finding := range findings {
		if findingFilters != nil && len(findingFilters.ConnectionID) > 0 {
			if !listContains(findingFilters.ConnectionID, connection) {
				continue
			}
		}
		if findingFilters != nil && len(findingFilters.ResourceTypeID) > 0 {
			findingsResult.Results = make(map[kaytuTypes.ConformanceStatus]int)
			for resourceType, result := range finding.ResourceTypes {
				if listContains(findingFilters.ResourceTypeID, resourceType) {
					for k, v := range result.QueryResult {
						if _, ok := findingsResult.Results[k]; ok {
							findingsResult.Results[k] += v
						} else {
							findingsResult.Results[k] = v
						}
					}
				}
			}
		} else {
			findingsResult.Results = finding.Result.QueryResult
		}
		findingsResult.ConnectionIDs = append(findingsResult.ConnectionIDs, connection)
	}
	return &findingsResult, nil
}

// getTablesUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (h *HttpHandler) getTablesUnderBenchmark(ctx context.Context, benchmarkId string) (map[string]bool, map[string]bool, error) {
	primaryTables := make(map[string]bool)
	listOfTables := make(map[string]bool)

	benchmark, err := h.db.GetBenchmarkWithControlQueries(ctx, benchmarkId)
	if err != nil {
		h.logger.Error("failed to fetch benchmarks", zap.Error(err))
		return nil, nil, err
	}
	for _, c := range benchmark.Controls {
		if c.Query != nil {
			if c.Query.PrimaryTable != nil {
				primaryTables[*c.Query.PrimaryTable] = true
			}
			for _, t := range c.Query.ListOfTables {
				listOfTables[t] = true
			}
		}
	}

	for _, child := range benchmark.Children {
		childPrimaryTables, childListOfTables, err := h.getTablesUnderBenchmark(ctx, child.ID)
		if err != nil {
			return nil, nil, err
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

		findings, evaluatedAt, err := es.BenchmarkConnectionSummary(ctx, h.logger, h.client, benchmark.ID)
		if err != nil {
			return nil, err
		}

		var findingsResult api.GetBenchmarkDetailsFindings
		findingsResult.LastEvaluatedAt = time.Unix(evaluatedAt, 0)
		for connection, finding := range findings {
			if req.FindingFilters != nil && len(req.FindingFilters.ConnectionID) > 0 {
				if !listContains(req.FindingFilters.ConnectionID, connection) {
					continue
				}
			}
			if req.FindingFilters != nil && len(req.FindingFilters.ResourceTypeID) > 0 {
				findingsResult.Results = make(map[kaytuTypes.ConformanceStatus]int)
				for resourceType, result := range finding.ResourceTypes {
					if listContains(req.FindingFilters.ResourceTypeID, resourceType) {
						for k, v := range result.QueryResult {
							if _, ok := findingsResult.Results[k]; ok {
								findingsResult.Results[k] += v
							} else {
								findingsResult.Results[k] = v
							}
						}
					}
				}
			} else {
				findingsResult.Results = finding.Result.QueryResult
			}
			findingsResult.ConnectionIDs = append(findingsResult.ConnectionIDs, connection)
		}

		benchmarks = append(benchmarks, api.GetBenchmarkDetailsChildren{
			ID:         child.ID,
			Title:      child.Title,
			Tags:       filterTagsByRegex(req.TagsRegex, model.TrimPrivateTags(child.GetTagsMap())),
			ControlIDs: controlIDs,
			Findings:   findingsResult,
			Children:   childChildren,
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
