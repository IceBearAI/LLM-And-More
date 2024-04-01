package datasettask

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) GenerationAnnotationContent(ctx context.Context, tenantId uint, modelName, taskId, taskSegmentId string) (res taskSegmentAnnotationRequest, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *logging) GetCheckTaskDatasetSimilarLog(ctx context.Context, tenantId uint, taskId string) (res string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetCheckTaskDatasetSimilarLog", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetCheckTaskDatasetSimilarLog(ctx, tenantId, taskId)
}

func (s *logging) CancelCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CancelCheckTaskDatasetSimilar", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CancelCheckTaskDatasetSimilar(ctx, tenantId, taskId)
}

func (s *logging) GetTaskInfo(ctx context.Context, tenantId uint, taskId string) (res taskDetail, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskInfo", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskInfo(ctx, tenantId, taskId)
}

func (s *logging) CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTask", "tenantId", tenantId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateTask(ctx, tenantId, req)
}

func (s *logging) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTasks", "tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListTasks(ctx, tenantId, name, page, pageSize)
}

func (s *logging) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteTask", "tenantId", tenantId, "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteTask(ctx, tenantId, uuid)
}

func (s *logging) GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskSegmentNext", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskSegmentNext(ctx, tenantId, taskId)
}

func (s *logging) AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AnnotationTaskSegment", "tenantId", tenantId, "taskId", taskId, "taskSegmentId", taskSegmentId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AnnotationTaskSegment(ctx, tenantId, taskId, taskSegmentId, req)
}

func (s *logging) AbandonTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AbandonTaskSegment", "tenantId", tenantId, "taskId", taskId, "taskSegmentId", taskSegmentId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AbandonTaskSegment(ctx, tenantId, taskId, taskSegmentId)
}

func (s *logging) AsyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AsyncCheckTaskDatasetSimilar", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AsyncCheckTaskDatasetSimilar(ctx, tenantId, taskId)
}

func (s *logging) SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SplitAnnotationDataSegment", "tenantId", tenantId, "taskId", taskId, "req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SplitAnnotationDataSegment(ctx, tenantId, taskId, req)
}

func (s *logging) ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (filePath string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ExportAnnotationData", "tenantId", tenantId, "taskId", taskId, "formatType", formatType, "filePath", filePath,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ExportAnnotationData(ctx, tenantId, taskId, formatType)
}

func (s *logging) DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteAnnotationTask", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteAnnotationTask(ctx, tenantId, taskId)
}

func (s *logging) CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CleanAnnotationTask", "tenantId", tenantId, "taskId", taskId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CleanAnnotationTask(ctx, tenantId, taskId)
}

func (s *logging) TaskDetectFinish(ctx context.Context, tenantId uint, taskId, testReport string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "TaskDetectFinish", "tenantId", tenantId, "taskId", taskId, "testReport", testReport,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.TaskDetectFinish(ctx, tenantId, taskId, testReport)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.datasettask", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
