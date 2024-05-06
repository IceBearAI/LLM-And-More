package chat

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"strings"
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

func (s *openAIService) getClient() *openai.Client {
	// 随机取一个endpoint
	totalEp := len(s.options.endpoints)
	ep := s.options.endpoints[rand.Intn(totalEp)]
	if !strings.HasSuffix(ep.Host, "/v1") {
		ep.Host += "/v1"
	}
	var config openai.ClientConfig
	if ep.Platform == "openai" {
		config = openai.DefaultConfig(ep.Token)
		config.BaseURL = ep.Host
	}
	if ep.Platform == "azure" {
		config = openai.DefaultAzureConfig(ep.Token, ep.Host)
		//config.AzureModelMapperFunc = func(model string) string {
		//	azureModelMapping := map[string]string{
		//		"gpt-3.5-turbo": "your gpt-3.5-turbo deployment name",
		//	}
		//	return azureModelMapping[model]
		//}
	}
	return openai.NewClientWithConfig(config)
}

func (s *openAIService) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, status int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *openAIService) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, status int, err error) {
	status = 200
	client := s.getClient()
	stream, err = client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = er
		}
		return nil, status, err
	}
	return stream, status, nil
}

func (s *openAIService) Models(ctx context.Context) (res []openai.Model, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *openAIService) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

// NewOpenAI creates a new OpenAI service.
func NewOpenAI(opts ...CreationOption) OpenAIService {
	options := &CreationOptions{
		endpoints: []Endpoint{
			{
				Platform: "openai",
				Host:     "https://api.openai.com",
				Token:    "",
			},
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &openAIService{
		options: options,
	}
}
