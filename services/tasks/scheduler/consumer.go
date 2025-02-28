package scheduler

import (
	"bytes"
	"encoding/json"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type TaskResponse struct {
	RunID          uint                 `json:"run_id"`
	Status         models.TaskRunStatus `json:"status"`
	FailureMessage string               `json:"failure_message"`
	Result         []byte               `json:"result"`
}

func (s *TaskScheduler) RunTaskResponseConsumer(ctx context.Context) error {
	if _, err := s.jq.Consume(ctx, s.NatsConfig.ResultConsumer, s.NatsConfig.Stream, []string{s.NatsConfig.ResultTopic},
		s.NatsConfig.ResultConsumer, func(msg jetstream.Msg) {
			if err := msg.Ack(); err != nil {
				s.logger.Error("Failed committing message", zap.Error(err))
			}

			var response TaskResponse
			if err := json.Unmarshal(msg.Data(), &response); err != nil {
				s.logger.Error("Failed to unmarshal ComplianceReportJob results", zap.Error(err))
				return
			}

			taskRunUpdate := models.TaskRun{
				Status:         response.Status,
				FailureMessage: response.FailureMessage,
			}
			emptyResult := []byte("")
			if response.Result == nil || len(response.Result) == 0 || bytes.Equal(response.Result, emptyResult) {
				response.Result = []byte("{}")
			}
			err := taskRunUpdate.Result.Set(response.Result)
			if err != nil {
				s.logger.Error("failed to set result", zap.Error(err))
				return
			}
			err = s.db.UpdateTaskRun(response.RunID, taskRunUpdate.Status, taskRunUpdate.Result, taskRunUpdate.FailureMessage)
			if err != nil {
				s.logger.Error("Failed to update the status of RunTaskResponse",
					zap.String("Task", s.TaskID),
					zap.Uint("RunID", response.RunID),
					zap.Error(err))
				return
			}
		}); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
