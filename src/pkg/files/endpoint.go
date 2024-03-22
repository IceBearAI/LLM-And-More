package files

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"mime/multipart"
	"time"
)

type (
	File struct {
		Size       int64     `json:"size"`
		CreatedAt  time.Time `json:"createdAt"`
		FileId     string    `json:"fileId"`
		FileName   string    `json:"filename"`
		Purpose    string    `json:"purpose"`
		FileType   string    `json:"fileType"`
		TenantId   uint      `json:"tenantId"`
		S3Url      string    `json:"s3Url"`
		LineCount  int       `json:"lineCount"`
		TokenCount int       `json:"tokenCount"`
	}
	FileRequest struct {
		TenantId   uint   `json:"tenantId"`
		Purpose    string `json:"purpose"`
		Header     *multipart.FileHeader
		File       multipart.File
		FileType   string `json:"fileType"`
		LineCount  int    `json:"lineCount"`
		TokenCount int    `json:"tokenCount"`
		IsPublic   bool   `json:"isPublic"`
	}

	FileList struct {
		Files []File `json:"list"`
		Total int64  `json:"total"`
	}

	GetFileRequest struct {
		FileId string `json:"fileId"`
	}

	ListFileRequest struct {
		TenantId uint   `json:"tenantId"`
		Purpose  string `json:"purpose"`
		FileName string `json:"filename"`
		FileType string `json:"fileType"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	}

	// Message 用于解析和验证每一行的JSON对象
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// MessagesWrapper 包含多个Message的结构体
	MessagesWrapper struct {
		Messages []Message `json:"messages"`
	}
)

type Endpoints struct {
	CreateFileEndpoint endpoint.Endpoint
	ListFilesEndpoint  endpoint.Endpoint
	GetFileEndpoint    endpoint.Endpoint
	DeleteFileEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateFileEndpoint: makeCreateFileEndpoint(s),
		ListFilesEndpoint:  makeListFilesEndpoint(s),
		GetFileEndpoint:    makeGetFileEndpoint(s),
		DeleteFileEndpoint: makeDeleteFileEndpoint(s),
	}

	for _, m := range dmw["File"] {
		eps.CreateFileEndpoint = m(eps.CreateFileEndpoint)
		eps.ListFilesEndpoint = m(eps.ListFilesEndpoint)
		eps.GetFileEndpoint = m(eps.GetFileEndpoint)
		eps.DeleteFileEndpoint = m(eps.DeleteFileEndpoint)
	}
	return eps
}

func makeCreateFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(FileRequest)
		req.TenantId = tenantId
		file, err := s.CreateFile(ctx, req)
		return encode.Response{Error: err, Data: file}, err
	}

}

func makeListFilesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListFileRequest)
		files, err := s.ListFiles(ctx, req)
		return encode.Response{
			Data:  files,
			Error: err,
		}, err
	}
}

func makeGetFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetFileRequest)
		file, err := s.GetFile(ctx, req.FileId)
		return encode.Response{
			Data:  file,
			Error: err,
		}, err
	}
}

func makeDeleteFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetFileRequest)
		err = s.DeleteFile(ctx, req.FileId)
		return encode.Response{
			Error: err,
		}, err
	}
}
