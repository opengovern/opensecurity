package query_runner

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/jetstream"
	queryrunner "github.com/opengovern/opencomply/jobs/query-runner-job"
	"go.uber.org/zap"
)

func (s *JobScheduler) RunQueryRunnerReportJobResultsConsumer(ctx context.Context) error {
	if _, err := s.jq.Consume(ctx, "scheduler-query-runner", queryrunner.StreamName, []string{queryrunner.JobResultQueueTopic}, "scheduler-query-runner", func(msg jetstream.Msg) {
		if err := msg.Ack(); err != nil {
			s.logger.Error("Failed committing message", zap.Error(err))
		}

		var result queryrunner.JobResult
		if err := json.Unmarshal(msg.Data(), &result); err != nil {
			s.logger.Error("Failed to unmarshal ComplianceReportJob results", zap.Error(err))
			return
		}

		s.logger.Info("Processing ReportJobResult for Job",
			zap.Uint("jobId", result.ID),
			zap.String("status", string(result.Status)),
		)
		err := s.db.UpdateQueryRunnerJobStatus(result.ID, result.Status, result.FailureMessage)
		if err != nil {
			s.logger.Error("Failed to update the status of QueryRunnerReportJob",
				zap.Uint("jobId", result.ID),
				zap.Error(err))
			return
		}
	}); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
