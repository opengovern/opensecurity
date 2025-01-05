package compliance

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	runner "github.com/opengovern/opencomply/jobs/compliance-runner-job"
	summarizer "github.com/opengovern/opencomply/jobs/compliance-summarizer-job"
	"go.uber.org/zap"
)

const JobTimeoutCheckInterval = 5 * time.Minute

func (s *JobScheduler) RunComplianceReportJobResultsConsumer(ctx context.Context) error {
	if _, err := s.jq.Consume(ctx, "scheduler-runner-compliance", runner.StreamName, []string{runner.ResultQueueTopic}, "scheduler-runner-compliance", func(msg jetstream.Msg) {
		if err := msg.Ack(); err != nil {
			s.logger.Error("Failed committing message", zap.Error(err))
		}

		var result runner.JobResult
		if err := json.Unmarshal(msg.Data(), &result); err != nil {
			s.logger.Error("Failed to unmarshal ComplianceReportJob results", zap.Error(err))
			return
		}

		s.logger.Info("Processing ReportJobResult for Job",
			zap.Uint("jobId", result.Job.ID),
			zap.String("status", string(result.Status)),
		)
		err := s.db.UpdateRunnerJob(result.Job.ID, result.Status, result.StartedAt, result.TotalComplianceResultCount, result.Error)
		if err != nil {
			s.logger.Error("Failed to update the status of ComplianceReportJob",
				zap.Uint("jobId", result.Job.ID),
				zap.Error(err))
			return
		}
	}); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *JobScheduler) RunComplianceSummarizerResultsConsumer(ctx context.Context) error {
	if _, err := s.jq.Consume(
		ctx,
		"scheduler-summarizer-compliance",
		summarizer.StreamName,
		[]string{summarizer.ResultQueueTopic},
		"scheduler-summarizer-compliance",
		func(msg jetstream.Msg) {
			var result summarizer.JobResult
			if err := json.Unmarshal(msg.Data(), &result); err != nil {
				s.logger.Error("Failed to unmarshal ComplianceSummarizer results", zap.Error(err))

				if err := msg.Ack(); err != nil {
					s.logger.Error("Failed committing message", zap.Error(err))
				}
				return
			}

			s.logger.Info("Processing SummarizerResult for Job",
				zap.Uint("jobId", result.Job.ID),
				zap.String("status", string(result.Status)),
			)
			err := s.db.UpdateSummarizerJob(result.Job.ID, result.Status, result.StartedAt, result.Error)
			if err != nil {
				s.logger.Error("Failed to update the status of Summarizer",
					zap.Uint("jobId", result.Job.ID),
					zap.Error(err))

				if err := msg.Ack(); err != nil {
					s.logger.Error("Failed committing message", zap.Error(err))
				}

				return
			}

			if err := msg.Ack(); err != nil {
				s.logger.Error("Failed committing message", zap.Error(err))
			}
		},
	); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
