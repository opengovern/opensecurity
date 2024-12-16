package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"time"

	"github.com/opengovern/og-util/pkg/api"
	es2 "github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	types2 "github.com/opengovern/opencomply/jobs/compliance-summarizer-job/types"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/compliance/es"
	es3 "github.com/opengovern/opencomply/services/describe/es"
	"go.uber.org/zap"
)

func (w *Worker) RunJob(ctx context.Context, j types2.Job) error {
	w.logger.Info("Running summarizer",
		zap.Uint("job_id", j.ID),
		zap.String("benchmark_id", j.BenchmarkID),
	)

	// We have to sort by platformResourceID to be able to optimize memory usage for resourceFinding generations
	// this way as soon as paginator switches to next resource we can send the previous resource to the queue and free up memory
	paginator, err := es.NewComplianceResultPaginator(w.esClient, types.ComplianceResultsIndex, []opengovernance.BoolFilter{
		opengovernance.NewTermFilter("stateActive", "true"),
	}, nil, []map[string]any{
		{"platformResourceID": "asc"},
		{"resourceType": "asc"},
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := paginator.Close(ctx); err != nil {
			w.logger.Error("failed to close paginator", zap.Error(err))
		}
	}()

	w.logger.Info("ComplianceResultsIndex paginator ready")

	jd := types2.JobDocs{
		BenchmarkSummary: types2.BenchmarkSummary{
			BenchmarkID:      j.BenchmarkID,
			JobID:            j.ID,
			EvaluatedAtEpoch: j.CreatedAt.Unix(),
			Integrations: types2.BenchmarkSummaryResult{
				BenchmarkResult: types2.ResultGroup{
					Result: types2.Result{
						QueryResult:    map[types.ComplianceStatus]int{},
						SeverityResult: map[types.ComplianceResultSeverity]int{},
						SecurityScore:  0,
					},
					ResourceTypes: map[string]types2.Result{},
					Controls:      map[string]types2.ControlResult{},
				},
				Integrations: map[string]types2.ResultGroup{},
			},
			ResourceCollections: map[string]types2.BenchmarkSummaryResult{},
		},
		ResourcesFindings:       make(map[string]types.ResourceFinding),
		ResourcesFindingsIsDone: make(map[string]bool),
	}

	controlView := &types.ComplianceJobReportControlView{
		Controls:          make(map[string]types.AuditControlResult),
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		JobSummary: types.JobSummary{
			JobID:        j.ID,
			FrameworkID:  j.BenchmarkID,
			Auditable:    true,
			JobStartedAt: time.Now(),
		},
	}
	controlSummary := &types.ComplianceJobReportControlSummary{
		Controls: make(map[string]*types.ControlSummary),
		ControlScore: &types.ControlScore{
			TotalControls:  0,
			FailedControls: 0,
		},
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		JobSummary: types.JobSummary{
			JobID:        j.ID,
			FrameworkID:  j.BenchmarkID,
			Auditable:    true,
			JobStartedAt: time.Now(),
		},
	}
	resourceView := &types.ComplianceJobReportResourceView{
		Integrations:      make(map[string]types.AuditIntegrationResult),
		ComplianceSummary: make(map[types.ComplianceStatus]uint64),
		JobSummary: types.JobSummary{
			JobID:        j.ID,
			FrameworkID:  j.BenchmarkID,
			Auditable:    true,
			JobStartedAt: time.Now(),
		},
	}

	totalControls := make(map[string]bool)
	failedControls := make(map[string]bool)
	integrationsMap := make(map[string]bool)
	for page := 1; paginator.HasNext(); page++ {
		w.logger.Info("Next page", zap.Int("page", page))
		page, err := paginator.NextPage(ctx)
		if err != nil {
			w.logger.Error("failed to fetch next page", zap.Error(err))
			return err
		}

		platformResourceIDs := make([]string, 0, len(page))
		for _, f := range page {
			platformResourceIDs = append(platformResourceIDs, f.PlatformResourceID)
		}

		lookupResourcesMap, err := es.FetchLookupByResourceIDBatch(ctx, w.esClient, platformResourceIDs)
		if err != nil {
			w.logger.Error("failed to fetch lookup resources", zap.Error(err))
			return err
		}

		w.logger.Info("resource lookup result", zap.Any("platformResourceIDs", platformResourceIDs),
			zap.Any("lookupResourcesMap", lookupResourcesMap))
		w.logger.Info("page size", zap.Int("pageSize", len(page)))
		for _, f := range page {
			var resource *es2.LookupResource
			potentialResources := lookupResourcesMap[f.PlatformResourceID]
			if len(potentialResources) > 0 {
				resource = &potentialResources[0]
			}
			w.logger.Info("Before adding resource finding", zap.String("platform_resource_id", f.PlatformResourceID),
				zap.Any("resource", resource))
			jd.AddComplianceResult(w.logger, j, f, resource)
			addJobSummary(controlSummary, controlView, resourceView, f)
			integrationsMap[f.IntegrationID] = true
			totalControls[f.ControlID] = true
			if f.ComplianceStatus == types.ComplianceStatusALARM {
				failedControls[f.ControlID] = true
			}
		}

		var docs []es2.Doc
		for resourceIdType, isReady := range jd.ResourcesFindingsIsDone {
			if !isReady {
				w.logger.Info("resource NOT DONE", zap.String("platform_resource_id", resourceIdType))
				continue
			}
			w.logger.Info("resource DONE", zap.String("platform_resource_id", resourceIdType))
			resourceFinding := jd.ResourcesFindings[resourceIdType]
			keys, idx := resourceFinding.KeysAndIndex()
			resourceFinding.EsID = es2.HashOf(keys...)
			resourceFinding.EsIndex = idx
			docs = append(docs, resourceFinding)
			delete(jd.ResourcesFindings, resourceIdType)
			delete(jd.ResourcesFindingsIsDone, resourceIdType)
		}
		w.logger.Info("Sending resource finding docs", zap.Int("docCount", len(docs)))

		if _, err := w.esSinkClient.Ingest(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, docs); err != nil {
			w.logger.Error("failed to send to ingest", zap.Error(err))
			return err
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return err
	}

	w.logger.Info("Starting to summarizer",
		zap.Uint("job_id", j.ID),
		zap.String("benchmark_id", j.BenchmarkID),
	)

	jd.Summarize(w.logger)

	keys, idx := jd.BenchmarkSummary.KeysAndIndex()
	jd.BenchmarkSummary.EsID = es2.HashOf(keys...)
	jd.BenchmarkSummary.EsIndex = idx

	docs := make([]es2.Doc, 0, len(jd.ResourcesFindings)+1)
	docs = append(docs, jd.BenchmarkSummary)
	resourceIds := make([]string, 0, len(jd.ResourcesFindings))
	for resourceId, rf := range jd.ResourcesFindings {
		resourceIds = append(resourceIds, resourceId)
		keys, idx := rf.KeysAndIndex()
		rf.EsID = es2.HashOf(keys...)
		rf.EsIndex = idx
		docs = append(docs, rf)
	}
	if _, err := w.esSinkClient.Ingest(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, docs); err != nil {
		w.logger.Error("failed to send to ingest", zap.Error(err))
		return err
	}

	// Delete old resource findings
	if len(resourceIds) > 0 {
		err = w.deleteOldResourceFindings(ctx, j, resourceIds)
		if err != nil {
			w.logger.Error("failed to delete old resource findings", zap.Error(err))
			return err
		}
	}

	w.logger.Info("Deleting compliance results and resource findings of removed integrations", zap.String("benchmark_id", j.BenchmarkID), zap.Uint("job_id", j.ID))

	currentInregrations, err := w.integrationClient.ListIntegrations(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, nil)
	if err != nil {
		w.logger.Error("failed to list integrations", zap.Error(err), zap.String("benchmark_id", j.BenchmarkID), zap.Uint("job_id", j.ID))
		return err
	}
	currentIntegrationIds := make([]string, 0, len(currentInregrations.Integrations))
	for _, i := range currentInregrations.Integrations {
		currentIntegrationIds = append(currentIntegrationIds, i.IntegrationID)
	}

	err = w.deleteComplianceResultsAndResourceFindingsOfRemovedIntegrations(ctx, j, currentIntegrationIds)
	if err != nil {
		w.logger.Error("failed to delete compliance results and resource findings of removed integrations", zap.Error(err), zap.String("benchmark_id", j.BenchmarkID), zap.Uint("job_id", j.ID))
		return err
	}

	w.logger.Info("Finished summarizer",
		zap.Uint("job_id", j.ID),
		zap.String("benchmark_id", j.BenchmarkID),
		zap.Int("resource_count", len(jd.ResourcesFindings)),
	)

	var integrations []string
	for i, _ := range integrationsMap {
		integrations = append(integrations, i)
	}

	controlView.JobSummary.IntegrationIDs = integrations
	keys, idx = controlView.KeysAndIndex()
	controlView.EsID = es2.HashOf(keys...)
	controlView.EsIndex = idx

	err = sendDataToOpensearch(w.esClient.ES(), controlView)

	resourceView.JobSummary.IntegrationIDs = integrations
	keys, idx = resourceView.KeysAndIndex()
	resourceView.EsID = es2.HashOf(keys...)
	resourceView.EsIndex = idx

	err = sendDataToOpensearch(w.esClient.ES(), resourceView)
	if err != nil {
		return err
	}

	controlSummary.JobSummary.IntegrationIDs = integrations
	controlSummary.ControlScore.TotalControls = int64(len(totalControls))
	controlSummary.ControlScore.FailedControls = int64(len(failedControls))
	keys, idx = controlSummary.KeysAndIndex()
	controlSummary.EsID = es2.HashOf(keys...)
	controlSummary.EsIndex = idx

	err = sendDataToOpensearch(w.esClient.ES(), controlSummary)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) deleteOldResourceFindings(ctx context.Context, j types2.Job, currentResourceIds []string) error {
	// Delete old resource findings
	filters := make([]opengovernance.BoolFilter, 0, 2)
	filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("platformResourceID", currentResourceIds)))
	filters = append(filters, opengovernance.NewRangeFilter("jobId", "", "", fmt.Sprintf("%d", j.ID), ""))

	root := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
	}
	rootJson, err := json.Marshal(root)
	if err != nil {
		w.logger.Error("failed to marshal root", zap.Error(err))
		return err
	}

	task := es3.DeleteTask{
		DiscoveryJobID: j.ID,
		IntegrationID:  j.BenchmarkID,
		ResourceType:   "resource-finding",
		TaskType:       es3.DeleteTaskTypeQuery,
		Query:          string(rootJson),
		QueryIndex:     types.ResourceFindingsIndex,
	}

	keys, idx := task.KeysAndIndex()
	task.EsID = es2.HashOf(keys...)
	task.EsIndex = idx
	if _, err := w.esSinkClient.Ingest(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, []es2.Doc{task}); err != nil {
		w.logger.Error("failed to send delete message to elastic", zap.Error(err))
		return err
	}

	return nil
}

func (w *Worker) deleteComplianceResultsAndResourceFindingsOfRemovedIntegrations(ctx context.Context, j types2.Job, currentIntegrationIds []string) error {
	// Delete compliance results
	filters := make([]opengovernance.BoolFilter, 0, 2)
	filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", currentIntegrationIds)))
	filters = append(filters, opengovernance.NewRangeFilter("jobId", "", "", fmt.Sprintf("%d", j.ID), ""))

	root := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
	}
	rootJson, err := json.Marshal(root)
	if err != nil {
		w.logger.Error("failed to marshal root", zap.Error(err))
		return err
	}

	task := es3.DeleteTask{
		DiscoveryJobID: j.ID,
		IntegrationID:  j.BenchmarkID,
		ResourceType:   "compliance-result-old-integrations-removal",
		TaskType:       es3.DeleteTaskTypeQuery,
		Query:          string(rootJson),
		QueryIndex:     types.ComplianceResultsIndex,
	}

	keys, idx := task.KeysAndIndex()
	task.EsID = es2.HashOf(keys...)
	task.EsIndex = idx
	if _, err := w.esSinkClient.Ingest(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, []es2.Doc{task}); err != nil {
		w.logger.Error("failed to send delete message to elastic", zap.Error(err))
		return err
	}

	// Delete resource findings
	filters = make([]opengovernance.BoolFilter, 0, 2)
	filters = append(filters, opengovernance.NewBoolMustNotFilter(opengovernance.NewTermsFilter("integrationID", currentIntegrationIds)))
	filters = append(filters, opengovernance.NewRangeFilter("jobId", "", "", fmt.Sprintf("%d", j.ID), ""))

	root = map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		},
	}
	rootJson, err = json.Marshal(root)
	if err != nil {
		w.logger.Error("failed to marshal root", zap.Error(err))
		return err
	}

	task = es3.DeleteTask{
		DiscoveryJobID: j.ID,
		IntegrationID:  j.BenchmarkID,
		ResourceType:   "resource-finding-old-integrations-removal",
		TaskType:       es3.DeleteTaskTypeQuery,
		Query:          string(rootJson),
		QueryIndex:     types.ResourceFindingsIndex,
	}

	keys, idx = task.KeysAndIndex()
	task.EsID = es2.HashOf(keys...)
	task.EsIndex = idx
	if _, err := w.esSinkClient.Ingest(&httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}, []es2.Doc{task}); err != nil {
		w.logger.Error("failed to send delete message to elastic", zap.Error(err))
		return err
	}

	return nil
}

func addJobSummary(controlSummary *types.ComplianceJobReportControlSummary,
	controlView *types.ComplianceJobReportControlView, resourceView *types.ComplianceJobReportResourceView,
	cr types.ComplianceResult) {
	if cr.ComplianceStatus != types.ComplianceStatusALARM {
		return
	}

	if _, ok := resourceView.Integrations[cr.IntegrationID]; !ok {
		resourceView.Integrations[cr.IntegrationID] = types.AuditIntegrationResult{
			ResourceTypes: make(map[string]types.AuditResourceTypesResult),
		}
	}

	if _, ok := resourceView.ComplianceSummary[cr.ComplianceStatus]; !ok {
		resourceView.ComplianceSummary[cr.ComplianceStatus] = 0
	}
	resourceView.ComplianceSummary[cr.ComplianceStatus] += 1
	if _, ok := resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType]; !ok {
		resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType] = types.AuditResourceTypesResult{
			Resources: make(map[string]types.AuditResourceResult),
		}
	}
	if _, ok := resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID]; !ok {
		resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID] = types.AuditResourceResult{
			ResourceSummary: make(map[types.ComplianceStatus]uint64),
			Results:         make(map[types.ComplianceStatus][]types.AuditControlFinding),
			ResourceName:    cr.ResourceName,
		}
	}
	if _, ok := resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].ResourceSummary[cr.ComplianceStatus]; !ok {
		resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].ResourceSummary[cr.ComplianceStatus] = 0
	}
	resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].ResourceSummary[cr.ComplianceStatus] += 1
	if _, ok := resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].Results[cr.ComplianceStatus]; !ok {
		resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].Results[cr.ComplianceStatus] = make([]types.AuditControlFinding, 0)
	}
	resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].Results[cr.ComplianceStatus] = append(
		resourceView.Integrations[cr.IntegrationID].ResourceTypes[cr.ResourceType].Resources[cr.ResourceID].Results[cr.ComplianceStatus], types.AuditControlFinding{
			Severity:  cr.Severity,
			ControlID: cr.ControlID,
			Reason:    cr.Reason,
		})

	// Audit Summary
	if _, ok := controlView.Controls[cr.ControlID]; !ok {
		controlView.Controls[cr.ControlID] = types.AuditControlResult{
			Severity:       cr.Severity,
			ControlSummary: make(map[types.ComplianceStatus]uint64),
			Results:        make(map[types.ComplianceStatus][]types.AuditResourceFinding),
		}
	}
	if _, ok := controlView.ComplianceSummary[cr.ComplianceStatus]; !ok {
		controlView.ComplianceSummary[cr.ComplianceStatus] = 0
	}
	controlView.ComplianceSummary[cr.ComplianceStatus] += 1

	if _, ok := controlView.Controls[cr.ControlID].ControlSummary[cr.ComplianceStatus]; !ok {
		controlView.Controls[cr.ControlID].ControlSummary[cr.ComplianceStatus] = 0
	}
	if _, ok := controlView.Controls[cr.ControlID].Results[cr.ComplianceStatus]; !ok {
		controlView.Controls[cr.ControlID].Results[cr.ComplianceStatus] = make([]types.AuditResourceFinding, 0)
	}
	controlView.Controls[cr.ControlID].ControlSummary[cr.ComplianceStatus] += 1
	controlView.Controls[cr.ControlID].Results[cr.ComplianceStatus] = append(controlView.Controls[cr.ControlID].Results[cr.ComplianceStatus],
		types.AuditResourceFinding{
			ResourceID:   cr.ResourceID,
			ResourceType: cr.ResourceType,
			Reason:       cr.Reason,
		})

	if _, ok := controlSummary.ComplianceSummary[cr.ComplianceStatus]; !ok {
		controlSummary.ComplianceSummary[cr.ComplianceStatus] = 0
	}
	controlSummary.ComplianceSummary[cr.ComplianceStatus] += 1
	if v, ok := controlSummary.Controls[cr.ControlID]; !ok || v == nil {
		controlSummary.Controls[cr.ControlID] = &types.ControlSummary{
			Severity: cr.Severity,
			Alarms:   0,
			Oks:      0,
		}
	}
	switch cr.ComplianceStatus {
	case types.ComplianceStatusALARM:
		controlSummary.Controls[cr.ControlID].Alarms += 1
	case types.ComplianceStatusOK:
		controlSummary.Controls[cr.ControlID].Oks += 1
	}
	return
}

func sendDataToOpensearch(client *opensearch.Client, doc es2.Doc) error {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	keys, index := doc.KeysAndIndex()

	// Use the opensearchapi.IndexRequest to index the document
	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: es2.HashOf(keys...),
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
