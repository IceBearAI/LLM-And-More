package ldapcli

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Connection(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Connection", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.ldapcli",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Connection(ctx)
}

func (s *tracing) Authenticate(ctx context.Context, account string, password string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Authenticate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.ldapcli",
	})
	defer func() {
		span.LogKV("account", account, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Authenticate(ctx, account, password)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
