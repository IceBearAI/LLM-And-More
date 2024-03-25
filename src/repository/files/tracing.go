package files

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) CreateFile(ctx context.Context, data *types.Files) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.files",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateFile(ctx, data)
}

func (s *tracing) FindFileByFileId(ctx context.Context, fileId string) (res types.Files, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindFileByFileId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindFileByFileId(ctx, fileId)
}

func (s *tracing) FindFileByMd5(ctx context.Context, md5 string) (res types.Files, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindFileByMd5", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.files",
	})
	defer func() {
		span.LogKV("md5", md5, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindFileByMd5(ctx, md5)
}

func (s *tracing) ListFiles(ctx context.Context, channelId uint, purpose string, fileName string, fileType string, page, pageSize int) (res []types.Files, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListFiles", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.files",
	})
	defer func() {
		span.LogKV(
			"channelId", channelId,
			"purpose", purpose,
			"fileName", fileName,
			"fileType", fileType,
			"page", page,
			"pageSize", pageSize,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListFiles(ctx, channelId, purpose, fileName, fileType, page, pageSize)
}

func (s *tracing) DeleteFile(ctx context.Context, fileId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
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
