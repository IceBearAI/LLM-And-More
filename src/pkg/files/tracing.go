package files

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"mime/multipart"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) UploadToStorage(ctx context.Context, file multipart.File, fileType string) (url string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UploadToStorage", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("file", file, "fileType", fileType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UploadToStorage(ctx, file, fileType)
}

func (s *tracing) UploadLocal(ctx context.Context, file multipart.File, fileType string) (localFile string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UploadLocal", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("file", file, "fileType", fileType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UploadLocal(ctx, file, fileType)
}

func (s *tracing) UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("file", file, "fileType", fileType, "isPublicBucket", isPublicBucket, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UploadToS3(ctx, file, fileType, isPublicBucket)
}

func (s *tracing) CreateFile(ctx context.Context, request FileRequest) (file File, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("channelId", request.TenantId, "purpose", request.Purpose, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.CreateFile(ctx, request)
}

func (s *tracing) ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListFiles", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("tenantId", request.TenantId, "purpose", request.Purpose, "fileName", request.FileName, "fileType", request.FileType, "page", request.Page, "pageSize", request.PageSize, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.ListFiles(ctx, request)
}

func (s *tracing) GetFile(ctx context.Context, fileId string) (file File, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.GetFile(ctx, fileId)
}

func (s *tracing) DeleteFile(ctx context.Context, fileId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return s.next.DeleteFile(ctx, fileId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
