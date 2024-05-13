package channel

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) ListChannels(ctx context.Context, request ListChannelRequest) (res []types.ChatChannels, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListChannels", "request", fmt.Sprintf("%v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListChannels(ctx, request)
}

func (s *logging) CreateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateChannel", "data", fmt.Sprintf("%v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateChannel(ctx, data)
}

func (s *logging) GetChannel(ctx context.Context, id uint) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetChannel", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetChannel(ctx, id)
}

func (s *logging) UpdateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateChannel", "data", fmt.Sprintf("%v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateChannel(ctx, data)
}

func (s *logging) DeleteChannel(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteChannel", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteChannel(ctx, id)
}

func (s *logging) AddChannelModels(ctx context.Context, channelId uint, models ...*types.Models) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddChannelModels", "channelId", channelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddChannelModels(ctx, channelId, models...)
}

func (s *logging) FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindChannelById", "id", id, "preloads", fmt.Sprintf("%v", preload),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindChannelById(ctx, id, preload...)
}

func (s *logging) RemoveChannelModels(ctx context.Context, channelId uint, models ...types.Models) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RemoveChannelModels", "channelId", channelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.RemoveChannelModels(ctx, channelId, models...)
}

func (s *logging) FindChannelByKey(ctx context.Context, apiKey string, preloads ...string) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindChannelById", "apiKey", apiKey, "preloads", fmt.Sprintf("%v", preloads),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindChannelByKey(ctx, apiKey, preloads...)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.channel", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
