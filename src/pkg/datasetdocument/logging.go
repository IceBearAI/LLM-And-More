package datasetdocument

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

func (s *logging) ListDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []datasetDocument, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListDocuments", "tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListDocuments(ctx, tenantId, name, page, pageSize)
}

func (s *logging) CreateDocument(ctx context.Context, tenantId uint, data documentCreateRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateDocument", "tenantId", tenantId, "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateDocument(ctx, tenantId, data)
}

func (s *logging) DeleteDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteDocument", "tenantId", tenantId, "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteDocument(ctx, tenantId, uuid)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.datasetdocument", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
