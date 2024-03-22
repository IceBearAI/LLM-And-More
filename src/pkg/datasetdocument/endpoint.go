package datasetdocument

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"mime/multipart"
	"time"
)

type (
	// listRequest 数据集列表请求结构
	listRequest struct {
		name           string
		page, pageSize int
	}

	// documentCreateRequest 创建文档请求结构
	documentCreateRequest struct {
		// Name 数据集名称
		Name string `json:"name" validate:"required,lt=50,gt=1"`
		// Remark 数据集描述
		Remark string `json:"remark"`
		// FormatType 数据集类型
		FormatType string `json:"formatType"`
		// SplitType 切割方式
		SplitType string `json:"splitType"`
		// SplitMax 切割的最大数据块
		SplitMax int `json:"splitMax"`
		// FileHeader 文件头
		FileHeader *multipart.FileHeader `json:"-"`
		// File 文件
		File multipart.File `json:"-"`
	}

	// documentDetail 文档详情返回结构
	datasetDocument struct {
		// UUID 数据集ID
		UUID string `json:"uuid"`
		// Name 数据集名称
		Name string `json:"name"`
		// Remark 数据集描述
		Remark string `json:"remark"`
		// SegmentCount 数据块数量
		SegmentCount int `json:"segmentCount"`
		// CreatorEmail 创建人邮箱
		CreatorEmail string `json:"creatorEmail"`
		// FormatType 数据集类型
		FormatType string `json:"formatType"`
		// SplitType 切割方式
		SplitType string `json:"splitType"`
		// SplitMax 切割的最大数据块
		SplitMax int `json:"splitMax"`
		// FileName 文件名
		FileName string `json:"fileName"`
		// CreatedAt 创建时间
		CreatedAt time.Time `json:"createdAt"`
	}
)

type Endpoints struct {
	CreateDocumentEndpoint endpoint.Endpoint
	ListDocumentsEndpoint  endpoint.Endpoint
	DeleteDocumentEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateDocumentEndpoint: makeCreateDocumentEndpoint(s),
		ListDocumentsEndpoint:  makeListDocumentsEndpoint(s),
		DeleteDocumentEndpoint: makeDeleteDocumentEndpoint(s),
	}

	for _, m := range mdw["DatasetDocument"] {
		eps.CreateDocumentEndpoint = m(eps.CreateDocumentEndpoint)
		eps.ListDocumentsEndpoint = m(eps.ListDocumentsEndpoint)
		eps.DeleteDocumentEndpoint = m(eps.DeleteDocumentEndpoint)
	}
	return eps
}

func makeCreateDocumentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(documentCreateRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		err := s.CreateDocument(ctx, tenantId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListDocumentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		res, total, err := s.ListDocuments(ctx, tenantId, req.name, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeDeleteDocumentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		uuid, _ := ctx.Value(contextKeyDatasetDocumentId).(string)
		err := s.DeleteDocument(ctx, tenantId, uuid)
		return encode.Response{
			Error: err,
		}, err
	}
}
