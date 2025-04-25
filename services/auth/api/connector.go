// Package api defines the data transfer objects (DTOs) used for the
// HTTP API layer of the auth service. These structs represent the expected
// request and response bodies for API interactions.
package api

// Use the correct import path for your project's api.Role if needed elsewhere
// "github.com/opengovern/og-util/pkg/api"
// Import time if needed for any types, e.g. if CreatedAt/LastUpdate change from 'any'

// CreateConnectorRequest represents the expected payload for creating a new Dex connector
// via the HTTP API. Validation tags are included.
type CreateConnectorRequest struct {
	// ConnectorType specifies the type of the connector (e.g., "oidc"). Currently only "oidc" is supported.
	ConnectorType string `json:"connector_type" validate:"required,oneof=oidc"`
	// ConnectorSubType specifies the vendor or configuration variant (e.g., "general", "google-workspace", "entraid"). Optional.
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"`
	// Issuer is the OIDC issuer URL. Required for "general" subtype.
	Issuer string `json:"issuer,omitempty" validate:"omitempty,url"`
	// TenantID is the Azure AD / Entra ID tenant ID. Required for "entraid" subtype.
	TenantID string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	// ClientID is the OAuth2 Client ID obtained from the identity provider.
	ClientID string `json:"client_id" validate:"required"`
	// ClientSecret is the OAuth2 Client Secret obtained from the identity provider.
	ClientSecret string `json:"client_secret" validate:"required"`
	// ID is the desired unique identifier for the connector within Dex. Optional, defaults may apply.
	ID string `json:"id,omitempty"`
	// Name is the desired display name for the connector within Dex. Optional, defaults may apply.
	Name string `json:"name,omitempty"`
}

// CreateAuth0ConnectorRequest defines the specialized payload for creating an Auth0 OIDC connector.
// It includes fields specific to Auth0 and URIs for configuring Dex OAuth clients.
type CreateAuth0ConnectorRequest struct {
	// Issuer is the Auth0 tenant's OIDC issuer URL (e.g., "https://your-domain.auth0.com/"). Optional if Domain is provided.
	Issuer string `json:"issuer,omitempty" validate:"omitempty,url"`
	// ClientID is the Auth0 Application's Client ID.
	ClientID string `json:"client_id" validate:"required"`
	// ClientSecret is the Auth0 Application's Client Secret.
	ClientSecret string `json:"client_secret" validate:"required"`
	// Domain is the Auth0 tenant domain (e.g., "your-domain.auth0.com").
	Domain string `json:"domain" validate:"required"`
	// PublicURIS is the list of redirect URIs to configure for the Dex 'public-client' OAuth client.
	PublicURIS []string `json:"public_uris" validate:"required"`
	// PrivateURIS is the list of redirect URIs to configure for the Dex 'private-client' OAuth client.
	PrivateURIS []string `json:"private_uris" validate:"required"`
}

// UpdateConnectorRequest represents the expected payload for updating an existing Dex connector's configuration.
type UpdateConnectorRequest struct {
	// ID is the local database primary key (uint) identifying the connector metadata record to update. Optional? Seems required by handler logic.
	ID uint `json:"id,omitempty" validate:"required"` // Changed from string based on handler logic using uint ID
	// ConnectorID is the unique identifier (string) of the connector within Dex that needs updating.
	ConnectorID string `json:"connector_id" validate:"required"`
	// ConnectorType specifies the type of the connector (e.g., "oidc"). Must match the existing connector.
	ConnectorType string `json:"connector_type" validate:"required,oneof=oidc"`
	// ConnectorSubType specifies the vendor or configuration variant.
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"`
	// Issuer is the new OIDC issuer URL. Required for "general" subtype.
	Issuer string `json:"issuer,omitempty" validate:"omitempty,url"`
	// TenantID is the new Azure AD / Entra ID tenant ID. Required for "entraid" subtype.
	TenantID string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	// ClientID is the new OAuth2 Client ID.
	ClientID string `json:"client_id" validate:"required"`
	// ClientSecret is the new OAuth2 Client Secret.
	ClientSecret string `json:"client_secret" validate:"required"`
	// Name is the optional new display name (Note: Dex gRPC update doesn't support changing name directly).
	Name string `json:"name,omitempty"`
}

// OIDCConfig represents a subset of fields that might be parsed from a Dex OIDC connector's
// JSON configuration, primarily for display or informational purposes in API responses.
// It does not represent the full configuration structure used internally or stored in Dex.
type OIDCConfig struct {
	// Issuer is the OIDC provider issuer URL.
	Issuer string `json:"issuer,omitempty"`
	// TenantID is the Azure AD / Entra ID tenant ID, specific to the 'entraid' subtype.
	TenantID string `json:"tenantID,omitempty"`
	// ClientID is the OAuth2 client ID configured for the connector.
	ClientID string `json:"clientID"`
	// ClientSecret is the OAuth2 client secret configured for the connector.
	ClientSecret string `json:"clientSecret"`
	// Note: Other fields like RedirectURIs, Scopes etc. exist in the actual Dex config
	// but might not be needed in this specific API response struct.
}

// GetConnectorsResponse defines the structure for representing a configured connector
// in list responses, combining data from Dex and the local database.
type GetConnectorsResponse struct {
	// ID is the unique database identifier for the local connector metadata record.
	ID uint `json:"id"`
	// ConnectorID is the unique identifier used by Dex for this connector.
	ConnectorID string `json:"connector_id"`
	// Type is the connector type (e.g., "oidc").
	Type string `json:"type"`
	// SubType is the specific subtype (e.g., "general", "entraid"), stored locally.
	SubType string `json:"sub_type"`
	// Name is the display name of the connector configured in Dex.
	Name string `json:"name"`
	// Issuer is the OIDC issuer URL (if applicable and parsed from config).
	Issuer string `json:"issuer,omitempty"`
	// ClientID is the OAuth2 Client ID (if applicable and parsed from config).
	ClientID string `json:"client_id,omitempty"`
	// TenantID is the Azure AD / Entra ID Tenant ID (if applicable and parsed from config).
	TenantID string `json:"tenant_id,omitempty"`
	// UserCount is a locally stored count of users associated with this connector (update mechanism TBD).
	UserCount uint `json:"user_count"`
	// CreatedAt is the timestamp when the local connector metadata record was created. Type 'any' suggests potential inconsistency; prefer time.Time.
	CreatedAt any `json:"created_at"` // Consider changing 'any' to time.Time
	// LastUpdate is the timestamp when the local connector metadata record was last updated. Type 'any' suggests potential inconsistency; prefer time.Time.
	LastUpdate any `json:"last_update"` // Consider changing 'any' to time.Time
}

// ConnectorSubTypes defines the structure for representing a supported connector subtype.
type ConnectorSubTypes struct {
	// ID is the technical identifier for the subtype (e.g., "google-workspace").
	ID string `json:"id"`
	// Name is the user-friendly display name for the subtype (e.g., "Google Workspaces").
	Name string `json:"name"`
}

// GetSupportedConnectorTypeResponse defines the structure for listing supported connector types
// and their available subtypes.
type GetSupportedConnectorTypeResponse struct {
	// ConnectorType is the main type identifier (e.g., "oidc").
	ConnectorType string `json:"connector_type"`
	// SubTypes is a list of supported subtypes for this ConnectorType.
	SubTypes []ConnectorSubTypes `json:"sub_types"`
}
