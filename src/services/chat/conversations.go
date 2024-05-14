package chat

import (
	"strings"
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
		ret = systemPrompt
		if c.SystemMessage != "" {
			ret += c.SystemMessage
		} else {
			ret += ""
		}
		for i, message := range c.Messages {
			if len(message) == 2 && message[1] != "" {
				ret += "<|start_header_id|>" + c.Roles[i%2] + "<|end_header_id|>\n\n"
				ret += message[1]
			} else {
				ret += "<|start_header_id|>" + c.Roles[i%2] + "<|end_header_id|>\n\n"
			}
		}
		return ret
	case CHATML:
		if c.SystemMessage != "" {
			ret = systemPrompt + c.Sep + "\n"
		} else {
			ret = ""
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
	}

	return ret
}
