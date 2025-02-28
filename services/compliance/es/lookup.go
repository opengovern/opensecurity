package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/es"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

const (
	InventorySummaryIndex = "inventory_summary"
)

type LookupQueryResponse struct {
	Aggregations struct {
		Resources struct {
			Buckets []struct {
				Key       string `json:"key"`
				HitSelect struct {
					Hits struct {
						Hits []struct {
							Source es.LookupResource `json:"_source"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"hit_select"`
			} `json:"buckets"`
		} `json:"resources"`
	} `json:"aggregations"`
}

func FetchLookupByResourceIDBatch(ctx context.Context, client opengovernance.Client, platformResourceIDs []string) (map[string][]es.LookupResource, error) {
	if len(platformResourceIDs) == 0 {
		return nil, nil
	}
	request := make(map[string]any)
	request["size"] = 0
	request["sort"] = []map[string]any{
		{
			"_id": "desc",
		},
	}
	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": map[string]any{
				"terms": map[string]any{
					"platform_id": platformResourceIDs,
				},
			},
		},
	}
	request["aggs"] = map[string]any{
		"resources": map[string]any{
			"terms": map[string]any{
				"field": "platform_id",
				"size":  len(platformResourceIDs),
			},
			"aggs": map[string]any{
				"hit_select": map[string]any{
					"top_hits": map[string]any{
						"size": 1000,
						"sort": map[string]any{
							"_id": "desc",
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("query=", string(b), "index=", InventorySummaryIndex)

	var response LookupQueryResponse
	err = client.Search(ctx, InventorySummaryIndex, string(b), &response)
	if err != nil {
		return nil, err
	}

	resources := make(map[string][]es.LookupResource)
	for _, bucket := range response.Aggregations.Resources.Buckets {
		if len(bucket.HitSelect.Hits.Hits) == 0 {
			continue
		}
		resources[bucket.Key] = make([]es.LookupResource, 0, len(bucket.HitSelect.Hits.Hits))
		for _, hit := range bucket.HitSelect.Hits.Hits {
			resources[bucket.Key] = append(resources[bucket.Key], hit.Source)
		}
	}

	return resources, nil
}

type ResourceQueryResponse struct {
	Hits struct {
		Total opengovernance.SearchTotal `json:"total"`
		Hits  []struct {
			ID      string      `json:"_id"`
			Score   float64     `json:"_score"`
			Index   string      `json:"_index"`
			Type    string      `json:"_type"`
			Version int64       `json:"_version,omitempty"`
			Source  es.Resource `json:"_source"`
			Sort    []any       `json:"sort"`
		}
	}
}

func FetchResourceByResourceIdAndType(ctx context.Context, client opengovernance.Client, platformResourceID string, resourceType string) (*es.Resource, error) {
	request := make(map[string]any)
	request["size"] = 1
	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": map[string]any{
				"term": map[string]any{
					"platform_id": platformResourceID,
				},
			},
		},
	}
	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	index := es.ResourceTypeToESIndex(resourceType)

	fmt.Println("query=", string(b), "index=", index)

	var response ResourceQueryResponse
	err = client.Search(ctx, index, string(b), &response)
	if err != nil {
		return nil, err
	}

	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	return &response.Hits.Hits[0].Source, nil
}
