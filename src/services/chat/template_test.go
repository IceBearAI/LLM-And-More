package chat

import (
	"context"
	"testing"
)

func initTemplate() Templates {
	return NewTemplates()
}

func TestTemplates_Register(t *testing.T) {
	tp := initTemplate()
	tp.Register(context.Background(), "llama-3", Conv{
		Name:     "llama-3",
		Template: "{{ if .System }}<|start_header_id|>system<|end_header_id|>\n\n{{ .System }}<|eot_id|>{{ end }}{{ if .Prompt }}<|start_header_id|>user<|end_header_id|>\n\n{{ .Prompt }}<|eot_id|>{{ end }}<|start_header_id|>assistant<|end_header_id|>\n\n{{ .Response }}<|eot_id|>",
	})
	tp.Register(context.Background(), "qwen", Conv{
		Name:     "qwen",
		Template: "{{ if .System }}<|im_start|>system\n{{ .System }}<|im_end|>{{ end }}<|im_start|>user\n{{ .Prompt }}<|im_end|>\n<|im_start|>assistant",
	})

	conv, ok := tp.GetConv(context.Background(), "qwen")
	if !ok {
		t.Fatal("GetConv failed")
		return
	}
	t.Log(conv.Template)
}
