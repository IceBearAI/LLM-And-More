package sys

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) TemplateList(ctx context.Context, page, pageSize int, name, templateType string) (res []templateListResult, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.sys",
	})
	defer func() {
		span.LogKV("page", page, "pageSize", pageSize, "name", name, "templateType", templateType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.TemplateList(ctx, page, pageSize, name, templateType)
}

func (t *tracing) TemplateCreate(ctx context.Context, req templateCreateRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.sys",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.TemplateCreate(ctx, req)
}

func (t *tracing) TemplateUpdate(ctx context.Context, req templateCreateRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.sys",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.TemplateUpdate(ctx, req)
}

func (t *tracing) TemplateDelete(ctx context.Context, name string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.sys",
	})
	defer func() {
		span.LogKV("name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.TemplateDelete(ctx, name)
}

func (t *tracing) ListAudit(ctx context.Context, request ListAuditRequest) (resp ListAuditResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListAudit", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListAudit(ctx, request)
}

func (t *tracing) ListDict(ctx context.Context, request ListDictRequest) (resp ListDictResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListDict(ctx, request)
}

func (t *tracing) CreateDict(ctx context.Context, data CreateDictRequest) (resp Dict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateDict(ctx, data)
}

func (t *tracing) DictTreeByCode(ctx context.Context, codes []string) (resp []Dict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DictTreeByCode", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("codes", codes, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DictTreeByCode(ctx, codes)
}

func (t *tracing) UpdateDict(ctx context.Context, data UpdateDictRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateDict(ctx, data)
}

func (t *tracing) DeleteDict(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "service.sys",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteDict(ctx, id)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
