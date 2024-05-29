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
	options   *CreationOptions
	template  Templates
	openaiSvc Service
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
				FinishReason: "stop",
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
	var fullContent string
	for {
		rs, ok := <-streamResp
		if !ok {
			break
		}
		if len(rs.Choices) > 0 {
			content = rs
			fullContent += rs.Choices[0].Delta.Content
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
					FinishReason: "stop",
					Message: openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: fullContent,
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
	dot := make(chan CompletionStreamResponse)
	//var maxTokens int

	//prompts := s.processInput(req.Model, req.Messages)
	//for _, prompt := range prompts {
	//	maxTokens, err = s.options.workerSvc.WorkerCheckLength(ctx, workerAddress, req.Model, req.MaxTokens, prompt)
	//	if err != nil {
	//		err = errors.WithMessage(err, "failed to check length")
	//		return dot, err
	//	}
	//}
	//if maxTokens != 0 && maxTokens < req.MaxTokens {
	//	req.MaxTokens = maxTokens
	//}

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
		var previousText string

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
			// 替换所有的Unicode替代字符\ufffd为空字符串
			decodedUnicode := strings.Replace(text, "\ufffd", "", -1)

			// 获取新的字符串，它是当前文本去掉与之前文本相同部分后的结果
			deltaText := decodedUnicode
			if len(previousText) < len(decodedUnicode) {
				deltaText = decodedUnicode[len(previousText):]
			}
			if len(decodedUnicode) > len(previousText) {
				previousText = decodedUnicode
			} else {
				deltaText = ""
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
			if content.FinishReason == "stop" {
				return
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
		template: Register(tp),
	}
}

func (s *fsChatApiClient) genParams(ctx context.Context, req openai.ChatCompletionRequest, workerAddress string) (params GenerateStreamParams, err error) {

	conv, ok := s.template.GetByModelName(ctx, req.Model)
	if !ok {
		err = errors.New("failed to get conv template")
		return
	}
	var imageList []string
	for _, v := range req.Messages {
		if v.Role == "system" {
			conv.SetSystemMessage(strings.TrimSpace(v.Content))
		} else if v.Role == "user" {
			if v.MultiContent != nil {
				var textList []string
				for _, item := range v.MultiContent {
					if item.Type == openai.ChatMessagePartTypeImageURL {
						if item.ImageURL != nil {
							imageList = append(imageList, item.ImageURL.URL)
						}
					} else {
						textList = append(textList, item.Text)
					}
				}
				text := strings.Repeat("<image>\n", len(imageList))
				text += strings.Join(textList, "\n")
				conv.AppendMessage(conv.Roles[0], text)
			} else {
				conv.AppendMessage(conv.Roles[0], strings.TrimSpace(v.Content))
			}
		} else if v.Role == "assistant" {
			conv.AppendMessage(conv.Roles[1], strings.TrimSpace(v.Content))
		}
	}
	conv.AppendMessage(conv.Roles[1], "")

	//convTemplate, err := s.options.workerSvc.WorkerGetConvTemplate(ctx, workerAddress, req.Model)
	//if err != nil {
	//	err = errors.WithMessage(err, "failed to get conv template")
	//	return
	//}
	prompt := conv.GetPrompt()

	if req.Stop == nil && conv.StopStr != nil {
		req.Stop = append(req.Stop, conv.StopStr...)
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
		StopTokenIds:     conv.StopTokenIds,
		Stop:             req.Stop,
		Echo:             false,
		Images:           imageList,
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

func (s *fsChatApiClient) processInput(modelName string, inp any) (newInp []string) {
	//var prompt string
	//switch inp.(type) {
	//case string:
	//	prompt = inp.(string)
	//case []string:
	//	prompt = strings.Join(inp.([]string), " ")
	//case []interface{}:
	//	prompts, _ := ConvertToSliceOfStrings(inp.([]interface{}))
	//	prompt = strings.Join(prompts, " ")
	//}

	fmt.Println(reflect.TypeOf(inp))
	if reflect.TypeOf(inp).Name() == "string" {
		newInp = []string{inp.(string)}
	} else if reflect.TypeOf(inp).Name() == "[]openai.ChatCompletionMessage" {
		fastInp := inp.([]openai.ChatCompletionMessage)
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

func ConvertToSliceOfStrings(data interface{}) ([]string, bool) {
	var result []string
	slice, ok := data.([]interface{})
	if !ok {
		return nil, false
	}
	for _, item := range slice {
		if num, ok := item.(string); ok {
			result = append(result, num)
		} else {
			return nil, false
		}
	}
	return result, true
}
