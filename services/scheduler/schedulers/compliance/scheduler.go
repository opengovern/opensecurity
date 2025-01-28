package compliance

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	runner "github.com/opengovern/opencomply/jobs/compliance-runner-job"
	integrationapi "github.com/opengovern/opencomply/services/integration/api/models"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.uber.org/zap"
	"time"
)

func (s *JobScheduler) runScheduler() error {
	s.logger.Info("scheduleComplianceJob")
	if s.complianceIntervalHours <= 0 {
		s.logger.Info("compliance interval is negative or zero, skipping compliance job scheduling")
		return nil
	}
	clientCtx := &httpclient.Context{UserRole: api.AdminRole}

	frameworks, err := s.complianceClient.ListBenchmarks(clientCtx, nil, nil)
	if err != nil {
		s.logger.Error("error while listing frameworks", zap.Error(err))
		return fmt.Errorf("error while listing frameworks: %v", err)
	}

	allIntegrations, err := s.integrationClient.ListIntegrations(clientCtx, nil)
	if err != nil {
		s.logger.Error("error while listing allConnections", zap.Error(err))
		return fmt.Errorf("error while listing allConnections: %v", err)
	}
	integrationsMap := make(map[string]*integrationapi.Integration)
	for _, connection := range allIntegrations.Integrations {
		connection := connection
		integrationsMap[connection.IntegrationID] = &connection
	}

	for _, framework := range frameworks {
		if !framework.Enabled {
			continue
		}
		var integrationIDs []string
		assignments, err := s.complianceClient.ListAssignmentsByBenchmark(clientCtx, framework.ID)
		if err != nil {
			s.logger.Error("error while listing assignments", zap.Error(err))
			return fmt.Errorf("error while listing assignments: %v", err)
		}

		for _, assignment := range assignments.Integrations {
			if !assignment.Status {
				continue
			}

			if _, ok := integrationsMap[assignment.IntegrationID]; !ok {
				continue
			}
			integration := integrationsMap[assignment.IntegrationID]

			if integration.State != integrationapi.IntegrationStateActive {
				continue
			}

			integrationIDs = append(integrationIDs, integration.IntegrationID)
		}

		if len(integrationIDs) == 0 {
			continue
		}

		complianceJob, err := s.db.GetLastComplianceJob(true, framework.ID)
		if err != nil {
			s.logger.Error("error while getting last compliance job", zap.Error(err))
			return err
		}

		timeAt := time.Now().Add(-s.complianceIntervalHours)
		if complianceJob == nil ||
			complianceJob.CreatedAt.Before(timeAt) {

			_, err := s.CreateComplianceReportJobs(true, framework.ID, complianceJob, integrationIDs, false, "system", nil)
			if err != nil {
				s.logger.Error("error while creating compliance job", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (s *JobScheduler) updateRunnersState() error {
	complianceJobs, err := s.db.ListComplianceJobsByStatus(aws.Bool(true), model.ComplianceJobRunnersInProgress)
	if err != nil {
		return fmt.Errorf("error while listing compliance jobs: %v", err)
	}
	for _, complianceJob := range complianceJobs {
		status := model.ComplianceRunnersStatus{}
		runners, err := s.db.ListComplianceJobRunnersWithParentID(complianceJob.ID)
		if err != nil {
			return fmt.Errorf("error while listing compliance runners: %v", err)
		}
		for _, r := range runners {
			switch r.Status {
			case runner.ComplianceRunnerCreated:
				status.RunnersCreated += 1
			case runner.ComplianceRunnerQueued:
				status.RunnersQueued += 1
			case runner.ComplianceRunnerInProgress:
				status.RunnersRunning += 1
			case runner.ComplianceRunnerFailed:
				status.RunnersFailed += 1
			case runner.ComplianceRunnerSucceeded:
				status.RunnersSucceeded += 1
			case runner.ComplianceRunnerTimeOut:
				status.RunnersTimedOut += 1
			}
		}
		statusJson, err := json.Marshal(status)
		if err != nil {
			return err
		}

		jp := pgtype.JSONB{}
		err = jp.Set(statusJson)
		if err != nil {
			return err
		}

		err = s.db.UpdateComplianceJobRunnersStatus(complianceJob.ID, jp)
		if err != nil {
			return err
		}
	}
	return nil
}
