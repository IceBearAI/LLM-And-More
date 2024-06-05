package auth

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) UpdateAccount(ctx context.Context, request UpdateAccountRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateAccount(ctx, request)
}

func (s *tracing) ListAccount(ctx context.Context, request ListAccountRequest) (res ListAccountResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListAccount(ctx, request)
}

func (s *tracing) CreateTenant(ctx context.Context, request CreateTenantRequest) (res TenantDetail, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateTenant(ctx, request)
}

func (s *tracing) ListTenants(ctx context.Context, request ListTenantRequest) (res ListTenantResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTenants", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListTenants(ctx, request)
}

func (s *tracing) CreateAccount(ctx context.Context, request CreateAccountRequest) (res Account, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateAccount(ctx, request)
}

func (s *tracing) Account(ctx context.Context, email string) (res accountResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Account", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("email", email, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Account(ctx, email)
}

func (s *tracing) Login(ctx context.Context, username, password string) (res loginResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Login", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.auth",
	})
	defer func() {
		span.LogKV("username", username, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Login(ctx, username, password)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
