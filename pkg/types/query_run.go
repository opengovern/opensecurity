package types

import (
	"github.com/opengovern/opencomply/services/core/api"
)

type QueryRunResult struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	RunId       string               `json:"runID"`
	CreatedBy   string               `json:"createdBy"`
	TriggeredAt int64                `json:"triggeredAt"`
	EvaluatedAt int64                `json:"evaluatedAt"`
	QueryID     string               `json:"queryID"`
	Parameters  []api.QueryParameter `json:"parameters"`
	ColumnNames []string             `json:"columnNames"`
	Result      [][]string           `json:"result"`
}

func (r QueryRunResult) KeysAndIndex() ([]string, string) {
	return []string{
		r.RunId,
	}, QueryRunIndex
}
