package channels

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) GetModel(ctx context.Context, modelName string) (res modelInfoResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModel")
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetModel(ctx, modelName)
}

func (t *tracing) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ChatCompletionStream")
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ChatCompletionStream(ctx, request)
}

func (t *tracing) GetChannel(ctx context.Context, id uint) (resp Channel, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetChannel")
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetChannel(ctx, id)
}

func (t *tracing) ListChannelModels(ctx context.Context, request ListChannelModelsRequest) (resp ChannelModelList, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListChannelModels")
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListChannelModels(ctx, request)
}

func (t *tracing) CreateChannel(ctx context.Context, request CreateChannelRequest) (resp Channel, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateChannel")
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CreateChannel(ctx, request)
}

func (t *tracing) UpdateChannel(ctx context.Context, request UpdateChannelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateChannel")
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.UpdateChannel(ctx, request)
}

func (t *tracing) ListChannel(ctx context.Context, request ListChannelRequest) (resp ChannelList, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListChannel")
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListChannel(ctx, request)
}

func (t *tracing) DeleteChannel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteChannel")
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DeleteChannel(ctx, id)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
