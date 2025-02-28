package utils

import "go.uber.org/zap"

func NewZapLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/home/steampipe/.steampipe/logs/opengovernance.log",
	}
	return cfg.Build()
}
