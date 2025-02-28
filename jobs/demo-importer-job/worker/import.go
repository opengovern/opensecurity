package worker

import (
	"context"
	"encoding/json"
	"fmt"
	// "os"
	// "os/exec"
	// "sync"
	"time"

	"github.com/jackc/pgtype"
	"github.com/opengovern/opensecurity/jobs/demo-importer-job/db"
	"github.com/opengovern/opensecurity/jobs/demo-importer-job/db/model"
	// "github.com/opengovern/opensecurity/jobs/demo-importer-job/types"
	"go.uber.org/zap"
)

func ImportJob(ctx context.Context, logger *zap.Logger, migratorDb db.Database) error {
	

	m, err := migratorDb.GetMigration(model.MigrationJobName)
	if err != nil {
		logger.Error("Error reading migration job", zap.Error(err))
		return err
	}
	if m == nil {
		jobsStatusJson, err := json.Marshal(model.ESImportProgress{
			Progress: 0,
		})
		if err != nil {
			return err
		}
		jp := pgtype.JSONB{}
		err = jp.Set(jobsStatusJson)
		if err != nil {
			return err
		}

		m = &model.Migration{
			ID:             model.MigrationJobName,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			AdditionalInfo: "",
			Status:         "Creating Indices",
			JobsStatus:     jp,
		}
		err = migratorDb.CreateMigration(m)
		if err != nil {
			return err
		}
	} else {
		jobsStatus := model.ESImportProgress{
			Progress: 0,
		}
		err = updateJob(migratorDb, m, "Creating Indices", jobsStatus)
		if err != nil {
			return err
		}
	}

	
	// var wg sync.WaitGroup
	// var totalTasks int64
	// var completedTasks int64

	
	// wg.Wait()

	fmt.Println("All indexing operations completed.")

	jobsStatus := model.ESImportProgress{
		Progress: float64(100) / float64(100),
	}
	err = updateJob(migratorDb, m, "COMPLETED", jobsStatus)
	if err != nil {
		return err
	}

	return nil
}

func updateJob(migratorDb db.Database, m *model.Migration, status string, jobsStatus model.ESImportProgress) error {
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

	err = migratorDb.UpdateMigrationJob(m.ID, m.Status, m.JobsStatus)
	if err != nil {
		return err
	}
	return nil
}

