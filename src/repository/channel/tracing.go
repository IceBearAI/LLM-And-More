package channel

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) ListChannels(ctx context.Context, request ListChannelRequest) (res []types.ChatChannels, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListChannels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListChannels(ctx, request)
}

func (s *tracing) CreateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateChannel(ctx, data)
}

func (s *tracing) GetChannel(ctx context.Context, id uint) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetChannel(ctx, id)
}

func (s *tracing) UpdateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateChannel(ctx, data)
}

func (s *tracing) DeleteChannel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteChannel(ctx, id)
}

func (s *tracing) AddChannelModels(ctx context.Context, channelId uint, models ...*types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddChannelModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("channelId", channelId, "models", models, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddChannelModels(ctx, channelId, models...)
}

func (s *tracing) FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindChannelById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("id", id, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindChannelById(ctx, id, preload...)
}

func (s *tracing) RemoveChannelModels(ctx context.Context, channelId uint, models ...types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveChannelModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("channelId", channelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RemoveChannelModels(ctx, channelId, models...)
}

func (s *tracing) FindChannelByKey(ctx context.Context, apiKey string, preloads ...string) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindChannelById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.channel",
	})
	defer func() {
		span.LogKV("apiKey", apiKey, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindChannelByKey(ctx, apiKey, preloads...)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
