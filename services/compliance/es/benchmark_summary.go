package es

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	types2 "github.com/opengovern/opencomply/jobs/compliance-summarizer-job/types"
	"github.com/opengovern/opencomply/pkg/types"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"go.uber.org/zap"
)

type BenchmarkTrendDatapoint struct {
	DateEpoch      int64
	QueryResult    map[types.ComplianceStatus]int
	SeverityResult map[types.ComplianceResultSeverity]int
	Controls       map[string]types2.ControlResult
}

func (t *BenchmarkTrendDatapoint) addResultGroupToTrendDataPoint(resultGroup types2.ResultGroup) {
	for k, v := range resultGroup.Result.QueryResult {
		t.QueryResult[k] += v
	}
	for k, v := range resultGroup.Result.SeverityResult {
		t.SeverityResult[k] += v
	}
	for controlId, control := range resultGroup.Controls {
		if _, ok := t.Controls[controlId]; !ok {
			t.Controls[controlId] = types2.ControlResult{
				Passed: true,
			}
		}
		v := t.Controls[controlId]
		v.FailedResourcesCount += control.FailedResourcesCount
		v.TotalResourcesCount += control.TotalResourcesCount
		v.FailedIntegrationCount += control.FailedIntegrationCount
		v.TotalIntegrationCount += control.TotalIntegrationCount
		v.Passed = v.Passed && control.Passed
		t.Controls[controlId] = v
	}
}

type ControlTrendDatapoint struct {
	DateEpoch              int64
	FailedResourcesCount   int
	TotalResourcesCount    int
	FailedIntegrationCount int
	TotalIntegrationCount  int
}

type FetchBenchmarkSummaryTrendAggregatedResponse struct {
	Aggregations struct {
		BenchmarkIDGroup struct {
			Buckets []struct {
				Key                   string `json:"key"`
				EvaluatedAtRangeGroup struct {
					Buckets []struct {
						From      float64 `json:"from"`
						To        float64 `json:"to"`
						DocCount  int     `json:"doc_count"`
						HitSelect struct {
							Hits struct {
								Hits []struct {
									Source types2.BenchmarkSummary `json:"_source"`
								} `json:"hits"`
							} `json:"hits"`
						} `json:"hit_select"`
					} `json:"buckets"`
				} `json:"evaluated_at_range_group"`
			} `json:"buckets"`
		} `json:"benchmark_id_group"`
	} `json:"aggregations"`
}

func FetchBenchmarkSummaryTrendByIntegrationID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkIDs []string, integrationIDs []string, from, to time.Time) (map[string][]BenchmarkTrendDatapoint, error) {
	pathFilters := make([]string, 0, len(integrationIDs)+4)
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.key")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.from")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.to")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.BenchmarkResult.Result")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.BenchmarkResult.Controls")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Result", integrationID),
			fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID),
		)
	}

	query := make(map[string]any)

	startTimeUnix := from.Truncate(24 * time.Hour).Unix()
	endTimeUnix := to.Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()
	step := int64((time.Hour * 24).Seconds())
	ranges := make([]map[string]any, 0, (endTimeUnix-startTimeUnix)/int64(step))
	for i := int64(0); i*step < endTimeUnix-startTimeUnix; i++ {
		ranges = append(ranges, map[string]any{
			"from": startTimeUnix + i*step,
			"to":   startTimeUnix + (i+1)*step - 1,
		})
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": to.Unix(),
				"gte": from.Unix(),
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}
	query["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}
	query["size"] = 0

	query["aggs"] = map[string]any{
		"benchmark_id_group": map[string]any{
			"terms": map[string]any{
				"field": "BenchmarkID",
				"size":  10000,
			},
			"aggs": map[string]any{
				"evaluated_at_range_group": map[string]any{
					"range": map[string]any{
						"field":  "EvaluatedAtEpoch",
						"ranges": ranges,
					},
					"aggs": map[string]any{
						"hit_select": map[string]any{
							"top_hits": map[string]any{
								"sort": map[string]any{
									"EvaluatedAtEpoch": "desc",
								},
								"size": 1,
							},
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err))
		return nil, err
	}

	logger.Info("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var response FetchBenchmarkSummaryTrendAggregatedResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &response)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
		return nil, err
	}

	trend := make(map[string][]BenchmarkTrendDatapoint)
	for _, bucket := range response.Aggregations.BenchmarkIDGroup.Buckets {
		benchmarkID := bucket.Key
		for _, rangeBucket := range bucket.EvaluatedAtRangeGroup.Buckets {
			date := int64(rangeBucket.To)
			if err != nil {
				logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
				return nil, err
			}
			trendDataPoint := BenchmarkTrendDatapoint{
				QueryResult:    make(map[types.ComplianceStatus]int),
				SeverityResult: make(map[types.ComplianceResultSeverity]int),
				Controls:       make(map[string]types2.ControlResult),
			}

			for _, hit := range rangeBucket.HitSelect.Hits.Hits {
				trendDataPoint.DateEpoch = date
				if len(integrationIDs) > 0 {
					for _, integrationID := range integrationIDs {
						if integration, ok := hit.Source.Integrations.Integrations[integrationID]; ok {
							trendDataPoint.addResultGroupToTrendDataPoint(integration)
						}
					}
				} else {
					trendDataPoint.addResultGroupToTrendDataPoint(hit.Source.Integrations.BenchmarkResult)
				}
			}
			if trendDataPoint.DateEpoch != 0 {
				trend[benchmarkID] = append(trend[benchmarkID], trendDataPoint)
			}
		}
		sort.Slice(trend[benchmarkID], func(i, j int) bool {
			return trend[benchmarkID][i].DateEpoch < trend[benchmarkID][j].DateEpoch
		})
	}

	return trend, nil
}

func FetchBenchmarkSummaryTrendByIntegrationIDV3(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	benchmarkIDs []string, integrationIDs []string, from, to, granularity int64) (map[string][]BenchmarkTrendDatapoint, error) {
	pathFilters := make([]string, 0, len(integrationIDs)+4)
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.key")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.from")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.to")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.BenchmarkResult.Result")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.BenchmarkResult.Controls")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Result", integrationID),
			fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID),
		)
	}

	query := make(map[string]any)

	startTimeUnix := from
	endTimeUnix := to
	ranges := make([]map[string]any, 0, (endTimeUnix-startTimeUnix)/int64(granularity))
	for i := int64(0); i*granularity < endTimeUnix-startTimeUnix; i++ {
		ranges = append(ranges, map[string]any{
			"from": startTimeUnix + i*granularity,
			"to":   startTimeUnix + (i+1)*granularity - 1,
		})
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": to,
				"gte": from,
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}
	query["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}
	query["size"] = 0

	query["aggs"] = map[string]any{
		"benchmark_id_group": map[string]any{
			"terms": map[string]any{
				"field": "BenchmarkID",
				"size":  10000,
			},
			"aggs": map[string]any{
				"evaluated_at_range_group": map[string]any{
					"range": map[string]any{
						"field":  "EvaluatedAtEpoch",
						"ranges": ranges,
					},
					"aggs": map[string]any{
						"hit_select": map[string]any{
							"top_hits": map[string]any{
								"sort": map[string]any{
									"EvaluatedAtEpoch": "desc",
								},
								"size": 1,
							},
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err))
		return nil, err
	}

	logger.Info("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var response FetchBenchmarkSummaryTrendAggregatedResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &response)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
		return nil, err
	}

	trend := make(map[string][]BenchmarkTrendDatapoint)
	trendMap := make(map[string]map[int64]BenchmarkTrendDatapoint)
	for _, bucket := range response.Aggregations.BenchmarkIDGroup.Buckets {
		benchmarkID := bucket.Key
		trendMap[benchmarkID] = make(map[int64]BenchmarkTrendDatapoint)
		for _, rangeBucket := range bucket.EvaluatedAtRangeGroup.Buckets {
			date := int64(rangeBucket.To)
			if err != nil {
				logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
				return nil, err
			}
			trendDataPoint := BenchmarkTrendDatapoint{
				QueryResult:    make(map[types.ComplianceStatus]int),
				SeverityResult: make(map[types.ComplianceResultSeverity]int),
				Controls:       make(map[string]types2.ControlResult),
			}

			for _, hit := range rangeBucket.HitSelect.Hits.Hits {
				trendDataPoint.DateEpoch = date
				if len(integrationIDs) > 0 {
					for _, integrationID := range integrationIDs {
						if integration, ok := hit.Source.Integrations.Integrations[integrationID]; ok {
							trendDataPoint.addResultGroupToTrendDataPoint(integration)
						}
					}
				} else {
					trendDataPoint.addResultGroupToTrendDataPoint(hit.Source.Integrations.BenchmarkResult)
				}
			}
			if trendDataPoint.DateEpoch != 0 {
				trendMap[benchmarkID][date] = trendDataPoint
			}
		}
		for i := int64(0); i*granularity < endTimeUnix-startTimeUnix; i++ {
			t := startTimeUnix + (i+1)*granularity - 1
			if v, ok := trendMap[benchmarkID][t]; ok {
				trend[benchmarkID] = append(trend[benchmarkID], v)
			} else {
				trend[benchmarkID] = append(trend[benchmarkID], BenchmarkTrendDatapoint{
					DateEpoch: t,
				})
			}
		}
		sort.Slice(trend[benchmarkID], func(i, j int) bool {
			return trend[benchmarkID][i].DateEpoch < trend[benchmarkID][j].DateEpoch
		})
	}

	return trend, nil
}

func FetchBenchmarkSummaryTrendByResourceCollectionAndIntegrationID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkIDs []string, integrationIDs []string, resourceCollections []string, from, to time.Time) (map[string][]BenchmarkTrendDatapoint, error) {
	if len(resourceCollections) == 0 {
		return nil, fmt.Errorf("resource collections cannot be empty")
	}
	pathFilters := make([]string, 0, (len(integrationIDs)+1)*len(resourceCollections)+3)
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.key")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.from")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.to")
	for _, resourceCollection := range resourceCollections {
		pathFilters = append(pathFilters, fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.ResourceCollections.%s.BenchmarkResult.Result.SecurityScore", resourceCollection))
		for _, integrationID := range integrationIDs {
			pathFilters = append(pathFilters,
				fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.ResourceCollections.%s.Integrations.%s.Result.SecurityScore", resourceCollection, integrationID))
		}
	}

	query := make(map[string]any)

	startTimeUnix := from.Truncate(24 * time.Hour).Unix()
	endTimeUnix := to.Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()
	step := int64((time.Hour * 24).Seconds())
	ranges := make([]map[string]any, 0, (endTimeUnix-startTimeUnix)/int64(step))
	for i := int64(0); i*step < endTimeUnix-startTimeUnix; i++ {
		ranges = append(ranges, map[string]any{
			"from": startTimeUnix + i*step,
			"to":   startTimeUnix + (i+1)*step - 1,
		})
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": to.Unix(),
				"gte": from.Unix(),
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}
	query["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}
	query["size"] = 0

	query["aggs"] = map[string]any{
		"benchmark_id_group": map[string]any{
			"terms": map[string]any{
				"field": "BenchmarkID",
				"size":  10000,
			},
			"aggs": map[string]any{
				"evaluated_at_range_group": map[string]any{
					"range": map[string]any{
						"field":  "EvaluatedAtEpoch",
						"ranges": ranges,
					},
					"aggs": map[string]any{
						"hit_select": map[string]any{
							"top_hits": map[string]any{
								"sort": map[string]any{
									"EvaluatedAtEpoch": "desc",
								},
								"size": 1,
							},
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err))
		return nil, err
	}

	logger.Info("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var response FetchBenchmarkSummaryTrendAggregatedResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &response)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
		return nil, err
	}

	trend := make(map[string][]BenchmarkTrendDatapoint)
	for _, bucket := range response.Aggregations.BenchmarkIDGroup.Buckets {
		benchmarkID := bucket.Key
		for _, rangeBucket := range bucket.EvaluatedAtRangeGroup.Buckets {
			date := int64(rangeBucket.To)
			if err != nil {
				logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDAtTime", zap.Error(err), zap.String("query", string(queryBytes)))
				return nil, err
			}
			trendDataPoint := BenchmarkTrendDatapoint{
				QueryResult:    make(map[types.ComplianceStatus]int),
				SeverityResult: make(map[types.ComplianceResultSeverity]int),
				Controls:       make(map[string]types2.ControlResult),
			}
			for _, hit := range rangeBucket.HitSelect.Hits.Hits {
				trendDataPoint.DateEpoch = date
				for _, resourceCollection := range hit.Source.ResourceCollections {
					if len(integrationIDs) > 0 {
						for _, integrationID := range integrationIDs {
							if integration, ok := resourceCollection.Integrations[integrationID]; ok {
								trendDataPoint.addResultGroupToTrendDataPoint(integration)
							}
						}
					} else {
						trendDataPoint.addResultGroupToTrendDataPoint(resourceCollection.BenchmarkResult)
					}
				}
			}
			if trendDataPoint.DateEpoch != 0 {
				trend[benchmarkID] = append(trend[benchmarkID], trendDataPoint)
			}
		}
		sort.Slice(trend[benchmarkID], func(i, j int) bool {
			return trend[benchmarkID][i].DateEpoch < trend[benchmarkID][j].DateEpoch
		})
	}

	return trend, nil
}

func FetchBenchmarkSummaryTrend(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkIDs []string, integrationIDs, resourceCollections []string, from, to time.Time) (map[string][]BenchmarkTrendDatapoint, error) {
	if len(resourceCollections) > 0 {
		return FetchBenchmarkSummaryTrendByResourceCollectionAndIntegrationID(ctx, logger, client, benchmarkIDs, integrationIDs, resourceCollections, from, to)
	}
	return FetchBenchmarkSummaryTrendByIntegrationID(ctx, logger, client, benchmarkIDs, integrationIDs, from, to)
}

func FetchBenchmarkSummaryTrendByIntegrationIDPerControl(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	benchmarkIDs []string, controlIDs []string, integrationIDs []string, from, to time.Time, stepDuration time.Duration) (map[string][]ControlTrendDatapoint, error) {
	pathFilters := make([]string, 0, len(integrationIDs)+4)
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.key")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.from")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.to")
	if len(integrationIDs) > 0 {
		if len(controlIDs) > 0 {
			for _, integrationID := range integrationIDs {
				for _, controlID := range controlIDs {
					pathFilters = append(pathFilters,
						fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Controls.%s", integrationID, controlID))
				}
			}
		} else {
			for _, integrationID := range integrationIDs {
				pathFilters = append(pathFilters,
					fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID))
			}
		}
	} else {
		if len(controlIDs) > 0 {
			for _, controlID := range controlIDs {
				pathFilters = append(pathFilters,
					fmt.Sprintf("aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.*.Controls.%s", controlID))
			}
		} else {
			pathFilters = append(pathFilters,
				"aggregations.benchmark_id_group.buckets.evaluated_at_range_group.buckets.hit_select.hits.hits._source.Integrations.Integrations.*.Controls")
		}
	}

	query := make(map[string]any)

	startTimeUnix := from.Truncate(stepDuration).Unix()
	endTimeUnix := to.Truncate(stepDuration).Add(stepDuration).Unix()
	step := int64(stepDuration.Seconds())
	ranges := make([]map[string]any, 0, (endTimeUnix-startTimeUnix)/int64(step))
	for i := int64(0); i*step < endTimeUnix-startTimeUnix; i++ {
		ranges = append(ranges, map[string]any{
			"from": startTimeUnix + i*step,
			"to":   startTimeUnix + (i+1)*step,
		})
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": to.Unix(),
				"gte": from.Unix(),
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}
	query["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}
	query["size"] = 0

	query["aggs"] = map[string]any{
		"benchmark_id_group": map[string]any{
			"terms": map[string]any{
				"field": "BenchmarkID",
				"size":  10000,
			},
			"aggs": map[string]any{
				"evaluated_at_range_group": map[string]any{
					"range": map[string]any{
						"field":  "EvaluatedAtEpoch",
						"ranges": ranges,
					},
					"aggs": map[string]any{
						"hit_select": map[string]any{
							"top_hits": map[string]any{
								"sort": map[string]any{
									"EvaluatedAtEpoch": "desc",
								},
								"size": 1,
							},
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDPerControl", zap.Error(err))
		return nil, err
	}

	logger.Info("FetchBenchmarkSummaryTrendByIntegrationIDPerControl", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var response FetchBenchmarkSummaryTrendAggregatedResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &response)
	if err != nil {
		logger.Error("FetchBenchmarkSummaryTrendByIntegrationIDPerControl", zap.Error(err), zap.String("query", string(queryBytes)))
		return nil, err
	}

	trendMap := make(map[string]map[int64]ControlTrendDatapoint)
	currentTimes := make(map[string]map[int64]int64)
	for _, bucket := range response.Aggregations.BenchmarkIDGroup.Buckets {
		for _, rangeBucket := range bucket.EvaluatedAtRangeGroup.Buckets {
			for _, hit := range rangeBucket.HitSelect.Hits.Hits {
				controlData := make(map[string][]ControlTrendDatapoint)
				for _, integration := range hit.Source.Integrations.Integrations {
					for controlId, control := range integration.Controls {
						trendDataPoint := ControlTrendDatapoint{}
						trendDataPoint.DateEpoch = int64(rangeBucket.To)
						trendDataPoint.FailedResourcesCount = control.FailedResourcesCount
						trendDataPoint.TotalResourcesCount = control.TotalResourcesCount
						trendDataPoint.FailedIntegrationCount = control.FailedIntegrationCount
						trendDataPoint.TotalIntegrationCount = control.TotalIntegrationCount
						controlData[controlId] = append(controlData[controlId], trendDataPoint)
					}
				}

				for controlId, controlTrendDataPoints := range controlData {
					if _, ok := trendMap[controlId]; !ok {
						trendMap[controlId] = make(map[int64]ControlTrendDatapoint)
						currentTimes[controlId] = make(map[int64]int64)
					}
					if _, ok := trendMap[controlId][int64(rangeBucket.To)]; !ok || currentTimes[controlId][int64(rangeBucket.To)] < hit.Source.EvaluatedAtEpoch {
						trendMap[controlId][int64(rangeBucket.To)] = ControlTrendDatapoint{
							DateEpoch:              int64(rangeBucket.To),
							FailedResourcesCount:   0,
							TotalResourcesCount:    0,
							FailedIntegrationCount: 0,
							TotalIntegrationCount:  0,
						}
						currentTimes[controlId][int64(rangeBucket.To)] = hit.Source.EvaluatedAtEpoch
						for _, controlTrendDataPoint := range controlTrendDataPoints {
							v := trendMap[controlId][int64(rangeBucket.To)]
							v.FailedResourcesCount += controlTrendDataPoint.FailedResourcesCount
							v.TotalResourcesCount += controlTrendDataPoint.TotalResourcesCount
							v.FailedIntegrationCount += controlTrendDataPoint.FailedIntegrationCount
							v.TotalIntegrationCount += controlTrendDataPoint.TotalIntegrationCount
							trendMap[controlId][int64(rangeBucket.To)] = v
						}
					}
				}
			}
		}
	}

	trend := make(map[string][]ControlTrendDatapoint)
	for controlId, trendMap := range trendMap {
		for _, trendDataPoint := range trendMap {
			trend[controlId] = append(trend[controlId], trendDataPoint)
		}
		sort.Slice(trend[controlId], func(i, j int) bool {
			return trend[controlId][i].DateEpoch < trend[controlId][j].DateEpoch
		})
	}

	return trend, nil
}

type ListBenchmarkSummariesAtTimeResponse struct {
	Aggregations struct {
		Summaries struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				LastResult struct {
					Hits struct {
						Hits []struct {
							Source types2.BenchmarkSummary `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"last_result"`
			} `json:"buckets"`
		} `json:"summaries"`
	} `json:"aggregations"`
}

func ListBenchmarkSummariesAtTime(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	benchmarkIDs []string,
	integrationIDs []string, resourceCollections []string,
	timeAt time.Time, fetchFullObject bool) (map[string]types2.BenchmarkSummary, error) {

	idx := types.BenchmarkSummaryIndex

	includes := []string{"Integrations.BenchmarkResult.Result", "EvaluatedAtEpoch", "Integrations.BenchmarkResult.Controls"}
	if len(integrationIDs) > 0 || fetchFullObject {
		includes = append(includes, "Integrations.Integrations")
	}
	if len(resourceCollections) > 0 || fetchFullObject {
		includes = append(includes, "ResourceCollections")
	}
	pathFilters := make([]string, 0, len(integrationIDs)+(len(resourceCollections)*(len(integrationIDs)+1))+2)
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.key")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.EvaluatedAtEpoch")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Result")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Controls")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.Integrations.%s.Result", integrationID))
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID))
	}
	for _, resourceCollection := range resourceCollections {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.BenchmarkResult.Result", resourceCollection))
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Controls", resourceCollection))
		for _, integrationID := range integrationIDs {
			pathFilters = append(pathFilters,
				fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Integrations.%s.Result", resourceCollection, integrationID))
			pathFilters = append(pathFilters,
				fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Integrations.%s.Controls", resourceCollection, integrationID))
		}
	}

	request := map[string]any{
		"aggs": map[string]any{
			"summaries": map[string]any{
				"terms": map[string]any{
					"field": "BenchmarkID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"last_result": map[string]any{
						"top_hits": map[string]any{
							"sort": []map[string]any{
								{
									"JobID": "desc",
								},
							},
							"_source": map[string]any{
								"includes": includes,
							},
							"size": 1,
						},
					},
				},
			},
		},
		"size": 0,
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": timeAt.Unix(),
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}

	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchBenchmarkSummariesByIntegrationIDAtTime", zap.String("query", string(query)), zap.String("index", idx))

	var response ListBenchmarkSummariesAtTimeResponse
	if fetchFullObject {
		err = client.Search(ctx, idx, string(query), &response)
	} else {
		err = client.SearchWithFilterPath(ctx, idx, string(query), pathFilters, &response)
	}
	if err != nil {
		return nil, err
	}

	benchmarkSummaries := make(map[string]types2.BenchmarkSummary)
	for _, summary := range response.Aggregations.Summaries.Buckets {
		for _, hit := range summary.LastResult.Hits.Hits {
			benchmarkSummaries[summary.Key] = hit.Source
		}
	}
	return benchmarkSummaries, nil
}

func GetComplianceSummaryByJobId(ctx context.Context, logger *zap.Logger, client opengovernance.Client, summaryJobIDs []string, fetchFullObject bool) (map[string]types2.BenchmarkSummary, error) {

	idx := types.BenchmarkSummaryIndex

	includes := []string{"Integrations.BenchmarkResult.Result", "EvaluatedAtEpoch", "Integrations.BenchmarkResult.Controls"}
	if fetchFullObject {
		includes = append(includes, "Integrations.Integrations")
	}
	if fetchFullObject {
		includes = append(includes, "ResourceCollections")
	}
	var pathFilters []string
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.key")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.EvaluatedAtEpoch")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Result")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Controls")

	request := map[string]any{
		"aggs": map[string]any{
			"summaries": map[string]any{
				"terms": map[string]any{
					"field": "BenchmarkID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"last_result": map[string]any{
						"top_hits": map[string]any{
							"sort": []map[string]any{
								{
									"JobID": "desc",
								},
							},
							"_source": map[string]any{
								"includes": includes,
							},
							"size": 1,
						},
					},
				},
			},
		},
		"size": 0,
	}

	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": map[string]any{
				"terms": map[string][]string{
					"JobID": summaryJobIDs,
				},
			},
		},
	}

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchBenchmarkSummariesByIntegrationIDAtTime", zap.String("query", string(query)), zap.String("index", idx))

	var response ListBenchmarkSummariesAtTimeResponse
	if fetchFullObject {
		err = client.Search(ctx, idx, string(query), &response)
	} else {
		err = client.SearchWithFilterPath(ctx, idx, string(query), pathFilters, &response)
	}
	if err != nil {
		return nil, err
	}

	benchmarkSummaries := make(map[string]types2.BenchmarkSummary)
	for _, summary := range response.Aggregations.Summaries.Buckets {
		for _, hit := range summary.LastResult.Hits.Hits {
			benchmarkSummaries[summary.Key] = hit.Source
		}
	}
	return benchmarkSummaries, nil
}

type BenchmarkSummaryResponse struct {
	Aggregations struct {
		LastResult struct {
			Hits struct {
				Hits []struct {
					Source types2.BenchmarkSummary `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"last_result"`
	} `json:"aggregations"`
}

func BenchmarkIntegrationSummary(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkID string) (map[string]types2.ResultGroup, int64, error) {
	includes := []string{"Integrations.Integrations", "EvaluatedAtEpoch"}
	pathFilters := make([]string, 0, 2)
	pathFilters = append(pathFilters, "aggregations.last_result.hits.hits._source.Integrations.Integrations.*.Result")
	pathFilters = append(pathFilters, "aggregations.last_result.hits.hits._source.EvaluatedAtEpoch")
	request := map[string]any{
		"aggs": map[string]any{
			"last_result": map[string]any{
				"top_hits": map[string]any{
					"sort": []map[string]any{
						{
							"JobID": "desc",
						},
					},
					"_source": map[string]any{
						"includes": includes,
					},
					"size": 1,
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"term": map[string]any{
							"BenchmarkID": benchmarkID,
						},
					},
				},
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, -1, err
	}

	logger.Info("BenchmarkIntegrationSummary", zap.String("query", string(queryBytes)))
	var resp BenchmarkSummaryResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &resp)
	if err != nil {
		return nil, -1, err
	}

	for _, res := range resp.Aggregations.LastResult.Hits.Hits {
		return res.Source.Integrations.Integrations, res.Source.EvaluatedAtEpoch, nil
	}
	return nil, 0, nil
}

func BenchmarkControlSummary(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkID string, integrationIDs []string, timeAt time.Time) (map[string]types2.ControlResult, int64, error) {
	includes := []string{"Integrations.BenchmarkResult.Controls", "EvaluatedAtEpoch"}
	if len(integrationIDs) > 0 {
		includes = append(includes, "Integrations.Integrations")
	}

	pathFilters := make([]string, 0, len(integrationIDs)+2)
	pathFilters = append(pathFilters, "aggregations.last_result.hits.hits._source.Integrations.BenchmarkResult.Controls")
	pathFilters = append(pathFilters, "aggregations.last_result.hits.hits._source.EvaluatedAtEpoch")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.last_result.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID))
	}

	request := map[string]any{
		"aggs": map[string]any{
			"last_result": map[string]any{
				"top_hits": map[string]any{
					"sort": []map[string]any{
						{
							"JobID": "desc",
						},
					},
					"_source": map[string]any{
						"includes": includes,
					},
					"size": 1,
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"term": map[string]any{
							"BenchmarkID": benchmarkID,
						},
					},
					{
						"range": map[string]any{
							"EvaluatedAtEpoch": map[string]any{
								"lte": timeAt.Unix(),
							},
						},
					},
				},
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, -1, err
	}

	logger.Info("BenchmarkControlSummary", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var resp BenchmarkSummaryResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex,
		string(queryBytes), pathFilters, &resp)
	if err != nil {
		return nil, -1, err
	}

	evAt := int64(0)
	result := make(map[string]types2.ControlResult)
	for _, res := range resp.Aggregations.LastResult.Hits.Hits {
		if len(integrationIDs) > 0 {
			for _, integrationID := range integrationIDs {
				if integration, ok := res.Source.Integrations.Integrations[integrationID]; ok {
					for key, controlRes := range integration.Controls {
						if v, ok := result[key]; !ok {
							result[key] = controlRes
						} else {
							v.FailedResourcesCount += controlRes.FailedResourcesCount
							v.FailedIntegrationCount += controlRes.FailedIntegrationCount
							v.TotalResourcesCount += controlRes.TotalResourcesCount
							v.TotalIntegrationCount += controlRes.TotalIntegrationCount
							v.Passed = v.Passed && controlRes.Passed
							result[key] = v
						}
					}
					evAt = res.Source.EvaluatedAtEpoch
				}
			}
		} else {
			result = res.Source.Integrations.BenchmarkResult.Controls
			evAt = res.Source.EvaluatedAtEpoch
			break
		}
	}
	return result, evAt, nil
}

type BenchmarksControlSummaryResponse struct {
	Aggregations struct {
		BenchmarkIDGroup struct {
			Buckets []struct {
				Key        string `json:"key"`
				LastResult struct {
					Hits struct {
						Hits []struct {
							Source types2.BenchmarkSummary `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"last_result"`
			} `json:"buckets"`
		} `json:"benchmark_id_group"`
	} `json:"aggregations"`
}

func BenchmarksControlSummary(ctx context.Context, logger *zap.Logger, client opengovernance.Client, benchmarkIDs []string, integrationIDs []string) (map[string]types2.ControlResult, map[string]int64, error) {
	includes := []string{"Integrations.BenchmarkResult.Controls", "EvaluatedAtEpoch"}
	if len(integrationIDs) > 0 {
		includes = append(includes, "Integrations.Integrations")
	}

	pathFilters := make([]string, 0, len(integrationIDs)+2)
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.key")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.last_result.hits.hits._source.EvaluatedAtEpoch")
	pathFilters = append(pathFilters, "aggregations.benchmark_id_group.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Controls")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.benchmark_id_group.buckets.last_result.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID))
	}

	request := map[string]any{
		"aggs": map[string]any{
			"benchmark_id_group": map[string]any{
				"terms": map[string]any{
					"field": "BenchmarkID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"last_result": map[string]any{
						"top_hits": map[string]any{
							"sort": []map[string]any{
								{
									"JobID": "desc",
								},
							},
							"_source": map[string]any{
								"includes": includes,
							},
							"size": 1,
						},
					},
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"terms": map[string][]string{
							"BenchmarkID": benchmarkIDs,
						},
					},
				},
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, nil, err
	}

	logger.Info("BenchmarksControlSummary", zap.String("query", string(queryBytes)), zap.String("pathFilters", strings.Join(pathFilters, ",")))
	var resp BenchmarksControlSummaryResponse
	err = client.SearchWithFilterPath(ctx, types.BenchmarkSummaryIndex, string(queryBytes), pathFilters, &resp)
	if err != nil {
		logger.Error("BenchmarksControlSummary", zap.Error(err))
		return nil, nil, err
	}

	perBenchmarkEvAt := make(map[string]int64)
	perBenchmarkResult := make(map[string]map[string]types2.ControlResult)
	for _, bucket := range resp.Aggregations.BenchmarkIDGroup.Buckets {
		benchmarkID := bucket.Key
		perBenchmarkEvAt[benchmarkID] = int64(0)
		perBenchmarkResult[benchmarkID] = make(map[string]types2.ControlResult)
		for _, res := range bucket.LastResult.Hits.Hits {
			if len(integrationIDs) > 0 {
				for _, integrationID := range integrationIDs {
					if integration, ok := res.Source.Integrations.Integrations[integrationID]; ok {
						for key, controlRes := range integration.Controls {
							if v, ok := perBenchmarkResult[benchmarkID][key]; !ok {
								perBenchmarkResult[benchmarkID][key] = controlRes
							} else {
								v.FailedResourcesCount += controlRes.FailedResourcesCount
								v.FailedIntegrationCount += controlRes.FailedIntegrationCount
								v.TotalResourcesCount += controlRes.TotalResourcesCount
								v.TotalIntegrationCount += controlRes.TotalIntegrationCount
								v.Passed = v.Passed && controlRes.Passed
								perBenchmarkResult[benchmarkID][key] = v
							}
						}
						perBenchmarkEvAt[benchmarkID] = res.Source.EvaluatedAtEpoch
					}
				}
			} else {
				perBenchmarkResult[benchmarkID] = res.Source.Integrations.BenchmarkResult.Controls
				perBenchmarkEvAt[benchmarkID] = res.Source.EvaluatedAtEpoch
				break
			}
		}
	}

	evAt := make(map[string]int64)
	result := make(map[string]types2.ControlResult)
	for benchmarkID, controlResults := range perBenchmarkResult {
		for key, controlRes := range controlResults {
			if _, ok := result[key]; !ok || evAt[key] < perBenchmarkEvAt[benchmarkID] {
				result[key] = controlRes
				evAt[key] = perBenchmarkEvAt[benchmarkID]
			}
		}
	}

	return result, evAt, nil
}

func ListJobsSummariesAtTime(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	benchmarkIDs []string, jobIDs []string,
	integrationIDs []string, resourceCollections []string,
	timeAt time.Time, fetchFullObject bool) (map[string]types2.BenchmarkSummary, error) {

	idx := types.BenchmarkSummaryIndex

	includes := []string{"Integrations.BenchmarkResult.Result", "EvaluatedAtEpoch", "Integrations.BenchmarkResult.Controls"}
	if len(integrationIDs) > 0 || fetchFullObject {
		includes = append(includes, "Integrations.Integrations")
	}
	if len(resourceCollections) > 0 || fetchFullObject {
		includes = append(includes, "ResourceCollections")
	}
	pathFilters := make([]string, 0, len(integrationIDs)+(len(resourceCollections)*(len(integrationIDs)+1))+2)
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.key")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.EvaluatedAtEpoch")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Result")
	pathFilters = append(pathFilters, "aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.BenchmarkResult.Controls")
	for _, integrationID := range integrationIDs {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.Integrations.%s.Result", integrationID))
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.Integrations.Integrations.%s.Controls", integrationID))
	}
	for _, resourceCollection := range resourceCollections {
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.BenchmarkResult.Result", resourceCollection))
		pathFilters = append(pathFilters,
			fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Controls", resourceCollection))
		for _, integrationID := range integrationIDs {
			pathFilters = append(pathFilters,
				fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Integrations.%s.Result", resourceCollection, integrationID))
			pathFilters = append(pathFilters,
				fmt.Sprintf("aggregations.summaries.buckets.last_result.hits.hits._source.ResourceCollections.%s.Integrations.%s.Controls", resourceCollection, integrationID))
		}
	}

	request := map[string]any{
		"aggs": map[string]any{
			"summaries": map[string]any{
				"terms": map[string]any{
					"field": "BenchmarkID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"last_result": map[string]any{
						"top_hits": map[string]any{
							"sort": []map[string]any{
								{
									"JobID": "desc",
								},
							},
							"_source": map[string]any{
								"includes": includes,
							},
							"size": 1,
						},
					},
				},
			},
		},
		"size": 0,
	}

	filters := make([]any, 0)
	filters = append(filters, map[string]any{
		"range": map[string]any{
			"EvaluatedAtEpoch": map[string]any{
				"lte": timeAt.Unix(),
			},
		},
	})
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": benchmarkIDs,
			},
		})
	}
	if len(jobIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"BenchmarkID": jobIDs,
			},
		})
	}

	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchBenchmarkSummariesByIntegrationIDAtTime", zap.String("query", string(query)), zap.String("index", idx))

	var response ListBenchmarkSummariesAtTimeResponse
	if fetchFullObject {
		err = client.Search(ctx, idx, string(query), &response)
	} else {
		err = client.SearchWithFilterPath(ctx, idx, string(query), pathFilters, &response)
	}
	if err != nil {
		return nil, err
	}

	jobsSummaries := make(map[string]types2.BenchmarkSummary)
	for _, summary := range response.Aggregations.Summaries.Buckets {
		for _, hit := range summary.LastResult.Hits.Hits {
			jobsSummaries[strconv.Itoa(int(hit.Source.JobID))] = hit.Source
		}
	}
	return jobsSummaries, nil
}
