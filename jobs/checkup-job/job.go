package checkup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	authAPI "github.com/opengovern/og-util/pkg/api"
	shared_entities "github.com/opengovern/og-util/pkg/api/shared-entities"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/jobs/checkup-job/config"
	authClient "github.com/opengovern/opensecurity/services/auth/client"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	"golang.org/x/net/context"

	"github.com/go-errors/errors"
	"github.com/opengovern/opensecurity/jobs/checkup-job/api"
	"github.com/opengovern/opensecurity/services/integration/client"
	"go.uber.org/zap"
)

var (
	UsageTrackerEndpoint = os.Getenv("USAGE_TRACKER_ENDPOINT")
)

type Job struct {
	JobID      uint
	ExecutedAt int64
}

type JobResult struct {
	JobID  uint
	Status api.CheckupJobStatus
	Error  string
}

func (j Job) Do(integrationClient client.IntegrationServiceClient, authClient authClient.AuthServiceClient,
	coreClient coreClient.CoreServiceClient, logger *zap.Logger, config config.WorkerConfig) (r JobResult) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("paniced with error:", err)
			fmt.Println(errors.Wrap(err, 2).ErrorStack())

			r = JobResult{
				JobID:  j.JobID,
				Status: api.CheckupJobFailed,
				Error:  fmt.Sprintf("paniced: %s", err),
			}
		}
	}()

	// Assume it succeeded unless it fails somewhere
	var (
		status         = api.CheckupJobSucceeded
		firstErr error = nil
	)

	fail := func(err error) {

		status = api.CheckupJobFailed
		if firstErr == nil {
			firstErr = err
		}
	}

	// Healthcheck
	logger.Info("starting healthcheck")

	counter := 0
	integrations, err := integrationClient.ListIntegrations(&httpclient.Context{
		UserRole: authAPI.EditorRole,
	}, nil)

	if err != nil {
		time.Sleep(3 * time.Minute)
		integrations, err = integrationClient.ListIntegrations(&httpclient.Context{
			UserRole: authAPI.EditorRole,
		}, nil)
		for {
			if err != nil {
				counter++
				if counter < 10 {
					logger.Warn("Waiting for status to be GREEN or YELLOW. Sleeping for 10 seconds...")
					time.Sleep(4 * time.Minute)
					continue
				}

				logger.Error("failed to check integration healthcheck", zap.Error(err))
				fail(fmt.Errorf("failed to check integration healthcheck: %w", err))
			}
			break
		}

	} else {
		for _, integrationObj := range integrations.Integrations {
			if integrationObj.LastCheck != nil && integrationObj.LastCheck.Add(8*time.Hour).After(time.Now()) {
				logger.Info("skipping integration health check", zap.String("integration_id", integrationObj.IntegrationID))
				continue
			}
			logger.Info("checking integration health", zap.String("integration_id", integrationObj.IntegrationID))
			_, err := integrationClient.IntegrationHealthcheck(&httpclient.Context{
				UserRole: authAPI.EditorRole,
			}, integrationObj.IntegrationID)
			if err != nil {
				logger.Error("failed to check integration health", zap.String("integration_id", integrationObj.IntegrationID), zap.Error(err))
				fail(fmt.Errorf("failed to check source health %s: %w", integrationObj.IntegrationID, err))
			}
		}
	}

	errMsg := ""
	if firstErr != nil {
		errMsg = firstErr.Error()
	}

	err = j.SendTelemetry(context.Background(), logger, config, integrationClient, authClient, coreClient)
	if err != nil {
		status = api.CheckupJobFailed
		errMsg = fmt.Sprintf("%s \n failed to send telemetry: %s", errMsg, err.Error())
	}

	return JobResult{
		JobID:  j.JobID,
		Status: status,
		Error:  errMsg,
	}
}

func (j *Job) SendTelemetry(ctx context.Context, logger *zap.Logger, workerConfig config.WorkerConfig,
	integrationClient client.IntegrationServiceClient, authClient authClient.AuthServiceClient, coreClient coreClient.CoreServiceClient) error {
	now := time.Now()

	httpCtx := httpclient.Context{Ctx: ctx, UserRole: authAPI.AdminRole}

	var plugins []shared_entities.UsageTrackerPluginInfo

	pluginsResponse, err := integrationClient.ListPlugins(&httpCtx)
	if err != nil {
		logger.Error("failed to list sources", zap.Error(err))
		return fmt.Errorf("failed to list sources: %w", err)
	}
	for _, p := range pluginsResponse.Items {
		plugins = append(plugins, shared_entities.UsageTrackerPluginInfo{
			Name:             p.Name,
			Version:          p.Version,
			IntegrationCount: p.Count.Total,
		})
	}

	users, err := authClient.ListUsers(&httpCtx)
	if err != nil {
		logger.Error("failed to list users", zap.Error(err))
		return fmt.Errorf("failed to list users: %w", err)
	}
	keys, err := authClient.ListApiKeys(&httpCtx)
	if err != nil {
		logger.Error("failed to list api keys", zap.Error(err))
		return fmt.Errorf("failed to list api keys: %w", err)
	}

	about, err := coreClient.GetAbout(&httpCtx)
	if err != nil {
		logger.Error("failed to get about", zap.Error(err))
		return fmt.Errorf("failed to get about: %w", err)
	}

	req := shared_entities.UsageTrackerRequest{
		InstanceID:      about.InstallID,
		Time:            now,
		Version:         about.AppVersion,
		Hostname:        workerConfig.TelemetryHostname,
		IsSsoConfigured: false,
		UserCount:       int64(len(users)),
		ApiKeyCount:     int64(len(keys)),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Error("failed to marshal telemetry request", zap.Error(err))
		return fmt.Errorf("failed to marshal telemetry request: %w", err)
	}
	var resp any
	if statusCode, err := httpclient.DoRequest(httpCtx.Ctx, http.MethodPost, UsageTrackerEndpoint, httpCtx.ToHeaders(), reqBytes, &resp); err != nil {
		logger.Error("failed to send telemetry", zap.Error(err), zap.Int("status_code", statusCode), zap.String("url", UsageTrackerEndpoint), zap.Any("req", req), zap.Any("resp", resp))
		return fmt.Errorf("failed to send telemetry request: %w", err)
	}

	logger.Info("sent telemetry", zap.String("url", UsageTrackerEndpoint))
	return nil
}
