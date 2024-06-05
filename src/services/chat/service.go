package chat

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/sashabaranov/go-openai"
)

type Endpoint struct {
	Platform string
	Host     string
	Token    string
}

// CreationOptions is the options for the chat service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	endpoints      []Endpoint
	workerSvc      WorkerService
}

// CreationOption is the option for the chat service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOpts is the option to set the http client options.
func WithHTTPClientOpts(opts ...kithttp.ClientOption) CreationOption {
	return func(o *CreationOptions) {
		o.httpClientOpts = opts
	}
}

// WithEndpoints is the option to set the endpoints.
func WithEndpoints(endpoints ...Endpoint) CreationOption {
	return func(o *CreationOptions) {
		o.endpoints = endpoints
	}
}

// WithWorkerService worker service interface
func WithWorkerService(svc WorkerService) CreationOption {
	return func(o *CreationOptions) {
		o.workerSvc = svc
	}
}

// CompletionStreamResponse 聊天处理流响应
type CompletionStreamResponse struct {
	Usage openai.Usage `json:"usage"`
	openai.ChatCompletionStreamResponse
}

// CompletionResponse 聊天处理响应
type CompletionResponse struct {
	Usage openai.Usage `json:"usage"`
	openai.ChatCompletionResponse
}

// CompletionResult 聊天处理结果
type CompletionResult struct {
	Usage openai.Usage `json:"usage"`
	openai.CompletionResponse
}

// Middleware 中间件
type Middleware func(Service) Service

// Service chat service interface
type Service interface {
	// Completion 生成
	Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error)
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error)
	// Models 模型列表
	Models(ctx context.Context) (res []openai.Model, err error)
	// Embeddings 创建embedding
	Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}
