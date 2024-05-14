package chat

import (
	"bytes"
	"context"
	"github.com/sashabaranov/go-openai"
	"strings"
	"text/template"
)

// TemplateMessage 模板消息
type TemplateMessage struct {
	// System 系统
	System string
	// Prompt 提示 user.content
	Prompt string
	// Response 响应 assistant.content
	Response string
}

type Templates interface {
	// Register 注册模板
	Register(ctx context.Context, name string, conv Conversation)
	// Get 获取模板
	Get(ctx context.Context, name string) (Conversation, bool)
	// GetByModelName 根据模型名称获取模版
	GetByModelName(ctx context.Context, modelName string) (Conversation, bool)
	// GetPrompt 获取提示
	GetPrompt(ctx context.Context, name string, messages []openai.ChatCompletionMessage) (string, bool)
}

type templates struct {
	conv map[string]Conversation
}

func (t *templates) GetPrompt(ctx context.Context, name string, messages []openai.ChatCompletionMessage) (string, bool) {
	name = strings.ToLower(name)
	conv, ok := t.conv[name]
	if !ok {
		return "", false
	}
	var res string
	for _, v := range messages {
		tpl, err := encodeTemplate(name, conv.Template, v)
		if err != nil {
			return "", false
		}
		res += tpl
	}
	res += conv.Roles[0]
	return res, true
}

// encodeTemplate 模版渲染
func encodeTemplate(name string, tpl string, data interface{}) (re string, err error) {
	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		panic(err)
	}
	re = buffer.String()
	return
}

func (t *templates) GetByModelName(ctx context.Context, modelName string) (Conversation, bool) {
	modelName = strings.ToLower(modelName)
	var models []Conversation
	for _, v := range t.conv {
		if strings.Contains(modelName, v.Name) {
			models = append(models, v)
		}
	}
	if len(models) > 0 {
		// 优先获取名称长的
		for i := 0; i < len(models); i++ {
			for j := i + 1; j < len(models); j++ {
				if len(models[i].Name) < len(models[j].Name) {
					models[i], models[j] = models[j], models[i]
				}
			}
		}
		return models[0], true
	}
	return Conversation{}, false
}

func (t *templates) Register(ctx context.Context, name string, conv Conversation) {
	name = strings.ToLower(name)
	if _, ok := t.conv[name]; ok {
		return
	}
	t.conv[name] = conv
}

func (t *templates) Get(ctx context.Context, name string) (Conversation, bool) {
	conv, ok := t.conv[name]
	return conv, ok
}

func NewTemplates() Templates {
	return &templates{conv: map[string]Conversation{}}
}

func ConvertMessages(messages []openai.ChatCompletionMessage) []TemplateMessage {
	var tms []TemplateMessage
	var currentTm TemplateMessage
	var expectRole string

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			// 当遇到 system 消息时，直接开始一个新的 TemplateMessage
			if currentTm.System != "" || currentTm.Prompt != "" || currentTm.Response != "" {
				// 如果当前 TemplateMessage 有内容，则先保存它
				tms = append(tms, currentTm)
			}
			// 重置 currentTm 并设置 System 字段
			currentTm = TemplateMessage{System: msg.Content}
			expectRole = "user" // 预期下一个角色是 user
		case "user":
			if expectRole == "user" {
				// 如果预期的是 user 消息，则设置 Prompt 并更新预期角色
				currentTm.Prompt = msg.Content
				expectRole = "assistant" // 下一个预期角色更新为 assistant
			} else if expectRole == "assistant" {
				// 如果上一个消息是 user 但没有对应的 assistant 回复，保存当前 TemplateMessage 并开始一个新的
				tms = append(tms, currentTm)
				currentTm = TemplateMessage{Prompt: msg.Content}
				expectRole = "assistant"
			}
		case "assistant":
			if expectRole == "assistant" {
				// 如果预期的是 assistant 消息，则设置 Response
				currentTm.Response = msg.Content
				// 保存当前 TemplateMessage 并准备开始一个新的
				tms = append(tms, currentTm)
				currentTm = TemplateMessage{}
				expectRole = "user" // 重置预期角色为 user，等待下一个消息
			}
		}
	}

	// 确保最后一个 TemplateMessage 也被保存
	if currentTm.System != "" || currentTm.Prompt != "" || currentTm.Response != "" {
		tms = append(tms, currentTm)
	}

	return tms
}
