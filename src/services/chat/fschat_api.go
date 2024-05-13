package chat

import (
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"strings"
	"time"
)

type fsChatApiClient struct {
	options  *CreationOptions
	template Templates
}

func (s *fsChatApiClient) Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	streamResp, err := s.ChatCompletionStream(ctx, openai.ChatCompletionRequest{})
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}
	var content CompletionStreamResponse
	for {
		rs, ok := <-streamResp
		if !ok {
			break
		}
		content = rs
	}
	res = openai.CompletionResponse{
		Usage: content.Usage,
		Choices: []openai.CompletionChoice{
			{
				FinishReason: string(content.Choices[0].FinishReason),
				Text:         content.Choices[0].Delta.Content,
			},
		},
	}
	return
}

func (s *fsChatApiClient) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error) {
	streamResp, err := s.ChatCompletionStream(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}
	var content CompletionStreamResponse
	for {
		rs, ok := <-streamResp
		if !ok {
			break
		}
		if len(rs.Choices) > 0 && rs.Choices[0].Delta.Content != "" {
			content = rs
		}
	}
	res = CompletionResponse{
		Usage: openai.Usage{
			PromptTokens:     content.Usage.PromptTokens,
			CompletionTokens: content.Usage.CompletionTokens,
			TotalTokens:      content.Usage.TotalTokens,
		},
		ChatCompletionResponse: openai.ChatCompletionResponse{
			ID:      fmt.Sprintf("cmpl-%s", shortuuid.New()),
			Object:  "chat.completion",
			Created: time.Now().UnixMilli(),
			Model:   req.Model,
			Choices: []openai.ChatCompletionChoice{
				{
					FinishReason: content.Choices[0].FinishReason,
					Message: openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: content.Choices[0].Delta.Content,
					},
				},
			},
		},
	}
	return
}

func (s *fsChatApiClient) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	models, err := s.options.workerSvc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return
	}
	var exists bool
	for _, model := range models {
		if model.ID == req.Model {
			exists = true
			break
		}
	}
	if !exists {
		err = errors.New("model not found")
		return
	}
	workerAddress, err := s.options.workerSvc.GetWorkerAddress(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		return
	}
	var maxTokens int
	//prompts := s.processInput(req.Model, req.Messages)
	//for _, prompt := range prompts {
	//	maxTokens, err = s.options.workerSvc.WorkerCheckLength(ctx, workerAddress, req.Model, req.MaxTokens, prompt)
	//	if err != nil {
	//		err = errors.WithMessage(err, "failed to check length")
	//		_ = level.Warn(logger).Log("msg", "failed to check length", "err", err)
	//		return res, err
	//	}
	//}
	//_ = level.Info(logger).Log("msg", "max tokens", "maxTokens", maxTokens)
	if maxTokens != 0 && maxTokens < req.MaxTokens {
		req.MaxTokens = maxTokens
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2048
	}

	dot := make(chan CompletionStreamResponse)
	genParams, err := s.genParams(ctx, req, workerAddress)
	if err != nil {
		err = errors.WithMessage(err, "failed to get gen params")
		return
	}
	streamResp, err := s.options.workerSvc.WorkerGenerateStream(ctx, workerAddress, genParams)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}

	go func() {
		now := time.Now().UnixMilli()
		defer close(dot)
		streamId := fmt.Sprintf("cmpl-%s", shortuuid.New())
		for {
			content, ok := <-streamResp
			if !ok {
				break
			}
			if content.ErrorCode != 0 {
				err = errors.New(content.Text)
				return
			}
			text := content.Text
			var previousText string
			// 替换所有的Unicode替代字符\ufffd为空字符串
			decodedUnicode := strings.Replace(text, "\ufffd", "", -1)

			// 获取新的字符串，它是当前文本去掉与之前文本相同部分后的结果
			deltaText := decodedUnicode
			if len(previousText) < len(decodedUnicode) {
				deltaText = decodedUnicode[len(previousText):]
			}

			// 更新previous_text变量为当前文本，但只在当前文本的长度大于previous_text的长度时
			if len(decodedUnicode) > len(previousText) {
				previousText = decodedUnicode
			}

			dot <- CompletionStreamResponse{
				Usage: struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				}{PromptTokens: content.Usage.PromptTokens, CompletionTokens: content.Usage.CompletionTokens, TotalTokens: content.Usage.TotalTokens},
				ChatCompletionStreamResponse: openai.ChatCompletionStreamResponse{
					ID:      streamId,
					Object:  "chat.completion.chunk",
					Created: now,
					Model:   req.Model,
					Choices: []openai.ChatCompletionStreamChoice{
						{
							FinishReason: openai.FinishReason(content.FinishReason),
							Delta: openai.ChatCompletionStreamChoiceDelta{
								Content: deltaText,
								Role:    "assistant",
							},
						},
					},
				},
			}
		}
	}()

	return dot, nil
}

func (s *fsChatApiClient) Models(ctx context.Context) (res []openai.Model, err error) {
	models, err := s.options.workerSvc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return nil, err
	}
	for _, model := range models {
		res = append(res, openai.Model{
			ID:   model.ID,
			Root: model.Root,
		})
	}
	return
}

func (s *fsChatApiClient) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func NewFsChatApi(opts ...CreationOption) Service {
	options := &CreationOptions{
		endpoints: []Endpoint{
			{
				Host:     "http://localhost:8000/v1",
				Token:    "",
				Platform: "localai",
			},
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	tp := NewTemplates()
	return &fsChatApiClient{
		options:  options,
		template: register(tp),
	}
}

func register(tp Templates) Templates {
	tp.Register(context.Background(), "llama-3", Conv{
		StopStr:        "<|eot_id|>",
		Sep:            "",
		Sep2:           "",
		StopTokenIds:   []int{128001, 128009},
		Name:           "llama-3",
		SepStyle:       int(LLAMA3),
		Roles:          []string{"user", "assistant"},
		SystemTemplate: "<|start_header_id|>system<|end_header_id|>\n\n{system_message}<|eot_id|>",
		SystemMessage:  "You are a helpful assistant.",
	})
	tp.Register(context.Background(), "qwen", Conv{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|im_start|>system\n{system_message}",
		StopTokenIds:   []int{151643, 151644, 151645},
		Name:           "qwen",
		SepStyle:       int(CHATML),
		Sep:            "<|im_end|>",
		Roles:          []string{"<|im_start|>user", "<|im_start|>assistant"},
		StopStr:        "<|endoftext|>",
	})
	tp.Register(context.Background(), "chatglm3", Conv{
		SystemMessage:  "You are a helpful assistant.",
		SystemTemplate: "<|system|>\n{system_message}",
		StopTokenIds:   []int{64795, 64797, 2},
		Name:           "chatglm3",
		SepStyle:       int(CHATGLM3),
		Sep:            "",
		Roles:          []string{"<|user|>", "<|assistant|>"},
		StopStr:        "",
	})
	tp.Register(context.Background(), "openbuddy-llama3", Conv{
		SystemMessage:  "<|role|>system<|says|>You(assistant) are a helpful, respectful and honest INTP-T AI Assistant named Buddy. You are talking to a human(user).\nAlways answer as helpfully and logically as possible, while being safe. Your answers should not include any harmful, political, religious, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature.\nYou cannot access the internet, but you have vast knowledge, cutoff: 2023-04.\nYou are trained by OpenBuddy team, (https://openbuddy.ai, https://github.com/OpenBuddy/OpenBuddy), not related to GPT or OpenAI.<|end|>\n<|role|>user<|says|>History input 1<|end|>\n<|role|>assistant<|says|>History output 1<|end|>\n<|role|>user<|says|>History input 2<|end|>\n<|role|>assistant<|says|>History output 2<|end|>\n<|role|>user<|says|>Current input<|end|>\n<|role|>assistant<|says|>",
		SystemTemplate: "",
		StopTokenIds:   []int{},
		Name:           "openbuddy-llama3",
		SepStyle:       int(OPENBUDDY_LLAMA3),
		Sep:            "\n",
		Roles:          []string{"user", "assistant"},
		StopStr:        "",
	})
	return tp
}

func (s *fsChatApiClient) genParams(ctx context.Context, req openai.ChatCompletionRequest, workerAddress string) (params GenerateStreamParams, err error) {

	conv, ok := s.template.GetByModelName(ctx, req.Model)
	if !ok {
		err = errors.New("failed to get conv template")
		return
	}
	convTemplate := ModelConvTemplate{Conv: conv}
	//convTemplate, err := s.options.workerSvc.WorkerGetConvTemplate(ctx, workerAddress, req.Model)
	//if err != nil {
	//	err = errors.WithMessage(err, "failed to get conv template")
	//	return
	//}
	prompt := s.getPrompt(ctx, req, convTemplate)
	//prompt += convTemplate.Conv.Roles[1] + "\n"

	if req.Stop == nil && convTemplate.Conv.StopStr != "" {
		req.Stop = append(req.Stop, convTemplate.Conv.StopStr)
	}

	genParams := GenerateStreamParams{
		Model:            req.Model,
		Prompt:           prompt,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		TopK:             -1,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxNewTokens:     req.MaxTokens,
		StopTokenIds:     convTemplate.Conv.StopTokenIds,
		Stop:             req.Stop,
		Echo:             false,
	}
	if req.N > 0 {
		genParams.N = &req.N
	}
	if req.LogProbs {
		logProbs := true
		genParams.Logprobs = &logProbs
	}

	return genParams, nil
}

func (s *fsChatApiClient) getPrompt(ctx context.Context, req openai.ChatCompletionRequest, convTemplate ModelConvTemplate) (prompt string) {
	// todo 应该从数据库模型模版获取
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role: "assistant",
	})
	var reqSystemMessage string
	for _, v := range req.Messages {
		if v.Role == "system" {
			reqSystemMessage = strings.TrimSpace(v.Content)
			break
		}
	}
	// 去除系统消息
	if convTemplate.Conv.SystemTemplate == "" {
		convTemplate.Conv.SystemTemplate = "{system_message}"
	}
	var systemMessage = strings.ReplaceAll(convTemplate.Conv.SystemTemplate, "{system_message}", reqSystemMessage)
	var ret string
	switch SeparatorStyle(convTemplate.Conv.SepStyle) {
	case LLAMA2:
		seps := []string{convTemplate.Conv.Sep, convTemplate.Conv.Sep2}
		if systemMessage != "" {
			ret = systemMessage
		} else {
			ret = "[INST] "
		}
		for i, v := range req.Messages {
			tag := convTemplate.Conv.Roles[i%2]
			if v.Content != "" {
				if i == 0 {
					ret += strings.TrimSpace(v.Content) + " "
				} else {
					ret += tag + " " + strings.TrimSpace(v.Content) + seps[i%2]
				}
			} else {
				ret += tag
			}
		}
	case LLAMA3:
		ret = systemMessage
		for _, v := range req.Messages {
			if v.Role == "system" && v.Content == "" {
				continue
			}
			if v.Content != "" {
				ret += fmt.Sprintf("<|start_header_id|>%s<|end_header_id|>\n\n%s<|eot_id|>", v.Role, strings.TrimSpace(v.Content))
			} else {
				ret += fmt.Sprintf("<|start_header_id|>%s<|end_header_id|>\n\n", v.Role)
			}
		}
	case CHATML:
		if systemMessage != "" {
			ret = systemMessage + convTemplate.Conv.Sep + "\n"
		}
		for _, v := range req.Messages {
			if v.Role == "system" {
				continue
			}
			if v.MultiContent != nil {
				for _, content := range v.MultiContent {
					if content.Type == openai.ChatMessagePartTypeImageURL {
						if content.ImageURL != nil {
							ret += ImagePlaceholderStr + content.ImageURL.URL
						}
					} else {
						ret += content.Text
					}
				}
			} else if v.Content != "" {
				if v.Role == "user" {
					ret += fmt.Sprintf("%s\n%s%s\n", convTemplate.Conv.Roles[0], v.Content, convTemplate.Conv.Sep)
				} else if v.Role == "assistant" {
					ret += fmt.Sprintf("%s\n%s%s\n", convTemplate.Conv.Roles[1], strings.TrimSpace(v.Content), convTemplate.Conv.Sep)
				} else {
					ret += fmt.Sprintf("%s\n%s%s\n", v.Role, strings.TrimSpace(v.Content), convTemplate.Conv.Sep)
				}
			} else {
				ret += fmt.Sprintf("%s\n", convTemplate.Conv.Roles[1])
			}
		}
	case OPENBUDDY_LLAMA3:
		//ret = systemMessage
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += fmt.Sprintf("<|role|>%s<|says|>%s<|end|>\n", v.Role, strings.TrimSpace(v.Content))
			} else {
				ret += fmt.Sprintf("<|role|>%s<|says|>\n", v.Role)
			}
		}
	case NO_COLON_SINGLE:
		ret = systemMessage
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += v.Role + strings.TrimSpace(v.Content) + convTemplate.Conv.Sep
			} else {
				ret += v.Role
			}
		}
	case CHATGLM3:
		if systemMessage != "" {
			ret += systemMessage + convTemplate.Conv.Sep
		}
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += fmt.Sprintf("%s\n%s%s\n", v.Role, strings.TrimSpace(v.Content), convTemplate.Conv.Sep)
			} else {
				ret += fmt.Sprintf("%s\n", v.Role)
			}
		}
	case CHATGLM:
		// source: https://huggingface.co/THUDM/chatglm-6b/blob/1d240ba371910e9282298d4592532d7f0f3e9f3e/modeling_chatglm.py#L1302-L1308
		// source2: https://huggingface.co/THUDM/chatglm2-6b/blob/e186c891cf64310ac66ef10a87e6635fa6c2a579/modeling_chatglm.py#L926
		roundAddN := 1
		if strings.Contains(req.Model, "chatglm2") {
			roundAddN = 0
		}
		if systemMessage != "" {
			ret = systemMessage + convTemplate.Conv.Sep
		}
		for i, v := range req.Messages {
			if i%2 == 0 {
				ret += fmt.Sprintf("[Round %d]%s", i/2+roundAddN, convTemplate.Conv.Sep)
			}
			if v.Content != "" {
				ret += fmt.Sprintf("%s：%s%s", v.Role, strings.TrimSpace(v.Content), convTemplate.Conv.Sep)
			} else {
				ret += fmt.Sprintf("%s：", v.Role)
			}
		}
	}

	return ret
}

func (s *fsChatApiClient) processInput(modelName string, inp any) (newInp []string) {
	fmt.Println(reflect.TypeOf(inp))
	if reflect.TypeOf(inp).Name() == "string" {
		newInp = []string{inp.(string)}
	} else if reflect.TypeOf(inp).Name() == "[]any" {
		fastInp := inp.([]any)
		if reflect.TypeOf(fastInp[0]).Name() == "int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				model := "cl100k_base"
				decoding, err = tiktoken.GetEncoding(model)
			}
			newInp = []string{decoding.Decode(inp.([]int))}
		} else if reflect.TypeOf(fastInp[0]).Name() == "[]int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				model := "cl100k_base"
				decoding, err = tiktoken.GetEncoding(model)
			}
			for _, text := range inp.([][]int) {
				newInp = append(newInp, decoding.Decode(text))
			}
		}
	}
	return
}
