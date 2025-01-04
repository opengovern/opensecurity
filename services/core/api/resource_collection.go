package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

type ResourceCollectionStatus string

const (
	ResourceCollectionStatusUnknown  ResourceCollectionStatus = ""
	ResourceCollectionStatusActive   ResourceCollectionStatus = "active"
	ResourceCollectionStatusInactive ResourceCollectionStatus = "inactive"
)

type ResourceCollection struct {
	ID          string                                    `json:"id"`
	Name        string                                    `json:"name"`
	Tags        map[string][]string                       `json:"tags"`
	Description string                                    `json:"description"`
	CreatedAt   time.Time                                 `json:"created_at"`
	Status      ResourceCollectionStatus                  `json:"status"`
	Filters     []opengovernance.ResourceCollectionFilter `json:"filters"`

	IntegrationTypes []integration.Type `json:"integration_types,omitempty"`
	LastEvaluatedAt  *time.Time         `json:"last_evaluated_at,omitempty"`
	ResourceCount    *int               `json:"resource_count,omitempty"`
	IntegrationCount *int               `json:"integration_count,omitempty"`
	MetricCount      *int               `json:"metric_count,omitempty"`
}

type ResourceCollectionLandscapeItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURI     string `json:"logo_uri"`
}

type ResourceCollectionLandscapeSubcategory struct {
	ID          string                            `json:"id"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Items       []ResourceCollectionLandscapeItem `json:"items"`
}

type ResourceCollectionLandscapeCategory struct {
	ID            string                                   `json:"id"`
	Name          string                                   `json:"name"`
	Description   string                                   `json:"description"`
	Subcategories []ResourceCollectionLandscapeSubcategory `json:"subcategories"`
}

type ResourceCollectionLandscape struct {
	Categories []ResourceCollectionLandscapeCategory `json:"categories"`
}
