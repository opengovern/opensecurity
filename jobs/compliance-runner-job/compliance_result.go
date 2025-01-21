package runner

import (
	"context"
	"fmt"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opencomply/services/integration/client"
	"reflect"
	"strconv"
	"strings"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/pkg/utils"
	"github.com/opengovern/opencomply/services/compliance/api"
	"go.uber.org/zap"
)

func (w *Job) GetResourceTypeFromTableName(integrationClient client.IntegrationServiceClient, tableName string, queryIntegrationType []string) (string, integration.Type, error) {
	var integrationType string
	if len(queryIntegrationType) == 1 {
		integrationType = queryIntegrationType[0]
	} else {
		integrationType = ""
	}

	tableResourceType, err := integrationClient.GetResourceTypeFromTableName(&httpclient.Context{Ctx: context.Background(), UserRole: authApi.ViewerRole}, integrationType, tableName)
	if err != nil {
		return "", "", err
	}

	return tableResourceType, integration.Type(integrationType), nil
}

func (w *Job) ExtractComplianceResults(logger *zap.Logger, benchmarkCache map[string]api.Benchmark, integrationClient client.IntegrationServiceClient, caller Caller, res *steampipe.Result, query api.Policy) ([]types.ComplianceResult, error) {
	var complianceResults []types.ComplianceResult
	var integrationType integration.Type
	var err error
	queryResourceType := ""
	if query.PrimaryResource != nil || len(query.ListOfResources) == 1 {
		tableName := ""
		if query.PrimaryResource != nil {
			tableName = *query.PrimaryResource
		} else {
			tableName = query.ListOfResources[0]
		}
		if tableName != "" {
			queryResourceType, integrationType, err = w.GetResourceTypeFromTableName(integrationClient, tableName, w.ExecutionPlan.Query.IntegrationType)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, record := range res.Data {
		if len(record) != len(res.Headers) {
			return nil, fmt.Errorf("invalid record length, record=%d headers=%d", len(record), len(res.Headers))
		}
		recordValue := make(map[string]any)
		for idx, header := range res.Headers {
			value := record[idx]
			recordValue[header] = value
		}
		resourceType := queryResourceType

		var platformResourceID, integrationID, resourceID, resourceName, reason string
		var costImpact *float64
		var status types.ComplianceStatus
		if v, ok := recordValue["platform_resource_id"].(string); ok {
			platformResourceID = v
		}
		if v, ok := recordValue["platform_account_id"].(string); ok {
			integrationID = v
		}
		if v, ok := recordValue["platform_table_name"].(string); ok && resourceType == "" {
			resourceType, integrationType, err = w.GetResourceTypeFromTableName(integrationClient, v, w.ExecutionPlan.Query.IntegrationType)
			if err != nil {
				return nil, err
			}
		}
		if v, ok := recordValue["resource"].(string); ok && v != "" && v != "null" {
			resourceID = v
		} else {
			continue
		}
		if v, ok := recordValue["name"].(string); ok {
			resourceName = v
		}
		if v, ok := recordValue["reason"].(string); ok {
			reason = v
		}
		if v, ok := recordValue["status"].(string); ok {
			status = types.ComplianceStatus(v)
		}
		if v, ok := recordValue["cost_optimization"]; ok {
			// cast to proper types
			reflectValue := reflect.ValueOf(v)
			switch reflectValue.Kind() {
			case reflect.Float32:
				costImpact = utils.GetPointer(float64(v.(float32)))
			case reflect.Float64:
				costImpact = utils.GetPointer(v.(float64))
			case reflect.String:
				c, err := strconv.ParseFloat(v.(string), 64)
				if err == nil {
					costImpact = &c
				} else {
					fmt.Printf("error parsing cost_optimization: %s\n", err)
					costImpact = utils.GetPointer(0.0)
				}
			case reflect.Int:
				costImpact = utils.GetPointer(float64(v.(int)))
			case reflect.Int8:
				costImpact = utils.GetPointer(float64(v.(int8)))
			case reflect.Int16:
				costImpact = utils.GetPointer(float64(v.(int16)))
			case reflect.Int32:
				costImpact = utils.GetPointer(float64(v.(int32)))
			case reflect.Int64:
				costImpact = utils.GetPointer(float64(v.(int64)))
			case reflect.Uint:
				costImpact = utils.GetPointer(float64(v.(uint)))
			case reflect.Uint8:
				costImpact = utils.GetPointer(float64(v.(uint8)))
			case reflect.Uint16:
				costImpact = utils.GetPointer(float64(v.(uint16)))
			case reflect.Uint32:
				costImpact = utils.GetPointer(float64(v.(uint32)))
			case reflect.Uint64:
				costImpact = utils.GetPointer(float64(v.(uint64)))
			default:
				fmt.Printf("error parsing cost_impact: unknown type %s\n", reflectValue.Kind())
			}
		}
		severity := caller.ControlSeverity
		if severity == "" {
			severity = types.ComplianceResultSeverityNone
		}

		if (integrationID == "" || integrationID == "null") && w.ExecutionPlan.IntegrationID != nil {
			integrationID = *w.ExecutionPlan.IntegrationID
		}

		benchmarkReferences := make([]string, 0, len([]string{caller.RootBenchmark}))
		for _, parentBenchmarkID := range []string{caller.RootBenchmark} {
			benchmarkReferences = append(benchmarkReferences, benchmarkCache[parentBenchmarkID].ReferenceCode)
		}

		if status != types.ComplianceStatusOK && status != types.ComplianceStatusALARM {
			continue
		}

		controlPath := strings.Join(append(benchmarkReferences, caller.ControlID), "/")

		complianceResults = append(complianceResults, types.ComplianceResult{
			BenchmarkID:        caller.RootBenchmark,
			ControlID:          caller.ControlID,
			IntegrationID:      integrationID,
			EvaluatedAt:        w.CreatedAt.UnixMilli(),
			StateActive:        true,
			ComplianceStatus:   status,
			Severity:           severity,
			IntegrationType:    integrationType,
			PlatformResourceID: platformResourceID,
			ResourceID:         resourceID,
			ResourceName:       resourceName,
			ResourceType:       resourceType,
			Reason:             reason,
			CostImpact:         costImpact,
			RunnerID:           w.ID,
			ComplianceJobID:    w.ParentJobID,
			ControlPath:        controlPath,
			ParentBenchmarks:   []string{caller.RootBenchmark},
			LastUpdatedAt:      w.CreatedAt.UnixMilli(),
		})
	}
	return complianceResults, nil
}
