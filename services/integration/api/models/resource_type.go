package models

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

type ResourceTypeConfiguration struct {
	Name            string           `json:"name"`
	IntegrationType integration.Type `json:"integration_type"`
	Description     string           `json:"description"`
	Params          []Param          `json:"params"`
	Table 			string			 `json:"table"`
}

func ApiResourceTypeConfiguration(configuration interfaces.ResourceTypeConfiguration) ResourceTypeConfiguration {
	params := make([]Param, 0, len(configuration.Params))
	for _, param := range configuration.Params {
		params = append(params, ApiParam(param))
	}

	return ResourceTypeConfiguration{
		Name:            configuration.Name,
		IntegrationType: configuration.IntegrationType,
		Description:     configuration.Description,
		Params:          params,
		Table:			 configuration.Table,
	}
}

type Param struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Default     *string `json:"default"`
}

func ApiParam(param interfaces.Param) Param {
	return Param{
		Name:        param.Name,
		Description: param.Description,
		Required:    param.Required,
		Default:     param.Default,
	}
}

type ListIntegrationTypeResourceTypesResponse struct {
	ResourceTypes []ResourceTypeConfiguration `json:"integration_types"`
	TotalCount    int                         `json:"total_count"`
}
