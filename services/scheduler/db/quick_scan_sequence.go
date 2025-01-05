package db

import (
	"errors"

	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) CreateQuickScanSequence(job *model.QuickScanSequence) (uint, error) {
	tx := db.ORM.
		Model(&model.QuickScanSequence{}).
		Create(job)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return job.ID, nil
}

func (db Database) FetchCreatedQuickScanSequences() ([]model.QuickScanSequence, error) {
	var jobs []model.QuickScanSequence
	tx := db.ORM.Model(&model.QuickScanSequence{}).Where("status = ?", model.QuickScanSequenceCreated).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) GetQuickScanSequenceByID(ID uint) (*model.QuickScanSequence, error) {
	var job model.QuickScanSequence
	tx := db.ORM.Where("id = ?", ID).Find(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &job, nil
}

func (db Database) ListQuickScanSequences() ([]model.QuickScanSequence, error) {
	var job []model.QuickScanSequence
	tx := db.ORM.Model(&model.QuickScanSequence{}).First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) UpdateQuickScanSequenceStatus(jobId uint, status model.QuickScanSequenceStatus, failureReason string) error {
	tx := db.ORM.Model(&model.QuickScanSequence{}).Where("id = ?", jobId).
		Updates(model.QuickScanSequence{Status: status, FailureMessage: failureReason})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
