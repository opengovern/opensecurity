package models

type IntegrationGroup struct {
	Name           string        `json:"name" example:"UltraSightApplication"`
	Query          string        `json:"query" example:"SELECT og_id FROM platform_integrations WHERE labels->'application' IS NOT NULL AND labels->'application' @> '\"UltraSight\"'"`
	IntegrationIds []string      `json:"integration_ids,omitempty" example:"[\"1e8ac3bf-c268-4a87-9374-ce04cc40a596\"]"`
	Integrations   []Integration `json:"integrations,omitempty"`
}
