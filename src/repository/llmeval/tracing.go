package llmeval

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

func (s *tracing) GetEarliestWaitingEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetEarliestWaitingEvalTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetEarliestWaitingEvalTask(ctx, preloads...)
}

func (s *tracing) GetRunningEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetRunningEvalTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetRunningEvalTask(ctx, preloads...)
}

func (s *tracing) UpdateEvalStatus(ctx context.Context, evalId uint, status types.EvalStatus, errMessage string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEvalStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("evalId", evalId, "status", status, "errMessage", errMessage, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateEvalStatus(ctx, evalId, status, errMessage)
}

func (s *tracing) UpdateEvalProgress(ctx context.Context, evalId uint, score, progress float64, status types.EvalStatus) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEvalProgress", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("evalId", evalId, "score", score, "progress", progress, "status", status, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateEvalProgress(ctx, evalId, score, progress, status)
}

func (s *tracing) UpdateEvalStartTime(ctx context.Context, evalId uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEvalStartTime", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("evalId", evalId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateEvalStartTime(ctx, evalId)
}

func (s *tracing) UpdateEvalDetail(ctx context.Context, evalId uint, detail string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateEvalDetail", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.llmeval",
	})
	defer func() {
		span.LogKV("evalId", evalId, "detail", detail, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateEvalDetail(ctx, evalId, detail)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
