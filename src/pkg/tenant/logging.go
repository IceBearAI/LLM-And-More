// Code generated . DO NOT EDIT.
package tenant

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}



func (s *logging) CreateTenant(ctx context.Context, req CreateTenantRequest) (res TenantDetail, err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTenant",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateTenant(ctx, req)

}



func (s *logging) DeleteTenant(ctx context.Context, id uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteTenant",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.DeleteTenant(ctx, id)

}



func (s *logging) ListTenants(ctx context.Context, req ListTenantRequest) (list []TenantDetail, total int64, err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTenants",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListTenants(ctx, req)

}



func (s *logging) UpdateTenant(ctx context.Context, req UpdateTenantRequest) (err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTenant",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateTenant(ctx, req)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.tenant", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
