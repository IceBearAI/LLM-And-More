package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAIService interface {
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, status int, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, status int, err error)
	// Models 模型列表
	Models(ctx context.Context) (res []openai.Model, err error)
	// Embeddings 创建图片
	Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}

type openAIService struct {
	options *CreationOptions
}

func (o openAIService) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, status int, err error) {
	//TODO implement me
	panic("implement me")
}

func (o openAIService) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, status int, err error) {
	//TODO implement me
	panic("implement me")
}

func (o openAIService) Models(ctx context.Context) (res []openai.Model, err error) {
	//TODO implement me
	panic("implement me")
}

func (o openAIService) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// NewOpenAI creates a new OpenAI service.
func NewOpenAI(opts ...CreationOption) OpenAIService {
	options := &CreationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return &openAIService{
		options: options,
	}
}
