// Package utils provides helper functions and shared logic for the auth service,
// including database interactions, connector configurations, and Kubernetes interactions.
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time" // Added for http client timeout

	dexapi "github.com/dexidp/dex/api/v2" // Correct Dex API v2 import path
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateConnectorRequest defines the structure for requesting the creation of a new Dex connector,
// typically via the HTTP API. It includes validation tags for use with validators like validator.v9.
type CreateConnectorRequest struct {
	ConnectorType    string `json:"connector_type" validate:"required,oneof=oidc"`                                  // Type of connector, currently only 'oidc'.
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"` // Specific subtype (e.g., 'entraid'). Determines required fields. Optional, defaults might apply.
	Issuer           string `json:"issuer,omitempty" validate:"omitempty,url"`                                      // OIDC issuer URL (required for 'general', derived for others if possible).
	TenantID         string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`                                  // Azure/Entra Tenant ID (required for 'entraid').
	ClientID         string `json:"client_id" validate:"required"`                                                  // OAuth2 Client ID provided by the IdP.
	ClientSecret     string `json:"client_secret" validate:"required"`                                              // OAuth2 Client Secret provided by the IdP.
	ID               string `json:"id,omitempty"`                                                                   // Desired Dex Connector ID (optional, defaults apply based on subtype).
	Name             string `json:"name,omitempty"`                                                                 // Desired Dex Connector Name (optional, defaults apply based on subtype).
}

// CreateAuth0ConnectorRequest defines the specialized structure for creating an Auth0 OIDC connector.
type CreateAuth0ConnectorRequest struct {
	Issuer       string `json:"issuer,omitempty" validate:"omitempty,url"` // Auth0 Issuer URL (e.g., https://your-domain.auth0.com/). Optional if Domain is provided.
	ClientID     string `json:"client_id" validate:"required"`             // Auth0 Application Client ID.
	ClientSecret string `json:"client_secret" validate:"required"`         // Auth0 Application Client Secret.
	Domain       string `json:"domain" validate:"required"`                // Auth0 Domain (e.g., your-domain.auth0.com). Used for configuration if Issuer is empty.
	// PublicURIS and PrivateURIS were part of the API request in http.go, not directly used in this util function's logic
	// PublicURIS   []string `json:"public_uris" validate:"required"`           // Redirect URIs for the Dex 'public-client'.
	// PrivateURIS  []string `json:"private_uris" validate:"required"`          // Redirect URIs for the Dex 'private-client'.
}

// UpdateConnectorRequest defines the structure for requesting updates to an existing Dex connector's configuration.
type UpdateConnectorRequest struct {
	// ID field from the original struct seems to be the Dex Connector ID string, used to identify the connector to update.
	ID string `json:"id" validate:"required"` // The ID of the Dex connector to update.

	ConnectorType    string `json:"connector_type" validate:"required,oneof=oidc"`                                  // Type of connector (currently only 'oidc'). Must match existing connector.
	ConnectorSubType string `json:"connector_sub_type" validate:"omitempty,oneof=general google-workspace entraid"` // Subtype for validation/logic.
	Issuer           string `json:"issuer,omitempty" validate:"omitempty,url"`                                      // New OIDC issuer URL (required for 'general').
	TenantID         string `json:"tenant_id,omitempty" validate:"omitempty,uuid"`                                  // New Azure/Entra Tenant ID (required for 'entraid').
	ClientID         string `json:"client_id" validate:"required"`                                                  // New OAuth2 Client ID.
	ClientSecret     string `json:"client_secret" validate:"required"`                                              // New OAuth2 Client Secret.
	// Name update via gRPC API is not supported by Dex (UpdateConnector only takes config).
	// Name             string `json:"name,omitempty"` // Optional new name (cannot be updated via Dex gRPC).
}

// OIDCConfig represents the structure of the JSON configuration stored within a Dex OIDC connector.
// Fields correspond to Dex's OIDC connector configuration options.
type OIDCConfig struct {
	Issuer                    string   `json:"issuer,omitempty"`          // OIDC provider issuer URL.
	TenantID                  string   `json:"tenantID,omitempty"`        // Specific field for Entra ID subtype.
	ClientID                  string   `json:"clientID"`                  // OAuth2 client ID.
	ClientSecret              string   `json:"clientSecret"`              // OAuth2 client secret.
	Name                      string   `json:"name,omitempty"`            // Display name (often set on Connector directly).
	RedirectURIs              []string `json:"redirectURIs,omitempty"`    // Correct field name for Dex config (list of allowed callback URLs for Dex).
	RedirectURI               string   `json:"redirectURI,omitempty"`     // Singular redirect URI (often derived, but included for parity with original code).
	InsecureEnableGroups      bool     `json:"insecureEnableGroups"`      // If true, enables the 'groups' scope. Usually requires Dex RBAC setup.
	InsecureSkipEmailVerified bool     `json:"insecureSkipEmailVerified"` // If true, allows users with unverified emails. Use with caution.
}

// ConnectorCreator defines a function type that takes creation parameters and returns
// a Dex gRPC CreateConnectorReq message or an error. Used for factory pattern.
type ConnectorCreator func(params CreateConnectorRequest) (*dexapi.CreateConnectorReq, error)

// connectorCreators maps connector types (lowercase string) to their corresponding ConnectorCreator functions.
var connectorCreators = map[string]ConnectorCreator{
	"oidc": CreateOIDCConnector,
	// Future connector types can be added here, e.g., "saml": CreateSAMLConnector
}

// SupportedConnectors maps connector types to a list of supported subtype identifiers (lowercase).
// Used for validation and providing options to the user.
var SupportedConnectors = map[string][]string{
	"oidc": {"general", "google-workspace", "entraid"},
	// Add more connector types and their sub-types here as needed.
}

// SupportedConnectorsNames maps connector types to a list of user-friendly names for supported subtypes.
// The order should correspond to the order in SupportedConnectors.
var SupportedConnectorsNames = map[string][]string{
	"oidc": {"General OIDC", "Google Workspaces", "AzureAD/EntraID"},
}

// CreateOIDCConnector builds a Dex CreateConnectorReq for an OIDC connector based on the provided parameters.
// It handles logic specific to subtypes like 'general', 'entraid', and 'google-workspace',
// setting default values and fetching required information like the Entra ID issuer URL.
// Redirect URIs are read from the DEX_CALLBACK_URL environment variable.
func CreateOIDCConnector(params CreateConnectorRequest) (*dexapi.CreateConnectorReq, error) {
	var oidcConfig OIDCConfig
	connectorID := params.ID     // Use provided ID
	connectorName := params.Name // Use provided Name

	// Read common redirect URI(s) from environment variable
	// Ensure DEX_CALLBACK_URL is set appropriately during deployment.
	dexCallbackURL := os.Getenv("DEX_CALLBACK_URL")
	if dexCallbackURL == "" {
		// Return an error if the crucial callback URL is missing
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable must be set")
	}
	redirectURIs := strings.Split(dexCallbackURL, ",")
	// Ensure no empty strings if multiple commas are used
	validRedirectURIs := []string{}
	for _, uri := range redirectURIs {
		trimmed := strings.TrimSpace(uri)
		if trimmed != "" {
			validRedirectURIs = append(validRedirectURIs, trimmed)
		}
	}
	if len(validRedirectURIs) == 0 {
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable is set but contains no valid URIs after splitting by comma")
	}
	// Set both singular and plural RedirectURI fields for parity with original code
	redirectURI := validRedirectURIs[0] // Use the first valid URI as the singular one

	// Default ID and Name based on subtype if not provided in request
	if connectorID == "" {
		switch params.ConnectorSubType {
		case "entraid":
			connectorID = "entra-id"
		case "google-workspace":
			connectorID = "google-oidc"
		default:
			connectorID = fmt.Sprintf("oidc-%s", strings.ToLower(params.Name)) // Fallback based on name
		}
		if connectorID == "oidc-" {
			connectorID = "oidc-default"
		} // Further fallback
	}
	if connectorName == "" {
		switch params.ConnectorSubType {
		case "entraid":
			connectorName = "AzureAD/EntraID"
		case "google-workspace":
			connectorName = "Google Workspace"
		default:
			connectorName = "OIDC" // Generic default
		}
		if params.Name != "" {
			connectorName = params.Name
		} // Use provided name if default wasn't specific
	}

	// Build OIDCConfig based on subtype
	switch params.ConnectorSubType {
	case "general":
		if params.Issuer == "" {
			return nil, fmt.Errorf("issuer URL is required for 'general' OIDC connector subtype")
		}
		oidcConfig = OIDCConfig{
			Issuer:                    params.Issuer,
			ClientID:                  params.ClientID,
			ClientSecret:              params.ClientSecret,
			RedirectURIs:              validRedirectURIs, // Use validated slice
			RedirectURI:               redirectURI,       // Set singular field
			InsecureEnableGroups:      true,
			InsecureSkipEmailVerified: true,
		}
	case "entraid":
		if params.TenantID == "" {
			return nil, fmt.Errorf("tenant_id is required for 'entraid' OIDC connector subtype")
		}
		if params.Issuer == "" {
			issuer, err := fetchEntraIDIssuer(params.TenantID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch issuer for entraid tenant %s: %w", params.TenantID, err)
			}
			params.Issuer = issuer
		}
		oidcConfig = OIDCConfig{
			Issuer:                    params.Issuer,
			TenantID:                  params.TenantID,
			ClientID:                  params.ClientID,
			ClientSecret:              params.ClientSecret,
			RedirectURIs:              validRedirectURIs, // Use validated slice
			RedirectURI:               redirectURI,       // Set singular field
			InsecureEnableGroups:      true,
			InsecureSkipEmailVerified: true,
		}
	case "google-workspace":
		oidcConfig = OIDCConfig{
			Issuer:                    "https://accounts.google.com",
			ClientID:                  params.ClientID,
			ClientSecret:              params.ClientSecret,
			RedirectURIs:              validRedirectURIs, // Use validated slice
			RedirectURI:               redirectURI,       // Set singular field
			InsecureEnableGroups:      true,
			InsecureSkipEmailVerified: true,
		}
	default:
		return nil, fmt.Errorf("unsupported or missing connector_sub_type: '%s'", params.ConnectorSubType)
	}

	// Serialize the specific OIDCConfig to JSON bytes for Dex.
	configBytes, err := json.Marshal(oidcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OIDC config: %w", err)
	}

	// Construct the Dex Connector message.
	connector := &dexapi.Connector{
		Id:     connectorID,
		Type:   "oidc", // Type is fixed for this function
		Name:   connectorName,
		Config: configBytes,
	}

	// Create the final gRPC request message.
	req := &dexapi.CreateConnectorReq{
		Connector: connector,
	}

	return req, nil
}

// CreateAuth0Connector builds a Dex CreateConnectorReq specifically for Auth0 OIDC.
// It uses a fixed connector ID "auth0" and name "Auth0".
// Redirect URIs are read from the DEX_CALLBACK_URL environment variable.
func CreateAuth0Connector(params CreateAuth0ConnectorRequest) (*dexapi.CreateConnectorReq, error) {
	connectorID := "auth0"
	connectorName := "Auth0"

	// Use provided Issuer, or derive from Domain if Issuer is missing
	issuer := params.Issuer
	if issuer == "" {
		if params.Domain == "" {
			return nil, fmt.Errorf("either issuer or domain must be provided for Auth0 connector")
		}
		issuer = fmt.Sprintf("https://%s/", strings.TrimSuffix(params.Domain, "/"))
	}

	// Read Dex's callback URL from environment.
	dexCallbackURL := os.Getenv("DEX_CALLBACK_URL")
	if dexCallbackURL == "" {
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable must be set for Auth0 connector")
	}
	redirectURIs := strings.Split(dexCallbackURL, ",")
	validRedirectURIs := []string{}
	for _, uri := range redirectURIs {
		trimmed := strings.TrimSpace(uri)
		if trimmed != "" {
			validRedirectURIs = append(validRedirectURIs, trimmed)
		}
	}
	if len(validRedirectURIs) == 0 {
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable contains no valid URIs")
	}
	redirectURI := validRedirectURIs[0] // Use first as singular

	oidcConfig := OIDCConfig{
		Issuer:                    issuer,
		ClientID:                  params.ClientID,
		ClientSecret:              params.ClientSecret,
		RedirectURIs:              validRedirectURIs, // Dex callback URI(s)
		RedirectURI:               redirectURI,       // Set singular field
		InsecureEnableGroups:      true,
		InsecureSkipEmailVerified: true,
	}

	configBytes, err := json.Marshal(oidcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Auth0 OIDC config: %w", err)
	}

	connector := &dexapi.Connector{Id: connectorID, Type: "oidc", Name: connectorName, Config: configBytes}
	req := &dexapi.CreateConnectorReq{Connector: connector}

	return req, nil
}

// fetchEntraIDIssuer retrieves the issuer URL from the Entra ID (Azure AD) OpenID configuration endpoint
// for a given tenant ID using an HTTP GET request.
func fetchEntraIDIssuer(tenantID string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", tenantID)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Entra ID OpenID configuration from %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d when fetching Entra ID OpenID configuration from %s: %s", resp.StatusCode, url, string(bodyBytes))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Entra ID OpenID configuration response: %w", err)
	}
	var config struct {
		Issuer string `json:"issuer"`
	}
	if err := json.Unmarshal(body, &config); err != nil {
		return "", fmt.Errorf("failed to parse Entra ID OpenID configuration JSON: %w", err)
	}
	if config.Issuer == "" {
		return "", fmt.Errorf("issuer field not found in Entra ID OpenID configuration")
	}
	return config.Issuer, nil
}

// UpdateOIDCConnector builds a Dex UpdateConnectorReq for an OIDC connector.
// It prepares the 'NewConfig' based on the provided parameters and subtype logic.
// Note: Dex's UpdateConnector gRPC only allows updating the config bytes.
func UpdateOIDCConnector(params UpdateConnectorRequest) (*dexapi.UpdateConnectorReq, error) {
	var oidcConfig OIDCConfig // Changed name to avoid conflict

	// Read common redirect URI(s) from environment - needed for config regeneration
	dexCallbackURL := os.Getenv("DEX_CALLBACK_URL")
	if dexCallbackURL == "" {
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable is required for updating OIDC connector")
	}
	redirectURIs := strings.Split(dexCallbackURL, ",")
	validRedirectURIs := []string{}
	for _, uri := range redirectURIs {
		trimmed := strings.TrimSpace(uri)
		if trimmed != "" {
			validRedirectURIs = append(validRedirectURIs, trimmed)
		}
	}
	if len(validRedirectURIs) == 0 {
		return nil, fmt.Errorf("DEX_CALLBACK_URL environment variable contains no valid URIs")
	}
	redirectURI := validRedirectURIs[0] // Use first as singular

	// Build the new OIDCConfig based on subtype and provided parameters
	switch strings.ToLower(params.ConnectorType) {
	case "oidc":
		switch strings.ToLower(params.ConnectorSubType) {
		case "general":
			if params.Issuer == "" {
				return nil, fmt.Errorf("issuer URL is required for 'general' OIDC connector subtype update")
			}
			oidcConfig = OIDCConfig{Issuer: params.Issuer, ClientID: params.ClientID, ClientSecret: params.ClientSecret, RedirectURIs: validRedirectURIs, RedirectURI: redirectURI, InsecureEnableGroups: true, InsecureSkipEmailVerified: true}
		case "entraid":
			if params.TenantID == "" {
				return nil, fmt.Errorf("tenant_id is required for 'entraid' OIDC connector subtype update")
			}
			if params.Issuer == "" {
				issuer, err := fetchEntraIDIssuer(params.TenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch issuer for entraid tenant %s update: %w", params.TenantID, err)
				}
				params.Issuer = issuer
			}
			oidcConfig = OIDCConfig{Issuer: params.Issuer, TenantID: params.TenantID, ClientID: params.ClientID, ClientSecret: params.ClientSecret, RedirectURIs: validRedirectURIs, RedirectURI: redirectURI, InsecureEnableGroups: true, InsecureSkipEmailVerified: true}
		case "google-workspace":
			oidcConfig = OIDCConfig{Issuer: "https://accounts.google.com", ClientID: params.ClientID, ClientSecret: params.ClientSecret, RedirectURIs: validRedirectURIs, RedirectURI: redirectURI, InsecureEnableGroups: true, InsecureSkipEmailVerified: true}
		default:
			return nil, fmt.Errorf("unsupported connector_sub_type for update: %s", params.ConnectorSubType)
		}
	default:
		return nil, fmt.Errorf("unsupported connector_type for update: %s", params.ConnectorType)
	}

	// Serialize the new OIDCConfig to JSON.
	configBytes, err := json.Marshal(oidcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new OIDC config for update: %w", err)
	}

	// Construct the UpdateConnectorReq message.
	// Note: Only Id and NewConfig are used by Dex gRPC.
	req := &dexapi.UpdateConnectorReq{
		Id:        params.ID, // Use the ID field from UpdateConnectorRequest
		NewConfig: configBytes,
	}

	return req, nil
}

// IsSupportedSubType checks if a given subtype identifier is valid for a given connector type.
// Uses the SupportedConnectors map for validation. Case-insensitive comparison.
func IsSupportedSubType(connectorType, subType string) bool {
	subTypes, exists := SupportedConnectors[strings.ToLower(connectorType)]
	if !exists {
		return false
	} // Connector type itself is not known/supported
	lowerSubType := strings.ToLower(subType)
	for _, st := range subTypes {
		if st == lowerSubType {
			return true
		}
	}
	return false // Subtype not found for the given connector type
}

// GetConnectorCreator returns the ConnectorCreator function for the specified connector type.
// Returns nil if the connector type is not supported.
func GetConnectorCreator(connectorType string) ConnectorCreator {
	return connectorCreators[strings.ToLower(connectorType)]
}

// GetSupportedConnectors returns a slice of supported subtype identifiers for a given connector type.
// Returns nil if the connector type is not supported.
func GetSupportedConnectors(connectorType string) []string {
	return SupportedConnectors[strings.ToLower(connectorType)]
}

// RestartDexPod attempts to restart the Dex deployment pod within the current Kubernetes cluster
// by deleting the pod associated with the 'dex' label/name within the configured namespace.
// This relies on a Kubernetes Deployment/StatefulSet to automatically recreate the pod.
// WARNING: This is a potentially disruptive operation and tightly couples the auth service
// to the Kubernetes environment and Dex's deployment specifics. Use with caution.
// It requires the service account running this auth service to have RBAC permissions
// to list and delete pods in the target namespace (read from NAMESPACE env var).
func RestartDexPod() error {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return fmt.Errorf("NAMESPACE environment variable not set, cannot determine where Dex pod lives")
	}

	kuberConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in-cluster Kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(kuberConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	// Consider using a label selector like "app=dex" for better targeting
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
	}

	deleted := false
	var deleteErrors []string
	for _, pod := range pods.Items {
		// Adjust this condition based on how Dex pods are named/labeled in your cluster
		if strings.Contains(pod.Name, "dex") {
			fmt.Printf("Attempting to delete Dex pod '%s' in namespace '%s' to reload configuration...\n", pod.Name, namespace)
			deletePolicy := metav1.DeletePropagationBackground
			err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
			if err != nil {
				fmt.Printf("WARNING: Failed to delete pod %s: %v\n", pod.Name, err)
				deleteErrors = append(deleteErrors, fmt.Sprintf("pod %s: %v", pod.Name, err))
				// Continue to try deleting other matching pods if any
			} else {
				fmt.Printf("Pod %s deleted successfully request sent.\n", pod.Name)
				deleted = true
				// If deleting one is sufficient (e.g., Deployment handles rollout), break here.
				// If all matching pods must be deleted, remove the break.
				break
			}
		}
	}

	if !deleted && len(deleteErrors) == 0 {
		fmt.Printf("WARNING: No pod containing 'dex' found in namespace '%s' to delete.\n", namespace)
		// Return nil or error depending on whether restart is critical
		// return fmt.Errorf("no Dex pod found in namespace '%s' to restart", namespace)
	}

	// If there were errors deleting pods, return a combined error
	if len(deleteErrors) > 0 {
		return fmt.Errorf("failed to delete one or more Dex pods: %s", strings.Join(deleteErrors, "; "))
	}

	return nil
}
