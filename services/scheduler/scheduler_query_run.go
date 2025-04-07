package describe

import (
	"context"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/ticker"
	coreApi "github.com/opengovern/opensecurity/services/core/api"
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
	namedQueries, err := s.coreClient.ListCacheEnabledQueries(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole})
	if err != nil {
		s.logger.Error("Failed to find the last job to check for CheckupJob", zap.Error(err))
		CheckupJobsCount.WithLabelValues("failure").Inc()
		return
	}

	for _, nq := range namedQueries {
		if nq.LastRun == nil || nq.LastRun.IsZero() || nq.LastRun.Before(time.Now().Add(-1*time.Hour-30*time.Minute)) {
			_, err = s.coreClient.RunQuery(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole},
				coreApi.RunQueryRequest{QueryId: nq.QueryID})
			if err != nil {
				s.logger.Error("Failed to run named query", zap.Error(err))
				return
			}
		}
	}
}
