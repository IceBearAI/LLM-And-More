// Code generated . DO NOT EDIT.
package account

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) CreateAccount(ctx context.Context, req CreateAccountRequest) (res Account, err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAccount",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateAccount(ctx, req)

}


func (s *logging) DeleteAccount(ctx context.Context, id uint) (err error) {

	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteAccount",

			"id", id,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.DeleteAccount(ctx, id)

}



func (s *logging) ListAccount(ctx context.Context, req ListAccountRequest) (list []Account, total int64, err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListAccount",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.ListAccount(ctx, req)

}


func (s *logging) UpdateAccount(ctx context.Context, req UpdateAccountRequest) (err error) {

	defer func(begin time.Time) {

		reqByte, _ := json.Marshal(req)
		reqJson := string(reqByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateAccount",

			"req", reqJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.UpdateAccount(ctx, req)

}


func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.auth", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
