package compliance

import (
	"fmt"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	integrationapi "github.com/opengovern/opencomply/services/integration/api/models"
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

	benchmarks, err := s.complianceClient.ListBenchmarks(clientCtx, nil)
	if err != nil {
		s.logger.Error("error while listing benchmarks", zap.Error(err))
		return fmt.Errorf("error while listing benchmarks: %v", err)
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

	for _, benchmark := range benchmarks {
		var integrationIDs []string
		assignments, err := s.complianceClient.ListAssignmentsByBenchmark(clientCtx, benchmark.ID)
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

		complianceJob, err := s.db.GetLastComplianceJob(true, benchmark.ID)
		if err != nil {
			s.logger.Error("error while getting last compliance job", zap.Error(err))
			return err
		}

		timeAt := time.Now().Add(-s.complianceIntervalHours)
		if complianceJob == nil ||
			complianceJob.CreatedAt.Before(timeAt) {

			_, err := s.CreateComplianceReportJobs(true, benchmark.ID, complianceJob, integrationIDs, false, "system", nil)
			if err != nil {
				s.logger.Error("error while creating compliance job", zap.Error(err))
				return err
			}

			ComplianceJobsCount.WithLabelValues("successful").Inc()
		}
	}

	return nil
}
