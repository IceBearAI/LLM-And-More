package chat

import (
	"context"
	"github.com/IceBearAI/aigc/tests"
	"github.com/go-kit/log"
	"github.com/sashabaranov/go-openai"
	"os"
	"testing"
)

func initSvc() Service {
	_ = os.Setenv("AIGC_DB_DRIVE", "sqlite")
	_ = os.Setenv("AIGC_FSCHAT_CONTROLLER_ADDRESS", "http://fschat-controller:21001")
	_ = os.Setenv("AIGC_ADMIN_SERVER_STORAGE_PATH", "~/go/src/github.com/icowan/LLM-And-More/storage")
	services, err := tests.Init()
	if err != nil {
		panic(err)
	}
	logger := log.NewLogfmtLogger(os.Stdout)
	return New(logger, "traceId", tests.Store, services)
}

func TestService_ChatCompletionStream(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	stream, err := svc.ChatCompletionStream(ctx, 1, openai.ChatCompletionRequest{
		Model: "qwen1.5-0.5b",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "What is the meaning of life?",
			},
		},
		Stream: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for v := range stream {
		t.Log(v.Choices[0].Delta.Content)
	}
}
