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

func (t *tracing) ModelCheckpoint(ctx context.Context, modelName string) (res []string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ModelCheckpoint", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ModelCheckpoint(ctx, modelName)
}

func (t *tracing) ModelCard(ctx context.Context, modelName string) (res modelCardResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ModelCard", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ModelCard(ctx, modelName)
}

func (t *tracing) ModelTree(ctx context.Context, modelName, catalog string) (res modelTreeResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ModelTree", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "catalog", catalog, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ModelTree(ctx, modelName, catalog)
}

func (t *tracing) ModelInfo(ctx context.Context, modelName string) (res modelInfoResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ModelInfo", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ModelInfo(ctx, modelName)
}

func (t *tracing) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ChatCompletionStream(ctx, request)
}

func (t *tracing) GetModelLogs(ctx context.Context, modelName, containerName string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModelLogs", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "containerName", containerName, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetModelLogs(ctx, modelName, containerName)
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
