// Code generated . DO NOT EDIT.
package runtime

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateDeployment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		span.LogKV(
			"config", configJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateDeployment(ctx, config)

}

func (s *tracing) CreateJob(ctx context.Context, config Config) (jobName string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		span.LogKV(
			"config", configJson,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.CreateJob(ctx, config)

}

func (s *tracing) GetDeploymentLogs(ctx context.Context, deploymentName string) (log string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetDeploymentLogs", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"deploymentName", deploymentName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetDeploymentLogs(ctx, deploymentName)

}

func (s *tracing) GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetDeploymentStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"deploymentName", deploymentName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetDeploymentStatus(ctx, deploymentName)

}

func (s *tracing) GetJobLogs(ctx context.Context, jobName string) (log string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetJobLogs", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"jobName", jobName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetJobLogs(ctx, jobName)

}

func (s *tracing) GetJobStatus(ctx context.Context, jobName string) (status string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetJobStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"jobName", jobName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.GetJobStatus(ctx, jobName)

}

func (s *tracing) RemoveDeployment(ctx context.Context, deploymentName string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveDeployment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"deploymentName", deploymentName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.RemoveDeployment(ctx, deploymentName)

}

func (s *tracing) RemoveJob(ctx context.Context, jobName string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "module.runtime",
	})
	defer func() {

		span.LogKV(
			"jobName", jobName,

			"err", err,
		)

		span.SetTag(string(ext.Error), err != nil)

		span.Finish()
	}()

	return s.next.RemoveJob(ctx, jobName)

}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
