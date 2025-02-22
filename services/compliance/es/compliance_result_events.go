package es

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/types"
	"github.com/opengovern/opensecurity/services/compliance/api"
	"go.uber.org/zap"
)

type ComplianceResultDriftEventsQueryHit struct {
	ID      string                           `json:"_id"`
	Score   float64                          `json:"_score"`
	Index   string                           `json:"_index"`
	Type    string                           `json:"_type"`
	Version int64                            `json:"_version,omitempty"`
	Source  types.ComplianceResultDriftEvent `json:"_source"`
	Sort    []any                            `json:"sort"`
}

type ComplianceResultDriftEventsQueryResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal            `json:"total"`
		Hits  []ComplianceResultDriftEventsQueryHit `json:"hits"`
	} `json:"hits"`
	PitID string `json:"pit_id"`
}

type FetchComplianceResultDriftEventsByComplianceResultIDResponse struct {
	Hits struct {
		Hits []ComplianceResultDriftEventsQueryHit `json:"hits"`
	} `json:"hits"`
}

func FetchComplianceResultDriftEventsByComplianceResultIDs(ctx context.Context, logger *zap.Logger, client opengovernance.Client, complianceResultID []string) ([]types.ComplianceResultDriftEvent, error) {
	request := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []any{
					map[string]any{
						"terms": map[string][]string{
							"complianceResultEsID": complianceResultID,
						},
					},
				},
			},
		},
		"sort": map[string]any{
			"evaluatedAt": "desc",
		},
		"size": 10000,
	}

	jsonReq, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetching complianceResult events", zap.String("request", string(jsonReq)), zap.String("index", types.ComplianceResultEventsIndex))

	var resp FetchComplianceResultDriftEventsByComplianceResultIDResponse
	err = client.Search(ctx, types.ComplianceResultEventsIndex, string(jsonReq), &resp)
	if err != nil {
		logger.Error("Failed to fetch complianceResult events", zap.Error(err), zap.String("request", string(jsonReq)), zap.String("index", types.ComplianceResultEventsIndex))
		return nil, err
	}
	result := make([]types.ComplianceResultDriftEvent, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		result = append(result, hit.Source)
	}
	return result, nil
}

func ComplianceResultDriftEventsQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	complianceResultIDs []string, platformResourceIDs []string,
	integrationType []integration.Type, integrationID []string, notIntegrationID []string,
	resourceTypes []string,
	benchmarkID []string, controlID []string, severity []types.ComplianceResultSeverity,
	evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time,
	stateActive []bool, complianceStatuses []types.ComplianceStatus,
	sorts []api.ComplianceResultDriftEventsSort, pageSizeLimit int, searchAfter []any) ([]ComplianceResultDriftEventsQueryHit, int64, error) {
	idx := types.ComplianceResultEventsIndex

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

		case sort.ResourceType != nil:
			requestSort = append(requestSort, map[string]any{
				"resourceType": *sort.ResourceType,
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
	if len(complianceResultIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("complianceResultEsID", complianceResultIDs))
	}
	if len(platformResourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("platformResourceID", platformResourceIDs))
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
	if len(integrationType) > 0 {
		var integrationTypes []string
		for _, p := range integrationType {
			integrationTypes = append(integrationTypes, p.String())
		}
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
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

	logger.Info("ComplianceResultDriftEventsQuery", zap.String("query", string(queryJson)), zap.String("index", idx))

	var response ComplianceResultDriftEventsQueryResponse
	err = client.SearchWithTrackTotalHits(ctx, idx, string(queryJson), nil, &response, true)
	if err != nil {
		return nil, 0, err
	}

	return response.Hits.Hits, response.Hits.Total.Value, err
}

type ComplianceResultDriftEventFiltersAggregationResponse struct {
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

func ComplianceResultDriftEventsFiltersQuery(ctx context.Context, logger *zap.Logger, client opengovernance.Client,
	complianceResultIDs []string, platformResourceIDs []string, integrationType []integration.Type, integrationID []string, notIntegrationID []string,
	resourceTypes []string, benchmarkID []string, controlID []string, severity []types.ComplianceResultSeverity,
	evaluatedAtFrom *time.Time, evaluatedAtTo *time.Time,
	stateActive []bool, complianceStatuses []types.ComplianceStatus,
) (*ComplianceResultDriftEventFiltersAggregationResponse, error) {
	idx := types.ComplianceResultEventsIndex

	var filters []opengovernance.BoolFilter
	if len(complianceResultIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("complianceResultEsID", complianceResultIDs))
	}
	if len(platformResourceIDs) > 0 {
		filters = append(filters, opengovernance.NewTermsFilter("platformResourceID", platformResourceIDs))
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
	if len(integrationType) > 0 {
		var integrationTypes []string
		for _, p := range integrationType {
			integrationTypes = append(integrationTypes, p.String())
		}
		filters = append(filters, opengovernance.NewTermsFilter("integrationType", integrationTypes))
	}
	if len(stateActive) > 0 {
		strStateActive := make([]string, 0)
		for _, s := range stateActive {
			strStateActive = append(strStateActive, fmt.Sprintf("%v", s))
		}
		filters = append(filters, opengovernance.NewTermsFilter("stateActive", strStateActive))
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
		logger.Error("ComplianceResultDriftEventsFiltersQuery", zap.Error(err), zap.String("query", string(queryBytes)), zap.String("index", idx))
		return nil, err
	}

	logger.Info("ComplianceResultDriftEventsFiltersQuery", zap.String("query", string(queryBytes)), zap.String("index", idx))

	var resp ComplianceResultDriftEventFiltersAggregationResponse
	err = client.Search(ctx, idx, string(queryBytes), &resp)
	if err != nil {
		logger.Error("ComplianceResultDriftEventsFiltersQuery", zap.Error(err), zap.String("query", string(queryBytes)), zap.String("index", idx))
		return nil, err
	}

	return &resp, nil
}

func FetchComplianceResultDriftEventByID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, driftEventId string) (*types.ComplianceResultDriftEvent, error) {
	query := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"es_id": driftEventId,
			},
		},
		"size": 1,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	logger.Info("FetchComplianceResultDriftEventByID", zap.String("query", string(queryBytes)))
	var resp ComplianceResultDriftEventsQueryResponse
	err = client.Search(ctx, types.ComplianceResultEventsIndex, string(queryBytes), &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Hits.Hits) == 0 {
		return nil, nil
	}

	return &resp.Hits.Hits[0].Source, nil
}

type ComplianceResultDriftEventsCountResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal `json:"total"`
	} `json:"hits"`
	PitID string `json:"pit_id"`
}

func ComplianceResultDriftEventsCount(ctx context.Context, client opengovernance.Client, benchmarkIDs []string, complianceStatuses []types.ComplianceStatus, stateActives []bool, startTime, endTime *time.Time) (int64, error) {
	idx := types.ComplianceResultEventsIndex

	filters := make([]map[string]any, 0)
	if len(complianceStatuses) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"complianceStatus": complianceStatuses,
			},
		})
	}
	if len(stateActives) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"stateActive": stateActives,
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
	if len(benchmarkIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"benchmarkID": benchmarkIDs,
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
