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
	"net/url"
	"os"
	"strings"
	"testing"
)

var systemContent = `以下是一个FAQ客服系统可能对应的用户意图口。

"询问基金是否保本;询问工作时间"

用户输入为: %s

请将用户输入归类为上述用户意图中的一种。不用说明理由和解释，直接输出最终答案。

如果你无法识别出意图，则返回{转人工}。`

type Intent struct {
	Intent    string   `json:"intent"`
	Answer    string   `json:"answer"`
	Questions []string `json:"questions"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type input struct {
	Messages []message `json:"messages"`
}

func TestConvert_Message2Intent(t *testing.T) {
	pwd, _ := os.Getwd()
	fileByte, err := os.ReadFile(fmt.Sprintf("%s/ft-20240118145626.jsonl", pwd))
	if err != nil {
		t.Error(err)
		return
	}
	var intentList []Intent
	var answers []string
	lines := strings.Split(string(fileByte), "\n")
	for _, line := range lines {
		var msg input
		_ = json.Unmarshal([]byte(line), &msg)
		if len(msg.Messages) > 2 {
			t.Log("len(msg.Messages) > 2", msg.Messages)
			continue
		}
		var answer, question string
		for _, m := range msg.Messages {
			if m.Role == "user" {
				question = m.Content
				continue
			}
			if m.Role == "assistant" {
				answer = m.Content
			}
		}
		if inArray(answers, answer) {
			intentList = appendIntent(intentList, question, answer)
		} else {
			answers = append(answers, answer)
			intentList = append(intentList, Intent{
				Intent:    "",
				Answer:    answer,
				Questions: []string{question},
			})
		}
	}
	b, _ := json.Marshal(intentList)
	_ = os.WriteFile(fmt.Sprintf("%s/intent-0124.json", pwd), b, 0666)
}

func appendIntent(intentList []Intent, question, answer string) []Intent {
	for k, v := range intentList {
		if strings.EqualFold(v.Answer, answer) {
			intentList[k].Questions = append(v.Questions, question)
			return intentList
		}
	}
	return intentList
}

func inArray(arr []string, str string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

func TestConvert_Intent(t *testing.T) {
	pwd, _ := os.Getwd()
	fileByte, err := os.ReadFile(fmt.Sprintf("%s/intent.json", pwd))
	if err != nil {
		t.Error(err)
		return
	}
	var intents []Intent
	err = json.Unmarshal(fileByte, &intents)
	if err != nil {
		t.Error(err)
		return
	}

	csvFile, err := os.Create(fmt.Sprintf("%s/intent-gpt-4-1106-preview.csv", pwd))
	if err != nil {
		t.Error(err)
		return
	}

	defer csvFile.Close()

	// 创建一个csv.Writer，并将其绑定到文件
	writer := csv.NewWriter(csvFile)
	defer writer.Flush() // 确保在函数结束时将所有缓冲的数据写入文件

	// 写入CSV标题行
	headers := []string{"意图", "模型意图", "是否匹配", "问题", "答案"}
	if err := writer.Write(headers); err != nil {
		t.Error(err)
		return
	}

	for _, v := range intents {
		for _, question := range v.Questions {
			res, err := httpPost(question)
			if err != nil {
				t.Error(err)
				continue
			}
			if !strings.EqualFold(res.Choices[0].Message.Content, v.Intent) {
				_ = writer.Write([]string{v.Intent, res.Choices[0].Message.Content, "否", question, v.Answer})
				t.Log(res.Choices[0].Message.Content)
			} else {
				_ = writer.Write([]string{v.Intent, res.Choices[0].Message.Content, "是", question, v.Answer})
			}
		}
	}
}

func httpPost(question string) (result openai.ChatCompletionResponse, err error) {
	tgt, _ := url.Parse("http://chat-api:8080/v1/chat/completions")
	ep := kithttp.NewClient(http.MethodPost, tgt, func(ctx context.Context, request *http.Request, i interface{}) error {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer sk-001")
		//b, _ := json.Marshal(request)
		//request.Body = io.NopCloser(bytes.NewReader(b))

		var b bytes.Buffer
		request.Body = io.NopCloser(&b)
		return json.NewEncoder(&b).Encode(i)
	}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
		if response2.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http status code is %d", response2.StatusCode)
		}
		return io.ReadAll(response2.Body)
	}).Endpoint()
	var messages []openai.ChatCompletionMessage
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "sys",
		Content: fmt.Sprintf(systemContent, question),
	}, openai.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})
	res, err := ep(context.Background(), openai.ChatCompletionRequest{
		Model:       "qwen-72b-chat",
		Messages:    messages,
		MaxTokens:   128,
		Temperature: 0,
		TopP:        0,
		N:           0,
		Stream:      false,
	})
	if err != nil {
		err = errors.Wrap(err, "ep(context.Background(), question)")
		return
	}
	if err = json.Unmarshal(res.([]byte), &result); err != nil {
		err = errors.Wrap(err, "json.Unmarshal(res.([]byte), &result)")
		return
	}
	return
}
