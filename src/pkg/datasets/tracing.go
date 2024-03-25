package datasets

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) ExportSample(ctx context.Context, tenantId uint, datasetId, format string) (samples []addSampleRequest, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ExportSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "format", format, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ExportSample(ctx, tenantId, datasetId, format)
}

func (s *tracing) List(ctx context.Context, tenantId uint, page, pageSize int, query string) (datasets []datasetResult, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "page", page, "pageSize", pageSize, "query", query, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, tenantId, page, pageSize, query)
}

func (s *tracing) Create(ctx context.Context, tenantId uint, name, remark string) (datasetId string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "remark", remark, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, tenantId, name, remark)
}

func (s *tracing) Update(ctx context.Context, tenantId uint, datasetId, name, remark string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "name", name, "remark", remark, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, tenantId, datasetId, name, remark)
}

func (s *tracing) Delete(ctx context.Context, tenantId uint, datasetId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, tenantId, datasetId)
}

func (s *tracing) Detail(ctx context.Context, tenantId uint, datasetId string) (dataset datasetResult, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Detail", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Detail(ctx, tenantId, datasetId)
}

func (s *tracing) AddSample(ctx context.Context, tenantId uint, uuid, sampleType string, samples []message) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "sampleType", sampleType, "samples", samples, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddSample(ctx, tenantId, uuid, sampleType, samples)
}

func (s *tracing) DeleteSample(ctx context.Context, tenantId uint, datasetId string, sampleUIds []string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "sampleUIds", sampleUIds, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteSample(ctx, tenantId, datasetId, sampleUIds)
}

func (s *tracing) SampleList(ctx context.Context, tenantId uint, datasetId string, page, pageSize int, title string) (samples []datasetSampleResult, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SampleList", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "page", page, "pageSize", pageSize, "title", title, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SampleList(ctx, tenantId, datasetId, page, pageSize, title)
}

func (s *tracing) UpdateSampleMessages(ctx context.Context, tenantId uint, datasetId string, sampleUId string, messages []message) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateSampleMessages", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "datasetId", datasetId, "sampleUId", sampleUId, "messages", messages, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateSampleMessages(ctx, tenantId, datasetId, sampleUId, messages)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
