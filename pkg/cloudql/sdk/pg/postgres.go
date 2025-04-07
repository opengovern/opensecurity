package pg

import (
	"context"
	"errors"
	integration "github.com/opengovern/opensecurity/services/integration/models"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) ListIntegrations(ctx context.Context) ([]integration.Integration, error) {
	var result []integration.Integration
	err := c.db.Preload(clause.Associations).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetIntegrationByID(ctx context.Context, opengovernanceId string, id string) (*integration.Integration, error) {
	var result integration.Integration
	var err error
	tx := c.db.Preload(clause.Associations).Model(&integration.Integration{})
	switch {
	case opengovernanceId != "" && id != "":
		err = tx.Where("integration_id = ? AND provider_id = ?", opengovernanceId, id).First(&result).Error
	case opengovernanceId != "" && id == "":
		err = tx.Where("integration_id = ?", opengovernanceId).First(&result).Error
	case opengovernanceId == "" && id != "":
		err = tx.Where("provider_id = ?", id).First(&result).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (c Client) ListIntegrationGroups(ctx context.Context) ([]integration.IntegrationGroup, error) {
	var result []integration.IntegrationGroup
	err := c.db.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetIntegrationGroupByName(ctx context.Context, name string) (*integration.IntegrationGroup, error) {
	var result integration.IntegrationGroup
	err := c.db.Where("name = ?", name).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (c Client) ListTasks(ctx context.Context) ([]models.Task, error) {
	var result []models.Task
	err := c.db.Model(&models.Task{}).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetTask(ctx context.Context, taskId string) (*models.Task, error) {
	var result models.Task
	err := c.db.Model(&models.Task{}).Where("id = ?", taskId).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c Client) GetLastTaskRun(id string) (*models.TaskRun, error) {
	var task models.TaskRun
	tx := c.db.Where("task_id = ?", id).
		Order("created_at desc").
		First(&task)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &task, nil
}
