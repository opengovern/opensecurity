package scheduler

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
)

func (s *TaskScheduler) createTask(ctx context.Context) error {
	ctx2 := &httpclient.Context{UserRole: api.AdminRole}
	ctx2.Ctx = ctx

	s.logger.Info("Create Task on schedule started")

	run, err := s.db.FetchLastCreatedTaskRunsByTaskID(s.TaskID)
	if err != nil {
		s.logger.Error("failed to get task runs", zap.Error(err))
		return err
	}

	newRun := models.TaskRun{
		TaskID: s.TaskID,
		Status: models.TaskRunStatusCreated,
	}

	err = newRun.Result.Set([]byte("{}"))
	if err != nil {
		return err
	}
	newRun.Params = run.Params

	if err = s.db.CreateTaskRun(&newRun); err != nil {
		return err
	}

	return nil
}
