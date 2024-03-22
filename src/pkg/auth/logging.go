package auth

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) UpdateAccount(ctx context.Context, request UpdateAccountRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateAccount", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateAccount(ctx, request)
}

func (s *logging) ListAccount(ctx context.Context, request ListAccountRequest) (res ListAccountResponse, err error) {
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

func (s *logging) CreateTenant(ctx context.Context, request CreateTenantRequest) (res TenantDetail, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTenant", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateTenant(ctx, request)
}

func (s *logging) ListTenants(ctx context.Context, request ListTenantRequest) (res ListTenantResponse, err error) {
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

func (s *logging) CreateAccount(ctx context.Context, request CreateAccountRequest) (res Account, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAccount", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateAccount(ctx, request)
}

func (s *logging) Account(ctx context.Context, email string) (res accountResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Account", "email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Account(ctx, email)
}

func (s *logging) Login(ctx context.Context, username, password string) (res loginResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Login", "username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Login(ctx, username, password)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.auth", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
