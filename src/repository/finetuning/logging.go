package finetuning

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) FindFineTuningJobRunning(ctx context.Context, preloads ...string) (jobs []types.FineTuningTrainJob, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningJobRunning",
			"jobs", jobs,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningJobRunning(ctx, preloads...)
}

func (l *logging) FindFineTuningTemplateByType(ctx context.Context, modelName string, templateType types.TemplateType) (tpl types.FineTuningTemplate, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningTemplateByType",
			"modelName", modelName,
			"templateType", templateType,
			"tpl", tpl,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningTemplateByType(ctx, modelName, templateType)
}

func (l *logging) FindFineTuningTemplateByModelType(ctx context.Context, modelName string, templateType types.TemplateType, preloads ...string) (template types.FineTuningTemplate, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningTemplateByModelType",
			"modelName", modelName,
			"templateType", templateType,
			"template", template,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningTemplateByModelType(ctx, modelName, templateType, preloads...)
}

func (l *logging) GetFineTuningJobByModelName(ctx context.Context, modelName string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetFineTuningJobByModelName",
			"modelName", modelName,
			"job", job,
			"err", err,
		)
	}(time.Now())
	return l.next.GetFineTuningJobByModelName(ctx, modelName, preloads...)
}

func (l *logging) ListFineTuningTemplate(ctx context.Context, request ListFineTuningTemplateRequest) (res []types.FineTuningTemplate, total int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListFineTuningTemplate",
			"request", fmt.Sprintf("%+v", request),
			"res", res,
			"total", total,
			"err", err,
		)
	}(time.Now())
	return l.next.ListFineTuningTemplate(ctx, request)
}

func (l *logging) CreateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateFineTuningJob",
			"request", fmt.Sprintf("%+v", job),
			"job", job,
			"err", err,
		)
	}(time.Now())
	return l.next.CreateFineTuningJob(ctx, job)
}

func (l *logging) FindFineTuningTemplateByModel(ctx context.Context, modelName string, preloads ...string) (template types.FineTuningTemplate, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningTemplateByModel",
			"modelName", modelName,
			"template", template,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningTemplateByModel(ctx, modelName, preloads...)
}

func (l *logging) FindFineTuningJobByJobId(ctx context.Context, jobId string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningJobByJobId",
			"jobId", jobId,
			"job", job,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningJobByJobId(ctx, jobId, preloads...)
}

func (l *logging) FindFineTuningJob(ctx context.Context, id uint, preloads ...string) (job types.FineTuningTrainJob, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningJob",
			"id", id,
			"job", job,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningJob(ctx, id, preloads...)
}

func (l *logging) EncodeFineTuningJobTemplate(ctx context.Context, tpl string, job *types.FineTuningTrainJob) (re string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "EncodeFineTuningJobTemplate",
			"tpl", tpl,
			"job", job,
			"re", re,
			"err", err,
		)
	}(time.Now())
	return l.next.EncodeFineTuningJobTemplate(ctx, tpl, job)
}

func (l *logging) FindFineTuningJobLastByStatus(ctx context.Context, status types.TrainStatus, orderBy string, preloads ...string) (jobs types.FineTuningTrainJob, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindFineTuningJobLastByStatus",
			"status", status,
			"orderBy", orderBy,
			"jobs", jobs,
			"err", err,
		)
	}(time.Now())
	return l.next.FindFineTuningJobLastByStatus(ctx, status, orderBy, preloads...)
}

func (l *logging) HasRunningJob(ctx context.Context) (has bool, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "HasRunningJob",
			"has", has,
			"err", err,
		)
	}(time.Now())
	return l.next.HasRunningJob(ctx)
}

func (l *logging) UpdateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UpdateFineTuningJob",
			"job", job,
			"err", err,
		)
	}(time.Now())
	return l.next.UpdateFineTuningJob(ctx, job)
}

func (l *logging) ListFindTuningJob(ctx context.Context, request ListFindTuningJobRequest) (jobs []types.FineTuningTrainJob, total int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListFindTuningJob",
			"request", fmt.Sprintf("%+v", request),
			"jobs", jobs,
			"total", total,
			"err", err,
		)
	}(time.Now())
	return l.next.ListFindTuningJob(ctx, request)
}

func (l *logging) CountFineTuningJobByStatus(ctx context.Context) (res map[string]int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CountFineTuningJobByStatus",
			"res", res,
			"err", err,
		)
	}(time.Now())
	return l.next.CountFineTuningJobByStatus(ctx)
}

func (l *logging) CountFineTuningJobDuration(ctx context.Context) (res int64, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CountFineTuningJobDuration",
			"res", res,
			"err", err,
		)
	}(time.Now())
	return l.next.CountFineTuningJobDuration(ctx)
}

func (l *logging) DeleteFineTuningJob(ctx context.Context, id uint) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteFineTuningJob",
			"id", id,
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteFineTuningJob(ctx, id)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
