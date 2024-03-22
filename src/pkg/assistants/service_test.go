package assistants

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/tests"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/tmc/langchaingo/llms/openai"
	"net/http"
	"net/http/httputil"
	"testing"
)

func initSvc() Service {
	_, _, err := tests.Init()
	if err != nil {
		panic(err)
	}
	return New(tests.Logger, "traceId", tests.Store, []kithttp.ClientOption{
		kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			return ctx
		}),
		kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
			dump, _ := httputil.DumpResponse(response, true)
			fmt.Println(string(dump))
			return ctx
		}),
	}, []openai.Option{
		openai.WithToken("sk-001"),
		openai.WithBaseURL("http://chat-api:8080/v1"),
	})
}

func TestService_Create(t *testing.T) {
	svc := initSvc()
	_, err := svc.Create(context.Background(), 1, createRequest{
		Name:         "测试智能体",
		Description:  "测试用的",
		ModelName:    "qwen-72b-chat-int4",
		Instructions: "To use a tool, please use the following format:\n\nThought: Do I need to use a tool? Yes\nAction: the action to take, should be one of [calculator, get_weather, get_my_shopping_cart, get_my_order_status, generate_image]\nAction Input: the input to the action\nObservation: the result of the action\n\nWhen you have a response to say to the Human, or if you do not need to use a tool, you MUST use the format:\n\nThought: Do I need to use a tool? No\nAI: [your response here]",
		Metadata:     ``,
	})
	if err != nil {
		t.Error(err)
		return
	}
	return
}

func TestService_AddTool(t *testing.T) {
	svc := initSvc()
	err := svc.AddTool(context.Background(), 1, "assistant-c681ff33-5932-433f-b780-257dd702347c", "tool-bb73b0ea-5123-4b5d-be89-a53adbf9e230")
	if err != nil {
		t.Error(err)
		return
	}
	return
}

func TestService_Playground(t *testing.T) {
	svc := initSvc()
	res, err := svc.Playground(context.Background(), 1, "assistant-1ed5da7d-5a7e-46e8-92aa-b57fd63de6bf", playgroundRequest{
		ModelName: "qwen-72b-chat-int4",
		//ModelName: "gpt-4-1106-preview",
		ToolIds: []string{"tool-09a9abdd-1a86-4f11-aac6-583e26a40f53", "tool-6ec81e83-8172-4feb-9a2f-2c6063032be1"},
		Messages: []message{
			{
				Role:    "user",
				Content: "您好！你是谁？",
			},
			{
				Role:    "assistant",
				Content: "我是智语超级智能体，您可以问我问题，我会尽力回答您的问题。",
			},
			{
				Role:    "user",
				Content: "现在几点了？",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for {
		select {
		case msg := <-res:
			t.Log(msg.FullContent)
		}
	}
}

func TestService_List(t *testing.T) {

}

func TestService_ListTool(t *testing.T) {
	svc := initSvc()
	tools, total, err := svc.ListTool(context.Background(), 1, "assistant-c681ff33-5932-433f-b780-257dd702347c", "", 1, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(total)
	for _, tool := range tools {
		t.Log(tool.Name)
	}
}

func TestService_Update(t *testing.T) {
	svc := initSvc()
	err := svc.Update(context.Background(), 1, "assistant-e21edba2-6717-4630-8d18-119ef2505cb3", updateRequest{
		Name:         "测试用的",
		Description:  "测试",
		ModelName:    "qwen-72b-chat-int4",
		Instructions: "To use a tool, please use the following format:\n\nThought: Do I need to use a tool? Yes\nAction: the action to take, should be one of [calculator, get_weather, get_my_shopping_cart, get_my_order_status, generate_image]\nAction Input: the input to the action\nObservation: the result of the action\n\nWhen you have a response to say to the Human, or if you do not need to use a tool, you MUST use the format:\n\nThought: Do I need to use a tool? No\nAI: [your response here]",
		Metadata:     ``,
		ToolIds: []string{
			"tool-09a9abdd-1a86-4f11-aac6-583e26a40f53",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	return
}
