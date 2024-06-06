// Code generated . DO NOT EDIT.
package auth

import (
	"context"
	"encoding/json"

	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) CreateAccount(ctx context.Context, data *types.Accounts) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateAccount(ctx, data)

}

func (s *tracing) CreateTenant(ctx context.Context, data *types.Tenants) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateTenant(ctx, data)

}

func (s *tracing) DeleteAccount(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
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
		Value: "repository.auth",
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

func (s *tracing) GetAccountByEmail(ctx context.Context, email string, preload ...string) (res types.Accounts, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetAccountByEmail", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		span.LogKV(
			"email", email, "preload", preloadJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetAccountByEmail(ctx, email, preload...)

}

func (s *tracing) GetAccountById(ctx context.Context, id uint) (res types.Accounts, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetAccountById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetAccountById(ctx, id)

}

func (s *tracing) GetTenantById(ctx context.Context, id uint, preload ...string) (res types.Tenants, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTenantById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		span.LogKV(
			"id", id, "preload", preloadJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetTenantById(ctx, id, preload...)

}

func (s *tracing) GetTenantByUuid(ctx context.Context, uuid string, preload ...string) (res types.Tenants, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTenantByUuid", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		span.LogKV(
			"uuid", uuid, "preload", preloadJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetTenantByUuid(ctx, uuid, preload...)

}

func (s *tracing) ListAccount(ctx context.Context, request ListAccountRequest) (res []types.Accounts, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		span.LogKV(
			"request", requestJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.ListAccount(ctx, request)

}

func (s *tracing) ListTenants(ctx context.Context, request ListTenantRequest) (res []types.Tenants, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTenants", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		span.LogKV(
			"request", requestJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.ListTenants(ctx, request)

}

func (s *tracing) UpdateAccount(ctx context.Context, data *types.Accounts) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateAccount(ctx, data)

}

func (s *tracing) UpdateTenant(ctx context.Context, data *types.Tenants) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTenant", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.auth",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateTenant(ctx, data)

}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
