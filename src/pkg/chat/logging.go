package chat

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/services/chat"
	"github.com/go-kit/log"
	"github.com/sashabaranov/go-openai"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Completion(ctx context.Context, channelId uint, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Completion", "channelId", channelId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Completion(ctx, channelId, req)
}

func (s *logging) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletion", "channelId", channelId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletion(ctx, channelId, req)
}

func (s *logging) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan chat.CompletionStreamResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletionStream", "channelId", channelId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletionStream(ctx, channelId, req)
}

func (s *logging) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Models", "channelId", channelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Models(ctx, channelId)
}

func (s *logging) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Embeddings", "channelId", channelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Embeddings(ctx, channelId, req)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.chat", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
