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
	"testing"
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
	config := openai.DefaultConfig("sk-0011")
	config.BaseURL = "http://localhost:8080/v1"
	client := openai.NewClientWithConfig(config)
	return client, openai.GPT3Dot5Turbo
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
