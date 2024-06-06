// Code generated . DO NOT EDIT.
package model

import (
	"context"
	"encoding/json"

	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) CancelModelDeploy(ctx context.Context, modelId uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CancelModelDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"modelId", modelId,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CancelModelDeploy(ctx, modelId)

}

func (s *tracing) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateEval(ctx, data)

}

func (s *tracing) CreateModel(ctx context.Context, data *types.Models) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateModel(ctx, data)

}

func (s *tracing) DeleteEval(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.DeleteEval(ctx, id)

}

func (s *tracing) DeleteModel(ctx context.Context, id uint) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.DeleteModel(ctx, id)

}

func (s *tracing) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByModelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		span.LogKV(
			"modelId", modelId, "preloads", preloadsJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.FindByModelId(ctx, modelId, preloads...)

}

func (s *tracing) FindByModelNames(ctx context.Context, modelNames []string) (models []types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByModelNames", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		modelNamesByte, _ := json.Marshal(modelNames)
		modelNamesJson := string(modelNamesByte)

		span.LogKV(
			"modelNames", modelNamesJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.FindByModelNames(ctx, modelNames)

}

func (s *tracing) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindDeployPendingModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.FindDeployPendingModels(ctx)

}

func (s *tracing) FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindModelDeployByModeId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"modelId", modelId,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.FindModelDeployByModeId(ctx, modelId)

}

func (s *tracing) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindModelsByTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"tenantId", tenantId,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.FindModelsByTenantId(ctx, tenantId)

}

func (s *tracing) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"id", id,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetEval(ctx, id)

}

func (s *tracing) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		preloadByte, _ := json.Marshal(preload)
		preloadJson := string(preloadByte)

		span.LogKV(
			"id", id, "preload", preloadJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetModel(ctx, id, preload...)

}

func (s *tracing) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetModelByModelName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"modelName", modelName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetModelByModelName(ctx, modelName)

}

func (s *tracing) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		span.LogKV(
			"request", requestJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.ListEval(ctx, request)

}

func (s *tracing) ListModels(ctx context.Context, request ListModelRequest, preloads ...string) (res []types.Models, total int64, err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		preloadsByte, _ := json.Marshal(preloads)
		preloadsJson := string(preloadsByte)

		span.LogKV(
			"request", requestJson, "preloads", preloadsJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.ListModels(ctx, request, preloads...)

}

func (s *tracing) SaveModel(ctx context.Context, model *types.Models) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		modelByte, _ := json.Marshal(model)
		modelJson := string(modelByte)

		span.LogKV(
			"model", modelJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.SaveModel(ctx, model)

}

func (s *tracing) SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveModelDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.SaveModelDeploy(ctx, data)

}

func (s *tracing) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SetModelEnabled", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		span.LogKV(
			"modelId", modelId, "enabled", enabled,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.SetModelEnabled(ctx, modelId, enabled)

}

func (s *tracing) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateDeployStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		statusByte, _ := json.Marshal(status)
		statusJson := string(statusByte)

		span.LogKV(
			"modelId", modelId, "status", statusJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateDeployStatus(ctx, modelId, status)

}

func (s *tracing) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		dataByte, _ := json.Marshal(data)
		dataJson := string(dataByte)

		span.LogKV(
			"data", dataJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateEval(ctx, data)

}

func (s *tracing) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {

		requestByte, _ := json.Marshal(request)
		requestJson := string(requestByte)

		span.LogKV(
			"request", requestJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.UpdateModel(ctx, request)

}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
