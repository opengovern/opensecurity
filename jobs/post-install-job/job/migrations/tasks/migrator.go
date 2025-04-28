package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"github.com/opengovern/opensecurity/services/tasks/utils"
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

var (
	ESAddress  = os.Getenv("ELASTICSEARCH_ADDRESS")
	ESUsername = os.Getenv("ELASTICSEARCH_USERNAME")
	ESPassword = os.Getenv("ELASTICSEARCH_PASSWORD")
	ESIsOnAks  = os.Getenv("ELASTICSEARCH_ISONAKS")

	InventoryBaseURL = os.Getenv("CORE_BASEURL")
	NatsURL          = os.Getenv("NATS_URL")
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

		return utils.ValidateAndLoadTask(orm, itOrm, logger, file)
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
