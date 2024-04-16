package terminal

import (
	"context"
	"github.com/igm/sockjs-go/v3/sockjs"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) HandleTerminalSession(session sockjs.Session) {
	s.next.HandleTerminalSession(session)
}

func (s *logging) Token(ctx context.Context, tenantId, userId uint, resourceType string, serviceName, containerName string) (res tokenResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Token",
			"tenantId", tenantId,
			"userId", userId,
			"resourceType", resourceType,
			"serviceName", serviceName,
			"containerName", containerName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.next.Token(ctx, tenantId, userId, resourceType, serviceName, containerName)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "terminal", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
