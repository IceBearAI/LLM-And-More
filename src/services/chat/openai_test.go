package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func initOpenAiSvc() Service {
	return NewOpenAI(WithEndpoints(Endpoint{
		Platform: "openai",
		Host:     "http://localhost:8000/v1",
		Token:    "sk-xxx",
	}))
}

func TestOpenAIService_ChatCompletionStream(t *testing.T) {
	svc := initOpenAiSvc()
	ctx := context.Background()
	stream, err := svc.ChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "您好！你是谁？",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for res := range stream {
		t.Log(res.Choices[0].Delta.Content)
	}
}
