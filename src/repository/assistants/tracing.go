package assistants

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Create(ctx context.Context, assistant *types.Assistants) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("assistant", fmt.Sprintf("%+v", assistant), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, assistant)
}

func (s *tracing) Update(ctx context.Context, assistant *types.Assistants) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("assistant", assistant, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, assistant)
}

func (s *tracing) Delete(ctx context.Context, tenantId uint, assistantId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantId", assistantId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, tenantId, assistantId)
}

func (s *tracing) Get(ctx context.Context, tenantId uint, assistantId string, preloads ...string) (assistant types.Assistants, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Get", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantId", assistantId, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Get(ctx, tenantId, assistantId, preloads...)
}

func (s *tracing) List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (assistants []types.Assistants, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, tenantId, name, page, pageSize, preloads...)
}

func (s *tracing) AddTool(ctx context.Context, tenantId uint, assistantId string, toolId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddTool", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantId", assistantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddTool(ctx, tenantId, assistantId, toolId)
}

func (s *tracing) RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveTool", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantId", assistantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RemoveTool(ctx, tenantId, assistantId, toolId)
}

func (s *tracing) ListTool(ctx context.Context, tenantId uint, assistantId, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTool", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantId", assistantId, "name", name, "page", page, "pageSize", pageSize, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListTool(ctx, tenantId, assistantId, name, page, pageSize, preloads...)
}

func (s *tracing) FindByAssistantName(ctx context.Context, tenantId uint, assistantName string) (assistant types.Assistants, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByAssistantName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "assistantName", assistantName, "assistant", assistant, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByAssistantName(ctx, tenantId, assistantName)
}

func (s *tracing) ReplaceTools(ctx context.Context, assistant *types.Assistants, tools []types.Tools) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ReplaceTools", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.assistants",
	})
	defer func() {
		span.LogKV("assistant", fmt.Sprintf("%+v", assistant), "tools", fmt.Sprintf("%+v", tools), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ReplaceTools(ctx, assistant, tools)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
