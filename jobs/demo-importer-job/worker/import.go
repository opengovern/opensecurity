package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgtype"
	"github.com/opengovern/opencomply/jobs/demo-importer-job/db"
	"github.com/opengovern/opencomply/jobs/demo-importer-job/db/model"
	"github.com/opengovern/opencomply/jobs/demo-importer-job/types"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.uber.org/zap"
)

func ImportJob(ctx context.Context, logger *zap.Logger, migratorDb db.Database, client *opensearchapi.Client, dir string) error {
	indexConfigs, err := ReadIndexConfigs(dir)
	if err != nil {
		logger.Error("Error reading index configs", zap.Error(err))
		return err
	}
	logger.Info("Read Index Configs Done")

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

	for indexName, config := range indexConfigs {
		err := CreateIndex(ctx, client, indexName, config.Settings, config.Mappings)
		if err != nil {
			logger.Error("Error creating index", zap.String("indexName", indexName), zap.Error(err))
			return err
		}
	}
	logger.Info("Create Indices Done")

	dataFiles, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		logger.Error("Error reading data files", zap.Error(err))
		return err
	}

	logger.Info("Read Data Files Done", zap.String("files", strings.Join(dataFiles, ",")))

	var wg sync.WaitGroup
	var totalTasks int64
	var completedTasks int64

	for _, file := range dataFiles {
		if strings.HasSuffix(file, ".mapping.json") || strings.HasSuffix(file, ".settings.json") {
			continue
		}

		indexName := strings.TrimSuffix(filepath.Base(file), ".json")
		if _, exists := indexConfigs[indexName]; exists {
			atomic.AddInt64(&totalTasks, 1)
			wg.Add(1)

			go func(file, indexName string) {
				defer wg.Done()
				ProcessJSONFile(ctx, logger, client, file, indexName)

				atomic.AddInt64(&completedTasks, 1)

				m.Status = fmt.Sprintf("Importing Indices")
				var progress float64
				if totalTasks > 0 {
					progress = float64(completedTasks) / float64(totalTasks)
				}
				jobsStatus := model.ESImportProgress{
					Progress: progress,
				}
				err = updateJob(migratorDb, m, m.Status, jobsStatus)
				if err != nil {
					fmt.Println("Error updating migration job:", err.Error())
				}
				fmt.Printf("Completed %d/%d tasks\n", completedTasks, totalTasks)
			}(file, indexName)
		} else {
			fmt.Println("No index config found for file: %s", file)
		}
	}

	wg.Wait()

	fmt.Println("All indexing operations completed.")

	jobsStatus := model.ESImportProgress{
		Progress: float64(completedTasks) / float64(totalTasks),
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

// execute psql command with configs
func ImportSQLFiles(config types.DemoImporterConfig, path string) error {

	cmd := exec.Command("psql", "-h", config.PostgreSQL.Host, "-p", config.PostgreSQL.Port, "-U", config.PostgreSQL.Username, "-d", "describe", "-f", path+"/describe.sql")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.PostgreSQL.Password))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing psql command: %s", string(out))
	}
	cmd = exec.Command("psql", "-h", config.PostgreSQL.Host, "-p", config.PostgreSQL.Port, "-U", config.PostgreSQL.Username, "-d", "onboard", "-f", path+"/onboard.sql")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.PostgreSQL.Password))
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing psql command: %s", string(out))
	}
	cmd = exec.Command("psql", "-h", config.PostgreSQL.Host, "-p", config.PostgreSQL.Port, "-U", config.PostgreSQL.Username, "-d", "core", "-f", path+"/core.sql")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.PostgreSQL.Password))
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing psql command: %s", string(out))
	}
	cmd = exec.Command("psql", "-h", config.PostgreSQL.Host, "-p", config.PostgreSQL.Port, "-U", config.PostgreSQL.Username, "-d", "core", "-f", path+"/core.sql")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.PostgreSQL.Password))
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing psql command: %s", string(out))
	}
	cmd = exec.Command("psql", "-h", config.PostgreSQL.Host, "-p", config.PostgreSQL.Port, "-U", config.PostgreSQL.Username, "-d", "onboard", "-c", "DELETE FROM credentials;")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.PostgreSQL.Password))
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing psql command: %s", string(out))
	}

	return nil
}
