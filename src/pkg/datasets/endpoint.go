package datasets

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	// listRequest 数据集列表请求结构
	listRequest struct {
		// Page 页码
		Page int `json:"page"`
		// PageSize 每页数量
		PageSize int `json:"pageSize"`
		// Name 数据集名称
		Name string `json:"name"`
	}

	// createRequest 创建数据集请求结构
	datasetRequest struct {
		// Name 数据集名称
		Name string `json:"name" validate:"required,lt=50,gt=1"`
		// Remark 数据集描述
		Remark string `json:"remark"`
	}

	// 数据集列表返回结构
	datasetResult struct {
		UUID string `json:"uuid"`
		// Name 名称
		Name string `json:"name"`
		// Remark 描述
		Remark string `json:"remark"`
		// Type 数据集类型
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		// Creator 创建人
		Creator string `json:"creator"`
		// Samples 样本数量
		Samples int `json:"samples"`
	}

	// 数据集样本列表返回结构
	datasetSampleResult struct {
		UUID string `json:"uuid"`
		// Title 标题
		Title string `json:"title"`
		// Content 内容
		Conversations string `json:"conversations" validate:"required"`
		// Label 标签
		Label string `json:"label,omitempty"`
		// Remark 备注
		Remark string `json:"remark,omitempty"`
		// Turns 对话轮次
		Turns     int       `json:"turns"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		// CreatorEmail 创建人邮箱
		CreatorEmail string `json:"creatorEmail"`
		// Messages 对话消息
		Messages []message `json:"messages,omitempty" validate:"required"`
	}

	// message 对话消息
	message struct {
		// Role 角色 user,assistant
		Role    string `json:"role" validate:"required,oneof=user assistant"`
		Content string `json:"content" validate:"required,lt=2048,gt=1"`
	}

	// addSampleRequest 添加数据集样本请求结构
	addSampleRequest struct {
		// Messages 对话消息
		Messages []message `json:"messages" validate:"required"`
		// System 系统
		System string `json:"sys,omitempty"`
		// Instruction 意图
		Instruction string `json:"instruction,omitempty"`
		// Input 输入
		Input string `json:"input,omitempty"`
		// Output 输出
		Output string `json:"output,omitempty"`
	}

	// exportSampleRequest 导出数据集样本请求结构
	exportSampleRequest struct {
		Format string `json:"format"`
	}
)

type Endpoints struct {
	ListEndpoint         endpoint.Endpoint
	CreateEndpoint       endpoint.Endpoint
	UpdateEndpoint       endpoint.Endpoint
	DeleteEndpoint       endpoint.Endpoint
	DetailEndpoint       endpoint.Endpoint
	AddSampleEndpoint    endpoint.Endpoint
	DeleteSampleEndpoint endpoint.Endpoint
	SampleListEndpoint   endpoint.Endpoint
	UpdateSampleEndpoint endpoint.Endpoint
	ExportSampleEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:         makeListEndpoint(s),
		CreateEndpoint:       makeCreateEndpoint(s),
		UpdateEndpoint:       makeUpdateEndpoint(s),
		DeleteEndpoint:       makeDeleteEndpoint(s),
		DetailEndpoint:       makeDetailEndpoint(s),
		AddSampleEndpoint:    makeAddSampleEndpoint(s),
		DeleteSampleEndpoint: makeDeleteSampleEndpoint(s),
		SampleListEndpoint:   makeSampleListEndpoint(s),
		UpdateSampleEndpoint: makeUpdateSampleEndpoint(s),
		ExportSampleEndpoint: makeExportSampleEndpoint(s),
	}

	for _, m := range mdw["Dataset"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.DetailEndpoint = m(eps.DetailEndpoint)
		eps.AddSampleEndpoint = m(eps.AddSampleEndpoint)
		eps.DeleteSampleEndpoint = m(eps.DeleteSampleEndpoint)
		eps.SampleListEndpoint = m(eps.SampleListEndpoint)
		eps.UpdateSampleEndpoint = m(eps.UpdateSampleEndpoint)
		eps.ExportSampleEndpoint = m(eps.ExportSampleEndpoint)
	}
	return eps
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		req := request.(datasetRequest)
		err = s.Update(ctx, tenantId, datasetId, req.Name, req.Remark)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		req := request.(datasetRequest)
		datasetId, err := s.Create(ctx, tenantId, req.Name, req.Remark)
		return encode.Response{
			Data: map[string]interface{}{
				"datasetId": datasetId,
			},
			Error: err,
		}, err
	}
}

func makeUpdateSampleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		datasetSampleId, _ := ctx.Value(contextKeyDatasetSampleId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}

		req := request.(addSampleRequest)
		err = s.UpdateSampleMessages(ctx, tenantId, datasetId, datasetSampleId, req.Messages)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		req := request.(listRequest)
		list, total, err := s.List(ctx, tenantId, req.Page, req.PageSize, req.Name)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  list,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		err = s.Delete(ctx, tenantId, datasetId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		dataset, err := s.Detail(ctx, tenantId, datasetId)
		return encode.Response{
			Data:  dataset,
			Error: err,
		}, err
	}
}

func makeAddSampleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		req := request.(addSampleRequest)
		err = s.AddSample(ctx, tenantId, datasetId, string(types.DatasetTypeText), req.Messages)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteSampleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		sampleId, _ := ctx.Value(contextKeyDatasetSampleId).(string)
		if sampleId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		err = s.DeleteSample(ctx, tenantId, datasetId, []string{sampleId})
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeSampleListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		tenantId, _ := middleware.GetTenantId(ctx)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		req := request.(listRequest)
		list, total, err := s.SampleList(ctx, tenantId, datasetId, req.Page, req.PageSize, req.Name)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  list,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeExportSampleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		datasetId, _ := ctx.Value(contextKeyDatasetId).(string)
		if datasetId == "" {
			return nil, encode.ErrDatasetNotFound.Error()
		}
		req := request.(exportSampleRequest)
		list, err := s.ExportSample(ctx, tenantId, datasetId, req.Format)
		return encode.Response{
			Data:  list,
			Error: err,
		}, err
	}
}
