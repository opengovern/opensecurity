package api

type IntegrationHealthCheckJobStatus string

const (
	IntegrationHealthCheckJobInProgress IntegrationHealthCheckJobStatus = "IN_PROGRESS"
	IntegrationHealthCheckJobFailed     IntegrationHealthCheckJobStatus = "FAILED"
	IntegrationHealthCheckJobSucceeded  IntegrationHealthCheckJobStatus = "SUCCEEDED"
)
