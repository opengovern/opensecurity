package types

import (
	"strconv"
)

type ComplianceJobReportResourceView struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	Integrations      map[string]AuditIntegrationResult `json:"integrations"`
	ComplianceSummary map[ComplianceStatus]uint64       `json:"compliance_summary"`
	JobSummary        JobSummary                        `json:"job_summary"`
}

func (r ComplianceJobReportResourceView) KeysAndIndex() ([]string, string) {
	return []string{
		strconv.Itoa(int(r.JobSummary.JobID)),
	}, ComplianceJobReportResourceViewIndex
}

type AuditControlFinding struct {
	Severity  ComplianceResultSeverity `json:"severity"`
	ControlID string                   `json:"control_id"`
	Reason    string                   `json:"reason"`
}

type AuditResourceResult struct {
	ResourceName    string                                     `json:"resource_name"`
	ResourceSummary map[ComplianceStatus]uint64                `json:"control_summary"`
	Results         map[ComplianceStatus][]AuditControlFinding `json:"results"`
}

type AuditResourceTypesResult struct {
	Resources map[string]AuditResourceResult `json:"resources"`
}

type AuditIntegrationResult struct {
	ResourceTypes map[string]AuditResourceTypesResult `json:"resource_types"`
}
