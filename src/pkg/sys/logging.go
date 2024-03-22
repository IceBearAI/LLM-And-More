package sys

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) TemplateList(ctx context.Context, page, pageSize int, name, templateType string) (res []templateListResult, total int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "page", page, "pageSize", pageSize, "name", name, "templateType", templateType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.TemplateList(ctx, page, pageSize, name, templateType)
}

func (l *logging) TemplateCreate(ctx context.Context, req templateCreateRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.TemplateCreate(ctx, req)
}

func (l *logging) TemplateUpdate(ctx context.Context, req templateCreateRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.TemplateUpdate(ctx, req)
}

func (l *logging) TemplateDelete(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.TemplateDelete(ctx, name)
}

func (l *logging) ListAudit(ctx context.Context, request ListAuditRequest) (resp ListAuditResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListAudit", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListAudit(ctx, request)
}

func (l *logging) ListDict(ctx context.Context, request ListDictRequest) (resp ListDictResponse, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListDict", "request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListDict(ctx, request)
}

func (l *logging) CreateDict(ctx context.Context, data CreateDictRequest) (resp Dict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateDict", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateDict(ctx, data)
}

func (l *logging) DictTreeByCode(ctx context.Context, codes []string) (resp []Dict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DictTreeByCode", "codes", fmt.Sprintf("%+v", codes),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DictTreeByCode(ctx, codes)
}

func (l *logging) UpdateDict(ctx context.Context, data UpdateDictRequest) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateDict", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateDict(ctx, data)
}

func (l *logging) DeleteDict(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteDict", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteDict(ctx, id)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "sys", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
