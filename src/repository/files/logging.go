package files

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) CreateFile(ctx context.Context, data *types.Files) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateFile", "tenantId", data.TenantID, "filename", data.Name, "purpose", data.Purpose, "fileType", data.Type,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateFile(ctx, data)
}

func (s *logging) FindFileByFileId(ctx context.Context, fileId string) (res types.Files, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindFileByFileId", "fileId", fileId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindFileByFileId(ctx, fileId)
}

func (s *logging) FindFileByMd5(ctx context.Context, md5 string) (res types.Files, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindFileByMd5", "md5", md5,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindFileByMd5(ctx, md5)
}

func (s *logging) ListFiles(ctx context.Context, tenantId uint, purpose string, fileName string, fileType string, page, pageSize int) (res []types.Files, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "ListFiles",
			"tenantId", tenantId, "purpose", purpose, "fileName", fileName, "fileType", fileType, "page", page, "pageSize", pageSize, "total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListFiles(ctx, tenantId, purpose, fileName, fileType, page, pageSize)
}

func (s *logging) DeleteFile(ctx context.Context, fileId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "DeleteFile", "fileId", fileId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteFile(ctx, fileId)
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
