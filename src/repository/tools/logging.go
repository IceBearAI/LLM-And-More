package tools

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

func (s *logging) Create(ctx context.Context, tool *types.Tools) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create", "tool", fmt.Sprintf("%+v", tool),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, tool)
}

func (s *logging) Update(ctx context.Context, tool *types.Tools) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update", "tool", fmt.Sprintf("%+v", tool),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, tool)
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

func (s *logging) Get(ctx context.Context, tenantId uint, toolId string) (tool types.Tools, err error) {
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

func (s *logging) List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "tenantId", tenantId, "page", page, "pageSize", pageSize, "name", name, "preloads", fmt.Sprintf("%+v", preloads),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, tenantId, name, page, pageSize, preloads...)
}

func (s *logging) GetByIds(ctx context.Context, tenantId uint, toolIds []string) (tools []types.Tools, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetByIds", "tenantId", tenantId, "toolIds", fmt.Sprintf("%+v", toolIds),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetByIds(ctx, tenantId, toolIds)
}

func (s *logging) ClearToolRelation(ctx context.Context, toolIds []uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ClearToolRelation", "toolIds", fmt.Sprintf("%+v", toolIds),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ClearToolRelation(ctx, toolIds)
}

func (s *logging) GetAssistants(ctx context.Context, tenantId, toolId uint) (resp []types.Assistants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetAssistants", "tenantId", tenantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetAssistants(ctx, tenantId, toolId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.tools", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
