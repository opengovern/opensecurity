package compliance

import (
	"fmt"
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/compliance/api"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.uber.org/zap"
)

func (s *JobScheduler) validateComplianceJob(framework api.Benchmark) error {
	validation, err := s.db.GetFrameworkValidation(framework.ID)
	if validation == nil {
		listOfTables, err := s.getTablesUnderBenchmark(framework, make(map[string]FrameworkTablesCache))
		if err != nil {
			_ = s.db.CreateFrameworkValidation(&model.FrameworkValidation{
				FrameworkID:    framework.ID,
				FailureMessage: err.Error(),
			})
			return err
		}

		var integrationTypes []interfaces.IntegrationType
		for _, itName := range framework.IntegrationTypes {
			if it, ok := integration_type.IntegrationTypes[integration.Type(itName)]; ok {
				integrationTypes = append(integrationTypes, it)
			} else {
				_ = s.db.CreateFrameworkValidation(&model.FrameworkValidation{
					FrameworkID:    framework.ID,
					FailureMessage: fmt.Errorf("integration type not valid: %s", itName).Error(),
				})
				return fmt.Errorf("integration type not valid: %s", itName)
			}
		}

		tablesMap := make(map[string]struct{})
		for _, it := range integrationTypes {
			tables, err := it.GetTablesByLabels(nil)
			if err != nil {
				_ = s.db.CreateFrameworkValidation(&model.FrameworkValidation{
					FrameworkID:    framework.ID,
					FailureMessage: err.Error(),
				})
				return err
			}
			for _, table := range tables {
				tablesMap[table] = struct{}{}
			}
		}

		for table := range listOfTables {
			if _, ok := tablesMap[table]; !ok {
				_ = s.db.CreateFrameworkValidation(&model.FrameworkValidation{
					FrameworkID:    framework.ID,
					FailureMessage: fmt.Sprintf("table %s not exist", table),
				})
				return fmt.Errorf("table %s not exist", table)
			}
		}

		_ = s.db.CreateFrameworkValidation(&model.FrameworkValidation{
			FrameworkID:    framework.ID,
			FailureMessage: "",
		})
	} else if validation.FailureMessage != "" {
		return fmt.Errorf("framework %s has failed validation: %s", framework.ID, validation.FailureMessage)
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
	ctx := &httpclient.Context{UserRole: api2.AdminRole}
	listOfTables := make(map[string]bool)

	controlIDsMap, err := s.getControlsUnderBenchmark(framework)
	if err != nil {
		s.logger.Error("failed to fetch controls", zap.Error(err))
		return nil, err
	}
	var controlIDs []string
	for controlID := range controlIDsMap {
		controlIDs = append(controlIDs, controlID)
	}

	controls, err := s.complianceClient.ListControl(ctx, controlIDs, nil)
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

	return listOfTables, nil
}

func (s *JobScheduler) getControlsUnderBenchmark(framework api.Benchmark) (map[string]bool, error) {
	ctx := &httpclient.Context{UserRole: api2.AdminRole}

	s.logger.Info("getting framework children", zap.String("framework_id", framework.ID), zap.Strings("children", framework.Children))

	controls := make(map[string]bool)
	for _, c := range framework.Controls {
		controls[c] = true
	}
	if len(framework.Children) > 0 {
		children, err := s.complianceClient.ListBenchmarks(ctx, framework.Children, nil)
		if err != nil {
			s.logger.Error("failed to fetch children", zap.Error(err))
			return nil, err
		}
		for _, child := range children {
			childControls, err := s.getControlsUnderBenchmark(child)
			if err != nil {
				s.logger.Error("failed to fetch controls", zap.Error(err))
				return nil, err
			}
			for k, _ := range childControls {
				controls[k] = true
			}
		}
	}

	s.logger.Info("got framework controls", zap.Any("controls", controls))
	return controls, nil
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
