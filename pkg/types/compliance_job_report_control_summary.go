package types

import (
	"strconv"
)

type ComplianceJobReportControlSummary struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	Controls          map[string]*ControlSummary  `json:"controls"`
	ControlScore      *ControlScore               `json:"control_score"`
	ComplianceSummary map[ComplianceStatus]uint64 `json:"compliance_summary"`
	JobSummary        JobSummary                  `json:"job_summary"`
}

func (r ComplianceJobReportControlSummary) KeysAndIndex() ([]string, string) {
	return []string{
		strconv.Itoa(int(r.JobSummary.JobID)),
	}, ComplianceJobReportControlSummaryIndex
}

type ControlScore struct {
	TotalControls  int64 `json:"total_controls"`
	FailedControls int64 `json:"failed_controls"`
}

type ControlSummary struct {
	Severity ComplianceResultSeverity `json:"severity"`
	Alarms   int64                    `json:"alarms"`
	Oks      int64                    `json:"oks"`
}
