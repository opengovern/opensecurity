package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
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

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
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

	itOrm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "integration_types",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	itDbm := db.Database{ORM: itOrm}

	dbm.ORM.Model(&models.Task{}).Where("1=1").Unscoped().Delete(&models.Task{})
	dbm.ORM.Model(&models.TaskRunSchedule{}).Where("1=1").Unscoped().Delete(&models.TaskRunSchedule{})
	itDbm.ORM.Model(&models.TaskBinary{}).Where("1=1").Unscoped().Delete(&models.TaskBinary{})

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

		scaleJsonData, err := json.Marshal(task.ScaleConfig)
		if err != nil {
			return err
		}

		var scaleJsonb pgtype.JSONB
		err = scaleJsonb.Set(scaleJsonData)
		if err != nil {
			return err
		}

		timeoutFloat, err := parseToTotalSeconds(task.Timeout)
		if err != nil {
			return err
		}

		if err = dbm.ORM.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&models.Task{
			ID:                  task.ID,
			Name:                task.Name,
			IsEnabled:           task.IsEnabled,
			Description:         task.Description,
			ImageUrl:            task.ImageURL,
			SteampipePluginName: task.SteampipePluginName,
			ArtifactsUrl:        task.ArtifactsURL,
			Command:             task.Command,
			Timeout:             timeoutFloat,
			NatsConfig:          natsJsonb,
			ScaleConfig:         scaleJsonb,
		}).Error; err != nil {
			return err
		}

		err = loadCloudqlBinary(itDbm, logger, task)

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

			if err = dbm.ORM.Create(&models.TaskRunSchedule{
				TaskID:    task.ID,
				Params:    paramsJsonb,
				Frequency: frequencyFloat,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	inClusterConfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("failed to get in cluster config", zap.Error(err))
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(inClusterConfig)
	if err != nil {
		logger.Error("failed to create clientset", zap.Error(err))
		return err
	}

	err = restartCloudQLEnabledServices(ctx, logger, clientset)
	if err != nil {
		logger.Error("failed to restart cloudQL enabled services", zap.Error(err))
		return err
	}

	err = restartTaskService(ctx, logger, clientset)
	if err != nil {
		logger.Error("failed to restart service", zap.Error(err))
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

func restartCloudQLEnabledServices(ctx context.Context, logger *zap.Logger, clientset *kubernetes.Clientset) error {
	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		logger.Error("current namespace lookup failed")
		return errors.New("current namespace lookup failed")
	}

	err := clientset.CoreV1().Pods(currentNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "cloudql-enabled=true"})
	if err != nil {
		logger.Error("failed to delete pods", zap.Error(err))
		return err
	}

	return nil
}

func restartTaskService(ctx context.Context, logger *zap.Logger, clientset *kubernetes.Clientset) error {
	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		logger.Error("current namespace lookup failed")
		return errors.New("current namespace lookup failed")
	}

	err := clientset.CoreV1().Pods(currentNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "app=task-service"})
	if err != nil {
		logger.Error("failed to delete pods", zap.Error(err))
		return err
	}

	return nil
}

func loadCloudqlBinary(itDbm db.Database, logger *zap.Logger, task worker.Task) (err error) {
	baseDir := "/tasks"

	// create tmp directory if not exists
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
			logger.Error("failed to create tmp directory", zap.Error(err))
			return err
		}
	}

	// download files from urls

	if task.ArtifactsURL == "" || task.SteampipePluginName == "" {
		return fmt.Errorf("task artifacts url or steampipe-plugin name is empty")
	}
	url := task.ArtifactsURL
	// remove existing files
	if err = os.RemoveAll(baseDir + "/" + task.ID); err != nil {
		logger.Error("failed to remove existing files", zap.Error(err), zap.String("id", task.ID), zap.String("path", baseDir+"/integration_type"))
		return err
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/" + task.ID,
		Mode: getter.ClientModeDir,
	}
	err = downloader.Get()
	if err != nil {
		logger.Error("failed to get integration binaries", zap.Error(err), zap.String("id", task.ID))
		return err
	}

	cloudqlPlugin, err := os.ReadFile(baseDir + "/" + task.ID + "/cloudql-plugin")
	if err != nil {
		logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("id", task.ID))
		return err
	}

	if err = itDbm.ORM.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "task_id"}},
		DoNothing: true,
	}).Create(&models.TaskBinary{
		TaskID:        task.ID,
		CloudQlPlugin: cloudqlPlugin,
	}).Error; err != nil {
		return err
	}

	return nil
}
