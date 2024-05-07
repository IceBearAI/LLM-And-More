package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestMakeHTTPHandler(t *testing.T) {
	t.Skip("Not implemented")
}

func TestHTTP_ChatCompletionStream(t *testing.T) {
	config := openai.DefaultConfig("sk-001")
	config.BaseURL = "http://localhost:8081/v1"
	client := openai.NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
		Model: "qwen1.5-0.5b",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "您好，你是谁?",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for {
		v, err := stream.Recv()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(v.Choices[0].Delta.Content)
	}
}
