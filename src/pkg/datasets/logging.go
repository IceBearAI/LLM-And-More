package datasets

import (
	"context"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) ExportSample(ctx context.Context, tenantId uint, datasetId, format string) (samples []addSampleRequest, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ExportSample", "tenantId", tenantId, "datasetId", datasetId, "format", format,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ExportSample(ctx, tenantId, datasetId, format)
}

func (s *logging) List(ctx context.Context, tenantId uint, page, pageSize int, query string) (datasets []datasetResult, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "tenantId", tenantId, "page", page, "pageSize", pageSize, "query", query,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, tenantId, page, pageSize, query)
}

func (s *logging) Create(ctx context.Context, tenantId uint, name, remark string) (datasetId string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create", "tenantId", tenantId, "name", name, "remark", remark,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, tenantId, name, remark)
}

func (s *logging) Update(ctx context.Context, tenantId uint, uuid, name, remark string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update", "tenantId", tenantId, "uuid", uuid, "name", name, "remark", remark,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, tenantId, uuid, name, remark)
}

func (s *logging) Delete(ctx context.Context, tenantId uint, datasetId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "tenantId", tenantId, "datasetId", datasetId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, tenantId, datasetId)
}

func (s *logging) Detail(ctx context.Context, tenantId uint, datasetId string) (dataset datasetResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Detail", "tenantId", tenantId, "datasetId", datasetId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Detail(ctx, tenantId, datasetId)
}

func (s *logging) AddSample(ctx context.Context, tenantId uint, uuid, sampleType string, samples []message) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddSample", "tenantId", tenantId, "uuid", uuid, "sampleType", sampleType, "samples", samples,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddSample(ctx, tenantId, uuid, sampleType, samples)
}

func (s *logging) DeleteSample(ctx context.Context, tenantId uint, datasetId string, sampleUIds []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteSample", "tenantId", tenantId, "datasetId", datasetId, "sampleUIds", sampleUIds,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteSample(ctx, tenantId, datasetId, sampleUIds)
}

func (s *logging) SampleList(ctx context.Context, tenantId uint, datasetId string, page, pageSize int, title string) (samples []datasetSampleResult, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SampleList", "tenantId", tenantId, "datasetId", datasetId, "page", page, "pageSize", pageSize, "title", title,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SampleList(ctx, tenantId, datasetId, page, pageSize, title)
}

func (s *logging) UpdateSampleMessages(ctx context.Context, tenantId uint, datasetId string, sampleUId string, messages []message) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateSampleMessages", "tenantId", tenantId, "datasetId", datasetId, "sampleUId", sampleUId, "messages", messages,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateSampleMessages(ctx, tenantId, datasetId, sampleUId, messages)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.datasets", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
