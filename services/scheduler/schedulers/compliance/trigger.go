package compliance

import (
	"strings"
	"time"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"
	complianceApi "github.com/opengovern/opencomply/services/compliance/api"
	integrationapi "github.com/opengovern/opencomply/services/integration/api/models"

	runner "github.com/opengovern/opencomply/jobs/compliance-runner-job"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.uber.org/zap"
)

func (s *JobScheduler) buildRunners(
	parentJobID uint,
	integrationID *string,
	integrationType *integration.Type,
	resourceCollectionID *string,
	rootBenchmarkID string,
	parentBenchmarkIDs []string,
	benchmarkID string,
	currentRunnerExistMap map[string]bool,
	triggerType model.ComplianceTriggerType,
) ([]*model.ComplianceRunner, []*model.ComplianceRunner, error) {
	ctx := &httpclient.Context{UserRole: api.AdminRole}
	var runners []*model.ComplianceRunner
	var globalRunners []*model.ComplianceRunner

	benchmark, err := s.complianceClient.GetBenchmark(ctx, benchmarkID)
	if err != nil {
		s.logger.Error("error while getting benchmark", zap.Error(err), zap.String("benchmarkID", benchmarkID))
		return nil, nil, err
	}
	if currentRunnerExistMap == nil {
		currentRunners, err := s.db.GetRunnersByParentJobID(parentJobID)
		if err != nil {
			s.logger.Error("error while getting current runners", zap.Error(err))
			return nil, nil, err
		}
		currentRunnerExistMap = make(map[string]bool)
		for _, r := range currentRunners {
			currentRunnerExistMap[r.GetKeyIdentifier()] = true
		}
	}

	for _, child := range benchmark.Children {
		childRunners, childGlobalRunners, err := s.buildRunners(parentJobID, integrationID, integrationType, resourceCollectionID, rootBenchmarkID, append(parentBenchmarkIDs, benchmarkID), child, currentRunnerExistMap, triggerType)
		if err != nil {
			s.logger.Error("error while building child runners", zap.Error(err))
			return nil, nil, err
		}

		runners = append(runners, childRunners...)
		globalRunners = append(globalRunners, childGlobalRunners...)
	}

	for _, controlID := range benchmark.Controls {
		control, err := s.complianceClient.GetControl(ctx, controlID)
		if err != nil {
			s.logger.Error("error while getting control", zap.Error(err), zap.String("controlID", controlID))
			return nil, nil, err
		}

		if control.Policy == nil {
			continue
		}
		if integrationType != nil && len(control.Policy.IntegrationType) > 0 {
			supportsConnector := false
			for _, c := range control.Policy.IntegrationType {
				if strings.ToLower(
					integrationType.String()) == strings.ToLower(c) {
					supportsConnector = true
					break
				}
			}
			if !supportsConnector {
				continue
			}
		}

		callers := runner.Caller{
			RootBenchmark:      rootBenchmarkID,
			TracksDriftEvents:  benchmark.TracksDriftEvents,
			ParentBenchmarkIDs: append(parentBenchmarkIDs, benchmarkID),
			ControlID:          control.ID,
			ControlSeverity:    control.Severity,
		}

		runnerJob := model.ComplianceRunner{
			FrameworkID:          rootBenchmarkID,
			PolicyID:             control.Policy.ID,
			ControlID:            control.ID,
			IntegrationID:        integrationID,
			ResourceCollectionID: resourceCollectionID,
			ParentJobID:          parentJobID,
			StartedAt:            time.Time{},
			RetryCount:           0,
			Status:               runner.ComplianceRunnerCreated,
			FailureMessage:       "",
			TriggerType:          triggerType,
		}
		err = runnerJob.SetCallers([]runner.Caller{callers})
		if err != nil {
			return nil, nil, err
		}
		runners = append(runners, &runnerJob)
	}

	uniqueMap := map[string]*model.ComplianceRunner{}
	for _, r := range runners {
		v, ok := uniqueMap[r.PolicyID]
		if ok {
			cr, err := r.GetCallers()
			if err != nil {
				s.logger.Error("error while getting callers", zap.Error(err))
				return nil, nil, err
			}

			cv, err := v.GetCallers()
			if err != nil {
				s.logger.Error("error while getting callers", zap.Error(err))
				return nil, nil, err
			}

			cv = append(cv, cr...)
			err = v.SetCallers(cv)
			if err != nil {
				s.logger.Error("error while setting callers", zap.Error(err))
				return nil, nil, err
			}
		} else {
			v = r
		}
		uniqueMap[r.PolicyID] = v
	}
	globalUniqueMap := map[string]*model.ComplianceRunner{}
	for _, r := range globalRunners {
		v, ok := globalUniqueMap[r.PolicyID]
		if ok {
			cr, err := r.GetCallers()
			if err != nil {
				s.logger.Error("error while getting callers", zap.Error(err))
				return nil, nil, err
			}

			cv, err := v.GetCallers()
			if err != nil {
				s.logger.Error("error while getting callers", zap.Error(err))
				return nil, nil, err
			}

			cv = append(cv, cr...)
			err = v.SetCallers(cv)
			if err != nil {
				s.logger.Error("error while setting callers", zap.Error(err))
				return nil, nil, err
			}
		} else {
			v = r
		}
		globalUniqueMap[r.PolicyID] = v
	}

	var jobs []*model.ComplianceRunner
	var globalJobs []*model.ComplianceRunner
	for _, v := range uniqueMap {
		if !currentRunnerExistMap[v.GetKeyIdentifier()] {
			jobs = append(jobs, v)
		}
	}
	for _, v := range globalUniqueMap {
		if !currentRunnerExistMap[v.GetKeyIdentifier()] {
			globalJobs = append(globalJobs, v)
		}
	}
	return jobs, globalJobs, nil
}

func (s *JobScheduler) CreateComplianceReportJobs(withIncident bool, frameworkID string,
	lastJob *model.ComplianceJob, integrationIDs []string, manual bool, createdBy string, parentJobID *uint) ([]model.ComplianceJob, error) {
	// delete old runners
	if lastJob != nil {
		err := s.db.DeleteOldRunnerJob(&lastJob.ID)
		if err != nil {
			s.logger.Error("error while deleting old runners", zap.Error(err))
			return nil, err
		}
	} else {
		err := s.db.DeleteOldRunnerJob(nil)
		if err != nil {
			s.logger.Error("error while deleting old runners", zap.Error(err))
			return nil, err
		}
	}
	triggerType := model.ComplianceTriggerTypeScheduled
	if manual {
		triggerType = model.ComplianceTriggerTypeManual
	}

	var jobs []model.ComplianceJob
	var integrationsEpoch []string

	for _, integrationID := range integrationIDs {
		integrationsEpoch = append(integrationsEpoch, integrationID)
		if len(integrationsEpoch) >= 10 {
			job := model.ComplianceJob{
				FrameworkID:         frameworkID,
				WithIncidents:       withIncident,
				Status:              model.ComplianceJobCreated,
				AreAllRunnersQueued: false,
				IntegrationIDs:      integrationsEpoch,
				TriggerType:         triggerType,
				CreatedBy:           createdBy,
				ParentID:            parentJobID,
			}
			err := s.db.CreateComplianceJob(nil, &job)
			if err != nil {
				s.logger.Error("error while creating compliance job", zap.Error(err))
				return nil, err
			}
			jobs = append(jobs, job)
			integrationsEpoch = []string{}
		}
	}
	if len(integrationsEpoch) > 0 {
		job := model.ComplianceJob{
			FrameworkID:         frameworkID,
			WithIncidents:       withIncident,
			Status:              model.ComplianceJobCreated,
			AreAllRunnersQueued: false,
			IntegrationIDs:      integrationsEpoch,
			TriggerType:         triggerType,
			CreatedBy:           createdBy,
			ParentID:            parentJobID,
		}
		err := s.db.CreateComplianceJob(nil, &job)
		if err != nil {
			s.logger.Error("error while creating compliance job", zap.Error(err))
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *JobScheduler) enqueueRunnersCycle() error {
	s.logger.Info("enqueue runners cycle started")
	var err error
	jobsWithUnqueuedRunners, err := s.db.ListComplianceJobsWithUnqueuedRunners(true)
	if err != nil {
		s.logger.Error("error while listing jobs with unqueued runners", zap.Error(err))
		return err
	}
	s.logger.Info("jobs with unqueued runners", zap.Int("count", len(jobsWithUnqueuedRunners)))
	for _, job := range jobsWithUnqueuedRunners {
		s.logger.Info("processing job with unqueued runners", zap.Uint("jobID", job.ID))
		var allRunners []*model.ComplianceRunner
		var assignments *complianceApi.BenchmarkAssignedEntities
		integrations, err := s.integrationClient.ListIntegrationsByFilters(&httpclient.Context{UserRole: api.AdminRole}, integrationapi.ListIntegrationsRequest{
			IntegrationID: job.IntegrationIDs,
		})
		if err != nil {
			s.logger.Error("error while getting integrations", zap.Error(err))
			continue
		}
		assignments = &complianceApi.BenchmarkAssignedEntities{}
		for _, integration := range integrations.Integrations {
			assignment := complianceApi.BenchmarkAssignedIntegration{
				IntegrationID:   integration.IntegrationID,
				ProviderID:      integration.ProviderID,
				IntegrationName: integration.Name,
				IntegrationType: integration.IntegrationType,
				Status:          true,
			}
			assignments.Integrations = append(assignments.Integrations, assignment)
		}

		var globalRunners []*model.ComplianceRunner
		var runners []*model.ComplianceRunner
		for _, it := range assignments.Integrations {
			if !it.Status {
				continue
			}
			connection := it
			runners, globalRunners, err = s.buildRunners(job.ID, &connection.IntegrationID, &connection.IntegrationType, nil, job.FrameworkID, nil, job.FrameworkID, nil, job.TriggerType)
			if err != nil {
				s.logger.Error("error while building runners", zap.Error(err))
				return err
			}
			allRunners = append(allRunners, runners...)
		}
		allRunners = append(allRunners, globalRunners...)
		if len(allRunners) > 0 {
			s.logger.Info("creating runners", zap.Int("count", len(allRunners)), zap.Uint("jobID", job.ID))
			err = s.db.CreateRunnerJobs(nil, allRunners)
			if err != nil {
				s.logger.Error("error while creating runners", zap.Error(err))
				return err
			}
		} else {
			s.logger.Info("no runners to create", zap.Uint("jobID", job.ID))
			err = s.db.UpdateComplianceJobAreAllRunnersQueued(job.ID, true)
			if err != nil {
				s.logger.Error("error while updating compliance job", zap.Error(err))
				return err
			}
		}
	}

	return nil
}
