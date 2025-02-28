package types

import "strings"

type ComplianceStatus string

const (
	ComplianceStatusOK    ComplianceStatus = "ok"
	ComplianceStatusALARM ComplianceStatus = "alarm"
	ComplianceStatusINFO  ComplianceStatus = "info"
	ComplianceStatusSKIP  ComplianceStatus = "skip"
	ComplianceStatusERROR ComplianceStatus = "error"
)

func GetComplianceStatuses() []ComplianceStatus {
	return complianceStatuses
}

func GetPassedComplianceStatuses() []ComplianceStatus {
	passed := make([]ComplianceStatus, 0)
	for _, status := range complianceStatuses {
		if status.IsPassed() {
			passed = append(passed, status)
		}
	}
	return passed
}

func GetFailedComplianceStatuses() []ComplianceStatus {
	failed := make([]ComplianceStatus, 0)
	for _, status := range complianceStatuses {
		if !status.IsPassed() {
			failed = append(failed, status)
		}
	}
	return failed
}

func (r ComplianceStatus) IsPassed() bool {
	return r == ComplianceStatusOK || r == ComplianceStatusINFO || r == ComplianceStatusSKIP
}

type ComplianceStatusSummaryWithTotal struct {
	ComplianceStatusSummary
	TotalCount int `json:"totalCount" example:"5"`
}

type ComplianceStatusSummary struct {
	OkCount    int `json:"okCount" example:"1"`
	AlarmCount int `json:"alarmCount" example:"1"`
	InfoCount  int `json:"infoCount" example:"1"`
	SkipCount  int `json:"skipCount" example:"1"`
	ErrorCount int `json:"errorCount" example:"1"`
}

func (c *ComplianceStatusSummary) AddComplianceStatusSummary(summary ComplianceStatusSummary) {
	c.OkCount += summary.OkCount
	c.AlarmCount += summary.AlarmCount
	c.InfoCount += summary.InfoCount
	c.SkipCount += summary.SkipCount
	c.ErrorCount += summary.ErrorCount
}

func (c *ComplianceStatusSummary) AddComplianceStatusMap(summary map[ComplianceStatus]int) {
	c.OkCount += summary[ComplianceStatusOK]
	c.AlarmCount += summary[ComplianceStatusALARM]
	c.InfoCount += summary[ComplianceStatusINFO]
	c.SkipCount += summary[ComplianceStatusSKIP]
	c.ErrorCount += summary[ComplianceStatusERROR]
}

type ComplianceResultShortSummary struct {
	Passed int `json:"passed"`
	Failed int `json:"failed"`
}

var complianceStatuses = []ComplianceStatus{
	ComplianceStatusOK,
	ComplianceStatusALARM,
	ComplianceStatusINFO,
	ComplianceStatusSKIP,
	ComplianceStatusERROR,
}

func ParseComplianceStatus(s string) ComplianceStatus {
	s = strings.ToLower(s)
	for _, status := range complianceStatuses {
		if s == strings.ToLower(string(status)) {
			return status
		}
	}
	return ""
}

func ParseComplianceStatuses(list []string) []ComplianceStatus {
	result := make([]ComplianceStatus, 0, len(list))
	for _, s := range list {
		result = append(result, ParseComplianceStatus(s))
	}
	return result
}
