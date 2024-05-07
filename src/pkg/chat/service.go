package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/chat"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"strings"
	"time"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	workerSvc      chat.WorkerService
	openAISvc      chat.OpenAIService
}

// CreationOption is the option for the chat service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOpts is the option to set the http client options.
func WithHTTPClientOpts(opts ...kithttp.ClientOption) CreationOption {
	return func(o *CreationOptions) {
		o.httpClientOpts = opts
	}
}

// WithWorkerService is the option to set the worker service.
func WithWorkerService(svc chat.WorkerService) CreationOption {
	return func(o *CreationOptions) {
		o.workerSvc = svc
	}
}

type Service interface {
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error)
	// Models 模型列表
	Models(ctx context.Context, channelId uint) (res []openai.Model, err error)
	// Embeddings 向量化处理
	Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}

type CompletionStreamResponse struct {
	Usage openai.Usage `json:"usage"`
	openai.ChatCompletionStreamResponse
}

type service struct {
	traceId    string
	logger     log.Logger
	options    *CreationOptions
	services   services.Service
	repository repository.Repository
}

func (s *service) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) processInput(modelName string, inp any) (newInp []string) {
	fmt.Println(reflect.TypeOf(inp))
	if reflect.TypeOf(inp).Name() == "string" {
		newInp = []string{inp.(string)}
	} else if reflect.TypeOf(inp).Name() == "[]any" {
		fastInp := inp.([]any)
		if reflect.TypeOf(fastInp[0]).Name() == "int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				_ = level.Warn(s.logger).Log("msg", "model not found. Using cl100k_base encoding.")
				model := "cl100k_base"
				decoding, err = tiktoken.GetEncoding(model)
			}
			newInp = []string{decoding.Decode(inp.([]int))}
		} else if reflect.TypeOf(fastInp[0]).Name() == "[]int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				_ = level.Warn(s.logger).Log("msg", "model not found. Using cl100k_base encoding.")
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

func (s *service) getPrompt(ctx context.Context, req openai.ChatCompletionRequest, convTemplate chat.ModelConvTemplate) (prompt string) {
	// todo 应该从数据库模型模版获取

	var reqSystemMessage string
	for _, v := range req.Messages {
		if v.Role == "system" {
			reqSystemMessage = v.Content
			break
		}
	}
	if convTemplate.Conv.SystemTemplate == "" {
		convTemplate.Conv.SystemTemplate = "{system_message}"
	}
	var systemMessage = strings.ReplaceAll(convTemplate.Conv.SystemTemplate, "{system_message}", reqSystemMessage)
	var ret string
	switch chat.SeparatorStyle(convTemplate.Conv.SepStyle) {
	case chat.LLAMA2:
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
					ret += v.Content + " "
				} else {
					ret += tag + " " + v.Content + seps[i%2]
				}
			} else {
				ret += tag
			}
		}
	case chat.LLAMA3:
		ret = systemMessage
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += fmt.Sprintf("<|start_header_id|>%s<|end_header_id|>\n\n%s<|eot_id|>", v.Role, v.Content)
			} else {
				ret += fmt.Sprintf("<|start_header_id|>%s<|end_header_id|>\n\n", v.Role)
			}
		}
	case chat.CHATML:
		if systemMessage != "" {
			ret += systemMessage + convTemplate.Conv.Sep + "\n"
		}
		for _, v := range req.Messages {
			if v.MultiContent != nil {
				for _, content := range v.MultiContent {
					if content.Type == openai.ChatMessagePartTypeImageURL {
						if content.ImageURL != nil {
							ret += chat.ImagePlaceholderStr + content.ImageURL.URL
						}
					} else {
						ret += content.Text
					}
				}
			} else if v.Content != "" {
				if v.Role == "user" {
					ret += fmt.Sprintf("%s\n%s%s\n", convTemplate.Conv.Roles[0], v.Content, convTemplate.Conv.Sep)
				} else if v.Role == "assistant" {
					ret += fmt.Sprintf("%s\n%s%s\n", convTemplate.Conv.Roles[1], v.Content, convTemplate.Conv.Sep)
				} else {
					ret += fmt.Sprintf("%s\n%s%s\n", v.Role, v.Content, convTemplate.Conv.Sep)
				}
			} else {
				ret += fmt.Sprintf("%s\n", convTemplate.Conv.Roles[1])
			}
		}
	case chat.OPENBUDDY_LLAMA3:
		ret = systemMessage + "\n"
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += fmt.Sprintf("<|role|>%s<|says|>%s<|end|>\n", v.Role, v.Content)
			} else {
				ret += fmt.Sprintf("<|role|>%s<|says|>\n", v.Role)
			}
		}
	case chat.NO_COLON_SINGLE:
		ret = systemMessage
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += v.Role + v.Content + convTemplate.Conv.Sep
			} else {
				ret += v.Role
			}
		}
	case chat.CHATGLM3:
		if systemMessage != "" {
			ret += systemMessage + convTemplate.Conv.Sep
		}
		for _, v := range req.Messages {
			if v.Content != "" {
				ret += fmt.Sprintf("%s\n%s%s\n", v.Role, v.Content, convTemplate.Conv.Sep)
			} else {
				ret += fmt.Sprintf("%s\n", v.Role)
			}
		}
	case chat.CHATGLM:
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
				ret += fmt.Sprintf("%s：%s%s", v.Role, v.Content, convTemplate.Conv.Sep)
			} else {
				ret += fmt.Sprintf("%s：", v.Role)
			}
		}
	}

	return ret
}

func (s *service) genParams(ctx context.Context, req openai.ChatCompletionRequest, workerAddress string) (params chat.GenerateStreamParams, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	convTemplate, err := s.options.workerSvc.WorkerGetConvTemplate(ctx, workerAddress, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to get conv template")
		_ = level.Warn(logger).Log("msg", "failed to get conv template", "err", err)
		return
	}
	_ = level.Debug(logger).Log("conv.SystemTemplate", convTemplate.Conv.SystemTemplate)
	prompt := s.getPrompt(ctx, req, convTemplate)
	prompt += convTemplate.Conv.Roles[1]

	_ = level.Info(logger).Log("conv.SepStyle", convTemplate.Conv.SepStyle, "prompt", prompt)

	return chat.GenerateStreamParams{
		Model:            req.Model,
		Prompt:           prompt,
		Temperature:      req.Temperature,
		Logprobs:         req.LogProbs,
		TopP:             req.TopP,
		TopK:             -1,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxNewTokens:     req.MaxTokens,
		StopTokenIds:     convTemplate.Conv.StopTokenIds,
		Images:           nil,
		UseBeamSearch:    false,
		Stop:             req.Stop,
	}, nil
}

func (s *service) localAIChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (res <-chan CompletionStreamResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	models, err := s.options.workerSvc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		_ = level.Warn(logger).Log("msg", "failed to list models", "err", err)
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
		_ = level.Warn(logger).Log("msg", "model not found", "err", err)
		return
	}
	workerAddress, err := s.options.workerSvc.GetWorkerAddress(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		_ = level.Warn(logger).Log("msg", "failed to get worker address", "err", err)
		return
	}
	_ = level.Info(logger).Log("msg", "worker address", "workerAddress", workerAddress)
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
		_ = level.Warn(logger).Log("msg", "failed to get gen params", "err", err)
		return
	}
	stream, err := s.options.workerSvc.WorkerGenerateStream(ctx, workerAddress, genParams)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		_ = level.Error(logger).Log("msg", "failed to generate stream", "err", err)
		return
	}

	go func() {
		now := time.Now().UnixMilli()
		defer close(dot)
		streamId := fmt.Sprintf("cmpl-%s", shortuuid.New())
		for {
			content, ok := <-stream
			if !ok {
				_ = level.Info(logger).Log("msg", "stream closed")
				break
			}
			if content.ErrorCode != 0 {
				err = errors.New(content.Text)
				_ = level.Warn(logger).Log("msg", "error code", "errorCode", content.ErrorCode, "text", content.Text)
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
								Content: content.Text,
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

func (s *service) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return
	}
	svc := chat.NewFastChatWorker()
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		models, err := svc.ListModels(ctx)
		if err != nil {
			err = errors.WithMessage(err, "failed to list models")
			_ = level.Warn(logger).Log("msg", "failed to list models", "err", err)
			return res, err
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
			_ = level.Warn(logger).Log("msg", "model not found", "err", err)
			return res, err
		}
		workerAddress, err := svc.GetWorkerAddress(ctx, modelInfo.ModelName)
		if err != nil {
			err = errors.WithMessage(err, "failed to get worker address")
			_ = level.Warn(logger).Log("msg", "failed to get worker address", "err", err)
			return res, err
		}
		workerResult, err := svc.WorkerGenerate(ctx, workerAddress, chat.GenerateParams{
			Model:            req.Model,
			Prompt:           "",
			Temperature:      req.Temperature,
			Logprobs:         req.TopLogProbs,
			TopP:             req.TopP,
			PresencePenalty:  req.PresencePenalty,
			FrequencyPenalty: req.FrequencyPenalty,
			MaxNewTokens:     0,
			Echo:             false,
			StopTokenIds:     nil,
			Images:           nil,
			BestOf:           0,
			UseBeamSearch:    false,
			Stop:             req.Stop,
		})
		fmt.Println(workerResult)
	}
	return
}

func (s *service) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	dot := make(chan CompletionStreamResponse)
	defer func() {
		if err != nil {
			close(dot)
		}
	}()
	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return dot, err
	}
	isError := true
	finished := false
	b, _ := json.Marshal(req.Messages)
	stop, _ := json.Marshal(req.Stop)
	msgData := &types.ChatMessages{
		ModelName:        req.Model,
		ChannelId:        channelId,
		Prompt:           req.Messages[len(req.Messages)-1].Content,
		Finished:         &finished,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		N:                req.N,
		User:             req.User,
		Messages:         string(b),
		Error:            &isError,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxTokens:        req.MaxTokens,
		Stop:             string(stop),
	}

	if err = s.repository.Messages().Create(ctx, msgData); err != nil {
		err = errors.WithMessage(err, "failed to create message")
		_ = level.Warn(logger).Log("msg", "failed to create message", "err", err)
		return
	}

	if modelInfo.BaseModelName != "" {
		req.Model = modelInfo.BaseModelName
	}
	completionStream, err := s.localAIChatCompletionStream(ctx, req)
	if err != nil {
		msgData.ErrorMessage = err.Error()
		_ = level.Warn(logger).Log("msg", "failed to get completion stream", "err", err)
		return dot, err
	}

	go func() {
		defer close(dot)
		for content := range completionStream {
			if content.Choices[0].FinishReason == openai.FinishReasonStop {
				isError = false
				finished = true
				// 更新数据库
				msgData.Response = content.Choices[0].Delta.Content
				msgData.Error = &isError
				msgData.Created = content.Created
				msgData.TimeCost = time.Since(msgData.CreatedAt).String()
				msgData.Finished = &finished
				msgData.PromptTokens = content.Usage.PromptTokens
				msgData.ResponseTokens = content.Usage.CompletionTokens
				if err = s.repository.Messages().Update(ctx, msgData); err != nil {
					_ = level.Warn(logger).Log("msg", "failed to update message", "err", err)
				}
				_ = level.Info(logger).Log("msg", "chat completion stream finished", "usage", fmt.Sprintf("%+v", content.Usage))
			}
			dot <- content
		}
	}()

	return dot, nil
}

func New(logger log.Logger, traceId string, store repository.Repository, services services.Service, opts ...CreationOption) Service {
	logger = log.With(logger, "service", "chat")
	options := &CreationOptions{}
	for _, o := range opts {
		o(options)
	}
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
		services:   services,
		options:    options,
	}
}
