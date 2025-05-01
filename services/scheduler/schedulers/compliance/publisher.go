package compliance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	integrationapi "github.com/opengovern/opensecurity/services/integration/api/models"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"golang.org/x/net/context"

	runner "github.com/opengovern/opensecurity/jobs/compliance-runner"
	complianceApi "github.com/opengovern/opensecurity/services/compliance/api"
	"go.uber.org/zap"
)

func (s *JobScheduler) runPublisher(ctx context.Context, manuals bool) error {
	s.logger.Info("runPublisher")
	ctx2 := &httpclient.Context{UserRole: api.AdminRole}
	ctx2.Ctx = ctx
	connectionsMap := make(map[string]*integrationapi.Integration)
	integrations, err := s.integrationClient.ListIntegrations(ctx2, nil)
	if err != nil {
		s.logger.Error("failed to get connections", zap.Error(err))
		return err
	}
	for _, integration := range integrations.Integrations {
		integration := integration
		connectionsMap[integration.IntegrationID] = &integration
	}

	queries, err := s.complianceClient.ListQueries(ctx2)
	if err != nil {
		s.logger.Error("failed to get queries", zap.Error(err))
		return err
	}
	queriesMap := make(map[string]*complianceApi.Policy)
	for _, query := range queries {
		query := query
		queriesMap[query.ID] = &query
	}

	for i := 0; i < 10; i++ {
		err := s.db.UpdateTimeoutQueuedRunnerJobs()
		if err != nil {
			s.logger.Error("failed to update timed out runners", zap.Error(err))
		}

		err = s.db.UpdateTimedOutInProgressRunners()
		if err != nil {
			s.logger.Error("failed to update timed out runners", zap.Error(err))
		}

		runners, err := s.db.FetchCreatedRunners(manuals)
		if err != nil {
			s.logger.Error("failed to fetch created runners", zap.Error(err))
			continue
		}

		if len(runners) == 0 {
			s.logger.Info("no created runners found skipping")
			break
		}

		for _, it := range runners {
			query, ok := queriesMap[it.PolicyID]
			if !ok || query == nil {
				s.logger.Error("query not found", zap.String("queryId", it.PolicyID), zap.Uint("runnerId", it.ID))
				_ = s.db.UpdateRunnerJob(it.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, "query not found", nil)
				continue
			}

			callers, err := it.GetCallers()
			if err != nil {
				s.logger.Error("failed to get callers", zap.Error(err), zap.Uint("runnerId", it.ID))
				_ = s.db.UpdateRunnerJob(it.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, "failed to get callers", nil)
				continue
			}
			var providerID *string
			if it.IntegrationID != nil && *it.IntegrationID != "" {
				if _, ok := connectionsMap[*it.IntegrationID]; ok {
					providerID = &connectionsMap[*it.IntegrationID].ProviderID
				} else {
					_ = s.db.UpdateRunnerJob(it.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, "integration does not exist", nil)
					continue
				}
			}
			job := runner.Job{
				ID:          it.ID,
				RetryCount:  it.RetryCount,
				ParentJobID: it.ParentJobID,
				CreatedAt:   it.CreatedAt,
				ExecutionPlan: runner.ExecutionPlan{
					Callers:       callers,
					Query:         *query,
					ControlID:     it.ControlID,
					IntegrationID: it.IntegrationID,
					ProviderID:    providerID,
				},
			}

			jobJson, err := json.Marshal(job)
			if err != nil {
				_ = s.db.UpdateRunnerJob(job.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, err.Error(), nil)
				s.logger.Error("failed to marshal job", zap.Error(err), zap.Uint("runnerId", it.ID))
				continue
			}

			s.logger.Info("publishing runner", zap.Uint("jobId", job.ID))
			topic := runner.JobQueueTopic
			if it.TriggerType == model.ComplianceTriggerTypeManual {
				topic = runner.JobQueueTopicManuals
			}
			seqNum, err := s.jq.Produce(ctx, topic, jobJson, fmt.Sprintf("job-%d-%d", job.ID, it.RetryCount))
			if err != nil {
				if err.Error() == "nats: no response from stream" {
					err = s.runSetupNatsStreams(ctx)
					if err != nil {
						s.logger.Error("Failed to setup nats streams", zap.Error(err))
						return err
					}
					seqNum, err = s.jq.Produce(ctx, topic, jobJson, fmt.Sprintf("job-%d-%d", job.ID, it.RetryCount))
					if err != nil {
						_ = s.db.UpdateRunnerJob(job.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, err.Error(), nil)
						s.logger.Error("failed to send job", zap.Error(err), zap.Uint("runnerId", it.ID))
						continue
					}
				} else {
					_ = s.db.UpdateRunnerJob(job.ID, model.ComplianceRunnerFailed, nil, nil, nil, nil, err.Error(), nil)
					s.logger.Error("failed to send job", zap.Error(err), zap.Uint("runnerId", it.ID), zap.String("error message", err.Error()))
					continue
				}
			}

			if seqNum != nil {
				_ = s.db.UpdateRunnerJobNatsSeqNum(job.ID, *seqNum)
			}
			now := time.Now()
			_ = s.db.UpdateRunnerJob(job.ID, model.ComplianceRunnerQueued, &now, nil, nil, nil, "", nil)
		}
	}

	err = s.db.RetryFailedRunners()
	if err != nil {
		s.logger.Error("failed to retry failed runners", zap.Error(err))
		return err
	}

	return nil
}
