package tokenizers

import (
	"bytes"
	"encoding/json"
)

type (
	// Message 用于解析和验证每一行的JSON对象
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// MessagesWrapper 包含多个Message的结构体
	MessagesWrapper struct {
		Messages []Message `json:"messages"`
	}

	// DataAnnotationSegment 数据标注的片段
	DataAnnotationSegment struct {
		Instruction string `json:"instruction,omitempty"`
		Input       string `json:"input,omitempty"`
		Output      string `json:"output,omitempty"`
		Intent      string `json:"intent,omitempty"`
		Document    string `json:"document,omitempty"`
		Question    string `json:"question,omitempty"`
	}
)

// GetFirstLineSystemPrompt 获取文件的第一行系统提示
func GetFirstLineSystemPrompt(fileBody []byte) (content string) {
	var fastLine []byte
	lines := bytes.Split(fileBody, []byte("\n"))
	if len(lines) > 0 {
		fastLine = bytes.TrimSpace(lines[0])
	}
	if len(fastLine) == 0 {
		return
	}
	var msg MessagesWrapper
	if json.Unmarshal(fastLine, &msg) == nil {
		for _, v := range msg.Messages {
			if v.Role == "system" {
				content = v.Content
				return
			}
		}
		return
	}

	var dataAnnotationSegment DataAnnotationSegment
	if json.Unmarshal(fastLine, &dataAnnotationSegment) == nil {
		content = dataAnnotationSegment.Instruction
		return
	}

	return

}

// ConvertToMessages 将文件内容转换为消息
func ConvertToMessages(fileBody []byte) (messages []MessagesWrapper, err error) {
	lines := bytes.Split(fileBody, []byte("\n"))
	if len(lines) == 0 {
		return
	}
	var msg MessagesWrapper
	var dataAnnotationSegment DataAnnotationSegment
	fastLine := bytes.TrimSpace(lines[0])
	if json.Unmarshal(fastLine, &msg) == nil {
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			var mw MessagesWrapper
			if json.Unmarshal(line, &mw) == nil {
				messages = append(messages, mw)
			}
		}
	} else if json.Unmarshal(fastLine, &dataAnnotationSegment) == nil {
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			var das DataAnnotationSegment
			if json.Unmarshal(line, &das) == nil {
				messages = append(messages, MessagesWrapper{Messages: []Message{
					{
						Content: das.Instruction,
						Role:    "system",
					},
					{
						Content: das.Input,
						Role:    "user",
					},
					{
						Content: das.Output,
						Role:    "assistant",
					},
				}})
			}
		}
	}
	return
}

// ConvertToDatasets 根据模型名称将消息转换为训练的数据集
func ConvertToDatasets(messages []MessagesWrapper, modelName string) (datasets []string) {
	return
}
