package model

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindByModelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindByModelId(ctx, modelId, preloads...)
}

func (t *tracing) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindDeployPendingModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindDeployPendingModels(ctx)
}

func (t *tracing) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateDeployStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "status", status, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateDeployStatus(ctx, modelId, status)
}

func (t *tracing) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "SetModelEnabled", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "enabled", enabled, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.SetModelEnabled(ctx, modelId, enabled)
}

func (t *tracing) FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindModelDeployByModeId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindModelDeployByModeId(ctx, modelId)
}

func (t *tracing) SaveModel(ctx context.Context, model *types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "SaveModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("model", model, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.SaveModel(ctx, model)
}

func (t *tracing) CancelModelDeploy(ctx context.Context, modelId uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CancelModelDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CancelModelDeploy(ctx, modelId)
}

func (t *tracing) SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "SaveModelDeploy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.SaveModelDeploy(ctx, data)
}

func (t *tracing) DeleteEval(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteEval(ctx, id)
}

func (t *tracing) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateEval(ctx, data)
}

func (t *tracing) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListEval(ctx, request)
}

func (t *tracing) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", fmt.Sprintf("%+v", data), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateEval(ctx, data)
}

func (t *tracing) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetEval", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetEval(ctx, id)
}

func (t *tracing) ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListModels(ctx, request)
}

func (t *tracing) CreateModel(ctx context.Context, data *types.Models) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateModel(ctx, data)
}

func (t *tracing) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetModel(ctx, id, preload...)
}

func (t *tracing) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateModel(ctx, request)
}

func (t *tracing) DeleteModel(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteModel(ctx, id)
}

func (t *tracing) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindModelsByTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindModelsByTenantId(ctx, tenantId)
}

func (t *tracing) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetModelByModelName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.model",
	})
	defer func() {
		span.LogKV("modelName", modelName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetModelByModelName(ctx, modelName)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
