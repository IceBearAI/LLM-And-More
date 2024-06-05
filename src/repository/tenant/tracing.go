package tenant

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

func (s *tracing) FindTenantByTenantId(ctx context.Context, tenantId string, preloads ...string) (res types.Tenants, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindTenantByTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tenant",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindTenantByTenantId(ctx, tenantId, preloads...)
}

func (s *tracing) AddModel(ctx context.Context, id uint, models ...*types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tenant",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddModel(ctx, id, models...)
}

func (s *tracing) FindTenant(ctx context.Context, id uint, preloads ...string) (res types.Tenants, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.tenant",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindTenant(ctx, id, preloads...)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
