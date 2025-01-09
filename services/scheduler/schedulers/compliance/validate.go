package compliance

import (
	"fmt"
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/compliance/api"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"go.uber.org/zap"
)

func (s *JobScheduler) ValidateComplianceJob(framework api.Benchmark) error {
	listOfTables, err := s.getTablesUnderBenchmark(framework, make(map[string]FrameworkTablesCache))
	if err != nil {
		return err
	}

	var integrationTypes []interfaces.IntegrationType
	for _, itName := range framework.IntegrationTypes {
		if it, ok := integration_type.IntegrationTypes[integration.Type(itName)]; ok {
			integrationTypes = append(integrationTypes, it)
		} else {
			return fmt.Errorf("integration type not valid: %s", itName)
		}
	}

	for table := range listOfTables {
		for _, it := range integrationTypes {
			if rt := it.GetResourceTypeFromTableName(table); rt == "" {
				return fmt.Errorf("invalid table name: %s", table)
			}
		}
	}

	listOfParameters, err := s.getParametersUnderFramework(framework, make(map[string]FrameworkParametersCache))
	if err != nil {
		return err
	}
	queryParams, err := s.coreClient.ListQueryParameters(&httpclient.Context{UserRole: api2.AdminRole})
	if err != nil {
		s.logger.Error("failed to get query parameters", zap.Error(err))
		return err
	}
	queryParamMap := make(map[string]string)
	for _, qp := range queryParams.Items {
		if qp.Value != "" {
			queryParamMap[qp.Key] = qp.Value
		}
	}

	for param := range listOfParameters {
		if _, ok := queryParamMap[param]; !ok {
			return fmt.Errorf("query parameter %s not exists", param)
		}
	}
	return nil
}

type FrameworkTablesCache struct {
	ListTables map[string]bool
}

type FrameworkParametersCache struct {
	ListParameters map[string]bool
}

// getTablesUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (s *JobScheduler) getTablesUnderBenchmark(framework api.Benchmark, benchmarkCache map[string]FrameworkTablesCache) (map[string]bool, error) {
	listOfTables := make(map[string]bool)

	controls, err := s.complianceClient.ListControl(&httpclient.Context{UserRole: api2.AdminRole}, framework.Controls, nil)
	if err != nil {
		s.logger.Error("failed to fetch controls", zap.Error(err))
		return nil, err
	}
	for _, c := range controls {
		if c.Policy != nil {
			for _, t := range c.Policy.ListOfResources {
				if t == "" {
					continue
				}
				listOfTables[t] = true
			}
		}
	}

	children, err := s.complianceClient.ListBenchmarks(&httpclient.Context{UserRole: api2.AdminRole}, framework.Children, nil)
	if err != nil {
		s.logger.Error("failed to fetch children", zap.Error(err))
		return nil, err
	}
	for _, child := range children {
		var childListOfTables map[string]bool
		if cache, ok := benchmarkCache[child.ID]; ok {
			childListOfTables = cache.ListTables
		} else {
			childListOfTables, err = s.getTablesUnderBenchmark(child, benchmarkCache)
			if err != nil {
				return nil, err
			}
			benchmarkCache[child.ID] = FrameworkTablesCache{
				ListTables: childListOfTables,
			}
		}

		for k, _ := range childListOfTables {
			childListOfTables[k] = true
		}
	}
	return listOfTables, nil
}

// getParametersUnderFramework ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (s *JobScheduler) getParametersUnderFramework(framework api.Benchmark, frameworkCache map[string]FrameworkParametersCache) (map[string]bool, error) {
	listOfParameters := make(map[string]bool)

	controls, err := s.complianceClient.ListControl(&httpclient.Context{UserRole: api2.AdminRole}, framework.Controls, nil)
	if err != nil {
		s.logger.Error("failed to fetch controls", zap.Error(err))
		return nil, err
	}
	for _, c := range controls {
		if c.Policy != nil {
			for _, t := range c.Policy.Parameters {
				listOfParameters[t.Key] = true
			}
		}
	}

	children, err := s.complianceClient.ListBenchmarks(&httpclient.Context{UserRole: api2.AdminRole}, framework.Children, nil)
	if err != nil {
		s.logger.Error("failed to fetch children", zap.Error(err))
		return nil, err
	}
	for _, child := range children {
		var childListOfParameters map[string]bool
		if cache, ok := frameworkCache[child.ID]; ok {
			childListOfParameters = cache.ListParameters
		} else {
			childListOfParameters, err = s.getParametersUnderFramework(child, frameworkCache)
			if err != nil {
				return nil, err
			}
			frameworkCache[child.ID] = FrameworkParametersCache{
				ListParameters: childListOfParameters,
			}
		}

		for k, _ := range childListOfParameters {
			childListOfParameters[k] = true
		}
	}
	return listOfParameters, nil
}
