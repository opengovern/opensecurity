package interfaces

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/api/models"
)

type ResourceTypeConfiguration struct {
	Name            string           `json:"name"`
	IntegrationType integration.Type `json:"integration_type"`
	Description     string           `json:"description"`
	Params          []Param          `json:"params"`
}

type Param struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Default     *string `json:"default"`
}

func (c *ResourceTypeConfiguration) ToAPI() models.ResourceTypeConfiguration {
	var params []models.Param
	for _, param := range c.Params {
		params = append(params, param.ToAPI())
	}
	return models.ResourceTypeConfiguration{
		Name:            c.Name,
		IntegrationType: c.IntegrationType,
		Description:     c.Description,
		Params:          params,
	}
}

func (p *ResourceTypeConfiguration) IsEmpty() bool {
	return p.Name == "" && p.Description == "" && len(p.Params) == 0 && p.IntegrationType == ""
}

func (p *Param) ToAPI() models.Param {
	return models.Param{
		Name:        p.Name,
		Description: p.Description,
		Required:    p.Required,
		Default:     p.Default,
	}
}
