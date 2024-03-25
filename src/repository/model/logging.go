package model

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"modelId", modelId,
			"method", "FindByModelId",
			"err", err,
		)
	}(time.Now())
	return l.next.FindByModelId(ctx, modelId, preloads...)
}

func (l *logging) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindDeployPendingModels",
			"err", err,
		)
	}(time.Now())
	return l.next.FindDeployPendingModels(ctx)
}

func (l *logging) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateDeployStatus",
			"modelId", modelId,
			"status", status,
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateDeployStatus(ctx, modelId, status)
}

func (l *logging) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SetModelEnabled",
			"modelId", modelId,
			"enabled", enabled,
			"err", err,
		)
	}(time.Now())
	return l.next.SetModelEnabled(ctx, modelId, enabled)
}

func (l *logging) FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindModelDeployByModeId",
			"modelId", modelId,
			"err", err,
		)
	}(time.Now())
	return l.next.FindModelDeployByModeId(ctx, modelId)
}

func (l *logging) SaveModel(ctx context.Context, model *types.Models) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SaveModel",
			"model", model,
			"err", err,
		)
	}(time.Now())
	return l.next.SaveModel(ctx, model)
}

func (l *logging) CancelModelDeploy(ctx context.Context, modelId uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "modelId", modelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CancelModelDeploy(ctx, modelId)
}

func (l *logging) SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SaveModelDeploy(ctx, data)
}

func (l *logging) DeleteEval(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteEval",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteEval(ctx, id)
}

func (l *logging) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateEval",
			"data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateEval(ctx, data)
}

func (l *logging) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetEval",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetEval(ctx, id)
}

func (l *logging) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateEval",
			"data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateEval(ctx, data)
}

func (l *logging) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListEval",
			"request", fmt.Sprintf("%+v", request),
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListEval(ctx, request)
}

func (l *logging) ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error) {
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

func (l *logging) CreateModel(ctx context.Context, data *types.Models) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateModel",
			"data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateModel(ctx, data)
}

func (l *logging) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetModel",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetModel(ctx, id, preload...)
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

func (l *logging) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindModelsByTenantId",
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindModelsByTenantId(ctx, tenantId)
}

func (l *logging) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetModelByModelName",
			"modelName", modelName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetModelByModelName(ctx, modelName)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
