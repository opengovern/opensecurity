package discovery

import (
	"context"
	"encoding/json"

	es2 "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/opensecurity/services/scheduler/api"
	"github.com/opengovern/opensecurity/services/scheduler/es"

	"strings"
	"time"

	"go.uber.org/zap"
)

const OldResourceDeleterInterval = 1 * time.Minute

func (s *Scheduler) OldResourceDeleter(ctx context.Context) {
	s.logger.Info("Scheduling OldResourceDeleter on a timer")

	t := ticker.NewTicker(OldResourceDeleterInterval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.runDeleter(ctx); err != nil {
			s.logger.Error("failed to run deleter", zap.Error(err))
			continue
		}
	}
}

func (s *Scheduler) runDeleter(ctx context.Context) error {
	s.logger.Info("runDeleter")

	tasks, err := es.GetDeleteTasks(ctx, s.esClient)
	if err != nil {
		s.logger.Error("failed to get delete tasks", zap.Error(err))
		return err
	}

	for _, task := range tasks.Hits.Hits {
		switch task.Source.TaskType {
		case es.DeleteTaskTypeResource:
			job, err := s.db.GetDescribeIntegrationJobByID(task.Source.DiscoveryJobID)
			if err != nil {
				s.logger.Error("failed to get describe connection job", zap.Error(err))
				continue
			}
			if job == nil || job.Status != api.DescribeResourceJobOldResourceDeletion {
				continue
			}
			s.logger.Info("deleting resources", zap.String("task", task.ID), zap.Uint("job", job.ID), zap.String("IntegrationID", job.IntegrationID), zap.String("resourceType", task.Source.ResourceType))
			for _, resource := range task.Source.DeletingResources {
				err = s.esClient.Delete(string(resource.Key), resource.Index)
				if err != nil {
					if strings.Contains(err.Error(), "[404 Not Found]") {
						s.logger.Warn("resource not found", zap.String("resource", string(resource.Key)), zap.String("index", resource.Index), zap.Error(err))
						continue
					}
					s.logger.Error("failed to delete resource", zap.Error(err))
					return err
				}
			}
			err = s.db.UpdateDescribeIntegrationJobStatus(job.ID, api.DescribeResourceJobSucceeded, job.FailureMessage, job.ErrorCode, job.DescribedResourceCount, job.DeletingCount)
			if err != nil {
				s.logger.Error("failed to update describe connection job status", zap.Error(err))
				continue
			}
		case es.DeleteTaskTypeQuery:
			var query any
			err = json.Unmarshal([]byte(task.Source.Query), &query)
			if err != nil {
				s.logger.Error("failed to unmarshal query", zap.Error(err))
				return err
			}
			s.logger.Info("deleting by query", zap.String("task", task.ID), zap.String("queryIndex", task.Source.QueryIndex), zap.Any("query", query))
			_, err = es2.DeleteByQuery(ctx, s.esClient.ES(), []string{task.Source.QueryIndex}, query)
			if err != nil {
				s.logger.Error("failed to delete by query", zap.Error(err))
				return err
			}
		}

		err = s.esClient.Delete(task.ID, es.DeleteTasksIndex)
		if err != nil {
			s.logger.Error("failed to delete task", zap.Error(err))
			return err
		}
	}

	return nil
}
