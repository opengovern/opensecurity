package compliance_quick_run_job

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"go.uber.org/zap"
)

type AuditJob struct {
	JobID          uint
	FrameworkID    string
	IntegrationIDs []string
	IncludeResult  []string

	JobReportControlSummary *types.ComplianceJobReportControlSummary
	JobReportControlView    *types.ComplianceJobReportControlView
	JobReportResourceView   *types.ComplianceJobReportResourceView
}

type JobResult struct {
	JobID          uint
	Status         model.ComplianceJobStatus
	FailureMessage string
}

func (w *Worker) RunJob(ctx context.Context, job *AuditJob) error {
	job.JobReportControlView = &types.ComplianceJobReportControlView{
		Controls:          make(map[string]types.AuditControlResult),
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		JobSummary: types.JobSummary{
			JobID:          job.JobID,
			FrameworkID:    job.FrameworkID,
			Auditable:      false,
			JobStartedAt:   time.Now(),
			IntegrationIDs: []string{},
		},
	}
	job.JobReportControlSummary = &types.ComplianceJobReportControlSummary{
		Controls:          make(map[string]*types.ControlSummary),
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		ControlScore: &types.ControlScore{
			TotalControls:  0,
			FailedControls: 0,
		},
		JobSummary: types.JobSummary{
			JobID:          job.JobID,
			FrameworkID:    job.FrameworkID,
			Auditable:      false,
			JobStartedAt:   time.Now(),
			IntegrationIDs: []string{},
		},
	}
	job.JobReportResourceView = &types.ComplianceJobReportResourceView{
		Integrations:      make(map[string]types.AuditIntegrationResult),
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		JobSummary: types.JobSummary{
			JobID:          job.JobID,
			FrameworkID:    job.FrameworkID,
			Auditable:      false,
			JobStartedAt:   time.Now(),
			IntegrationIDs: []string{},
		},
	}

	totalControls := make(map[string]bool)
	failedControls := make(map[string]bool)

	for _, integrationID := range job.IntegrationIDs {
		err := w.RunJobForIntegration(ctx, job, integrationID, &totalControls, &failedControls)
		if err != nil {
			w.logger.Error("failed to run audit job for integration", zap.String("integration_id", integrationID), zap.Error(err))
			return err
		}
		w.logger.Info("audit job for integration completed", zap.String("integration_id", integrationID))
	}

	keys, idx := job.JobReportControlView.KeysAndIndex()
	job.JobReportControlView.EsID = es.HashOf(keys...)
	job.JobReportControlView.EsIndex = idx

	err := sendDataToOpensearch(w.esClient.ES(), *job.JobReportControlView)
	if err != nil {
		return err
	}

	keys, idx = job.JobReportResourceView.KeysAndIndex()
	job.JobReportResourceView.EsID = es.HashOf(keys...)
	job.JobReportResourceView.EsIndex = idx

	err = sendDataToOpensearch(w.esClient.ES(), *job.JobReportResourceView)
	if err != nil {
		return err
	}

	job.JobReportControlSummary.ControlScore.FailedControls = int64(len(failedControls))
	job.JobReportControlSummary.ControlScore.TotalControls = int64(len(totalControls))
	keys, idx = job.JobReportControlSummary.KeysAndIndex()
	job.JobReportControlSummary.EsID = es.HashOf(keys...)
	job.JobReportControlSummary.EsIndex = idx

	err = sendDataToOpensearch(w.esClient.ES(), *job.JobReportControlSummary)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) RunJobForIntegration(ctx context.Context, job *AuditJob, integrationId string, totalControls, failedControls *map[string]bool) error {
	include := make(map[string]bool)
	if len(job.IncludeResult) > 0 {
		for _, result := range job.IncludeResult {
			include[result] = true
		}
	} else {
		include[string(types.ComplianceStatusOK)] = true
		include[string(types.ComplianceStatusALARM)] = true
	}

	job.JobReportControlView.JobSummary.IntegrationIDs = append(job.JobReportControlView.JobSummary.IntegrationIDs, integrationId)
	job.JobReportResourceView.JobSummary.IntegrationIDs = append(job.JobReportResourceView.JobSummary.IntegrationIDs, integrationId)
	job.JobReportControlSummary.JobSummary.IntegrationIDs = append(job.JobReportControlSummary.JobSummary.IntegrationIDs, integrationId)

	job.JobReportResourceView.Integrations[integrationId] = types.AuditIntegrationResult{
		ResourceTypes: make(map[string]types.AuditResourceTypesResult),
	}

	err := w.Initialize(ctx, integrationId)
	if err != nil {
		return err
	}

	defer w.steampipeConn.UnsetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyIntegrationID)
	defer w.steampipeConn.UnsetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyClientType)

	return w.runJobForFramework(ctx, job, integrationId, job.FrameworkID, totalControls, failedControls, include)
}

func (w *Worker) runJobForFramework(ctx context.Context, job *AuditJob, integrationId string, frameworkId string,
	totalControls, failedControls *map[string]bool, include map[string]bool) error {
	ctx2 := &httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole}

	framework, err := w.complianceClient.GetBenchmark(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole}, frameworkId)
	if err != nil {
		return err
	}
	controls, err := w.complianceClient.ListControl(ctx2, framework.Controls, nil)
	if err != nil {
		return err
	}

	for _, control := range controls {
		if control.Policy == nil {
			continue
		}
		queryJob := QueryJob{
			AuditJobID: job.JobID,
			ExecutionPlan: ExecutionPlan{
				Policy:         *control.Policy,
				ControlID:      control.ID,
				IntegrationIDs: job.IntegrationIDs,
			},
		}
		queryResults, err := w.RunQuery(ctx, queryJob)
		if err != nil {
			w.logger.Error("failed to run query", zap.String("jobID", strconv.Itoa(int(job.JobID))),
				zap.String("frameworkID", job.FrameworkID), zap.String("integrationID", integrationId),
				zap.String("controlID", control.ID), zap.Error(err))
			continue
		}
		for _, qr := range queryResults {
			(*totalControls)[control.ID] = true
			if qr.ComplianceStatus == types.ComplianceStatusALARM {
				(*failedControls)[control.ID] = true
			}
			if _, ok := include[string(qr.ComplianceStatus)]; !ok {
				continue
			}
			if _, ok := job.JobReportResourceView.ComplianceSummary[qr.ComplianceStatus]; !ok {
				job.JobReportResourceView.ComplianceSummary[qr.ComplianceStatus] = 0
			}
			job.JobReportResourceView.ComplianceSummary[qr.ComplianceStatus] += 1
			if _, ok := job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType]; !ok {
				job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType] = types.AuditResourceTypesResult{
					Resources: make(map[string]types.AuditResourceResult),
				}
			}
			if _, ok := job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID]; !ok {
				job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID] = types.AuditResourceResult{
					ResourceSummary: make(map[types.ComplianceStatus]uint64),
					Results:         make(map[types.ComplianceStatus][]types.AuditControlFinding),
					ResourceName:    qr.ResourceName,
				}
			}
			if _, ok := job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].ResourceSummary[qr.ComplianceStatus]; !ok {
				job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].ResourceSummary[qr.ComplianceStatus] = 0
			}
			job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].ResourceSummary[qr.ComplianceStatus] += 1
			if _, ok := job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].Results[qr.ComplianceStatus]; !ok {
				job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].Results[qr.ComplianceStatus] = make([]types.AuditControlFinding, 0)
			}
			job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].Results[qr.ComplianceStatus] = append(
				job.JobReportResourceView.Integrations[integrationId].ResourceTypes[qr.ResourceType].Resources[qr.ResourceID].Results[qr.ComplianceStatus], types.AuditControlFinding{
					Severity:  control.Severity,
					ControlID: control.ID,
					Reason:    qr.Reason,
				})

			// Audit Summary
			if _, ok := job.JobReportControlView.Controls[control.ID]; !ok {
				job.JobReportControlView.Controls[control.ID] = types.AuditControlResult{
					Severity:       control.Severity,
					ControlSummary: make(map[types.ComplianceStatus]uint64),
					Results:        make(map[types.ComplianceStatus][]types.AuditResourceFinding),
				}
			}
			if _, ok := job.JobReportControlView.ComplianceSummary[qr.ComplianceStatus]; !ok {
				job.JobReportControlView.ComplianceSummary[qr.ComplianceStatus] = 0
			}
			job.JobReportControlView.ComplianceSummary[qr.ComplianceStatus] += 1

			if _, ok := job.JobReportControlView.Controls[control.ID].ControlSummary[qr.ComplianceStatus]; !ok {
				job.JobReportControlView.Controls[control.ID].ControlSummary[qr.ComplianceStatus] = 0
			}
			if _, ok := job.JobReportControlView.Controls[control.ID].Results[qr.ComplianceStatus]; !ok {
				job.JobReportControlView.Controls[control.ID].Results[qr.ComplianceStatus] = make([]types.AuditResourceFinding, 0)
			}
			job.JobReportControlView.Controls[control.ID].ControlSummary[qr.ComplianceStatus] += 1
			job.JobReportControlView.Controls[control.ID].Results[qr.ComplianceStatus] = append(job.JobReportControlView.Controls[control.ID].Results[qr.ComplianceStatus],
				types.AuditResourceFinding{
					ResourceID:   qr.ResourceID,
					ResourceType: qr.ResourceType,
					Reason:       qr.Reason,
				})

			if _, ok := job.JobReportControlSummary.ComplianceSummary[qr.ComplianceStatus]; !ok {
				job.JobReportControlSummary.ComplianceSummary[qr.ComplianceStatus] = 0
			}
			job.JobReportControlSummary.ComplianceSummary[qr.ComplianceStatus] += 1
			if v, ok := job.JobReportControlSummary.Controls[control.ID]; !ok || v == nil {
				job.JobReportControlSummary.Controls[control.ID] = &types.ControlSummary{
					Severity: control.Severity,
					Alarms:   0,
					Oks:      0,
				}
			}
			switch qr.ComplianceStatus {
			case types.ComplianceStatusALARM:
				job.JobReportControlSummary.Controls[control.ID].Alarms += 1
			case types.ComplianceStatusOK:
				job.JobReportControlSummary.Controls[control.ID].Oks += 1
			}
		}
	}

	for _, child := range framework.Children {
		err = w.runJobForFramework(ctx, job, integrationId, child, totalControls, failedControls, include)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) Initialize(ctx context.Context, integrationId string) error {
	err := w.steampipeConn.SetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyIntegrationID, integrationId)
	if err != nil {
		w.logger.Error("failed to set account id", zap.Error(err))
		return err
	}
	err = w.steampipeConn.SetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyClientType, "compliance")
	if err != nil {
		w.logger.Error("failed to set client type", zap.Error(err))
		return err
	}

	return nil
}

func sendDataToOpensearch(client *opensearch.Client, doc es.Doc) error {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	keys, index := doc.KeysAndIndex()

	// Use the opensearchapi.IndexRequest to index the document
	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: es.HashOf(keys...),
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true", // Makes the document immediately available for search
	}
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}
