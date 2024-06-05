package chat

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"math/rand"
	"strings"
)

type openAIService struct {
	options *CreationOptions
}

func (s *openAIService) Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	client := s.getClient()
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			err = er
		}
		return res, err
	}
	return resp, nil
}

func (s *openAIService) getClient() *openai.Client {
	// 随机取一个endpoint
	var openaiEndpoints []Endpoint
	for _, ep := range s.options.endpoints {
		if ep.Platform == "openai" {
			openaiEndpoints = append(openaiEndpoints, ep)
		}
	}
	var ep Endpoint
	totalEp := len(openaiEndpoints)
	if totalEp > 1 {
		ep = openaiEndpoints[rand.Intn(totalEp)]
	} else {
		ep = openaiEndpoints[0]
	}
	if !strings.HasSuffix(ep.Host, "/v1") {
		ep.Host += "/v1"
	}
	var config openai.ClientConfig
	switch ep.Platform {
	case "openai":
		config = openai.DefaultConfig(ep.Token)
		config.BaseURL = ep.Host
	case "azure":
		config = openai.DefaultAzureConfig(ep.Token, ep.Host)
	//config.AzureModelMapperFunc = func(model string) string {
	//	azureModelMapping := map[string]string{
	//		"gpt-3.5-turbo": "your gpt-3.5-turbo deployment name",
	//	}
	//	return azureModelMapping[model]
	//}
	case "localai":
		config = openai.DefaultConfig(ep.Token)
		config.BaseURL = ep.Host
	}
	return openai.NewClientWithConfig(config)
}

func (s *openAIService) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error) {
	client := s.getClient()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			err = er
		}
		return res, err
	}
	res = CompletionResponse{
		ChatCompletionResponse: resp,
		Usage:                  resp.Usage,
	}

	return res, nil
}

func (s *openAIService) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	client := s.getClient()
	resp, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			err = er
		}
		return nil, err
	}
	dot := make(chan CompletionStreamResponse)
	go func() {
		defer close(dot)
		for {
			recv, err := resp.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("stream error", err)
				break
			}
			dot <- CompletionStreamResponse{
				ChatCompletionStreamResponse: recv,
			}
		}
	}()
	return dot, nil
}

func (s *openAIService) Models(ctx context.Context) (res []openai.Model, err error) {
	client := s.getClient()
	models, err := client.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return nil, err
	}
	for _, model := range models.Models {
		res = append(res, model)
	}
	return
}

func (s *openAIService) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	client := s.getClient()
	res, err = client.CreateEmbeddings(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			err = er
		}
		return res, nil
	}
	return res, nil
}

// NewOpenAI creates a new OpenAI service.
func NewOpenAI(opts ...CreationOption) Service {
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
