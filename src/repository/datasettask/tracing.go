package datasettask

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) GetSegmentFaqIntentInSegmentId(ctx context.Context, segmentIds []uint, annotationStatus types.DatasetAnnotationStatus, annotationType types.DatasetAnnotationType) (res []types.DatasetAnnotationTaskSegment, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *tracing) GetTaskByDetection(ctx context.Context, status types.DatasetAnnotationStatus, detectionStatus types.DatasetAnnotationDetectionStatus, preload ...string) (res []types.DatasetAnnotationTask, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskByDetection", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("status", status, "detectionStatus", detectionStatus, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskByDetection(ctx, status, detectionStatus, preload...)
}

func (s *tracing) GetTaskSegmentPrev(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus) (res types.DatasetAnnotationTaskSegment, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskSegmentPrev", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("taskId", taskId, "status", status, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskSegmentPrev(ctx, taskId, status)
}

func (s *tracing) GetDatasetDocumentByUUID(ctx context.Context, tenantId uint, uuid string, preload ...string) (res *types.DatasetDocument, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetDatasetDocumentByUUID", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetDatasetDocumentByUUID(ctx, tenantId, uuid, preload...)
}

func (s *tracing) CreateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateTask(ctx, data)
}

func (s *tracing) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (res []types.DatasetAnnotationTask, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTasks", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListTasks(ctx, tenantId, name, page, pageSize, preloads...)
}

func (s *tracing) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteTask(ctx, tenantId, uuid)
}

func (s *tracing) UpdateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateTask(ctx, data)
}

func (s *tracing) GetTask(ctx context.Context, tenantId uint, uuid string, preloads ...string) (res *types.DatasetAnnotationTask, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "preloads", fmt.Sprintf("%v", preloads), "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTask(ctx, tenantId, uuid, preloads...)
}

func (s *tracing) AddTaskSegments(ctx context.Context, data []types.DatasetAnnotationTaskSegment) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddTaskSegments", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddTaskSegments(ctx, data)
}

func (s *tracing) GetTaskOneSegment(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskOneSegment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("taskId", taskId, "status", status, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskOneSegment(ctx, taskId, status, preload...)
}

func (s *tracing) GetTaskSegmentByUUID(ctx context.Context, taskId uint, uuid string, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskSegmentByUUID", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("taskId", taskId, "uuid", uuid, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskSegmentByUUID(ctx, taskId, uuid, preload...)
}

func (s *tracing) UpdateTaskSegment(ctx context.Context, data *types.DatasetAnnotationTaskSegment) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTaskSegment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateTaskSegment(ctx, data)
}

func (s *tracing) GetTaskSegments(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, page, pageSize int, preload ...string) (res []types.DatasetAnnotationTaskSegment, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskSegments", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("taskId", taskId, "status", status, "page", page, "pageSize", pageSize, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskSegments(ctx, taskId, status, page, pageSize, preload...)
}

func (s *tracing) GetTaskSegmentByRand(ctx context.Context, taskId uint, testPercent float64, status types.DatasetAnnotationStatus, segmentType types.DatasetAnnotationSegmentType, preload ...string) (res []types.DatasetAnnotationTaskSegment, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskSegmentByRand", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("taskId", taskId, "testPercent", testPercent, "status", status, "segmentType", segmentType, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskSegmentByRand(ctx, taskId, testPercent, status, segmentType, preload...)
}

func (s *tracing) UpdateTaskSegmentType(ctx context.Context, segmentIds []uint, segmentType types.DatasetAnnotationSegmentType) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateTaskSegmentType", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("segmentIds", segmentIds, "segmentType", segmentType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateTaskSegmentType(ctx, segmentIds, segmentType)
}

func (s *tracing) GetDatasetDocumentSegmentByRange(ctx context.Context, datasetDocumentId uint, start, end int, preload ...string) (res []types.DatasetDocumentSegment, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetDatasetDocumentSegmentByRange", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("datasetDocumentId", datasetDocumentId, "start", start, "end", end, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetDatasetDocumentSegmentByRange(ctx, datasetDocumentId, start, end, preload...)
}

func (s *tracing) DeleteDatasetDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteDatasetDocument", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteDatasetDocument(ctx, tenantId, uuid)
}

func (s *tracing) DeleteDatasetDocumentById(ctx context.Context, id uint, unscoped bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteDatasetDocumentById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("id", id, "unscoped", unscoped, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteDatasetDocumentById(ctx, id, unscoped)
}

func (s *tracing) ListDatasetDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []types.DatasetDocument, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListDatasetDocuments", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListDatasetDocuments(ctx, tenantId, name, page, pageSize)
}

func (s *tracing) CreateDatasetDocument(ctx context.Context, data *types.DatasetDocument) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateDatasetDocument", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateDatasetDocument(ctx, data)
}

func (s *tracing) AddDatasetDocumentSegments(ctx context.Context, data []types.DatasetDocumentSegment) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddDatasetDocumentSegments", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AddDatasetDocumentSegments(ctx, data)
}

func (s *tracing) UpdateDatasetDocumentSegmentCount(ctx context.Context, datasetDocumentId uint, count int) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateDatasetDocumentSegmentCount", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.datasettask",
	})
	defer func() {
		span.LogKV("datasetDocumentId", datasetDocumentId, "count", count, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateDatasetDocumentSegmentCount(ctx, datasetDocumentId, count)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
