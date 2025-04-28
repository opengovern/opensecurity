package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/platformspec"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"github.com/opengovern/opensecurity/services/tasks/worker/consts"
	"github.com/xhit/go-str2duration/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

var (
	ESAddress  = os.Getenv("ELASTICSEARCH_ADDRESS")
	ESUsername = os.Getenv("ELASTICSEARCH_USERNAME")
	ESPassword = os.Getenv("ELASTICSEARCH_PASSWORD")
	ESIsOnAks  = os.Getenv("ELASTICSEARCH_ISONAKS")

	InventoryBaseURL = os.Getenv("CORE_BASEURL")
	NatsURL          = os.Getenv("NATS_URL")
)

func ValidateAndLoadTask(orm *gorm.DB, itOrm *gorm.DB, logger *zap.Logger, data []byte) error {
	validator := platformspec.NewDefaultValidator()

	// --- Process the Specification (Full Validation including Artifacts) ---
	validatedSpecInterface, err := validator.ProcessSpecification(
		data,
		"",
		"",
		platformspec.ArtifactTypeAll,
		false,
	)
	if err != nil {
		return err
	}
	switch spec := validatedSpecInterface.(type) {
	case *platformspec.TaskSpecification:
		if spec == nil {
			return errors.New("nil plugin specification")
		}
		err = LoadTask(orm, itOrm, logger, *spec)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid type for ValidateAndLoadPlugin")
	}
	return nil
}

func LoadTask(orm *gorm.DB, itOrm *gorm.DB, logger *zap.Logger, task platformspec.TaskSpecification) error {
	if strings.ToLower(task.Type) != "task" {
		return nil
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

	defaultEnvVars := defaultEnvs(&task)
	logger.Info("env variables", zap.Any("variables", defaultEnvVars))
	envVarsJsonData, err := json.Marshal(defaultEnvVars)
	if err != nil {
		return err
	}

	var envVarsJsonb pgtype.JSONB
	err = envVarsJsonb.Set(envVarsJsonData)
	if err != nil {
		return err
	}

	timeoutFloat, err := parseToTotalSeconds(task.Timeout)
	if err != nil {
		return err
	}
	var command string
	if task.Command != nil && len(task.Command) > 0 {
		command = task.Command[0]
	}

	var configs []string
	for _, config := range task.Configs {
		configs = append(configs, fmt.Sprintf("%v", config))
	}

	if err = orm.Clauses(clause.OnConflict{
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
		Command:             command,
		Timeout:             timeoutFloat,
		NatsConfig:          natsJsonb,
		ScaleConfig:         scaleJsonb,
		EnvVars:             envVarsJsonb,
		Params:              task.Params,
		Configs:             configs,
	}).Error; err != nil {
		return err
	}

	logger.Info("task added", zap.String("id", task.ID))

	err = loadCloudqlBinary(itOrm, logger, task)
	if err != nil {
		return err
	}

	logger.Info("cloudql binary loaded", zap.String("id", task.ID))

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

		if err = orm.Create(&models.TaskRunSchedule{
			ID:        runSchedule.ID,
			TaskID:    task.ID,
			Params:    paramsJsonb,
			Frequency: frequencyFloat,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func loadCloudqlBinary(orm *gorm.DB, logger *zap.Logger, task platformspec.TaskSpecification) (err error) {
	if task.ArtifactsURL == "" || task.SteampipePluginName == "" {
		logger.Warn("task artifacts url or steampipe-plugin name is empty", zap.String("task", task.ID))
		return nil
	}

	baseDir := "/tasks"

	// create tmp directory if not exists
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
			logger.Error("failed to create tmp directory", zap.Error(err))
			return err
		}
	}

	// download files from urls
	url := task.ArtifactsURL
	// remove existing files
	if err = os.RemoveAll(baseDir + "/" + task.SteampipePluginName); err != nil {
		logger.Error("failed to remove existing files", zap.Error(err), zap.String("id", task.ID), zap.String("path", baseDir+"/"+task.SteampipePluginName))
		return err
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/" + task.SteampipePluginName,
		Mode: getter.ClientModeDir,
	}
	err = downloader.Get()
	if err != nil {
		logger.Error("failed to get integration binaries", zap.Error(err), zap.String("id", task.ID))
		return err
	}

	cloudqlPlugin, err := os.ReadFile(baseDir + "/" + task.SteampipePluginName + "/cloudql-plugin")
	if err != nil {
		logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("id", task.ID))
		return err
	}

	if err = orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "task_id"}},
		DoNothing: true,
	}).Create(&models.TaskBinary{
		TaskID:        task.ID,
		CloudQlPlugin: cloudqlPlugin,
	}).Error; err != nil {
		logger.Error("failed to create task binary", zap.Error(err), zap.String("id", task.ID))
		return err
	}

	return nil
}

func defaultEnvs(taskConfig *platformspec.TaskSpecification) map[string]string {
	return map[string]string{
		consts.NatsURLEnv:                    NatsURL,
		consts.NatsConsumerEnv:               taskConfig.NatsConfig.Consumer,
		consts.NatsStreamNameEnv:             taskConfig.NatsConfig.Stream,
		consts.NatsTopicNameEnv:              taskConfig.NatsConfig.Topic,
		consts.NatsResultTopicNameEnv:        taskConfig.NatsConfig.ResultTopic,
		consts.ElasticSearchAddressEnv:       ESAddress,
		consts.ElasticSearchUsernameEnv:      ESUsername,
		consts.ElasticSearchPasswordEnv:      ESPassword,
		consts.ElasticSearchIsOnAksNameEnv:   ESIsOnAks,
		consts.ElasticSearchIsOpenSearch:     "false",
		consts.ElasticSearchAwsRegionEnv:     "",
		consts.ElasticSearchAssumeRoleArnEnv: "",
		consts.InventoryBaseURL:              InventoryBaseURL,
	}
}

func fillMissedConfigs(taskConfig *platformspec.TaskSpecification) {
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
