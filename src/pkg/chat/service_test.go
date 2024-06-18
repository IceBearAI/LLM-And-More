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

func TestService_ChatEmbeddings(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	res, err := svc.Embeddings(ctx, 1, openai.EmbeddingRequest{
		Input: [][]int{
			{161,
				102,
				112,
				161,
				226,
				123,
				16205,
				227,
				72406,
				253,
				11883,
				21990,
				44659,
				223,
				21007,
				227,
				49691,
				248,
				98806,
				21043,
				163,
				228,
				253,
				44659,
				223,
				21007,
				227,
				49691,
				248},
		},
		//Input: "您好！",
		Model:          "llama-3-8b-instruct",
		EncodingFormat: openai.EmbeddingEncodingFormatBase64,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, emb := range res.Data {
		t.Log(emb.Embedding)
	}
}
