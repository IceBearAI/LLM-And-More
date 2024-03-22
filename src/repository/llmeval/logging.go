package llmeval

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"strings"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) GetEarliestWaitingEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetEarliestWaitingEvalTask", "preloads", strings.Join(preloads, ","),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetEarliestWaitingEvalTask(ctx, preloads...)
}

func (s *logging) GetRunningEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetRunningEvalTask", "preloads", strings.Join(preloads, ","),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetRunningEvalTask(ctx, preloads...)
}

func (s *logging) UpdateEvalStatus(ctx context.Context, evalId uint, status types.EvalStatus, errMessage string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateEvalStatus", "evalId", evalId, "status", status, "errMessage", errMessage,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateEvalStatus(ctx, evalId, status, errMessage)
}

func (s *logging) UpdateEvalProgress(ctx context.Context, evalId uint, score, progress float64, status types.EvalStatus) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateEvalProgress", "evalId", evalId, "score", score, "progress", progress, "status", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateEvalProgress(ctx, evalId, score, progress, status)
}

func (s *logging) UpdateEvalStartTime(ctx context.Context, evalId uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateEvalStartTime", "evalId", evalId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateEvalStartTime(ctx, evalId)
}

func (s *logging) UpdateEvalDetail(ctx context.Context, evalId uint, detail string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateEvalDetail", "evalId", evalId, "detail", detail,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateEvalDetail(ctx, evalId, detail)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.llmeval", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
