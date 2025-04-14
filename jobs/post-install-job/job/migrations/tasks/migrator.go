package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"github.com/opengovern/opensecurity/services/tasks/worker"
	"github.com/xhit/go-str2duration/v2"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/clause"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/opengovern/og-util/pkg/postgres"
	"go.uber.org/zap"
)

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}
func (m Migration) AttachmentFolderPath() string {
	return config.TasksGitPath
}

func (m Migration) Run(_ context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "task",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	dbm.ORM.Model(&models.Task{}).Where("1=1").Unscoped().Delete(&models.Task{})
	dbm.ORM.Model(&models.TaskRunSchedule{}).Where("1=1").Unscoped().Delete(&models.TaskRunSchedule{})

	err = filepath.WalkDir(m.AttachmentFolderPath(), func(path string, d fs.DirEntry, err error) error {
		if !(strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var task worker.Task
		err = yaml.Unmarshal(file, &task)
		if err != nil {
			return err
		}

		fillMissedConfigs(&task)

		natsJsonData, err := json.Marshal(task.NatsConfig)
		if err != nil {
			return err
		}

		var natsJsonb pgtype.JSONB
		err = natsJsonb.Set(natsJsonData)
		if err != nil {
			return err
		}

		scaleJsonData, err := json.Marshal(task.NatsConfig)
		if err != nil {
			return err
		}

		var scaleJsonb pgtype.JSONB
		err = scaleJsonb.Set(scaleJsonData)
		if err != nil {
			return err
		}

		if err = dbm.ORM.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&models.Task{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			ImageUrl:    task.ImageURL,
			Command:     task.Command,
			NatsConfig:  natsJsonb,
			ScaleConfig: scaleJsonb,
		}).Error; err != nil {
			return err
		}

		for _, runSchedule := range task.RunSchedule {
			paramsJsonData, err := json.Marshal(runSchedule.Params)
			if err != nil {
				return err
			}

			var paramsJsonb pgtype.JSONB
			err = paramsJsonb.Set(paramsJsonData)
			if err != nil {
				return err
			}

			frequencyFloat, err := parseToTotalSeconds(runSchedule.Frequency)
			if err != nil {
				return err
			}
			timeoutFloat, err := parseToTotalSeconds(runSchedule.Timeout)
			if err != nil {
				return err
			}

			if err = dbm.ORM.Create(&models.TaskRunSchedule{
				TaskID:    task.ID,
				Params:    paramsJsonb,
				Frequency: frequencyFloat,
				Timeout:   timeoutFloat,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func fillMissedConfigs(taskConfig *worker.Task) {
	if taskConfig.NatsConfig.Stream == "" {
		taskConfig.NatsConfig.Stream = taskConfig.ID
	}
	if taskConfig.NatsConfig.Consumer == "" {
		taskConfig.NatsConfig.Consumer = taskConfig.ID
	}
	if taskConfig.NatsConfig.Topic == "" {
		taskConfig.NatsConfig.Topic = taskConfig.ID
	}
	if taskConfig.NatsConfig.ResultConsumer == "" {
		taskConfig.NatsConfig.ResultConsumer = taskConfig.ID + "-result"
	}
	if taskConfig.NatsConfig.ResultTopic == "" {
		taskConfig.NatsConfig.ResultTopic = taskConfig.ID + "-result"
	}

	if taskConfig.ScaleConfig.Stream == "" {
		taskConfig.ScaleConfig.Stream = taskConfig.ID
	}
	if taskConfig.ScaleConfig.Consumer == "" {
		taskConfig.ScaleConfig.Consumer = taskConfig.ID
	}

	if taskConfig.ScaleConfig.PollingInterval == 0 {
		taskConfig.ScaleConfig.PollingInterval = 30
	}
	if taskConfig.ScaleConfig.CooldownPeriod == 0 {
		taskConfig.ScaleConfig.CooldownPeriod = 300
	}
}

func parseToTotalSeconds(input string) (float64, error) {
	duration, err := str2duration.ParseDuration(input)
	if err != nil {
		return 0, err
	}
	return duration.Seconds(), nil
}
