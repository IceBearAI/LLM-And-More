// Code generated . DO NOT EDIT.
package tenant

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) AddModel(ctx context.Context, id uint, models ...*types.Models) (err error) {

	defer func(begin time.Time) {

		modelsByte, _ := json.Marshal(models)
		modelsJson := string(modelsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddModel",

			"id", id,

			"models", modelsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.AddModel(ctx, id, models...)

}

func (s *logging) FindByPublicTenantIdItems(ctx context.Context, publicTenantIdItems []string, preloads ...string) (res []types.Tenants, err error) {

	defer func(begin time.Time) {

		publicTenantIdItemsByte, _ := json.Marshal(publicTenantIdItems)
		publicTenantIdItemsJson := string(publicTenantIdItemsByte)

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByPublicTenantIdItems",

			"publicTenantIdItems", publicTenantIdItemsJson,

			"preloads", preloadsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindByPublicTenantIdItems(ctx, publicTenantIdItems, preloads...)

}

func (s *logging) FindTenant(ctx context.Context, id uint, preloads ...string) (res types.Tenants, err error) {

	defer func(begin time.Time) {

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindTenant",

			"id", id,

			"preloads", preloadsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindTenant(ctx, id, preloads...)

}

func (s *logging) FindTenantByTenantId(ctx context.Context, tenantId string, preloads ...string) (res types.Tenants, err error) {

	defer func(begin time.Time) {

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindTenantByTenantId",

			"tenantId", tenantId,

			"preloads", preloadsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindTenantByTenantId(ctx, tenantId, preloads...)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.tenant", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
