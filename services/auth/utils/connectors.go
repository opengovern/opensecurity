package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	dexapi "github.com/dexidp/dex/api/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
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
	ID               string `json:"id,omitempty"`   // Optional
	Name             string `json:"name,omitempty"` // Optional
}

type OIDCConfig struct {
	Issuer       string `json:"issuer,omitempty"`
	TenantID     string `json:"tenantID,omitempty"` // Added TenantID for entraid sub-type
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Name 			string `json:"name,omitempty"`
	RedirectURIs		[]string `json:"redirect_uris,omitempty"`
	RedirectURI 		string `json:"redirectURI,omitempty"`
	InsecureEnableGroups bool     `json:"insecureEnableGroups"`
	InsecureSkipEmailVerified bool `json:"insecureSkipEmailVerified"`



}

type ConnectorCreator func( params CreateConnectorRequest) (*dexapi.CreateConnectorReq, error)

var  connectorCreators = map[string]ConnectorCreator{
	"oidc": CreateOIDCConnector,
	// Future connector types can be added here, e.g., "saml": (*DexClient).CreateSAMLConnector
}
var SupportedConnectors = map[string][]string{
	"oidc": {"general", "google-workspace", "entraid"},
	// Add more connector types and their sub-types here as needed.
}
var SupportedConnectorsNames = map[string][]string{
	"oidc": {"General OIDC", "Google Workspaces", "AzureAD/EntraID"},

}

func  CreateOIDCConnector(params CreateConnectorRequest) (*dexapi.CreateConnectorReq, error) {


	var oidcConfig OIDCConfig
	var connectorID, connectorName string
	connectorID = params.ID
	connectorName = params.Name
	switch params.ConnectorSubType {
	case "general":
		// Required: issuer, clientID, clientSecret
		oidcConfig = OIDCConfig{
			Issuer:       params.Issuer,
			ClientID:     params.ClientID,
			ClientSecret: params.ClientSecret,
			RedirectURIs: strings.Split(os.Getenv("DEX_CALLBACK_URL"),","),
			RedirectURI: strings.Split(os.Getenv("DEX_CALLBACK_URL"),",")[0],
			InsecureEnableGroups: true,
			InsecureSkipEmailVerified: true,


		}
		

		

	case "entraid":
		// Required: tenantID, clientID, clientSecret
		if   params.TenantID != "" && params.Issuer == "" {
				issuer, err := fetchEntraIDIssuer(params.TenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch issuer for entraid: %w", err)
				}
				params.Issuer = issuer
			}
		oidcConfig = OIDCConfig{
			Issuer:       params.Issuer,
			TenantID:     params.TenantID,
			ClientID:     params.ClientID,
			ClientSecret: params.ClientSecret,
			RedirectURIs: strings.Split(os.Getenv("DEX_CALLBACK_URL"),","),
			RedirectURI: strings.Split(os.Getenv("DEX_CALLBACK_URL"),",")[0],
			InsecureEnableGroups: true,
			InsecureSkipEmailVerified: true,


		}
		

		

	case "google-workspace":
		// Required: clientID, clientSecret
		oidcConfig = OIDCConfig{
			ClientID:     params.ClientID,
			ClientSecret: params.ClientSecret,
			Issuer:       "https://accounts.google.com",
			RedirectURIs: strings.Split(os.Getenv("DEX_CALLBACK_URL"),","),
			RedirectURI: strings.Split(os.Getenv("DEX_CALLBACK_URL"),",")[0],
			InsecureEnableGroups: true,
			InsecureSkipEmailVerified: true,


		}
	

		
	default:
		return nil, fmt.Errorf("unsupported connector_sub_type: %s", params.ConnectorSubType)
	}

	// Serialize the OIDCConfig to JSON.
	configBytes, err := json.Marshal(oidcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OIDC config: %w", err)
	}

	// Construct the Connector message.
	connector := &dexapi.Connector{
		Id:     connectorID,
		Type:   "oidc",
		Name:   connectorName,
		Config: configBytes,
	}

	// Create the CreateConnectorReq message.
	req := &dexapi.CreateConnectorReq{
		Connector: connector,
	}

	

	// Execute the CreateConnector RPC.
	

	return req, nil
}
func fetchEntraIDIssuer(tenantID string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", tenantID)
	resp, err := http.Get(url)
	if err != nil {
		
		return "", fmt.Errorf("failed to fetch OpenID configuration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
	
		return "", fmt.Errorf("unexpected status code %d when fetching OpenID configuration", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenID configuration response: %w", err)
	}

	var config struct {
		Issuer string `json:"issuer"`
	}
	if err := json.Unmarshal(body, &config); err != nil {
		return "", fmt.Errorf("failed to parse OpenID configuration: %w", err)
	}

	if config.Issuer == "" {
		return "", fmt.Errorf("issuer not found in OpenID configuration")
	}


	return config.Issuer, nil
}


func UpdateOIDCConnector(params UpdateConnectorRequest) (*dexapi.UpdateConnectorReq, error) {
	var newOIDCConfig OIDCConfig

	switch params.ConnectorType {
	case "oidc":
		switch params.ConnectorSubType {
		case "google-workspace", "entraid", "general":
			if params.ConnectorSubType == "entraid" && params.TenantID != "" && params.Issuer == "" {
				issuer, err := fetchEntraIDIssuer(params.TenantID)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch issuer for entraid: %w", err)
				}
				params.Issuer = issuer
			}
			if params.ConnectorSubType == "google-workspace" {
				params.Issuer = "https://accounts.google.com"
			}
			
				newOIDCConfig = OIDCConfig{
				Issuer:       params.Issuer,
				TenantID:     params.TenantID, // Ensure TenantID is set for entraid
				ClientID:     params.ClientID,
				ClientSecret: params.ClientSecret,
				}
			
			
			
		default:
			return nil, fmt.Errorf("unsupported connector_sub_type: %s", params.ConnectorSubType)
		}
	default:
		return nil, fmt.Errorf("unsupported connector_type: %s", params.ConnectorType)
	}
	configBytes, err := json.Marshal(newOIDCConfig)
	if err != nil {
	
		return nil, fmt.Errorf("failed to marshal new OIDC config: %w", err)
	}
	
	req := &dexapi.UpdateConnectorReq{
		Id:        params.ID,
		NewConfig: configBytes,
	}
	

	
	return req, nil

	
}
func IsSupportedSubType(connectorType, subType string) bool {
	subTypes, exists := SupportedConnectors[connectorType]
	if !exists {
		return false
	}
	for _, st := range subTypes {
		if strings.ToLower(st) == subType {
			return true
		}
	}
	return false
}

func GetConnectorCreator(connectorType string) ConnectorCreator {
	return connectorCreators[connectorType]
}
func GetSupportedConnectors(connectorType string) ([]string ) {
	return SupportedConnectors[connectorType]
}



func RestartDexPod() error {
	// Restart Dex pod by deleting it.
	// The pod will be recreated by the deployment.
	// This is a workaround to reload the connectors.
	kuberConfig, err := rest.InClusterConfig()
	if err != nil {
		
		return  fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(kuberConfig)
	if err != nil {
		
		return  fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}
	pods, err := clientset.CoreV1().Pods(os.Getenv("NAMESPACE")).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "dex") {
			err := clientset.CoreV1().Pods(os.Getenv("NAMESPACE")).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("failed to delete pod %s: %v", pod.Name, err)
			}
			fmt.Printf("Pod %s deleted successfully\n", pod.Name)
		}
	}
	

	return nil


}