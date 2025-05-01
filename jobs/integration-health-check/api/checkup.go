package api

type CheckupJobStatus string

const (
	CheckupJobInProgress CheckupJobStatus = "IN_PROGRESS"
	CheckupJobFailed     CheckupJobStatus = "FAILED"
	CheckupJobSucceeded  CheckupJobStatus = "SUCCEEDED"
)
