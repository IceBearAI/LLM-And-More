package datasettask

import (
	"context"
	"fmt"
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

func (s *logging) GetSegmentFaqIntentInSegmentId(ctx context.Context, segmentIds []uint, annotationStatus types.DatasetAnnotationStatus, annotationType types.DatasetAnnotationType) (res []types.DatasetAnnotationTaskSegment, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *logging) GetTaskByDetection(ctx context.Context, status types.DatasetAnnotationStatus, detectionStatus types.DatasetAnnotationDetectionStatus, preload ...string) (res []types.DatasetAnnotationTask, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskByDetection", "status", status, "detectionStatus", detectionStatus, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskByDetection(ctx, status, detectionStatus, preload...)
}

func (s *logging) GetTaskSegmentPrev(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus) (res types.DatasetAnnotationTaskSegment, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskSegmentPrev", "taskId", taskId, "status", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskSegmentPrev(ctx, taskId, status)
}

func (s *logging) GetDatasetDocumentByUUID(ctx context.Context, tenantId uint, uuid string, preload ...string) (res *types.DatasetDocument, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetDatasetDocumentByUUID", "tenantId", tenantId, "uuid", uuid, "preload", fmt.Sprintf("%+v", preload),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetDatasetDocumentByUUID(ctx, tenantId, uuid, preload...)
}

func (s *logging) CreateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateTask", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateTask(ctx, data)
}

func (s *logging) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (res []types.DatasetAnnotationTask, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListTasks", "tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListTasks(ctx, tenantId, name, page, pageSize, preloads...)
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

func (s *logging) UpdateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTask", "data", fmt.Sprintf("%+v", data),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateTask(ctx, data)
}

func (s *logging) GetTask(ctx context.Context, tenantId uint, uuid string, preloads ...string) (res *types.DatasetAnnotationTask, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTask", "tenantId", tenantId, "uuid", uuid, "preloads", fmt.Sprintf("%v", preloads),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTask(ctx, tenantId, uuid, preloads...)
}

func (s *logging) AddTaskSegments(ctx context.Context, data []types.DatasetAnnotationTaskSegment) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddTaskSegments", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddTaskSegments(ctx, data)
}

func (s *logging) GetTaskOneSegment(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskOneSegment", "taskId", taskId, "status", status, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskOneSegment(ctx, taskId, status, preload...)
}

func (s *logging) GetTaskSegmentByUUID(ctx context.Context, taskId uint, uuid string, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskSegmentByUUID", "taskId", taskId, "uuid", uuid, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskSegmentByUUID(ctx, taskId, uuid, preload...)
}

func (s *logging) UpdateTaskSegment(ctx context.Context, data *types.DatasetAnnotationTaskSegment) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTaskSegment", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateTaskSegment(ctx, data)
}

func (s *logging) GetTaskSegments(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, page, pageSize int, preload ...string) (res []types.DatasetAnnotationTaskSegment, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskSegments", "taskId", taskId, "status", status, "page", page, "pageSize", pageSize, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskSegments(ctx, taskId, status, page, pageSize, preload...)
}

func (s *logging) GetTaskSegmentByRand(ctx context.Context, taskId uint, testPercent float64, status types.DatasetAnnotationStatus, segmentType types.DatasetAnnotationSegmentType, preload ...string) (res []types.DatasetAnnotationTaskSegment, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetTaskSegmentByRand", "taskId", taskId, "testPercent", testPercent, "status", status, "segmentType", segmentType, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetTaskSegmentByRand(ctx, taskId, testPercent, status, segmentType, preload...)
}

func (s *logging) UpdateTaskSegmentType(ctx context.Context, segmentIds []uint, segmentType types.DatasetAnnotationSegmentType) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateTaskSegmentType", "segmentIds", segmentIds, "segmentType", segmentType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateTaskSegmentType(ctx, segmentIds, segmentType)
}

func (s *logging) GetDatasetDocumentSegmentByRange(ctx context.Context, datasetDocumentId uint, start, end int, preload ...string) (res []types.DatasetDocumentSegment, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetDatasetDocumentSegmentByRange", "datasetDocumentId", datasetDocumentId, "start", start, "end", end, "preload", preload,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetDatasetDocumentSegmentByRange(ctx, datasetDocumentId, start, end, preload...)
}

func (s *logging) DeleteDatasetDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteDatasetDocument", "tenantId", tenantId, "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteDatasetDocument(ctx, tenantId, uuid)
}

func (s *logging) DeleteDatasetDocumentById(ctx context.Context, id uint, unscoped bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteDatasetDocumentById", "id", id, "unscoped", unscoped,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteDatasetDocumentById(ctx, id, unscoped)
}

func (s *logging) ListDatasetDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []types.DatasetDocument, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListDatasetDocuments", "tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListDatasetDocuments(ctx, tenantId, name, page, pageSize)
}

func (s *logging) CreateDatasetDocument(ctx context.Context, data *types.DatasetDocument) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateDatasetDocument", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateDatasetDocument(ctx, data)
}

func (s *logging) AddDatasetDocumentSegments(ctx context.Context, data []types.DatasetDocumentSegment) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AddDatasetDocumentSegments", "data", data,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AddDatasetDocumentSegments(ctx, data)
}

func (s *logging) UpdateDatasetDocumentSegmentCount(ctx context.Context, datasetDocumentId uint, count int) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateDatasetDocumentSegmentCount", "datasetDocumentId", datasetDocumentId, "count", count,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateDatasetDocumentSegmentCount(ctx, datasetDocumentId, count)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.datasettask", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
