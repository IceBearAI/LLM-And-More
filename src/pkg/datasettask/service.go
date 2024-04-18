package datasettask

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/repository"
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
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Middleware func(Service) Service

// Service is the interface that provides datasettask methods.
type Service interface {
	// CreateTask creates a new task.
	CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error)
	// ListTasks returns all tasks.
	ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error)
	// DeleteTask deletes a task.
	DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error)
	// GetTaskSegmentNext 获取一条待标注任务样本
	GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error)
	// AnnotationTaskSegment 标注一条任务样本
	AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error)
	// AbandonTaskSegment 放弃一条标注任务样本
	AbandonTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string) (err error)
	// AsyncCheckTaskDatasetSimilar 同步检查标注任务的数据集相似
	AsyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error)
	// GetCheckTaskDatasetSimilarLog 获取检测日志
	GetCheckTaskDatasetSimilarLog(ctx context.Context, tenantId uint, taskId string) (res string, err error)
	// CancelCheckTaskDatasetSimilar 取消检测
	CancelCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error)
	// SplitAnnotationDataSegment 将标注数据拆分成训练集和测试集
	SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error)
	// ExportAnnotationData 导出标注任务数据
	ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (filePath string, err error)
	// DeleteAnnotationTask 删除标注任务
	DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error)
	// CleanAnnotationTask 清理标注任务
	CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error)
	// TaskDetectFinish 任务检测完成
	TaskDetectFinish(ctx context.Context, tenantId uint, taskId, testReport string) (err error)
	// GetTaskInfo 获取任务详情
	GetTaskInfo(ctx context.Context, tenantId uint, taskId string) (res taskDetail, err error)
	// GenerationAnnotationContent 智能生成标注内容
	GenerationAnnotationContent(ctx context.Context, tenantId uint, modelName, taskId, taskSegmentId string) (res taskSegmentAnnotationRequest, err error)
	// GetTaskFAQIntents 获取所有意图标签
	GetTaskFAQIntents(ctx context.Context, tenantId uint, taskId string) (res []string, err error)
}

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts     []kithttp.ClientOption
	datasetImage       string
	datasetModel       string
	datasetDrive       string
	callbackHost       string
	gpuTolerationValue string
	volumeName         string
	fileSvc            files.Service
	convertUrlFun      func(fileUrl string) string
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithDatasetImage returns a CreationOption that sets the base url.
func WithDatasetImage(image string) CreationOption {
	return func(co *CreationOptions) {
		co.datasetImage = image
	}
}

// WithDatasetModel returns a CreationOption that sets the dataset model.
func WithDatasetModel(model string) CreationOption {
	return func(co *CreationOptions) {
		co.datasetModel = model
	}
}

// WithCallbackHost returns a CreationOption that sets the callback host.
func WithCallbackHost(host string) CreationOption {
	return func(co *CreationOptions) {
		co.callbackHost = host
	}
}

// WithDatasetDrive returns a CreationOption  that sets the dataset drive.
func WithDatasetDrive(drive string) CreationOption {
	return func(co *CreationOptions) {
		co.datasetDrive = drive
	}
}

// WithDatasetGpuTolerationValue returns a CreationOption  that sets the dataset drive.
func WithDatasetGpuTolerationValue(gpuTolerationValue string) CreationOption {
	return func(co *CreationOptions) {
		co.gpuTolerationValue = gpuTolerationValue
	}
}

// WithFileSvc returns a CreationOption that sets the file service.
func WithFileSvc(fileSvc files.Service) CreationOption {
	return func(co *CreationOptions) {
		co.fileSvc = fileSvc
	}
}

// WithVolumeName returns a CreationOption  that sets the dataset drive.
func WithVolumeName(volumeName string) CreationOption {
	return func(co *CreationOptions) {
		co.volumeName = volumeName
	}
}

// WithConvertUrl returns a CreationOption that sets the callback host.
func WithConvertUrl(convertFunc func(fileUrl string) string) CreationOption {
	return func(co *CreationOptions) {
		co.convertUrlFun = convertFunc
	}
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	apiSvc     services.Service
	options    *CreationOptions
	fileSvc    files.Service
}

func (s *service) GetTaskFAQIntents(ctx context.Context, tenantId uint, taskId string) (res []string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	taskInfo, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId, "DatasetDocument.DatasetDocumentSegments")
	if err != nil {
		_ = level.Warn(logger).Log("msg", "get task error", "err", err.Error())
		return
	}
	var documentSegmentIds []uint
	for _, v := range taskInfo.DatasetDocument.DatasetDocumentSegments {
		documentSegmentIds = append(documentSegmentIds, v.ID)
	}

	segments, err := s.repository.DatasetTask().GetSegmentFaqIntentInSegmentId(ctx, documentSegmentIds, types.DatasetAnnotationStatusCompleted, types.DatasetAnnotationTypeFAQ)
	if err != nil {
		_ = level.Warn(logger).Log("repository.DatasetTask", "GetSegmentFaqIntentInSegmentId", "err", err.Error())
		return res, nil
	}
	for _, v := range segments {
		res = append(res, v.Intent)
	}
	return
}

func (s *service) GenerationAnnotationContent(ctx context.Context, tenantId uint, modelName, taskId, taskSegmentId string) (res taskSegmentAnnotationRequest, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	taskInfo, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId, "Segments")
	if err != nil {
		_ = level.Warn(logger).Log("msg", "get task error", "err", err.Error())
		return res, nil
	}
	taskSegment, err := s.repository.DatasetTask().GetTaskSegmentByUUID(ctx, taskInfo.ID, taskSegmentId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.DatasetTask", "GetTaskSegmentByUUID", "err", err.Error())
		return res, nil
	}
	if strings.TrimSpace(taskSegment.Instruction) == "" {
		if prevSegment, err := s.repository.DatasetTask().GetTaskSegmentPrev(ctx, taskSegment.ID, types.DatasetAnnotationStatusCompleted); err == nil {
			taskSegment.Instruction = prevSegment.Instruction
		}
	}
	// 获取之前所有意图标签
	var intents []string
	var systemPrompt, userPrompt string
	if taskSegment.Instruction != "" {
		userPrompt = fmt.Sprintf("instruction: %s", taskSegment.Instruction)
	}
	switch types.DatasetAnnotationType(taskInfo.AnnotationType) {
	case types.DatasetAnnotationTypeGeneral:
		systemPrompt = fmt.Sprintf(`请仔细阅读下面的对话，并根据对话内容填充相应的字段。首先，识别客户提出的问题，并将其填充到"input"字段。然后，根据问题的内容，将回答填充到"output"字段。最后，将您的角色填充到"instruction"字段。请确保您的输出格式严格遵循以下JSON格式：

{
    "input": "<填写客户的问题>",
    "output": "<填写客服的回答>",
    "instruction": "<您的角色>"
}

示例：

输入：

instruction: %s  
document: 雇主责任险的医疗费用是不可重复报销的。如果您已经通过社保或其他商业保险公司报销了部分医疗费用，您可以提供分割单及理赔单据的复印件，以申请雇主责任险对剩余未报销部分进行索赔。这有助于确保您获得合理的赔付。
关于“企业无忧”雇主责任险的医疗费用，以下是您需要了解的信息：
1. 自费药费用包含：医疗费用中包含自费药费用。在保险单明细表所载明的每人每次医疗费用赔偿限额内，自费药费用将零免赔，百分百报销，无比例限制。
2. 医院级别要求：针对医院级别，并无特定等级限制要求。该保险涵盖不同级别的医院，包括一级医院（如社区医院）和二级医院（如县人民医院）等。保险将支付合理且必要的门诊和住院费用。

输出：

{
    "input": "雇主责任险的医疗费怎么报销？",
    "output": "雇主责任险的医疗费用是不可重复报销的，如果您已经通过社保或其他商业保险公司报销了部分医疗费用，您可以提供分割单及理赔单据的复印件，以申请雇主责任险对剩余未报销部分进行索赔。这有助于确保您获得合理的赔付。",
    "instruction": "%s"
}`, taskSegment.Instruction, taskSegment.Instruction)
	case types.DatasetAnnotationTypeFAQ:
		for _, segment := range taskInfo.Segments {
			if segment.Status != types.DatasetAnnotationStatusCompleted {
				continue
			}
			if util.StringInArray(intents, segment.Intent) {
				continue
			}
			if segment.Intent != "" {
				intents = append(intents, segment.Intent)
			}
		}
		systemPrompt = fmt.Sprintf(`请仔细阅读下面的对话，并根据对话内容填充相应的字段。首先，识别客户提出的问题，并将其填充到"question"字段。然后，根据问题的内容，从以下候选意图类别中选择一个最合适的意图，并填充到"intent"字段。如果提出的问题不符合任何候选类别，请基于您的理解提供一个合适的类别。
候选类别有：%s

接下来，将客服的回答填充到"output"字段。最后，将您的角色填充到"instruction"字段。请确保您的输出格式严格遵循以下JSON格式：

{
    "question": "<填写客户的问题>",
    "intent": "<填写问题的意图类别>",
    "output": "<填写客服的回答>",
    "instruction": "<您的角色>"
}

示例：

输入：

instruction: %s  
document: 客户: 怎么申请贷款？  
客服: 您好！我们这是基金，没有贷款业务  

输出：

{
    "question": "怎么申请贷款？",
    "intent": "咨询贷款申请",
    "output": "您好！我们这是基金，没有贷款业务",
    "instruction": "%s"
}`, strings.Join(intents, "、"), taskSegment.Instruction, taskSegment.Instruction)
	case types.DatasetAnnotationTypeRAG:
		systemPrompt = fmt.Sprintf(`请仔细阅读下面的对话，并根据对话内容填充相应的字段。首先，识别客户提出的问题，并将其填充到"input"字段。然后，根据问题的内容，将回答填充到"output"字段。最后，将您的角色填充到"instruction"字段。请确保您的输出格式严格遵循以下JSON格式：

{
    "question": "<填写客户的问题>",
    "document": "<检索出来的内容>",
    "output": "<填写客服的回答>",
    "instruction": "<您的角色>"
}

示例：

输入：

instruction: %s  
document: 雇主责任险的医疗费用是否可以重复报销？\n雇主责任险的医疗费用是不可重复报销的。如果您已经通过社保或其他商业保险公司报销了部分医疗费用，您可以提供分割单及理赔单据的复印件，以申请雇主责任险对剩余未报销部分进行索赔。这有助于确保您获得合理的赔付。

输出：

{
    "question": "雇主责任险的医疗费怎么报销？",
    "document": "雇主责任险的医疗费用是否可以重复报销？\n雇主责任险的医疗费用是不可重复报销的。如果您已经通过社保或其他商业保险公司报销了部分医疗费用，您可以提供分割单及理赔单据的复印件，以申请雇主责任险对剩余未报销部分进行索赔。这有助于确保您获得合理的赔付。",
    "output": "雇主责任险的医疗费用是不可重复报销的。如果您已经通过社保或其他商业保险公司报销了部分医疗费用，您可以提供分割单及理赔单据的复印件，以申请雇主责任险对剩余未报销部分进行索赔。这有助于确保您获得合理的赔付。",
    "instruction": "%s"
}`, taskSegment.Instruction, taskSegment.Instruction)
	}
	userPrompt += fmt.Sprintf("document: %s", taskSegment.SegmentContent)

	chatStream, _, err := s.apiSvc.FastChat().ChatCompletion(ctx, modelName, []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}, 0.7, 0.9, 0.5, 0.5, 1024, 1, nil, "", nil, nil)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "chat completion failed", "err", err)
		return
	}
	if chatStream.Choices == nil || len(chatStream.Choices) == 0 {
		//err = errors.New("chat completion failed")
		_ = level.Warn(logger).Log("msg", "chat completion failed", "err", err)
		return res, nil
	}
	content := chatStream.Choices[0].Message.Content
	if err = json.Unmarshal([]byte(content), &res); err != nil {
		_ = level.Warn(logger).Log("msg", "json.Unmarshal", "content", content, "err", err.Error())
		return res, nil
	}
	res.Instruction = taskSegment.Instruction
	return
}

func (s *service) GetCheckTaskDatasetSimilarLog(ctx context.Context, tenantId uint, taskId string) (res string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.DetectionStatus == types.DatasetAnnotationDetectionStatusProcessing && task.DetectionLog == "" {
		task.DetectionLog, err = s.apiSvc.Runtime().GetJobLogs(ctx, task.JobName)
		if err != nil {
			_ = level.Warn(logger).Log("msg", "GetJobLogs failed", "err", err.Error())
		}
	}
	return task.DetectionLog, nil
}

func (s *service) CancelCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	//if task.DetectionStatus != types.DatasetAnnotationDetectionStatusProcessing {
	//	err = errors.New("the annotation task is not processing, cannot be canceled")
	//	_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be canceled", "err", err)
	//	return
	//}
	if err = s.apiSvc.Runtime().RemoveJob(ctx, task.JobName); err != nil {
		err = errors.Wrap(err, "cancel job failed")
		_ = level.Error(logger).Log("msg", "cancel job failed", "err", err)
		return
	}
	task.DetectionStatus = types.DatasetAnnotationDetectionStatusCanceled
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}
	return
}

func (s *service) GetTaskInfo(ctx context.Context, tenantId uint, taskId string) (res taskDetail, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	sequence := strings.Split(task.DataSequence, "-")
	start, _ := strconv.Atoi(sequence[0])
	end, _ := strconv.Atoi(sequence[1])
	res = taskDetail{
		UUID:           task.UUID,
		Name:           task.Name,
		Remark:         task.Remark,
		AnnotationType: task.AnnotationType,
		Principal:      task.Principal,
		Status:         string(task.Status),
		Total:          task.Total,
		Completed:      task.Completed,
		DataSequence:   []int{start, end},
		CreatedAt:      task.CreatedAt,
		CompletedAt:    task.CompletedAt,
		Abandoned:      task.Abandoned,
		TrainTotal:     task.TrainTotal,
		TestTotal:      task.TestTotal,
		TestReport:     task.TestReport,
	}
	return
}

func (s *service) TaskDetectFinish(ctx context.Context, tenantId uint, taskId, testReport string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	detectionLog, err := s.apiSvc.Runtime().GetJobLogs(ctx, task.JobName)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "GetJobLogs failed", "err", err.Error())
	}
	if detectionLog != "" {
		task.DetectionLog = detectionLog
	}
	task.TestReport = testReport
	task.DetectionStatus = types.DatasetAnnotationDetectionStatusCompleted
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}
	if err = s.apiSvc.Runtime().RemoveJob(ctx, task.JobName); err != nil {
		_ = level.Error(logger).Log("msg", "remove job failed", "err", err)
	}
	return
}

func (s *service) CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}

	if task.Status != types.DatasetAnnotationStatusProcessing {
		err = errors.New("the annotation task is not processing, cannot be cleaned")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be cleaned", "err", err)
		return
	}
	now := time.Now()
	task.Status = types.DatasetAnnotationStatusCleaned
	task.CompletedAt = &now
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}

	return
}

func (s *service) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tasks, total, err := s.repository.DatasetTask().ListTasks(ctx, tenantId, name, page, pageSize, "DatasetDocument")
	if err != nil {
		err = errors.Wrap(err, "list tasks failed")
		_ = level.Warn(logger).Log("msg", "list tasks failed", "err", err)
		return
	}
	for _, task := range tasks {
		sequence := strings.Split(task.DataSequence, "-")
		start, _ := strconv.Atoi(sequence[0])
		end, _ := strconv.Atoi(sequence[1])
		res = append(res, taskDetail{
			UUID:            task.UUID,
			Name:            task.Name,
			Remark:          task.Remark,
			AnnotationType:  task.AnnotationType,
			Principal:       task.Principal,
			Status:          string(task.Status),
			Total:           task.Total,
			Completed:       task.Completed,
			DataSequence:    []int{start, end},
			CreatedAt:       task.CreatedAt,
			CompletedAt:     task.CompletedAt,
			Abandoned:       task.Abandoned,
			TrainTotal:      task.TrainTotal,
			TestTotal:       task.TestTotal,
			TestReport:      task.TestReport,
			DatasetName:     task.DatasetDocument.Name,
			DetectionStatus: string(task.DetectionStatus),
			JobName:         task.JobName,
		})
	}
	return
}

func (s *service) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, uuid)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status == types.DatasetAnnotationStatusProcessing {
		err = errors.New("the annotation task is processing, cannot be deleted")
		_ = level.Warn(logger).Log("msg", "the annotation task is processing, cannot be deleted", "err", err)
		return
	}
	if err = s.repository.DatasetTask().DeleteTask(ctx, tenantId, uuid); err != nil {
		err = errors.Wrap(err, "delete task failed")
		_ = level.Error(logger).Log("msg", "delete task failed", "err", err)
		return
	}
	return
}

func (s *service) GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status == types.DatasetAnnotationStatusCompleted {
		_ = level.Warn(logger).Log("msg", "the annotation task is completed, cannot be annotated")
		return
	}
	if task.Status != types.DatasetAnnotationStatusProcessing && task.Status != types.DatasetAnnotationStatusPending {
		err = errors.New("the annotation task is not processing, cannot be annotated")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be annotated", "err", err)
		return
	}
	segment, err := s.repository.DatasetTask().GetTaskOneSegment(ctx, task.ID, types.DatasetAnnotationStatusPending)
	if err != nil {
		err = errors.Wrap(err, "get task segment next failed")
		_ = level.Warn(logger).Log("msg", "get task segment next failed", "err", err)
		return
	}
	if strings.TrimSpace(segment.Instruction) == "" {
		if prevSegment, err := s.repository.DatasetTask().GetTaskSegmentPrev(ctx, task.ID, types.DatasetAnnotationStatusCompleted); err == nil {
			segment.Instruction = prevSegment.Instruction
		}
	}

	res = taskSegmentDetail{
		UUID:           segment.UUID,
		AnnotationType: string(segment.AnnotationType),
		SegmentContent: segment.SegmentContent,
		Status:         string(segment.Status),
		CreatedAt:      segment.CreatedAt,
		Document:       segment.Document,
		Instruction:    segment.Instruction,
		Input:          segment.Input,
		Question:       segment.Question,
		Intent:         segment.Intent,
		Output:         segment.Output,
		CreatorEmail:   segment.CreatorEmail,
	}
	return
}

func (s *service) AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status != types.DatasetAnnotationStatusProcessing && task.Status != types.DatasetAnnotationStatusPending {
		err = errors.New("the annotation task is not processing, cannot be annotated")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be annotated", "err", err)
		return
	}
	segment, err := s.repository.DatasetTask().GetTaskSegmentByUUID(ctx, task.ID, taskSegmentId)
	if err != nil {
		err = errors.Wrap(err, "get task segment by uuid failed")
		_ = level.Warn(logger).Log("msg", "get task segment by uuid failed", "err", err)
		return
	}
	segment.Document = req.Document
	segment.Instruction = req.Instruction
	segment.Input = req.Input
	segment.Question = req.Question
	segment.Intent = req.Intent
	segment.Output = req.Output
	segment.CreatorEmail = email
	segment.SegmentType = types.DatasetAnnotationSegmentTypeTrain
	segment.Status = types.DatasetAnnotationStatusCompleted
	if err = s.repository.DatasetTask().UpdateTaskSegment(ctx, segment); err != nil {
		err = errors.Wrap(err, "update task segment failed")
		_ = level.Error(logger).Log("msg", "update task segment failed", "err", err)
		return
	}

	task.Completed = task.Completed + 1
	task.TrainTotal = task.TrainTotal + 1

	if task.Status != types.DatasetAnnotationStatusProcessing {
		task.Status = types.DatasetAnnotationStatusProcessing
	}
	if task.TrainTotal+task.Abandoned >= task.Total {
		task.Status = types.DatasetAnnotationStatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	}
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return err
	}

	return
}

func (s *service) AbandonTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}

	if task.Status != types.DatasetAnnotationStatusProcessing && task.Status != types.DatasetAnnotationStatusPending {
		err = errors.New("the annotation task is not processing, cannot be abandoned")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be abandoned", "err", err)
		return
	}
	segment, err := s.repository.DatasetTask().GetTaskSegmentByUUID(ctx, task.ID, taskSegmentId)
	if err != nil {
		err = errors.Wrap(err, "get task segment by uuid failed")
		_ = level.Warn(logger).Log("msg", "get task segment by uuid failed", "err", err)
		return
	}
	segment.Status = types.DatasetAnnotationStatusAbandoned
	if err = s.repository.DatasetTask().UpdateTaskSegment(ctx, segment); err != nil {
		err = errors.Wrap(err, "update task segment failed")
		_ = level.Error(logger).Log("msg", "update task segment failed", "err", err)
		return
	}
	task.Abandoned += 1
	if task.Abandoned+task.TrainTotal >= task.Total {
		task.Status = types.DatasetAnnotationStatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	}
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}
	return
}

func (s *service) AsyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted {
		err = errors.New("the annotation task is not completed, cannot be split")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed, cannot be split", "err", err)
		return
	}

	// 获取所有标注好的数据
	segments, _, err := s.repository.DatasetTask().GetTaskSegments(ctx, annotationTask.ID, types.DatasetAnnotationStatusCompleted, 1, 100000)
	if err != nil {
		err = errors.Wrap(err, "get task segments failed")
		return
	}
	var datasetBody string
	for _, segment := range segments {
		line, _ := json.Marshal(DataAnnotationSegment{
			Instruction: segment.Instruction,
			Input:       segment.Input,
			Output:      segment.Output,
			Intent:      segment.Intent,
			Document:    segment.Document,
			Question:    segment.Question,
		})
		datasetBody += string(line) + "\n"
	}

	fileUrl, err := s.fileSvc.UploadToStorage(ctx, util.NewFile([]byte(datasetBody)), "jsonl")
	if err != nil {
		_ = level.Error(logger).Log("msg", "upload to storage failed", "err", err)
		return err
	}

	_ = level.Info(logger).Log("fileUrl", fileUrl, "msg", "upload to storage success")

	var jobName string
	tenantUUid, _ := ctx.Value(middleware.ContextKeyPublicTenantId).(string)
	auth, _ := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)

	// 组装脚本，调用api创建容器执行
	var envs []runtime.Env
	var envVars []string
	envs = append(envs, runtime.Env{
		Name:  "DATASET_ANALYZE_MODEL",
		Value: s.options.datasetModel,
	}, runtime.Env{
		Name:  "DATASET_PATH",
		Value: s.options.convertUrlFun(fileUrl),
	}, runtime.Env{
		Name:  "DATASET_TYPE",
		Value: annotationTask.AnnotationType,
	}, runtime.Env{
		Name:  "TENANT_ID",
		Value: tenantUUid,
	}, runtime.Env{
		Name:  "DATA_TASK_JOB_ID",
		Value: annotationTask.UUID,
	}, runtime.Env{
		Name:  "API_URL",
		Value: fmt.Sprintf("%s/api/mgr/annotation/task/%s/detect/finish", s.options.callbackHost, annotationTask.UUID), // 回调
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
	}, runtime.Env{
		Name:  "NO_PROXY",
		Value: os.Getenv("NO_PROXY"),
	})
	for _, v := range envs {
		envVars = append(envVars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	var gpuNum int
	if s.options.datasetDrive == "cuda" {
		gpuNum = 1
	}
	jobName, err = s.apiSvc.Runtime().CreateJob(ctx, runtime.Config{
		ServiceName: fmt.Sprintf("dataset-similar-task-%d", annotationTask.ID),
		Image:       s.options.datasetImage,
		GPU:         gpuNum,
		Command: []string{
			"/bin/bash",
			"/app/eval/analyze_similar_questions_and_intents.sh",
		},
		EnvVars:            envVars,
		GpuTolerationValue: s.options.gpuTolerationValue,
		Volumes:            []runtime.Volume{{Key: s.options.volumeName, Value: "/data"}},
	})
	if err != nil {
		err = errors.Wrap(err, "create job failed")
		_ = level.Error(logger).Log("msg", "create job failed", "err", err)
		return
	}
	_ = level.Info(logger).Log("msg", "create job success", "jobName", jobName)
	annotationTask.DetectionStatus = types.DatasetAnnotationDetectionStatusProcessing
	annotationTask.JobName = jobName
	if err = s.repository.DatasetTask().UpdateTask(ctx, annotationTask); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}
	return
}

func (s *service) SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted {
		err = errors.New("the annotation task is not completed, cannot be split")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed, cannot be split", "err", err)
		return
	}
	if req.TestPercent <= 0 || req.TestPercent >= 1 {
		err = errors.New("the test percent must be between 0 and 1")
		_ = level.Warn(logger).Log("msg", "the test percent must be between 0 and 1", "err", err)
		return
	}

	taskSegments, err := s.repository.DatasetTask().GetTaskSegmentByRand(ctx, annotationTask.ID, req.TestPercent,
		types.DatasetAnnotationStatusCompleted, types.DatasetAnnotationSegmentTypeTrain)
	if err != nil {
		err = errors.Wrap(err, "get dataset document segment by rand failed")
		_ = level.Warn(logger).Log("msg", "get dataset document segment by rand failed", "err", err)
		return
	}

	segmentIds := make([]uint, 0)
	for _, segment := range taskSegments {
		segmentIds = append(segmentIds, segment.ID)
	}

	if err = s.repository.DatasetTask().UpdateTaskSegmentType(ctx, segmentIds, types.DatasetAnnotationSegmentTypeTest); err != nil {
		err = errors.Wrap(err, "update task segment type failed")
		_ = level.Error(logger).Log("msg", "update task segment type failed", "err", err)
		return
	}
	annotationTask.TestTotal = len(segmentIds)
	annotationTask.TrainTotal = annotationTask.Total - annotationTask.TestTotal
	if err = s.repository.DatasetTask().UpdateTask(ctx, annotationTask); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}

	return
}

func (s *service) ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (filePath string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted {
		err = errors.New("the annotation task is not completed, cannot be exported")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed, cannot be exported", "err", err)
		return
	}
	segments, total, err := s.repository.DatasetTask().GetTaskSegments(ctx, annotationTask.ID, types.DatasetAnnotationStatusCompleted, 1, 100000)
	if err != nil {
		err = errors.Wrap(err, "get task segments failed")
		return
	}
	if total == 0 {
		err = errors.New("segments not found")
		return
	}
	trainSegments := make([]types.DatasetAnnotationTaskSegment, 0)
	testSegments := make([]types.DatasetAnnotationTaskSegment, 0)
	for _, segment := range segments {
		if segment.SegmentType == types.DatasetAnnotationSegmentTypeTrain {
			trainSegments = append(trainSegments, segment)
		} else {
			testSegments = append(testSegments, segment)
		}
	}

	storageDir := "."
	_ = os.MkdirAll(fmt.Sprintf("%s/temp_files", storageDir), os.ModePerm)
	trainFile := fmt.Sprintf("%s/temp_files/%s-train.jsonl", storageDir, annotationTask.UUID)
	testFile := fmt.Sprintf("%s/temp_files/%s-test.jsonl", storageDir, annotationTask.UUID)
	err = writeSegmentsToFile(trainSegments, trainFile, formatType, types.DatasetAnnotationType(annotationTask.AnnotationType))
	if err != nil {
		_ = level.Error(logger).Log("msg", "write segments to file failed", "err", err)
		return
	}
	if len(testSegments) > 0 {
		err = writeSegmentsToFile(testSegments, testFile, formatType, types.DatasetAnnotationType(annotationTask.AnnotationType))
		if err != nil {
			_ = level.Error(logger).Log("msg", "write segments to file failed", "err", err)
			return
		}
	}

	zipFilename := fmt.Sprintf("%s/temp_files/%s-files.zip", storageDir, annotationTask.UUID)
	err = createZip(zipFilename, []string{trainFile, testFile})
	if err != nil {
		err = errors.Wrap(err, "create zip failed")
		_ = level.Error(logger).Log("msg", "create zip failed", "err", err)
		return
	}
	_ = os.Remove(trainFile)
	if len(testSegments) > 0 {
		_ = os.Remove(testFile)
	}

	return zipFilename, err
}

func (s *service) DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted && annotationTask.Status != types.DatasetAnnotationStatusCleaned {
		err = errors.New("the annotation task is not completed or cleaned, cannot be deleted")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed or cleaned, cannot be deleted", "err", err)
		return
	}
	return
}

func (s *service) CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	//email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)

	datasetDocument, err := s.repository.DatasetTask().GetDatasetDocumentByUUID(ctx, tenantId, req.DatasetId)
	if err != nil {
		err = errors.Wrap(err, "get dataset document by uuid failed")
		_ = level.Warn(logger).Log("msg", "get dataset document by uuid failed", "err", err)
		return
	}

	if datasetDocument == nil {
		err = errors.New("dataset document not found")
		_ = level.Warn(logger).Log("msg", "dataset document not found", "err", err)
		return
	}

	total := req.DataSequence[1] - req.DataSequence[0]
	if total <= 0 {
		err = errors.Wrap(err, "total less than or equal to 0")
		_ = level.Warn(logger).Log("msg", "total less than or equal to 0", "err", err)
		return
	}

	datasetTask := types.DatasetAnnotationTask{
		DatasetDocumentId: datasetDocument.ID,
		UUID:              "task-" + uuid.New().String(),
		Name:              req.Name,
		Remark:            req.Remark,
		AnnotationType:    req.AnnotationType,
		TenantID:          tenantId,
		Principal:         req.Principal,
		Status:            types.DatasetAnnotationStatusPending,
		DetectionStatus:   types.DatasetAnnotationDetectionStatusPending,
		DataSequence:      fmt.Sprintf("%d-%d", req.DataSequence[0], req.DataSequence[1]),
	}

	if err = s.repository.DatasetTask().CreateTask(ctx, &datasetTask); err != nil {
		err = errors.Wrap(err, "create task failed")
		_ = level.Warn(logger).Log("msg", "create task failed", "err", err)
		return
	}

	documentSegments, err := s.repository.DatasetTask().GetDatasetDocumentSegmentByRange(ctx, datasetDocument.ID,
		req.DataSequence[0], req.DataSequence[1])
	if err != nil {
		err = errors.Wrap(err, "get dataset document segment by range failed")
		_ = level.Warn(logger).Log("msg", "get dataset document segment by range failed", "err", err)
		return
	}
	datasetTask.Total = len(documentSegments)

	var taskSegments []types.DatasetAnnotationTaskSegment
	for _, segment := range documentSegments {
		var dataSegment DataAnnotationSegment
		dats := types.DatasetAnnotationTaskSegment{
			DataAnnotationID: datasetTask.ID,
			UUID:             "das-" + uuid.New().String(),
			AnnotationType:   types.DatasetAnnotationType(datasetTask.AnnotationType),
			SegmentContent:   segment.SegmentContent,
			Status:           types.DatasetAnnotationStatusPending,
			SegmentID:        segment.ID,
		}
		if err = json.Unmarshal([]byte(segment.SegmentContent), &dataSegment); err == nil {
			dats.Instruction = dataSegment.Instruction
			dats.Input = dataSegment.Input
			dats.Output = dataSegment.Output
			dats.Intent = dataSegment.Intent
			dats.Document = dataSegment.Document
			dats.Question = dataSegment.Question
		}
		taskSegments = append(taskSegments, dats)
	}

	if err = s.repository.DatasetTask().AddTaskSegments(ctx, taskSegments); err != nil {
		err = errors.Wrap(err, "add task segments failed")
		_ = level.Error(logger).Log("msg", "add task segments failed", "err", err)
		return
	}

	if err = s.repository.DatasetTask().UpdateTask(ctx, &datasetTask); err != nil {
		_ = level.Warn(logger).Log("msg", "update task failed", "err", err)
		return err
	}

	return
}

func New(traceId string, logger log.Logger, repository repository.Repository, apiSvc services.Service, fileSvc files.Service, opts ...CreationOption) Service {
	logger = log.With(logger, "service", "datasettask")
	options := &CreationOptions{
		datasetImage: "dudulu/llmops:latest",
		datasetModel: "uer/sbert-base-chinese-nli",
		callbackHost: "http://localhost:8080",
		//volumeName:   "aigc-data-cfs",
		convertUrlFun: func(fileUrl string) string {
			return fileUrl
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
		options:    options,
		apiSvc:     apiSvc,
		fileSvc:    fileSvc,
	}
}

type DataAnnotationSegment struct {
	Instruction string `json:"instruction,omitempty"`
	Input       string `json:"input,omitempty"`
	Output      string `json:"output,omitempty"`
	Intent      string `json:"intent,omitempty"`
	Document    string `json:"document,omitempty"`
	Question    string `json:"question,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptMessage struct {
	Messages []Message `json:"messages"`
}

func writeSegmentsToFile(segments []types.DatasetAnnotationTaskSegment, filePath string, formatType string, annotationType types.DatasetAnnotationType) (err error) {
	//file, err := os.CreateTemp("", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "create file failed")
	}
	defer file.Close()

	for _, segment := range segments {
		var line interface{}

		if formatType == "conversation" {
			var messages []Message

			// 使用switch语句替代多个if-else
			switch annotationType {
			case types.DatasetAnnotationTypeFAQ:
				messages = append(messages, Message{"system", segment.Instruction},
					Message{"user", segment.Question + "\n" + segment.Intent},
					Message{"assistant", segment.Output})
			case types.DatasetAnnotationTypeGeneral:
				messages = append(messages, Message{"system", segment.Instruction},
					Message{"user", segment.Input},
					Message{"assistant", segment.Output})
			case types.DatasetAnnotationTypeRAG:
				messages = append(messages, Message{"system", segment.Instruction},
					Message{"user", fmt.Sprintf("背景知识: %s\n 问题: %s", segment.Document, segment.Question)},
					Message{"assistant", segment.Output})
			}

			line = PromptMessage{Messages: messages}
		} else {
			line = DataAnnotationSegment{
				Instruction: segment.Instruction,
				Input:       segment.Input,
				Output:      segment.Output,
				Intent:      segment.Intent,
				Document:    segment.Document,
				Question:    segment.Question,
			}
		}

		lineJSON, err := json.Marshal(line)
		if err != nil {
			return errors.Wrap(err, "marshalling line to json failed")
		}

		if _, err = file.WriteString(string(lineJSON) + "\n"); err != nil {
			return errors.Wrap(err, "write string to file failed")
		}
	}

	return nil
}

// createZip 创建一个包含指定文件的ZIP文件。
func createZip(zipFileName string, files []string) error {
	// 创建ZIP文件
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		err = errors.Wrap(err, "create zip file failed")
		return err
	}
	defer newZipFile.Close()

	// 创建一个新的zip.Writer
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 遍历所有文件，将它们添加到ZIP中
	for _, file := range files {
		if !util.FileExists(file) {
			continue
		}
		if err := addFileToZip(zipWriter, file); err != nil {
			err = errors.Wrap(err, "add file to zip failed")
			return err
		}
	}

	return nil
}

// addFileToZip 将单个文件添加到zip.Writer中
func addFileToZip(zipWriter *zip.Writer, fileName string) error {
	// 打开要添加的文件
	fileToZip, err := os.Open(fileName)
	if err != nil {
		err = errors.Wrap(err, "open file to zip failed")
		return err
	}
	defer fileToZip.Close()

	// 获取文件信息
	info, err := fileToZip.Stat()
	if err != nil {
		err = errors.Wrap(err, "get file info failed")
		return err
	}

	// 创建zip文件中的文件头信息
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		err = errors.Wrap(err, "create file info header failed")
		return err
	}

	// 设置压缩方法
	header.Method = zip.Deflate

	// 创建zip文件中的文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		err = errors.Wrap(err, "create file in zip failed")
		return err
	}

	// 将文件内容写入zip文件
	_, err = io.Copy(writer, fileToZip)
	return err
}
