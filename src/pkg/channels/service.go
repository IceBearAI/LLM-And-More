package channels

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/tokenizers"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/channel"
	"github.com/IceBearAI/aigc/src/repository/model"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
)

type Service interface {
	CreateChannel(ctx context.Context, request CreateChannelRequest) (resp Channel, err error)
	UpdateChannel(ctx context.Context, request UpdateChannelRequest) (err error)
	ListChannel(ctx context.Context, request ListChannelRequest) (resp ChannelList, err error)
	DeleteChannel(ctx context.Context, id uint) (err error)
	GetChannel(ctx context.Context, id uint) (resp Channel, err error)
	ListChannelModels(ctx context.Context, request ListChannelModelsRequest) (resp ChannelModelList, err error)
	ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error)
	// GetModel 获取模型基本信息
	GetModel(ctx context.Context, modelName string) (res modelInfoResult, err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  services.Service
	fileSvc files.Service
}

func (s *service) GetModel(ctx context.Context, modelName string) (res modelInfoResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	modelInfo, err := s.store.Model().FindByModelId(ctx, modelName, "ModelDeploy", "FineTuningTrainJob")
	if err != nil {
		_ = level.Warn(logger).Log("store.Model", "FindByModelId", "err", err.Error())
		return
	}
	res = modelInfoResult{
		ModelName:     modelInfo.ModelName,
		ModelType:     string(modelInfo.ModelType),
		ProviderName:  string(modelInfo.ProviderName),
		BaseModelName: modelInfo.BaseModelName,
		MaxTokens:     modelInfo.MaxTokens,
		Enabled:       modelInfo.Enabled,
		Remark:        modelInfo.Remark,
		CreatedAt:     modelInfo.CreatedAt,
	}
	if modelInfo.ModelDeploy.ModelID > 0 {
		res.Deployment = modelDeploymentResult{
			VLLM:         modelInfo.ModelDeploy.Vllm,
			Status:       modelInfo.ModelDeploy.Status,
			ModelPath:    modelInfo.ModelDeploy.ModelPath,
			Replicas:     modelInfo.ModelDeploy.Replicas,
			InferredType: modelInfo.ModelDeploy.InferredType,
			GPU:          modelInfo.ModelDeploy.Gpu,
			Quantization: modelInfo.ModelDeploy.Quantization,
		}
	}
	if modelInfo.IsFineTuning {
		res.FineTuned = &modelFineTuneResult{
			JobId:   modelInfo.FineTuningTrainJob.JobId,
			FileId:  modelInfo.FineTuningTrainJob.FileId,
			FileUrl: modelInfo.FineTuningTrainJob.FileUrl,
		}
		if fileBody, err := util.GetHttpFileBody(modelInfo.FineTuningTrainJob.FileUrl); err == nil {
			// 取body第一行数据解析成json
			res.SystemPrompt = tokenizers.GetFirstLineSystemPrompt(fileBody)
		} else {
			_ = level.Warn(logger).Log("util.GetHttpFileBody", "err", err.Error())
		}
	}
	return res, nil
}

func (s *service) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ChatCompletionStream")
	completionStream, err := s.apiSvc.FastChat().CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       request.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
	})
	if err != nil {
		_ = level.Error(logger).Log("apiSvc.PaasChat", "ChatCompletionStream", "err", err.Error())
		return stream, err
	}

	dot := make(chan CompletionsStreamResult)
	go func(completionStream *openai.ChatCompletionStream, dot chan CompletionsStreamResult) {
		var fullContent string
		defer func() {
			completionStream.Close()
			close(dot)
		}()
		for {
			completion, err := completionStream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				_ = level.Error(logger).Log("completionStream", "Recv", "err", err.Error())
				return
			}
			begin := time.Now()
			fullContent += completion.Choices[0].Delta.Content
			dot <- CompletionsStreamResult{
				FullContent: fullContent,
				Content:     completion.Choices[0].Delta.Content,
				CreatedAt:   begin,
				ContentType: "text",
				MessageId:   "",
				Model:       "",
				TopP:        0,
				Temperature: 0,
				MaxTokens:   0,
			}
		}
	}(completionStream, dot)
	return dot, nil
}

func (s *service) ListChannelModels(ctx context.Context, request ListChannelModelsRequest) (resp ChannelModelList, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListChannelModels")
	if request.TenantId == types.SystemTenant {

		//五维图时，不判断状态
		var enabled *bool
		if request.EvalTag != string(types.EvaluateTargetTypeFive) {
			e := true
			enabled = &e
		}

		req := model.ListModelRequest{
			Page:          -1,
			PageSize:      -1,
			Enabled:       enabled,
			ProviderName:  request.ProviderName,
			ModelType:     request.ModelType,
			BaseModelName: request.BaseModelName,
		}
		res, total, err := s.store.Model().ListModels(ctx, req)
		if err != nil {
			_ = level.Error(logger).Log("store.Model", "ListModels", "err", err.Error())
			return resp, err
		}
		resp.Total = total
		resp.Models = make([]Model, 0)
		for _, v := range res {

			// 查询有五维图评测试数据的模型
			if request.EvalTag == string(types.EvaluateTargetTypeFive) {
				//判断是否有有五维图数据
				isFive, _ := s.store.ModelEvaluate().IsExistFiveByModelId(ctx, v.ID)
				if isFive == false {
					continue
				}
			}

			resp.Models = append(resp.Models, convertModel(&v))
		}
		return resp, nil
	}
	res, err := s.store.Model().FindModelsByTenantId(ctx, request.TenantId)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "FindModelsByTenantId", "err", err.Error())
		return resp, err
	}
	resp.Models = make([]Model, 0)
	for _, v := range res {
		if v.Enabled {
			// 查询有五维图评测试数据的模型
			if request.EvalTag == string(types.EvaluateTargetTypeFive) {
				//判断是否有有五维图数据
				isFive, _ := s.store.ModelEvaluate().IsExistFiveByModelId(ctx, v.ID)
				if isFive == false {
					continue
				}
			}

			resp.Models = append(resp.Models, convertModel(&v))
		}
	}
	resp.Total = int64(len(resp.Models))
	return
}

func (s *service) GetChannel(ctx context.Context, id uint) (resp Channel, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetChannel")
	res, err := s.store.Channel().GetChannel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "GetChannel", "err", err.Error(), "id", id)
		return
	}
	resp = convert(&res)
	return
}

func (s *service) CreateChannel(ctx context.Context, request CreateChannelRequest) (resp Channel, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateChannel")
	data := &types.ChatChannels{
		Name:         uuid.New().String(),
		Alias:        request.Alias,
		Remark:       request.Remark,
		Quota:        request.Quota,
		ApiKey:       "sk-" + string(util.Krand(48, util.KC_RAND_KIND_ALL)),
		Email:        request.Email,
		LastOperator: request.LastOperator,
		TenantId:     request.TenantId,
		ModelId:      request.ModelId,
	}
	err = s.store.Channel().CreateChannel(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "CreateChannel", "err", err.Error())
		return
	}
	resp = convert(data)
	return
}

func (s *service) UpdateChannel(ctx context.Context, request UpdateChannelRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateChannel")
	res, err := s.store.Channel().GetChannel(ctx, request.Id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "FindChannelById", "err", err.Error())
		return
	}
	if request.Name != nil {
		res.Name = *request.Name
	}
	if request.Alias != nil {
		res.Alias = *request.Alias
	}
	if request.Quota != nil {
		res.Quota = *request.Quota
	}
	if request.Email != nil {
		res.Email = *request.Email
	}
	if request.Remark != nil {
		res.Remark = *request.Remark
	}
	res.ModelId = request.ModelId
	res.UpdatedAt = time.Now()
	err = s.store.Channel().UpdateChannel(ctx, &res)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "UpdateChannel", "err", err.Error(), "id", request.Id)
		return
	}
	return
}

func (s *service) ListChannel(ctx context.Context, request ListChannelRequest) (resp ChannelList, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListChannel")
	resp.Channels = make([]Channel, 0)
	listReq := channel.ListChannelRequest{
		Page:        request.Page,
		PageSize:    request.PageSize,
		Name:        request.Name,
		Email:       request.Email,
		ProjectName: request.ProjectName,
		ServiceName: request.ServiceName,
		TenantId:    request.TenantId,
	}
	res, total, err := s.store.Channel().ListChannels(ctx, listReq)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "ListChannels", "err", err.Error())
		return
	}
	for _, v := range res {
		resp.Channels = append(resp.Channels, convert(&v))
	}
	resp.Total = total
	return
}

func convert(data *types.ChatChannels) Channel {
	c := Channel{
		Id:           data.ID,
		Name:         data.Name,
		Alias:        data.Alias,
		Quota:        data.Quota,
		ApiKey:       data.ApiKey,
		Email:        data.Email,
		Remark:       data.Remark,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		TenantId:     data.TenantId,
		LastOperator: data.LastOperator,
	}
	models := make([]Model, 0)
	for _, v := range data.ChannelModels {
		models = append(models, Model{
			Id:           v.ID,
			ProviderName: v.ProviderName.String(),
			ModelType:    v.ModelType.String(),
			ModelName:    v.ModelName,
			MaxTokens:    v.MaxTokens,
			//IsPrivate:    v.IsPrivate,
			Remark:       v.Remark,
			Enabled:      v.Enabled,
			IsFineTuning: v.IsFineTuning,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}
	c.Model.Num = len(models)
	c.Model.List = models
	return c
}

func convertModel(data *types.Models) Model {
	m := Model{
		Id:           data.ID,
		ProviderName: data.ProviderName.String(),
		ModelType:    data.ModelType.String(),
		ModelName:    data.ModelName,
		MaxTokens:    data.MaxTokens,
		//IsPrivate:     data.IsPrivate,
		IsFineTuning:  data.IsFineTuning,
		Enabled:       data.Enabled,
		Remark:        data.Remark,
		BaseModelName: data.BaseModelName,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}
	return m
}

func (s *service) DeleteChannel(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteChannel")
	err = s.store.Channel().DeleteChannel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "DeleteChannel", "err", err.Error(), "id", id)
		return
	}
	return
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc services.Service) Service {
	return &service{
		logger:  log.With(logger, "service", "channels"),
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
	}
}
