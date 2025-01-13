//go:build go1.18
// +build go1.18

package healthcheck

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Global constants for required permissions and required licenses
var requiredPermissions = []string{"Directory.Read.All", "User.Read.All"}
var requiredLicenses = []string{"EMSPREMIUM", "AAD_PREMIUM_P1", "AAD_PREMIUM_P2"}

// Config represents the JSON input configuration
type Config struct {
	ObjectID     string `json:"objectId,omitempty"`
	TenantID     string `json:"tenantId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"`
	CertPath     string `json:"certPath,omitempty"`
	CertContent  string `json:"certContent,omitempty"`
	CertPassword string `json:"certPassword,omitempty"`
}

// HealthCheckDetails contains details about the health check
type HealthCheckDetails struct {
	LicensesAssigned    []string `json:"licensesAssigned"`
	LicensesRequired    []string `json:"licensesRequired"`
	PermissionsRequired []string `json:"permissionsRequired"`
	PermissionsAssigned []string `json:"permissionsAssigned"`
}

// TenantInfo represents the comprehensive information of a tenant (directory)
type TenantInfo struct {
	Name               string             `json:"name"`
	TenantID           string             `json:"tenantId"`
	PrimaryDomain      string             `json:"primaryDomain"`
	Healthy            bool               `json:"healthy"`
	HealthCheckDetails HealthCheckDetails `json:"healthCheckDetails"`
}

func EntraidIntegrationHealthcheck(config Config) (bool, error) {
	ctx := context.Background()

	// Validate required fields
	if config.TenantID == "" || config.ClientID == "" {
		return false, fmt.Errorf("tenantId and clientId are required in the configuration.")
	}

	// Authenticate and get credentials
	credential, err := authenticate(ctx, &config)
	if err != nil {
		return false, fmt.Errorf("Authentication failed: %v", err)
	}

	// Get tenant info
	tenantInfo, err := getTenantInfo(ctx, credential, config)
	if err != nil {
		return false, fmt.Errorf("Failed to get tenant info: %v", err)
	}

	return tenantInfo.Healthy, nil
}

// authenticate creates an azcore.TokenCredential based on the provided configuration.
func authenticate(ctx context.Context, config *Config) (azcore.TokenCredential, error) {
	var cred azcore.TokenCredential
	var err error

	if config.CertPath != "" || config.CertContent != "" {
		// Use certificate-based authentication

		var certData []byte
		if config.CertPath != "" {
			// Read the certificate file
			certData, err = ioutil.ReadFile(config.CertPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read certificate file: %v", err)
			}
		} else {
			// Use certificate content provided directly
			certData = []byte(config.CertContent)
		}

		var password []byte
		if config.CertPassword != "" {
			password = []byte(config.CertPassword)
		}

		// Parse the certificate using azidentity.ParseCertificates
		certs, key, err := azidentity.ParseCertificates(certData, password)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}

		// Create the ClientCertificateCredential
		cred, err = azidentity.NewClientCertificateCredential(
			config.TenantID,
			config.ClientID,
			certs,
			key,
			nil, // Additional options can be set here if needed
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create certificate credential: %v", err)
		}
	} else if config.ClientSecret != "" {
		// Use client secret authentication
		cred, err = azidentity.NewClientSecretCredential(
			config.TenantID,
			config.ClientID,
			config.ClientSecret,
			nil, // Additional options can be set here if needed
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create client secret credential: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no valid authentication method found. Provide clientSecret or certPath/certContent in the configuration.")
	}

	return cred, nil
}

// parseJWT parses a JWT token and returns the claims.
func parseJWT(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token: not enough parts")
	}
	payload := parts[1]

	// Add padding if necessary
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	case 0:
		// No padding needed
	default:
		// Invalid base64 string
		return nil, fmt.Errorf("invalid base64 payload")
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode payload: %v", err)
	}

	var claims map[string]interface{}
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %v", err)
	}

	return claims, nil
}

// getTenantInfo retrieves tenant information, license details, and determines the health status.
func getTenantInfo(ctx context.Context, credential azcore.TokenCredential, config Config) (*TenantInfo, error) {
	// Initialize TenantInfo
	tenantInfo := &TenantInfo{}
	healthDetails := HealthCheckDetails{}

	// Get an access token for Microsoft Graph API
	token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	client := &http.Client{}

	// Step 1: Get organization info
	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/organization", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("failed to get organization info: %s", bodyString)
	}

	// Parse the response
	var result struct {
		Value []struct {
			DisplayName     string `json:"displayName"`
			ID              string `json:"id"`
			VerifiedDomains []struct {
				Capabilities string `json:"capabilities"`
				IsDefault    bool   `json:"isDefault"`
				Name         string `json:"name"`
			} `json:"verifiedDomains"`
		} `json:"value"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse organization response: %v", err)
	}

	if len(result.Value) == 0 {
		return nil, fmt.Errorf("no organization info found")
	}

	org := result.Value[0]

	tenantInfo.Name = org.DisplayName
	tenantInfo.TenantID = org.ID

	for _, domain := range org.VerifiedDomains {
		if domain.IsDefault {
			tenantInfo.PrimaryDomain = domain.Name
			break
		}
	}

	// Step 2: Fetch License Information
	licenses, err := getLicenseInfo(ctx, client, token.Token)
	if err != nil {
		log.Printf("Warning: Failed to get license info for tenant %s: %v", tenantInfo.TenantID, err)
		licenses = []string{"N/A"}
	}
	healthDetails.LicensesAssigned = licenses
	healthDetails.LicensesRequired = requiredLicenses

	// Step 3: Check if SPN has required permissions
	assignedPermissions, err := getAssignedPermissions(ctx, credential)
	if err != nil {
		return nil, fmt.Errorf("failed to check SPN permissions: %v", err)
	}
	healthDetails.PermissionsAssigned = assignedPermissions
	healthDetails.PermissionsRequired = requiredPermissions

	// Step 4: Determine Healthy Status
	// Healthy if primary domain is present, and both required permissions and licenses are assigned
	hasPrimaryDomain := tenantInfo.PrimaryDomain != ""
	hasAllPermissions := hasAllElements(requiredPermissions, assignedPermissions)
	hasAllLicenses := hasAnyElements(requiredLicenses, licenses)

	tenantInfo.Healthy = hasPrimaryDomain && hasAllPermissions && hasAllLicenses
	tenantInfo.HealthCheckDetails = healthDetails

	if !hasPrimaryDomain {
		return tenantInfo, fmt.Errorf("needs to have primary domain")
	}
	if !hasAllPermissions {
		return tenantInfo, fmt.Errorf("needs to have all permissions: %v", requiredPermissions)
	}
	if !hasAllLicenses {
		return tenantInfo, fmt.Errorf("needs to have all licenses: %v", requiredLicenses)
	}

	return tenantInfo, nil
}

// getAssignedPermissions retrieves the permissions assigned to the SPN.
func getAssignedPermissions(ctx context.Context, credential azcore.TokenCredential) ([]string, error) {
	// Get an access token for Microsoft Graph API
	token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	// Parse the token
	claims, err := parseJWT(token.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract roles claim
	rolesInterface, ok := claims["roles"]
	if !ok {
		return nil, fmt.Errorf("roles claim not found in token")
	}

	var assignedPermissions []string

	switch rolesValue := rolesInterface.(type) {
	case []interface{}:
		// roles is an array
		for _, role := range rolesValue {
			if roleStr, ok := role.(string); ok {
				assignedPermissions = append(assignedPermissions, roleStr)
			}
		}
	case string:
		// roles is a single string
		assignedPermissions = []string{rolesValue}
	default:
		return nil, fmt.Errorf("unexpected type for roles claim")
	}

	return assignedPermissions, nil
}

// getLicenseInfo fetches the license information of the tenant.
func getLicenseInfo(ctx context.Context, client *http.Client, token string) ([]string, error) {
	// Make a GET request to Microsoft Graph API to get subscribed SKUs
	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/subscribedSkus", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribed SKUs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("failed to get subscribed SKUs: %s", bodyString)
	}

	// Parse the response
	var skusResult struct {
		Value []struct {
			SkuPartNumber string `json:"skuPartNumber"`
		} `json:"value"`
	}

	err = json.NewDecoder(resp.Body).Decode(&skusResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscribed SKUs response: %v", err)
	}

	// Collect the SKU part numbers
	var licenses []string
	for _, sku := range skusResult.Value {
		licenses = append(licenses, sku.SkuPartNumber)
	}

	if len(licenses) == 0 {
		return []string{"N/A"}, nil
	}

	return licenses, nil
}

// hasAllElements checks if all elements of required are present in assigned.
func hasAllElements(required, assigned []string) bool {
	assignedSet := make(map[string]struct{})
	for _, item := range assigned {
		assignedSet[item] = struct{}{}
	}

	for _, req := range required {
		if _, found := assignedSet[req]; !found {
			return false
		}
	}

	return true
}

// hasAnyElements checks if any elements of required are present in assigned.
func hasAnyElements(required, assigned []string) bool {
	assignedSet := make(map[string]struct{})
	for _, item := range assigned {
		assignedSet[item] = struct{}{}
	}

	for _, req := range required {
		if _, found := assignedSet[req]; found {
			return true
		}
	}

	return false
}
