package types

import (
	"go.uber.org/zap"
)

type Worker struct {
	Logger *zap.Logger
	Conf   DemoImporterConfig
}
