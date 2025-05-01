package models

import (
	"time"

	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/model"
)

type Query struct {
	ID              string `gorm:"primaryKey"`
	QueryToExecute  string
	IntegrationType pq.StringArray `gorm:"type:text[]"`
	PrimaryTable    *string
	ListOfTables    pq.StringArray `gorm:"type:text[]"`
	Engine          string
	QueryViews      []QueryView `gorm:"foreignKey:QueryID"`
	//NamedQuery      []NamedQuery     `gorm:"foreignKey:QueryID"`
	Parameters []QueryParameter `gorm:"foreignKey:QueryID"`
	Global     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type QueryParameter struct {
	QueryID  string `gorm:"primaryKey"`
	Key      string `gorm:"primaryKey"`
	Required bool   `gorm:"default:false"`
}
type QueryViewTag struct {
	model.Tag
	QueryViewID string `gorm:"primaryKey"`
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
	IsView           bool
	// default is system
	Owner      string `gorm:"type:text;default:system"`
	Visibility string `gorm:"type:text;default:public"`
}
type NamedQueryTag struct {
	model.Tag
	NamedQueryID string `gorm:"primaryKey"`
}

type NamedQueryTagsResult struct {
	Key          string
	UniqueValues pq.StringArray `gorm:"type:text[]"`
}
type NamedQueryWithCacheStatus struct {
	NamedQuery
	LastRun *time.Time `gorm:"column:last_run"`
}
