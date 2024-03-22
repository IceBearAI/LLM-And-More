package tools

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Assistants(ctx context.Context, tenantId uint, toolId string) (resp []assistantResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "opentracing", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Assistants(ctx, tenantId, toolId)
}

func (s *tracing) Create(ctx context.Context, tenantId uint, req createRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, tenantId, req)
}

func (s *tracing) Update(ctx context.Context, tenantId uint, toolId string, req createRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, tenantId, toolId, req)
}

func (s *tracing) Delete(ctx context.Context, tenantId uint, toolId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, tenantId, toolId)
}

func (s *tracing) Get(ctx context.Context, tenantId uint, toolId string) (resp toolResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Get", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Get(ctx, tenantId, toolId)
}

func (s *tracing) List(ctx context.Context, tenantId uint, name string, page, pageSize int) (resp []toolResult, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, tenantId, name, page, pageSize)
}

func (s *tracing) Test(ctx context.Context, tenantId uint, toolId string, input string) (resp string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Test", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "toolId", toolId, "input", input, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Test(ctx, tenantId, toolId, input)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
