package db

import "github.com/kaytu-io/kaytu-engine/pkg/describe/db/model"

func (db Database) CreateJobSequencer(job *model.JobSequencer) error {
	tx := db.ORM.
		Model(&model.JobSequencer{}).
		Create(job)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListWaitingJobSequencers() ([]model.JobSequencer, error) {
	var jobs []model.JobSequencer
	tx := db.ORM.Model(&model.JobSequencer{}).
		Where("status = ?", model.JobSequencerWaitingForDependencies).
		Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) ListLast20JobSequencers() ([]model.JobSequencer, error) {
	var jobs []model.JobSequencer
	tx := db.ORM.Model(&model.JobSequencer{}).Limit(20).Order("created_at desc").Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) UpdateJobSequencerFailed(id uint) error {
	tx := db.ORM.Model(&model.JobSequencer{}).
		Where("id = ?", id).
		Update("status", model.JobSequencerFailed)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

func (db Database) UpdateJobSequencerFinished(id uint, nextJobIDs []int64) error {
	tx := db.ORM.Model(&model.JobSequencer{}).
		Where("id = ?", id).
		Update("status", model.JobSequencerFinished).
		Update("next_job_ids", nextJobIDs)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}
