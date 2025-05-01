package models

import (
	"github.com/opengovern/og-util/pkg/model"
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
