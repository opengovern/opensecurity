package models

import "github.com/opengovern/og-util/pkg/integration"

type IntegrationTypeBinaries struct {
	IntegrationType integration.Type `gorm:"primaryKey"`
	URL             string

	IntegrationPlugin []byte `gorm:"type:bytea;not null"`
	CloudQlPlugin     []byte `gorm:"type:bytea;not null"`
}
