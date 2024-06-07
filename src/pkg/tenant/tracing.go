// Code generated . DO NOT EDIT.
package tenant

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}




func (s *tracing) CreateTenant(ctx context.Context, req CreateTenantRequest) (res TenantDetail, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.tenant",
	})
	defer func() {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		span.LogKV(
			"req", reqJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateTenant(ctx, req)

}



func (s *tracing) DeleteTenant(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.tenant",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.DeleteTenant(ctx, id)

}



func (s *tracing) ListTenants(ctx context.Context, req ListTenantRequest) (list []TenantDetail, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTenants", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.tenant",
	})
	defer func() {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		span.LogKV(
			"req", reqJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.ListTenants(ctx, req)

}



func (s *tracing) UpdateTenant(ctx context.Context, req UpdateTenantRequest) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.tenant",
	})
	defer func() {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		span.LogKV(
			"req", reqJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateTenant(ctx, req)

}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
