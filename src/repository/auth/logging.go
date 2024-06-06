// Code generated . DO NOT EDIT.
package auth

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

func (s *logging) CreateAccount(ctx context.Context, data *types.Accounts, tenantId uint) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAccount",

			"data", dataJson,

			"tenantId", tenantId,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateAccount(ctx, data, tenantId)

}

func (s *logging) CreateAccountV2(ctx context.Context, data *types.Accounts) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAccountV2",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateAccountV2(ctx, data)

}

func (s *logging) CreateTenant(ctx context.Context, data *types.Tenants) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTenant",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateTenant(ctx, data)

}

func (s *logging) DeleteAccount(ctx context.Context, id uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteAccount",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.DeleteAccount(ctx, id)

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

func (s *logging) GetAccountByEmail(ctx context.Context, email string, preload ...string) (res types.Accounts, err error) {

	defer func(begin time.Time) {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetAccountByEmail",

			"email", email,

			"preload", preloadJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetAccountByEmail(ctx, email, preload...)

}

func (s *logging) GetAccountById(ctx context.Context, id uint) (res types.Accounts, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetAccountById",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetAccountById(ctx, id)

}

func (s *logging) GetTenantById(ctx context.Context, id uint, preload ...string) (res types.Tenants, err error) {

	defer func(begin time.Time) {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTenantById",

			"id", id,

			"preload", preloadJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetTenantById(ctx, id, preload...)

}

func (s *logging) GetTenantByUuid(ctx context.Context, uuid string, preload ...string) (res types.Tenants, err error) {

	defer func(begin time.Time) {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTenantByUuid",

			"uuid", uuid,

			"preload", preloadJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetTenantByUuid(ctx, uuid, preload...)

}

func (s *logging) ListAccount(ctx context.Context, request ListAccountRequest) (res []types.Accounts, total int64, err error) {

	defer func(begin time.Time) {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListAccount",

			"request", requestJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListAccount(ctx, request)

}

func (s *logging) ListTenants(ctx context.Context, request ListTenantRequest) (res []types.Tenants, total int64, err error) {

	defer func(begin time.Time) {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTenants",

			"request", requestJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListTenants(ctx, request)

}

func (s *logging) UpdateAccount(ctx context.Context, data *types.Accounts) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateAccount",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateAccount(ctx, data)

}

func (s *logging) UpdateTenant(ctx context.Context, data *types.Tenants) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTenant",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateTenant(ctx, data)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.auth", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
