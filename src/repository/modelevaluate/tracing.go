package modelevaluate

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

func (s *tracing) CountEvaluate(ctx context.Context, modelId int, evalTargetType string, status []string) (res int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CountEvaluate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("modelId", modelId, "evalTargetType", evalTargetType, "status", status, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CountEvaluate(ctx, modelId, evalTargetType, status)
}

func (s *tracing) IsExistFiveByModelId(ctx context.Context, modelId uint) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "IsExistFiveByModelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.IsExistFiveByModelId(ctx, modelId)
}

func (s *tracing) GetByUuid(ctx context.Context, uuid string) (res types.ModelEvaluate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetByUuid", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetByUuid(ctx, uuid)
}

func (s *tracing) ListModelEvaluate(ctx context.Context, page, pageSize int, modelId uint, status string, evalTargetType string) (res []types.ModelEvaluate, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListModelEvaluate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("page", page, "pageSize", pageSize, "modelId", modelId, "status", status, "evalTargetType", evalTargetType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListModelEvaluate(ctx, page, pageSize, modelId, status, evalTargetType)
}

func (s *tracing) FindFiveGraphLastByModelId(ctx context.Context, modelId, evaluateId uint) (res types.ModelEvaluate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindFiveGraphLastByModelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("modelId", modelId, "evaluateId", evaluateId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindFiveGraphLastByModelId(ctx, modelId, evaluateId)
}

func (s *tracing) DeleteById(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteById(ctx, id)
}

func (s *tracing) Save(ctx context.Context, data *types.ModelEvaluate) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func (s *tracing) GetById(ctx context.Context, id uint) (res types.ModelEvaluate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.modelevaluate",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetById(ctx, id)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
