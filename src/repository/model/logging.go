// Code generated . DO NOT EDIT.
package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) CancelModelDeploy(ctx context.Context, modelId uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CancelModelDeploy",

			"modelId", modelId,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CancelModelDeploy(ctx, modelId)

}

func (s *logging) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateEval",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateEval(ctx, data)

}

func (s *logging) CreateModel(ctx context.Context, data *types.Models) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateModel",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateModel(ctx, data)

}

func (s *logging) DeleteEval(ctx context.Context, id uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteEval",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.DeleteEval(ctx, id)

}

func (s *logging) DeleteModel(ctx context.Context, id uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteModel",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.DeleteModel(ctx, id)

}

func (s *logging) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {

	defer func(begin time.Time) {

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByModelId",

			"modelId", modelId,

			"preloads", preloadsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindByModelId(ctx, modelId, preloads...)

}

func (s *logging) FindByModelNames(ctx context.Context, modelNames []string) (models []types.Models, err error) {

	defer func(begin time.Time) {

		modelNamesByte, _ := json.Marshal(modelNames)
		modelNamesJson := string(modelNamesByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByModelNames",

			"modelNames", modelNamesJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindByModelNames(ctx, modelNames)

}

func (s *logging) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindDeployPendingModels",

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindDeployPendingModels(ctx)

}

func (s *logging) FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindModelDeployByModeId",

			"modelId", modelId,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindModelDeployByModeId(ctx, modelId)

}

func (s *logging) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindModelsByTenantId",

			"tenantId", tenantId,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.FindModelsByTenantId(ctx, tenantId)

}

func (s *logging) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetEval",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetEval(ctx, id)

}

func (s *logging) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {

	defer func(begin time.Time) {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetModel",

			"id", id,

			"preload", preloadJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetModel(ctx, id, preload...)

}

func (s *logging) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetModelByModelName",

			"modelName", modelName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetModelByModelName(ctx, modelName)

}

func (s *logging) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {

	defer func(begin time.Time) {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListEval",

			"request", requestJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListEval(ctx, request)

}

func (s *logging) ListModels(ctx context.Context, request ListModelRequest, preloads ...string) (res []types.Models, total int64, err error) {

	defer func(begin time.Time) {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListModels",

			"request", requestJson,

			"preloads", preloadsJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListModels(ctx, request, preloads...)

}

func (s *logging) SaveModel(ctx context.Context, model *types.Models) (err error) {

	defer func(begin time.Time) {

		modelByte, _ := json.Marshal(model)
		modelJson := string(modelByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveModel",

			"model", modelJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.SaveModel(ctx, model)

}

func (s *logging) SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveModelDeploy",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.SaveModelDeploy(ctx, data)

}

func (s *logging) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SetModelEnabled",

			"modelId", modelId,

			"enabled", enabled,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.SetModelEnabled(ctx, modelId, enabled)

}

func (s *logging) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {

	defer func(begin time.Time) {

		statusByte, _ := json.Marshal(status)
		statusJson := string(statusByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateDeployStatus",

			"modelId", modelId,

			"status", statusJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateDeployStatus(ctx, modelId, status)

}

func (s *logging) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {

	defer func(begin time.Time) {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateEval",

			"data", dataJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateEval(ctx, data)

}

func (s *logging) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {

	defer func(begin time.Time) {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateModel",

			"request", requestJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateModel(ctx, request)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.model", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
