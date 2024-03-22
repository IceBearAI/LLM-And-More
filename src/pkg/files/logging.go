package files

import (
	"context"
	"github.com/go-kit/log"
	"mime/multipart"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) UploadToStorage(ctx context.Context, file multipart.File, fileType string) (url string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UploadToStorage", "fileType", fileType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UploadToStorage(ctx, file, fileType)
}

func (l *logging) UploadLocal(ctx context.Context, file multipart.File, fileType string) (localFile string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "UploadLocal", "fileType", fileType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UploadLocal(ctx, file, fileType)
}

func (l *logging) UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "", "file", file, "fileType", fileType, "isPublicBucket", isPublicBucket,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.UploadToS3(ctx, file, fileType, isPublicBucket)
}

func (l *logging) CreateFile(ctx context.Context, request FileRequest) (file File, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreateFile",
			"tenantId", request.TenantId,
			"purpose", request.Purpose,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateFile(ctx, request)
}

func (l *logging) ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "ListFiles",
			"tenantId", request.TenantId,
			"purpose", request.Purpose,
			"fileName", request.FileName,
			"fileType", request.FileType,
			"page", request.Page,
			"pageSize", request.PageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.ListFiles(ctx, request)
}

func (l *logging) GetFile(ctx context.Context, fileId string) (file File, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "GetFile",
			"fileId", fileId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.GetFile(ctx, fileId)
}

func (l *logging) DeleteFile(ctx context.Context, fileId string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "DeleteFile",
			"fileId", fileId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.DeleteFile(ctx, fileId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "pkg.files", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
