package datasettask

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) GetTaskFAQIntents(ctx context.Context, tenantId uint, taskId string) (res []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *tracing) GenerationAnnotationContent(ctx context.Context, tenantId uint, modelName, taskId, taskSegmentId string) (res taskSegmentAnnotationRequest, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetCheckTaskDatasetSimilarLog", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "modelName", modelName, "taskId", taskId, "taskSegmentId", taskSegmentId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GenerationAnnotationContent(ctx, tenantId, modelName, taskId, taskSegmentId)
}

func (s *tracing) GetCheckTaskDatasetSimilarLog(ctx context.Context, tenantId uint, taskId string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetCheckTaskDatasetSimilarLog", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "res", res, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetCheckTaskDatasetSimilarLog(ctx, tenantId, taskId)
}

func (s *tracing) CancelCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CancelCheckTaskDatasetSimilar", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CancelCheckTaskDatasetSimilar(ctx, tenantId, taskId)
}

func (s *tracing) GetTaskInfo(ctx context.Context, tenantId uint, taskId string) (res taskDetail, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskInfo", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskInfo(ctx, tenantId, taskId)
}

func (s *tracing) CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateTask(ctx, tenantId, req)
}

func (s *tracing) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListTasks", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListTasks(ctx, tenantId, name, page, pageSize)
}

func (s *tracing) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteTask(ctx, tenantId, uuid)
}

func (s *tracing) GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetTaskSegmentNext", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetTaskSegmentNext(ctx, tenantId, taskId)
}

func (s *tracing) AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AnnotationTaskSegment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "taskSegmentId", taskSegmentId, "req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AnnotationTaskSegment(ctx, tenantId, taskId, taskSegmentId, req)
}

func (s *tracing) AbandonTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AbandonTaskSegment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "taskSegmentId", taskSegmentId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AbandonTaskSegment(ctx, tenantId, taskId, taskSegmentId)
}

func (s *tracing) AsyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AsyncCheckTaskDatasetSimilar", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AsyncCheckTaskDatasetSimilar(ctx, tenantId, taskId)
}

func (s *tracing) SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SplitAnnotationDataSegment", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SplitAnnotationDataSegment(ctx, tenantId, taskId, req)
}

func (s *tracing) ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (filePath string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ExportAnnotationData", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "formatType", formatType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ExportAnnotationData(ctx, tenantId, taskId, formatType)
}

func (s *tracing) DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteAnnotationTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteAnnotationTask(ctx, tenantId, taskId)
}

func (s *tracing) CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CleanAnnotationTask", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CleanAnnotationTask(ctx, tenantId, taskId)
}

func (s *tracing) TaskDetectFinish(ctx context.Context, tenantId uint, taskId, testReport string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TaskDetectFinish", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.datasettask",
	})
	defer func() {
		span.LogKV("tenantId", tenantId, "taskId", taskId, "testReport", testReport, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.TaskDetectFinish(ctx, tenantId, taskId, testReport)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
