package repository

import (
	"github.com/IceBearAI/aigc/src/pkg/tenant"
	"github.com/IceBearAI/aigc/src/repository/assistants"
	"github.com/IceBearAI/aigc/src/repository/auth"
	"github.com/IceBearAI/aigc/src/repository/channel"
	"github.com/IceBearAI/aigc/src/repository/datasets"
	"github.com/IceBearAI/aigc/src/repository/datasettask"
	"github.com/IceBearAI/aigc/src/repository/files"
	"github.com/IceBearAI/aigc/src/repository/finetuning"
	"github.com/IceBearAI/aigc/src/repository/llmeval"
	"github.com/IceBearAI/aigc/src/repository/model"
	"github.com/IceBearAI/aigc/src/repository/modelevaluate"
	"github.com/IceBearAI/aigc/src/repository/sys"
	"github.com/IceBearAI/aigc/src/repository/tools"
	"gorm.io/gorm"

	"github.com/IceBearAI/aigc/src/repository/chat"
	"github.com/go-kit/log"
	"github.com/opentracing/opentracing-go"
)

type Repository interface {
	// Chat 模块
	Chat() chat.Service

	// Channel 模块
	Channel() channel.Service
	// Model 模块
	Model() model.Service
	// Auth 模块
	Auth() auth.Service
	// Files 模块
	Files() files.Service
	// FineTuning 微调模块
	FineTuning() finetuning.Service
	// Sys 模块
	Sys() sys.Service
	// Dataset 模块
	Dataset() datasets.Service
	// Tools 工具
	Tools() tools.Service
	// Assistants 助手
	Assistants() assistants.Service
	// LLMEval 评估模块
	LLMEval() llmeval.Service
	// Tenants 租户模块
	Tenants() tenant.Service
	// DatasetTask 数据集任务模块
	DatasetTask() datasettask.Service
	// ModelEvaluate 模型评测
	ModelEvaluate() modelevaluate.Service
}

type repository struct {
	chatSvc          chat.Service
	channelSvc       channel.Service
	modelSvc         model.Service
	authSvc          auth.Service
	filesSvc         files.Service
	fineTuningSvc    finetuning.Service
	sysSvc           sys.Service
	datasetSvc       datasets.Service
	toolsSvc         tools.Service
	assistantsSvc    assistants.Service
	llmEvalSvc       llmeval.Service
	tenantSvc        tenant.Service
	datasetTaskSvc   datasettask.Service
	modelEvaluateSvc modelevaluate.Service
}

func (r *repository) DatasetTask() datasettask.Service {
	return r.datasetTaskSvc
}

func (r *repository) Tenants() tenant.Service {
	return r.tenantSvc
}

func (r *repository) LLMEval() llmeval.Service {
	return r.llmEvalSvc
}

func (r *repository) Assistants() assistants.Service {
	return r.assistantsSvc
}

func (r *repository) Tools() tools.Service {
	return r.toolsSvc
}

func (r *repository) Dataset() datasets.Service {
	return r.datasetSvc
}

func (r *repository) Chat() chat.Service {
	return r.chatSvc
}

func (r *repository) Channel() channel.Service {
	return r.channelSvc
}

func (r *repository) Model() model.Service {
	return r.modelSvc
}

func (r *repository) Auth() auth.Service {
	return r.authSvc
}

func (r *repository) Files() files.Service {
	return r.filesSvc
}

func (r *repository) FineTuning() finetuning.Service {
	return r.fineTuningSvc
}

func (r *repository) Sys() sys.Service {
	return r.sysSvc
}

func (r *repository) ModelEvaluate() modelevaluate.Service {
	return r.modelEvaluateSvc
}

var _ Repository = (*repository)(nil)

func New(db *gorm.DB, logger log.Logger, traceId string, tracer opentracing.Tracer) Repository {
	chatSvc := chat.New(db)
	channelSvc := channel.New(db)
	modelSvc := model.New(db)
	authSvc := auth.New(db)
	filesSvc := files.New(db)
	fineTuningSvc := finetuning.New(db)
	sysSvc := sys.New(db)
	datasetSvc := datasets.New(db)
	toolsSvc := tools.New(db)
	assistantsSvc := assistants.New(db)
	llmEvalSvc := llmeval.New(db)
	tenantSvc := tenant.New(db)
	datasetTaskSvc := datasettask.New(db)
	modelEvaluateSvc := modelevaluate.New(db)

	if logger != nil {
		chatSvc = chat.NewLogging(logger, traceId)(chatSvc)
		authSvc = auth.NewLogging(logger, traceId)(authSvc)
		filesSvc = files.NewLogging(logger, traceId)(filesSvc)
		modelSvc = model.NewLogging(logger, traceId)(modelSvc)
		fineTuningSvc = finetuning.NewLogging(logger, traceId)(fineTuningSvc)
		sysSvc = sys.NewLogging(logger, traceId)(sysSvc)
		datasetSvc = datasets.NewLogging(logger, traceId)(datasetSvc)
		toolsSvc = tools.NewLogging(logger, traceId)(toolsSvc)
		assistantsSvc = assistants.NewLogging(logger, traceId)(assistantsSvc)
		llmEvalSvc = llmeval.NewLogging(logger, traceId)(llmEvalSvc)
		tenantSvc = tenant.NewLogging(logger, traceId)(tenantSvc)
		datasetTaskSvc = datasettask.NewLogging(logger, traceId)(datasetTaskSvc)
		modelEvaluateSvc = modelevaluate.NewLogging(logger, traceId)(modelEvaluateSvc)
	}

	if tracer != nil {
		chatSvc = chat.NewTracing(tracer)(chatSvc)
		authSvc = auth.NewTracing(tracer)(authSvc)
		filesSvc = files.NewTracing(tracer)(filesSvc)
		modelSvc = model.NewTracing(tracer)(modelSvc)
		fineTuningSvc = finetuning.NewTracing(tracer)(fineTuningSvc)
		sysSvc = sys.NewTracing(tracer)(sysSvc)
		datasetSvc = datasets.NewTracing(tracer)(datasetSvc)
		toolsSvc = tools.NewTracing(tracer)(toolsSvc)
		assistantsSvc = assistants.NewTracing(tracer)(assistantsSvc)
		llmEvalSvc = llmeval.NewTracing(tracer)(llmEvalSvc)
		tenantSvc = tenant.NewTracing(tracer)(tenantSvc)
		datasetTaskSvc = datasettask.NewTracing(tracer)(datasetTaskSvc)
		modelEvaluateSvc = modelevaluate.NewTracing(tracer)(modelEvaluateSvc)
	}

	return &repository{
		chatSvc:          chatSvc,
		channelSvc:       channelSvc,
		modelSvc:         modelSvc,
		authSvc:          authSvc,
		filesSvc:         filesSvc,
		fineTuningSvc:    fineTuningSvc,
		sysSvc:           sysSvc,
		datasetSvc:       datasetSvc,
		toolsSvc:         toolsSvc,
		assistantsSvc:    assistantsSvc,
		llmEvalSvc:       llmEvalSvc,
		tenantSvc:        tenantSvc,
		datasetTaskSvc:   datasetTaskSvc,
		modelEvaluateSvc: modelEvaluateSvc,
	}
}
