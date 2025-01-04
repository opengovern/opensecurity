package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type ResourceType struct {
	IntegrationType integration.Type    `json:"integration_type" example:"Azure"`                                                                                                     // Cloud Provider
	ResourceType    string              `json:"resource_type" example:"Microsoft.Compute/virtualMachines"`                                                                            // Resource Type
	ResourceLabel   string              `json:"resource_name" example:"VM"`                                                                                                           // Resource Name
	ServiceName     string              `json:"service_name" example:"compute"`                                                                                                       // Service Name
	Tags            map[string][]string `json:"tags,omitempty" swaggertype:"array,string" example:"category:[Data and Analytics,Database,Integration,Management Governance,Storage]"` // Tags
	LogoURI         *string             `json:"logo_uri,omitempty" example:"https://opengovernance.io/logo.png"`                                                                      // Logo URI

	Count    *int `json:"count" example:"100" minimum:"0"`    // Number of Resources of this Resource Type - Metric
	OldCount *int `json:"old_count" example:"90" minimum:"0"` // Number of Resources of this Resource Type in the past - Metric

	ComplianceCount *int     `json:"compliance_count" minimum:"0"` // Number of Compliance that use this Resource Type - Metadata
	Compliance      []string `json:"compliance"`                   // List of Compliance that support this Resource Type - Metadata (GET only)
	Attributes      []string `json:"attributes"`                   // List supported steampipe Attributes (columns) for this resource type - Metadata (GET only)
}

type ResourceTypeV2 struct {
	IntegrationType integration.Type `json:"integration_type"`
	ResourceName    string           `json:"resource_name"`
	ResourceID      string           `json:"resource_id"`
	SteampipeTable  string           `json:"steampipe_table"`
	Category        string           `json:"category"`
}

type CategoriesTables struct {
	Category string   `json:"category"`
	Tables   []string `json:"tables"`
}

type ListResourceTypeMetadataResponse struct {
	TotalResourceTypeCount int            `json:"total_resource_type_count" example:"100" minimum:"0"`
	ResourceTypes          []ResourceType `json:"resource_types"`
}

type CountPair struct {
	OldCount int `json:"old_count" minimum:"0"`
	Count    int `json:"count" minimum:"0"`
}

type ListResourceTypeCompositionResponse struct {
	TotalCount      int                  `json:"total_count" minimum:"0"`
	TotalValueCount int                  `json:"total_value_count" minimum:"0"`
	TopValues       map[string]CountPair `json:"top_values"`
	Others          CountPair            `json:"others"`
}

type ResourceTypeTrendDatapoint struct {
	Count                                   int                        `json:"count" example:"100" minimum:"0"`
	CountStacked                            []ResourceCountStackedItem `json:"countStacked"`
	TotalDescribedConnectionCount           int64                      `json:"totalIntegrationCount"`
	TotalSuccessfulDescribedConnectionCount int64                      `json:"totalSuccessfulDescribedConnectionCount"`
	Date                                    time.Time                  `json:"date" format:"date-time"`
}

type LocationResponse struct {
	Location         string `json:"location" example:"na-west"`                        // Region
	ResourceCount    *int   `json:"resourceCount,omitempty" example:"100" minimum:"0"` // Number of resources in the region
	ResourceOldCount *int   `json:"resourceOldCount,omitempty" example:"50" minimum:"0"`
}

type RegionsResourceCountResponse struct {
	TotalCount int                `json:"totalCount" minimum:"0"`
	Regions    []LocationResponse `json:"regions"`
}

type AnalyticsCategoriesResponse struct {
	CategoryResourceType map[string][]string `json:"categoryResourceType"`
}

type ResourceCategory struct {
	Category  string           `json:"category"`
	Resources []ResourceTypeV2 `json:"resources"`
}

type GetResourceCategoriesResult struct {
	Categories []ResourceCategory `json:"categories"`
}

type ServiceQueries struct {
	Service string
	Queries []NamedQueryItemV2 `json:"queries"`
}

type CategoryQueries struct {
	Category string           `json:"category"`
	Services []ServiceQueries `json:"services"`
}

type GetCategoriesControlsResponse struct {
	Categories []CategoryQueries `json:"categories"`
}

type GetResourceCategoriesTables struct {
	Name         string `json:"name"`
	Table        string `json:"table"`
	ResourceType string `json:"resource_type"`
}

type GetResourceCategoriesCategory struct {
	Category string                        `json:"category"`
	Tables   []GetResourceCategoriesTables `json:"tables"`
}

type GetResourceCategoriesResponse struct {
	Categories []GetResourceCategoriesCategory `json:"categories"`
}
