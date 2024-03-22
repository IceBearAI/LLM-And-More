package finetuning

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
)

type Middleware func(Service) Service

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) _createJob(ctx context.Context, tenantId, channelId uint, trainingFileId, model, suffix, validationFile string, epochs int) (res jobResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "_createJob", "tenantId", tenantId, "channelId", channelId, "trainingFileId", trainingFileId, "model", model, "suffix", suffix, "validationFile", validationFile, "epochs", epochs,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next._createJob(ctx, tenantId, channelId, trainingFileId, model, suffix, validationFile, epochs)
}

func (s *logging) _cancelJob(ctx context.Context, channelId uint, fineTuningJob string) (res jobResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "_cancelJob", "channelId", channelId, "fineTuningJob", fineTuningJob,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next._cancelJob(ctx, channelId, fineTuningJob)
}

func (s *logging) UpdateJobFinishedStatus(ctx context.Context, fineTuningJob string, status types.TrainStatus, message string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateJobFinishedStatus", "fineTuningJob", fineTuningJob, "status", status, "message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateJobFinishedStatus(ctx, fineTuningJob, status, message)
}

func (s *logging) RunWaitingTrain(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RunWaitingTrain",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.RunWaitingTrain(ctx)
}

func (s *logging) _createFineTuningJob(ctx context.Context, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "_createFineTuningJob", "jobId", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next._createFineTuningJob(ctx, jobId)
}

func (s *logging) Estimate(ctx context.Context, tenantId uint, request CreateJobRequest) (response EstimateResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Estimate",
			"request", fmt.Sprintf("%+v", request),
			"response", response,
			"tenantId", tenantId,
			"err", err,
		)
	}(time.Now())
	return s.next.Estimate(ctx, tenantId, request)
}

func (s *logging) ListTemplate(ctx context.Context, tenantId uint, request ListTemplateRequest) (response ListTemplateResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTemplate",
			"request", fmt.Sprintf("%+v", request),
			"tenantId", tenantId,
			"response", response,
			"err", err,
		)
	}(time.Now())
	return s.next.ListTemplate(ctx, tenantId, request)
}

func (s *logging) GetJob(ctx context.Context, tenantId uint, jobId string) (response JobResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetJob",
			"jobId", jobId,
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetJob(ctx, tenantId, jobId)
}

func (s *logging) DeleteJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteJob",
			"jobId", jobId,
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteJob(ctx, tenantId, jobId)
}

func (s *logging) DashBoard(ctx context.Context, tenantId uint) (res DashBoardResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DashBoard",
			"tenantId", tenantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DashBoard(ctx, tenantId)
}

func (s *logging) CreateJob(ctx context.Context, tenantId uint, request CreateJobRequest) (response JobResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateJob",
			"tenantId", tenantId,
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateJob(ctx, tenantId, request)
}

func (s *logging) ListJob(ctx context.Context, tenantId uint, request ListJobRequest) (response ListJobResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListJob",
			"tenantId", tenantId,
			"request", fmt.Sprintf("%+v", request),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListJob(ctx, tenantId, request)
}

func (s *logging) CancelJob(ctx context.Context, tenantId uint, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CancelJob",
			"tenantId", tenantId,
			"jobId", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CancelJob(ctx, tenantId, jobId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.finetuning", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
