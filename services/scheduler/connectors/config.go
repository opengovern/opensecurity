package connectors

import "encoding/json"

type AWSAccountConfig struct {
	AccountID            string   `json:"accountId"`
	Regions              []string `json:"regions"`
	SecretKey            string   `json:"secretKey"`
	AccessKey            string   `json:"accessKey"`
	SessionToken         string   `json:"sessionToken"`
	AssumeRoleName       string   `json:"assumeRoleName"`
	AssumeAdminRoleName  string   `json:"assumeAdminRoleName"`
	AssumeRolePolicyName string   `json:"assumeRolePolicyName"`
	ExternalID           *string  `json:"externalID,omitempty"`
}

func (asc AWSAccountConfig) ToMap() map[string]any {
	jsonCnf, err := json.Marshal(asc)
	if err != nil {
		return nil
	}
	res := make(map[string]any)
	err = json.Unmarshal(jsonCnf, &res)
	if err != nil {
		return nil
	}
	return res
}

func AWSAccountConfigFromMap(m map[string]any) (AWSAccountConfig, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return AWSAccountConfig{}, err
	}

	var c AWSAccountConfig
	err = json.Unmarshal(mj, &c)
	if err != nil {
		return AWSAccountConfig{}, err
	}

	return c, nil
}

type AzureSubscriptionConfig struct {
	SubscriptionID  string `json:"subscriptionId"`
	TenantID        string `json:"tenantId"`
	ObjectID        string `json:"objectId"`
	SecretID        string `json:"secretId"`
	ClientID        string `json:"clientId"`
	ClientSecret    string `json:"clientSecret"`
	CertificatePath string `json:"certificatePath"`
	CertificatePass string `json:"certificatePass"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}

func (asc AzureSubscriptionConfig) ToMap() map[string]any {
	jsonCnf, err := json.Marshal(asc)
	if err != nil {
		return nil
	}
	res := make(map[string]any)
	err = json.Unmarshal(jsonCnf, &res)
	if err != nil {
		return nil
	}
	return res
}

func AzureSubscriptionConfigFromMap(m map[string]interface{}) (AzureSubscriptionConfig, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return AzureSubscriptionConfig{}, err
	}

	var c AzureSubscriptionConfig
	err = json.Unmarshal(mj, &c)
	if err != nil {
		return AzureSubscriptionConfig{}, err
	}

	return c, nil
}
