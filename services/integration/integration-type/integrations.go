package integration_type

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	
	"github.com/opengovern/opencomply/services/integration/models"
	hczap "github.com/zaffka/zap-to-hclog"
)

var integrationTypes = map[integration.Type]interfaces.IntegrationType{
	
}

type IntegrationTypeManager struct {
	logger            *zap.Logger
	hcLogger          hclog.Logger
	IntegrationTypeDb *gorm.DB
	IntegrationTypes  map[integration.Type]interfaces.IntegrationType

	clients  map[integration.Type]*plugin.Client
	retryMap map[integration.Type]int
	// mutex
	pingLocks  map[integration.Type]*sync.Mutex
	maxRetries int
}

func NewIntegrationTypeManager(logger *zap.Logger, integrationTypeDb *gorm.DB, maxRetries int, pingInterval time.Duration) *IntegrationTypeManager {
	if maxRetries == 0 {
		maxRetries = 1
	}
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

	// create directory for plugins if not exists
	baseDir := "/plugins"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			logger.Error("failed to create plugins directory", zap.Error(err))
			return nil
		}
	}

	plugins := make(map[string]string)
	for _, t := range types {
		// write the plugin to the file system
		pluginPath := filepath.Join(baseDir, t.IntegrationType.String()+".so")
		err := os.WriteFile(pluginPath, t.IntegrationPlugin, 0755)
		if err != nil {
			logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", t.IntegrationType.String()))
			continue
		}
		plugins[t.IntegrationType.String()] = pluginPath
	}

	var clients = make(map[integration.Type]*plugin.Client)
	var pingLocks = make(map[integration.Type]*sync.Mutex)
	for pluginName, pluginPath := range plugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: interfaces.HandshakeConfig,
			Plugins:         map[string]plugin.Plugin{pluginName: &interfaces.IntegrationTypePlugin{}},
			Cmd:             exec.Command(pluginPath),
			Logger:          hcLogger,
			Managed:         true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			logger.Error("failed to create plugin client", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		// Request the plugin
		raw, err := rpcClient.Dispense(pluginName)
		if err != nil {
			logger.Error("failed to dispense plugin", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		// Cast the raw interface to the appropriate interface
		itInterface, ok := raw.(interfaces.IntegrationType)
		if !ok {
			logger.Error("failed to cast plugin to integration type", zap.String("plugin", pluginName), zap.String("path", pluginPath))
			client.Kill()
			continue
		}

		iType, err := itInterface.GetIntegrationType()
		if err != nil {
			logger.Error("failed to get integration type from plugin", zap.Error(err))
			client.Kill()
			continue
		}

		integrationTypes[iType] = itInterface
		clients[iType] = client
		pingLocks[iType] = &sync.Mutex{}
	}

	manager := IntegrationTypeManager{
		logger:            logger,
		hcLogger:          hcLogger,
		IntegrationTypes:  integrationTypes,
		IntegrationTypeDb: integrationTypeDb,

		clients:    clients,
		retryMap:   make(map[integration.Type]int),
		pingLocks:  pingLocks,
		maxRetries: maxRetries,
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
			lock, ok := m.pingLocks[t]
			// Just in case, shouldn't ever happen but if happens since we init it in the new manage func this is will safeguard 99.99% of the time, the other 0.01 is when an uninitialized in the new manager integration type (which shouldn't exist) ping get called and reaches this line at teh same time in 2 parallel go routines
			if !ok {
				lock = &sync.Mutex{}
				m.pingLocks[t] = lock
			}
			lock.Lock()
			if m.retryMap[t] < m.maxRetries {
				var current models.IntegrationPlugin
				err := m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", current).First(&current).Error
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
	client, ok := m.clients[t.IntegrationType]
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
		t.OperationalStatusUpdates = append(t.OperationalStatusUpdates, string(updateJson))
		if len(t.OperationalStatusUpdates) > 20 {
			t.OperationalStatusUpdates = t.OperationalStatusUpdates[len(t.OperationalStatusUpdates)-20:]
		}
		err = m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", t).Updates(&t).Error
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
	m.clients[t.IntegrationType] = client
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
	t.OperationalStatus = models.IntegrationPluginOperationalStatusEnabled //TODO remember enabled/disabled and change back to it here
	t.OperationalStatusUpdates = append(t.OperationalStatusUpdates, string(updateJson))
	if len(t.OperationalStatusUpdates) > 20 {
		t.OperationalStatusUpdates = t.OperationalStatusUpdates[len(t.OperationalStatusUpdates)-20:]
	}
	err = m.IntegrationTypeDb.Model(&models.IntegrationPlugin{}).Where("integration_type = ?", t).Updates(&t).Error
	if err != nil {
		m.logger.Error("failed to update integration plugin operational status", zap.Error(err), zap.String("integration_type", t.IntegrationType.String()))
	}

	return nil
}
