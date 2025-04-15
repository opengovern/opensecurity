package models

import (
	"time"

	"github.com/opengovern/og-util/pkg/integration"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"gorm.io/gorm"
)

type ResourceTypeTag struct {
	model.Tag
	ResourceType string `gorm:"primaryKey; type:citext"`
}

type NamedQueryTag struct {
	model.Tag
	NamedQueryID string `gorm:"primaryKey"`
}

type NamedQueryTagsResult struct {
	Key          string
	UniqueValues pq.StringArray `gorm:"type:text[]"`
}

type NamedQuery struct {
	ID               string         `gorm:"primarykey"`
	IntegrationTypes pq.StringArray `gorm:"type:text[]"`
	Title            string
	Description      string
	QueryID          *string
	Query            *Query `gorm:"foreignKey:QueryID;references:ID;constraint:OnDelete:SET NULL"`
	IsBookmarked     bool
	Tags             []NamedQueryTag `gorm:"foreignKey:NamedQueryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CacheEnabled     bool
	// default is system
	Owner 		 string		`gorm:"type:text:default:system"`
	Visibility	 string         `gorm:"type:text:default:public"`
}

type NamedQueryHistory struct {
	Query      string `gorm:"type:citext; primaryKey"`
	ExecutedAt time.Time
}

type ResourceType struct {
	IntegrationType integration.Type `json:"integration_type" gorm:"index"`
	ResourceType    string           `json:"resource_type" gorm:"primaryKey; type:citext"`
	ResourceLabel   string           `json:"resource_name"`
	ServiceName     string           `json:"service_name" gorm:"index"`
	DoSummarize     bool             `json:"do_summarize"`
	LogoURI         *string          `json:"logo_uri,omitempty"`

	Tags    []ResourceTypeTag   `gorm:"foreignKey:ResourceType;references:ResourceType;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	tagsMap map[string][]string `gorm:"-:all"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ResourceCollectionTag struct {
	model.Tag
	ResourceCollectionID string `gorm:"primaryKey"`
}

type ResourceCollectionStatus string

const (
	ResourceCollectionStatusActive   ResourceCollectionStatus = "active"
	ResourceCollectionStatusInactive ResourceCollectionStatus = "inactive"
)

type ResourceCollection struct {
	ID          string `gorm:"primarykey"`
	Name        string
	FiltersJson pgtype.JSONB `gorm:"type:jsonb"`
	Description string
	Status      ResourceCollectionStatus

	Tags    []ResourceCollectionTag `gorm:"foreignKey:ResourceCollectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	tagsMap map[string][]string     `gorm:"-:all"`

	Created   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Filters []opengovernance.ResourceCollectionFilter `gorm:"-:all"`
}

type ResourceTypeV2 struct {
	IntegrationType integration.Type `gorm:"column:integration_type"`
	ResourceName    string           `gorm:"column:resource_name"`
	ResourceID      string           `gorm:"primaryKey"`
	SteampipeTable  string           `gorm:"column:steampipe_table"`
	Category        string           `gorm:"column:category"`
}

type CategoriesTables struct {
	Category string   `json:"category"`
	Tables   []string `json:"tables"`
}

type RunNamedQueryRunCache struct {
	QueryID    string `gorm:"primaryKey"`
	ParamsHash string `gorm:"primaryKey"`
	LastRun    time.Time
	Result     pgtype.JSONB
}

type NamedQueryWithCacheStatus struct {
	NamedQuery
	LastRun *time.Time `gorm:"column:last_run"`
}
