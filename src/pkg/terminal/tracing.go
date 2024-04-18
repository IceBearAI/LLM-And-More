package terminal

import (
	"context"
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) HandleTerminalSession(session sockjs.Session) {
	t.next.HandleTerminalSession(session)
}

func (t *tracing) Token(ctx context.Context, tenantId, userId uint, resourceType string, name string) (res tokenResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "Token", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "userId", userId, "resourceType", resourceType, "name", name, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.Token(ctx, tenantId, userId, resourceType, name)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
