package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Database struct {
	Orm *gorm.DB
}

func (db Database) Initialize() error {
	err := db.Orm.AutoMigrate(
		&models.Task{},
		&models.TaskRun{},
		&models.TaskConfigSecret{},
		&models.TaskRunSchedule{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) CreateTask(task *models.Task) error {
	tx := db.Orm.FirstOrCreate(task)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTask(id string, task *models.Task) error {
	tx := db.Orm.
		Model(&models.Task{}).
		Where("id = ?", id).
		Updates(task)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// GetTask retrieves a task by Task name
func (db Database) GetTask(id string) (*models.Task, error) {
	var task models.Task
	tx := db.Orm.Where("id = ?", id).
		First(&task)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &task, nil
}

// GetTaskRunResult retrieves a task result by Task ID
func (db Database) GetTaskRunResult(id string) (*models.TaskRun, error) {
	var task models.TaskRun
	tx := db.Orm.Model(&models.TaskRun{}).Where("id = ?", id).
		First(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &task, nil
}

// ListTaskRunResult retrieves a task result by Task ID
func (db Database) ListTaskRunResult(taskId string) ([]models.TaskRun, error) {
	var task []models.TaskRun
	tx := db.Orm.
		Model(&models.TaskRun{}).
		Where("task_id = ?", taskId).
		Order("created_at desc").
		Find(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return task, nil
}

// FetchCreatedTaskRunsByTaskID retrieves a list of task runs
func (db Database) FetchCreatedTaskRunsByTaskID(taskID string) ([]models.TaskRun, error) {
	var tasks []models.TaskRun
	tx := db.Orm.Model(&models.TaskRun{}).
		Where("task_id = ?", taskID).
		Where("status = ?", models.TaskRunStatusCreated).
		Find(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tasks, nil
}

// FetchLastCreatedTaskRunsByTaskID retrieves last task runs
func (db Database) FetchLastCreatedTaskRunsByTaskID(taskID string) (*models.TaskRun, error) {
	var task models.TaskRun
	tx := db.Orm.Model(&models.TaskRun{}).
		Where("task_id = ?", taskID).
		Where("status = ?", models.TaskRunStatusCreated).
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

// TimeoutTaskRunsByTaskID Timeout task runs for given task id by given timeout interval
func (db Database) TimeoutTaskRunsByTaskID(taskID string, timeoutInterval uint64) error {
	tx := db.Orm.
		Model(&models.TaskRun{}).
		Where(fmt.Sprintf("created_at < NOW() - INTERVAL '%d MINUTES'", timeoutInterval)).
		Where("status IN ?", []string{string(models.TaskRunStatusCreated),
			string(models.TaskRunStatusQueued),
			string(models.TaskRunStatusInProgress),
		}).
		Where("task_id = ?", taskID).
		Updates(models.TaskRun{Status: models.TaskRunStatusTimeout})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) CreateTaskRun(taskRun *models.TaskRun) error {
	tx := db.Orm.Create(taskRun)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// UpdateTaskRun creates a task result
func (db Database) UpdateTaskRun(runID uint, status models.TaskRunStatus, result pgtype.JSONB, failureMessage string) error {
	tx := db.Orm.Where("id = ?", runID).Updates(&models.TaskRun{
		Status: status, Result: result, FailureMessage: failureMessage,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// GetTaskList retrieves a list of tasks
func (db Database) GetTaskList() ([]models.Task, error) {
	var tasks []models.Task
	tx := db.Orm.Order("created_at desc").Find(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tasks, nil
}

// GetEnabledTaskList retrieves a list of tasks
func (db Database) GetEnabledTaskList() ([]models.Task, error) {
	var tasks []models.Task
	tx := db.Orm.Model(models.Task{}).Where("is_enabled = ?", true).Find(&tasks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tasks, nil
}

// GetTaskRunSchedules retrieves a list of task run schedules
func (db Database) GetTaskRunSchedules(taskId string) ([]models.TaskRunSchedule, error) {
	var runSchedules []models.TaskRunSchedule
	tx := db.Orm.Model(models.TaskRunSchedule{}).Where("task_id = ?", taskId).Find(&runSchedules)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return runSchedules, nil
}

func (db Database) SetTaskConfigSecret(configSecret models.TaskConfigSecret) error {
	tx := db.Orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "task_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"secret", "health_status"}),
	}).Create(&configSecret)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTaskConfigSecretHealthStatus(taskId string, healthStatus models.TaskSecretHealthStatus) error {
	tx := db.Orm.Model(&models.TaskConfigSecret{}).Where("task_id = ?", taskId).Update("health_status", healthStatus)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetTaskConfigSecret(taskId string) (*models.TaskConfigSecret, error) {
	var configSecret models.TaskConfigSecret
	tx := db.Orm.Model(&models.TaskConfigSecret{}).Where("task_id = ?", taskId).First(&configSecret)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &configSecret, nil
}

func (db Database) UpdateTaskRunScheduleLastRun(id uint) error {
	tx := db.Orm.Model(&models.TaskRunSchedule{}).Where("id = ?", id).Update("last_run", time.Now())
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
