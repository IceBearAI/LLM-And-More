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

func (s *tracing) _createJob(ctx context.Context, tenantId, channelId uint, trainingFileId, model, suffix, validationFile string, epochs int) (res jobResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "_createJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "channelId", channelId, "trainingFileId", trainingFileId, "model", model, "suffix", suffix, "validationFile", validationFile, "epochs", epochs, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next._createJob(ctx, tenantId, channelId, trainingFileId, model, suffix, validationFile, epochs)
}

func (s *tracing) _cancelJob(ctx context.Context, channelId uint, fineTuningJob string) (res jobResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "_cancelJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("channelId", channelId, "fineTuningJob", fineTuningJob, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next._cancelJob(ctx, channelId, fineTuningJob)
}

func (s *tracing) UpdateJobFinishedStatus(ctx context.Context, fineTuningJob string, status types.TrainStatus, message string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateJobFinishedStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("fineTuningJob", fineTuningJob, "status", status, "message", message, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateJobFinishedStatus(ctx, fineTuningJob, status, message)
}

func (s *tracing) RunWaitingTrain(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RunWaitingTrain", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RunWaitingTrain(ctx)
}

func (s *tracing) _createFineTuningJob(ctx context.Context, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "_createFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next._createFineTuningJob(ctx, jobId)
}

func (s *tracing) Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Estimate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "request", fmt.Sprintf("%+v", request), "response", response, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.Estimate(ctx, tenantId, request)
}

func (s *tracing) ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "request", fmt.Sprintf("%+v", request), "response", response, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.ListTemplate(ctx, tenantId, request)
}

func (s *tracing) GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "jobId", jobId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.GetJob(ctx, tenantId, jobId)
}

func (s *tracing) DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.DeleteJob(ctx, tenantId, jobId)
}

func (s *tracing) DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DashBoard", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.DashBoard(ctx, tenantId)
}

func (s *tracing) CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.CreateJob(ctx, tenantId, request)
}

func (s *tracing) ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("request", fmt.Sprintf("%+v", request), "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.ListJob(ctx, tenantId, request)
}

func (s *tracing) CancelJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CancelJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.finetuning",
	})
	defer func() {
		span.LogKV("jobId", jobId, "tenantId", tenantId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.CancelJob(ctx, tenantId, jobId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
