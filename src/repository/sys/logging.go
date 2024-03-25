package sys

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

func (l *logging) ListFineTuningTemplate(ctx context.Context, page, pageSize int, name, templateType string) (res []types.FineTuningTemplate, total int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "page", page, "pageSize", pageSize, "name", name, "templateType", templateType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListFineTuningTemplate(ctx, page, pageSize, name, templateType)
}

func (l *logging) SaveFineTuningTemplate(ctx context.Context, data *types.FineTuningTemplate) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SaveFineTuningTemplate(ctx, data)
}

func (l *logging) DeleteFineTuningTemplate(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteFineTuningTemplate(ctx, name)
}

func (l *logging) IsExistTuningTemplate(ctx context.Context, name string) (isExist bool, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.IsExistTuningTemplate(ctx, name)
}

func (l *logging) GetFineTuningTemplate(ctx context.Context, name string) (res types.FineTuningTemplate, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetFineTuningTemplate(ctx, name)
}

func (l *logging) CreateAudit(ctx context.Context, data *types.SysAudit) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateAudit", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateAudit(ctx, data)
}

func (l *logging) ListAudit(ctx context.Context, request ListAuditRequest) (items []types.SysAudit, total int64, err error) {
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

func (l *logging) ListDict(ctx context.Context, request ListDictRequest) (items []types.SysDict, total int64, err error) {
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

func (l *logging) CreateDict(ctx context.Context, data *types.SysDict) (err error) {
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

func (l *logging) GetDict(ctx context.Context, id uint) (res types.SysDict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetDict", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetDict(ctx, id)
}

func (l *logging) UpdateDict(ctx context.Context, data *types.SysDict, updateChildrenCode bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateDict", "data", fmt.Sprintf("%+v", data), "updateChildrenCode", updateChildrenCode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateDict(ctx, data, updateChildrenCode)
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

func (l *logging) GetDictByCode(ctx context.Context, code string) (res types.SysDict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetDictByCode", "code", code,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetDictByCode(ctx, code)
}

func (l *logging) FindDictTreeByParentId(ctx context.Context, parentId uint, parentDictType ...string) (res []types.SysDict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindDictTreeByParentId", "parentId", parentId, "parentDictType", parentDictType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindDictTreeByParentId(ctx, parentId, parentDictType...)
}

func (l *logging) GetDictByDictValue(ctx context.Context, parentId uint, dictValue string) (res types.SysDict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetDictByDictValue", "code", parentId, "dictValue", dictValue,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetDictByDictValue(ctx, parentId, dictValue)
}

func (l *logging) FindDictTreeByCode(ctx context.Context, code []string) (res []types.SysDict, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindDictTreeByCode", "code", code,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindDictTreeByCode(ctx, code)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.sys", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
