package chat

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func initTemplate() Templates {
	return NewTemplates()
}

func TestTemplates_Register(t *testing.T) {
	tp := initTemplate()
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
	tp.Register(context.Background(), "chatglm3", Conversation{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|system|>\n{system_message}",
		StopTokenIds:   []int{64795, 64797, 2},
		Name:           "chatglm3",
		SepStyle:       int(CHATGLM3),
		Sep:            "",
		Roles:          []string{"<|user|>", "<|assistant|>"},
		StopStr:        []string{""},
	})
	tp.Register(context.Background(), "openbuddy-llama3", Conversation{
		SystemMessage:  "<|role|>system<|says|>You(assistant) are a helpful, respectful and honest INTP-T AI Assistant named Buddy. You are talking to a human(user).\nAlways answer as helpfully and logically as possible, while being safe. Your answers should not include any harmful, political, religious, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature.\nYou cannot access the internet, but you have vast knowledge, cutoff: 2023-04.\nYou are trained by OpenBuddy team, (https://openbuddy.ai, https://github.com/OpenBuddy/OpenBuddy), not related to GPT or OpenAI.<|end|>\n<|role|>user<|says|>History input 1<|end|>\n<|role|>assistant<|says|>History output 1<|end|>\n<|role|>user<|says|>History input 2<|end|>\n<|role|>assistant<|says|>History output 2<|end|>\n<|role|>user<|says|>Current input<|end|>\n<|role|>assistant<|says|>",
		SystemTemplate: "",
		StopTokenIds:   []int{},
		Name:           "openbuddy-llama3",
		SepStyle:       int(OPENBUDDY_LLAMA3),
		Sep:            "\n",
		Roles:          []string{"user", "assistant"},
		StopStr:        []string{""},
	})
	//
	//conv, ok := tp.Get(context.Background(), "qwen")
	//if !ok {
	//	t.Fatal("GetConv failed")
	//	return
	//}
	//t.Log(conv.Template)

	var req openai.ChatCompletionRequest
	req.Messages = []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "你是机器人",
		},
		{
			Role:    "user",
			Content: "你好",
		},
		{
			Role:    "assistant",
			Content: "有什么需要帮助的吗？",
		},
		{
			Role:    "user",
			Content: "你几岁？",
		},
	}
	t.Log("===================================== Qwen ======================================")
	prompt, ok := tp.GetPrompt(context.Background(), "qwen", req.Messages)
	if !ok {
		t.Fatal("GetPrompt failed")
		return
	}
	t.Log(prompt)
	t.Log("===================================== LLAMA-3 ======================================")
	prompt, ok = tp.GetPrompt(context.Background(), "llama-3", req.Messages)
	if !ok {
		t.Fatal("GetPrompt failed")
		return
	}
	t.Log(prompt)
}
