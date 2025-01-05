package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/es"
	"strings"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

type ResourceIdentifierFetchResponse struct {
	Hits ResourceIdentifierFetchHits `json:"hits"`
}
type ResourceIdentifierFetchHits struct {
	Total opengovernance.SearchTotal   `json:"total"`
	Hits  []ResourceIdentifierFetchHit `json:"hits"`
}
type ResourceIdentifierFetchHit struct {
	ID      string            `json:"_id"`
	Score   float64           `json:"_score"`
	Index   string            `json:"_index"`
	Type    string            `json:"_type"`
	Version int64             `json:"_version,omitempty"`
	Source  es.LookupResource `json:"_source"`
	Sort    []any             `json:"sort"`
}

func GetResourceIDsForAccountResourceTypeFromES(ctx context.Context, client opengovernance.Client, integrationID, resourceType string, additionalFilters []map[string]any, searchAfter []any, size int) (*ResourceIdentifierFetchResponse, error) {
	root := map[string]any{}
	root["query"] = map[string]any{
		"bool": map[string]any{
			"filter": append([]map[string]any{
				{"term": map[string]string{"integration_id": integrationID}},
				{"term": map[string]string{"resource_type": strings.ToLower(resourceType)}},
			}, additionalFilters...),
		},
	}
	if searchAfter != nil {
		root["search_after"] = searchAfter
	}
	root["size"] = size
	root["sort"] = []map[string]any{
		{"described_at": "asc"},
		{"_id": "desc"},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	var response ResourceIdentifierFetchResponse
	err = client.Search(ctx, es.InventorySummaryIndex,
		string(queryBytes), &response)
	if err != nil {
		fmt.Println("query=", string(queryBytes))
		return nil, err
	}

	return &response, nil
}

func GetResourceIDsNotInIntegrationsFromES(ctx context.Context, client opengovernance.Client, integrationIDs []string, searchAfter []any, size int) (*ResourceIdentifierFetchResponse, error) {
	root := map[string]any{}
	if len(integrationIDs) > 0 {
		root["query"] = map[string]any{
			"bool": map[string]any{
				"must_not": []map[string]any{
					{"terms": map[string]any{"integration_id": integrationIDs}},
				},
			},
		}
	} else {
		root["query"] = map[string]any{
			"match_all": map[string]any{},
		}
	}
	if searchAfter != nil {
		root["search_after"] = searchAfter
	}
	root["size"] = size
	root["sort"] = []map[string]any{
		{"described_at": "asc"},
		{"_id": "desc"},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	var response ResourceIdentifierFetchResponse
	err = client.Search(ctx, es.InventorySummaryIndex,
		string(queryBytes), &response)
	if err != nil {
		fmt.Println("query=", string(queryBytes))
		return nil, err
	}

	return &response, nil
}

func GetResourceIDsForIntegrationFromES(ctx context.Context, client opengovernance.Client, integrationID string, searchAfter []any, size int) (*ResourceIdentifierFetchResponse, error) {
	root := map[string]any{}
	root["query"] = map[string]any{
		"bool": map[string]any{
			"filter": []map[string]any{
				{"term": map[string]string{"integration_id": integrationID}},
			},
		},
	}
	if searchAfter != nil {
		root["search_after"] = searchAfter
	}
	root["size"] = size
	root["sort"] = []map[string]any{
		{"described_at": "asc"},
		{"_id": "desc"},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	var response ResourceIdentifierFetchResponse
	err = client.Search(ctx, es.InventorySummaryIndex,
		string(queryBytes), &response)
	if err != nil {
		fmt.Println("query=", string(queryBytes))
		return nil, err
	}

	return &response, nil
}

type InventoryCountResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{}   `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
}

func GetInventoryCountResponse(ctx context.Context, client opengovernance.Client, resourceType string) (int64, error) {
	query := fmt.Sprintf(`{"size": 0, "query": {"bool": {"filter": [{"term": {"resource_type": "%s"}}]}}}`, resourceType)

	fmt.Println("GetInventoryCountResponse, query=", query)
	var response InventoryCountResponse
	err := client.SearchWithTrackTotalHits(ctx, es.InventorySummaryIndex,
		query, nil, &response, true)
	if err != nil {
		return 0, err
	}
	fmt.Println("GetInventoryCountResponse, response=", response)

	return response.Hits.Total.Value, nil
}
