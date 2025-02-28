package integration

import (
	"github.com/opengovern/opensecurity/services/integration/models"
)

var defaultIntegrationGroups = []models.IntegrationGroup{
	{
		Name:  "active",
		Query: `SELECT integration_id FROM platform_integrations WHERE state = 'ACTIVE'`,
	},
	{
		Name:  "inactive",
		Query: `SELECT integration_id FROM platform_integrations WHERE state = 'INACTIVE'`,
	},
	{
		Name:  "archived",
		Query: `SELECT integration_id FROM platform_integrations WHERE state = 'ARCHIVED'`,
	},
}
