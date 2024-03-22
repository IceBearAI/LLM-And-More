package models

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

func (t *tracing) DeleteEval(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DeleteEval(ctx, id)
}

func (t *tracing) CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CreateEval(ctx, request)
}

func (t *tracing) ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListEval(ctx, request)
}

func (t *tracing) CancelEval(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CancelEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CancelEval(ctx, id)
}

func (t *tracing) Undeploy(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "Undeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.Undeploy(ctx, id)
}

func (t *tracing) ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListModels(ctx, request)
}

func (t *tracing) CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CreateModel(ctx, request)
}

func (t *tracing) GetModel(ctx context.Context, id uint) (res Model, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetModel(ctx, id)
}

func (t *tracing) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.UpdateModel(ctx, request)
}

func (t *tracing) DeleteModel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DeleteModel(ctx, id)
}

func (t *tracing) Deploy(ctx context.Context, request ModelDeployRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "Deploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("id", request, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.Deploy(ctx, request)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
