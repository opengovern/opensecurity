package es

import (
	"encoding/json"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/inventory/api"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type GetAsyncQueryRunResultSource struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	RunId       string               `json:"runID"`
	QueryID     string               `json:"queryID"`
	Parameters  []api.QueryParameter `json:"parameters"`
	ColumnNames []string             `json:"columnNames"`
	CreatedBy   string               `json:"createdBy"`
	TriggeredAt int64                `json:"triggeredAt"`
	EvaluatedAt int64                `json:"evaluatedAt"`
	Result      [][]string           `json:"result"`
}

type GetAsyncQueryRunResultHit struct {
	Index  string                       `json:"_index"`
	ID     string                       `json:"_id"`
	Score  float64                      `json:"_score"`
	Source GetAsyncQueryRunResultSource `json:"_source"`
}

type GetAsyncQueryRunResultResponse struct {
	Hits struct {
		Hits []GetAsyncQueryRunResultHit `json:"hits"`
	} `json:"hits"`
}

func GetAsyncQueryRunResult(ctx context.Context, logger *zap.Logger, client opengovernance.Client, runID string) (*GetAsyncQueryRunResultSource, error) {
	idx := types.QueryRunIndex

	request := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"runID": runID,
			},
		},
	}

	jsonReq, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetching async query run result", zap.String("request", string(jsonReq)), zap.String("index", idx))

	var resp GetAsyncQueryRunResultResponse
	err = client.Search(ctx, idx, string(jsonReq), &resp)
	if err != nil {
		logger.Error("Failed to fetch async query run result", zap.Error(err), zap.String("request", string(jsonReq)), zap.String("index", idx))
		return nil, err
	}

	var result GetAsyncQueryRunResultSource
	for _, hit := range resp.Hits.Hits {
		result = hit.Source
	}
	return &result, nil
}
