package ldapcli

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

func (s *logging) Connection(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Connection",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Connection(ctx)
}

func (s *logging) Authenticate(ctx context.Context, account string, password string) (bool, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Authenticate", "account", account,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Authenticate(ctx, account, password)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.ldapcli", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
