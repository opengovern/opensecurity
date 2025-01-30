package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	summarizer "github.com/opengovern/opencomply/jobs/compliance-summarizer-job"
	types2 "github.com/opengovern/opencomply/jobs/compliance-summarizer-job/types"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.uber.org/zap"
)

const SummarizerSchedulingInterval = 20 * time.Second

type SankDocumentCountResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		}
	}
}

func (s *JobScheduler) getSankDocumentCountBenchmark(ctx context.Context, benchmarkId string, parentJobID uint) (int, error) {
	request := make(map[string]any)
	filters := make([]map[string]any, 0)
	filters = append(filters, map[string]any{
		"term": map[string]any{
			"benchmarkID": benchmarkId,
		},
	})
	filters = append(filters, map[string]any{
		"term": map[string]any{
			"parentComplianceJobID": parentJobID,
		},
	})
	request["query"] = map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}
	request["size"] = 0

	query, err := json.Marshal(request)
	if err != nil {
		s.logger.Error("failed to marshal request", zap.Error(err))
		return 0, err
	}

	s.logger.Info("GetSankDocumentCountBenchmark", zap.String("benchmarkId", benchmarkId), zap.String("query", string(query)))

	sankDocumentCountResponse := SankDocumentCountResponse{}
	err = s.esClient.SearchWithTrackTotalHits(
		ctx, types.ComplianceResultsIndex,
		string(query),
		nil,
		&sankDocumentCountResponse, true,
	)
	if err != nil {
		s.logger.Error("failed to get sank document count", zap.Error(err), zap.String("benchmarkId", benchmarkId))
		return 0, err
	}

	return sankDocumentCountResponse.Hits.Total.Value, nil
}

func (s *JobScheduler) runSummarizer(ctx context.Context, manuals bool) error {
	s.logger.Info("checking for benchmarks to summarize")

	err := s.db.SetJobToRunnersInProgress()
	if err != nil {
		s.logger.Error("failed to set jobs to runners in progress", zap.Error(err))
		return err
	}

	err = s.db.UpdateComplianceJobsTimedOut(true, 24*60)
	if err != nil {
		s.logger.Error("failed to update compliance jobs timed out", zap.Error(err))
		return err
	}

	jobs, err := s.db.ListJobsWithRunnersCompleted(manuals)
	if err != nil {
		s.logger.Error("failed to list jobs with runners completed", zap.Error(err))
		return err
	}
	if len(jobs) == 0 {
		s.logger.Info("no jobs with runners completed, skipping this summarizer scheduling")
	}
	for _, job := range jobs {
		for _, framework := range job.FrameworkIds {
			sankDocCount, err := s.getSankDocumentCountBenchmark(ctx, framework, job.ID)
			if err != nil {
				s.logger.Error("failed to get sank document count", zap.Error(err), zap.String("framework", framework))
				return err
			}
			totalDocCount, err := s.db.FetchTotalFindingCountForComplianceJob(job.ID)
			if err != nil {
				s.logger.Error("failed to get total document count", zap.Error(err), zap.String("framework", framework))
				return err
			}

			lastUpdatedRunner, err := s.db.GetLastUpdatedRunnerForParent(job.ID)
			if err != nil {
				s.logger.Error("failed to get last updated runner", zap.Error(err), zap.String("framework", framework))
				return err
			}

			if sankDocCount != totalDocCount {
				// do not summarize if all docs are not sank
				// do not summarize if either less than 90% of the docs are sank or last job update is in less than an hour ago
				s.logger.Info("waiting for documents to sink",
					zap.String("framework", framework),
					zap.Int("sankDocCount", sankDocCount),
					zap.Int("totalDocCount", totalDocCount),
					zap.Time("lastUpdatedRunner", lastUpdatedRunner.UpdatedAt),
				)
				if job.Status != model.ComplianceJobSinkInProgress {
					err = s.db.UpdateComplianceJob(job.ID, model.ComplianceJobSinkInProgress, "", nil)
					if err != nil {
						s.logger.Error("failed to update compliance job status", zap.Error(err), zap.String("framework", framework))
						return err
					}
					continue
				} else if time.Now().Add(-10 * time.Minute).Before(job.UpdatedAt) {
					continue
				}
			}
			s.logger.Info("documents are sank, creating summarizer", zap.String("framework", framework), zap.Int("sankDocCount", sankDocCount), zap.Int("totalDocCount", totalDocCount))

			err = s.CreateSummarizer(framework, job.IntegrationIDs, &job.ID, job.TriggerType)
			if err != nil {
				s.logger.Error("failed to create summarizer", zap.Error(err), zap.String("benchmarkId", framework))
				return err
			}
		}
	}

	createds, err := s.db.FetchCreatedSummarizers(manuals)
	if err != nil {
		s.logger.Error("failed to fetch created summarizers", zap.Error(err))
		return err
	}

	for _, job := range createds {
		err = s.triggerSummarizer(ctx, job)
		if err != nil {
			s.logger.Error("failed to trigger summarizer", zap.Error(err), zap.String("benchmarkId", job.BenchmarkID))
			return err
		}
	}

	jobs, err = s.db.ListJobsToFinish()
	for _, job := range jobs {
		err = s.finishComplianceJob(job)
		if err != nil {
			s.logger.Error("failed to finish compliance job", zap.Error(err), zap.Uint("job-id", job.ID))
			return err
		}
	}

	err = s.db.RetryFailedSummarizers()
	if err != nil {
		s.logger.Error("failed to retry failed runners", zap.Error(err))
		return err
	}

	return nil
}

func (s *JobScheduler) finishComplianceJob(job model.ComplianceJob) error {
	err := s.updateJobRunnersState(job)
	if err != nil {
		return err
	}

	failedRunners, err := s.db.ListFailedRunnersWithParentID(job.ID)
	if err != nil {
		return err
	}

	if len(failedRunners) > 0 {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("%d runners failed: [", len(failedRunners)))
		for i, runner := range failedRunners {
			identify := fmt.Sprintf("query[%s]", runner.PolicyID)
			if callers, cErr := runner.GetCallers(); cErr == nil {
				ids := map[string]any{}
				for _, c := range callers {
					ids[c.ControlID] = struct{}{}
				}
				var uniqIDs []string
				for c := range ids {
					uniqIDs = append(uniqIDs, c)
				}
				identify = fmt.Sprintf("controls[%s]", strings.Join(uniqIDs, ","))
			}
			builder.WriteString(fmt.Sprintf("%s: %s", identify, runner.FailureMessage))
			if i != len(failedRunners)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString("]")
		return s.db.UpdateComplianceJob(job.ID, model.ComplianceJobSucceeded, builder.String(), nil)
	}

	failedSummarizers, err := s.db.ListFailedSummarizersWithParentID(job.ID)
	if err != nil {
		return err
	}

	if len(failedSummarizers) > 0 {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("%d summarizers failed: [", len(failedSummarizers)))
		for i, summarizer := range failedSummarizers {
			builder.WriteString(fmt.Sprintf("%s", summarizer.FailureMessage))
			if i != len(failedSummarizers)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString("]")
		status := model.ComplianceJobSummarizerInProgress
		return s.db.UpdateComplianceJob(job.ID, model.ComplianceJobFailed, builder.String(), &status)
	}

	return s.db.UpdateComplianceJob(job.ID, model.ComplianceJobSucceeded, "", nil)
}

func (s *JobScheduler) CreateSummarizer(benchmarkId string, integrationIDs []string, jobId *uint, triggerType model.ComplianceTriggerType) error {
	// run summarizer
	dbModel := model.ComplianceSummarizer{
		BenchmarkID:    benchmarkId,
		IntegrationIDs: integrationIDs,
		StartedAt:      time.Now(),
		Status:         summarizer.ComplianceSummarizerCreated,
		TriggerType:    triggerType,
	}
	if jobId != nil {
		dbModel.ParentJobID = *jobId
	}
	err := s.db.CreateSummarizerJob(&dbModel)
	if err != nil {
		return err
	}
	if jobId != nil {
		return s.db.UpdateComplianceJob(*jobId, model.ComplianceJobSummarizerInProgress, "", nil)
	}
	return nil
}

func (s *JobScheduler) triggerSummarizer(ctx context.Context, job model.ComplianceSummarizer) error {
	summarizerJob := types2.Job{
		ID:              job.ID,
		ComplianceJobID: job.ParentJobID,
		IntegrationIDs:  job.IntegrationIDs,
		RetryCount:      job.RetryCount,
		BenchmarkID:     job.BenchmarkID,
		CreatedAt:       job.CreatedAt,
	}
	jobJson, err := json.Marshal(summarizerJob)
	if err != nil {
		_ = s.db.UpdateSummarizerJob(job.ID, summarizer.ComplianceSummarizerFailed, job.CreatedAt, err.Error())
		return err
	}

	topic := summarizer.JobQueueTopic
	if job.TriggerType == model.ComplianceTriggerTypeManual {
		topic = summarizer.JobQueueTopicManuals
	}
	if _, err := s.jq.Produce(ctx, topic, jobJson, fmt.Sprintf("job-%d-%d", job.ID, job.RetryCount)); err != nil {
		_ = s.db.UpdateSummarizerJob(job.ID, summarizer.ComplianceSummarizerFailed, job.CreatedAt, err.Error())
		return err
	}

	return s.db.UpdateSummarizerJob(job.ID, summarizer.ComplianceSummarizerInProgress, job.CreatedAt, "")
}
