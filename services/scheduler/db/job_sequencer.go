package db

import (
	"fmt"
	"strings"

	"github.com/opengovern/opencomply/services/scheduler/db/model"
)

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

func (db Database) ListJobSequencersOfTypeOfToday(dependencySource, nextJob model.JobSequencerJobType) ([]model.JobSequencer, error) {
	var jobs []model.JobSequencer
	tx := db.ORM.Model(&model.JobSequencer{}).
		Where("dependency_source", dependencySource).
		Where("next_job", nextJob).
		Where("created_at > NOW() - interval '1 day'").Order("created_at desc").Find(&jobs)
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
	var nid []string
	for _, i := range nextJobIDs {
		nid = append(nid, fmt.Sprintf("%d", i))
	}
	tx := db.ORM.Model(&model.JobSequencer{}).
		Where("id = ?", id).
		Update("status", model.JobSequencerFinished).
		Update("next_job_ids", strings.Join(nid, ","))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
