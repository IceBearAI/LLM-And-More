package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/tests"
	"github.com/go-kit/kit/endpoint"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func initEndpoint() Service {
	apiSvc, _ := tests.Init()
	return New(tests.Logger, "traceId", tests.Store, apiSvc)
}

func TestEndpoint_Embeddings(t *testing.T) {
	svc := initEndpoint()
	ep := MakeEndpoints(svc, map[string][]endpoint.Middleware{})
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyChannelId, uint(1))

	res, err := ep.EmbeddingsEndpoint(ctx, openai.EmbeddingRequest{
		//Input: [][]int{
		//	{161,
		//		102,
		//		112,
		//		161,
		//		226,
		//		123,
		//		16205,
		//		227,
		//		72406,
		//		253,
		//		11883,
		//		21990,
		//		44659,
		//		223,
		//		21007,
		//		227,
		//		49691,
		//		248,
		//		98806,
		//		21043,
		//		163,
		//		228,
		//		253,
		//		44659,
		//		223,
		//		21007,
		//		227,
		//		49691,
		//		248},
		//},
		Input:          "您好！",
		Model:          "llama-3-8b-instruct",
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
	})
	if err != nil {
		t.Error(err)
		return
	}
	resp, ok := res.(encode.Response)
	if !ok {
		t.Error("response type not match")
		return
	}
	t.Log(resp.Data)
}
