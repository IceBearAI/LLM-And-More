package chat

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/services/chat"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sashabaranov/go-openai"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Completion(ctx context.Context, channelId uint, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Completion", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "req", fmt.Sprintf("%+v", req), "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Completion(ctx, channelId, req)
}

func (s *tracing) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletion", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "req", fmt.Sprintf("%+v", req), "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletion(ctx, channelId, req)
}

func (s *tracing) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan chat.CompletionStreamResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "req", fmt.Sprintf("%+v", req), "stream", stream, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletionStream(ctx, channelId, req)
}

func (s *tracing) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Models", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Models(ctx, channelId)
}

func (s *tracing) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Embeddings", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "req", fmt.Sprintf("%+v", req), "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Embeddings(ctx, channelId, req)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
