package models


import (
	"github.com/lib/pq"
	"time"
)


type Query struct {
	ID              string `gorm:"primaryKey"`
	QueryToExecute  string
	IntegrationType pq.StringArray `gorm:"type:text[]"`
	PrimaryTable    *string
	ListOfTables    pq.StringArray `gorm:"type:text[]"`
	Engine          string
	QueryViews      []QueryView      `gorm:"foreignKey:QueryID"`
	NamedQuery     []NamedQuery     `gorm:"foreignKey:QueryID"`
	Parameters      []QueryParameter `gorm:"foreignKey:QueryID"`
	Global          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type QueryParameter struct {
	QueryID  string `gorm:"primaryKey"`
	Key      string `gorm:"primaryKey"`
	Required bool   `gorm:"default:false"`
}
