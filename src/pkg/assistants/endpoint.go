package assistants

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	// createRequest 创建助手请求
	createRequest struct {
		// Name 名称
		Name string `json:"name" validate:"required"`
		// Remark 描述
		Remark string `json:"remark"`
		// ModelName 模型名称
		ModelName string `json:"modelName" validate:"required,max=255,min=1"`
		// Description 助手描述
		Description string `json:"description"`
		// Instructions 助手使用说明
		Instructions string `json:"instructions"`
		// Metadata 助手元数据
		Metadata string `json:"metadata"`
		// Avatar 头像
		Avatar string `json:"avatar"`
	}

	// updateRequest 更新助手请求
	updateRequest struct {
		// Name 名称
		Name string `json:"name" validate:"required,max=255,min=1"`
		// Remark 描述
		Remark string `json:"remark"`
		// ModelName 模型名称
		ModelName string `json:"modelName" validate:"required,max=255,min=1"`
		// Description 助手描述
		Description string `json:"description"`
		// Instructions 助手使用说明
		Instructions string `json:"instructions"`
		// Metadata 助手元数据
		Metadata string `json:"metadata"`
		// ToolIds 工具ID集合
		ToolIds []string `json:"toolIds"`
	}

	// assistantResult 获取助手响应
	assistantResult struct {
		AssistantId string `json:"assistantId"`
		// Name 名称
		Name string `json:"name"`
		// Remark 描述
		Remark string `json:"remark"`
		// ModelName 模型名称
		ModelName string `json:"modelName"`
		// Description 助手描述
		Description string `json:"description"`
		// Instructions 助手使用说明
		Instructions string `json:"instructions"`
		// Metadata 助手元数据
		Metadata  string                `json:"metadata"`
		CreatedAt time.Time             `json:"createdAt"`
		UpdatedAt time.Time             `json:"updatedAt"`
		Tools     []assistantToolResult `json:"tools"`
		// Avatar 头像
		Avatar   string `json:"avatar"`
		Operator string `json:"operator"`
	}

	// addToolRequest 添加工具请求
	addToolRequest struct {
		// ToolId 工具ID
		ToolId string `json:"toolId"`
	}

	// assistantToolResult  助手工具响应
	assistantToolResult struct {
		// ToolId 工具ID
		ToolId string `json:"toolId"`
		// Name 工具名称
		Name string `json:"name"`
		// Description 工具描述
		Description string `json:"description"`
		// ToolType 工具类型
		ToolType string `json:"toolType"`
		// Metadata 工具元数据
		Metadata string `json:"metadata"`
		// CreatedAt 创建时间
		CreatedAt time.Time `json:"createdAt"`
		// UpdatedAt 更新时间
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// message 消息
	message struct {
		Role    string `json:"role" validate:"required,max=9,min=4,oneof=user assistant"`
		Content string `json:"content" validate:"required,max=255,min=1"`
	}

	// playgroundRequest 操场测试对话请求
	playgroundRequest struct {
		// Messages 消息
		Messages []message `json:"messages" validate:"required"`
		// Instructions 助手使用说明
		Instructions string `json:"instructions"`
		// ModelName 模型名称
		ModelName string `json:"modelName" validate:"required"`
		// 工具集IDS
		ToolIds []string `json:"toolIds" validate:"required"`
		// Stream 是否流式
		Stream bool `json:"stream"`
	}

	// playgroundResult 操场测试对话响应
	playgroundResult struct {
		FullContent  string    `json:"fullContent"`
		Content      string    `json:"content"`
		CreatedAt    time.Time `json:"createdAt"`
		FinishReason string    `json:"finishReason"`
		ContentType  string    `json:"contentType"`
		MessageId    string    `json:"messageId"`
	}
	listRequest struct {
		Name     string `json:"name"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	}
)

type Endpoints struct {
	PlaygroundEndpoint endpoint.Endpoint
	CreateEndpoint     endpoint.Endpoint
	UpdateEndpoint     endpoint.Endpoint
	GetEndpoint        endpoint.Endpoint
	ListEndpoint       endpoint.Endpoint
	DeleteEndpoint     endpoint.Endpoint
	AddToolEndpoint    endpoint.Endpoint
	RemoveToolEndpoint endpoint.Endpoint
	ListToolEndpoint   endpoint.Endpoint
	PublishEndpoint    endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		PlaygroundEndpoint: makePlaygroundEndpoint(s),
		CreateEndpoint:     makeCreateEndpoint(s),
		UpdateEndpoint:     makeUpdateEndpoint(s),
		GetEndpoint:        makeGetEndpoint(s),
		ListEndpoint:       makeListEndpoint(s),
		AddToolEndpoint:    makeAddToolEndpoint(s),
		RemoveToolEndpoint: makeRemoveToolEndpoint(s),
		ListToolEndpoint:   makeListToolEndpoint(s),
		DeleteEndpoint:     makeDeleteEndpoint(s),
		PublishEndpoint:    makePublishEndpoint(s),
	}

	for _, m := range mdw["Assistants"] {
		eps.PlaygroundEndpoint = m(eps.PlaygroundEndpoint)
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
		eps.GetEndpoint = m(eps.GetEndpoint)
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.AddToolEndpoint = m(eps.AddToolEndpoint)
		eps.RemoveToolEndpoint = m(eps.RemoveToolEndpoint)
		eps.ListToolEndpoint = m(eps.ListToolEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.PublishEndpoint = m(eps.PublishEndpoint)
	}
	return eps
}

func makePublishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		err = s.Publish(ctx, tenantId, assistantId)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		err = s.Delete(ctx, tenantId, assistantId)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeListToolEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		tools, total, err := s.ListTool(ctx, tenantId, assistantId, req.Name, req.Page, req.PageSize)
		data := make(map[string]interface{})
		data["list"] = tools
		data["total"] = total
		return encode.Response{
			Data:  data,
			Error: err,
		}, err
	}
}

func makeRemoveToolEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		toolId, _ := ctx.Value(contextKeyToolId).(string)
		err = s.RemoveTool(ctx, tenantId, assistantId, toolId)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeAddToolEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addToolRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		err = s.AddTool(ctx, tenantId, assistantId, req.ToolId)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistants, total, err := s.List(ctx, tenantId, req.Name, req.Page, req.PageSize)
		data := make(map[string]interface{})
		data["list"] = assistants
		data["total"] = total
		return encode.Response{
			Data:  data,
			Error: err,
		}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		assistant, err := s.Get(ctx, tenantId, assistantId)
		return encode.Response{
			Data:  assistant,
			Error: err,
		}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		err = s.Update(ctx, tenantId, assistantId, req)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistant, err := s.Create(ctx, tenantId, req)
		return encode.Response{
			Data:  assistant,
			Error: err,
		}, err
	}
}

func makePlaygroundEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		assistantId, _ := ctx.Value(contextKeyAssistantId).(string)
		if assistantId == "" {
			return nil, encode.ErrAssistantNotFound.Error()
		}
		req := request.(playgroundRequest)
		assistants, err := s.Playground(ctx, tenantId, assistantId, req)
		return encode.Response{
			Stream: req.Stream,
			Data:   assistants,
			Error:  err,
		}, err
	}
}
