package rego_runner

import (
	"context"
	"regexp"
	"strings"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

type Client struct {
	ES es.Client
}

type Metadata struct {
	ID               string
	SubscriptionID   string
	Location         string
	CloudEnvironment string

	Name         string
	AccountID    string
	SourceID     string
	Region       string
	Partition    string
	ResourceType string
}

type ResourceHit struct {
	ID      string                 `json:"_id"`
	Score   float64                `json:"_score"`
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Version int64                  `json:"_version,omitempty"`
	Source  map[string]interface{} `json:"_source"`
	Sort    []any                  `json:"sort"`
}

type ResourceHits struct {
	Total es.SearchTotal `json:"total"`
	Hits  []ResourceHit  `json:"hits"`
}

type ResourceSearchResponse struct {
	PitID string       `json:"pit_id"`
	Hits  ResourceHits `json:"hits"`
}

type ResourcePaginator struct {
	paginator *es.BaseESPaginator
}

func (k Client) NewResourcePaginator(filters []es.BoolFilter, limit *int64, index string) (ResourcePaginator, error) {
	paginator, err := es.NewPaginator(k.ES.ES(), index, filters, limit)
	if err != nil {
		return ResourcePaginator{}, err
	}

	paginator.UpdatePageSize(100)

	p := ResourcePaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ResourcePaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ResourcePaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ResourcePaginator) NextPage(ctx context.Context) ([]map[string]interface{}, error) {
	var response ResourceSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []map[string]interface{}
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var resourceTypeMap = map[string]string{
	"aws::ec2::volumegp3": "aws::ec2::volume",
}

var stopWordsRe = regexp.MustCompile(`\W+`)

func ResourceTypeToESIndex(t string) string {
	if rt, ok := resourceTypeMap[t]; ok {
		t = rt
	}
	t = stopWordsRe.ReplaceAllString(t, "_")
	return strings.ToLower(t)
}
