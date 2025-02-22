package job

import (
	"crypto/sha256"
	"encoding/base64"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opengovern/opensecurity/jobs/post-install-job/db/model"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/types"
	"go.uber.org/zap"
)

func (w *Job) CheckIfUpdateIsNeeded(name string, mig types.Migration) (bool, error) {
	m, err := w.db.GetMigration(name)
	if err != nil {
		return false, err
	}

	if m == nil {
		m = &model.Migration{
			ID:             name,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			AdditionalInfo: "",
		}
		err = w.db.CreateMigration(m)
		if err != nil {
			return false, err
		}
	}

	if mig.IsGitBased() {
		w.logger.Info("Git based migration", zap.String("name", name), zap.String("commit_refs", w.commitRefs), zap.String("additional_info", m.AdditionalInfo))
		return m.AdditionalInfo != w.commitRefs, nil
	} else {
		hashes, err := w.FindFilesHashes(mig)
		if err != nil {
			return false, err
		}
		w.logger.Info("File based migration", zap.String("name", name), zap.String("hashes", hashes), zap.String("additional_info", m.AdditionalInfo))
		return m.AdditionalInfo != hashes, nil
	}
}

func (w *Job) UpdateMigration(name string, mig types.Migration) error {
	if mig.IsGitBased() {
		return w.db.UpdateMigrationAdditionalInfo(name, w.commitRefs)
	} else {
		hashes, err := w.FindFilesHashes(mig)
		if err != nil {
			return err
		}

		return w.db.UpdateMigrationAdditionalInfo(name, hashes)
	}
}

func (w *Job) FindFilesHashes(mig types.Migration) (string, error) {
	h := sha256.New()
	err := filepath.Walk(mig.AttachmentFolderPath(), func(path string, info fs.FileInfo, err error) error {
		if info == nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = h.Write(content)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
