package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func initOpenAiSvc() Service {
	return NewOpenAI(WithEndpoints(Endpoint{
		Platform: "openai",
		Host:     "http://paas-chat-api.paas.paas.test/v1",
		Token:    "sk-xxx",
	}))
}

func TestOpenAIService_ChatCompletionStream(t *testing.T) {
	svc := initOpenAiSvc()
	ctx := context.Background()
	stream, _, err := svc.ChatCompletionStream(ctx, openai.ChatCompletionRequest{
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
	for {
		res, err := stream.Recv()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res.Choices[0].Delta.Content)
	}
}
