package models

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/tokenizers"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/model"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/runtime"
	"github.com/IceBearAI/aigc/src/util"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

type Middleware func(Service) Service

// Service is the interface for the model service.
// Service 模型服务接口
// 应该从中间件去判断该租户是否有权限访问相关模型，后续多租户的时候再处理
type Service interface {
	// ListModels 模型分页列表
	ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error)
	// CreateModel 创建模型
	CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error)
	// GetModel 获取模型
	GetModel(ctx context.Context, id uint) (res Model, err error)
	// UpdateModel 更新模型
	UpdateModel(ctx context.Context, request UpdateModelRequest) (err error)
	// DeleteModel 删除模型
	DeleteModel(ctx context.Context, id uint) (err error)
	// Deploy 模型部署
	Deploy(ctx context.Context, request ModelDeployRequest) (err error)
	// Undeploy 模型取消部署
	Undeploy(ctx context.Context, id uint) (err error)
	// CreateEval 创建评估任务
	// Deprecated: use modelevaluate.Create instead
	CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error)
	// ListEval 评估任务分页列表
	// Deprecated: use modelevaluate.List instead
	ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error)
	// CancelEval 取消评估任务
	// Deprecated: use modelevaluate.Cancel instead
	CancelEval(ctx context.Context, id uint) (err error)
	// DeleteEval 删除评估任务
	// Deprecated: use modelevaluate.Delete instead
	DeleteEval(ctx context.Context, id uint) (err error)
	// GetModelLogs 获取模型输出日志
	GetModelLogs(ctx context.Context, modelName, containerName string) (res string, err error)
	// ChatCompletionStream 聊天补全流
	ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error)
	// ModelInfo 获取模型信息
	ModelInfo(ctx context.Context, modelName string) (res modelInfoResult, err error)
}

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts     []kithttp.ClientOption
	volumeName         string
	gpuTolerationValue string
	controllerAddress  string
	runtimePlatform    string
	convertUrlFun      func(fileUrl string) string
	defaultTrainImage  string
	dataFs             embed.FS
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOpts returns a CreationOption that sets the http client options.
func WithHTTPClientOpts(opts ...kithttp.ClientOption) CreationOption {
	return func(co *CreationOptions) {
		co.httpClientOpts = opts
	}
}

// WithVolumeName returns a CreationOption that sets the volume name.
func WithVolumeName(volumeName string) CreationOption {
	return func(co *CreationOptions) {
		co.volumeName = volumeName
	}
}

// WithGPUTolerationValue returns a CreationOption that sets the gpu toleration value.
func WithGPUTolerationValue(gpuTolerationValue string) CreationOption {
	return func(co *CreationOptions) {
		co.gpuTolerationValue = gpuTolerationValue
	}
}

// WithControllerAddress returns a CreationOption that sets the controller address.
func WithControllerAddress(controllerAddress string) CreationOption {
	return func(co *CreationOptions) {
		co.controllerAddress = controllerAddress
	}
}

// WithRuntimePlatform returns a CreationOption that sets the controller address.
func WithRuntimePlatform(platform string) CreationOption {
	return func(co *CreationOptions) {
		co.runtimePlatform = platform
	}
}

// WithConvertUrlFun returns a CreationOption that sets the convert url function.
func WithConvertUrlFun(fun func(fileUrl string) string) CreationOption {
	return func(co *CreationOptions) {
		co.convertUrlFun = fun
	}
}

// WithDefaultTrainImage returns a CreationOption that sets the default train image.
func WithDefaultTrainImage(image string) CreationOption {
	return func(co *CreationOptions) {
		co.defaultTrainImage = image
	}
}

// WithDataFs returns a CreationOption that sets the data fs.
func WithDataFs(dataFs embed.FS) CreationOption {
	return func(co *CreationOptions) {
		co.dataFs = dataFs
	}
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  services.Service
	options *CreationOptions
}

func (s *service) ModelInfo(ctx context.Context, modelName string) (res modelInfoResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	m, err := s.store.Model().FindByModelId(ctx, modelName, "ModelDeploy", "FineTuningTrainJob")
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModelByModelName", "err", err.Error(), "modelName", modelName)
		return res, encode.ErrSystem.Wrap(errors.New("查询模型失败"))
	}
	res = modelInfoResult{
		ModelName:    m.ModelName,
		ProviderName: m.ProviderName.String(),
		ModelType:    m.ModelType.String(),
		MaxTokens:    m.MaxTokens,
		Remark:       m.Remark,
		Enabled:      m.Enabled,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if m.ModelDeploy.ID > 0 {
		res.Deployment = &modelDeploymentResult{
			Status:       m.ModelDeploy.Status,
			Replicas:     m.ModelDeploy.Replicas,
			InferredType: m.ModelDeploy.InferredType,
			GPU:          m.ModelDeploy.Gpu,
			Quantization: m.ModelDeploy.Quantization,
			VLLM:         m.ModelDeploy.Vllm,
			ModelPath:    m.ModelDeploy.ModelPath,
		}
	}

	if m.FineTuningTrainJob.ID > 0 {
		res.FineTuned = &modelFineTuneResult{
			FileId:   m.FineTuningTrainJob.FileId,
			JobId:    m.FineTuningTrainJob.JobId,
			FileUrl:  m.FineTuningTrainJob.FileUrl,
			Scenario: string(m.FineTuningTrainJob.Scenario),
		}
		// todo: 如果是标注的内容则直接从标注表中获取，否则读取文件解析
		// 暂时默认从训练文件中获取
		if fileBody, err := util.GetHttpFileBody(s.options.convertUrlFun(m.FineTuningTrainJob.FileUrl)); err == nil {
			// 取body第一行数据解析成json
			prompts, err := tokenizers.GetSystemContent(fileBody, 50)
			if err != nil {
				_ = level.Warn(logger).Log("tokenizers.GetSystemContent", "err", err.Error())
			}
			res.FineTuned.SystemPrompts = prompts
		} else {
			_ = level.Warn(logger).Log("util.GetHttpFileBody", "err", err.Error())
		}
	}

	return res, nil
}

func (s *service) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (stream <-chan CompletionsStreamResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ChatCompletionStream")
	modelInfo, err := s.store.Model().FindByModelId(ctx, request.Model)
	if err != nil {
		_ = level.Warn(logger).Log("store.Model", "FindByModelId", "err", err.Error())
		return stream, err
	}
	if modelInfo.BaseModelName != "" {
		request.Model = modelInfo.BaseModelName
	}
	pn := services.ProviderOpenAI
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		pn = services.ProviderFsChat
	}
	completionStream, err := s.apiSvc.Chat(pn).ChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       request.Model,
		Messages:    request.Messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
	})
	if err != nil {
		_ = level.Error(logger).Log("apiSvc.Chat", "ChatCompletionStream", "err", err.Error())
		return stream, err
	}

	dot := make(chan CompletionsStreamResult)
	go func() {
		var fullContent string
		defer func() {
			close(dot)
		}()
		begin := time.Now()
		for {
			result, ok := <-completionStream
			if !ok {
				return
			}
			fullContent += result.Choices[0].Delta.Content
			dot <- CompletionsStreamResult{
				FullContent: fullContent,
				Content:     result.Choices[0].Delta.Content,
				CreatedAt:   begin,
				ContentType: "text",
				MessageId:   "",
				Model:       "",
				TopP:        0,
				Temperature: 0,
				MaxTokens:   0,
			}
		}
	}()
	return dot, nil
}

func (s *service) GetModelLogs(ctx context.Context, modelName, containerName string) (res string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	modelInfo, err := s.store.Model().GetModelByModelName(ctx, modelName)
	if err != nil {
		_ = level.Warn(logger).Log("store.Model", "GetModelByModelName", "err", err.Error())
		return
	}
	status, err := s.apiSvc.Runtime().GetDeploymentStatus(ctx, fmt.Sprintf("%s-%d", util.ReplacerServiceName(modelName), modelInfo.ID))
	if err != nil {
		_ = level.Error(logger).Log("api.Runtime", "GetDeploymentStatus", "err", err.Error())
		return res, nil
	}
	if strings.ToLower(status) != "running" {
		_ = level.Warn(logger).Log("msg", "deployment not running", "status", status)
		return res, nil
	}
	return s.apiSvc.Runtime().GetDeploymentLogs(ctx, fmt.Sprintf("%s-%d", util.ReplacerServiceName(modelName), modelInfo.ID), containerName)
}

func (s *service) DeleteEval(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteEval")
	eval, err := s.store.Model().GetEval(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetEval", "err", err.Error(), "id", id)
		return encode.ErrSystem.Wrap(errors.New("查询评估任务失败"))
	}

	if utils.Contains([]string{types.EvalStatusPending.String(), types.EvalStatusRunning.String()}, eval.Status.String()) {
		_ = level.Error(logger).Log("msg", "等待中和运行中的任务不可删除", "status", eval.Status.String())
		return encode.InvalidParams.Wrap(errors.New("等待中和运行中的任务不可删除"))
	}
	err = s.store.Model().DeleteEval(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "DeleteEval", "err", err.Error(), "id", id)
		return encode.ErrSystem.Wrap(errors.New("删除评估任务失败"))
	}
	return
}

func (s *service) CancelEval(ctx context.Context, id uint) (err error) {
	eval, err := s.store.Model().GetEval(ctx, id)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.New("查询模型失败"))
	}
	eval.Status = types.EvalStatusCancel
	err = s.store.Model().UpdateEval(ctx, &eval)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.New("取消评估任务失败"))
	}
	return
}

func (s *service) ListEval(ctx context.Context, request ListEvalRequest) (res ListEvalResponse, err error) {
	eval, total, err := s.store.Model().ListEval(ctx, model.ListEvalRequest{
		Page:        request.Page,
		PageSize:    request.PageSize,
		ModelName:   request.ModelName,
		MetricName:  request.MetricName,
		Status:      request.Status,
		DatasetType: request.DatasetType,
	})
	if err != nil {
		return res, encode.ErrSystem.Wrap(errors.New("查询评估任务失败"))
	}
	list := make([]Eval, 0)
	for _, v := range eval {
		list = append(list, convertEval(&v))
	}
	res = ListEvalResponse{
		Total: total,
		List:  list,
	}
	return
}

func (s *service) CreateEval(ctx context.Context, request CreateEvalRequest) (res Eval, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateEval")
	if request.DatasetType == types.EvalDataSetTypeTrain.String() {
		request.EvalPercent = request.EvalPercent / 100.0
		if request.EvalPercent <= 0 || request.EvalPercent > 1 {
			_ = level.Error(logger).Log("msg", "评估比例不正确", "evalPercent", request.EvalPercent)
			return res, encode.InvalidParams.Wrap(errors.New("评估比例不正确"))
		}
		job, err := s.store.FineTuning().GetFineTuningJobByModelName(ctx, request.ModelName)
		if err != nil {
			_ = level.Error(logger).Log("store.FineTuning", "GetFineTuningJobByModelName", "err", err.Error())
			return res, encode.ErrSystem.Wrap(errors.New("查询微调任务失败"))
		}
		request.DatasetFileId = job.FileId
	}

	file, err := s.store.Files().FindFileByFileId(ctx, request.DatasetFileId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, encode.InvalidParams.Wrap(errors.New("文件不存在"))
		}
		_ = level.Error(logger).Log("store.Files", "FindFileByFileId", "err", err.Error())
		return res, encode.ErrSystem.Wrap(errors.New("查询文件失败"))
	}

	if request.DatasetType == types.EvalDataSetTypeTrain.String() && file.Purpose != types.FilePurposeFineTune.String() {
		_ = level.Error(logger).Log("msg", "文件用途不正确", "purpose", file.Purpose)
		return res, encode.InvalidParams.Wrap(errors.New("文件用途不正确"))
	}

	if request.DatasetType == types.EvalDataSetTypeCustom.String() && file.Purpose != types.FilePurposeFineTuneEval.String() {
		_ = level.Error(logger).Log("msg", "文件用途不正确", "purpose", file.Purpose)
		return res, encode.InvalidParams.Wrap(errors.New("文件用途不正确"))
	}

	req := types.LLMEvalResults{
		Status:        types.EvalStatusPending,
		ModelName:     request.ModelName,
		MetricName:    request.MetricName,
		DatasetFileId: file.ID,
		DatasetType:   request.DatasetType,
		Remark:        request.Remark,
		UUid:          uuid.NewString(),
	}
	if request.DatasetType == types.EvalDataSetTypeTrain.String() {
		req.EvalTotal = int(float64(file.LineCount) * request.EvalPercent)
	}
	err = s.store.Model().CreateEval(ctx, &req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "CreateEval", "err", err.Error())
		return
	}
	res = convertEval(&req)
	return
}

func (s *service) Undeploy(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Undeploy")
	m, err := s.store.Model().GetModel(ctx, id, "ModelDeploy")
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "id", id)
		return
	}
	//if !m.IsPrivate {
	//	// 公有模型不需要部署
	//	_ = level.Warn(logger).Log("msg", "public model not need undeploy", "modelName", m.ModelName)
	//	return encode.Invalid.Wrap(errors.Errorf("public model not need undeploy, model:%s", m.ModelName))
	//}
	//deploy, err := s.store.Model().FindModelDeployByModeId(ctx, m.ID)
	//if err != nil {
	//	_ = level.Warn(logger).Log("store.Model", "FindModelDeployByModeId", "err", err.Error())
	//	return err
	//}
	if err = s.store.Model().CancelModelDeploy(ctx, m.ID); err != nil {
		_ = level.Warn(logger).Log("repository.ModelDeploy", "CancelModelDeploy", "err", err.Error())
		return err
	}
	var deploymentName = fmt.Sprintf("%s-%d", util.ReplacerServiceName(m.ModelName), m.ID)
	if m.ModelDeploy.DeploymentName != "" {
		deploymentName = m.ModelDeploy.DeploymentName
	}
	// 调用API取消部署模型
	err = s.apiSvc.Runtime().RemoveDeployment(ctx, deploymentName)
	if err != nil {
		_ = level.Error(logger).Log("api.Runtime", "UndeployModel", "err", err.Error(), "modelName", m.ModelName)
	}
	_ = level.Info(logger).Log("msg", "undeploy model success")

	return nil
}

func (s *service) Deploy(ctx context.Context, request ModelDeployRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Deploy")
	m, err := s.store.Model().GetModel(ctx, request.Id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "request", fmt.Sprintf("%+v", request))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return encode.InvalidParams.Wrap(errors.New("模型不存在"))
		}
		return encode.ErrSystem.Wrap(errors.New("查询模型异常"))
	}
	if m.ProviderName != types.ModelProviderLocalAI {
		// 非本地模型不需要部署
		_ = level.Warn(logger).Log("msg", "non-local model not need deploy", "modelName", m.ModelName)
		return encode.Invalid.Wrap(errors.Errorf("non-local model not need deploy, model:%s", m.ModelName))
	}
	if m.BaseModelName != "" {
		// 别名模型不需要部署
		_ = level.Warn(logger).Log("msg", "BaseModelName model not need deploy", "modelName", m.ModelName)
		return encode.Invalid.Wrap(errors.Errorf("BaseModelName model not need deploy, model:%s", m.ModelName))
	}

	serviceName := util.ReplacerServiceName(m.ModelName)
	_ = level.Info(logger).Log("serviceName", serviceName)
	var (
		baseModel = m.ModelName
		modelPath = fmt.Sprintf("/data/base-model/%s", serviceName)
	)
	if m.IsFineTuning {
		finetunedModel, err := s.store.FineTuning().GetFineTuningJobByModelName(ctx, m.ModelName)
		if err != nil {
			err = errors.Wrap(err, "GetFineTuningJobByModelName failed")
			_ = level.Warn(logger).Log("msg", "GetFineTuningJobByModelName failed", "err", err)
			return err
		}
		baseModel = finetunedModel.BaseModel
		modelPath = finetunedModel.OutputDir
	}
	baseModelTemplate, err := s.store.FineTuning().FindFineTuningTemplateByModelType(ctx, baseModel, types.TemplateTypeInference)
	if err != nil {
		_ = level.Warn(logger).Log("repository.FineTuning", "FindFineTuningTemplateByModel", "err", err.Error())
		return err
	}
	if !m.IsFineTuning {
		modelPath = baseModelTemplate.BaseModelPath
	}
	minPort := 1024
	maxPort := 65535
	randomPort := rand.Intn(maxPort-minPort+1) + minPort

	var envs []runtime.Env
	var envVars []string
	envs = append(envs, runtime.Env{
		Name:  "MODEL_PATH",
		Value: modelPath,
	}, runtime.Env{
		Name:  "MODEL_NAME",
		Value: m.ModelName,
	}, runtime.Env{
		Name:  "HTTP_PORT",
		Value: strconv.Itoa(randomPort),
	}, runtime.Env{
		Name:  "QUANTIZATION",
		Value: request.Quantization,
	}, runtime.Env{
		Name:  "NUM_GPUS",
		Value: strconv.Itoa(request.Gpu),
	}, runtime.Env{
		Name:  "MAX_GPU_MEMORY",
		Value: strconv.Itoa(request.MaxGpuMemory),
	}, runtime.Env{
		Name:  "MODEL_WORKER_TYPE",
		Value: strings.ToLower(request.ModelWorker),
	}, runtime.Env{
		Name:  "INFERRED_TYPE",
		Value: request.InferredType,
	}, runtime.Env{
		Name:  "HF_HOME",
		Value: "/data/hf",
	}, runtime.Env{
		Name:  "CONTROLLER_ADDRESS",
		Value: s.options.controllerAddress,
	}, runtime.Env{
		Name:  "HF_ENDPOINT",
		Value: os.Getenv("HF_ENDPOINT"),
	}, runtime.Env{
		Name:  "HTTP_PROXY",
		Value: os.Getenv("HTTP_PROXY"),
	}, runtime.Env{
		Name:  "HTTPS_PROXY",
		Value: os.Getenv("HTTPS_PROXY"),
	}, runtime.Env{
		Name:  "NO_PROXY",
		Value: os.Getenv("NO_PROXY"),
	})
	for _, v := range envs {
		envVars = append(envVars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}

	runtimeConfig := runtime.Config{
		ServiceName: fmt.Sprintf("%s-%d", serviceName, m.ID),
		Image:       baseModelTemplate.TrainImage,
		Command: []string{
			"/bin/bash",
			"/app/start.sh",
		},
		EnvVars:            envVars,
		Volumes:            []runtime.Volume{{Key: s.options.volumeName, Value: "/data"}},
		GpuTolerationValue: request.Label,
		GPU:                request.Gpu,
		ConfigData: map[string]string{
			"/app/start.sh": baseModelTemplate.Content,
		},
		Replicas: int32(request.Replicas),
		Ports: map[string]string{
			strconv.Itoa(randomPort): strconv.Itoa(randomPort),
		},
	}

	deploymentName, err := s.apiSvc.Runtime().CreateDeployment(ctx, runtimeConfig)
	if err != nil {
		_ = level.Error(logger).Log("api.Runtime", "DeployModel", "err", err.Error(), "modelName", m.Model)
		return err
	}
	_ = level.Info(logger).Log("msg", "create deployment success", "deploymentName", deploymentName)
	// 更新模型部署状态
	if err = s.store.Model().SaveModelDeploy(ctx, &types.ModelDeploy{
		ModelID:        uint64(m.ID),
		ModelPath:      modelPath,
		Status:         types.ModelDeployStatusPending.String(),
		Replicas:       request.Replicas,
		InferredType:   request.InferredType,
		Label:          request.Label,
		Gpu:            request.Gpu,
		Cpu:            request.Cpu,
		Quantization:   request.Quantization,
		Vllm:           request.Vllm,
		DeploymentName: deploymentName,
	}); err != nil {
		_ = level.Warn(logger).Log("repository.ModelDeploy", "SaveModelDeploy", "err", err.Error())
		return err
	}
	return
}

func (s *service) ListModels(ctx context.Context, request ListModelRequest) (res ListModelResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListModels")
	req := model.ListModelRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
		Enabled:  request.Enabled,
		//IsPrivate:    request.IsPrivate,
		IsFineTuning: request.IsFineTuning,
		ModelName:    request.ModelName,
		ProviderName: request.ProviderName,
		ModelType:    request.ModelType,
	}
	models, total, err := s.store.Model().ListModels(ctx, req, "ModelDeploy")
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "ListModels", "err", err.Error())
		return res, encode.ErrSystem.Wrap(errors.New("查询模型列表失败"))
	}

	list := make([]Model, 0)
	for _, v := range models {
		var containers []runtime.Container
		var containerNames []string
		if v.ModelDeploy.ID > 0 {
			containers, err = s.apiSvc.Runtime().GetContainers(ctx, fmt.Sprintf("%s-%d", util.ReplacerServiceName(v.ModelName), v.ID))
			if err != nil {
				_ = level.Warn(logger).Log("api.Runtime", "GetDeploymentContainerNames", "err", err.Error())
				continue
			}
		}

		for _, container := range containers {
			containerNames = append(containerNames, container.Name)
		}

		c := convert(&v)
		c.ContainerNames = containerNames
		list = append(list, c)
	}
	res = ListModelResponse{
		Total:  total,
		Models: list,
	}
	return
}

func (s *service) CreateModel(ctx context.Context, request CreateModelRequest) (res Model, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateModel")
	modelPath := fmt.Sprintf("/data/base-model/%s", request.ModelName)
	if request.ModelPath != "" {
		modelPath = request.ModelPath
	}
	m, err := s.store.Model().GetModelByModelName(ctx, request.ModelName)
	if err == nil && m.ID > 0 {
		_ = level.Warn(logger).Log("store.Model", "GetModelByModelName", "err", "模型名称已存在", "modelName", request.ModelName)
		return res, encode.InvalidParams.Wrap(errors.Errorf("%s 模型已存在", request.ModelName))
	}
	if request.BaseModelName != "" {
		m, err = s.store.Model().GetModelByModelName(ctx, request.BaseModelName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = level.Warn(logger).Log("store.Model", "GetModelByModelName", "err", "别名模型不存在", "BaseModelName", request.BaseModelName)
			return res, encode.InvalidParams.Wrap(errors.Errorf("%s 别名模型不存在", request.BaseModelName))
		}
	}
	provider := request.ProviderName
	if provider == "" {
		provider = providerName(request.ModelName).String()
	}
	req := types.Models{
		ProviderName: types.ModelProvider(provider),
		ModelType:    types.ModelType(request.ModelType),
		ModelName:    request.ModelName,
		MaxTokens:    request.MaxTokens,
		//IsPrivate:    request.IsPrivate,
		IsFineTuning:  request.IsFineTuning,
		Enabled:       request.Enabled,
		Remark:        request.Remark,
		TenantId:      request.TenantId,
		Parameters:    request.Parameters,
		LastOperator:  request.Email,
		BaseModelName: request.BaseModelName,
		Replicas:      request.Replicas,     //并行/实例数量
		Label:         request.Label,        //调度标签
		K8sCluster:    request.K8sCluster,   //k8s集群
		InferredType:  request.InferredType, //推理类型cpu,gpu
		Gpu:           request.Gpu,          //GPU数
		Cpu:           request.Cpu,          //CPU核数
		Memory:        request.Memory,       //内存G
	}
	err = s.store.Model().CreateModel(ctx, &req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "CreateModel", "err", err.Error())
		return
	}

	if request.ProviderName == types.ModelProviderLocalAI.String() {
		inferenceShell, _ := s.options.dataFs.ReadFile("data/inference.sh")
		// 自动添加模版
		if err = s.store.Sys().SaveFineTuningTemplate(ctx, &types.FineTuningTemplate{
			TemplateType:  types.TemplateTypeInference,
			Content:       string(inferenceShell),
			TrainImage:    s.options.defaultTrainImage,
			Name:          fmt.Sprintf("%s-%s", request.ModelName, types.TemplateTypeInference),
			BaseModel:     request.ModelName,
			BaseModelPath: modelPath,
			ScriptFile:    "/app/start.sh",
			OutputDir:     "/data/ft-model",
			MaxTokens:     request.MaxTokens,
			Enabled:       true,
		}); err != nil {
			_ = level.Warn(logger).Log("store.FineTuning", "CreateFineTuningTemplate", "err", err.Error())
		}
		trainShell, _ := s.options.dataFs.ReadFile("data/train.sh")
		// 自动添加模版
		if err = s.store.Sys().SaveFineTuningTemplate(ctx, &types.FineTuningTemplate{
			TemplateType:  types.TemplateTypeTrain,
			Content:       string(trainShell),
			TrainImage:    s.options.defaultTrainImage,
			Name:          fmt.Sprintf("%s-%s", request.ModelName, types.TemplateTypeTrain),
			BaseModel:     request.ModelName,
			BaseModelPath: modelPath,
			ScriptFile:    "/app/train.sh",
			OutputDir:     "/data/ft-model",
			MaxTokens:     request.MaxTokens,
			Enabled:       true,
		}); err != nil {
			_ = level.Warn(logger).Log("store.FineTuning", "CreateFineTuningTemplate", "err", err.Error())
		}
	}
	res = convert(&req)
	return
}

func (s *service) GetModel(ctx context.Context, id uint) (res Model, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetModel")
	m, err := s.store.Model().GetModel(ctx, id, "FineTuningTrainJob", "ModelDeploy")
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "GetModel", "err", err.Error(), "id", id)
		return
	}
	if m.BaseModelName != "" {
		m, err = s.store.Model().GetModelByModelName(ctx, m.BaseModelName)
		if err != nil {
			_ = level.Error(logger).Log("store.Model", "GetModelByModelName", "err", err.Error(), "BaseModelName", m.BaseModelName)
			return res, errors.New("别名模型不存在")
		}
	}
	res = convert(&m)
	var containers []runtime.Container
	var containerNames []string
	if m.ModelDeploy.ModelID > 0 {
		containers, err = s.apiSvc.Runtime().GetContainers(ctx, fmt.Sprintf("%s-%d", util.ReplacerServiceName(m.ModelName), m.ID))
		if err != nil {
			_ = level.Warn(logger).Log("api.Runtime", "GetDeploymentContainerNames", "err", err.Error())
		}
		res.Deployment = modelDeploymentResult{
			VLLM:         m.ModelDeploy.Vllm,
			Status:       m.ModelDeploy.Status,
			ModelPath:    m.ModelDeploy.ModelPath,
			Replicas:     m.ModelDeploy.Replicas,
			InferredType: m.ModelDeploy.InferredType,
			GPU:          m.ModelDeploy.Gpu,
			Quantization: m.ModelDeploy.Quantization,
		}
	}

	//if m.IsFineTuning {
	//	res.FineTuned = &modelFineTuneResult{
	//		JobId:   m.FineTuningTrainJob.JobId,
	//		FileId:  m.FineTuningTrainJob.FileId,
	//		FileUrl: m.FineTuningTrainJob.FileUrl,
	//	}
	//	if fileBody, err := util.GetHttpFileBody(m.FineTuningTrainJob.FileUrl); err == nil {
	//		// 取body第一行数据解析成json
	//		res.FineTuned.SystemPrompt = tokenizers.GetFirstLineSystemPrompt(fileBody)
	//	} else {
	//		_ = level.Warn(logger).Log("util.GetHttpFileBody", "err", err.Error())
	//	}
	//}

	for _, v := range containers {
		containerNames = append(containerNames, v.Name)
	}

	res.ContainerNames = containerNames
	return res, nil
}

func (s *service) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateModel")
	var req model.UpdateModelRequest
	req.Id = request.Id
	if request.TenantId != nil {
		req.TenantId = request.TenantId
	}
	if request.MaxTokens != nil {
		req.MaxTokens = request.MaxTokens
	}
	if request.Enabled != nil {
		req.Enabled = request.Enabled
	}
	if request.Remark != nil {
		req.Remark = request.Remark
	}
	req.BaseModelName = request.BaseModelName
	req.Replicas = request.Replicas         //并行/实例数量
	req.Label = request.Label               //调度标签
	req.K8sCluster = request.K8sCluster     //k8s集群
	req.InferredType = request.InferredType //推理类型cpu,gpu
	req.Gpu = request.Gpu                   //GPU数
	req.Cpu = request.Cpu                   //CPU核数
	req.Memory = request.Memory             //内存G

	err = s.store.Model().UpdateModel(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "UpdateModel", "err", err.Error())
		return
	}
	return
}

func (s *service) DeleteModel(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteModel")
	// todo 删除之前判断是否有绑定的渠道
	err = s.store.Model().DeleteModel(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Model", "DeleteModel", "err", err.Error())
		return
	}
	return
}

func convert(data *types.Models) Model {
	m := Model{
		Id:           data.ID,
		ProviderName: string(data.ProviderName),
		ModelType:    string(data.ModelType),
		ModelName:    data.ModelName,
		MaxTokens:    data.MaxTokens,
		//IsPrivate:    data.IsPrivate,
		IsFineTuning:  data.IsFineTuning,
		Enabled:       data.Enabled,
		Remark:        data.Remark,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
		DeployStatus:  data.ModelDeploy.Status,
		Parameters:    data.Parameters,
		LastOperator:  data.LastOperator,
		BaseModelName: data.BaseModelName,
		Replicas:      data.Replicas,
		Label:         data.Label,
		K8sCluster:    data.K8sCluster,
		InferredType:  data.InferredType,
		Gpu:           data.Gpu,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ServiceName:   fmt.Sprintf("%s-%d", util.ReplacerServiceName(data.ModelName), data.ID),
	}
	tenants := make([]Tenant, 0)
	for _, t := range data.Tenants {
		tenants = append(tenants, Tenant{
			Id:   t.ID,
			Name: t.Name,
		})
	}

	m.Tenants = tenants
	operation := make([]string, 0)
	operation = append(operation, "edit")
	if data.CanDelete() {
		operation = append(operation, "delete")
	}
	if data.CanDeploy() {
		operation = append(operation, "deploy")
	}
	if data.CanUndeploy() {
		operation = append(operation, "undeploy")
	}
	m.Operation = operation
	if data.FineTuningTrainJob.ID > 0 {
		m.JobId = data.FineTuningTrainJob.JobId
	}
	return m
}

func providerName(m string) types.ModelProvider {
	openAIModels := []string{
		"gpt-4-turbo-preview",
		"gpt-4-0125-preview",
		"gpt-4-1106-preview",
		"gpt-4-vision-preview",
		"gpt-4",
		"gpt-4-32k",
		"gpt-4-0613",
		"gpt-4-32k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-3.5-turbo-1106",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
		"gpt-3.5-turbo-instruct",
		"gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-3.5-turbo-0301",
		"text-davinci-003",
		"text-davinci-002",
		"code-davinci-002",
	}
	for _, v := range openAIModels {
		if strings.EqualFold(v, m) {
			return types.ModelProviderOpenAI
		}
	}
	return types.ModelProviderLocalAI
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc services.Service, opts ...CreationOption) Service {
	logger = log.With(logger, "service", "models")
	options := &CreationOptions{
		volumeName:        "aigc-data",
		defaultTrainImage: "dudulu/llmops:latest",
		convertUrlFun: func(fileUrl string) string {
			return fileUrl
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &service{
		logger:  log.With(logger, "service", "models"),
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
		options: options,
	}
}

func convertEval(data *types.LLMEvalResults) Eval {
	e := Eval{
		Id:          data.ID,
		ModelName:   data.ModelName,
		DatasetType: data.DatasetType,
		Progress:    data.Progress,
		Score:       data.Score,
		CreatedAt:   data.CreatedAt,
		Status:      data.Status.String(),
		EvalTotal:   data.EvalTotal,
		Remark:      data.Remark,
		MetricName:  data.MetricName,
	}
	if data.Status == types.EvalStatusRunning && data.StartedAt != nil {
		start := data.StartedAt.Time
		e.Duration = util.FormatDuration(float64(time.Now().Sub(start)), util.PrecisionMinutes)
	}
	if data.Status == types.EvalStatusSuccess && data.CompletedAt != nil && data.StartedAt != nil {
		start := data.StartedAt.Time
		end := data.CompletedAt.Time
		e.Duration = util.FormatDuration(float64(end.Sub(start)), util.PrecisionMinutes)
	}
	if data.StartedAt != nil {
		e.StartedAt = data.StartedAt.Time.Format(time.RFC3339)
	}
	return e
}
