package api

import (
	"github.com/opengovern/opensecurity/pkg/types"
	"time"
)

type AuditSummary struct {
	Controls     map[string]types.AuditControlResult     `json:"controls"`
	Integrations map[string]types.AuditIntegrationResult `json:"integrations"`
	AuditSummary map[types.ComplianceStatus]uint64       `json:"audit_summary"`
	JobSummary   types.JobSummary                        `json:"job_summary"`
}

type GetJobReportSummaryJobDetails struct {
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Framework struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"framework"`
	IntegrationIDs []string                             `json:"integration_ids"`
	JobScore       JobScore                             `json:"job_score"`
	Results        GetJobReportSummaryJobDetailsResults `json:"results"`
}

type GetJobReportSummaryJobDetailsResults struct {
	Alarm uint64 `json:"alarm"`
	Ok    uint64 `json:"ok"`
}

type JobScore struct {
	TotalControls  int64               `json:"total_controls"`
	FailedControls int64               `json:"failed_controls"`
	ControlView    JobScoreControlView `json:"control_view"`
}

type JobScoreControlView struct {
	BySeverity map[string]*JobScoreControlViewBySeverityScore `json:"by_severity"`
}

type JobScoreControlViewBySeverityScore struct {
	TotalControls  uint64 `json:"total_controls"`
	FailedControls uint64 `json:"failed_controls"`
}

type GetJobReportSummaryResponse struct {
	JobID         uint                          `json:"job_id"`
	WithIncidents bool                          `json:"with_incidents"`
	JobDetails    GetJobReportSummaryJobDetails `json:"job_details"`
}
