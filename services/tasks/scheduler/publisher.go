package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/tasks"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func (s *TaskScheduler) runPublisher(ctx context.Context) error {
	ctx2 := &httpclient.Context{UserRole: api.AdminRole}
	ctx2.Ctx = ctx

	s.logger.Info("Policy Runner publisher started")

	err := s.db.TimeoutTaskRunsByTaskID(s.TaskID, s.Timeout)
	if err != nil {
		s.logger.Error("failed to timeout task runs", zap.String("task_id", s.TaskID),
			zap.Uint64("timeout", s.Timeout), zap.Error(err))
		return err
	}

	runs, err := s.db.FetchCreatedTaskRunsByTaskID(s.TaskID)
	if err != nil {
		s.logger.Error("failed to get task runs", zap.Error(err))
		return err
	}

	for _, run := range runs {
		params, err := JSONBToMap(run.Params)
		if err != nil {
			result := pgtype.JSONB{}
			_ = result.Set([]byte("{}"))
			_ = s.db.UpdateTaskRun(run.ID, models.TaskRunStatusFailed, result, "failed to get params")
			s.logger.Error("failed to get params", zap.Error(err), zap.Uint("runId", run.ID))
			return err
		}
		req := tasks.TaskRequest{
			EsDeliverEndpoint:         s.cfg.ESSinkEndpoint,
			IngestionPipelineEndpoint: s.cfg.ElasticSearch.IngestionEndpoint,
			UseOpenSearch:             s.cfg.ElasticSearch.IsOpenSearch,
			TaskDefinition: tasks.TaskDefinition{
				RunID:      run.ID,
				TaskType:   s.TaskID,
				ResultType: s.ResultType,
				Params:     params,
			},
			ExtraInputs: nil,
		}
		reqJson, err := json.Marshal(req)
		if err != nil {
			result := pgtype.JSONB{}
			_ = result.Set([]byte("{}"))
			_ = s.db.UpdateTaskRun(run.ID, models.TaskRunStatusFailed, result, "failed to marshal run")
			s.logger.Error("failed to marshal Task Run", zap.Error(err), zap.Uint("runId", run.ID))
			return err
		}

		s.logger.Info("publishing audit job", zap.Uint("runId", run.ID), zap.String("topic", s.NatsConfig.Topic))
		_, err = s.jq.Produce(ctx, s.NatsConfig.Topic, reqJson, fmt.Sprintf("run-%d", run.ID))
		if err != nil {
			if err.Error() == "nats: no response from stream" {
				err = s.runSetupNatsStreams(ctx)
				if err != nil {
					s.logger.Error("Failed to setup nats streams", zap.Error(err))
					return err
				}
				_, err = s.jq.Produce(ctx, s.NatsConfig.Topic, reqJson, fmt.Sprintf("run-%d", run.ID))
				if err != nil {
					result := pgtype.JSONB{}
					_ = result.Set([]byte("{}"))
					_ = s.db.UpdateTaskRun(run.ID, models.TaskRunStatusFailed, result, err.Error())
					s.logger.Error("failed to send run", zap.Error(err), zap.Uint("runId", run.ID))
					continue
				}
			} else {
				result := pgtype.JSONB{}
				_ = result.Set([]byte("{}"))
				_ = s.db.UpdateTaskRun(run.ID, models.TaskRunStatusFailed, result, err.Error())
				s.logger.Error("failed to send run", zap.Error(err), zap.Uint("runId", run.ID), zap.String("error message", err.Error()))
				continue
			}
		} else {
			result := pgtype.JSONB{}
			_ = result.Set([]byte("{}"))
			_ = s.db.UpdateTaskRun(run.ID, models.TaskRunStatusQueued, result, "")
		}
	}

	return nil
}

func JSONBToMap(jsonb pgtype.JSONB) (map[string]any, error) {
	if jsonb.Status != pgtype.Present {
		return nil, fmt.Errorf("JSONB data is not present")
	}

	var result map[string]any
	if err := json.Unmarshal(jsonb.Bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	return result, nil
}
