package modelevaluate

import (
	"context"
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

func (s *logging) IsExistFiveByModelId(ctx context.Context, modelId uint) (res bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "IsExistFiveByModelId", "modelId", modelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.IsExistFiveByModelId(ctx, modelId)
}

func (s *logging) GetByUuid(ctx context.Context, uuid string) (res types.ModelEvaluate, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetByUuid", "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetByUuid(ctx, uuid)
}

func (s *logging) ListModelEvaluate(ctx context.Context, page, pageSize int, modelId uint, status string, evalTargetType string) (res []types.ModelEvaluate, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CountEvaluate", "page", page, "pageSize", pageSize, "modelId", modelId, "status", status, "evalTargetType", evalTargetType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListModelEvaluate(ctx, page, pageSize, modelId, status, evalTargetType)
}

func (s *logging) CountEvaluate(ctx context.Context, evalTargetType string, status []string) (res int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CountEvaluate", "evalTargetType", evalTargetType, "status", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CountEvaluate(ctx, evalTargetType, status)
}

func (s *logging) FindFiveGraphLastByModelId(ctx context.Context, modelId, evaluateId uint) (res types.ModelEvaluate, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindFiveGraphLastByModelId", "modelId", modelId, "evaluateId", evaluateId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindFiveGraphLastByModelId(ctx, modelId, evaluateId)
}

func (s *logging) DeleteById(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteById", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteById(ctx, id)
}

func (s *logging) Save(ctx context.Context, data *types.ModelEvaluate) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, data)
}

func (s *logging) GetById(ctx context.Context, id uint) (res types.ModelEvaluate, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetById", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetById(ctx, id)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.modelEvaluate", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
