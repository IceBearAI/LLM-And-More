package modelevaluate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/runtime"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts     []kithttp.ClientOption
	callbackHost       string
	gpuTolerationValue string
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithCallbackHost returns a CreationOption that sets the callback host.
func WithCallbackHost(host string) CreationOption {
	return func(co *CreationOptions) {
		co.callbackHost = host
	}
}

// WithDatasetGpuTolerationValue returns a CreationOption  that sets the dataset drive.
func WithDatasetGpuTolerationValue(gpuTolerationValue string) CreationOption {
	return func(co *CreationOptions) {
		co.gpuTolerationValue = gpuTolerationValue
	}
}

type Middleware func(Service) Service

type Service interface {
	// List 模型评估列表
	List(ctx context.Context, req listRequest) (res []listResult, total int64, err error)
	// Create 模型评估创建
	Create(ctx context.Context, req createRequest) (err error)
	// Cancel 模型评估取消
	Cancel(ctx context.Context, req cancelRequest) (err error)
	// Delete 模型评估取消
	Delete(ctx context.Context, req deleteRequest) (err error)
	// FiveGraph 五维图
	FiveGraph(ctx context.Context, req fiveGraphRequest) (res1, res2, res3 fiveGraphResult, err error)
	// EvalFinish 评测结果完成
	EvalFinish(ctx context.Context, req finishRequest) (err error)
}

type service struct {
	options    *CreationOptions
	traceId    string
	logger     log.Logger
	repository repository.Repository
	apiSvc     services.Service
	filesSvc   files.Service
}

func (s *service) List(ctx context.Context, req listRequest) (res []listResult, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "List")
	list, total, err := s.repository.ModelEvaluate().ListModelEvaluate(ctx, req.Page, req.PageSize, uint(req.ModelId), req.Status, req.EvalTargetType)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "ListModelEvaluate", "err", err.Error())
		return
	}

	for _, v := range list {
		dd := listResult{
			Id:             int(v.ID),
			Uuid:           v.Uuid,
			ModelID:        v.ModelId,
			Status:         v.Status,
			EvalTargetType: v.EvalTargetType,
			Score:          v.Score,
			DataType:       v.DataType,
			DataSize:       v.DataSize,
			Complete:       fmt.Sprintf("%.0f%%", v.CompleteRate*100),
			Remark:         v.Remark,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		}

		res = append(res, dd)
	}

	return
}

func (s *service) Create(ctx context.Context, req createRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Create")

	//如果有同类型的评测，不能重复提交
	count, err := s.repository.ModelEvaluate().CountEvaluate(ctx, req.EvalTargetType, []string{string(types.EvaluateStatusWaiting), string(types.EvaluateStatusRunning)})
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "CountEvaluate", "err", err.Error())
		return
	}

	if count > 0 {
		return errors.New("有未结束的评测任务，请稍后再试！")
	}

	var dataSize int
	var fileId, s3Url string
	if req.EvalTargetType != string(types.EvaluateTargetTypeFive) {
		//获取文件s3地址
		fileInfo, err := s.repository.Files().FindFileByFileId(ctx, req.FileId)
		if err != nil {
			_ = level.Error(logger).Log("s.repository.Files", "FindFileByFileId", "err", err.Error())
			return err
		}

		dataSize = fileInfo.LineCount
		fileId = fileInfo.FileID
		s3Url = fileInfo.S3Url
	}

	modelInfo, err := s.repository.Model().GetModel(ctx, uint(req.ModelId))
	if err != nil {
		_ = level.Error(logger).Log("s.repository.Model", "GetModel", "err", err.Error())
		return
	}

	serviceName := strings.ReplaceAll(modelInfo.ModelName, ":", "-")
	serviceName = strings.ReplaceAll(modelInfo.ModelName, "::", "-")
	serviceName = strings.ReplaceAll(modelInfo.ModelName, ".", "-")

	var (
		baseModel = modelInfo.ModelName
		modelPath = fmt.Sprintf("/data/base-model/%s", serviceName)
	)
	if modelInfo.IsFineTuning {
		finetunedModel, err := s.repository.FineTuning().GetFineTuningJobByModelName(ctx, modelInfo.ModelName)
		if err != nil {
			err = errors.Wrap(err, "GetFineTuningJobByModelName failed")
			_ = level.Warn(logger).Log("msg", "GetFineTuningJobByModelName failed", "err", err)
			return err
		}
		baseModel = finetunedModel.BaseModel
		modelPath = finetunedModel.OutputDir
	}
	baseModelTemplate, err := s.repository.FineTuning().FindFineTuningTemplateByModel(ctx, baseModel)
	if err != nil {
		_ = level.Warn(logger).Log("repository.FineTuning", "FindFineTuningTemplateByModel", "err", err.Error())
		return err
	}

	operatorEmail, _ := middleware.GetEmail(ctx)
	data := &types.ModelEvaluate{
		ModelId:        req.ModelId,
		ModelPath:      modelPath,
		Status:         string(types.EvaluateStatusWaiting),
		Replicas:       req.Replicas,
		InferredType:   req.InferredType,
		Label:          req.Label,
		Gpu:            req.Gpu,
		Cpu:            req.Cpu,
		MaxGpuMemory:   req.maxGpuMemory,
		Memory:         req.Memory,
		EvalTargetType: req.EvalTargetType,
		MaxLength:      req.MaxLength,
		BatchSize:      req.BatchSize,
		DataSize:       dataSize,
		DataType:       string(types.EvaluateDataTypeCustom),
		DataFileId:     fileId,
		DataUrl:        s3Url,
		Remark:         req.Remark,
		CompleteRate:   0,
		OperatorEmail:  operatorEmail,
		Uuid:           uuid.New().String(),
	}

	err = s.repository.ModelEvaluate().Save(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "Save", "err", err.Error())
		return
	}

	tenantUUid, _ := ctx.Value(middleware.ContextKeyPublicTenantId).(string)
	auth, _ := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)

	// 默认指标脚本
	shellName := "model_performance_evaluation.sh"
	if req.EvalTargetType == string(types.EvaluateTargetTypeFive) {
		// 五维图指标脚本
		shellName = "evaluate_model_from_five_dimensions.sh"
	}

	var envs []runtime.Env
	var envVars []string
	envs = append(envs, runtime.Env{
		Name:  "MODEL_PATH",
		Value: modelPath,
	}, runtime.Env{
		Name:  "JOB_ID",
		Value: data.Uuid, // JOB uuid
	}, runtime.Env{
		Name:  "TENANT_ID",
		Value: tenantUUid,
	}, runtime.Env{
		Name:  "API_URL",
		Value: fmt.Sprintf("%s/api/evaluate/finish/%s", s.options.callbackHost, data.Uuid), // 回调
	}, runtime.Env{
		Name:  "AUTH",
		Value: auth,
	}, runtime.Env{
		Name:  "HF_HOME",
		Value: "/data/hf",
	}, runtime.Env{
		Name:  "HF_ENDPOINT",
		Value: os.Getenv("HF_ENDPOINT"),
	}, runtime.Env{
		Name:  "HTTP_PROXY",
		Value: os.Getenv("HTTP_PROXY"),
	}, runtime.Env{
		Name:  "HTTPS_PROXY",
		Value: os.Getenv("HTTPS_PROXY"),
	})
	// 如果是其他指标，则传以下参数
	envs = append(envs, runtime.Env{
		Name:  "DATASET_PATH", // 数据集路径如果是url则传url如果是文件则传文件路径
		Value: data.DataUrl,
	}, runtime.Env{
		Name:  "EVALUATION_METRICS", // 评估的指标
		Value: data.EvalTargetType,
	}, runtime.Env{
		Name:  "MAX_SEQ_LEN", // 最大长度
		Value: strconv.Itoa(data.MaxLength),
	}, runtime.Env{
		Name:  "PER_DEVICE_BATCH_SIZE", // 评估批次
		Value: strconv.Itoa(data.BatchSize),
	})

	for _, v := range envs {
		envVars = append(envVars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}

	config := runtime.Config{
		ServiceName: fmt.Sprintf("%s-%s", serviceName, strings.ToLower(req.EvalTargetType)),
		Image:       baseModelTemplate.TrainImage,
		Replicas:    int32(req.Replicas),
		GPU:         req.Gpu,
		Command: []string{
			"/bin/bash",
			"/app/eval/" + shellName,
		},
		EnvVars: envVars,
	}

	if req.Label != "" {
		config.GpuTolerationValue = req.Label
	} else {
		config.GpuTolerationValue = s.options.gpuTolerationValue
	}

	jobName, err := s.apiSvc.Runtime().CreateJob(ctx, config)

	status := string(types.EvaluateStatusRunning)
	statusMsg := ""
	if err != nil {
		status = string(types.EvaluateStatusFailed)
		statusMsg = err.Error()
	}
	// 更新数据库状态
	data.JobName = jobName
	data.Status = status
	data.StatusMsg = statusMsg

	err = s.repository.ModelEvaluate().Save(ctx, data)

	return
}

func (s *service) Cancel(ctx context.Context, req cancelRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Delete")
	info, err := s.repository.ModelEvaluate().GetByUuid(ctx, req.Uuid)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "GetById", "err", err.Error())
		return
	}
	if info.ID == 0 {
		return fmt.Errorf("该条信息不存在")
	}

	//取消job
	if err = s.apiSvc.Runtime().RemoveJob(ctx, info.JobName); err != nil {
		_ = level.Error(logger).Log("s.apiSvc.Runtime", "RemoveJob", "err", err.Error())
		return
	}

	// 更新db
	info.Status = string(types.EvaluateStatusCancel)
	info.OperatorEmail, _ = middleware.GetEmail(ctx)
	err = s.repository.ModelEvaluate().Save(ctx, &info)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "Save", "err", err.Error())
		return
	}
	return
}

func (s *service) Delete(ctx context.Context, req deleteRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Delete")
	info, err := s.repository.ModelEvaluate().GetByUuid(ctx, req.Uuid)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "GetById", "err", err.Error())
		return
	}
	if err = s.apiSvc.Runtime().RemoveJob(ctx, info.JobName); err != nil {
		_ = level.Warn(logger).Log("s.apiSvc.Runtime", "RemoveJob", "err", err.Error())
	}
	if info.ID == 0 {
		return fmt.Errorf("该条信息不存在")
	}
	err = s.repository.ModelEvaluate().DeleteById(ctx, info.ID)
	return
}

func (s *service) FiveGraph(ctx context.Context, req fiveGraphRequest) (res1, res2, res3 fiveGraphResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Delete")

	res1, err = s.getFiveGraphData(ctx, req.CurrentModelId, req.CurrentModelEvaluateId)
	if err != nil {
		_ = level.Error(logger).Log("s.getFiveGraphData CurrentModelId modelId", req.CurrentModelId, "err", err.Error())
	}
	res2, err = s.getFiveGraphData(ctx, req.Compare1ModelId, 0)
	if err != nil {
		_ = level.Error(logger).Log("s.getFiveGraphData Compare1ModelId modelId", req.Compare1ModelId, "err", err.Error())
	}
	res3, err = s.getFiveGraphData(ctx, req.Compare2ModelId, 0)
	if err != nil {
		_ = level.Error(logger).Log("s.getFiveGraphData Compare2ModelId modelId", req.Compare2ModelId, "err", err.Error())
	}

	return res1, res2, res3, nil
}

type logEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Loss         float64   `json:"loss"`
	LearningRate float64   `json:"learning_rate"`
	Epoch        float64   `json:"epoch"`
}

func (s *service) getFiveGraphData(ctx context.Context, modelId, evaluateId int) (res fiveGraphResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "getFiveGraphData")

	data, err := s.repository.ModelEvaluate().FindFiveGraphLastByModelId(ctx, uint(modelId), uint(evaluateId))
	if err != nil {
		return res, err
	}

	var lastLoss float64

	// 如果是微调模型，获取损失率
	if data.Models.IsFineTuning == true {
		fineTuningJob, err := s.repository.FineTuning().GetFineTuningJobByModelName(ctx, data.Models.ModelName)

		if err == nil && fineTuningJob.TrainLog != "" {
			lineArr := strings.Split(fineTuningJob.TrainLog, "\n")
			re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z) (\{.*?\})`)

			var logEntryList []logEntry

			for _, log := range lineArr {
				matches := re.FindStringSubmatch(log)
				if len(matches) == 3 {
					timestampStr, jsonStr := matches[1], matches[2]

					timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
					if err != nil {
						_ = level.Warn(logger).Log("msg", "parse timestamp failed", "err", err.Error())
						continue
					}

					jsonStr = strings.Replace(jsonStr, "'", "\"", -1)        // 将单引号替换为双引号
					jsonStr = strings.Replace(jsonStr, "False", "false", -1) // 将 False 替换为 false
					jsonStr = strings.Replace(jsonStr, "True", "true", -1)   // 将 True 替换为 true

					var entry logEntry
					err = json.Unmarshal([]byte(jsonStr), &entry)
					if err != nil {
						_ = level.Warn(logger).Log("msg", "unmarshal json failed", "err", err.Error())
						continue
					}

					entry.Timestamp = timestamp
					logEntryList = append(logEntryList, entry)
				}
			}
			lastLine := logEntryList[len(logEntryList)-1]
			lastLoss = lastLine.Loss
		}

	}

	var riskOver, riskUnder, riskDisaster bool
	// 过拟合 <0.02 过拟合  建议降低训练轮次/ >0.5  欠拟合  建议提高训练轮次
	if lastLoss < 0.02 {
		riskOver = true
	} else if lastLoss > 0.5 {
		riskUnder = true
	}

	_ = level.Info(logger).Log("modelEvaluate", "getFiveGraphData", "IsFineTuning", data.Models.IsFineTuning, "lastLoss", lastLoss)

	var remind string
	remind = "建议"
	if data.RiskOver == true {
		remind += "降低训练轮次，"
	}
	if data.RiskUnder == true {
		remind += "提高训练轮次，"
	}
	if data.RiskDisaster == true {

	}
	remind += "并重新训练"

	score, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", data.Five1+data.Five2+data.Five3+data.Five4+data.Five5), 64)

	res = fiveGraphResult{
		RiskOver:     riskOver,
		RiskUnder:    riskUnder,
		RiskDisaster: riskDisaster,
		Remind:       remind,
		Value:        []float64{data.Five1, data.Five2, data.Five3, data.Five4, data.Five5},
		Score:        score,
		ModelId:      data.ModelId,
		Name:         data.Models.ModelName,
		IsFineTuning: data.Models.IsFineTuning,
	}

	return res, nil
}

func (s *service) EvalFinish(ctx context.Context, req finishRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "EvalFinish")
	info, err := s.repository.ModelEvaluate().GetByUuid(ctx, req.JobId)
	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "GetById", "err", err.Error())
		return
	}

	_ = level.Info(logger).Log("modelEvaluate", "EvalFinish", "uuid", req.JobId, "result", req.Result)

	// 五维图
	if info.EvalTargetType == string(types.EvaluateTargetTypeFive) {
		var data0 fiveResult
		// 解码JSON到map
		err = json.Unmarshal([]byte(req.Result), &data0)
		if err != nil {
			_ = level.Error(logger).Log("EvalFinish", "json.Unmarshal", "err", err.Error())
		}

		//InferenceAbility  推理能力
		//ReadingComprehension 阅读理解能力
		//ChineseLanguageSkill 中文能力
		//CommandCompliance 指令遵从能力
		//InnovationCapacity 创新能力
		//five1 中文能力
		//five2 推理能力
		//five3 指令遵从能力
		//five4 创新能力
		//five5 阅读理解

		if data0.Status == "success" {
			info.Status = string(types.EvaluateStatusSuccess)
			info.Five1 = data0.Data.ChineseLanguageSkill
			info.Five2 = data0.Data.InferenceAbility
			info.Five3 = data0.Data.CommandCompliance
			info.Five4 = data0.Data.InnovationCapacity
			info.Five5 = data0.Data.ReadingComprehension
			info.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", info.Five1+info.Five2+info.Five3+info.Five4+info.Five5), 64)
		} else {
			info.Status = string(types.EvaluateStatusFailed)
		}

	}

	// 非五维图
	if info.EvalTargetType != string(types.EvaluateTargetTypeFive) {
		var data callbackResult
		// 解码JSON到map
		err = json.Unmarshal([]byte(req.Result), &data)
		if err != nil {
			_ = level.Error(logger).Log("EvalFinish", "json.Unmarshal", "err", err.Error())
		}

		if data.Status == "success" {
			info.Status = string(types.EvaluateStatusSuccess)
			// 提取得分结果
			for key, value := range data.Data.OverallEvaluationMetrics {
				_ = level.Info(logger).Log("modelEvaluate", "EvalFinish", "Key", key, "value", value)
				if key != string(types.EvaluateTargetTypeFive) {
					if floatValue, ok := value.(float64); ok {
						info.Score = floatValue
					}
				}
			}
		} else {
			info.Status = string(types.EvaluateStatusFailed)
		}
	}

	info.Result = req.Result

	err = s.repository.ModelEvaluate().Save(ctx, &info)

	if err != nil {
		_ = level.Error(logger).Log("s.repository.ModelEvaluate", "Save", "err", err.Error())
		return
	}

	//移除job
	if err = s.apiSvc.Runtime().RemoveJob(ctx, info.JobName); err != nil {
		_ = level.Error(logger).Log("s.apiSvc.Runtime", "RemoveJob", "err", err.Error())
		return
	}

	return
}

var _ Service = &service{}

func New(logger log.Logger, traceId string, repository repository.Repository, apiSvc services.Service, filesSvc files.Service, opts ...CreationOption) Service {
	logger = log.With(logger, "evaluate", "service")
	options := &CreationOptions{
		callbackHost: "http://localhost:8080",
	}
	for _, opt := range opts {
		opt(options)
	}
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		apiSvc:     apiSvc,
		filesSvc:   filesSvc,
		options:    options,
	}
}
