package inventory

import (
	"context"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gitlab.com/keibiengine/keibi-engine/pkg/describe/kafka"

	"gitlab.com/keibiengine/keibi-engine/pkg/cloudservice"

	"gitlab.com/keibiengine/keibi-engine/pkg/describe"
	"gitlab.com/keibiengine/keibi-engine/pkg/inventory/api"
	"gitlab.com/keibiengine/keibi-engine/pkg/inventory/es"
	"gitlab.com/keibiengine/keibi-engine/pkg/keibi-es-sdk"
	"gitlab.com/keibiengine/keibi-engine/pkg/source"
)

func GetCategories(client keibi.Client, provider source.Type, sourceID *string) ([]api.CategoriesResponse, error) {
	var searchAfter []interface{}
	categoryMap := map[string]api.CategoriesResponse{}
	for {
		query, err := es.GetCategoriesQuery(string(provider), sourceID, EsFetchPageSize, searchAfter)
		if err != nil {
			return nil, err
		}

		var response es.CategoriesQueryResponse
		err = client.Search(context.Background(), describe.SourceResourcesSummary, query, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Hits.Hits) == 0 {
			break
		}

		for _, hit := range response.Hits.Hits {
			if v, ok := categoryMap[hit.Source.CategoryName]; ok {
				v.ResourceCount += hit.Source.ResourceCount
				categoryMap[hit.Source.CategoryName] = v
			} else {
				categoryMap[hit.Source.CategoryName] = api.CategoriesResponse{
					CategoryName:     hit.Source.CategoryName,
					ResourceCount:    hit.Source.ResourceCount,
					LastDayCount:     hit.Source.LastDayCount,
					LastWeekCount:    hit.Source.LastWeekCount,
					LastQuarterCount: hit.Source.LastQuarterCount,
					LastYearCount:    hit.Source.LastYearCount,
				}
			}
			searchAfter = hit.Sort
		}
	}

	var res []api.CategoriesResponse
	for _, v := range categoryMap {
		res = append(res, v)
	}

	return res, nil
}

func GetServices(client keibi.Client, provider source.Type, sourceID *string) ([]api.TopServicesResponse, error) {
	var searchAfter []interface{}
	serviceResponse := map[string]api.TopServicesResponse{}
	for {
		query, err := es.FetchServicesQuery(string(provider), sourceID, EsFetchPageSize, searchAfter)
		if err != nil {
			return nil, err
		}

		var response es.FetchServicesQueryResponse
		err = client.Search(context.Background(), describe.SourceResourcesSummary, query, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Hits.Hits) == 0 {
			break
		}

		for _, hit := range response.Hits.Hits {
			if v, ok := serviceResponse[hit.Source.ServiceName]; ok {
				v.ResourceCount += hit.Source.ResourceCount
				serviceResponse[hit.Source.ServiceName] = v
			} else {
				serviceResponse[hit.Source.ServiceName] = api.TopServicesResponse{
					ServiceName:      hit.Source.ServiceName,
					Provider:         string(hit.Source.SourceType),
					ResourceCount:    hit.Source.ResourceCount,
					LastDayCount:     hit.Source.LastDayCount,
					LastWeekCount:    hit.Source.LastWeekCount,
					LastQuarterCount: hit.Source.LastQuarterCount,
					LastYearCount:    hit.Source.LastYearCount,
				}
			}
			searchAfter = hit.Sort
		}
	}

	var res []api.TopServicesResponse
	for _, v := range serviceResponse {
		res = append(res, v)
	}
	return res, nil
}

func GetResources(client keibi.Client, rcache *redis.Client, cache *cache.Cache, provider source.Type, sourceID *string, resourceTypes []string) ([]api.ResourceTypeResponse, error) {
	var providerPtr *string
	if provider != "" {
		v := string(provider)
		providerPtr = &v
	}

	var hits []kafka.SourceResourcesSummary
	for _, resourceType := range resourceTypes {
		if cached, err := es.FetchResourceLastSummaryCached(rcache, cache, providerPtr, sourceID, &resourceType); err == nil && len(cached) > 0 {
			hits = append(hits, cached...)
		} else {
			//TODO-Saleh performance issue: use list of resource types instead
			result, err := es.FetchResourceLastSummary(client, providerPtr, sourceID, &resourceType)
			if err != nil {
				return nil, err
			}
			hits = append(hits, result...)
		}
	}

	resourceTypeResponse := map[string]api.ResourceTypeResponse{}
	for _, hit := range hits {
		if v, ok := resourceTypeResponse[hit.ResourceType]; ok {
			v.ResourceCount += hit.ResourceCount
			resourceTypeResponse[hit.ResourceType] = v
		} else {
			resourceTypeResponse[hit.ResourceType] = api.ResourceTypeResponse{
				ResourceType:     cloudservice.ResourceTypeName(hit.ResourceType),
				ResourceCount:    hit.ResourceCount,
				LastDayCount:     hit.LastDayCount,
				LastWeekCount:    hit.LastWeekCount,
				LastQuarterCount: hit.LastQuarterCount,
				LastYearCount:    hit.LastYearCount,
			}
		}
	}

	var res []api.ResourceTypeResponse
	for _, v := range resourceTypeResponse {
		if v.ResourceCount == 0 {
			continue
		}

		res = append(res, v)
	}
	return res, nil
}
