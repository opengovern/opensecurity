package types

import (
	"context"
	"github.com/kaytu-io/kaytu-engine/services/migrator/config"
	"go.uber.org/zap"
)

type Migration interface {
	Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error
	IsGitBased() bool
	AttachmentFolderPath() string
}
