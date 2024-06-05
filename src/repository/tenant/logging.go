package tenant

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) FindTenantByTenantId(ctx context.Context, tenantId string, preloads ...string) (res types.Tenants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindTenantByTenantId",
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindTenantByTenantId(ctx, tenantId, preloads...)
}

func (s *logging) AddModel(ctx context.Context, id uint, models ...*types.Models) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddModel(ctx, id, models...)
}

func (s *logging) FindTenant(ctx context.Context, id uint, preloads ...string) (res types.Tenants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindTenant",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindTenant(ctx, id, preloads...)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.tenant", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
