package auth

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) UpdateTenant(ctx context.Context, data *types.Tenants) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTenant", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateTenant(ctx, data)
}

func (s *logging) DeleteTenant(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteTenant", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteTenant(ctx, id)
}

func (s *logging) ListAccount(ctx context.Context, request ListAccountRequest) (res []types.Accounts, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListAccount", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListAccount(ctx, request)
}

func (s *logging) UpdateAccount(ctx context.Context, data *types.Accounts) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateAccount", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateAccount(ctx, data)
}

func (s *logging) GetAccountById(ctx context.Context, id uint) (res types.Accounts, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetAccountById", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetAccountById(ctx, id)
}

func (s *logging) DeleteAccount(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "DeleteAccount", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteAccount(ctx, id)
}

func (s *logging) CreateTenant(ctx context.Context, data *types.Tenants) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTenant", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateTenant(ctx, data)
}

func (s *logging) ListTenants(ctx context.Context, request ListTenantRequest) (res []types.Tenants, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTenants", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListTenants(ctx, request)
}

func (s *logging) CreateAccount(ctx context.Context, data *types.Accounts, tenantId uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAccount", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateAccount(ctx, data, tenantId)
}

func (s *logging) GetTenantByUuid(ctx context.Context, uuid string, preload ...string) (res types.Tenants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTenantByUuid", "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTenantByUuid(ctx, uuid, preload...)
}

func (s *logging) GetAccountByEmail(ctx context.Context, email string, preload ...string) (res types.Accounts, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "GetAccountByEmail", "email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetAccountByEmail(ctx, email, preload...)
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
