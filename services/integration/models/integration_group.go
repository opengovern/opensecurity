package models

type IntegrationGroup struct {
	Name  string `gorm:"primaryKey" json:"name"`
	Query string `json:"query"`
}
