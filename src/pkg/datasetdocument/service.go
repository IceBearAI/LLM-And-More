package datasetdocument

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type Middleware func(Service) Service

// Service is the interface that provides datasetdocument methods.
type Service interface {
	// ListDocuments returns all documents.
	ListDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []datasetDocument, total int64, err error)
	// CreateDocument creates a new document.
	CreateDocument(ctx context.Context, tenantId uint, data documentCreateRequest) (err error)
	// DeleteDocument deletes a document.
	DeleteDocument(ctx context.Context, tenantId uint, uuid string) (err error)
}

type service struct {
	traceId    string
	repository repository.Repository
	logger     log.Logger
}

func (s *service) ListDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []datasetDocument, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	documents, total, err := s.repository.DatasetTask().ListDatasetDocuments(ctx, tenantId, name, page, pageSize)
	if err != nil {
		err = errors.Wrap(err, "list dataset documents failed")
		_ = level.Warn(logger).Log("msg", "list dataset documents failed", "err", err)
		return
	}

	for _, document := range documents {
		res = append(res, datasetDocument{
			UUID:         document.UUID,
			Name:         document.Name,
			Remark:       document.Remark,
			SegmentCount: document.SegmentCount,
			CreatorEmail: document.CreatorEmail,
			FormatType:   document.FormatType,
			SplitType:    document.SplitType,
			SplitMax:     document.SplitMax,
			FileName:     document.FileName,
			CreatedAt:    document.CreatedAt,
		})
	}

	return res, total, err
}

func (s *service) CreateDocument(ctx context.Context, tenantId uint, data documentCreateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	defer data.File.Close()
	fileName := data.FileHeader.Filename

	document := types.DatasetDocument{
		TenantID:     tenantId,
		UUID:         uuid.New().String(),
		Name:         data.Name,
		Remark:       data.Remark,
		SegmentCount: 0,
		CreatorEmail: email,
		FormatType:   data.FormatType,
		SplitType:    data.SplitType,
		SplitMax:     data.SplitMax,
		FileName:     fileName,
	}
	if err = s.repository.DatasetTask().CreateDatasetDocument(ctx, &document); err != nil {
		err = errors.Wrap(err, "create dataset document failed")
		_ = level.Warn(logger).Log("msg", "create dataset document failed", "err", err)
		return
	}

	// 读取文件内容
	fileBytes, err := io.ReadAll(data.File)
	if err != nil {
		err = errors.Wrap(err, "read file failed")
		_ = level.Warn(logger).Log("msg", "read file failed", "err", err)
		return
	}

	// 将文件内容转换为字符串
	fileContent := string(fileBytes)

	// 按 "\n\n" 分隔符切割文件内容
	splitType := util.UnescapeString(data.SplitType)
	parts := strings.Split(fileContent, splitType)

	var documentSegments []types.DatasetDocumentSegment
	// 处理切割后的每部分
	for i, part := range parts {
		// 过滤空字符串
		if len(strings.Fields(part)) == 0 {
			continue
		}
		documentSegments = append(documentSegments, types.DatasetDocumentSegment{
			DatasetDocumentId: document.ID,
			UUID:              fmt.Sprintf("doc-%s", uuid.New().String()),
			SegmentContent:    part,
			WordCount:         len(strings.Fields(part)),
			SerialNumber:      i,
		})
	}

	// 批量插入切割后的每部分
	if err = s.repository.DatasetTask().AddDatasetDocumentSegments(ctx, documentSegments); err != nil {
		err = errors.Wrap(err, "add dataset document segments failed")
		_ = level.Warn(logger).Log("msg", "add dataset document segments failed", "err", err)
		return
	}

	if err = s.repository.DatasetTask().UpdateDatasetDocumentSegmentCount(ctx, document.ID, len(documentSegments)); err != nil {
		err = errors.Wrap(err, "update dataset document segment count failed")
		_ = level.Warn(logger).Log("msg", "update dataset document segment count failed", "err", err)
		return
	}

	return err
}

func (s *service) DeleteDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	document, err := s.repository.DatasetTask().GetDatasetDocumentByUUID(ctx, tenantId, uuid)
	if err != nil {
		err = errors.Wrap(err, "get dataset document failed")
		_ = level.Warn(logger).Log("msg", "get dataset document failed", "err", err)
		return
	}
	if err = s.repository.DatasetTask().DeleteDatasetDocumentById(ctx, document.ID, false); err != nil {
		err = errors.Wrap(err, "delete dataset document failed")
		_ = level.Error(logger).Log("msg", "delete dataset document failed", "err", err)
		return
	}

	return err
}

func New(traceId string, logger log.Logger, repository repository.Repository) Service {
	logger = log.With(logger, "service", "datasetdocument")
	return &service{
		traceId:    traceId,
		repository: repository,
		logger:     logger,
	}
}
