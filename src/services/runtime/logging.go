// Code generated . DO NOT EDIT.
package runtime

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

func (s *logging) CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error) {
	defer func(begin time.Time) {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateDeployment",

			"config", configJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateDeployment(ctx, config)

}

func (s *logging) CreateJob(ctx context.Context, config Config) (jobName string, err error) {
	defer func(begin time.Time) {

		configByte, _ := json.Marshal(config)
		configJson := string(configByte)

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateJob",

			"config", configJson,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.CreateJob(ctx, config)

}

func (s *logging) GetDeploymentLogs(ctx context.Context, deploymentName string) (log string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetDeploymentLogs",

			"deploymentName", deploymentName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetDeploymentLogs(ctx, deploymentName)

}

func (s *logging) GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetDeploymentStatus",

			"deploymentName", deploymentName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetDeploymentStatus(ctx, deploymentName)

}

func (s *logging) GetJobLogs(ctx context.Context, jobName string) (log string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetJobLogs",

			"jobName", jobName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetJobLogs(ctx, jobName)

}

func (s *logging) GetJobStatus(ctx context.Context, jobName string) (status string, err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetJobStatus",

			"jobName", jobName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.GetJobStatus(ctx, jobName)

}

func (s *logging) RemoveDeployment(ctx context.Context, deploymentName string) (err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RemoveDeployment",

			"deploymentName", deploymentName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.RemoveDeployment(ctx, deploymentName)

}

func (s *logging) RemoveJob(ctx context.Context, jobName string) (err error) {
	defer func(begin time.Time) {

		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RemoveJob",

			"jobName", jobName,

			"took", time.Since(begin),

			"err", err,
		)
	}(time.Now())

	return s.next.RemoveJob(ctx, jobName)

}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "module.runtime", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
