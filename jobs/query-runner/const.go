package query_runner

import "time"

const (
	JobQueueTopic       = "query-runner-jobs-queue"
	JobResultQueueTopic = "query-runner-results-queue"
	ConsumerGroup       = "query-runner-worker"
	StreamName          = "query-runner-worker"

	JobTimeoutMinutes = 5
	JobTimeout        = JobTimeoutMinutes * time.Minute
)
