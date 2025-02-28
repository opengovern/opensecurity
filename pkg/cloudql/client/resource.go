package opengovernance_client

import (
	"context"
	"regexp"
	"runtime"
	"strings"

	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

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

type Resource struct {
	PlatformID      string   `json:"platform_id"`
	Description     any      `json:"description"`
	Metadata        Metadata `json:"metadata"`
	ResourceID      string   `json:"resource_id"`
	IntegrationID   string   `json:"integration_id"`
	IntegrationType string   `json:"integration_type"`
	ResourceType    string   `json:"resource_type"`
	ResourceName    string   `json:"resource_name"`
	DescribedBy     string   `json:"described_by"`
	DescribedAt     int64    `json:"described_at"`
}

type ResourceHit struct {
	ID      string   `json:"_id"`
	Score   float64  `json:"_score"`
	Index   string   `json:"_index"`
	Type    string   `json:"_type"`
	Version int64    `json:"_version,omitempty"`
	Source  Resource `json:"_source"`
	Sort    []any    `json:"sort"`
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

func (p ResourcePaginator) NextPage(ctx context.Context) ([]Resource, error) {
	var response ResourceSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Resource
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

var resourceMapping = map[string]string{
	"platform_id":      "platform_id",
	"resource_id":      "resource_id",
	"integration_id":   "integration_id",
	"integration_type": "integration_type",
	"resource_type":    "resource_type",
	"resource_name":    "resource_name",
	"described_by":     "described_by",
	"described_at":     "described_at",
	"name":             "metadata.Name",
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

func ListResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListResources 1", d)
	runtime.GC()
	// create service
	cfg := config.GetConfig(d.Connection)

	plugin.Logger(ctx).Trace("ListResources 2", cfg)
	ke, err := config.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{ES: ke}

	plugin.Logger(ctx).Trace("ListResources 3", k)
	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListResources NewSelfClientCached", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("ListResources 4", sc)
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListResources GetConfigTableValueOrNil for resource_collection_filters", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("ListResources 5", encodedResourceCollectionFilters)
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListResources GetConfigTableValueOrNil for client_type", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Columns", d.EqualsQuals)
	var indexes []string
	for column, q := range d.EqualsQuals {
		if column == "resource_type" {
			if s, ok := q.GetValue().(*proto.QualValue_StringValue); ok && s != nil {
				indexes = []string{ResourceTypeToESIndex(s.StringValue)}
			} else if l := q.GetListValue(); l != nil {
				for _, v := range l.GetValues() {
					if v == nil {
						continue
					}
					indexes = append(indexes, v.GetStringValue())
				}
			}
		}
	}

	for _, index := range indexes {
		paginator, err := k.NewResourcePaginator(es.BuildFilterWithDefaultFieldName(ctx, d.QueryContext, resourceMapping,
			nil, encodedResourceCollectionFilters, clientType, true), d.QueryContext.Limit, index)
		if err != nil {
			plugin.Logger(ctx).Error("ListResources NewResourcePaginator", "error", err)
			return nil, err
		}

		for paginator.HasNext() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("ListResources NextPage", "error", err)
				return nil, err
			}
			plugin.Logger(ctx).Trace("ListResources", "next page")

			for _, v := range page {
				d.StreamListItem(ctx, v)
			}
		}
		err = paginator.Close(ctx)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
