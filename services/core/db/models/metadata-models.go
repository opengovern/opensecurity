package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/model"
)

// Metadata models

type Filter struct {
	Name     string            `json:"name" gorm:"primary_key"`
	KeyValue map[string]string `json:"kayValue" gorm:"key_values"`
}


type ConfigMetadata struct {
	Key   MetadataKey        `json:"key" gorm:"primary_key"`
	Type  ConfigMetadataType `json:"type" gorm:"default:'string'"`
	Value string             `json:"value" gorm:"type:text;not null"`
}
type PlatformConfiguration struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	InstallID  uuid.UUID `json:"install_id"`
	Configured bool      `json:"configured"`
}

type PolicyParameterValues struct {
	Key       string `gorm:"primaryKey"`
	ControlID string `gorm:"primaryKey"`
	Value     string `gorm:"type:text;not null"`
}


type QueryViewTag struct {
	model.Tag
	QueryViewID string `gorm:"primaryKey"`
}

type QueryView struct {
	ID           string `json:"id" gorm:"primary_key"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	QueryID      *string
	Query        *Query         `gorm:"foreignKey:QueryID;references:ID;constraint:OnDelete:SET NULL"`
	Dependencies pq.StringArray `gorm:"type:text[]"`
	Tags         []QueryViewTag `gorm:"foreignKey:QueryViewID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type UserLayout struct{
	UserID string `gorm:"primaryKey" ,json:"user_id"`
	LayoutConfig pgtype.JSONB `gorm:"type:jsonb" ,json:"layout_config"`
	Name string `gorm:"type:text" ,json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPrivate bool `json:"is_private"`
}
// Array of widgets 
// example: 
// [{"id":"table1","type":"table","title":"Table 1","query_id":"query1","layout":{"x":0,"y":0,"w":6,"h":4}},{"id":"table2","type":"table","title":"Table 2","query_id":"query2","layout":{"x":6,"y":0,"w":6,"h":4}}]}]
// 


type MetadataKey string
type ConfigMetadataType string


type StringConfigMetadata struct {
	ConfigMetadata
}
type JSONConfigMetadata struct {
	ConfigMetadata
	Value any
}
type BoolConfigMetadata struct {
	ConfigMetadata
	Value bool
}
type IntConfigMetadata struct {
	ConfigMetadata
	Value int
}