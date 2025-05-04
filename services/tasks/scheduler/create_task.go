package scheduler

import (
	"context"
	"time"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
)

func (s *MainScheduler) CreateTaskScheduler(ctx context.Context) {
	s.logger.Info("Scheduling publisher on a timer")

	t := ticker.NewTicker(1*time.Minute, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.createTasks(ctx); err != nil {
			s.logger.Error("failed to run compliance publisher", zap.Error(err))
			continue
		}
	}
}

func (s *MainScheduler) createTasks(ctx context.Context) error {
	ctx2 := &httpclient.Context{UserRole: api.AdminRole}
	ctx2.Ctx = ctx

	s.logger.Info("Create Task on schedule started")
	tasks, err := s.db.GetEnabledTaskList()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		runSchedules, err := s.db.GetTaskRunSchedules(task.ID)
		if err != nil {
			return err
		}
		for _, runSchedule := range runSchedules {
			lastRun, err := s.db.FetchLastTaskRunsByTaskSchedulerID(task.ID, runSchedule.ID)
			if err != nil {
				return err
			}
			if lastRun != nil {
				if time.Now().Before(lastRun.CreatedAt.Add(time.Duration(runSchedule.Frequency) * time.Second)) {
					continue
				}
			}

			newRun := models.TaskRun{
				TaskID:      task.ID,
				Status:      models.TaskRunStatusCreated,
				TriggerType: models.TriggerTypeScheduled,
				TriggeredBy: runSchedule.ID,
			}

			err = newRun.Result.Set([]byte("{}"))
			if err != nil {
				return err
			}
			newRun.Params = runSchedule.Params

			if err = s.db.CreateTaskRun(&newRun); err != nil {
				return err
			}
		}
	}

	return nil
}
