package tools

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	// createRequest 请求参数
	createRequest struct {
		// Name 名称
		Name string `json:"name" validate:"required,max=255,min=1"`
		// Description 描述
		Description string `json:"description" validate:"required,max=5000,min=1"`
		// ToolType 工具类型
		ToolType string `json:"toolType" validate:"required,max=255,min=1,oneof=function retrieval code_interpreter"`
		// Metadata 元数据
		Metadata string `json:"metadata"`
		// Operator 操作人
		Operator string `json:"operator"`
		// Remark 备注
		Remark string `json:"remark"`
	}

	// updateRequest 请求参数
	updateRequest struct {
	}

	// toolResult 工具结果
	toolResult struct {
		ToolId string `json:"toolId"`
		// Name 名称
		Name string `json:"name"`
		// Description 描述
		Description string `json:"description"`
		// ToolType 工具类型
		ToolType string `json:"toolType"`
		// Metadata 元数据
		Metadata string `json:"metadata"`
		// 操作人
		Operator string `json:"operator"`
		// Remark 备注
		Remark string `json:"remark"`
		// Assistants 助手
		Assistants []assistantResult `json:"assistants"`
		// UpdatedAt 更新时间
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// listRequest 请求参数
	listRequest struct {
		name           string
		page, pageSize int
	}

	// assistantResult 助手结果
	assistantResult struct {
		// AssistantId 助手ID
		AssistantId string `json:"assistantId"`
		// Name 名称
		Name string `json:"name"`
		// Description 描述
		Description string `json:"description"`
	}
)

type Endpoints struct {
	ListEndpoint       endpoint.Endpoint
	CreateEndpoint     endpoint.Endpoint
	UpdateEndpoint     endpoint.Endpoint
	DeleteEndpoint     endpoint.Endpoint
	DetailEndpoint     endpoint.Endpoint
	AssistantsEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:       makeListEndpoint(s),
		CreateEndpoint:     makeCreateEndpoint(s),
		UpdateEndpoint:     makeUpdateEndpoint(s),
		DeleteEndpoint:     makeDeleteEndpoint(s),
		DetailEndpoint:     makeDetailEndpoint(s),
		AssistantsEndpoint: makeAssistantsEndpoint(s),
	}

	for _, m := range mdw["Tools"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.DetailEndpoint = m(eps.DetailEndpoint)
		eps.AssistantsEndpoint = m(eps.AssistantsEndpoint)
	}
	return eps
}

func makeAssistantsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		toolId, _ := ctx.Value(contextKeyToolId).(string)
		if toolId == "" {
			return nil, encode.ErrToolNotFound.Error()
		}
		assistants, err := s.Assistants(ctx, tenantId, toolId)
		return encode.Response{
			Data:  assistants,
			Error: err,
		}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		toolId, _ := ctx.Value(contextKeyToolId).(string)
		if toolId == "" {
			return nil, encode.ErrToolNotFound.Error()
		}
		tool, err := s.Get(ctx, tenantId, toolId)
		return encode.Response{
			Data:  tool,
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		toolId, _ := ctx.Value(contextKeyToolId).(string)
		if toolId == "" {
			return nil, encode.ErrToolNotFound.Error()
		}
		err = s.Delete(ctx, tenantId, toolId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		req := request.(createRequest)
		email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
		req.Operator = email
		err = s.Create(ctx, tenantId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		req := request.(listRequest)
		tools, total, err := s.List(ctx, tenantId, req.name, req.page, req.pageSize)
		return encode.Response{
			Data:  map[string]interface{}{"total": total, "list": tools},
			Error: err,
		}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		toolId, _ := ctx.Value(contextKeyToolId).(string)
		if toolId == "" {
			return nil, encode.ErrToolNotFound.Error()
		}
		req := request.(createRequest)
		email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
		req.Operator = email
		err = s.Update(ctx, tenantId, toolId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}
