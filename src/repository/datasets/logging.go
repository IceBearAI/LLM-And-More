package datasets

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) FindSampleByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.DatasetSample, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindSampleByUUID", "uuid", uuid, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.FindSampleByUUID(ctx, uuid, preloads...)
}

func (s *logging) DeleteSampleByUUID(ctx context.Context, uuid []string) (err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteSampleByUUID", "uuid", fmt.Sprintf("%+v", uuid),
			"err", err,
		)
	}()
	return s.next.DeleteSampleByUUID(ctx, uuid)
}

func (s *logging) FindSampleByTitle(ctx context.Context, datasetId uint, title string, preloads ...string) (dataset types.DatasetSample, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindSampleByTitle", "datasetId", datasetId, "title", title, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.FindSampleByTitle(ctx, datasetId, title, preloads...)
}

func (s *logging) List(ctx context.Context, tenantId uint, page, pageSize int, name string, preloads ...string) (datasets []types.Dataset, total int64, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "tenantId", tenantId, "page", page, "pageSize", pageSize, "name", name, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.List(ctx, tenantId, page, pageSize, name, preloads...)
}

func (s *logging) Create(ctx context.Context, dataset *types.Dataset) (err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create", "dataset", fmt.Sprintf("%+v", dataset),
			"err", err,
		)
	}()
	return s.next.Create(ctx, dataset)
}

func (s *logging) Update(ctx context.Context, dataset *types.Dataset) (err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update", "dataset", fmt.Sprintf("%+v", dataset),
			"err", err,
		)
	}()
	return s.next.Update(ctx, dataset)
}

func (s *logging) Delete(ctx context.Context, id uint) (err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "id", id,
			"err", err,
		)
	}()
	return s.next.Delete(ctx, id)
}

func (s *logging) FindById(ctx context.Context, id uint, preloads ...string) (dataset types.Dataset, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindById", "id", id, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.FindById(ctx, id, preloads...)
}

func (s *logging) FindByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.Dataset, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByUUID", "uuid", uuid, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.FindByUUID(ctx, uuid, preloads...)
}

func (s *logging) FindByUUIDAndTenantId(ctx context.Context, uuid string, tenantId uint, preloads ...string) (dataset types.Dataset, err error) {
	defer func() {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByUUIDAndTenantId", "uuid", uuid, "tenantId", tenantId, "preloads", fmt.Sprintf("%+v", preloads),
			"err", err,
		)
	}()
	return s.next.FindByUUIDAndTenantId(ctx, uuid, tenantId, preloads...)
}

func (s *logging) CreateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateSample", "sample", fmt.Sprintf("%+v", sample),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateSample(ctx, sample)
}

func (s *logging) UpdateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateSample", "sample", fmt.Sprintf("%+v", sample),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateSample(ctx, sample)
}

func (s *logging) DeleteSample(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteSample", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteSample(ctx, id)
}

func (s *logging) SampleList(ctx context.Context, datasetId uint, page, pageSize int, title string, preloads ...string) (samples []types.DatasetSample, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SampleList", "datasetId", datasetId, "page", page, "pageSize", pageSize, "title", title, "preloads", fmt.Sprintf("%+v", preloads),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SampleList(ctx, datasetId, page, pageSize, title, preloads...)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.datasets", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
