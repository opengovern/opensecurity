package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opencomply/services/compliance/api"

	"github.com/opengovern/og-util/pkg/ticker"
	describeApi "github.com/opengovern/opencomply/services/scheduler/api"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.uber.org/zap"
)

func (s *Scheduler) RunJobSequencer(ctx context.Context) {
	s.logger.Info("Scheduling job sequencer")

	t := ticker.NewTicker(JobSequencerInterval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		err := s.checkJobSequences(ctx)
		if err != nil {
			s.logger.Error("failed to run checkJobSequences", zap.Error(err))
			continue
		}
	}
}

func (s *Scheduler) checkJobSequences(ctx context.Context) error {
	jobs, err := s.db.ListWaitingJobSequencers()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		switch job.DependencySource {
		case model.JobSequencerJobTypeBenchmark:
			err := s.resolveBenchmarkDependency(ctx, job)
			if err != nil {
				s.logger.Error("failed to resolve benchmark dependency", zap.Uint("jobID", job.ID), zap.Error(err))
				if err := s.db.UpdateJobSequencerFailed(job.ID); err != nil {
					return err
				}
				continue
			}
		case model.JobSequencerJobTypeDescribe:
			err := s.resolveDescribeDependency(ctx, job)
			if err != nil {
				s.logger.Error("failed to resolve describe dependency", zap.Uint("jobID", job.ID), zap.Error(err))
				if err := s.db.UpdateJobSequencerFailed(job.ID); err != nil {
					return err
				}
				continue
			}
		default:
			s.logger.Error("job dependency %s not supported", zap.Uint("jobID", job.ID), zap.String("dependencySource", string(job.DependencySource)))
		}
	}
	return nil
}

func getControlPaths(benchmarkID, controlID string, currentPath []string, allBenchmarksCache []api.Benchmark) [][]string {
	for _, b := range allBenchmarksCache {
		if b.ID != benchmarkID {
			continue
		}
		paths := make([][]string, 0)
		for _, control := range b.Controls {
			if control == controlID {
				paths = append(paths, append(currentPath, benchmarkID))
				break
			}
		}

		for _, child := range b.Children {
			paths = append(paths, getControlPaths(child, controlID, append(currentPath, benchmarkID), allBenchmarksCache)...)
		}
		return paths
	}
	return nil
}

func (s *Scheduler) getParentBenchmarkPaths(rootBenchmark, controlID string) ([][]string, error) {
	benchmarks, err := s.complianceClient.ListAllBenchmarks(&httpclient.Context{
		UserRole: authApi.AdminRole,
	}, false)
	if err != nil {
		return nil, err
	}

	paths := getControlPaths(rootBenchmark, controlID, nil, benchmarks)
	paths = UniqueArray(paths)
	return paths, nil
}

func (s *Scheduler) runNextJob(ctx context.Context, job model.JobSequencer) error {
	switch job.NextJob {
	case model.JobSequencerJobTypeBenchmarkRunner:
		parameters := model.JobSequencerJobTypeBenchmarkRunnerParameters{}
		if job.NextJobParameters == nil {
			s.logger.Error("job parameters not found", zap.Uint("jobID", job.ID))
			return fmt.Errorf("job parameters not found")
		}
		err := json.Unmarshal(job.NextJobParameters.Bytes, &parameters)
		if err != nil {
			s.logger.Error("failed to unmarshal benchmark runner parameters", zap.Error(err))
			return err
		}
		controls, err := s.complianceClient.ListControl(&httpclient.Context{UserRole: authApi.AdminRole}, parameters.ControlIDs, nil)

		rootBenchmark, err := s.complianceClient.GetBenchmark(&httpclient.Context{UserRole: authApi.AdminRole}, parameters.BenchmarkID)
		if err != nil {
			s.logger.Error("failed to get benchmark", zap.Error(err))
			return err
		}

		runners := make([]*model.ComplianceRunner, 0, len(parameters.ConnectionIDs)*len(controls))
		for _, control := range controls {
			parentPaths, err := s.getParentBenchmarkPaths(parameters.BenchmarkID, control.ID)
			if len(parentPaths) == 0 {
				s.logger.Error("no parent paths found", zap.String("benchmarkID", parameters.BenchmarkID), zap.String("controlID", control.ID))
				continue
			}

			for _, integrationID := range parameters.ConnectionIDs {
				callers := make([]model.Caller, 0, len(parentPaths))
				for _, path := range parentPaths {
					caller := model.Caller{
						RootBenchmark:      parameters.BenchmarkID,
						TracksDriftEvents:  rootBenchmark.TracksDriftEvents,
						ParentBenchmarkIDs: path,
						ControlID:          control.ID,
						ControlSeverity:    control.Severity,
					}
					callers = append(callers, caller)
				}

				runnerJob := model.ComplianceRunner{
					FrameworkID:    parameters.BenchmarkID,
					PolicyID:       control.Policy.ID,
					ControlID:      control.ID,
					IntegrationID:  &integrationID,
					StartedAt:      time.Time{},
					RetryCount:     0,
					Status:         model.ComplianceRunnerCreated,
					FailureMessage: "",
				}
				err = runnerJob.SetCallers(callers)
				if err != nil {
					s.logger.Error("failed to set callers", zap.Error(err))
					return err
				}
				runners = append(runners, &runnerJob)
				if err != nil {
					return err
				}
			}
		}

		if len(runners) == 0 {
			s.logger.Error("no runners found", zap.String("benchmarkID", parameters.BenchmarkID), zap.Strings("controlIDs", parameters.ControlIDs))
			err = s.db.UpdateJobSequencerFinished(job.ID, nil)
			if err != nil {
				s.logger.Error("error while updating job sequencer", zap.Error(err))
				return err
			}
		}

		err = s.db.CreateRunnerJobs(nil, runners)
		if err != nil {
			s.logger.Error("error while creating runners", zap.Error(err))
			return err
		}

		var runnerJobIDs []int64
		for _, j := range runners {
			runnerJobIDs = append(runnerJobIDs, int64(j.ID))
		}

		err = s.db.UpdateJobSequencerFinished(job.ID, runnerJobIDs)
		if err != nil {
			s.logger.Error("error while updating job sequencer", zap.Error(err))
			return err
		}
	default:
		return fmt.Errorf("job type %s not supported", job.NextJob)
	}
	return nil
}

func (s *Scheduler) resolveBenchmarkDependency(ctx context.Context, job model.JobSequencer) error {
	allDependencyResolved := true
	for _, id := range job.DependencyList {
		complianceJob, err := s.db.GetComplianceJobByID(uint(id))
		if err != nil {
			return err
		}

		if complianceJob == nil {
			return fmt.Errorf("job not found: %v", id)
		}

		if complianceJob.Status == model.ComplianceJobCreated || complianceJob.Status == model.ComplianceJobRunnersInProgress {
			allDependencyResolved = false
			break
		}
	}

	if allDependencyResolved {
		err := s.runNextJob(ctx, job)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResourceCountResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
	} `json:"hits"`
}

func (s *Scheduler) resolveDescribeDependency(ctx context.Context, job model.JobSequencer) error {
	allDependencyResolved := true
	for _, id := range job.DependencyList {
		describeConnectionJob, err := s.db.GetDescribeIntegrationJobByID(uint(id))
		if err != nil {
			return err
		}

		if describeConnectionJob == nil {
			return fmt.Errorf("job not found: %v", id)
		}

		if describeConnectionJob.Status != describeApi.DescribeResourceJobSucceeded &&
			describeConnectionJob.Status != describeApi.DescribeResourceJobFailed &&
			describeConnectionJob.Status != describeApi.DescribeResourceJobTimeout {
			allDependencyResolved = false
			break
		}

		// Ignore sink count if the job is older than 24 hours
		if describeConnectionJob.UpdatedAt.Before(time.Now().Add(-time.Hour * 24)) {
			continue
		}

		root := make(map[string]any)
		root["query"] = map[string]any{
			"bool": map[string]any{
				"filter": []any{
					map[string]any{
						"term": map[string]any{
							"resource_job_id": id,
						},
					},
				},
			},
		}
		root["size"] = 0

		rootJson, err := json.Marshal(root)
		if err != nil {
			s.logger.Error("failed to marshal root", zap.Error(err))
			return err
		}

		var resourceCountResponse ResourceCountResponse
		err = s.es.SearchWithTrackTotalHits(ctx, InventorySummaryIndex, string(rootJson), nil, &resourceCountResponse, true)
		if err != nil {
			s.logger.Error("failed to search resource count", zap.Error(err))

		}

		if resourceCountResponse.Hits.Total.Value < int(float64(describeConnectionJob.DescribedResourceCount)*0.9) {
			allDependencyResolved = false
			break
		}

	}

	if allDependencyResolved {
		err := s.runNextJob(ctx, job)
		if err != nil {
			return err
		}
	}
	return nil
}
