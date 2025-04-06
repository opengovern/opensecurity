package describe

import (
	"context"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/ticker"
	"go.uber.org/zap"
	"time"
)

const (
	NamedQueryCacheInterval = 1 * time.Minute
)

func (s *Scheduler) RunNamedQueryCache(ctx context.Context) {
	s.logger.Info("Scheduling named query cache run on a timer")

	t := ticker.NewTicker(NamedQueryCacheInterval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		s.scheduleNamedQueryCache(ctx)
	}
}

func (s *Scheduler) scheduleNamedQueryCache(ctx context.Context) {
	_, err := s.coreClient.ListCacheEnabledQueries(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole})
	if err != nil {
		s.logger.Error("Failed to find the last job to check for CheckupJob", zap.Error(err))
		CheckupJobsCount.WithLabelValues("failure").Inc()
		return
	}

}
