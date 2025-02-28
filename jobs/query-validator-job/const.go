package query_validator

import "time"

const (
	JobQueueTopic       = "query-validator-jobs-queue"
	JobResultQueueTopic = "query-validator-results-queue"
	ConsumerGroup       = "query-validator-worker"
	StreamName          = "query-validator"

	JobTimeoutMinutes = 5
	JobTimeout        = JobTimeoutMinutes * time.Minute
)

type QueryError string

const (
	MissingPlatformResourceIDQueryError QueryError = "missing_platform_resource_id"
	MissingAccountIDQueryError          QueryError = "missing_account_id"
	MissingResourceTypeQueryError       QueryError = "missing_resource_type"
	ResourceNotFoundQueryError          QueryError = "resource_not_found"
)
