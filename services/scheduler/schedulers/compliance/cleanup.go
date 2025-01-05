package compliance

import (
	"encoding/json"
	"fmt"
	es2 "github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opencomply/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func (s *JobScheduler) cleanupComplianceResultsNotInIntegrations(ctx context.Context, integrationIDs []string) {
	var searchAfter []any
	totalDeletedCount := 0
	deletedIntegrationIDs := make(map[string]bool)
	for {
		esResp, err := getComplianceResultsNotInIntegrationsFromES(ctx, s.esClient, integrationIDs, searchAfter, 1000)
		if err != nil {
			s.logger.Error("failed to get resource ids from es", zap.Error(err))
			break
		}
		totalDeletedCount += len(esResp.Hits.Hits)
		if len(esResp.Hits.Hits) == 0 {
			break
		}
		deletedCount := 0
		for _, hit := range esResp.Hits.Hits {
			deletedIntegrationIDs[hit.Source.IntegrationID] = true
			searchAfter = hit.Sort

			complianceResult := types.ComplianceResult{
				PlatformResourceID: hit.Source.PlatformResourceID,
				ResourceID:         hit.Source.ResourceID,
				IntegrationID:      hit.Source.IntegrationID,
				ControlID:          hit.Source.ControlID,
				BenchmarkID:        hit.Source.BenchmarkID,
			}
			keys, idx := complianceResult.KeysAndIndex()
			deletedCount += 1
			key := es2.HashOf(keys...)
			complianceResult.EsID = key
			complianceResult.EsIndex = idx
			err = s.esClient.Delete(key, idx)
			if err != nil {
				s.logger.Error("failed to delete complianceResult from open-search", zap.Error(err))
				return
			}
		}
	}
	s.logger.Info("total deleted resource count", zap.Int("count", totalDeletedCount),
		zap.Any("deleted integrations", deletedIntegrationIDs))
	return
}

func getComplianceResultsNotInIntegrationsFromES(ctx context.Context, client opengovernance.Client, integrationIDs []string, searchAfter []any, size int) (*ComplianceResultFetchResponse, error) {
	root := map[string]any{}
	root["query"] = map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{"terms": map[string]any{"integrationID": integrationIDs}},
			},
		},
	}
	if searchAfter != nil {
		root["search_after"] = searchAfter
	}
	root["size"] = size
	root["sort"] = []map[string]any{
		{"evaluatedAt": "asc"},
		{"_id": "desc"},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	var response ComplianceResultFetchResponse
	err = client.Search(ctx, types.ComplianceResultsIndex,
		string(queryBytes), &response)
	if err != nil {
		fmt.Println("query=", string(queryBytes))
		return nil, err
	}

	return &response, nil
}

type ComplianceResultFetchResponse struct {
	Hits ComplianceResultFetchHits `json:"hits"`
}
type ComplianceResultFetchHits struct {
	Total opengovernance.SearchTotal `json:"total"`
	Hits  []ComplianceResultFetchHit `json:"hits"`
}
type ComplianceResultFetchHit struct {
	ID      string                 `json:"_id"`
	Score   float64                `json:"_score"`
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Version int64                  `json:"_version,omitempty"`
	Source  types.ComplianceResult `json:"_source"`
	Sort    []any                  `json:"sort"`
}
