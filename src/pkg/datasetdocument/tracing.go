package datasetdocument

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) ListDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []datasetDocument, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListDocuments", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datsetdocument",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListDocuments(ctx, tenantId, name, page, pageSize)
}

func (s *tracing) CreateDocument(ctx context.Context, tenantId uint, data documentCreateRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateDocument", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datsetdocument",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateDocument(ctx, tenantId, data)
}

func (s *tracing) DeleteDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteDocument", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datsetdocument",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteDocument(ctx, tenantId, uuid)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
