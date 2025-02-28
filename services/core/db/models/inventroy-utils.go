package models

import (
	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/opensecurity/services/core/api"
)





func (p NamedQuery) GetTagsMap() map[string][]string {
	var tagsMap map[string][]string
	if p.Tags != nil {
		tagLikeArr := make([]model.TagLike, 0, len(p.Tags))
		for _, tag := range p.Tags {
			tagLikeArr = append(tagLikeArr, tag)
		}
		tagsMap = model.GetTagsMap(tagLikeArr)
	}
	return tagsMap
}

func (s NamedQueryTagsResult) ToApi() api.NamedQueryTagsResult {
	return api.NamedQueryTagsResult{
		Key:          s.Key,
		UniqueValues: s.UniqueValues,
	}
}

func (s NamedQueryHistory) ToApi() api.NamedQueryHistory {
	return api.NamedQueryHistory{
		Query:      s.Query,
		ExecutedAt: s.ExecutedAt,
	}
}
func (r ResourceType) GetTagsMap() map[string][]string {
	if r.tagsMap == nil {
		tagLikeArr := make([]model.TagLike, 0, len(r.Tags))
		for _, tag := range r.Tags {
			tagLikeArr = append(tagLikeArr, tag)
		}
		r.tagsMap = model.GetTagsMap(tagLikeArr)
	}
	return r.tagsMap
}

func (r ResourceCollectionStatus) ToApi() api.ResourceCollectionStatus {
	switch r {
	case ResourceCollectionStatusActive:
		return api.ResourceCollectionStatusActive
	case ResourceCollectionStatusInactive:
		return api.ResourceCollectionStatusInactive
	default:
		return api.ResourceCollectionStatusUnknown
	}
}

func (r ResourceCollection) ToApi() api.ResourceCollection {
	apiResourceCollection := api.ResourceCollection{
		ID:          r.ID,
		Name:        r.Name,
		Tags:        model.TrimPrivateTags(r.GetTagsMap()),
		Description: r.Description,
		CreatedAt:   r.Created,
		Status:      r.Status.ToApi(),
		Filters:     r.Filters,
	}
	return apiResourceCollection
}


func (r ResourceCollection) GetTagsMap() map[string][]string {
	if r.tagsMap == nil {
		tagLikeArr := make([]model.TagLike, 0, len(r.Tags))
		for _, tag := range r.Tags {
			tagLikeArr = append(tagLikeArr, tag)
		}
		r.tagsMap = model.GetTagsMap(tagLikeArr)
	}
	return r.tagsMap
}

func (r CategoriesTables) ToApi() api.CategoriesTables {
	apiResourceType := api.CategoriesTables{
		Tables:   r.Tables,
		Category: r.Category,
	}
	return apiResourceType
}


func (r ResourceTypeV2) ToApi() api.ResourceTypeV2 {
	apiResourceType := api.ResourceTypeV2{
		IntegrationType: r.IntegrationType,
		ResourceName:    r.ResourceName,
		ResourceID:      r.ResourceID,
		SteampipeTable:  r.SteampipeTable,
		Category:        r.Category,
	}
	return apiResourceType
}
