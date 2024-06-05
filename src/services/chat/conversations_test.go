package chat

import (
	"context"
	"testing"
)

func TestConversation_GetPrompt(t *testing.T) {
	tpl := NewTemplates()
	tpl.Register(context.Background(), "phi-3", Conversation{
		StopTokenIds: []int{},
		Name:         "phi-3",
		SepStyle:     int(CHATML),
		Sep:          "<|end|>",
		Roles:        []string{"<|user|>", "<|assistant|>"},
		StopStr:      []string{"<|endoftext|>", "<|user|>"},
	})

	conv, ok := tpl.Get(context.Background(), "phi-3")
	if !ok {
		t.Error("failed to get conv template")
		return
	}

	conv.SetSystemMessage("you are a helpful assistant")
	conv.AppendMessage(conv.Roles[0], "hello")
	conv.AppendMessage(conv.Roles[1], "hi")
	conv.AppendMessage(conv.Roles[0], "how are you")
	conv.AppendMessage(conv.Roles[1], "")

	t.Log(conv.GetPrompt())
}

func TestConversation_GetPrompt2(t *testing.T) {
	tpl := Register(NewTemplates())

	t.Log("Llama-3-8B-Instruct")
	conv, ok := tpl.GetByModelName(context.Background(), "Llama-3-8B-Instruct")
	if !ok {
		t.Error("failed to get conv template")
		return
	}
	//conv.SetSystemMessage("you are a helpful assistant")
	conv.AppendMessage(conv.Roles[0], "hello")
	conv.AppendMessage(conv.Roles[1], "hi")
	conv.AppendMessage(conv.Roles[0], "how are you")
	conv.AppendMessage(conv.Roles[1], "")

	t.Log(conv.GetPrompt())
	t.Log(conv.Name)

	t.Log("\n\nLlama3-8B-Chinese-Chat")
	conv, ok = tpl.GetByModelName(context.Background(), "Llama3-8B-Chinese-Chat")
	if !ok {
		t.Error("failed to get conv template")
		return
	}
	conv.SetSystemMessage("you are a helpful assistant")
	conv.AppendMessage(conv.Roles[0], "hello")
	conv.AppendMessage(conv.Roles[1], "hi")
	conv.AppendMessage(conv.Roles[0], "how are you")
	conv.AppendMessage(conv.Roles[1], "")

	t.Log(conv.GetPrompt())
	t.Log(conv.Name)

	t.Log("\n\nglm-4-9b-chat")
	conv, ok = tpl.GetByModelName(context.Background(), "glm-4-9b-chat")
	if !ok {
		t.Error("failed to get conv template")
		return
	}
	conv.SetSystemMessage("you are a helpful assistant")
	conv.AppendMessage(conv.Roles[0], "hello")
	conv.AppendMessage(conv.Roles[1], "hi")
	conv.AppendMessage(conv.Roles[0], "how are you")
	conv.AppendMessage(conv.Roles[1], "")
	t.Log(conv.GetPrompt())
	t.Log(conv.Name)
}
