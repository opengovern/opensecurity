package api


// CreateConnectorRequest represents the expected payload for creating or updating a connector.
type CreateConnectorRequest struct {

	ConnectorType    string `json:"connector_type" validate:"required,oneof=oidc"`                                  // 'oidc' is supported for now
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"` // Optional sub-type
	Issuer           string `json:"issuer,omitempty" validate:"omitempty,url"`
	TenantID         string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	ClientID         string `json:"client_id" validate:"required"`
	ClientSecret     string `json:"client_secret" validate:"required"`
	ID               string `json:"id,omitempty"`   // Optional
	Name             string `json:"name,omitempty"` // Optional
}
type UpdateConnectorRequest struct {
	ConnectorID 	string `json:"connector_id" validate:"required"`
	ConnectorType    string `json:"connector_type" validate:"required,oneof=oidc"`                                  // 'oidc' is supported for now
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"` // Optional sub-type
	Issuer           string `json:"issuer,omitempty" validate:"omitempty,url"`
	TenantID         string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	ClientID         string `json:"client_id" validate:"required"`
	ClientSecret     string `json:"client_secret" validate:"required"`
	ID               uint `json:"id,omitempty"`   // Optional
	Name             string `json:"name,omitempty"` // Optional

}

type OIDCConfig struct {
	Issuer       string `json:"issuer,omitempty"`
	TenantID     string `json:"tenantID,omitempty"` // Added TenantID for entraid sub-type
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type GetConnectorsResponse struct {
		ID       uint `json:"id"`
		ConnectorID string `json:"connector_id"`
		Type     string `json:"type"`
		SubType  string `json:"sub_type"`
		Name     string `json:"name"`
		Issuer   string `json:"issuer,omitempty"`
		ClientID string `json:"client_id,omitempty"`
		TenantID string `json:"tenant_id,omitempty"`
		UserCount uint `json:"user_count"`
		CreatedAt any `json:"created_at"`
		LastUpdate any `json:"last_update"`
	}


type ConnectorSubTypes struct {
	ID string  `json:"id"`
	Name string `json:"name"`
}
	type GetSupportedConnectorTypeResponse struct {
		ConnectorType string   `json:"connector_type"`
		SubTypes      []ConnectorSubTypes `json:"sub_types"`

	}
