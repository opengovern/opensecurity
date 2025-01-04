package api

import (
	"time"

	"github.com/opengovern/opencomply/jobs/post-install-job/db/model"
	authApi "github.com/opengovern/opencomply/services/auth/api"
	api "github.com/opengovern/opencomply/services/integration/api/models"
)

type SetConfigMetadataRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type DexConnectorInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type About struct {
	InstallID             string                       `json:"install_id"`
	DexConnectors         []DexConnectorInfo           `json:"dex_connectors"`
	AppVersion            string                       `json:"app_version"`
	WorkspaceCreationTime time.Time                    `json:"workspace_creation_time"`
	Users                 []authApi.GetUsersResponse   `json:"users"`
	PrimaryDomainURL      string                       `json:"primary_domain_url"`
	APIKeys               []authApi.APIKeyResponse     `json:"api_keys"`
	Integrations          map[string][]api.Integration `json:"integrations"`
	SampleData            bool                         `json:"sample_data"`
	TotalSpendGoverned    float64                      `json:"total_spend_governed"`
}

type GetMigrationStatusResponse struct {
	Status     string                   `json:"status"`
	JobsStatus map[string]model.JobInfo `json:"jobs_status"`
	Summary    struct {
		TotalJobs          int     `json:"total_jobs"`
		CompletedJobs      int     `json:"completed_jobs"`
		ProgressPercentage float64 `json:"progress_percentage"`
	}
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	
}

type GetSampleSyncStatusResponse struct {
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
}
