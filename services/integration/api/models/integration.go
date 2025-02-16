package models

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type IntegrationState string

const (
	IntegrationStateActive   IntegrationState = "ACTIVE"
	IntegrationStateInactive IntegrationState = "INACTIVE"
	IntegrationStateArchived IntegrationState = "ARCHIVED"
	IntegrationStateSample   IntegrationState = "SAMPLE_INTEGRATION"
)

type DiscoverIntegrationRequest struct {
	IntegrationType integration.Type `json:"integration_type"`
	CredentialType  string           `json:"credential_type"`
	Description     string           `json:"description"`
	CredentialID    *string          `json:"credential_id"`
	Credentials     map[string]any   `json:"credentials"`
}

type DiscoverIntegrationResponse struct {
	CredentialID string        `json:"credential_id"`
	Integrations []Integration `json:"integrations"`
}

type AddIntegrationsRequest struct {
	IntegrationType integration.Type `json:"integration_type"`
	ProviderIDs     []string         `json:"provider_ids"`
	CredentialID    string           `json:"credential_id"`
}

type UpdateRequest struct {
	Credentials map[string]any `json:"credentials"`
	Description string         `json:"description"`
}

type Integration struct {
	IntegrationID   string            `json:"integration_id"`
	ProviderID      string            `json:"provider_id"`
	Name            string            `json:"name"`
	IntegrationType integration.Type  `json:"integration_type"`
	Annotations     map[string]string `json:"annotations"`
	Labels          map[string]string `json:"labels"`

	CredentialID string `json:"credential_id"`

	State     IntegrationState `json:"state"`
	LastCheck *time.Time       `json:"last_check,omitempty"`
}

type ListIntegrationsResponse struct {
	Integrations []Integration `json:"integrations"`
	TotalCount   int           `json:"total_count"`
}

type ListIntegrationsRequest struct {
	IntegrationID   []string `json:"integration_id"`
	IntegrationType []string `json:"integration_type"`
	ProviderIDRegex *string  `json:"provider_id_regex"`
	NameRegex       *string  `json:"integration_name_regex"`
	Cursor          *int64   `json:"cursor"`
	PerPage         *int64   `json:"per_page"`
}

type SetResourceTypesForIntegration struct {
	ResourceTypes []string `json:"resource_types"`
}
