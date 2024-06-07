// Code generated . DO NOT EDIT.
package account

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


func (s *tracing) CreateAccount(ctx context.Context, req CreateAccountRequest) (res Account, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.account",
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

func (s *tracing) DeleteAccount(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.account",
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



func (s *tracing) ListAccount(ctx context.Context, req ListAccountRequest) (list []Account, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.account",
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




func (s *tracing) UpdateAccount(ctx context.Context, req UpdateAccountRequest) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateAccount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.account",
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



func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
