package query_validator

type QueryValidatorStatus string

const (
	QueryValidatorCreated    QueryValidatorStatus = "CREATED"
	QueryValidatorQueued     QueryValidatorStatus = "QUEUED"
	QueryValidatorInProgress QueryValidatorStatus = "IN_PROGRESS"
	QueryValidatorSucceeded  QueryValidatorStatus = "SUCCEEDED"
	QueryValidatorFailed     QueryValidatorStatus = "FAILED"
	QueryValidatorTimeOut    QueryValidatorStatus = "TIMEOUT"
)

type JobResult struct {
	ID             uint                 `json:"id"`
	QueryType      QueryType            `json:"query_type"`
	ControlId      string               `json:"control_id"`
	QueryId        string               `json:"query_id"`
	Status         QueryValidatorStatus `json:"status"`
	FailureMessage string               `json:"failure_message"`
}
