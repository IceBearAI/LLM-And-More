package channels

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"github.com/sashabaranov/go-openai"
	"time"
)

type (
	Channel struct {
		Id     uint   `json:"id"`
		Name   string `json:"name"`
		Alias  string `json:"alias"`
		Quota  int    `json:"quota"`
		ApiKey string `json:"apiKey"`
		Email  string `json:"email"`
		Model  struct {
			List []Model `json:"list"`
			Num  int     `json:"num"`
		} `json:"model"`
		ProjectName  string    `json:"projectName"`
		ServiceName  string    `json:"serviceName"`
		Remark       string    `json:"remark"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		LastOperator string    `json:"lastOperator"`
		TenantId     uint      `json:"tenantId"`
	}
	CreateChannelRequest struct {
		Name         string `json:"name"`
		Alias        string `json:"alias" validate:"required"`
		Quota        int    `json:"quota" validate:"required"`
		Email        string `json:"email" validate:"required"`
		ModelId      []uint `json:"modelId" validate:"required"`
		ProjectName  string `json:"projectName"`
		ServiceName  string `json:"serviceName"`
		Remark       string `json:"remark"`
		LastOperator string `json:"lastOperator"`
		TenantId     uint   `json:"tenantId"`
	}
	UpdateChannelRequest struct {
		Id           uint    `json:"id"`
		TenantId     uint    `json:"tenantId"`
		LastOperator string  `json:"lastOperator"`
		ModelId      []uint  `json:"modelId"`
		Name         *string `json:"name"`
		Alias        *string `json:"alias"`
		Quota        *int    `json:"quota"`
		Email        *string `json:"email"`
		ProjectName  *string `json:"projectName"`
		ServiceName  *string `json:"serviceName"`
		Remark       *string `json:"remark"`
	}
	IdRequest struct {
		Id uint `json:"id" validate:"required"`
	}
	ListChannelRequest struct {
		TenantId    uint    `json:"tenantId"`
		Name        *string `json:"name"`
		Alias       *string `json:"alias"`
		Email       *string `json:"email"`
		ProjectName *string `json:"projectName"`
		ServiceName *string `json:"serviceName"`
		Page        int     `json:"page"`
		PageSize    int     `json:"pageSize"`
	}
	ChannelList struct {
		Channels []Channel `json:"list"`
		Total    int64     `json:"total"`
	}

	Model struct {
		Id           uint      `json:"id"`
		ProviderName string    `json:"providerName"`
		ModelName    string    `json:"modelName"`
		ModelType    string    `json:"modelType"`
		MaxTokens    int       `json:"maxTokens"`
		IsPrivate    bool      `json:"isPrivate"`
		IsFineTuning bool      `json:"isFineTuning"`
		Enabled      bool      `json:"enabled"`
		Remark       string    `json:"remark"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	CreateModelRequest struct {
		ProviderName string `json:"providerName" validate:"required"`
		ModelName    string `json:"modelName" validate:"required"`
		ModelType    string `json:"modelType" validate:"required"`
		MaxTokens    int    `json:"maxTokens" validate:"required"`
		IsPrivate    bool   `json:"isPrivate"`
		IsFineTuning bool   `json:"isFineTuning"`
		Enabled      bool   `json:"enabled"`
		Remark       string `json:"remark"`
	}

	UpdateModelRequest struct {
		Id           uint   `json:"id" validate:"required"`
		ModelName    string `json:"modelName"`
		ModelType    string `json:"modelType"`
		MaxTokens    int    `json:"maxTokens"`
		IsPrivate    *bool  `json:"isPrivate"`
		IsFineTuning *bool  `json:"isFineTuning"`
		Enabled      *bool  `json:"enabled"`
	}
	DeleteChannelModelRequest struct {
		Id uint `json:"id" validate:"required"`
	}
	ListModelRequest struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		ModelName string `json:"modelName"`
		Enabled   *bool  `json:"enabled"`
	}
	ModelList struct {
		Models []Model `json:"list"`
		Total  int64   `json:"total"`
	}
	ListChannelModelsRequest struct {
		TenantId uint `json:"tenantId"`
	}
	ChannelModelList struct {
		Models []Model `json:"list"`
		Total  int64   `json:"total"`
	}
	ChatCompletionRequest struct {
		Messages    []openai.ChatCompletionMessage `json:"messages"`
		MaxTokens   int                            `json:"maxTokens"`
		Temperature float32                        `json:"temperature"`
		TopP        float32                        `json:"topP"`
		Model       string                         `json:"model"`
	}

	CompletionsStreamResult struct {
		FullContent  string    `json:"fullContent"`
		Content      string    `json:"content"`
		CreatedAt    time.Time `json:"createdAt"`
		FinishReason string    `json:"finishReason"`
		ContentType  string    `json:"contentType"`
		MessageId    string    `json:"messageId"`
		Model        string    `json:"model"`       // 模型唯一标识
		TopP         float64   `json:"topP"`        // 生成文本的多样性
		Temperature  float64   `json:"temperature"` // 生成文本的多样性
		MaxTokens    int       `json:"maxTokens"`   // 生成文本的最大长度
	}
)

type Endpoints struct {
	CreateChannelEndpoint        endpoint.Endpoint
	ListChannelsEndpoint         endpoint.Endpoint
	UpdateChannelEndpoint        endpoint.Endpoint
	DeleteChannelEndpoint        endpoint.Endpoint
	ListChannelModelsEndpoint    endpoint.Endpoint
	ChatCompletionStreamEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateChannelEndpoint:        makeCreateChannelEndpoint(s),
		ListChannelsEndpoint:         makeListChannelsEndpoint(s),
		UpdateChannelEndpoint:        makeUpdateChannelEndpoint(s),
		DeleteChannelEndpoint:        makeDeleteChannelEndpoint(s),
		ListChannelModelsEndpoint:    makeListChannelModelsEndpoint(s),
		ChatCompletionStreamEndpoint: makeChatCompletionStreamEndpoint(s),
	}

	for _, m := range dmw["Channel"] {
		eps.CreateChannelEndpoint = m(eps.CreateChannelEndpoint)
		eps.ListChannelsEndpoint = m(eps.ListChannelsEndpoint)
		eps.UpdateChannelEndpoint = m(eps.UpdateChannelEndpoint)
		eps.DeleteChannelEndpoint = m(eps.DeleteChannelEndpoint)
		eps.ListChannelModelsEndpoint = m(eps.ListChannelModelsEndpoint)
		eps.ChatCompletionStreamEndpoint = m(eps.ChatCompletionStreamEndpoint)
	}
	return eps
}

func makeCreateChannelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateChannelRequest)
		if email, ok := middleware.GetEmail(ctx); ok {
			req.LastOperator = email
		}
		tenantId, _ := middleware.GetTenantId(ctx)
		req.TenantId = tenantId
		resp, err := s.CreateChannel(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListChannelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListChannelRequest)
		tenantId, _ := middleware.GetTenantId(ctx)
		req.TenantId = tenantId
		resp, err := s.ListChannel(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeUpdateChannelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateChannelRequest)
		if email, ok := middleware.GetEmail(ctx); ok {
			req.LastOperator = email
		}
		err = s.UpdateChannel(ctx, req)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeDeleteChannelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.DeleteChannel(ctx, req.Id)
		return encode.Response{
			Data:  make(map[string]interface{}),
			Error: err,
		}, err
	}
}

func makeListChannelModelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListChannelModelsRequest)
		tenantId, _ := middleware.GetTenantId(ctx)
		req.TenantId = tenantId
		resp, err := s.ListChannelModels(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeChatCompletionStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ChatCompletionRequest)
		resp, err := s.ChatCompletionStream(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}
