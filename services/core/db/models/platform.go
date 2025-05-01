package models

import (
	"encoding/json"
	"strconv"

	"strings"

	metadataErrors "github.com/opengovern/opensecurity/services/core/errors"

	"errors"

	"github.com/opengovern/og-util/pkg/api"
)

const (
	ConfigMetadataTypeString ConfigMetadataType = "string"
	ConfigMetadataTypeInt    ConfigMetadataType = "int"
	ConfigMetadataTypeBool   ConfigMetadataType = "bool"
	ConfigMetadataTypeJSON   ConfigMetadataType = "json"
)
const (
	MetadataKeyInstallID                    MetadataKey = "platform_install_id"
	MetadataKeyAppVersion                   MetadataKey = "platform_app_version"
	MetadataKeyCreationTime                 MetadataKey = "platform_creation_time"
	MetadataKeyMaxAPIKeys                   MetadataKey = "platform_max_keys"
	MetadataKeyAssetDiscoveryJobInterval    MetadataKey = "platform_discovery_interval"
	MetadataKeyComplianceEnabled            MetadataKey = "platform_compliance_enabled"
	MetadataKeyDefaultComplianceJobInterval MetadataKey = "platform_compliance_interval"
	MetadataKeyDomain                       MetadataKey = "platform_domain"
	MetadataKeyHTTPSEnabled                 MetadataKey = "platform_https_enabled"
	MetadataKeyUserLimit                    MetadataKey = "platform_user_limit"
	MetadataKeyUsageReportingInterval       MetadataKey = "platform_usage_reporting"
	MetadataKeyDataRetention                MetadataKey = "platformdata_retention_duration"
	MetadataKeyDateTimeFormat               MetadataKey = "platform_date_time_format"
	MetadataKeyPlatformConfigurationGitURL  MetadataKey = "platform_git_url"
)

var MetadataKeys = []MetadataKey{
	MetadataKeyInstallID,
	MetadataKeyAppVersion,
	MetadataKeyCreationTime,
	MetadataKeyMaxAPIKeys,
	MetadataKeyAssetDiscoveryJobInterval,
	MetadataKeyComplianceEnabled,
	MetadataKeyDefaultComplianceJobInterval,
	MetadataKeyDomain,
	MetadataKeyHTTPSEnabled,
	MetadataKeyUserLimit,
	MetadataKeyUsageReportingInterval,
	MetadataKeyDataRetention,
	MetadataKeyDateTimeFormat,
	MetadataKeyPlatformConfigurationGitURL,
}

func (k MetadataKey) String() string {
	return string(k)
}

func (k MetadataKey) GetConfigMetadataType() ConfigMetadataType {
	switch k {
	case MetadataKeyInstallID:
		return ConfigMetadataTypeString
	case MetadataKeyAppVersion:
		return ConfigMetadataTypeString
	case MetadataKeyCreationTime:
		return ConfigMetadataTypeString
	case MetadataKeyMaxAPIKeys:
		return ConfigMetadataTypeInt
	case MetadataKeyComplianceEnabled:
		return ConfigMetadataTypeBool
	case MetadataKeyAssetDiscoveryJobInterval:
		return ConfigMetadataTypeInt
	case MetadataKeyDefaultComplianceJobInterval:
		return ConfigMetadataTypeInt
	case MetadataKeyDomain:
		return ConfigMetadataTypeString
	case MetadataKeyHTTPSEnabled:
		return ConfigMetadataTypeBool
	case MetadataKeyUserLimit:
		return ConfigMetadataTypeInt
	case MetadataKeyUsageReportingInterval:
		return ConfigMetadataTypeInt
	case MetadataKeyDataRetention:
		return ConfigMetadataTypeInt
	case MetadataKeyDateTimeFormat:
		return ConfigMetadataTypeString
	case MetadataKeyPlatformConfigurationGitURL:
		return ConfigMetadataTypeString
	}
	return ""
}

func (k MetadataKey) GetMinAuthRole() api.Role {
	switch k {
	case MetadataKeyInstallID:
		return api.ViewerRole
	case MetadataKeyAppVersion:
		return api.ViewerRole
	case MetadataKeyCreationTime:
		return api.ViewerRole
	case MetadataKeyMaxAPIKeys:
		return api.ViewerRole
	case MetadataKeyComplianceEnabled:
		return api.ViewerRole
	case MetadataKeyDefaultComplianceJobInterval:
		return api.ViewerRole
	case MetadataKeyDomain:
		return api.ViewerRole
	case MetadataKeyHTTPSEnabled:
		return api.ViewerRole
	case MetadataKeyUserLimit:
		return api.ViewerRole
	case MetadataKeyUsageReportingInterval:
		return api.ViewerRole
	case MetadataKeyDataRetention:
		return api.ViewerRole
	case MetadataKeyDateTimeFormat:
		return api.ViewerRole
	case MetadataKeyPlatformConfigurationGitURL:
		return api.ViewerRole
	case MetadataKeyAssetDiscoveryJobInterval:
		return api.ViewerRole
	}
	return ""
}

func ParseMetadataKey(key string) (MetadataKey, error) {
	lowerKey := strings.ToLower(key)
	for _, k := range MetadataKeys {
		if lowerKey == strings.ToLower(k.String()) {
			return k, nil
		}
	}
	return "", metadataErrors.ErrMetadataKeyNotFound
}

func (t ConfigMetadataType) SerializeValue(value any) (string, error) {
	switch t {
	case ConfigMetadataTypeString:
		valueStr, ok := value.(string)
		if !ok {
			return "", metadataErrors.ErrorMetadataValueTypeMismatch
		}
		return valueStr, nil
	case ConfigMetadataTypeInt:
		switch value.(type) {
		case int:
			return strconv.Itoa(value.(int)), nil
		case string:
			valueM, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				return "", err
			}
			return strconv.Itoa(int(valueM)), nil
		default:
			return "", metadataErrors.ErrorMetadataValueTypeMismatch
		}
	case ConfigMetadataTypeBool:
		switch value.(type) {
		case bool:
			return strconv.FormatBool(value.(bool)), nil
		case string:
			valueM, err := strconv.ParseBool(value.(string))
			if err != nil {
				return "", err
			}
			return strconv.FormatBool(valueM), nil
		default:
			return "", metadataErrors.ErrorMetadataValueTypeMismatch
		}
	case ConfigMetadataTypeJSON:
		valueJson, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(valueJson), nil
	}
	return "", metadataErrors.ErrorMetadataValueTypeMismatch
}

func (t ConfigMetadataType) DeserializeValue(value string) (any, error) {
	switch t {
	case ConfigMetadataTypeString:
		return value, nil
	case ConfigMetadataTypeInt:
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return valueInt, nil
	case ConfigMetadataTypeBool:
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return valueBool, nil
	case ConfigMetadataTypeJSON:
		var valueJson any
		err := json.Unmarshal([]byte(value), &valueJson)
		if err != nil {
			return nil, err
		}
		return valueJson, nil
	}
	return nil, metadataErrors.ErrorMetadataValueTypeMismatch
}

func (c *ConfigMetadata) ParseToType() (IConfigMetadata, error) {
	value, err := c.Type.DeserializeValue(c.Value)
	if err != nil {
		return nil, err
	}
	switch c.Type {
	case ConfigMetadataTypeString:
		return &StringConfigMetadata{ConfigMetadata: *c}, nil
	case ConfigMetadataTypeInt:
		return &IntConfigMetadata{ConfigMetadata: *c, Value: value.(int)}, nil
	case ConfigMetadataTypeBool:
		return &BoolConfigMetadata{ConfigMetadata: *c, Value: value.(bool)}, nil
	case ConfigMetadataTypeJSON:
		return &JSONConfigMetadata{ConfigMetadata: *c, Value: value}, nil
	}

	return nil, metadataErrors.ErrorConfigMetadataTypeNotSupported
}

var ErrIncompatibleType = errors.New("given config metadata interface does not have the required type")

func HasType(cfg IConfigMetadata, typ ConfigMetadataType) error {
	if cfg.GetType() != typ {
		return ErrIncompatibleType
	}

	return nil
}

type IConfigMetadata interface {
	GetKey() MetadataKey
	GetType() ConfigMetadataType
	GetValue() any
	GetCore() ConfigMetadata
}

func (c *StringConfigMetadata) GetKey() MetadataKey {
	return c.Key
}

func (c *StringConfigMetadata) GetType() ConfigMetadataType {
	return ConfigMetadataTypeString
}

func (c *StringConfigMetadata) GetValue() any {
	return c.Value
}

func (c *StringConfigMetadata) GetCore() ConfigMetadata {
	return c.ConfigMetadata
}

func (c *IntConfigMetadata) GetKey() MetadataKey {
	return c.Key
}

func (c *IntConfigMetadata) GetType() ConfigMetadataType {
	return ConfigMetadataTypeInt
}

func (c *IntConfigMetadata) GetValue() any {
	return c.Value
}

func (c *IntConfigMetadata) GetCore() ConfigMetadata {
	return c.ConfigMetadata
}

func (c *BoolConfigMetadata) GetKey() MetadataKey {
	return c.Key
}

func (c *BoolConfigMetadata) GetType() ConfigMetadataType {
	return ConfigMetadataTypeBool
}

func (c *BoolConfigMetadata) GetValue() any {
	return c.Value
}

func (c *BoolConfigMetadata) GetCore() ConfigMetadata {
	return c.ConfigMetadata
}

func (c *JSONConfigMetadata) GetKey() MetadataKey {
	return c.Key
}

func (c *JSONConfigMetadata) GetType() ConfigMetadataType {
	return ConfigMetadataTypeJSON
}

func (c *JSONConfigMetadata) GetValue() any {
	return c.Value
}

func (c *JSONConfigMetadata) GetCore() ConfigMetadata {
	return c.ConfigMetadata
}

func (qp PolicyParameterValues) GetKey() string {
	return qp.Key
}

func (qp PolicyParameterValues) GetValue() string {
	return qp.Value
}
