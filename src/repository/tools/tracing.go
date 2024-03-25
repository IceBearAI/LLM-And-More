package tools

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Create(ctx context.Context, tool *types.Tools) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tool", tool, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, tool)
}

func (s *tracing) Update(ctx context.Context, tool *types.Tools) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tool", tool, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, tool)
}

func (s *tracing) Delete(ctx context.Context, tenantId uint, toolId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, tenantId, toolId)
}

func (s *tracing) Get(ctx context.Context, tenantId uint, toolId string) (tool types.Tools, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Get", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Get(ctx, tenantId, toolId)
}

func (s *tracing) List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "page", page, "pageSize", pageSize, "name", name, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, tenantId, name, page, pageSize, preloads...)
}

func (s *tracing) GetByIds(ctx context.Context, tenantId uint, toolIds []string) (tools []types.Tools, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolIds", toolIds, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetByIds(ctx, tenantId, toolIds)
}

func (s *tracing) ClearToolRelation(ctx context.Context, toolIds []uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ClearToolRelation", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("toolIds", toolIds, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ClearToolRelation(ctx, toolIds)
}

func (s *tracing) GetAssistants(ctx context.Context, tenantId, toolId uint) (resp []types.Assistants, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetAssistants", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tools",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetAssistants(ctx, tenantId, toolId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
