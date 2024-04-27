package tests

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

type QA struct {
	Questions []string `json:"questions"`
}

func getOpenAIClient() (*openai.Client, string) {
	config := openai.DefaultConfig("sk-7H3vT1nB4gM9xK2jD6qF8sZ0cW5rY7lP2mN1bV3kX9yG4wJ6")
	config.BaseURL = "http://aigc-server:8080/v1"
	return openai.NewClientWithConfig(config), openai.GPT3Dot5Turbo16K
}

func getLocalOpenAIClient() (*openai.Client, string) {
	config := openai.DefaultConfig("sk-q8oSE4F7ANJPI7L60NBEENAGYXbYdS6J7gPFDPIFx24")
	config.BaseURL = "http://paas-chat-api.paas.paas.test/v1"
	client := openai.NewClientWithConfig(config)
	return client, openai.GPT3Dot5Turbo
}

func TestLocalAiPing(t *testing.T) {
	client, _ := getLocalOpenAIClient()
	var loopNum = 20
	var totalNum = 0
	now := time.Now()
	for {
		_now := time.Now()
		var concurrent = 4
		var wg sync.WaitGroup
		wg.Add(concurrent)
		for {
			totalNum++
			go func() {
				ccr, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
					Model: "qwen1.5-14b-chat",
					Messages: []openai.ChatCompletionMessage{
						{
							Role: "user",
							Content: `Answer the following questions as best you can. You have access to the following tools:

分贝通快捷入口通用工具: 
                本工具提供一个通用的接口来访问分贝通的各种快捷入口功能。
                                
                使用方法很简单，只需按照“分贝通，入参：[你的参数]”的格式提问即可。
                例如，“分贝通，入参：emailPrefix=xxx”，此时，系统会调用fbt_quick_access函数，并将“emailPrefix=xxx”作为参数传入。
                
获取助手形象: 
                当你获取助手形象时很有用。
                例如示例1，若获取助手形象，可以这样提问：
                '助手形象，入参：emailPrefix=xxx'，
                此时，调用query_assistant_images函数，输入参数格式为：emailPrefix=xxx。
                
全系统待办综合查询: 
                亲爱的用户，当您询问‘待审批’时，我将为您提供一个全面的待办事项视图，覆盖所有系统的待办事项。
                这个功能特别适合于您需要全局了解待办事项的场合。
                您无需指定具体系统，只需告诉我您的需求，比如直接说‘待审批’或‘我的待办’，我就会立刻为您展示来自所有系统的待办事项清单。
                我会通过调用query_all_systems_todos函数，系统自带的输入参数为：emailPrefix=xxx，来处理您的查询，以便精确地为您提供服务。
                现在，请告诉我您具体想查询什么类型的信息。我在这里，随时准备为您服务！
                
特定系统待办查询: 
                用于查询指定系统中的待办事项。请在查询时明确提供系统名称及必要的参数信息，并且本工具支持灵活的查询格式。
                例如示例1，若需查询费用管理系统的待办事项，可以这样提问：
                '费用管理待审批，入参：emailPrefix=xxx'，
                此时，调用query_specific_system_todos函数，输入参数格式为：system_name=费用管理&emailPrefix=xxx。
                例如示例2，如果您想查询BPM系统的待办事项，无论是使用
                'bpm待审批，入参：emailPrefix=xxx'还是
                'bpm系统待审批，入参：emailPrefix=xxx'，
                都能够被正确理解和处理，进而调用query_specific_system_todos函数，输入参数格式为：system_name=bpm&emailPrefix=xxx。
                
                【注意事项】
                - 当输入'资管审批流待审批'时对应的system_name为'资管审批'。
                - 当输入'芝麻待审批'时对应的system_name为'芝麻开门'。

                通过这种方式，您可以快速获取任何指定系统的待办事项清单。
                
全系统已审批综合查询: 
                此工具提供了一个快速获取全系统已审批事项的综合视图，特别适合于在没有指定系统且查询所有已完成审批任务的场景。
                
                【核心功能】
                - 汇总查询：集成并展示所有系统中的已审批事项。
                
                【操作指南】
                1. 构造查询请求：明确您的查询意图，例如查询所有系统的已审批事项。
                2. 发起查询：使用'已审批，入参：emailPrefix=xxx'格式发起请求。
                
                【示例查询】
                如果您希望查询所有系统中的已审批事项列表，您可以这样提问：
                '查询全系统已审批事项，入参：emailPrefix=xxx'。
                
                【输入参数说明】
                - emailPrefix: 邮箱前缀。
                
                依照以上步骤，您将能够获得一个全面的已审批任务列表，帮助您全局了解审批状态。
            
特定系统已审批查询: 
                此工具专为查询指定系统中已完成审批的事项而设计。使用时，您需要提供明确的系统名称和相关的身份验证信息，以确保查询的准确性和安全性。
                
                【快速上手】
                - 明确指定需要查询的系统名称，例如“费用管理系统”或“BPM系统”。
                
                【查询示例】
                - 查询费用管理系统已审批事项：
                  '费用管理已审批，入参：emailPrefix=xxx'。
                  调用query_specific_system_dones函数时，参数为：system_name=费用管理&emailPrefix=xxx。
                  
                - 查询BPM系统已审批事项：
                  'bpm已审批，入参：emailPrefix=xxx' 或
                  'bpm系统已审批，入参：emailPrefix=xxx'。
                  这两种表述方式都将被正确处理，调用相同函数，参数为：system_name=bpm&emailPrefix=xxx。
                
                【注意事项】
                - 当输入'资管审批流已审批'时对应的system_name为'资管审批'。
                - 当输入'芝麻已审批'时对应的system_name为'芝麻开门'。

                【输入参数详解】
                - system_name: 指定查询的系统名称，确保名称准确无误。
                - emailPrefix: 邮箱前缀。
                
                依照这些指南，您可以有效地获取指定系统中已审批事项的详细清单，无论是常规审批还是特定流程审批。
            
查询待办详情: 
                 使用"查询待办详情"工具时，简单提及待办类型和人名即可触发查询。为了更精确的操作，建议在询问时使用“查询”这样的动词。
        
                【使用示例1】
                假设您需要查询报销审批的待办详情，您应该这样提问：
                '报销审批详情，入参：emailPrefix=xxx'。
                在此示例中，您将调用query_todo_details函数，并按以下格式提供输入参数：
                todo_type_name=报销审批&emailPrefix=xxx。
                
                【使用示例2】
                假设您需要查询某个人的报销审批的待办详情，您应该这样提问：
                '许赛赛的报销申请，入参：emailPrefix=xxx'。
                在此示例中，您将调用query_todo_details函数，并按以下格式提供输入参数：
                todo_type_name=报销审批&emailPrefix=xxx&userName=许赛赛。
                
                【参数说明】
                - todo_type_name: 待查询的待办类型名称（例如：报销审批、出差审批等）。
                - emailPrefix: 邮箱前缀。
                
                【注意事项】
                - 请确保入参格式正确，包括待办类型名称和所需的身份验证信息。
                - 本工具支持灵活的查询格式，但建议遵循上述示例以确保准确性。
                - 当输入'资管审批流待办详情'时对应的system_name为'资管审批'。
                - 当输入'芝麻待审批详情'时对应的system_name为'芝麻开门'。

                通过遵循上述指南，您可以有效且准确地查询到指定系统中待办事项的详细信息。
                
查询已审批详情: 
                 使用"查询已审批详情"工具时，简单提及已审批类型和人名即可触发查询。为了更精确的操作，建议在询问时使用“查询”这样的动词。
        
                【使用示例1】
                假设您需要查询报销审批的已审批详情，您应该这样提问：
                '报销审批已审批详情，入参：emailPrefix=xxx'。
                在此示例中，您将调用query_dones_details函数，并按以下格式提供输入参数：
                todo_type_name=报销审批&emailPrefix=xxx。
                
                【使用示例2】
                假设您需要查询某个人的报销审批的已审批详情，您应该这样提问：
                '许赛赛的已审批申请，入参：emailPrefix=xxx'。
                在此示例中，您将调用query_dones_details函数，并按以下格式提供输入参数：
                todo_type_name=报销审批&emailPrefix=xxx&userName=许赛赛。
                
                【参数说明】
                - todo_type_name: 待查询的已审批类型名称（例如：报销审批、出差审批、合同续签等）。
                - emailPrefix: 邮箱前缀。

                【注意事项】
                - 请确保入参格式正确，包括已审批类型名称和所需的身份验证信息。
                - 本工具支持灵活的查询格式，但建议遵循上述示例以确保准确性。
                - 当输入'资管审批流已审批详情'时对应的system_name为'资管审批'。
                - 当输入'芝麻已审批详情'时对应的system_name为'芝麻开门'。
                - 当输入'合同续签已审批详情'时对应的system_name为'合同续签'。

                通过遵循上述指南，您可以有效且准确地查询到指定系统中已审批事项的详细信息。
                
查询我发起的审批列表: 
                本工具旨在帮助您快速查询跨所有系统中您发起的审批事项。它集中展示了您的审批发起记录，无需逐一检查各个系统。
        
                【操作步骤】
                1. 根据需要指定查询类型：
                   - 若查询特定类型的审批（如费用管理），请明确指出审批类型。
                   - 若查询您发起的所有审批，请仅使用“我发起的”进行查询。
                2. 发起查询，使用格式“我发起的[审批类型, 可选]，入参：emailPrefix=xxx”。
                
                【示例查询】
                - 查询所有审批事项： “我发起的，入参：emailPrefix=xxx”。此时，应调用query_initiated_approvals函数，并提供以下参数：emailPrefix=xxx。
                - 查询特定类型审批事项（如费用管理）： “我发起的费用管理，入参：emailPrefix=xxx”。此时，应调用query_initiated_approvals函数，并提供以下参数：
                emailPrefix=xxx,  system_name='费用管理'。
                
                【注意事项】
                - 请根据您的查询需求选择适当的格式。系统将根据提供的信息执行相应的查询操作。
                - 当输入'我发起的资管审批流'时对应的system_name为'资管审批'。
                - 当输入'我发起的芝麻审批'时对应的system_name为'芝麻开门'。

                【输入参数详解】
                - system_name: 审批类型，可选。
                - emailPrefix: 邮箱前缀。

                依此指导，您将能轻松获取您发起的审批任务清单，助您高效管理和跟踪审批进度。
            
查询我发起的审批详情: 
                 使用"查询我发起的审批详情"工具时，简单提及审批类型和人名即可触发查询。为了更精确的操作，建议在询问时使用“查询”这样的动词。
        
                【使用示例1】
                假设您需要查询我发起的报销审批详情时，您应该这样提问：
                '我发起的报销审批详情，入参：emailPrefix=xxx'。
                在此示例中，您将调用query_initiated_approvals_details函数，并按以下格式提供输入参数：
                todo_type_name=报销审批&emailPrefix=xxx。
                
                【参数说明】
                - todo_type_name: 待查询的我发起的审批类型名称（例如：报销审批、出差审批等）。
                - emailPrefix: 邮箱前缀。
                
                【注意事项】
                - 请确保入参格式正确，包括我发起的审批类型名称和所需的身份验证信息。
                - 本工具支持灵活的查询格式，但建议遵循上述示例以确保准确性。
                - 当输入'我发起的资管审批流详情'时对应的system_name为'资管审批'。
                - 当输入'我发起的芝麻审批详情'时对应的system_name为'芝麻开门'。

                通过遵循上述指南，您可以有效且准确地查询到指定系统中我发起的审批事项的详细信息。
                
单个审批处理: 
                **单个审批**工具允许您快速审批或拒绝单个提交的申请。使用时，请明确指出您的动作意图，如“同意”或“拒绝”，并提供相应的单号和其他必要信息。

                **如何使用**：
                - **同意示例**：说出“同意5544667的报销申请”，并提供"emailPrefix"。调用"approve_one"函数时，使用参数：emailPrefix=xxx&action=1&orderId=5544667。
                - **拒绝示例**：说出“拒绝5544667的报销申请”，并提供"emailPrefix"。调用"approve_one"函数时，使用参数：emailPrefix=xxx&action=2&orderId=5544667。
        
                **参数说明**：
                - "emailPrefix"：用户的邮箱前缀。
                - "action"：操作类型，1代表同意，2代表拒绝。
                - "orderId"：申请的单号。
                
批量审批处理: 
                **批量审批**工具使您能够快速审批或拒绝由特定员工提交的多个申请。在提出请求时，请使用明确的动作词，如“同意”或“拒绝”，并指定待审批的系统或者员工名。
        
                **如何使用**：
                - **同意示例**：说出“同意outman的所有报销申请”，并提供"emailPrefix"。调用"approve_more"函数时，使用参数：system_name=报销审批&emailPrefix=xxx&userName=陈涛&action=1。
                - **拒绝示例**：说出“拒绝outman的所有报销申请”，并提供"emailPrefix"。调用"approve_more"函数时，使用参数：system_name=报销审批&emailPrefix=xxx&userName=陈涛&action=2。
                
                - **同意示例**：说出“同意所有报销申请”，并提供"emailPrefix"。调用"approve_more"函数时，使用参数：system_name=报销审批&emailPrefix=xxx&action=1。
                - **拒绝示例**：说出“拒绝所有报销申请”，并提供"emailPrefix"。调用"approve_more"函数时，使用参数：system_name=报销审批&emailPrefix=xxx&action=2。

                **参数说明**：
                - "system_name"：系统名称，例如'报销审批'。
                - "emailPrefix"：用户的邮箱前缀。
                - "userName"：员工的姓名【可选择传入，非必须参数】。
                - "action"：操作类型，1代表同意，2代表拒绝。
                
                **注意**
                当前此功能只支持费用报销批量审批。
                
切换助理风格: 
                当你想切换助理风格时很有用。
                例如示例1，若想切换助理风格时，可以这样提问：
                '切换助理风格，入参：emailPrefix=xxx'，
                此时，调用change_system_style函数，输入参数格式为：emailPrefix=xxx。
                
                例如示例2，若想切换助理风格时，可以这样提问：
                '魔法，入参：emailPrefix=xxx'，
                此时，调用change_system_style函数，输入参数格式为：emailPrefix=xxx。
                

Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of [分贝通快捷入口通用工具, 获取助手形象, 全系统待办综合查询, 特定系统待办查询, 全系统已审批综合查询, 特定系统已审批查询, 查询待办详情, 查询已审批详情, 查询我发起的审批列表, 查询我发起的审批详情, 单个审批处理, 批量审批处理, 切换助理风格]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question

Begin!

Question: 报销审批待审批详情，入参：xiaofengma2
Thought:"`,
						},
					},
				})
				defer wg.Done()
				if err != nil {
					t.Error(err)
					return
				}
				if len(ccr.Choices) == 0 {
					t.Error("ccr.Choices is null")
					return
				}
			}()

			concurrent--
			if concurrent == 0 {
				break
			}
		}
		wg.Wait()
		loopNum--
		if loopNum == 0 {
			break
		}
		t.Log("并发：", concurrent, "耗时: ", time.Since(_now))
	}
	t.Log("totalNum", totalNum)
	t.Log("总耗时", time.Since(now))
}

func TestChatCompletionsNoStream(t *testing.T) {
	//config := openai.DefaultConfig("sk-001")
	//config.BaseURL = "http://localhost:8080/v1"
	//client := openai.NewClientWithConfig(config)

	//pwd, _ := os.Getwd()
	//fmt.Println(pwd)
	//return
	//data, err := os.ReadFile("./question.json")
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	var qa QA
	//_ = json.Unmarshal(data, &qa)

	// 打开文件
	file, err := os.Open("~/Downloads/question.csv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	var n = 0
	// 解析 CSV 文件
	csvReader := csv.NewReader(file)
	for {
		// 读取一行
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			return
		}
		if record[2] == "mp3_content" {
			continue
		}
		// 输出数据
		qa.Questions = append(qa.Questions, record[2])
		n++
		if n > 50 {
			break
		}
	}

	type resQa struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}

	var response []resQa

	client, model := getOpenAIClient()
	//model = "vicuna-13b-16k"
	for _, v := range qa.Questions {
		req := openai.ChatCompletionRequest{
			Model:     model,
			MaxTokens: 8192,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "sys",
					Content: "我希望你充当语言检测器。我会用任何语言输入一个句子，你会回答我，我写的句子里是否有进入语音邮箱。不要写任何解释或其他文字，只需回复True或False即可。",
				},
				{
					Role:    "user",
					Content: v,
				},
			},
			Stream: false,
		}
		ctx := context.Background()
		res, err := client.CreateChatCompletion(ctx, req)
		if err != nil {
			t.Error(err)
			continue
		}
		if len(res.Choices) < 1 {
			t.Log(res)
			continue
		}
		response = append(response, resQa{
			Question: v,
			Answer:   res.Choices[0].Message.Content,
		})
	}

	b, _ := json.Marshal(response)
	_ = os.WriteFile("./answer3.json", b, 0644)
}

func TestChatCompletions(t *testing.T) {
	client, model := getOpenAIClient()
	req := openai.ChatCompletionRequest{
		Model:     model,
		MaxTokens: 2048,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello, I'm a human. who are you?",
			},
		},
		Stream: true,
	}
	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer stream.Close()
	for {
		resp, err := stream.Recv()
		if err != nil {
			t.Log(err)
			break
		}
		if len(resp.Choices) > 0 {
			t.Log(resp.Choices[0].Delta.Content)
		}
	}
	t.Log("done")
}

func TestAudioTranslations(t *testing.T) {
	client, model := getLocalOpenAIClient()
	//client, model := getOpenAIClient()
	req := openai.AudioRequest{
		Model:    model,
		FilePath: "~/Downloads/111.mp3",
		Format:   "vtt",
	}
	ctx := context.Background()
	res, err := client.CreateTranslation(ctx, req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Text)
}

func TestAudioTranscriptions(t *testing.T) {
	client, model := getLocalOpenAIClient()
	//client, model := getOpenAIClient()
	req := openai.AudioRequest{
		Model:    model,
		FilePath: "/Users/leng/Downloads/111.mp3",
		Language: "zh",
	}
	ctx := context.Background()
	res, err := client.CreateTranscription(ctx, req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Text)
}

func TestTexToAudio(t *testing.T) {
	//u, _ := url.Parse("http://localhost:8080/v1/text/audio")
	u, _ := url.Parse("http://aigc-server:8080/v1/text/audio")
	var opts []kithttp.ClientOption
	opts = append(opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
		return ctx
	}),
		kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
			dump, _ := httputil.DumpResponse(response, true)
			fmt.Println(string(dump))
			return ctx
		}),
	)
	ep := kithttp.NewClient(http.MethodPost, u, func(ctx context.Context, r *http.Request, request interface{}) error {
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Authorization", "Bearer sk-0011")
		var b bytes.Buffer
		r.Body = io.NopCloser(&b)
		return json.NewEncoder(&b).Encode(request)
	}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
		if response2.StatusCode != http.StatusOK {
			return nil, errors.New("http status code not 200")
		}
		b, err := io.ReadAll(response2.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}, opts...).Endpoint()
	type textToAudioParams struct {
		Text string `json:"text"`
	}

	var params textToAudioParams
	params.Text = "你好，我是一个人类。你叫什么名字？"
	ctx := context.Background()
	res, err := ep(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res.([]byte)))
}

func TestDifyChat(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/v1/dify/chat")
	var opts []kithttp.ClientOption
	opts = append(opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
		return ctx
	}),
		kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
			dump, _ := httputil.DumpResponse(response, true)
			fmt.Println(string(dump))
			return ctx
		}),
	)
	ep := kithttp.NewClient(http.MethodPost, u, func(ctx context.Context, r *http.Request, request interface{}) error {
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Authorization", "Bearer sk-0011")
		var b bytes.Buffer
		r.Body = io.NopCloser(&b)
		return json.NewEncoder(&b).Encode(request)
	}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
		if response2.StatusCode != http.StatusOK {
			return nil, errors.New("http status code not 200")
		}
		b, err := io.ReadAll(response2.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}, opts...).Endpoint()
	type difyChatParams struct {
		Query          string `json:"query" validate:"required"`
		ConversationId string `json:"conversation_id"`
		User           string `json:"user"`
	}

	var params difyChatParams
	params.Query = "你好你好"
	//params.User = "bx-aics"
	ctx := context.Background()
	res, err := ep(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res.([]byte)))
}
