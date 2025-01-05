package compliance_quick_run

import (
	"context"
	"encoding/json"
	auditjob "github.com/opengovern/opencomply/jobs/compliance-quick-run-job"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *JobScheduler) RunAuditJobResultsConsumer(ctx context.Context) error {
	if _, err := s.jq.Consume(ctx, "scheduler-compliance-quick-run", auditjob.StreamName, []string{auditjob.ResultQueueTopic},
		"scheduler-compliance-quick-run", func(msg jetstream.Msg) {
			if err := msg.Ack(); err != nil {
				s.logger.Error("Failed committing message", zap.Error(err))
			}

			var result auditjob.JobResult
			if err := json.Unmarshal(msg.Data(), &result); err != nil {
				s.logger.Error("Failed to unmarshal ComplianceReportJob results", zap.Error(err))
				return
			}

			s.logger.Info("Processing ReportJobResult for Job",
				zap.Uint("jobId", result.JobID),
				zap.String("status", string(result.Status)),
			)
			err := s.db.UpdateComplianceJob(result.JobID, result.Status, result.FailureMessage)
			if err != nil {
				s.logger.Error("Failed to update the status of QueryRunnerReportJob",
					zap.Uint("jobId", result.JobID),
					zap.Error(err))
				return
			}
		}); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
