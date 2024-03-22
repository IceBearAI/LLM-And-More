package datasets

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

func (s *tracing) FindSampleByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.DatasetSample, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindSampleByUUID", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("uuid", uuid, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindSampleByUUID(ctx, uuid, preloads...)
}

func (s *tracing) DeleteSampleByUUID(ctx context.Context, uuid []string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteSampleByUUID", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteSampleByUUID(ctx, uuid)
}

func (s *tracing) FindSampleByTitle(ctx context.Context, datasetId uint, title string, preloads ...string) (dataset types.DatasetSample, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindSampleByTitle", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("datasetId", datasetId, "title", title, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindSampleByTitle(ctx, datasetId, title, preloads...)
}

func (s *tracing) List(ctx context.Context, tenantId uint, page, pageSize int, name string, preloads ...string) (datasets []types.Dataset, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "page", page, "pageSize", pageSize, "name", name, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, tenantId, page, pageSize, name, preloads...)
}

func (s *tracing) Create(ctx context.Context, dataset *types.Dataset) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("dataset", dataset, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, dataset)
}

func (s *tracing) Update(ctx context.Context, dataset *types.Dataset) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("dataset", dataset, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, dataset)
}

func (s *tracing) Delete(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, id)
}

func (s *tracing) FindById(ctx context.Context, id uint, preloads ...string) (dataset types.Dataset, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("id", id, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindById(ctx, id, preloads...)
}

func (s *tracing) FindByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.Dataset, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByUUID", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("uuid", uuid, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByUUID(ctx, uuid, preloads...)
}

func (s *tracing) FindByUUIDAndTenantId(ctx context.Context, uuid string, tenantId uint, preloads ...string) (dataset types.Dataset, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByUUIDAndTenantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("uuid", uuid, "tenantId", tenantId, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByUUIDAndTenantId(ctx, uuid, tenantId, preloads...)
}

func (s *tracing) CreateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("sample", sample, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateSample(ctx, sample)
}

func (s *tracing) UpdateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("sample", sample, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateSample(ctx, sample)
}

func (s *tracing) DeleteSample(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteSample", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteSample(ctx, id)
}

func (s *tracing) SampleList(ctx context.Context, datasetId uint, page, pageSize int, title string, preloads ...string) (samples []types.DatasetSample, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SampleList", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasets",
	})
	defer func() {
		span.LogKV("datasetId", datasetId, "page", page, "pageSize", pageSize, "title", title, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SampleList(ctx, datasetId, page, pageSize, title, preloads...)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
