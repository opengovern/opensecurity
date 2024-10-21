package interfaces

import (
	"github.com/opengovern/opengovernance/services/integration-v2/models"
)

type CredentialType interface {
	HealthCheck() error
	DiscoverIntegrations() ([]models.Integration, error)
}

// IntegrationCreator CredentialType interface, credentials, error
type CredentialCreator func(jsonData []byte) (CredentialType, error)
