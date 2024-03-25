package modelevaluate

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) EvalFinish(ctx context.Context, req finishRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "EvalFinish", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.EvalFinish(ctx, req)
}

func (s *logging) FiveGraph(ctx context.Context, req fiveGraphRequest) (res1, res2, res3 fiveGraphResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FiveGraph", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FiveGraph(ctx, req)
}

func (s *logging) Create(ctx context.Context, req createRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, req)
}

func (s *logging) Cancel(ctx context.Context, req cancelRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Cancel", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Cancel(ctx, req)
}

func (s *logging) Delete(ctx context.Context, req deleteRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, req)
}

func (s *logging) List(ctx context.Context, req listRequest) (res []listResult, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, req)
}

var _ Service = &logging{}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "modelEvaluate", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
