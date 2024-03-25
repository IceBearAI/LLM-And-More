package tools

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Assistants(ctx context.Context, tenantId uint, toolId string) (resp []assistantResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Assistants", "tenantId", tenantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Assistants(ctx, tenantId, toolId)
}

func (s *logging) Create(ctx context.Context, tenantId uint, req createRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Test", "tenantId", tenantId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, tenantId, req)
}

func (s *logging) Update(ctx context.Context, tenantId uint, toolId string, req createRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update", "tenantId", tenantId, "toolId", toolId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, tenantId, toolId, req)
}

func (s *logging) Delete(ctx context.Context, tenantId uint, toolId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "tenantId", tenantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, tenantId, toolId)
}

func (s *logging) Get(ctx context.Context, tenantId uint, toolId string) (resp toolResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Get", "tenantId", tenantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx, tenantId, toolId)
}

func (s *logging) Test(ctx context.Context, tenantId uint, toolId string, input string) (resp string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Test", "tenantId", tenantId, "toolId", toolId, "input", input,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Test(ctx, tenantId, toolId, input)
}

func (s *logging) List(ctx context.Context, tenantId uint, query string, page, pageSize int) (datasets []toolResult, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "tenantId", tenantId, "page", page, "pageSize", pageSize, "query", query,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, tenantId, query, page, pageSize)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.tools", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
