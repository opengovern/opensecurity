package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jackc/pgtype"
	"strings"
	"time"

	"github.com/opengovern/opensecurity/pkg/types"

	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"

	"github.com/nats-io/nats.go/jetstream"
	es2 "github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/opensecurity/services/scheduler/api"
	"github.com/opengovern/opensecurity/services/scheduler/es"
	"go.uber.org/zap"
)

func (s *Scheduler) RunDescribeJobResultsConsumer(ctx context.Context) error {
	s.logger.Info("Consuming messages from the JobResults queue")

	consumeCtx, err := s.jq.Consume(
		ctx,
		"describe-receiver",
		DescribeStreamName,
		[]string{DescribeResultsQueueName},
		"describe-receiver",
		func(msg jetstream.Msg) {
			var result DescribeJobResult
			if err := json.Unmarshal(msg.Data(), &result); err != nil {
				ResultsProcessedCount.WithLabelValues("", "failure").Inc()

				s.logger.Error("failed to consume message from describeJobResult", zap.Error(err))

				// the job cannot be parsed into json, so send ack and throw message away.
				if err := msg.Ack(); err != nil {
					s.logger.Error("failure while sending ack for message", zap.Error(err))
				}

				return
			}

			s.logger.Info("Processing JobResult for Job",
				zap.Uint("jobId", result.JobID),
				zap.String("status", string(result.Status)),
			)

			job, err := s.db.GetDescribeIntegrationJobByID(result.JobID)
			if err != nil {
				s.logger.Error("failed to get describe integration job by id", zap.Uint("job_id", result.JobID))
			}
			if job == nil {
				return
			}
			var params map[string]string
			if job.Parameters.Status == pgtype.Present {
				if err := json.Unmarshal(job.Parameters.Bytes, &params); err != nil {
					s.logger.Error("failed to unmarshal parameters", zap.Uint("job_id", result.JobID))
				}
			} else {
				params = make(map[string]string)
			}

			var deletedCount int64
			if s.DoDeleteOldResources && result.Status == api.DescribeResourceJobSucceeded {
				result.Status = api.DescribeResourceJobOldResourceDeletion

				dlc, err := s.cleanupOldResources(ctx, result, params)
				if err != nil {
					ResultsProcessedCount.WithLabelValues(string(result.DescribeJob.IntegrationType), "failure").Inc()
					s.logger.Error("failed to cleanupOldResources", zap.Error(err))

					if err := msg.Nak(); err != nil {
						s.logger.Error("failure while sending not-ack for message", zap.Error(err))
					}

					return
				}

				deletedCount = dlc
				if deletedCount == 0 {
					result.Status = api.DescribeResourceJobSucceeded
				}
			}

			errStr := strings.ReplaceAll(result.Error, "\x00", "")
			errCodeStr := strings.ReplaceAll(result.ErrorCode, "\x00", "")
			if errCodeStr == "" {
				if strings.Contains(errStr, "exceeded maximum number of attempts") {
					errCodeStr = "TooManyRequestsException"
				} else if strings.Contains(errStr, "context deadline exceeded") {
					errCodeStr = "ContextDeadlineExceeded"
				}
			}

			s.logger.Info("updating job status", zap.Uint("jobID", result.JobID), zap.String("status", string(result.Status)))
			if err := s.db.UpdateDescribeIntegrationJobStatus(result.JobID, result.Status, errStr, errCodeStr, int64(len(result.DescribedResourceIDs)), deletedCount); err != nil {
				ResultsProcessedCount.WithLabelValues(string(result.DescribeJob.IntegrationType), "failure").Inc()

				s.logger.Error("failed to UpdateDescribeResourceJobStatus", zap.Error(err))

				if err := msg.Nak(); err != nil {
					s.logger.Error("failure while sending not-ack for message", zap.Error(err))
				}

				return
			}

			ResultsProcessedCount.WithLabelValues(string(result.DescribeJob.IntegrationType), "successful").Inc()

			if err := msg.Ack(); err != nil {
				s.logger.Error("failure while sending ack for message", zap.Error(err))
			}
		},
	)
	if err != nil {
		return err
	}

	t := ticker.NewTicker(JobTimeoutCheckInterval, time.Minute*1)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			s.handleTimeoutForDiscoveryJobs()
		case <-ctx.Done():
			consumeCtx.Drain()
			consumeCtx.Stop()

			return nil
		}
	}
}

func (s *Scheduler) handleTimeoutForDiscoveryJobs() {
	err := s.db.UpdateDescribeIntegrationJobsTimedOut(int64(s.discoveryIntervalHours.Hours()))
	if err != nil {
		s.logger.Error("failed to UpdateDescribeConnectionJobsTimedOut", zap.Error(err))
	}
}

func (s *Scheduler) cleanupOldResources(ctx context.Context, res DescribeJobResult, jobParams map[string]string) (int64, error) {
	var searchAfter []any

	var additionalFilters []map[string]any

	deletedCount := 0

	s.logger.Info("starting to schedule deleting old resources",
		zap.Uint("jobId", res.JobID),
		zap.String("integration_id", res.DescribeJob.IntegrationID),
		zap.String("resource_type", res.DescribeJob.ResourceType),
	)

	for {
		esResp, err := es.GetResourceIDsForAccountResourceTypeFromES(
			ctx,
			s.es,
			res.DescribeJob.IntegrationID,
			res.DescribeJob.ResourceType,
			es2.ConvertMapToString(jobParams),
			additionalFilters,
			searchAfter,
			1000)
		if err != nil {
			CleanupJobCount.WithLabelValues("failure").Inc()
			s.logger.Error("CleanJob failed",
				zap.Error(err))
			return 0, err
		}

		if len(esResp.Hits.Hits) == 0 {
			break
		}
		task := es.DeleteTask{
			DiscoveryJobID:  res.JobID,
			IntegrationID:   res.DescribeJob.IntegrationID,
			ResourceType:    res.DescribeJob.ResourceType,
			IntegrationType: res.DescribeJob.IntegrationType,
			TaskType:        es.DeleteTaskTypeResource,
		}

		for _, hit := range esResp.Hits.Hits {
			searchAfter = hit.Sort
			esResourceID := hit.Source.ResourceID

			exists := false
			for _, describedResourceID := range res.DescribedResourceIDs {
				if esResourceID == describedResourceID {
					exists = true
					break
				}
			}

			if !exists {
				OldResourcesDeletedCount.WithLabelValues(string(res.DescribeJob.IntegrationType)).Inc()
				resource := es2.Resource{
					ResourceID:      esResourceID,
					IntegrationID:   res.DescribeJob.IntegrationID,
					ResourceType:    res.DescribeJob.ResourceType,
					IntegrationType: res.DescribeJob.IntegrationType,
				}
				keys, idx := resource.KeysAndIndex()
				deletedCount += 1
				task.DeletingResources = append(task.DeletingResources, es.DeletingResource{
					Key:        []byte(es2.HashOf(keys...)),
					ResourceID: esResourceID,
					Index:      idx,
				})

				lookupResource := es2.LookupResource{
					ResourceID:      esResourceID,
					IntegrationID:   res.DescribeJob.IntegrationID,
					ResourceType:    res.DescribeJob.ResourceType,
					IntegrationType: res.DescribeJob.IntegrationType,
				}
				lookUpKeys, lookUpIdx := lookupResource.KeysAndIndex()
				deletedCount += 1
				task.DeletingResources = append(task.DeletingResources, es.DeletingResource{
					Key:        []byte(es2.HashOf(lookUpKeys...)),
					ResourceID: esResourceID,
					Index:      lookUpIdx,
				})
			}
		}

		i := 0
		for {
			taskKeys, taskIdx := task.KeysAndIndex()
			task.EsID = es2.HashOf(taskKeys...)
			task.EsIndex = taskIdx

			if len(task.DeletingResources) > 0 {
				if _, err := s.sinkClient.Ingest(&httpclient.Context{UserRole: authApi.AdminRole}, []es2.Doc{task}); err != nil {
					s.logger.Error("failed to send delete message to elastic",
						zap.Uint("jobId", res.JobID),
						zap.String("integration_id", res.DescribeJob.IntegrationID),
						zap.String("resource_type", res.DescribeJob.ResourceType),
						zap.Error(err))
					if i > 10 {
						CleanupJobCount.WithLabelValues("failure").Inc()
						return 0, err
					}
					i++
					continue
				}
			}
			break
		}
	}

	s.logger.Info("scheduled deleting old resources",
		zap.Uint("jobId", res.JobID),
		zap.String("connection_id", res.DescribeJob.IntegrationID),
		zap.String("resource_type", res.DescribeJob.ResourceType),
		zap.Int("deleted_count", deletedCount))

	CleanupJobCount.WithLabelValues("successful").Inc()
	return int64(deletedCount), nil
}

func (s *Scheduler) cleanupDescribeResourcesNotInIntegrations(ctx context.Context, integrationIDs []string) {
	s.logger.Info("starting cleaning up resources")
	var searchAfter []any
	totalDeletedCount := 0
	deletedIntegrationIDs := make(map[string]bool)
	for {
		esResp, err := es.GetResourceIDsNotInIntegrationsFromES(ctx, s.es, integrationIDs, searchAfter, 1000)
		if err != nil {
			s.logger.Error("failed to get resource ids from es", zap.Error(err))
			break
		}
		totalDeletedCount += len(esResp.Hits.Hits)
		if len(esResp.Hits.Hits) == 0 {
			s.logger.Info("no old resource found")
			break
		} else {
			s.logger.Info("deleting resources", zap.Int("number of resources", len(esResp.Hits.Hits)))
		}
		deletedCount := 0
		for _, hit := range esResp.Hits.Hits {
			deletedIntegrationIDs[hit.Source.IntegrationID] = true
			searchAfter = hit.Sort

			resource := es2.Resource{
				ResourceID:      hit.Source.ResourceID,
				IntegrationID:   hit.Source.IntegrationID,
				ResourceType:    strings.ToLower(hit.Source.ResourceType),
				IntegrationType: hit.Source.IntegrationType,
			}
			keys, idx := resource.KeysAndIndex()
			deletedCount += 1
			key := es2.HashOf(keys...)
			resource.EsID = key
			resource.EsIndex = idx
			err = s.es.Delete(key, idx)
			if err != nil {
				if !strings.Contains(err.Error(), "404 Not Found") {
					s.logger.Error("failed to delete resource from open-search", zap.Error(err))
				}
			}

			lookupResource := es2.LookupResource{
				ResourceID:      hit.Source.ResourceID,
				IntegrationID:   hit.Source.IntegrationID,
				ResourceType:    strings.ToLower(hit.Source.ResourceType),
				IntegrationType: hit.Source.IntegrationType,
			}
			deletedCount += 1
			keys, idx = lookupResource.KeysAndIndex()
			key = es2.HashOf(keys...)
			lookupResource.EsID = key
			lookupResource.EsIndex = idx
			err = s.es.Delete(key, idx)
			if err != nil {
				if !strings.Contains(err.Error(), "404 Not Found") {
					s.logger.Error("failed to delete lookup from open-search", zap.Error(err))
				}
			}

			resourceFinding := types.ResourceFinding{
				PlatformResourceID: hit.Source.PlatformID,
				ResourceType:       strings.ToLower(hit.Source.ResourceType),
			}
			deletedCount += 1
			keys, idx = resourceFinding.KeysAndIndex()
			key = es2.HashOf(keys...)
			resourceFinding.EsID = key
			resourceFinding.EsIndex = idx
			err = s.es.Delete(key, idx)
			if err != nil {
				if !strings.Contains(err.Error(), "404 Not Found") {
					s.logger.Error("failed to delete resource finding from open-search", zap.Error(err))
				}
			}
		}
		s.logger.Info("deleted resource count", zap.Int("count", deletedCount),
			zap.Any("deleted integrations", deletedIntegrationIDs))
	}
	s.logger.Info("total deleted resource count", zap.Int("count", totalDeletedCount),
		zap.Any("deleted integrations", deletedIntegrationIDs))
	return
}

func (s *Scheduler) cleanupDescribeResourcesForConnectionAndResourceType(IntegrationID, resourceType string) error {
	root := make(map[string]any)
	root["query"] = map[string]any{
		"bool": map[string]any{
			"filter": []any{
				map[string]any{
					"term": map[string]any{
						"integration_id": IntegrationID,
					},
				},
				map[string]any{
					"term": map[string]any{
						"resource_type": strings.ToLower(resourceType),
					},
				},
			},
		},
	}
	query, err := json.Marshal(root)
	if err != nil {
		return err
	}

	index := es2.ResourceTypeToESIndex(resourceType)
	res, err := s.es.ES().DeleteByQuery([]string{index}, bytes.NewReader(query))
	if err != nil {
		return err
	}

	opengovernance.CloseSafe(res)

	res, err = s.es.ES().DeleteByQuery([]string{InventorySummaryIndex}, bytes.NewReader(query))
	if err != nil {
		return err
	}

	opengovernance.CloseSafe(res)
	return nil
}
