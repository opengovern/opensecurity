package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db/model"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/types"
	"go.uber.org/zap"
)

var (
	ComplianceEnabled = os.Getenv("COMPLIANCE_ENABLED")
)

type GitConfig struct {
	AnalyticsGitURL         string
	ControlEnrichmentGitURL string
	githubToken             string
}

type Job struct {
	db         db.Database
	logger     *zap.Logger
	conf       config.MigratorConfig
	commitRefs string
}

func InitializeJob(
	conf config.MigratorConfig,
	logger *zap.Logger,
) (w *Job, err error) {
	w = &Job{
		conf:   conf,
		logger: logger,
	}

	cfg := postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      conf.PostgreSQL.DB,
		SSLMode: conf.PostgreSQL.SSLMode,
	}
	orm, err := postgres.NewClient(&cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new postgres client: %w", err)
	}
	w.db = db.Database{ORM: orm}
	fmt.Println("Connected to the postgres database: ", conf.PostgreSQL.DB)

	err = w.db.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failure while initializing database: %w", err)
	}

	m, err := w.db.GetMigration("main")
	if err != nil {
		return nil, err
	}

	if m == nil {
		jp := pgtype.JSONB{}
		err = jp.Set([]byte(""))
		if err != nil {
			return nil, err
		}
		m = &model.Migration{
			ID:             "main",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			AdditionalInfo: "",
			Status:         "Fetching Data",
			JobsStatus:     jp,
		}
		err = w.db.CreateMigration(m)
		if err != nil {
			return nil, err
		}
	} else {
		jp := pgtype.JSONB{}
		err = jp.Set([]byte(""))
		if err != nil {
			return nil, err
		}
		err = w.db.UpdateMigrationJob("main", "Fetching Data", jp)
		if err != nil {
			return nil, err
		}
	}
	if conf.IsManual {
		w.commitRefs, err = GitClone(conf, logger)
	} else {
		var refs string
		refs, err = GitClone(conf, logger)
		w.commitRefs = fmt.Sprintf("automated/%s", refs)
	}
	if err != nil {
		return nil, fmt.Errorf("failure while running git clone: %w", err)
	}

	return w, nil
}

func (w *Job) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("paniced with error", zap.Error(fmt.Errorf("%v", r)))
		}
	}()
	w.logger.Info("Starting migrator job")

	m, err := w.db.GetMigration("main")
	if err != nil {
		return err
	}
	w.logger.Info("Starting migration job", zap.String("manual status", m.ID))

	jobsStatus := make(map[string]model.JobInfo)
	var migrationList map[string]types.Migration
	var order []string
	if w.conf.IsManual {
		migrationList = manualMigrations
		order = ManualOrder
	} else {
		migrationList = migrations
		order = Order
	}
	for _, key := range order {
		if key == "compliance" && !(strings.ToLower(ComplianceEnabled) == "true") {
			continue
		}
		jobsStatus[key] = model.JobInfo{
			MigrationJobName: key,
			Status:           model.JobStatusPending,
			FailureReason:    "",
		}
	}

	m.Status = "Started"
	err = w.updateJob(m, m.Status, jobsStatus)

	hasFailed := false
	for _, name := range order {
		if name == "compliance" && !(strings.ToLower(ComplianceEnabled) == "true") {
			continue
		}
		mig := migrationList[name]
		w.logger.Info("running migration", zap.String("migrationName", name))

		jobsStatus, err = getJobsStatus(m)
		if err != nil {
			w.logger.Error("failed to get job status", zap.Error(err), zap.String("migrationName", name))
		}
		jobsStatus[name] = model.JobInfo{
			MigrationJobName: name,
			Status:           model.JobStatusInProgress,
			FailureReason:    "",
		}
		m.Status = fmt.Sprintf("Running migration %s", name)

		err = w.updateJob(m, m.Status, jobsStatus)
		if err != nil {
			w.logger.Error("failed to update job status", zap.Error(err), zap.String("migrationName", name))
		}

		updateFailed := false
		migErr := mig.Run(ctx, w.conf, w.logger)
		if migErr != nil {
			w.logger.Error("failed to run migration", zap.Error(migErr), zap.String("migrationName", name))
			updateFailed = true
		}

		jobsStatus, err = getJobsStatus(m)
		if err != nil {
			w.logger.Error("failed to get job status", zap.Error(err), zap.String("migrationName", name))
		}
		if updateFailed {
			hasFailed = true
			jobsStatus[name] = model.JobInfo{
				MigrationJobName: name,
				Status:           model.JobStatusFailed,
				FailureReason:    migErr.Error(),
			}
		} else {
			jobsStatus[name] = model.JobInfo{
				MigrationJobName: name,
				Status:           model.JobStatusCompleted,
				FailureReason:    "",
			}
		}
		err = w.updateJob(m, m.Status, jobsStatus)
		if err != nil {
			w.logger.Error("failed to update job status", zap.Error(err), zap.String("migrationName", name))
		}

	}

	if hasFailed {
		err = w.db.UpdateMigrationJob(m.ID, "FAILED", m.JobsStatus)
		if err != nil {
			return err
		}
	} else {
		err = w.db.UpdateMigrationJob(m.ID, "SUCCEEDED", m.JobsStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Job) updateJob(m *model.Migration, status string, jobsStatus map[string]model.JobInfo) error {
	jobsStatusJson, err := json.Marshal(jobsStatus)
	if err != nil {
		return err
	}

	jp := pgtype.JSONB{}
	err = jp.Set(jobsStatusJson)
	if err != nil {
		return err
	}
	m.JobsStatus = jp
	m.Status = status

	err = w.db.UpdateMigrationJob(m.ID, m.Status, m.JobsStatus)
	if err != nil {
		return err
	}
	return nil
}

func getJobsStatus(m *model.Migration) (map[string]model.JobInfo, error) {
	jobsStatus := make(map[string]model.JobInfo)

	if len(m.JobsStatus.Bytes) > 0 {
		err := json.Unmarshal(m.JobsStatus.Bytes, &jobsStatus)
		if err != nil {
			return nil, err
		}
	}
	return jobsStatus, nil
}
