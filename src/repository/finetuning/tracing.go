package finetuning

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

func (t *tracing) FindFineTuningJobRunning(ctx context.Context, preloads ...string) (jobs []types.FineTuningTrainJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningJobRunning", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("jobs", jobs, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningJobRunning(ctx, preloads...)
}

func (t *tracing) FindFineTuningTemplateByType(ctx context.Context, modelName string, templateType types.TemplateType) (tpl types.FineTuningTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningTemplateByType", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("modelName", modelName, "templateType", templateType, "tpl", tpl, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningTemplateByType(ctx, modelName, templateType)
}

func (t *tracing) FindFineTuningTemplateByModelType(ctx context.Context, modelName string, templateType types.TemplateType, preloads ...string) (template types.FineTuningTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningTemplateByModelType", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("modelName", modelName, "templateType", templateType, "template", template, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningTemplateByModelType(ctx, modelName, templateType, preloads...)
}

func (t *tracing) GetFineTuningJobByModelName(ctx context.Context, modelName string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetFineTuningJobByModelName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("modelName", modelName, "job", job, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetFineTuningJobByModelName(ctx, modelName, preloads...)
}

func (t *tracing) ListFineTuningTemplate(ctx context.Context, request ListFineTuningTemplateRequest) (res []types.FineTuningTemplate, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListFineTuningTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "res", res, "total", total, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListFineTuningTemplate(ctx, request)
}

func (t *tracing) CreateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", job), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateFineTuningJob(ctx, job)
}

func (t *tracing) FindFineTuningTemplateByModel(ctx context.Context, modelName string, preloads ...string) (template types.FineTuningTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningTemplateByModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("modelName", modelName, "template", template, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningTemplateByModel(ctx, modelName, preloads...)
}

func (t *tracing) FindFineTuningJobByJobId(ctx context.Context, jobId string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningJobByJobId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "job", job, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningJobByJobId(ctx, jobId, preloads...)
}

func (t *tracing) FindFineTuningJob(ctx context.Context, id uint, preloads ...string) (job types.FineTuningTrainJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningJob(ctx, id, preloads...)
}

func (t *tracing) EncodeFineTuningJobTemplate(ctx context.Context, tpl string, job *types.FineTuningTrainJob) (re string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "EncodeFineTuningJobTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", job), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.EncodeFineTuningJobTemplate(ctx, tpl, job)
}

func (t *tracing) FindFineTuningJobLastByStatus(ctx context.Context, status types.TrainStatus, orderBy string, preloads ...string) (jobs types.FineTuningTrainJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindFineTuningJobLastByStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("status", status, "orderBy", orderBy, "jobs", jobs, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindFineTuningJobLastByStatus(ctx, status, orderBy, preloads...)
}

func (t *tracing) HasRunningJob(ctx context.Context) (has bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "HasRunningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.HasRunningJob(ctx)
}

func (t *tracing) UpdateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", job), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateFineTuningJob(ctx, job)
}

func (t *tracing) ListFindTuningJob(ctx context.Context, request ListFindTuningJobRequest) (jobs []types.FineTuningTrainJob, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListFindTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListFindTuningJob(ctx, request)
}

func (t *tracing) CountFineTuningJobByStatus(ctx context.Context) (res map[string]int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CountFineTuningJobByStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CountFineTuningJobByStatus(ctx)
}

func (t *tracing) CountFineTuningJobDuration(ctx context.Context) (res int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CountFineTuningJobDuration", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CountFineTuningJobDuration(ctx)
}

func (t *tracing) DeleteFineTuningJob(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.finetuning",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteFineTuningJob(ctx, id)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
