package models

import (
	"time"

	"github.com/opengovern/og-util/pkg/integration"

	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/model"
	"gorm.io/gorm"
)

type ResourceTypeTag struct {
	model.Tag
	ResourceType string `gorm:"primaryKey; type:citext"`
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
