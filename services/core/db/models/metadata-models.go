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
type Dashboard struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	IsDefault   bool      `json:"is_default"`
	UserID      string    `gorm:"type:text" json:"user_id"`
	Widgets     []Widget  `gorm:"foreignKey:DashboardID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // One-to-many
	Name        string    `gorm:"type:text" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPrivate   bool      `json:"is_private"`
}

type Widget struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:text" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	WidgetType    string         `gorm:"type:text" json:"widget_type"`
	WidgetProps   pgtype.JSONB   `json:"widget_props" gorm:"type:jsonb"`
	RowSpan       int            `json:"row_span"`
	ColumnSpan    int            `json:"column_span"`
	ColumnOffset  int            `json:"column_offset"`
	IsPublic      bool           `json:"is_public"`
	UserID        string         `gorm:"type:text" json:"user_id"`
	DashboardID   string         `gorm:"type:text" json:"dashboard_id"` // Correct: this links to Dashboard.ID
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