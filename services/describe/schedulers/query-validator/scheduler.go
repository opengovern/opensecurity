package query_validator

import (
	"fmt"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	queryvalidator "github.com/opengovern/opencomply/jobs/query-validator-job"
	"github.com/opengovern/opencomply/services/describe/db/model"
	"go.uber.org/zap"
)

func (s *JobScheduler) runScheduler() error {
	clientCtx := &httpclient.Context{UserRole: api.AdminRole}

	controls, err := s.complianceClient.ListControl(clientCtx, nil, nil)
	if err != nil {
		s.logger.Error("error while listing benchmarks", zap.Error(err))
		return fmt.Errorf("error while listing benchmarks: %v", err)
	}
	for _, c := range controls {
		hasParams := false
		if len(c.Policy.Parameters) > 0 {
			hasParams = true
		}
		_, err = s.db.CreateQueryValidatorJob(&model.QueryValidatorJob{
			QueryId:        c.ID,
			QueryType:      queryvalidator.QueryTypeComplianceControl,
			Status:         queryvalidator.QueryValidatorCreated,
			HasParams:      hasParams,
			FailureMessage: "",
		})
		if err != nil {
			s.logger.Error("error while creating query-validator job", zap.Error(err))
		}
	}

	namedQueries, err := s.coreClient.ListQueriesV2(clientCtx, nil)
	if err != nil {
		s.logger.Error("error while listing benchmarks", zap.Error(err))
		return fmt.Errorf("error while listing benchmarks: %v", err)
	}
	for _, nq := range namedQueries.Items {
		hasParams := false
		if len(nq.Query.Parameters) > 0 {
			hasParams = true
		}
		_, err = s.db.CreateQueryValidatorJob(&model.QueryValidatorJob{
			QueryId:        nq.ID,
			QueryType:      queryvalidator.QueryTypeNamedQuery,
			Status:         queryvalidator.QueryValidatorCreated,
			HasParams:      hasParams,
			FailureMessage: "",
		})
		if err != nil {
			s.logger.Error("error while creating query-validator job", zap.Error(err))
		}
	}

	return nil
}
