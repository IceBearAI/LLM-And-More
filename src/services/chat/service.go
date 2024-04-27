package chat

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/sashabaranov/go-openai"
)

type endpoint struct {
	Platform string
	Host     string
	Token    string
}

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	endpoints      []endpoint
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
func WithEndpoints(endpoints ...endpoint) CreationOption {
	return func(o *CreationOptions) {
		o.endpoints = endpoints
	}
}

// Service chat service interface
type Service interface {
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, status int, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, status int, err error)
	// Models 模型列表
	Models(ctx context.Context) (res []openai.Model, err error)
	// Embeddings 创建图片
	Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}
