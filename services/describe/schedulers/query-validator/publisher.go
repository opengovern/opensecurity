package query_validator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	queryvalidator "github.com/opengovern/opencomply/jobs/query-validator-job"
	coreApi "github.com/opengovern/opencomply/services/core/api"
	"go.uber.org/zap"
)

func (s *JobScheduler) runPublisher(ctx context.Context) error {
	ctx2 := &httpclient.Context{UserRole: api.AdminRole}
	ctx2.Ctx = ctx

	s.logger.Info("Query Runner publisher started")

	err := s.db.UpdateTimedOutQueuedQueryRunners()
	if err != nil {
		s.logger.Error("failed to update timed out query runners", zap.Error(err))
	}

	err = s.db.UpdateTimedOutInProgressQueryRunners()
	if err != nil {
		s.logger.Error("failed to update timed out query runners", zap.Error(err))
	}

	count, err := s.db.GetInProgressJobsCount()
	if err != nil {
		s.logger.Error("GetInProgressJobsCount Error", zap.Error(err))
		return err
	}
	jobs, err := s.db.FetchCreatedQueryValidatorJobs(200 - count)
	if err != nil {
		s.logger.Error("List Queries Error", zap.Error(err))
		return err
	}
	for _, job := range jobs {
		jobMsg := &queryvalidator.Job{
			ID: job.ID,
		}
		if job.QueryType == queryvalidator.QueryTypeNamedQuery {
			jobMsg.QueryType = queryvalidator.QueryTypeNamedQuery
			jobMsg.QueryId = job.QueryId
			namedQuery, err := s.coreClient.GetQuery(ctx2, job.QueryId)
			if err != nil {
				s.logger.Error("Get Query Error", zap.Error(err))
			}
			jobMsg.Query = namedQuery.Query.QueryToExecute
			jobMsg.Parameters = namedQuery.Query.Parameters
			jobMsg.ListOfTables = namedQuery.Query.ListOfTables
			jobMsg.PrimaryTable = namedQuery.Query.PrimaryTable
			jobMsg.IntegrationType = namedQuery.IntegrationTypes
		} else if job.QueryType == queryvalidator.QueryTypeComplianceControl {
			jobMsg.QueryType = queryvalidator.QueryTypeComplianceControl
			jobMsg.QueryId = job.QueryId
			controlQuery, err := s.complianceClient.GetControlDetails(ctx2, job.QueryId)
			if err != nil {
				s.logger.Error("Get Control Error", zap.Error(err))
			}
			jobMsg.Query = controlQuery.Query.QueryToExecute
			var parameters []coreApi.QueryParameter
			for _, qp := range controlQuery.Query.Parameters {
				parameters = append(parameters, coreApi.QueryParameter{
					Key:      qp.Key,
					Required: qp.Required,
				})
			}
			jobMsg.Parameters = parameters
			jobMsg.ListOfTables = controlQuery.Query.ListOfTables
			jobMsg.PrimaryTable = controlQuery.Query.PrimaryTable
			jobMsg.IntegrationType = controlQuery.IntegrationType
		} else {
			_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, "query ID not found")
			continue
		}

		queryParams, err := s.coreClient.ListQueryParameters(&httpclient.Context{UserRole: api.AdminRole})
		if err != nil {
			_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, fmt.Sprintf("failed to list parameters: %s", err.Error()))
			return err
		}
		queryParamMap := make(map[string]string)
		for _, qp := range queryParams.Items {
			queryParamMap[qp.Key] = qp.Value
		}
		queryTemplate, err := template.New(jobMsg.QueryId).Parse(jobMsg.Query)
		if err != nil {
			return err
		}
		var queryOutput bytes.Buffer
		if err := queryTemplate.Execute(&queryOutput, queryParamMap); err != nil {
			_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, fmt.Sprintf("failed to execute query template: %s", err.Error()))
			return fmt.Errorf("failed to execute query template: %w", err)
		}

		jobMsg.Query = queryOutput.String()

		jobJson, err := json.Marshal(jobMsg)
		if err != nil {
			_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, "failed to marshal job")
			s.logger.Error("failed to marshal Query Runner Job", zap.Error(err), zap.Uint("runnerId", job.ID))
			continue
		}

		s.logger.Info("publishing query runner job", zap.Uint("jobId", job.ID))
		topic := queryvalidator.JobQueueTopic
		_, err = s.jq.Produce(ctx, topic, jobJson, fmt.Sprintf("job-%d", job.ID))
		if err != nil {
			if err.Error() == "nats: no response from stream" {
				err = s.runSetupNatsStreams(ctx)
				if err != nil {
					s.logger.Error("Failed to setup nats streams", zap.Error(err))
					return err
				}
				_, err = s.jq.Produce(ctx, topic, jobJson, fmt.Sprintf("job-%d", job.ID))
				if err != nil {
					_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, err.Error())
					s.logger.Error("failed to send job", zap.Error(err), zap.Uint("runnerId", job.ID))
					continue
				}
			} else {
				_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorFailed, err.Error())
				s.logger.Error("failed to send query runner job", zap.Error(err), zap.Uint("runnerId", job.ID), zap.String("error message", err.Error()))
				continue
			}
		}

		_ = s.db.UpdateQueryValidatorJobStatus(job.ID, queryvalidator.QueryValidatorQueued, "")
	}
	return nil
}
