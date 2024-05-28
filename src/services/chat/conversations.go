package chat

import (
	"context"
	"strings"
)

type SeparatorStyle int

const (
	_ SeparatorStyle = iota
	ADD_COLON_SINGLE
	ADD_COLON_TWO
	ADD_COLON_SPACE_SINGLE
	NO_COLON_SINGLE
	NO_COLON_TWO
	ADD_NEW_LINE_SINGLE
	LLAMA2
	LLAMA3
	CHATGLM
	CHATML
	CHATINTERN
	DOLLY
	RWKV
	PHOENIX
	ROBIN
	FALCON_CHAT
	CHATGLM3
	DEEPSEEK_CHAT
	METAMATH
	YUAN2
	GEMMA
	CLLM
	DEFAULT
	OPENBUDDY_LLAMA3
	PHI3
)

// Conversation 对话
type Conversation struct {
	// Name 名称
	Name string `json:"name"`
	// SystemTemplate 系统模板
	SystemTemplate string `json:"system_template"`
	// SystemMessage 系统消息
	SystemMessage string `json:"system_message"`
	// Roles 角色
	Roles []string `json:"roles"`
	// Messages 消息 [role, message]
	Messages [][]string `json:"messages"`
	// Offset 偏移
	Offset int `json:"offset"`
	// SepStyle 分隔符样式
	SepStyle int `json:"sep_style"`
	// Sep 分隔符
	Sep string `json:"sep"`
	// Sep2 分隔符2
	Sep2 string `json:"sep2"`
	// StopStr 停止字符串
	StopStr []string `json:"stop_str"`
	// StopTokenIds 停止令牌ID
	StopTokenIds []int `json:"stop_token_ids"`
	// Template 模板
	Template string `json:"template,omitempty"`
}

func (c *Conversation) SetSystemMessage(systemMessage string) {
	c.SystemMessage = systemMessage
}

func (c *Conversation) AppendMessage(role string, message string) {
	c.Messages = append(c.Messages, []string{role, message})
}

func (c *Conversation) GetImages() []string {
	var images []string
	for i, message := range c.Messages[c.Offset:] {
		if i%2 == 0 {
			if len(message) == 2 {
				images = append(images, message[1])
			}
		}
	}
	return nil
}

func (c *Conversation) GetPrompt() (ret string) {
	var systemPrompt = strings.ReplaceAll(c.SystemTemplate, "{system_message}", c.SystemMessage)
	switch SeparatorStyle(c.SepStyle) {
	case ADD_COLON_SINGLE:
		ret = systemPrompt + c.Sep
		for _, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += message[0] + ": " + message[1] + c.Sep
			} else {
				ret += message[0] + ":"
			}
		}
		return ret
	case NO_COLON_SINGLE:
		ret = systemPrompt
		for _, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += message[0] + message[1] + c.Sep
			} else {
				ret += message[0]
			}
		}
		return ret
	case LLAMA2:
		seps := []string{c.Sep, c.Sep2}
		if c.SystemMessage != "" {
			ret = systemPrompt
		} else {
			ret = "[INST] "
		}
		for i, message := range c.Messages {
			tag := c.Roles[i%2]
			if len(message) == 2 && message[1] != "" {
				if i == 0 {
					ret += message[1] + " "
				} else {
					ret += tag + " " + message[1] + seps[i%2]
				}
			} else {
				ret += tag
			}
		}
		return ret
	case LLAMA3:
		ret = "<|begin_of_text|>"
		ret += systemPrompt
		for i, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += "<|start_header_id|>" + c.Roles[i%2] + "<|end_header_id|>\n\n"
				ret += message[1] + "<|eot_id|>"
			} else {
				ret += "<|start_header_id|>" + c.Roles[i%2] + "<|end_header_id|>\n\n"
			}
		}
		return ret
	case CHATML:
		if systemPrompt != "" {
			ret = systemPrompt + c.Sep + "\n"
		}
		for _, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += message[0] + "\n" + message[1] + c.Sep + "\n"
			} else {
				ret += message[0] + "\n"
			}
		}
		return ret
	case CHATGLM3:
		if c.SystemMessage != "" {
			ret = systemPrompt
		} else {
			ret = ""
		}
		for _, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += message[0] + "\n" + message[1]
			} else {
				ret += message[0]
			}
		}
		return ret
	case OPENBUDDY_LLAMA3:
		ret = systemPrompt + "\n"
		for _, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += "<|role|>" + message[0] + "<|says|>" + message[1] + "<|end|>\n"
			} else {
				ret += "<|role|>" + message[0] + "<|says|>\n"
			}
		}
		return ret
	default:
		ret = ""
	}

	return ret
}

func Register(tp Templates) Templates {
	tp.Register(context.Background(), "llama-3", Conversation{
		StopStr:        []string{"<|eot_id|>"},
		Sep:            "",
		Sep2:           "",
		StopTokenIds:   []int{128001, 128009},
		Name:           "llama-3",
		SepStyle:       int(LLAMA3),
		Roles:          []string{"user", "assistant"},
		SystemTemplate: "<|start_header_id|>system<|end_header_id|>\n\n{system_message}<|eot_id|>",
		SystemMessage:  "You are a helpful assistant.",
	})
	tp.Register(context.Background(), "qwen", Conversation{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|im_start|>system\n{system_message}",
		StopTokenIds:   []int{151643, 151644, 151645},
		Name:           "qwen",
		SepStyle:       int(CHATML),
		Sep:            "<|im_end|>",
		Roles:          []string{"<|im_start|>user", "<|im_start|>assistant"},
		StopStr:        []string{"<|endoftext|>"},
	})
	tp.Register(context.Background(), "chatglm2", Conversation{
		StopTokenIds: []int{},
		Name:         "chatglm2",
		SepStyle:     int(CHATGLM),
		Sep:          "\n\n",
		Roles:        []string{"问", "答"},
		StopStr:      []string{},
	})
	tp.Register(context.Background(), "chatglm3", Conversation{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|system|>\n{system_message}",
		StopTokenIds:   []int{64795, 64797, 2},
		Name:           "chatglm3",
		SepStyle:       int(CHATGLM3),
		Sep:            "",
		Roles:          []string{"<|user|>", "<|assistant|>"},
		StopStr:        []string{"<|observation|>", "<|user|>", "</s>"},
	})
	tp.Register(context.Background(), "openbuddy-llama3", Conversation{
		SystemMessage:  "<|role|>system<|says|>You(assistant) are a helpful, respectful and honest INTP-T AI Assistant named Buddy. You are talking to a human(user).\nAlways answer as helpfully and logically as possible, while being safe. Your answers should not include any harmful, political, religious, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature.\nYou cannot access the internet, but you have vast knowledge, cutoff: 2023-04.\nYou are trained by OpenBuddy team, (https://openbuddy.ai, https://github.com/OpenBuddy/OpenBuddy), not related to GPT or OpenAI.<|end|>\n<|role|>user<|says|>History input 1<|end|>\n<|role|>assistant<|says|>History output 1<|end|>\n<|role|>user<|says|>History input 2<|end|>\n<|role|>assistant<|says|>History output 2<|end|>\n<|role|>user<|says|>Current input<|end|>\n<|role|>assistant<|says|>",
		SystemTemplate: "",
		StopTokenIds:   []int{},
		Name:           "openbuddy-llama3",
		SepStyle:       int(OPENBUDDY_LLAMA3),
		Sep:            "\n",
		Roles:          []string{"user", "assistant"},
		StopStr:        nil,
	})
	tp.Register(context.Background(), "baichuan2", Conversation{
		StopTokenIds: []int{},
		Name:         "baichuan2",
		SepStyle:     int(NO_COLON_SINGLE),
		Sep:          "",
		Roles:        []string{"<reserved_106>", "<reserved_107>"},
		StopStr:      nil,
	})
	tp.Register(context.Background(), "phi-3", Conversation{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|system|>\n{system_message}",
		StopTokenIds:   []int{32000, 32007},
		Name:           "phi-3",
		SepStyle:       int(CHATML),
		Sep:            "<|end|>",
		Roles:          []string{"<|user|>", "<|assistant|>"},
		StopStr:        []string{"<|endoftext|>"},
	})
	tp.Register(context.Background(), "yi-", Conversation{
		StopTokenIds: []int{2, 6, 7, 8}, // "<|endoftext|>", "<|im_start|>", "<|im_end|>", "<|im_sep|>"
		Name:         "yi-",
		SepStyle:     int(CHATML),
		Sep:          "<|im_end|>",
		Roles:        []string{"<|im_start|>user", "<|im_start|>assistant"},
		StopStr:      []string{"<|endoftext|>", "<|im_start|>", "<|im_end|>"},
	})
	tp.Register(context.Background(), "llama3", Conversation{
		StopStr:        []string{"<|start_header_id|>", "<|end_header_id|>", "<|eot_id|>"},
		Sep:            "",
		Sep2:           "",
		StopTokenIds:   []int{128001, 128009},
		Name:           "llama3",
		SepStyle:       int(LLAMA3),
		Roles:          []string{"user", "assistant"},
		SystemTemplate: "<|start_header_id|>system<|end_header_id|>\n\n{system_message}<|eot_id|>",
		SystemMessage:  "You are a helpful assistant.",
	})
	return tp
}
