package types

import (
	"strings"
)

type ComplianceResultSeverity string

const (
	ComplianceResultSeverityNone     ComplianceResultSeverity = "none"
	ComplianceResultSeverityLow      ComplianceResultSeverity = "low"
	ComplianceResultSeverityMedium   ComplianceResultSeverity = "medium"
	ComplianceResultSeverityHigh     ComplianceResultSeverity = "high"
	ComplianceResultSeverityCritical ComplianceResultSeverity = "critical"
)

func (s ComplianceResultSeverity) Level() int {
	switch s {
	case ComplianceResultSeverityNone:
		return 1
	case ComplianceResultSeverityLow:
		return 2
	case ComplianceResultSeverityMedium:
		return 3
	case ComplianceResultSeverityHigh:
		return 4
	case ComplianceResultSeverityCritical:
		return 5
	default:
		return 0
	}
}

func (s ComplianceResultSeverity) String() string {
	return string(s)
}

var complianceResultSeveritiesSeverities = []ComplianceResultSeverity{
	ComplianceResultSeverityNone,
	ComplianceResultSeverityLow,
	ComplianceResultSeverityMedium,
	ComplianceResultSeverityHigh,
	ComplianceResultSeverityCritical,
}

func ParseComplianceResultSeverity(s string) ComplianceResultSeverity {
	s = strings.ToLower(s)
	for _, sev := range complianceResultSeveritiesSeverities {
		if s == strings.ToLower(sev.String()) {
			return sev
		}
	}
	return ""
}

func ParseComplianceResultSeverities(list []string) []ComplianceResultSeverity {
	result := make([]ComplianceResultSeverity, 0, len(list))
	for _, s := range list {
		result = append(result, ParseComplianceResultSeverity(s))
	}
	return result
}

type SeverityResultWithTotal struct {
	SeverityResult
	TotalCount int `json:"totalCount" example:"5"`
}

type SeverityResult struct {
	NoneCount     int `json:"noneCount" example:"1"`
	LowCount      int `json:"lowCount" example:"1"`
	MediumCount   int `json:"mediumCount" example:"1"`
	HighCount     int `json:"highCount" example:"1"`
	CriticalCount int `json:"criticalCount" example:"1"`
}

type SeverityResultV2 struct {
	None     int `json:"none"`
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
	Total    int `json:"total"`
}

func (r *SeverityResultV2) AddSeverityResult(severity SeverityResult) {
	r.None += severity.NoneCount
	r.Low += severity.LowCount
	r.Medium += severity.MediumCount
	r.High += severity.HighCount
	r.Critical += severity.CriticalCount
	r.Total += severity.NoneCount + severity.LowCount + severity.MediumCount + severity.HighCount + severity.CriticalCount
}

func (r *SeverityResultV2) AddResultMap(result map[ComplianceResultSeverity]int) {
	r.None += result[ComplianceResultSeverityNone]
	r.Low += result[ComplianceResultSeverityLow]
	r.Medium += result[ComplianceResultSeverityMedium]
	r.High += result[ComplianceResultSeverityHigh]
	r.Critical += result[ComplianceResultSeverityCritical]
	r.Total += result[ComplianceResultSeverityNone] + result[ComplianceResultSeverityLow] + result[ComplianceResultSeverityMedium] + result[ComplianceResultSeverityHigh] + result[ComplianceResultSeverityCritical]
}

func (r *SeverityResult) AddSeverityResult(severity SeverityResult) {
	r.NoneCount += severity.NoneCount
	r.LowCount += severity.LowCount
	r.MediumCount += severity.MediumCount
	r.HighCount += severity.HighCount
	r.CriticalCount += severity.CriticalCount
}

func (r *SeverityResult) AddResultMap(result map[ComplianceResultSeverity]int) {
	r.NoneCount += result[ComplianceResultSeverityNone]
	r.LowCount += result[ComplianceResultSeverityLow]
	r.MediumCount += result[ComplianceResultSeverityMedium]
	r.HighCount += result[ComplianceResultSeverityHigh]
	r.CriticalCount += result[ComplianceResultSeverityCritical]
}

func (r *SeverityResult) IncreaseBySeverity(severity ComplianceResultSeverity) {
	switch severity {
	case ComplianceResultSeverityCritical:
		r.CriticalCount++
	case ComplianceResultSeverityHigh:
		r.HighCount++
	case ComplianceResultSeverityMedium:
		r.MediumCount++
	case ComplianceResultSeverityLow:
		r.LowCount++
	case ComplianceResultSeverityNone:
		r.LowCount++
	}
}

func (r *SeverityResult) IncreaseBySeverityByAmount(severity ComplianceResultSeverity, amount int) {
	switch severity {
	case ComplianceResultSeverityCritical:
		r.CriticalCount += amount
	case ComplianceResultSeverityHigh:
		r.HighCount += amount
	case ComplianceResultSeverityMedium:
		r.MediumCount += amount
	case ComplianceResultSeverityLow:
		r.LowCount += amount
	case ComplianceResultSeverityNone:
		r.LowCount += amount
	}
}
