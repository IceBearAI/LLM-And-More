// Code generated . DO NOT EDIT.
package auth

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

func (s *tracing) Account(ctx context.Context, email string) (res AccountResult, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Account", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {

		span.LogKV(
			"email", email,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.Account(ctx, email)

}

func (s *tracing) CreateAccount(ctx context.Context, req CreateAccountRequest) (res Account, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

	return s.next.CreateAccount(ctx, req)

}

func (s *tracing) CreateTenant(ctx context.Context, req CreateTenantRequest) (res TenantDetail, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

func (s *tracing) DeleteAccount(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.DeleteAccount(ctx, id)

}

func (s *tracing) DeleteTenant(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

func (s *tracing) ListAccount(ctx context.Context, req ListAccountRequest) (list []Account, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

	return s.next.ListAccount(ctx, req)

}

func (s *tracing) ListTenants(ctx context.Context, req ListTenantRequest) (list []TenantDetail, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTenants", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

func (s *tracing) Login(ctx context.Context, req LoginRequest) (res LoginResult, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Login", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

	return s.next.Login(ctx, req)

}

func (s *tracing) UpdateAccount(ctx context.Context, req UpdateAccountRequest) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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

	return s.next.UpdateAccount(ctx, req)

}

func (s *tracing) UpdateTenant(ctx context.Context, req UpdateTenantRequest) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
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
