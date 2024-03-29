package modelevaluate

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) GetEvalLog(ctx context.Context, tenantId uint, modelUUID, evalJobId string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetEvalLog", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "modelUUID", modelUUID, "evalJobId", evalJobId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetEvalLog(ctx, tenantId, modelUUID, evalJobId)
}

func (s *tracing) EvalFinish(ctx context.Context, req finishRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "EvalFinish", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.EvalFinish(ctx, req)
}

func (s *tracing) FiveGraph(ctx context.Context, req fiveGraphRequest) (res1, res2, res3 fiveGraphResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FiveGraph", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FiveGraph(ctx, req)
}

func (s *tracing) Create(ctx context.Context, req createRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, req)
}

func (s *tracing) Cancel(ctx context.Context, req cancelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Cancel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Cancel(ctx, req)
}

func (s *tracing) Delete(ctx context.Context, req deleteRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, req)
}

func (s *tracing) List(ctx context.Context, req listRequest) (res []listResult, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.modelevaluate",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, req)
}

var _ Service = &tracing{}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
