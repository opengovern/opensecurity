package es

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/compliance/api"
	"go.uber.org/zap"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

type ComplianceResultsQueryResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal  `json:"total"`
		Hits  []ComplianceResultsQueryHit `json:"hits"`
	} `json:"hits"`
	PitID string `json:"pit_id"`
}

type ComplianceResultsQueryHit struct {
	ID      string                 `json:"_id"`
	Score   float64                `json:"_score"`
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Version int64                  `json:"_version,omitempty"`
	Source  types.ComplianceResult `json:"_source"`
	Sort    []any                  `json:"sort"`
}

type ComplianceResultPaginator struct {
	paginator *opengovernance.BaseESPaginator
}

func NewComplianceResultPaginator(client opengovernance.Client, idx string, filters []opengovernance.BoolFilter, limit *int64, sort []map[string]any) (ComplianceResultPaginator, error) {
	paginator, err := opengovernance.NewPaginatorWithSort(client.ES(), idx, filters, limit, sort)
	if err != nil {
		return ComplianceResultPaginator{}, err
	}

	p := ComplianceResultPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ComplianceResultPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ComplianceResultPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ComplianceResultPaginator) NextPage(ctx context.Context) ([]types.ComplianceResult, error) {
	var response ComplianceResultsQueryResponse
	err := p.paginator.SearchWithLog(ctx, &response, true)
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

type ComplianceResultsCountQueryHit struct {
	Aggregations struct {
		ControlIDCount struct {
			Buckets []struct {
				Key                   string `json:"key"`
				DocCount              int64  `json:"doc_count"`
				ComplianceStatusCount struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int64  `json:"doc_count"`
					} `json:"buckets"`
				} `json:"complianceStatus_count"`
			} `json:"buckets"`
		} `json:"controlID_count"`
	} `json:"aggregations"`
}

func ComplianceResultsCountByControlID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, resourceIDs []string, integrationTypes []string, integrationID []string, notIntegrationID []string, resourceTypes []string, benchmarkID []string, controlID []string, severity []types.ComplianceResultSeverity, lastTransitionFrom *time.Time, lastTransitionTo *time.Time, evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time, stateActive []bool, complianceStatuses []types.ComplianceStatus) (map[string]map[string]int64, error) {
	idx := types.ComplianceResultsIndex
	var filters []opengovernance.BoolFilter
	if len(resourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceID", resourceIDs))
	}
	if len(resourceTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceType", resourceTypes))
	}
	if len(benchmarkID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("parentBenchmarks", benchmarkID))
	}
	if len(controlID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("controlID", controlID))
	}
	if len(severity) > 0 {
		strSeverity := make([]string, 0)
		for _, s := range severity {
			strSeverity = append(strSeverity, string(s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("severity", strSeverity))
	}
	if len(complianceStatuses) > 0 {
		strComplianceStatus := make([]string, 0)
		for _, cr := range complianceStatuses {
			strComplianceStatus = append(strComplianceStatus, string(cr))
		}
		filters = append(filters, opengovernance.NewTermsFilter("complianceStatus", strComplianceStatus))
	}
	if len(integrationID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationID", integrationID))
	}
	if len(notIntegrationID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", notIntegrationID)))
	}
	if len(integrationTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
	}
	if lastTransitionFrom != nil && lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	} else if lastTransitionFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", ""))
	} else if lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", "",
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	}
	if evaluatedAtFrom != nil && evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	} else if evaluatedAtFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", ""))
	} else if evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", "",
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	}

	query := map[string]any{
		"size": 0,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
		"aggs": map[string]any{
			"controlID_count": map[string]any{
				"terms": map[string]any{
					"field": "controlID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"complianceStatus_count": map[string]any{
						"terms": map[string]any{
							"field": "complianceStatus",
							"size":  10,
						},
					},
				},
			},
		},
	}

	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	logger.Info("ComplianceResultsCountByControlID", zap.String("query", string(queryJson)), zap.String("index", idx))

	var response ComplianceResultsCountQueryHit
	err = client.SearchWithTrackTotalHits(ctx, idx, string(queryJson), nil, &response, true)
	if err != nil {
		return nil, err
	}

	controlIDCount := make(map[string]map[string]int64)
	for _, bucket := range response.Aggregations.ControlIDCount.Buckets {
		controlIDCount[bucket.Key] = make(map[string]int64)
		for _, complianceBucket := range bucket.ComplianceStatusCount.Buckets {
			controlIDCount[bucket.Key][complianceBucket.Key] = complianceBucket.DocCount
		}
	}

	return controlIDCount, nil
}

func ComplianceResultsQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client, resourceIDs []string, integrationTypes []string,
	integrationID []string, notIntegrationID []string, resourceTypes []string, benchmarkID []string, controlID []string,
	severity []types.ComplianceResultSeverity, lastTransitionFrom *time.Time, lastTransitionTo *time.Time,
	evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time, stateActive []bool, complianceStatuses []types.ComplianceStatus,
	sorts []api.ComplianceResultsSort, pageSizeLimit int, searchAfter []any, jobIDs []string) ([]ComplianceResultsQueryHit, int64, error) {
	idx := types.ComplianceResultsIndex

	requestSort := make([]map[string]any, 0, len(sorts)+1)
	for _, sort := range sorts {
		switch {
		case sort.IntegrationType != nil:
			requestSort = append(requestSort, map[string]any{
				"integrationType": *sort.IntegrationType,
			})
		case sort.PlatformResourceID != nil:
			requestSort = append(requestSort, map[string]any{
				"platformResourceID": *sort.PlatformResourceID,
			})
		case sort.ResourceID != nil:
			requestSort = append(requestSort, map[string]any{
				"resourceID": *sort.ResourceID,
			})
		case sort.ResourceTypeID != nil:
			requestSort = append(requestSort, map[string]any{
				"resourceType": *sort.ResourceTypeID,
			})
		case sort.IntegrationID != nil:
			requestSort = append(requestSort, map[string]any{
				"integrationID": *sort.IntegrationID,
			})
		case sort.BenchmarkID != nil:
			requestSort = append(requestSort, map[string]any{
				"benchmarkID": *sort.BenchmarkID,
			})
		case sort.ControlID != nil:
			requestSort = append(requestSort, map[string]any{
				"controlID": *sort.ControlID,
			})
		case sort.Severity != nil:
			scriptSource :=
				`if (params['_source']['severity'] == 'critical') {
					return 5
				} else if (params['_source']['severity'] == 'high') {
					return 4
				} else if (params['_source']['severity'] == 'medium') {
					return 3
				} else if (params['_source']['severity'] == 'low') {
					return 2
				} else if (params['_source']['severity'] == 'none') {
					return 1
				} else {
					return 1
				}`
			requestSort = append(requestSort, map[string]any{
				"_script": map[string]any{
					"type": "number",
					"script": map[string]any{
						"lang":   "painless",
						"source": scriptSource,
					},
					"order": *sort.Severity,
				},
			})
		case sort.ComplianceStatus != nil:
			scriptSource :=
				`if (params['_source']['complianceStatus'] == 'alarm') {
					return 5
				} else if (params['_source']['complianceStatus'] == 'error') {
					return 4
				} else if (params['_source']['complianceStatus'] == 'info') {
					return 3
				} else if (params['_source']['complianceStatus'] == 'skip') {
					return 2
				} else if (params['_source']['complianceStatus'] == 'ok') {
					return 1
				} else {
					return 1
				}`
			requestSort = append(requestSort, map[string]any{
				"_script": map[string]any{
					"type": "number",
					"script": map[string]any{
						"lang":   "painless",
						"source": scriptSource,
					},
					"order": *sort.ComplianceStatus,
				},
			})
		case sort.StateActive != nil:
			requestSort = append(requestSort, map[string]any{
				"stateActive": *sort.StateActive,
			})
		}
	}
	requestSort = append(requestSort, map[string]any{
		"_id": "asc",
	})

	var filters []opengovernance.BoolFilter
	if len(resourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceID", resourceIDs))
	}
	if len(resourceTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceType", resourceTypes))
	}
	if len(benchmarkID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("benchmarkID", benchmarkID))
	}
	if len(controlID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("controlID", controlID))
	}
	if len(jobIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("complianceJobID", jobIDs))
	}
	if len(severity) > 0 {
		strSeverity := make([]string, 0)
		for _, s := range severity {
			strSeverity = append(strSeverity, string(s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("severity", strSeverity))
	}
	if len(complianceStatuses) > 0 {
		strComplianceStatus := make([]string, 0)
		for _, cr := range complianceStatuses {
			strComplianceStatus = append(strComplianceStatus, string(cr))
		}
		filters = append(filters, opengovernance.NewTermsFilter("complianceStatus", strComplianceStatus))
	}
	if len(integrationID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationID", integrationID))
	}
	if len(notIntegrationID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", notIntegrationID)))
	}
	if len(integrationTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
	}
	if lastTransitionFrom != nil && lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	} else if lastTransitionFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", ""))
	} else if lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", "",
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	}
	if evaluatedAtFrom != nil && evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	} else if evaluatedAtFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", ""))
	} else if evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", "",
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	}

	query := make(map[string]any)
	if len(filters) > 0 {
		query["query"] = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}
	query["sort"] = requestSort
	if len(searchAfter) > 0 {
		query["search_after"] = searchAfter
	}
	if pageSizeLimit == 0 {
		pageSizeLimit = 1000
	}
	query["size"] = pageSizeLimit
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	logger.Info("ComplianceResultsQuery", zap.String("query", string(queryJson)), zap.String("index", idx))

	var response ComplianceResultsQueryResponse
	err = client.SearchWithTrackTotalHits(ctx, idx, string(queryJson), nil, &response, true)
	if err != nil {
		return nil, 0, err
	}

	return response.Hits.Hits, response.Hits.Total.Value, err
}

type ComplianceResultsCountResponse struct {
	Hits  ComplianceResultsCountHits `json:"hits"`
	PitID string                     `json:"pit_id"`
}
type ComplianceResultsCountHits struct {
	Total opengovernance.SearchTotal `json:"total"`
}

func ComplianceResultsCount(ctx context.Context, client opengovernance.Client, complianceStatuses []types.ComplianceStatus, stateActive []bool) (int64, error) {
	idx := types.ComplianceResultsIndex

	filters := make([]map[string]any, 0)
	if len(complianceStatuses) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"complianceStatus": complianceStatuses,
			},
		})
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"stateActive": strStateActive,
			},
		})
	}

	query := make(map[string]any)
	query["size"] = 0

	if len(filters) > 0 {
		query["query"] = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}

	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	var response ComplianceResultsCountResponse
	err = client.SearchWithTrackTotalHits(ctx, idx, string(queryJson), nil, &response, true)
	if err != nil {
		return 0, err
	}

	return response.Hits.Total.Value, err
}

type AggregationResult struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}

func (a AggregationResult) GetBucketsKeys() []string {
	var keys []string
	for _, bucket := range a.Buckets {
		keys = append(keys, bucket.Key)
	}
	return keys
}

type ComplianceResultFiltersAggregationResponse struct {
	Aggregations struct {
		ControlIDFilter          AggregationResult `json:"control_id_filter"`
		SeverityFilter           AggregationResult `json:"severity_filter"`
		IntegrationTypeFilter    AggregationResult `json:"integration_type_filter"`
		IntegrationIDFilter      AggregationResult `json:"integration_id_filter"`
		BenchmarkIDFilter        AggregationResult `json:"benchmark_id_filter"`
		ResourceTypeFilter       AggregationResult `json:"resource_type_filter"`
		ResourceCollectionFilter AggregationResult `json:"resource_collection_filter"`
		ComplianceStatusFilter   AggregationResult `json:"compliance_status_filter"`
		StateActiveFilter        struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				KeyAsString string `json:"key_as_string"`
				DocCount    int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"state_active_filter"`
	} `json:"aggregations"`
}

func ComplianceResultsFiltersQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	resourceIDs []string, integrationTypes []string, integrationID []string, notIntegrationID []string,
	resourceTypes []string, benchmarkID []string, controlID []string, severity []types.ComplianceResultSeverity,
	lastTransitionFrom *time.Time, lastTransitionTo *time.Time,
	evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time,
	stateActive []bool, complianceStatuses []types.ComplianceStatus,
) (*ComplianceResultFiltersAggregationResponse, error) {
	idx := types.ComplianceResultsIndex

	var filters []opengovernance.BoolFilter
	if len(resourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceID", resourceIDs))
	}
	if len(resourceTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceType", resourceTypes))
	}
	if len(benchmarkID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("parentBenchmarks", benchmarkID))
	}
	if len(controlID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("controlID", controlID))
	}
	if len(severity) > 0 {
		strSeverity := make([]string, 0)
		for _, s := range severity {
			strSeverity = append(strSeverity, string(s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("severity", strSeverity))
	}
	if len(complianceStatuses) > 0 {
		strComplianceStatus := make([]string, 0)
		for _, cr := range complianceStatuses {
			strComplianceStatus = append(strComplianceStatus, string(cr))
		}
		filters = append(filters, opengovernance.NewTermsFilter("complianceStatus", strComplianceStatus))
	}
	if len(integrationID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationID", integrationID))
	}
	if len(notIntegrationID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", notIntegrationID)))
	}
	if len(integrationTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
	}
	if lastTransitionFrom != nil && lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	} else if lastTransitionFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", ""))
	} else if lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", "",
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	}
	if evaluatedAtFrom != nil && evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	} else if evaluatedAtFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", ""))
	} else if evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", "",
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	}

	root := map[string]any{}
	root["size"] = 0

	aggs := map[string]any{
		"integration_type_filter":    map[string]any{"terms": map[string]any{"field": "integrationType", "size": 1000}},
		"resource_type_filter":       map[string]any{"terms": map[string]any{"field": "resourceType", "size": 1000}},
		"integration_id_filter":      map[string]any{"terms": map[string]any{"field": "integrationID", "size": 1000}},
		"resource_collection_filter": map[string]any{"terms": map[string]any{"field": "resourceCollection", "size": 1000}},
		"benchmark_id_filter":        map[string]any{"terms": map[string]any{"field": "benchmarkID", "size": 1000}},
		"control_id_filter":          map[string]any{"terms": map[string]any{"field": "controlID", "size": 1000}},
		"severity_filter":            map[string]any{"terms": map[string]any{"field": "severity", "size": 1000}},
		"compliance_status_filter":   map[string]any{"terms": map[string]any{"field": "complianceStatus", "size": 1000}},
		"state_active_filter":        map[string]any{"terms": map[string]any{"field": "stateActive", "size": 1000}},
	}
	root["aggs"] = aggs

	if len(filters) > 0 {
		root["query"] = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		logger.Error("ComplianceResultsFiltersQuery", zap.Error(err), zap.String("query", string(queryBytes)), zap.String("index", idx))
		return nil, err
	}

	logger.Info("ComplianceResultsFiltersQuery", zap.String("query", string(queryBytes)), zap.String("index", idx))

	var resp ComplianceResultFiltersAggregationResponse
	err = client.Search(ctx, idx, string(queryBytes), &resp)
	if err != nil {
		logger.Error("ComplianceResultsFiltersQuery", zap.Error(err), zap.String("query", string(queryBytes)), zap.String("index", idx))
		return nil, err
	}

	return &resp, nil
}

type ComplianceResultKPIResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
	} `json:"hits"`
	Aggregations struct {
		ResourceCount struct {
			Value int64 `json:"value"`
		} `json:"resource_count"`
		ControlCount struct {
			Value int64 `json:"value"`
		} `json:"control_count"`
		IntegrationCount struct {
			Value int64 `json:"value"`
		} `json:"integration_count"`
	} `json:"aggregations"`
}

func ComplianceResultKPIQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client) (*ComplianceResultKPIResponse, error) {
	root := make(map[string]any)
	root["size"] = 0
	root["track_total_hits"] = true

	filters := make([]map[string]any, 0)
	filters = append(filters, map[string]any{
		"terms": map[string]any{
			"complianceStatus": types.GetFailedComplianceStatuses(),
		},
	})
	root["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}

	root["aggs"] = map[string]any{
		"resource_count": map[string]any{
			"cardinality": map[string]any{
				"field": "platformResourceID",
			},
		},
		"control_count": map[string]any{
			"cardinality": map[string]any{
				"field": "controlID",
			},
		},
		"integration_count": map[string]any{
			"cardinality": map[string]any{
				"field": "integrationID",
			},
		},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	logger.Info("ComplianceResultKPIQuery", zap.String("query", string(queryBytes)))
	var resp ComplianceResultKPIResponse
	err = client.SearchWithTrackTotalHits(ctx, types.ComplianceResultsIndex, string(queryBytes), nil, &resp, true)
	if err != nil {
		logger.Error("ComplianceResultKPIQuery", zap.Error(err), zap.String("query", string(queryBytes)))
		return nil, err
	}
	return &resp, err
}

type ComplianceResultsTopFieldResponse struct {
	Aggregations struct {
		FieldFilter struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"field_filter"`
		BucketCount struct {
			Value int `json:"value"`
		} `json:"bucket_count"`
	} `json:"aggregations"`
}

func ComplianceResultsTopFieldQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	field string, integrationTypes []string, resourceTypeID []string, integrationIDs []string, notIntegrationIDs []string, jobIDs []string,
	benchmarkID []string, controlID []string, severity []types.ComplianceResultSeverity, complianceStatuses []types.ComplianceStatus, stateActives []bool,
	size int, startTime, endTime *time.Time) (*ComplianceResultsTopFieldResponse, error) {
	filters := make([]map[string]any, 0)

	idx := types.ComplianceResultsIndex
	if len(benchmarkID) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"benchmarkID": benchmarkID,
			},
		})
	}

	if len(controlID) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"controlID": controlID,
			},
		})
	}

	if len(complianceStatuses) > 0 {
		cfStrs := make([]string, 0, len(complianceStatuses))
		for _, cf := range complianceStatuses {
			cfStrs = append(cfStrs, string(cf))
		}
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"complianceStatus": cfStrs,
			},
		})
	}

	if len(severity) > 0 {
		sevStrs := make([]string, 0, len(severity))
		for _, s := range severity {
			sevStrs = append(sevStrs, string(s))
		}
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"severity": sevStrs,
			},
		})
	}

	if len(integrationIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"integrationID": integrationIDs,
			},
		})
	}

	if len(jobIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"complianceJobID": jobIDs,
			},
		})
	}

	if len(notIntegrationIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"integrationID": notIntegrationIDs,
			},
		})
	}

	if len(resourceTypeID) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"resourceType": resourceTypeID,
			},
		})
	}

	if len(integrationTypes) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"integrationType": integrationTypes,
			},
		})
	}

	if len(stateActives) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActives {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"stateActive": strStateActive,
			},
		})
	} else {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"stateActive": "true",
			},
		})
	}

	if endTime != nil && startTime != nil {
		filters = append(filters, map[string]any{
			"range": map[string]any{
				"evaluatedAt": map[string]any{
					"gte": startTime.UnixMilli(),
					"lte": endTime.UnixMilli(),
				},
			},
		})
	} else if endTime != nil {
		filters = append(filters, map[string]any{
			"range": map[string]any{
				"evaluatedAt": map[string]any{
					"lte": endTime.UnixMilli(),
				},
			},
		})
	} else if startTime != nil {
		filters = append(filters, map[string]any{
			"range": map[string]any{
				"evaluatedAt": map[string]any{
					"gte": startTime.UnixMilli(),
				},
			},
		})
	}

	root := map[string]any{}
	root["size"] = 0

	root["aggs"] = map[string]any{
		"field_filter": map[string]any{
			"terms": map[string]any{
				"field": field,
				"size":  size,
			},
		},
		"bucket_count": map[string]any{
			"cardinality": map[string]any{
				"field": field,
			},
		},
	}

	if len(filters) > 0 {
		root["query"] = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	logger.Info("ComplianceResultsTopFieldQuery", zap.String("query", string(queryBytes)), zap.String("index", idx))
	var resp ComplianceResultsTopFieldResponse
	err = client.Search(ctx, idx, string(queryBytes), &resp)
	return &resp, err
}

type ResourceTypesComplianceResultsSummaryResponse struct {
	Aggregations struct {
		Summaries struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				Severity struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"severity"`
				ComplianceStatus struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"complianceStatus"`
			} `json:"buckets"`
		} `json:"summaries"`
	} `json:"aggregations"`
}

func ResourceTypesComplianceResultsSummary(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	integrationIDs []string, benchmarkID string) (*ResourceTypesComplianceResultsSummaryResponse, error) {
	var filters []map[string]any

	filters = append(filters, map[string]any{
		"term": map[string]any{
			"parentBenchmarks": benchmarkID,
		},
	})

	if len(integrationIDs) != 0 {
		filters = append(filters, map[string]any{
			"term": map[string]any{
				"integrationID": integrationIDs,
			},
		})
	}

	filters = append(filters, map[string]any{
		"term": map[string]any{
			"stateActive": true,
		},
	})

	request := map[string]any{
		"aggs": map[string]any{
			"summaries": map[string]any{
				"terms": map[string]any{
					"field": "resourceType",
				},
				"aggs": map[string]any{
					"severity": map[string]any{
						"terms": map[string]any{
							"field": "severity",
							"size":  1000,
						},
					},
					"complianceStatus": map[string]any{
						"terms": map[string]any{
							"field": "complianceStatus",
							"size":  1000,
						},
					},
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("ResourceTypesComplianceResultsSummary", zap.String("query", string(queryBytes)))
	var resp ResourceTypesComplianceResultsSummaryResponse
	err = client.Search(ctx, types.ComplianceResultsIndex, string(queryBytes), &resp)
	return &resp, err
}

type ComplianceResultsFieldCountByControlResponse struct {
	Aggregations struct {
		ControlCount struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                string `json:"key"`
				DocCount           int    `json:"doc_count"`
				ComplianceStatuses struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					Buckets                 []struct {
						Key        string `json:"key"`
						DocCount   int    `json:"doc_count"`
						FieldCount struct {
							Value int `json:"value"`
						} `json:"field_count"`
					} `json:"buckets"`
				} `json:"complianceStatus"`
			} `json:"buckets"`
		} `json:"control_count"`
	} `json:"aggregations"`
}

func ComplianceResultsFieldCountByControl(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	field string, integrationTypes []string, resourceTypeID []string, integrationIDs []string, benchmarkID []string, controlID []string,
	severity []types.ComplianceResultSeverity, complianceStatuses []types.ComplianceStatus) (*ComplianceResultsFieldCountByControlResponse, error) {
	terms := make(map[string]any)
	idx := types.ComplianceResultsIndex
	if len(benchmarkID) > 0 {
		terms["benchmarkID"] = benchmarkID
	}

	if len(controlID) > 0 {
		terms["controlID"] = controlID
	}

	if len(severity) > 0 {
		terms["severity"] = severity
	}

	if len(complianceStatuses) > 0 {
		terms["complianceStatus"] = complianceStatuses
	}

	if len(integrationIDs) > 0 {
		terms["integrationID"] = integrationIDs
	}

	if len(resourceTypeID) > 0 {
		terms["resourceType"] = resourceTypeID
	}

	if len(integrationTypes) > 0 {
		terms["integrationType"] = integrationTypes
	}

	terms["stateActive"] = []bool{true}

	root := map[string]any{}
	root["size"] = 0

	root["aggs"] = map[string]any{
		"control_count": map[string]any{
			"terms": map[string]any{
				"field": "controlID",
			},
			"aggs": map[string]any{
				"complianceStatus": map[string]any{
					"terms": map[string]any{
						"field": "complianceStatus",
					},
					"aggs": map[string]any{
						"field_count": map[string]any{
							"cardinality": map[string]any{
								"field": field,
							},
						},
					},
				},
			},
		},
	}

	boolQuery := make(map[string]any)
	if terms != nil && len(terms) > 0 {
		var filters []map[string]any
		for k, vs := range terms {
			filters = append(filters, map[string]any{
				"terms": map[string]any{
					k: vs,
				},
			})
		}

		boolQuery["filter"] = filters
	}
	if len(boolQuery) > 0 {
		root["query"] = map[string]any{
			"bool": boolQuery,
		}
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	logger.Info("ComplianceResultsFieldCountByControl", zap.String("query", string(queryBytes)), zap.String("index", idx))
	var resp ComplianceResultsFieldCountByControlResponse
	err = client.Search(ctx, idx, string(queryBytes), &resp)
	return &resp, err
}

type ComplianceResultsComplianceStatusCountByControlPerIntegrationResponse struct {
	Aggregations struct {
		IntegrationGroup struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key          string `json:"key"`
				ControlCount struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					Buckets                 []struct {
						Key                string `json:"key"`
						DocCount           int    `json:"doc_count"`
						ComplianceStatuses struct {
							Key     string `json:"key"`
							Buckets []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"complianceStatus"`
					} `json:"buckets"`
				} `json:"control_count"`
			} `json:"buckets"`
		} `json:"integration_group"`
	} `json:"aggregations"`
}

func ComplianceResultsComplianceStatusCountByControlPerIntegration(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	integrationTypes []string, resourceTypeID []string, integrationIDs []string, benchmarkID []string, controlID []string,
	severity []types.ComplianceResultSeverity, complianceStatuses []types.ComplianceStatus, startTime, endTime *time.Time) (*ComplianceResultsComplianceStatusCountByControlPerIntegrationResponse, error) {
	terms := make(map[string]any)
	idx := types.ComplianceResultsIndex
	if len(benchmarkID) > 0 {
		terms["benchmarkID"] = benchmarkID
	}

	if len(controlID) > 0 {
		terms["controlID"] = controlID
	}

	if len(severity) > 0 {
		terms["severity"] = severity
	}

	if len(complianceStatuses) > 0 {
		terms["complianceStatus"] = complianceStatuses
	}

	if len(integrationIDs) > 0 {
		terms["integrationID"] = integrationIDs
	}

	if len(resourceTypeID) > 0 {
		terms["resourceType"] = resourceTypeID
	}

	if len(integrationTypes) > 0 {
		terms["integrationType"] = integrationTypes
	}

	terms["stateActive"] = []bool{true}

	root := map[string]any{}
	root["size"] = 0

	root["aggs"] = map[string]any{
		"integration_group": map[string]any{
			"terms": map[string]any{
				"field": "integrationID",
				"size":  10000,
			},
			"aggs": map[string]any{
				"control_count": map[string]any{
					"terms": map[string]any{
						"field": "controlID",
						"size":  10000,
					},
					"aggs": map[string]any{
						"complianceStatus": map[string]any{
							"terms": map[string]any{
								"field": "complianceStatus",
								"size":  10000,
							},
						},
					},
				},
			},
		},
	}

	boolQuery := make(map[string]any)
	if (terms != nil && len(terms) > 0) || (endTime != nil || startTime != nil) {
		var filters []map[string]any

		if terms != nil && len(terms) > 0 {
			for k, vs := range terms {
				filters = append(filters, map[string]any{
					"terms": map[string]any{
						k: vs,
					},
				})
			}
		}
		if endTime != nil && startTime != nil {
			filters = append(filters, map[string]any{
				"range": map[string]any{
					"evaluatedAt": map[string]any{
						"gte": startTime.UnixMilli(),
						"lte": endTime.UnixMilli(),
					},
				},
			})
		} else if endTime != nil {
			filters = append(filters, map[string]any{
				"range": map[string]any{
					"evaluatedAt": map[string]any{
						"lte": endTime.UnixMilli(),
					},
				},
			})
		} else if startTime != nil {
			filters = append(filters, map[string]any{
				"range": map[string]any{
					"evaluatedAt": map[string]any{
						"gte": startTime.UnixMilli(),
					},
				},
			})
		}

		boolQuery["filter"] = filters
	}
	if len(boolQuery) > 0 {
		root["query"] = map[string]any{
			"bool": boolQuery,
		}
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	logger.Info("ComplianceResultsFieldCountByControl", zap.String("query", string(queryBytes)), zap.String("index", idx))
	var resp ComplianceResultsComplianceStatusCountByControlPerIntegrationResponse
	err = client.Search(ctx, idx, string(queryBytes), &resp)
	if err != nil {
		logger.Error("ComplianceResultsFieldCountByControl", zap.Error(err), zap.String("query", string(queryBytes)), zap.String("index", idx))
		return nil, err
	}
	return &resp, nil
}

type ComplianceResultCountPerPlatformResourceIDsResponse struct {
	Aggregations struct {
		PlatformResourceIDGroup struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"og_resource_id_group"`
	} `json:"aggregations"`
}

func FetchComplianceResultCountPerPlatformResourceIDs(ctx context.Context, logger *zap.Logger, client opengovernance.Client, platformResourceIDs []string,
	severities []types.ComplianceResultSeverity, complianceStatuses []types.ComplianceStatus,
) (map[string]int, error) {
	var filters []map[string]any

	if len(platformResourceIDs) == 0 {
		return map[string]int{}, nil
	}

	filters = append(filters, map[string]any{
		"terms": map[string]any{
			"platformResourceID": platformResourceIDs,
		},
	})
	if len(severities) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"severity": severities,
			},
		})
	}
	if len(complianceStatuses) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"complianceStatus": complianceStatuses,
			},
		})
	}
	filters = append(filters, map[string]any{
		"term": map[string]any{
			"stateActive": true,
		},
	})

	request := map[string]any{
		"aggs": map[string]any{
			"og_resource_id_group": map[string]any{
				"terms": map[string]any{
					"field": "platformResourceID",
					"size":  len(platformResourceIDs),
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchComplianceResultCountPerPlatformResourceIDs", zap.String("query", string(queryBytes)))
	var resp ComplianceResultCountPerPlatformResourceIDsResponse
	err = client.Search(ctx, types.ComplianceResultsIndex, string(queryBytes), &resp)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, bucket := range resp.Aggregations.PlatformResourceIDGroup.Buckets {
		result[bucket.Key] = bucket.DocCount
	}

	return result, nil
}

type ComplianceResultsPerControlForResourceIdResponse struct {
	Aggregations struct {
		ControlIDGroup struct {
			Buckets []struct {
				Key       string `json:"key"`
				HitSelect struct {
					Hits struct {
						Hits []struct {
							Source types.ComplianceResult `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"hit_select"`
			} `json:"buckets"`
		} `json:"control_id_group"`
	} `json:"aggregations"`
}

func FetchComplianceResultsPerControlForResourceId(ctx context.Context, logger *zap.Logger, client opengovernance.Client, platformResourceID string) ([]types.ComplianceResult, error) {
	request := map[string]any{
		"aggs": map[string]any{
			"control_id_group": map[string]any{
				"terms": map[string]any{
					"field": "controlID",
					"size":  10000,
				},
				"aggs": map[string]any{
					"hit_select": map[string]any{
						"top_hits": map[string]any{
							"sort": map[string]any{
								"complianceJobID": "desc",
							},
							"size": 1,
						},
					},
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": map[string]any{
					"term": map[string]any{
						"platformResourceID": platformResourceID,
					},
				},
			},
		},
		"size": 0,
	}

	queryBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchComplianceResultsPerControlForResourceId", zap.String("query", string(queryBytes)))
	var resp ComplianceResultsPerControlForResourceIdResponse
	err = client.Search(ctx, types.ComplianceResultsIndex, string(queryBytes), &resp)
	if err != nil {
		return nil, err
	}

	var complianceResults []types.ComplianceResult
	for _, bucket := range resp.Aggregations.ControlIDGroup.Buckets {
		for _, hit := range bucket.HitSelect.Hits.Hits {
			complianceResults = append(complianceResults, hit.Source)
		}
	}

	return complianceResults, nil
}

func FetchComplianceResultByID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, complianceResultId string) (*types.ComplianceResult, error) {
	query := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"es_id": complianceResultId,
			},
		},
		"size": 1,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchComplianceResultByID", zap.String("query", string(queryBytes)))
	var resp ComplianceResultsQueryResponse
	err = client.Search(ctx, types.ComplianceResultsIndex, string(queryBytes), &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Hits.Hits) == 0 {
		return nil, nil
	}

	return &resp.Hits.Hits[0].Source, nil
}

func ComplianceResultsQueryV2(ctx context.Context, logger *zap.Logger, client opengovernance.Client, resourceIDs []string, notResourceIDs []string,
	integrationTypes []string, integrationID []string, notIntegrationID []string, resourceTypes []string, notResourceTypes []string,
	benchmarkID []string, notBenchmarkID []string, controlID []string, notControlID []string, severity []types.ComplianceResultSeverity,
	notSeverity []types.ComplianceResultSeverity, lastTransitionFrom *time.Time, lastTransitionTo *time.Time, notLastTransitionFrom *time.Time,
	notLastTransitionTo *time.Time, evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time, stateActive []bool,
	complianceStatuses []types.ComplianceStatus, sorts []api.ComplianceResultsSortV2, pageSizeLimit int, searchAfter []any) ([]ComplianceResultsQueryHit, int64, error) {
	idx := types.ComplianceResultsIndex

	requestSort := make([]map[string]any, 0, len(sorts)+1)
	for _, sort := range sorts {
		switch {
		case sort.ResourceType != nil:
			requestSort = append(requestSort, map[string]any{
				"resourceType": *sort.ResourceType,
			})
		case sort.BenchmarkID != nil:
			requestSort = append(requestSort, map[string]any{
				"benchmarkID": *sort.BenchmarkID,
			})
		case sort.ControlID != nil:
			requestSort = append(requestSort, map[string]any{
				"controlID": *sort.ControlID,
			})
		case sort.Severity != nil:
			scriptSource :=
				`if (params['_source']['severity'] == 'critical') {
					return 5
				} else if (params['_source']['severity'] == 'high') {
					return 4
				} else if (params['_source']['severity'] == 'medium') {
					return 3
				} else if (params['_source']['severity'] == 'low') {
					return 2
				} else if (params['_source']['severity'] == 'none') {
					return 1
				} else {
					return 1
				}`
			requestSort = append(requestSort, map[string]any{
				"_script": map[string]any{
					"type": "number",
					"script": map[string]any{
						"lang":   "painless",
						"source": scriptSource,
					},
					"order": *sort.Severity,
				},
			})
		case sort.ComplianceStatus != nil:
			scriptSource :=
				`if (params['_source']['complianceStatus'] == 'alarm') {
					return 5
				} else if (params['_source']['complianceStatus'] == 'error') {
					return 4
				} else if (params['_source']['complianceStatus'] == 'info') {
					return 3
				} else if (params['_source']['complianceStatus'] == 'skip') {
					return 2
				} else if (params['_source']['complianceStatus'] == 'ok') {
					return 1
				} else {
					return 1
				}`
			requestSort = append(requestSort, map[string]any{
				"_script": map[string]any{
					"type": "number",
					"script": map[string]any{
						"lang":   "painless",
						"source": scriptSource,
					},
					"order": *sort.ComplianceStatus,
				},
			})
		case sort.LastUpdated != nil:
			requestSort = append(requestSort, map[string]any{
				"lastTransition": *sort.LastUpdated,
			})
		}
	}
	requestSort = append(requestSort, map[string]any{
		"_id": "asc",
	})

	var filters []opengovernance.BoolFilter
	if len(resourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceID", resourceIDs))
	}
	if len(notResourceIDs) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("resourceID", notResourceIDs)))
	}
	if len(resourceTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("resourceType", resourceTypes))
	}
	if len(notResourceTypes) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("resourceType", notResourceTypes)))
	}
	if len(benchmarkID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("parentBenchmarks", benchmarkID))
	}
	if len(notBenchmarkID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("parentBenchmarks", notBenchmarkID)))
	}
	if len(controlID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("controlID", controlID))
	}
	if len(notControlID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("controlID", notControlID)))
	}
	if len(severity) > 0 {
		strSeverity := make([]string, 0)
		for _, s := range severity {
			strSeverity = append(strSeverity, string(s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("severity", strSeverity))
	}
	if len(notSeverity) > 0 {
		strSeverity := make([]string, 0)
		for _, s := range notSeverity {
			strSeverity = append(strSeverity, string(s))
		}
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("severity", strSeverity)))
	}
	if len(complianceStatuses) > 0 {
		strComplianceStatus := make([]string, 0)
		for _, cr := range complianceStatuses {
			strComplianceStatus = append(strComplianceStatus, string(cr))
		}
		filters = append(filters, opengovernance.NewTermsFilter("complianceStatus", strComplianceStatus))
	}
	if len(integrationID) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationID", integrationID))
	}
	if len(notIntegrationID) > 0 {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", notIntegrationID)))
	}
	if len(integrationTypes) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
	}
	if lastTransitionFrom != nil && lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	} else if lastTransitionFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", lastTransitionFrom.UnixMilli()),
			"", ""))
	} else if lastTransitionTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("lastTransition",
			"", "",
			"", fmt.Sprintf("%d", lastTransitionTo.UnixMilli())))
	}
	if notLastTransitionFrom != nil && notLastTransitionTo != nil {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", notLastTransitionFrom.UnixMilli()),
			"", fmt.Sprintf("%d", notLastTransitionTo.UnixMilli()))))
	} else if notLastTransitionFrom != nil {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewRangeFilter("lastTransition",
			"", fmt.Sprintf("%d", notLastTransitionFrom.UnixMilli()),
			"", "")))
	} else if notLastTransitionTo != nil {
		filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewRangeFilter("lastTransition",
			"", "",
			"", fmt.Sprintf("%d", notLastTransitionTo.UnixMilli()))))
	}
	if evaluatedAtFrom != nil && evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	} else if evaluatedAtFrom != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", fmt.Sprintf("%d", evaluatedAtFrom.UnixMilli()),
			"", ""))
	} else if evaluatedAtTo != nil {
		filters = append(filters, opengovernance.NewRangeFilter("evaluatedAt",
			"", "",
			"", fmt.Sprintf("%d", evaluatedAtTo.UnixMilli())))
	}

	query := make(map[string]any)
	if len(filters) > 0 {
		query["query"] = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}
	query["sort"] = requestSort
	if len(searchAfter) > 0 {
		query["search_after"] = searchAfter
	}
	if pageSizeLimit == 0 {
		pageSizeLimit = 1000
	}
	query["size"] = pageSizeLimit
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	logger.Info("ComplianceResultsQuery", zap.String("query", string(queryJson)), zap.String("index", idx))

	var response ComplianceResultsQueryResponse
	err = client.SearchWithTrackTotalHits(ctx, idx, string(queryJson), nil, &response, true)
	if err != nil {
		return nil, 0, err
	}

	return response.Hits.Hits, response.Hits.Total.Value, err
}
