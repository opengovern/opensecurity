package es

import (
	"context"
	"encoding/json"

	"github.com/opengovern/og-util/pkg/integration"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opencomply/pkg/utils"
	"github.com/opengovern/opencomply/services/describe"
)

type ResourceTypeCountsResponse struct {
	Aggregations struct {
		ResourceTypeGroup struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"resource_type_group"`
	} `json:"aggregations"`
}

func GetResourceTypeCounts(ctx context.Context, client opengovernance.Client, integrationTypes []integration.Type, integrationIDs []string, resourceTypes []string, size int) (map[string]int, error) {
	var filters []any
	if len(integrationTypes) > 0 {
		integrationTypesStr := make([]string, 0, len(integrationTypes))
		for _, integrationType := range integrationTypes {
			integrationTypesStr = append(integrationTypesStr, integrationType.String())
		}
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"integration_type": integrationTypesStr,
			},
		})
	}
	if len(integrationIDs) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"integration_id": integrationIDs,
			},
		})
	}
	if len(resourceTypes) > 0 {
		resourceTypes = utils.ToLowerStringSlice(resourceTypes)
		filters = append(filters, map[string]any{
			"terms": map[string][]string{
				"resource_type": resourceTypes,
			},
		})
	}

	query := map[string]any{
		"size": 0,
		"aggs": map[string]any{
			"resource_type_group": map[string]any{
				"terms": map[string]any{
					"field": "resource_type",
					"size":  size,
				},
			},
		},
		"query": map[string]any{
			"bool": map[string][]any{
				"filter": filters,
			},
		},
	}

	queryStr, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	var response ResourceTypeCountsResponse
	err = client.Search(ctx, describe.InventorySummaryIndex,
		string(queryStr), &response)
	if err != nil {
		return nil, err
	}

	res := make(map[string]int)
	for _, bucket := range response.Aggregations.ResourceTypeGroup.Buckets {
		res[bucket.Key] = bucket.DocCount
	}
	return res, nil
}
