package assistants

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

func (s *logging) Create(ctx context.Context, assistant *types.Assistants) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create", "assistant", fmt.Sprintf("%+v", assistant),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, assistant)
}

func (s *logging) Update(ctx context.Context, assistant *types.Assistants) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "", "assistant", assistant,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, assistant)
}

func (s *logging) Delete(ctx context.Context, tenantId uint, assistantId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "tenantId", tenantId, "assistantId", assistantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, tenantId, assistantId)
}

func (s *logging) Get(ctx context.Context, tenantId uint, assistantId string, preloads ...string) (assistant types.Assistants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Get", "tenantId", tenantId, "assistantId", assistantId, "preloads", preloads,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx, tenantId, assistantId, preloads...)
}

func (s *logging) List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (assistants []types.Assistants, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "preloads", preloads,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, tenantId, name, page, pageSize, preloads...)
}

func (s *logging) AddTool(ctx context.Context, tenantId uint, assistantId string, toolId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddTool", "tenantId", tenantId, "assistantId", assistantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddTool(ctx, tenantId, assistantId, toolId)
}

func (s *logging) RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RemoveTool", "tenantId", tenantId, "assistantId", assistantId, "toolId", toolId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.RemoveTool(ctx, tenantId, assistantId, toolId)
}

func (s *logging) ListTool(ctx context.Context, tenantId uint, assistantId, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTool", "tenantId", tenantId, "assistantId", assistantId, "name", name, "page", page, "pageSize", pageSize, "preloads", preloads,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListTool(ctx, tenantId, assistantId, name, page, pageSize, preloads...)
}

func (s *logging) FindByAssistantName(ctx context.Context, tenantId uint, assistantName string) (assistant types.Assistants, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByAssistantName", "tenantId", tenantId, "assistantName", assistantName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByAssistantName(ctx, tenantId, assistantName)
}

func (s *logging) ReplaceTools(ctx context.Context, assistant *types.Assistants, tools []types.Tools) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ReplaceTools", "assistant", assistant, "tools", tools,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ReplaceTools(ctx, assistant, tools)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.assistants", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
