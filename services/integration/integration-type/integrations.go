package integration_type

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/opencomply/services/integration/config"
	"github.com/opengovern/opencomply/services/integration/db"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"sync"
	"time"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"

	"github.com/opengovern/opencomply/services/integration/models"
	hczap "github.com/zaffka/zap-to-hclog"
)

const (
	TemplateDeploymentPath          string = "/integrations/deployment-template.yaml"
	TemplateManualsDeploymentPath   string = "/integrations/deployment-template-manuals.yaml"
	TemplateScaledObjectPath        string = "/integrations/scaled-object-template.yaml"
	TemplateManualsScaledObjectPath string = "/integrations/scaled-object-template-manuals.yaml"
)

var integrationTypes = map[integration.Type]interfaces.IntegrationType{}

type IntegrationTypeManager struct {
	cnf config.IntegrationPluginsConfig

	logger            *zap.Logger
	hcLogger          hclog.Logger
	IntegrationTypeDb *gorm.DB
	IntegrationTypes  map[integration.Type]interfaces.IntegrationType

	kubeClient    client.Client
	kubeClientset *kubernetes.Clientset
	database      db.Database

	Clients    map[integration.Type]*plugin.Client
	retryMap   map[integration.Type]int
	PingLocks  map[integration.Type]*sync.Mutex
	maxRetries int
}

func NewIntegrationTypeManager(logger *zap.Logger, database db.Database, integrationTypeDb *gorm.DB, kubeClient client.Client, kubeClientset *kubernetes.Clientset, cnf config.IntegrationPluginsConfig) *IntegrationTypeManager {
	maxRetries := cnf.MaxAutoRebootRetries
	if maxRetries == 0 {
		maxRetries = 1
	}
	pingInterval := time.Duration(cnf.PingIntervalSeconds) * time.Second
	if pingInterval == 0 {
		pingInterval = 5 * time.Minute
	}

	hcLogger := hczap.Wrap(logger)

	err := integrationTypeDb.AutoMigrate(&models.IntegrationPlugin{})
	if err != nil {
		logger.Error("failed to auto migrate integration plugin model", zap.Error(err))
		return nil
	}

	var types []models.IntegrationPlugin
	err = integrationTypeDb.Where("install_state = ?", models.IntegrationTypeInstallStateInstalled).Find(&types).Error
	if err != nil {
		logger.Error("failed to fetch integration types", zap.Error(err))
		return nil
	}
	var typesMap = make(map[integration.Type]models.IntegrationPlugin)
	for _, t := range types {
		typesMap[t.IntegrationType] = t
	}

	// create directory for plugins if not exists
	baseDir := "/plugins"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			logger.Error("failed to create plugins directory", zap.Error(err))
			return nil
		}
	}

	plugins := make(map[integration.Type]models.IntegrationPlugin)
	for _, t := range types {
		var pluginBinary models.IntegrationPluginBinary
		err = integrationTypeDb.Model(models.IntegrationPluginBinary{}).Where("plugin_id = ?", t.PluginID).Find(&pluginBinary).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error("could not find the binary", zap.String("plugin_id", t.PluginID))
			} else {
				logger.Error("failed to fetch plugin binary", zap.String("plugin_id", t.PluginID), zap.Error(err))
			}
			continue
		}
		// write the plugin to the file system
		pluginPath := filepath.Join(baseDir, t.IntegrationType.String()+".so")
		err := os.WriteFile(pluginPath, pluginBinary.IntegrationPlugin, 0755)
		if err != nil {
			logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", t.IntegrationType.String()))
			continue
		}
		plugins[t.IntegrationType] = t
	}

	var clients = make(map[integration.Type]*plugin.Client)
	var pingLocks = make(map[integration.Type]*sync.Mutex)

	manager := IntegrationTypeManager{
		cnf:               cnf,
		logger:            logger,
		hcLogger:          hcLogger,
		IntegrationTypes:  integrationTypes,
		IntegrationTypeDb: integrationTypeDb,

		Clients:    clients,
		retryMap:   make(map[integration.Type]int),
		PingLocks:  pingLocks,
		maxRetries: maxRetries,

		database:      database,
		kubeClient:    kubeClient,
		kubeClientset: kubeClientset,
	}

	for integrationType, p := range plugins {
		pluginPath := filepath.Join(baseDir, p.IntegrationType.String()+".so")

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: interfaces.HandshakeConfig,
			Plugins:         map[string]plugin.Plugin{integrationType.String(): &interfaces.IntegrationTypePlugin{}},
			Cmd:             exec.Command(pluginPath),
			Logger:          hcLogger,
			Managed:         true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			logger.Error("failed to create plugin client", zap.Error(err), zap.String("plugin", integrationType.String()), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		// Request the plugin
		raw, err := rpcClient.Dispense(integrationType.String())
		if err != nil {
			logger.Error("failed to dispense plugin", zap.Error(err), zap.String("plugin", integrationType.String()), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		// Cast the raw interface to the appropriate interface
		itInterface, ok := raw.(interfaces.IntegrationType)
		if !ok {
			logger.Error("failed to cast plugin to integration type", zap.String("plugin", integrationType.String()), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		manager.IntegrationTypes[integrationType] = itInterface
		manager.Clients[integrationType] = client
		manager.PingLocks[integrationType] = &sync.Mutex{}

		err = manager.EnableIntegrationTypeHelper(context.Background(), &p)
		if err != nil {
			logger.Error("failed to enable integration type", zap.Error(err))
			continue
		}

		pType, ok := typesMap[integrationType]
		if !ok {
			logger.Error("could not find the plugin from the type map", zap.String("plugin", integrationType.String()))
			continue
		}
		currentOperationalStatusUpdates, err := pType.GetStringOperationalStatusUpdates()
		if err != nil {
			logger.Error("failed to get operational status updates", zap.Error(err), zap.String("integration_type", integrationType.String()))
			continue
		}
		if len(currentOperationalStatusUpdates) == 0 {
			update := models.OperationalStatusUpdate{
				Time:      time.Now(),
				OldStatus: "",
				NewStatus: models.IntegrationPluginOperationalStatusEnabled,
				Reason:    "Successfully enabled",
			}
			updateJson, err := json.Marshal(update)
			if err != nil {
				logger.Error("failed to marshal operational status update", zap.Error(err), zap.String("integration_type", integrationType.String()))
				continue
			}
			currentOperationalStatusUpdates = append(currentOperationalStatusUpdates, string(updateJson))

			currentOperationalStatusUpdatesStr, err := json.Marshal(currentOperationalStatusUpdates)
			if err != nil {
				logger.Error("failed to marshal operational status updates", zap.Error(err), zap.String("integration_type", integrationType.String()))
				continue
			}

			err = pType.OperationalStatusUpdates.Set(currentOperationalStatusUpdatesStr)
			if err != nil {
				logger.Error("failed to set operational status updates", zap.Error(err), zap.String("integration_type", integrationType.String()))
				continue
			}
			err = integrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", pType.IntegrationType).Updates(models.IntegrationPlugin{
				OperationalStatusUpdates: pType.OperationalStatusUpdates,
			}).Error
			if err != nil {
				logger.Error("failed to update integration plugin initial operational status", zap.Error(err), zap.String("integration_type", integrationType.String()))
				continue
			}
		}

	}

	go func() {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		for range ticker.C {
			manager.PingRoutine()
		}
	}()

	return &manager
}

func (m *IntegrationTypeManager) GetIntegrationTypes() []integration.Type {
	types := make([]integration.Type, 0, len(m.IntegrationTypes))
	for t := range m.IntegrationTypes {
		types = append(types, t)
	}
	return types
}

func (m *IntegrationTypeManager) GetIntegrationType(t integration.Type) interfaces.IntegrationType {
	return m.IntegrationTypes[t]
}

func (m *IntegrationTypeManager) GetIntegrationTypeMap() map[integration.Type]interfaces.IntegrationType {
	return m.IntegrationTypes
}

func (m *IntegrationTypeManager) ParseType(str string) integration.Type {
	str = strings.ToLower(str)
	for t, _ := range m.IntegrationTypes {
		if str == strings.ToLower(t.String()) {
			return t
		}
	}
	return ""
}

func (m *IntegrationTypeManager) ParseTypes(str []string) []integration.Type {
	result := make([]integration.Type, 0, len(str))
	for _, s := range str {
		t := m.ParseType(s)
		if t == "" {
			continue
		}
		result = append(result, t)
	}
	return result
}

func (m *IntegrationTypeManager) UnparseTypes(types []integration.Type) []string {
	result := make([]string, 0, len(types))
	for _, t := range types {
		result = append(result, t.String())
	}
	return result
}

func (m *IntegrationTypeManager) PingRoutine() {
	m.logger.Info("running plugin ping routine")
	for t, it := range m.IntegrationTypes {
		err := it.Ping()
		if err != nil {
			m.logger.Warn("failed to ping integration type attemoting restart", zap.Error(err), zap.String("integration_type", t.String()), zap.Int("retry_count", m.retryMap[t]))
			lock, ok := m.PingLocks[t]
			// Just in case, shouldn't ever happen but if happens since we init it in the new manage func this is will safeguard 99.99% of the time, the other 0.01 is when an uninitialized in the new manager integration type (which shouldn't exist) ping get called and reaches this line at teh same time in 2 parallel go routines
			if !ok {
				lock = &sync.Mutex{}
				m.PingLocks[t] = lock
			}
			lock.Lock()
			if m.retryMap[t] < m.maxRetries {
				var current models.IntegrationPlugin
				err := m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", current.IntegrationType).First(&current).Error
				if err != nil {
					m.logger.Error("failed to fetch integration plugin", zap.Error(err), zap.String("integration_type", current.IntegrationType.String()))
					lock.Unlock()
					continue
				}
				m.retryMap[current.IntegrationType]++
				err = m.RetryRebootIntegrationType(&current)
				if err != nil {
					m.logger.Error("failed to restart integration type", zap.Error(err), zap.String("integration_type", current.IntegrationType.String()), zap.Int("retry_count", m.retryMap[t]))
				} else {
					m.retryMap[t] = 0
				}
			}
			lock.Unlock()
		}
	}
}

func (m *IntegrationTypeManager) RetryRebootIntegrationType(t *models.IntegrationPlugin) error {
	m.logger.Info("rebooting integration type", zap.String("integration_type", t.IntegrationType.String()), zap.String("plugin_id", t.PluginID), zap.Int("retry_count", m.retryMap[t.IntegrationType]))
	client, ok := m.Clients[t.IntegrationType]
	if ok {
		client.Kill()
	}

	pluginPath := filepath.Join("/plugins", t.IntegrationType.String()+".so")
	client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: interfaces.HandshakeConfig,
		Plugins:         map[string]plugin.Plugin{t.IntegrationType.String(): &interfaces.IntegrationTypePlugin{}},
		Cmd:             exec.Command(pluginPath),
		Logger:          m.hcLogger,
		Managed:         true,
	})

	changeToFailed := func(err error) {
		if t.OperationalStatus == models.IntegrationPluginOperationalStatusFailed {
			return
		}
		update := models.OperationalStatusUpdate{
			Time:      time.Now(),
			OldStatus: t.OperationalStatus,
			NewStatus: models.IntegrationPluginOperationalStatusFailed,
			Reason:    err.Error(),
		}
		updateJson, err := json.Marshal(update)
		if err != nil {
			m.logger.Error("failed to marshal operational status update", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
			return
		}

		currentOperationalStatusUpdates, err := t.GetStringOperationalStatusUpdates()
		if err != nil {
			m.logger.Error("failed to get operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
			return
		}

		currentOperationalStatusUpdates = append(currentOperationalStatusUpdates, string(updateJson))
		if len(currentOperationalStatusUpdates) > 20 {
			currentOperationalStatusUpdates = currentOperationalStatusUpdates[len(currentOperationalStatusUpdates)-20:]
		}
		if currentOperationalStatusUpdates == nil {
			currentOperationalStatusUpdates = []string{}
		}

		currentOperationalStatusUpdatesStr, err := json.Marshal(currentOperationalStatusUpdates)
		if err != nil {
			m.logger.Error("failed to marshal operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
			return
		}

		err = t.OperationalStatusUpdates.Set(currentOperationalStatusUpdatesStr)
		if err != nil {
			m.logger.Error("failed to set operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
			return
		}
		err = m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", t.IntegrationType).Updates(models.IntegrationPlugin{
			OperationalStatus:        models.IntegrationPluginOperationalStatusFailed,
			OperationalStatusUpdates: t.OperationalStatusUpdates,
		}).Error
		if err != nil {
			m.logger.Error("failed to update integration plugin operational status", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
		}
	}

	rpcClient, err := client.Client()
	if err != nil {
		m.logger.Error("failed to create plugin client", zap.Error(err), zap.String("plugin", t.IntegrationType.String()), zap.String("path", pluginPath))
		client.Kill()
		changeToFailed(err)
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(t.IntegrationType.String())
	if err != nil {
		m.logger.Error("failed to dispense plugin", zap.Error(err), zap.String("plugin", t.IntegrationType.String()), zap.String("path", pluginPath))
		client.Kill()
		changeToFailed(err)
		return err
	}

	// Cast the raw interface to the appropriate interface
	itInterface, ok := raw.(interfaces.IntegrationType)
	if !ok {
		m.logger.Error("failed to cast plugin to integration type", zap.String("plugin", t.IntegrationType.String()), zap.String("path", pluginPath))
		client.Kill()
		changeToFailed(err)
		return err
	}

	m.IntegrationTypes[t.IntegrationType] = itInterface
	m.Clients[t.IntegrationType] = client
	update := models.OperationalStatusUpdate{
		Time:      time.Now(),
		OldStatus: t.OperationalStatus,
		NewStatus: models.IntegrationPluginOperationalStatusEnabled,
		Reason:    "Successfully rebooted after detecting failed state",
	}
	updateJson, err := json.Marshal(update)
	if err != nil {
		m.logger.Error("failed to marshal operational status update", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
		return err
	}

	currentOperationalStatusUpdates, err := t.GetStringOperationalStatusUpdates()
	if err != nil {
		m.logger.Error("failed to get operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
		return err
	}

	t.OperationalStatus = models.IntegrationPluginOperationalStatusEnabled //TODO remember enabled/disabled and change back to it here
	currentOperationalStatusUpdates = append(currentOperationalStatusUpdates, string(updateJson))
	if len(currentOperationalStatusUpdates) > 20 {
		currentOperationalStatusUpdates = currentOperationalStatusUpdates[len(currentOperationalStatusUpdates)-20:]
	}
	if currentOperationalStatusUpdates == nil {
		currentOperationalStatusUpdates = []string{}
	}

	currentOperationalStatusUpdatesStr, err := json.Marshal(currentOperationalStatusUpdates)
	if err != nil {
		m.logger.Error("failed to marshal operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
		return err
	}

	err = t.OperationalStatusUpdates.Set(currentOperationalStatusUpdatesStr)
	if err != nil {
		m.logger.Error("failed to set operational status updates", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
		return err
	}

	err = m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", t.IntegrationType).Updates(models.IntegrationPlugin{
		OperationalStatus:        models.IntegrationPluginOperationalStatusEnabled,
		OperationalStatusUpdates: t.OperationalStatusUpdates,
	}).Error
	if err != nil {
		m.logger.Error("failed to update integration plugin operational status", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
	}

	return nil
}

func (a *IntegrationTypeManager) DisableIntegrationTypeHelper(ctx context.Context, integrationTypeName string) error {
	plugin, err := a.database.GetPluginByIntegrationType(integrationTypeName)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	integrations, err := a.database.ListIntegration([]integration.Type{integration.Type(integrationTypeName)})
	if err != nil {
		a.logger.Error("failed to list credentials", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list credential")
	}
	if len(integrations) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "integration type contains integrations, you can not disable it")
	}

	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "current namespace lookup failed")
	}
	integrationType, ok := a.GetIntegrationTypeMap()[integration.Type(integrationTypeName)]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "invalid integration type")
	}
	cnf, err := integrationType.GetConfiguration()
	if err != nil {
		a.logger.Error("failed to get configuration", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get configuration"+err.Error())
	}

	// Scheduled deployment
	var describerDeployment appsv1.Deployment
	err = a.kubeClient.Get(ctx, client.ObjectKey{
		Namespace: currentNamespace,
		Name:      cnf.DescriberDeploymentName,
	}, &describerDeployment)
	if err != nil {
		a.logger.Error("failed to get manual deployment", zap.Error(err))
	} else {
		err = a.kubeClient.Delete(ctx, &describerDeployment)
		if err != nil {
			a.logger.Error("failed to delete deployment", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete deployment")
		}
	}

	// Manual deployment
	var describerDeploymentManuals appsv1.Deployment
	err = a.kubeClient.Get(ctx, client.ObjectKey{
		Namespace: currentNamespace,
		Name:      cnf.DescriberDeploymentName + "-manuals",
	}, &describerDeploymentManuals)
	if err != nil {
		a.logger.Error("failed to get manual deployment", zap.Error(err))
	} else {
		err = a.kubeClient.Delete(ctx, &describerDeploymentManuals)
		if err != nil {
			a.logger.Error("failed to delete manual deployment", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete manual deployment")
		}
	}

	kedaEnabled, ok := os.LookupEnv("KEDA_ENABLED")
	if !ok {
		kedaEnabled = "false"
	}
	if strings.ToLower(kedaEnabled) == "true" {
		// Scheduled ScaledObject
		var describerScaledObject kedav1alpha1.ScaledObject
		err = a.kubeClient.Get(ctx, client.ObjectKey{
			Namespace: currentNamespace,
			Name:      cnf.DescriberDeploymentName + "-scaled-object",
		}, &describerScaledObject)
		if err != nil {
			a.logger.Error("failed to get scaled object", zap.Error(err))
		} else {
			err = a.kubeClient.Delete(ctx, &describerScaledObject)
			if err != nil {
				a.logger.Error("failed to delete scaled object", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete scaled object")
			}
		}

		// Manual ScaledObject
		var describerScaledObjectManuals kedav1alpha1.ScaledObject
		err = a.kubeClient.Get(ctx, client.ObjectKey{
			Namespace: currentNamespace,
			Name:      cnf.DescriberDeploymentName + "-manuals-scaled-object",
		}, &describerScaledObjectManuals)
		if err != nil {
			a.logger.Error("failed to get manual scaled object", zap.Error(err))
		} else {
			err = a.kubeClient.Delete(ctx, &describerScaledObjectManuals)
			if err != nil {
				a.logger.Error("failed to delete manual scaled object", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete manual scaled object")
			}
		}
	}
	return nil
}

func (a *IntegrationTypeManager) EnableIntegrationTypeHelper(ctx context.Context, plugin *models.IntegrationPlugin) error {
	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "current namespace lookup failed")
	}

	kedaEnabled, ok := os.LookupEnv("KEDA_ENABLED")
	if !ok {
		kedaEnabled = "false"
	}

	// Scheduled deployment
	var describerDeployment appsv1.Deployment
	templateDeploymentFile, err := os.Open(TemplateDeploymentPath)
	if err != nil {
		a.logger.Error("failed to open template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template deployment file")
	}
	defer templateDeploymentFile.Close()

	data, err := ioutil.ReadAll(templateDeploymentFile)
	if err != nil {
		a.logger.Error("failed to read template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template deployment file")
	}

	err = yaml.Unmarshal(data, &describerDeployment)
	if err != nil {
		a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
	}

	integrationType, ok := a.GetIntegrationTypeMap()[plugin.IntegrationType]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "invalid integration type")
	}
	cnf, err := integrationType.GetConfiguration()
	if err != nil {
		a.logger.Error("failed to get integration type configuration", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration type configuration")
	}

	describerDeployment.ObjectMeta.Name = cnf.DescriberDeploymentName
	describerDeployment.ObjectMeta.Namespace = currentNamespace
	if kedaEnabled == "true" {
		describerDeployment.Spec.Replicas = aws.Int32(0)
	} else {
		describerDeployment.Spec.Replicas = aws.Int32(5)
	}
	describerDeployment.Spec.Selector.MatchLabels["app"] = cnf.DescriberDeploymentName
	describerDeployment.Spec.Template.ObjectMeta.Labels["app"] = cnf.DescriberDeploymentName
	describerDeployment.Spec.Template.Spec.ServiceAccountName = "og-describer"

	container := describerDeployment.Spec.Template.Spec.Containers[0]
	container.Name = cnf.DescriberDeploymentName
	container.Image = fmt.Sprintf("%s:%s", plugin.DescriberURL, plugin.DescriberTag)
	container.Command = []string{cnf.DescriberRunCommand}
	natsUrl, ok := os.LookupEnv("NATS_URL")
	if ok {
		container.Env = append(container.Env, v1.EnvVar{
			Name:  "NATS_URL",
			Value: natsUrl,
		})
	}
	describerDeployment.Spec.Template.Spec.Containers[0] = container

	newDeployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cnf.DescriberDeploymentName,
			Namespace: currentNamespace,
			Labels: map[string]string{
				"app": cnf.DescriberDeploymentName,
			},
		},
		Spec: describerDeployment.Spec,
	}

	err = a.kubeClient.Create(ctx, &newDeployment)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		} else {
			existingDeployment := &appsv1.Deployment{}
			err = a.kubeClient.Get(ctx, client.ObjectKey{
				Name:      cnf.DescriberDeploymentName,
				Namespace: currentNamespace,
			}, existingDeployment)
			if err != nil {
				return err // Return if fetching fails
			}

			// Update the existing deployment's spec
			existingDeployment.Spec = describerDeployment.Spec

			// Apply the update
			err = a.kubeClient.Update(ctx, existingDeployment)
			if err != nil {
				return err // Return if updating fails
			}
		}
	}

	// Manual deployment
	var describerDeploymentManuals appsv1.Deployment
	templateManualsDeploymentFile, err := os.Open(TemplateManualsDeploymentPath)
	if err != nil {
		a.logger.Error("failed to open template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template manuals deployment file")
	}
	defer templateManualsDeploymentFile.Close()

	data, err = ioutil.ReadAll(templateManualsDeploymentFile)
	if err != nil {
		a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template manuals deployment file")
	}

	err = yaml.Unmarshal(data, &describerDeploymentManuals)
	if err != nil {
		a.logger.Error("failed to unmarshal template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template manuals deployment file")
	}

	describerDeploymentManuals.ObjectMeta.Name = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.ObjectMeta.Namespace = currentNamespace
	if kedaEnabled == "true" {
		describerDeploymentManuals.Spec.Replicas = aws.Int32(0)
	} else {
		describerDeploymentManuals.Spec.Replicas = aws.Int32(2)
	}
	describerDeploymentManuals.Spec.Selector.MatchLabels["app"] = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.Spec.Template.ObjectMeta.Labels["app"] = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.Spec.Template.Spec.ServiceAccountName = "og-describer"

	containerManuals := describerDeploymentManuals.Spec.Template.Spec.Containers[0]
	containerManuals.Name = cnf.DescriberDeploymentName
	containerManuals.Image = fmt.Sprintf("%s:%s", plugin.DescriberURL, plugin.DescriberTag)
	containerManuals.Command = []string{cnf.DescriberRunCommand}
	natsUrl, ok = os.LookupEnv("NATS_URL")
	if ok {
		containerManuals.Env = append(containerManuals.Env, v1.EnvVar{
			Name:  "NATS_URL",
			Value: natsUrl,
		})
	}
	describerDeploymentManuals.Spec.Template.Spec.Containers[0] = containerManuals

	newDeploymentManuals := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cnf.DescriberDeploymentName + "-manuals",
			Namespace: currentNamespace,
			Labels: map[string]string{
				"app": cnf.DescriberDeploymentName + "-manuals",
			},
		},
		Spec: describerDeploymentManuals.Spec,
	}

	err = a.kubeClient.Create(ctx, &newDeploymentManuals)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		} else {
			existingDeployment := &appsv1.Deployment{}
			err = a.kubeClient.Get(ctx, client.ObjectKey{
				Name:      cnf.DescriberDeploymentName + "-manuals",
				Namespace: currentNamespace,
			}, existingDeployment)
			if err != nil {
				return err // Return if fetching fails
			}

			// Update the existing deployment's spec
			existingDeployment.Spec = describerDeploymentManuals.Spec

			// Apply the update
			err = a.kubeClient.Update(ctx, existingDeployment)
			if err != nil {
				return err // Return if updating fails
			}
		}
	}

	if strings.ToLower(kedaEnabled) == "true" {
		// Scheduled ScaledObject
		var describerScaledObject kedav1alpha1.ScaledObject
		templateScaledObjectFile, err := os.Open(TemplateScaledObjectPath)
		if err != nil {
			a.logger.Error("failed to open template scaledobject file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template scaledobject file")
		}
		defer templateScaledObjectFile.Close()

		data, err = ioutil.ReadAll(templateScaledObjectFile)
		if err != nil {
			a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template scaledobject file")
		}

		err = yaml.Unmarshal(data, &describerScaledObject)
		if err != nil {
			a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
		}

		describerScaledObject.Spec.ScaleTargetRef.Name = cnf.DescriberDeploymentName

		trigger := describerScaledObject.Spec.Triggers[0]
		trigger.Metadata["stream"] = cnf.NatsStreamName
		soNatsUrl, _ := os.LookupEnv("SCALED_OBJECT_NATS_URL")
		trigger.Metadata["natsServerMonitoringEndpoint"] = soNatsUrl
		trigger.Metadata["consumer"] = cnf.NatsConsumerGroup + "-service"
		describerScaledObject.Spec.Triggers[0] = trigger

		newScaledObject := kedav1alpha1.ScaledObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cnf.DescriberDeploymentName + "-scaled-object",
				Namespace: currentNamespace,
			},
			Spec: describerScaledObject.Spec,
		}

		err = a.kubeClient.Create(ctx, &newScaledObject)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}

		// Manual ScaledObject
		var describerScaledObjectManuals kedav1alpha1.ScaledObject
		templateManualsScaledObjectFile, err := os.Open(TemplateManualsScaledObjectPath)
		if err != nil {
			a.logger.Error("failed to open template manuals scaledobject file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template manuals scaledobject file")
		}
		defer templateManualsScaledObjectFile.Close()

		data, err = ioutil.ReadAll(templateManualsScaledObjectFile)
		if err != nil {
			a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template manuals scaledobject file")
		}

		err = yaml.Unmarshal(data, &describerScaledObjectManuals)
		if err != nil {
			a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
		}

		describerScaledObjectManuals.Spec.ScaleTargetRef.Name = cnf.DescriberDeploymentName + "-manuals"

		triggerManuals := describerScaledObjectManuals.Spec.Triggers[0]
		triggerManuals.Metadata["stream"] = cnf.NatsStreamName
		triggerManuals.Metadata["natsServerMonitoringEndpoint"] = soNatsUrl
		triggerManuals.Metadata["consumer"] = cnf.NatsConsumerGroupManuals + "-service"
		describerScaledObjectManuals.Spec.Triggers[0] = triggerManuals

		newScaledObjectManuals := kedav1alpha1.ScaledObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cnf.DescriberDeploymentName + "-manuals-scaled-object",
				Namespace: currentNamespace,
			},
			Spec: describerScaledObjectManuals.Spec,
		}

		err = a.kubeClient.Create(ctx, &newScaledObjectManuals)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}

	return nil
}

func (a *IntegrationTypeManager) RestartCloudQLEnabledServices(ctx context.Context) error {
	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		a.logger.Error("current namespace lookup failed")
		return errors.New("current namespace lookup failed")
	}

	err := a.kubeClientset.CoreV1().Pods(currentNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "cloudql-enabled=true"})
	if err != nil {
		a.logger.Error("failed to delete pods", zap.Error(err))
		return err
	}

	return nil
}
