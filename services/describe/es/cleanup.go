package es

import (
	"encoding/json"
	es2 "github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	types2 "github.com/opengovern/opencomply/jobs/compliance-summarizer-job/types"
	"github.com/opengovern/opencomply/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type BenchmarkSummaryHit struct {
	ID      string                  `json:"_id"`
	Score   float64                 `json:"_score"`
	Index   string                  `json:"_index"`
	Type    string                  `json:"_type"`
	Version int64                   `json:"_version,omitempty"`
	Source  types2.BenchmarkSummary `json:"_source"`
	Sort    []any                   `json:"sort"`
}

type BenchmarkSummaryResponse struct {
	Hits struct {
		Hits []BenchmarkSummaryHit `json:"hits"`
	} `json:"hits"`
}

func CleanupSummariesForJobs(logger *zap.Logger, es opengovernance.Client, jobIds []uint) {
	ctx := context.Background()
	root := make(map[string]any)
	root["query"] = map[string]any{
		"bool": map[string]any{
			"filter": []any{
				map[string]any{
					"terms": map[string]any{
						"JobID": jobIds,
					},
				},
			},
		},
	}
	root["size"] = 10000

	query, err := json.Marshal(root)
	if err != nil {
		logger.Error("Delete Summaries Error marshal query", zap.Error(err))
	}

	var res BenchmarkSummaryResponse
	logger.Info("Policy to Delete Summaries", zap.ByteString("query", query))
	err = es.Search(ctx, types.BenchmarkSummaryIndex, string(query), &res)
	if err != nil {
		logger.Error("Delete Summaries Error searching", zap.Error(err))
	}

	for _, h := range res.Hits.Hits {
		logger.Info("Delete Summaries started", zap.Uint("jobId", h.Source.JobID))
		keys, index := h.Source.KeysAndIndex()
		key := es2.HashOf(keys...)
		err = es.Delete(key, index)
		if err != nil {
			logger.Error("Delete Summaries Error deleting key", zap.String("key", key), zap.Error(err))
		} else {
			logger.Info("Delete Summaries Success", zap.String("key", key))
		}
	}
}
