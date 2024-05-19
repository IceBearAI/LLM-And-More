package models

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) ModelCheckpoint(ctx context.Context, modelName string) (res []string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelName", modelName,
			"method", "ModelCheckpoint",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ModelCheckpoint(ctx, modelName)
}

func (l *logging) ModelCard(ctx context.Context, modelName string) (res modelCardResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelName", modelName,
			"method", "ModelCard",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ModelCard(ctx, modelName)
}

func (l *logging) ModelTree(ctx context.Context, modelName, catalog string) (res modelTreeResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelName", modelName,
			"catalog", catalog,
			"method", "ModelTree",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ModelTree(ctx, modelName, catalog)
}

func (l *logging) ModelInfo(ctx context.Context, modelName string) (res modelInfoResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelName", modelName,
			"method", "ModelInfo",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ModelInfo(ctx, modelName)
}

func (l *logging) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ChatCompletionStream",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ChatCompletionStream(ctx, request)
}

func (l *logging) GetModelLogs(ctx context.Context, modelName, containerName string) (res string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelName", modelName,
			"containerName", containerName,
			"method", "GetModelLogs",
			"err", err,
		)
	}(time.Now())
	return l.next.GetModelLogs(ctx, modelName, containerName)
}

func (l *logging) Undeploy(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Undeploy",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Undeploy(ctx, id)
}

func (l *logging) ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListModels",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListModels(ctx, request)
}

func (l *logging) CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateModel",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateModel(ctx, request)
}

func (l *logging) GetModel(ctx context.Context, id uint) (res Model, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetModel(ctx, id)
}

func (l *logging) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateModel",
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateModel(ctx, request)
}

func (l *logging) DeleteModel(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteModel(ctx, id)
}

func (l *logging) Deploy(ctx context.Context, request ModelDeployRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Deploy",
			"id", request,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Deploy(ctx, request)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.models", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
