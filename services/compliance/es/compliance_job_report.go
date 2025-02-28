package es

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type ControlViewResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal `json:"total"`
		Hits  []struct {
			ID      string                               `json:"_id"`
			Score   float64                              `json:"_score"`
			Index   string                               `json:"_index"`
			Type    string                               `json:"_type"`
			Version int64                                `json:"_version,omitempty"`
			Source  types.ComplianceJobReportControlView `json:"_source"`
			Sort    []any                                `json:"sort"`
		}
	}
}

func GetJobReportControlViewByJobID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, jobID string, auditable bool,
	controls []string) (*types.ComplianceJobReportControlView, error) {
	request := make(map[string]any)
	request["size"] = 1
	request["query"] = map[string]any{
		"bool": map[string]any{
			"must": []map[string]any{
				{
					"term": map[string]any{
						"job_summary.job_id": jobID,
					},
				},
				{
					"term": map[string]any{
						"job_summary.auditable": auditable,
					},
				},
			},
		},
	}
	if len(controls) > 0 {
		sourceItems := []string{
			"es_id",
			"es_index",
			"compliance_summary",
			"job_summary",
		}
		for _, control := range controls {
			sourceItems = append(sourceItems, fmt.Sprintf("controls.%s", control))
		}
		request["_source"] = sourceItems
	}
	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("ES Policy", zap.String("index", types.ComplianceJobReportControlViewIndex), zap.String("query", string(b)))

	var response ControlViewResponse
	err = client.Search(ctx, types.ComplianceJobReportControlViewIndex, string(b), &response)
	if err != nil {
		return nil, err
	}

	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	return &response.Hits.Hits[0].Source, nil
}

type ControlSummaryResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal `json:"total"`
		Hits  []struct {
			ID      string                                  `json:"_id"`
			Score   float64                                 `json:"_score"`
			Index   string                                  `json:"_index"`
			Type    string                                  `json:"_type"`
			Version int64                                   `json:"_version,omitempty"`
			Source  types.ComplianceJobReportControlSummary `json:"_source"`
			Sort    []any                                   `json:"sort"`
		}
	}
}

func GetJobReportControlSummaryByJobID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, jobID string,
	controls []string) (*types.ComplianceJobReportControlSummary, error) {
	request := make(map[string]any)
	request["size"] = 1
	request["query"] = map[string]any{
		"bool": map[string]any{
			"must": []map[string]any{
				{
					"term": map[string]any{
						"job_summary.job_id": jobID,
					},
				},
			},
		},
	}
	if len(controls) > 0 {
		sourceItems := []string{
			"es_id",
			"es_index",
			"control_score",
			"compliance_summary",
			"job_summary",
		}
		for _, control := range controls {
			sourceItems = append(sourceItems, fmt.Sprintf("controls.%s", control))
		}
		request["_source"] = sourceItems
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("ES Policy to get compliance job report control summary", zap.String("index", types.ComplianceJobReportControlSummaryIndex), zap.String("query", string(b)))

	var response ControlSummaryResponse
	err = client.Search(ctx, types.ComplianceJobReportControlSummaryIndex, string(b), &response)
	if err != nil {
		return nil, err
	}

	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	return &response.Hits.Hits[0].Source, nil
}

type ResourceViewResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal `json:"total"`
		Hits  []struct {
			ID      string                                `json:"_id"`
			Score   float64                               `json:"_score"`
			Index   string                                `json:"_index"`
			Type    string                                `json:"_type"`
			Version int64                                 `json:"_version,omitempty"`
			Source  types.ComplianceJobReportResourceView `json:"_source"`
			Sort    []any                                 `json:"sort"`
		}
	}
}

func GetJobReportResourceViewByJobID(ctx context.Context, logger *zap.Logger, client opengovernance.Client, jobID string, auditable bool) (*types.ComplianceJobReportResourceView, error) {
	request := make(map[string]any)
	request["size"] = 1
	request["query"] = map[string]any{
		"bool": map[string]any{
			"must": []map[string]any{
				{
					"term": map[string]any{
						"job_summary.job_id": jobID,
					},
				},
				{
					"term": map[string]any{
						"job_summary.auditable": auditable,
					},
				},
			},
		},
	}
	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	logger.Info("ES Policy", zap.String("index", types.ComplianceJobReportResourceViewIndex), zap.String("query", string(b)))

	var response ResourceViewResponse
	err = client.Search(ctx, types.ComplianceJobReportResourceViewIndex, string(b), &response)
	if err != nil {
		return nil, err
	}

	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	return &response.Hits.Hits[0].Source, nil
}
