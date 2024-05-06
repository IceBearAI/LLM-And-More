package chat

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type fsChatApiClient struct {
	options      *CreationOptions
	workerClient WorkerService
}

func (s *fsChatApiClient) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, status int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *fsChatApiClient) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, status int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *fsChatApiClient) Models(ctx context.Context) (res []openai.Model, err error) {
	models, err := s.workerClient.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return nil, err
	}
	for _, model := range models {
		res = append(res, openai.Model{
			ID:   model.ID,
			Root: model.Root,
		})
	}
	return
}

func (s *fsChatApiClient) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func NewFsChatApi(opts ...CreationOption) Service {
	options := &CreationOptions{
		endpoints: []Endpoint{
			{
				Host:     "http://localhost:8000/v1",
				Token:    "",
				Platform: "localai",
			},
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &fsChatApiClient{
		options: options,
	}
}
