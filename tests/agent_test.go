package tests

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
	"strings"
	"testing"
)

func TestAgent_Example(t *testing.T) {
	_, _ = Init()
	llm, testErr := openai.New(openai.WithBaseURL("http://chat-api:8080/v1"),
		openai.WithToken(""),
		//openai.WithModel("qwen-14b-chat"),
		openai.WithModel("qwen-72b-chat-int4"), // int8, int4
		//openai.WithModel("gpt-4-1106-preview"),
		openai.WithEmbeddingModel("bge-large-zh"),
	)
	if testErr != nil {
		t.Error("openai.New", testErr)
		return
	}
	agentTools := []tools.Tool{
		tools.Calculator{},
		NewDynamicTool(ToolOptions{
			Name: "get_weather",
			Description: `对于获取天气很有用。
	输入应该是一个城市或地点。`,
			FunctionCall: func(ctx context.Context, input string) (string, error) {
				return fmt.Sprintf("%s的天气是： %s", input, "晴天"), nil
			},
		}),
		NewDynamicTool(ToolOptions{
			Name: "get_my_shopping_cart",
			Description: `对于获取我的购物车列表很有用。
	输入应该是我的UserId。`,
			FunctionCall: func(ctx context.Context, input string) (string, error) {
				//args := strings.Split(input, ",")
				//fmt.Println(args)
				// 根据 toolType 确认应该创建httpClient实例还是DB实例
				// 调用接口获取购物车列表
				// func api 表需要有 http 地址, auth
				// 返回 苹果，香蕉，橘子
				return fmt.Sprintf("您的购物车有： %s", "苹果，香蕉，橘子"), nil
			},
		}),
		NewDynamicTool(ToolOptions{
			Name: "get_my_order_status",
			//			Description: `对于获取我的订单发货状态很有用。
			//输入应该是两个参数，第一个是我的UserId，第二个是订单OrderId。`,
			Description: `对于获取我的订单发货状态很有用。OrderId要求输入。
	第一个输入应该是我的的UserId
	第二个输入应该是订单的OrderId。`,
			FunctionCall: func(ctx context.Context, input string) (string, error) {
				args := strings.Split(input, ",")
				if len(args) == 1 {
					args = strings.Split(input, " ")
				}
				//fmt.Println(args)
				//fmt.Println(len(args))
				if len(args) != 2 {
					return "抱歉订单编号不正确，请输入新的订单编号", nil
				}
				if !strings.HasPrefix(strings.TrimSpace(args[1]), "SF") {
					return "抱歉订单编号不正确，请输入新的订单编号", nil
				}
				return fmt.Sprintf("%s的订单发送状态是： %s", args[0], "运输中，已经到达了南昌市"), nil
			},
		}),
		NewDynamicTool(ToolOptions{
			Name: "generate_image",
			Description: `对于生成一张图片很有帮助。
	输入应该是相要生成的图片的描述。`,
			FunctionCall: func(ctx context.Context, input string) (string, error) {
				return fmt.Sprintf("生成好的图片地址: %s", "http://localhost:8080/s/BPN_1BKcSg0sQ/origin.png"), nil
			},
		}),
	}

	var previousMessages []schema.ChatMessage
	previousMessages = append(previousMessages,
		schema.SystemChatMessage{Content: "你是智语AI助手！我的基本信息是\nUsername:王聪\nUserID:10001"},
		schema.HumanChatMessage{Content: "您好！"},
		schema.AIChatMessage{Content: "您好！请问有什么可以帮助你的吗？"},
		schema.HumanChatMessage{Content: "今天天气怎么样？"},
		schema.AIChatMessage{Content: "很抱歉，我无法确定您所询问的具体位置。请提供一个具体的城市或地点，我才能帮您查询天气。"},
		//schema.AIChatMessage{Content: "很抱歉，我无法提供具体的地点天气预报，因为我没有收到具体的地点信息。如果您提供一个具体的城市或地区，我可以帮您查询明天的天气。您想查询哪个地方的天气呢？"},
	)

	executor, testErr := agents.Initialize(
		llm,
		agentTools,
		agents.ConversationalReactDescription,
		agents.WithMemory(
			memory.NewConversationBuffer(
				memory.WithChatHistory(
					memory.NewChatMessageHistory(
						memory.WithPreviousMessages(previousMessages),
					),
				),
			),
		),
		//agents.WithReturnIntermediateSteps(),
		//agents.WithCallbacksHandler(callbacks.StreamLogHandler{}),
		//agents.ZeroShotReactDescription,
		agents.WithMaxIterations(3),
		agents.WithParserErrorHandler(agents.NewParserErrorHandler(func(s string) string {
			fmt.Println("ParserErrorHandler", s)
			return s
		})),
	)
	if testErr != nil {
		t.Error("agents.Initialize", testErr)
		return
	}
	var chainsCalls []chains.ChainCallOption
	chainsCalls = append(chainsCalls,
		chains.WithTemperature(0),
		chains.WithTopP(0),
		chains.WithMaxTokens(1024),
		//chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		//	t.Log(string(chunk))
		//	return nil
		//}),
	)
	var answer string
	ctx := context.Background()
	//question := "明天天气怎么样"
	//answer, testErr = chains.Run(ctx, executor, question, chainsCalls...)
	//if testErr != nil {
	//	t.Error("chains.Run", testErr)
	//	return
	//}
	//t.Log(answer)
	t.Log("北京")
	answer, testErr = chains.Run(ctx, executor, "北京", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	//answer, testErr = chains.Run(ctx, executor, "88*2+14-2=?", chainsCalls...)
	//if testErr != nil {
	//	t.Error("chains.Run", testErr)
	//	return
	//}
	//t.Log(answer)
	t.Log("室外应该穿什么衣服？")
	answer, testErr = chains.Run(ctx, executor, "室外应该穿什么衣服？", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	t.Log("看一下我的购物车")
	answer, testErr = chains.Run(ctx, executor, "看一下我的购物车", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	t.Log("帮我查一下订单发货状态")
	answer, testErr = chains.Run(ctx, executor, "帮我查一下订单发货状态", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	answer, testErr = chains.Run(ctx, executor, "SF202401142231938420", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	answer, testErr = chains.Run(ctx, executor, "帮我画一张小白兔的吃草的图", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
	answer, testErr = chains.Run(ctx, executor, "谢谢", chainsCalls...)
	if testErr != nil {
		t.Error("chains.Run", testErr)
		return
	}
	t.Log(answer)
}
