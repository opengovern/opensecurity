package models



import (
	"github.com/google/uuid"
	"time"
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

type QueryParameterValues struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"type:text;not null"`
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