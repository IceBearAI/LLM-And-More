package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func initFsChatSvc() Service {
	return NewFsChatApi(WithWorkerService(NewFastChatWorker(WithControllerAddress("http://paasgpt.paas.corp/controller"))))
}

func TestOpenAIService_Embeddings(t *testing.T) {
	svc := initFsChatSvc()
	ctx := context.Background()

	embRes, err := svc.Embeddings(ctx, openai.EmbeddingRequest{
		Model: "glm-4-9b-chat",
		Input: []string{"Hello, my name is"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(embRes)
}
