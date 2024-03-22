package datasettask

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"time"
)

type (
	// taskListRequest 任务列表请求结构
	taskListRequest struct {
		name           string
		page, pageSize int
	}

	// taskCreateRequest 创建任务请求结构
	taskCreateRequest struct {
		// Name 任务名称
		Name string `json:"name" validate:"required,lt=50,gt=1"`
		// Remark 任务描述
		Remark string `json:"remark"`
		// DatasetId 数据集ID
		DatasetId string `json:"datasetId" validate:"required"`
		// TaskType 任务类型
		DataSequence []int `json:"dataSequence" validate:"required"`
		// TaskType 任务类型
		AnnotationType string `json:"annotationType" validate:"required"`
		// principal 负责人
		Principal string `json:"principal"`
	}

	// taskDetail 任务详情返回结构
	taskDetail struct {
		// UUID 任务ID
		UUID string `json:"uuid"`
		// Name 任务名称
		Name string `json:"name"`
		// DatasetName 数据集名称
		DatasetName string `json:"datasetName"`
		// Remark 任务描述
		Remark string `json:"remark"`
		// AnnotationType 标注类型
		AnnotationType string `json:"annotationType"`
		// Operator 操作人
		Operator string `json:"operator"`
		// Principal 负责人
		Principal string `json:"principal"`
		// Status 任务状态
		Status string `json:"status"`
		// DataSequence 数据序列
		DataSequence []int `json:"dataSequence"`
		// Total 标注总量
		Total int `json:"total"`
		// CreatedAt 创建时间
		CreatedAt time.Time `json:"createdAt"`
		// CompletedAt 完成时间
		CompletedAt *time.Time `json:"completedAt"`
		// Completed 完成标注量
		Completed int `json:"completed"`
		// Abandoned 废弃标注量
		Abandoned int `json:"abandoned"`
		// TrainTotal 训练数据总量
		TrainTotal int `json:"trainTotal"`
		// TestTotal 测试数据总量
		TestTotal int `json:"testTotal"`
		// TestReport 检测报告
		TestReport string `json:"testReport"`
		// DetectionStatus 检测状态
		DetectionStatus string `json:"detectionStatus"`
	}

	// taskSegmentDetail 任务样本详情返回结构
	taskSegmentDetail struct {
		// UUID 样本ID
		UUID string `json:"uuid"`
		// Status 样本状态
		Status string `json:"status"`
		// CreatedAt 创建时间
		CreatedAt time.Time `json:"createdAt"`
		// UpdatedAt 更新时间
		UpdatedAt time.Time `json:"updatedAt"`
		// AnnotationType 标注类型
		AnnotationType string `json:"annotationType"`
		// SegmentContent 样本内容
		SegmentContent string `json:"segmentContent"`
		// Document 标注文本
		Document string `json:"document"`
		// Instruction 标注说明
		Instruction string `json:"instruction"`
		// Input 标注输入
		Input string `json:"input"`
		// Question 标注问题
		Question string `json:"question"`
		// Intent 标注意图
		Intent string `json:"intent"`
		// Output 输出结果
		Output string `json:"output"`
		// CreatorEmail 创建人邮箱
		CreatorEmail string `json:"creatorEmail"`
		// Index 样本序号
		Index int `json:"index,omitempty"`
	}

	// taskSegmentAnnotationRequest 标注任务样本请求结构
	taskSegmentAnnotationRequest struct {
		Document    string `json:"document"`
		Instruction string `json:"instruction"`
		Input       string `json:"input"`
		Question    string `json:"question"`
		Intent      string `json:"intent"`
		Output      string `json:"output" validate:"lt=5000,gt=1"`
	}

	// taskExportAnnotationDataRequest 导出标注任务数据请求结构
	taskSplitAnnotationDataRequest struct {
		TestPercent float64 `json:"testPercent" validate:"required,min=0,max=1"`
	}

	// taskExportAnnotationDataRequest 导出标注任务数据请求结构
	taskExportAnnotationDataRequest struct {
		FormatType string `json:"formatType" validate:"required"`
	}

	// taskDetectFinishRequest 标注任务检测完成请求结构
	taskDetectFinishRequest struct {
		Status     string      `json:"status"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		TestReport string      `json:"-"`
	}
)

type Endpoints struct {
	CreateTaskEndpoint                    endpoint.Endpoint
	ListTasksEndpoint                     endpoint.Endpoint
	DeleteTaskEndpoint                    endpoint.Endpoint
	GetTaskSegmentNextEndpoint            endpoint.Endpoint
	CleanAnnotationTaskEndpoint           endpoint.Endpoint
	AnnotationTaskSegmentEndpoint         endpoint.Endpoint
	SplitAnnotationDataSegmentEndpoint    endpoint.Endpoint
	AbandonTaskSegmentEndpoint            endpoint.Endpoint
	ExportAnnotationDataEndpoint          endpoint.Endpoint
	TaskDetectFinishEndpoint              endpoint.Endpoint
	TaskInfoEndpoint                      endpoint.Endpoint
	CancelCheckTaskDatasetSimilarEndpoint endpoint.Endpoint
	AsyncCheckTaskDatasetSimilarEndpoint  endpoint.Endpoint
}

func NewEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateTaskEndpoint:                    makeCreateTaskEndpoint(s),
		ListTasksEndpoint:                     makeListTasksEndpoint(s),
		DeleteTaskEndpoint:                    makeDeleteTaskEndpoint(s),
		GetTaskSegmentNextEndpoint:            makeGetTaskSegmentNextEndpoint(s),
		CleanAnnotationTaskEndpoint:           makeCleanAnnotationTaskEndpoint(s),
		AnnotationTaskSegmentEndpoint:         makeAnnotationTaskSegmentEndpoint(s),
		SplitAnnotationDataSegmentEndpoint:    makeSplitAnnotationDataSegmentEndpoint(s),
		AbandonTaskSegmentEndpoint:            makeAbandonTaskSegmentEndpoint(s),
		ExportAnnotationDataEndpoint:          makeExportAnnotationDataEndpoint(s),
		TaskDetectFinishEndpoint:              makeTaskDetectFinishEndpoint(s),
		TaskInfoEndpoint:                      makeTaskInfoEndpoint(s),
		CancelCheckTaskDatasetSimilarEndpoint: makeCancelCheckTaskDatasetSimilarEndpoint(s),
		AsyncCheckTaskDatasetSimilarEndpoint:  makeAsyncCheckTaskDatasetSimilarEndpoint(s),
	}

	for _, m := range mdw["DatasetTask"] {
		eps.CreateTaskEndpoint = m(eps.CreateTaskEndpoint)
		eps.ListTasksEndpoint = m(eps.ListTasksEndpoint)
		eps.DeleteTaskEndpoint = m(eps.DeleteTaskEndpoint)
		eps.GetTaskSegmentNextEndpoint = m(eps.GetTaskSegmentNextEndpoint)
		eps.CleanAnnotationTaskEndpoint = m(eps.CleanAnnotationTaskEndpoint)
		eps.AnnotationTaskSegmentEndpoint = m(eps.AnnotationTaskSegmentEndpoint)
		eps.SplitAnnotationDataSegmentEndpoint = m(eps.SplitAnnotationDataSegmentEndpoint)
		eps.AbandonTaskSegmentEndpoint = m(eps.AbandonTaskSegmentEndpoint)
		eps.ExportAnnotationDataEndpoint = m(eps.ExportAnnotationDataEndpoint)
		eps.TaskDetectFinishEndpoint = m(eps.TaskDetectFinishEndpoint)
		eps.TaskInfoEndpoint = m(eps.TaskInfoEndpoint)
		eps.CancelCheckTaskDatasetSimilarEndpoint = m(eps.CancelCheckTaskDatasetSimilarEndpoint)
		eps.AsyncCheckTaskDatasetSimilarEndpoint = m(eps.AsyncCheckTaskDatasetSimilarEndpoint)
	}
	return eps
}

func makeCreateTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(taskCreateRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		err := s.CreateTask(ctx, tenantId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeListTasksEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(taskListRequest)
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		res, total, err := s.ListTasks(ctx, tenantId, req.name, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeDeleteTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return encode.Response{
				Error: encode.InvalidParams.Wrap(err),
			}, err
		}
		err := s.DeleteTask(ctx, tenantId, taskId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeGetTaskSegmentNextEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		res, err := s.GetTaskSegmentNext(ctx, tenantId, taskId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCleanAnnotationTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		err := s.CleanAnnotationTask(ctx, tenantId, taskId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAnnotationTaskSegmentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		segmentId, ok := ctx.Value(contextKeyDatasetTaskSegmentId).(string)
		if !ok {
			err := errors.New("invalid segment id")
			return nil, err
		}
		req := request.(taskSegmentAnnotationRequest)
		err := s.AnnotationTaskSegment(ctx, tenantId, taskId, segmentId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeSplitAnnotationDataSegmentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		req := request.(taskSplitAnnotationDataRequest)
		err := s.SplitAnnotationDataSegment(ctx, tenantId, taskId, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAbandonTaskSegmentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		segmentId, ok := ctx.Value(contextKeyDatasetTaskSegmentId).(string)
		if !ok {
			err := errors.New("invalid segment id")
			return nil, err
		}
		err := s.AbandonTaskSegment(ctx, tenantId, taskId, segmentId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeExportAnnotationDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		req, _ := request.(taskExportAnnotationDataRequest)
		filePath, err := s.ExportAnnotationData(ctx, tenantId, taskId, req.FormatType)
		return encode.Response{
			Data:  filePath,
			Error: err,
		}, err
	}
}

func makeTaskDetectFinishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		req := request.(taskDetectFinishRequest)
		b, _ := json.Marshal(req.Data)
		if req.Status != "success" {
			req.TestReport = req.Message
		} else {
			req.TestReport = string(b)
		}
		err := s.TaskDetectFinish(ctx, tenantId, taskId, req.TestReport)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeTaskInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		res, err := s.GetTaskInfo(ctx, tenantId, taskId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCancelCheckTaskDatasetSimilarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		err := s.CancelCheckTaskDatasetSimilar(ctx, tenantId, taskId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAsyncCheckTaskDatasetSimilarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		taskId, ok := ctx.Value(contextKeyDatasetTaskId).(string)
		if !ok {
			err := errors.New("invalid task id")
			return nil, err
		}
		err := s.AsyncCheckTaskDatasetSimilar(ctx, tenantId, taskId)
		return encode.Response{
			Error: err,
		}, err
	}
}
