package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/chat"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"time"
)

// Service 聊天服务接口
type Service interface {
	// Completion 生成
	Completion(ctx context.Context, channelId uint, req openai.CompletionRequest) (res openai.CompletionResponse, err error)
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan chat.CompletionStreamResponse, err error)
	// Models 模型列表
	Models(ctx context.Context, channelId uint) (res []openai.Model, err error)
	// Embeddings 向量化处理
	Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	services   services.Service
	repository repository.Repository
}

func (s *service) Completion(ctx context.Context, channelId uint, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	now := time.Now()
	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return
	}

	if modelInfo.BaseModelName != "" {
		req.Model = modelInfo.BaseModelName
	}
	providerName := services.ProviderOpenAI
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		providerName = services.ProviderFsChat
	}
	res, err = s.services.Chat(providerName).Completion(ctx, req)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "failed to get completion", "err", err)
		return
	}
	var finished = true
	prompt, _ := json.Marshal(req.Prompt)
	var stop []byte
	if req.Stop != nil {
		stop, _ = json.Marshal(req.Stop)
	}
	msgData := &types.ChatMessages{
		ModelName:        req.Model,
		ChannelId:        channelId,
		Prompt:           string(prompt),
		Finished:         &finished,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		N:                req.N,
		User:             req.User,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxTokens:        req.MaxTokens,
		Stop:             string(stop),
		Response:         res.Choices[0].Text,
		Created:          res.Created,
		TimeCost:         time.Since(now).String(),
		PromptTokens:     res.Usage.PromptTokens,
		ResponseTokens:   res.Usage.CompletionTokens,
	}

	if err = s.repository.Messages().Create(ctx, msgData); err != nil {
		err = errors.WithMessage(err, "failed to create message")
		_ = level.Warn(logger).Log("msg", "failed to create message", "err", err)
		return
	}

	return
}

func (s *service) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	channelInfo, err := s.repository.Channel().FindChannelById(ctx, channelId, "ChannelModels")
	if err != nil {
		err = errors.WithMessage(err, "failed to find channel")
		_ = level.Warn(logger).Log("msg", "failed to find channel", "err", err)
		return nil, err
	}

	for _, v := range channelInfo.ChannelModels {
		res = append(res, openai.Model{
			ID:        v.ModelName,
			Object:    "model",
			Root:      v.ModelName,
			CreatedAt: v.CreatedAt.Unix(),
		})
	}

	return
}

func (s *service) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return res, err
	}
	isError := true
	finished := false
	b, _ := json.Marshal(req.Messages)
	var stop []byte
	if req.Stop != nil {
		stop, _ = json.Marshal(req.Stop)
	}
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
	providerName := services.ProviderOpenAI
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		providerName = services.ProviderFsChat
	}
	completionResult, err := s.services.Chat(providerName).ChatCompletion(ctx, req)
	if err != nil {
		msgData.ErrorMessage = err.Error()
		// todo: 更新数据库
		_ = level.Warn(logger).Log("msg", "failed to get completion stream", "err", err)
		return
	}

	var resContent string
	res = openai.ChatCompletionResponse{
		ID:                completionResult.ID,
		Object:            completionResult.Object,
		Created:           completionResult.Created,
		Model:             completionResult.Model,
		Choices:           completionResult.Choices,
		Usage:             completionResult.Usage,
		SystemFingerprint: completionResult.SystemFingerprint,
	}

	if len(res.Choices) > 0 {
		resContent = res.Choices[0].Message.Content
	}

	isError = false
	finished = true
	// 更新数据库
	msgData.Response = resContent
	msgData.Error = &isError
	msgData.Created = time.Now().Unix()
	msgData.TimeCost = time.Since(msgData.CreatedAt).String()
	msgData.Finished = &finished
	msgData.PromptTokens = res.Usage.PromptTokens
	msgData.ResponseTokens = res.Usage.CompletionTokens
	if err = s.repository.Messages().Update(ctx, msgData); err != nil {
		_ = level.Warn(logger).Log("msg", "failed to update message", "err", err)
	}
	_ = level.Info(logger).Log("msg", "chat completion stream finished", "usage", fmt.Sprintf("%+v", res.Usage))
	return
}

func (s *service) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan chat.CompletionStreamResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	dot := make(chan chat.CompletionStreamResponse)
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
	var stop []byte
	if req.Stop != nil {
		stop, _ = json.Marshal(req.Stop)
	}
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
	if modelInfo.BaseModelName != "" {
		req.Model = modelInfo.BaseModelName
	}
	_ = level.Info(logger).Log("model", req.Model, "providerName", modelInfo.ProviderName)
	providerName := services.ProviderOpenAI
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		providerName = services.ProviderFsChat
	}
	if modelInfo.BaseModelName != "" {
		req.Model = modelInfo.BaseModelName
	}
	completionStream, err := s.services.Chat(providerName).ChatCompletionStream(ctx, req)
	if err != nil {
		msgData.ErrorMessage = err.Error()
		// 更新错误信息到数据表
		_ = level.Warn(logger).Log("msg", "failed to get completion stream", "err", err)
		return
	}

	go func() {
		defer close(dot)
		var resContent string
		for content := range completionStream {
			if util.StringInArray([]string{string(services.ProviderFsChat), string(services.ProviderFsChat)}, string(providerName)) {
				if len(content.Choices) > 0 && len(content.Choices[0].Delta.Content) >= len(resContent) && len(resContent) > 0 {
					content.Choices[0].Delta.Content = content.Choices[0].Delta.Content[len(resContent):]
					resContent += content.Choices[0].Delta.Content
				}
			}
			if len(content.Choices) > 0 && content.Choices[0].FinishReason == openai.FinishReasonStop {
				isError = false
				finished = true
				// 更新数据库
				msgData.Response = resContent
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

func New(logger log.Logger, traceId string, store repository.Repository, services services.Service) Service {
	logger = log.With(logger, "service", "chat")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
		services:   services,
	}
}
