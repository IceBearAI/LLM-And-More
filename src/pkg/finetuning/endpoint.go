package finetuning

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	CreateJobRequest struct {
		FileId            string  `json:"fileId" validate:"required"`     // 文件ID
		BaseModel         string  `json:"baseModel" validate:"required"`  // 基础模型
		TrainEpoch        int     `json:"trainEpoch" validate:"required"` // 训练轮次
		Suffix            string  `json:"suffix"`                         // 后缀
		Remark            string  `json:"remark"`                         // 备注
		TrainBatchSize    int     `json:"trainBatchSize"`                 // 训练批次大小
		EvalBatchSize     int     `json:"evalBatchSize"`                  //
		AccumulationSteps int     `json:"accumulationSteps"`              // 梯度累加步数
		ProcPerNode       int     `json:"procPerNode" validate:"gt=0"`    // 使用GPU数量
		LearningRate      float64 `json:"learningRate"`                   // 学习率
		ModelMaxLength    int     `json:"modelMaxLength"`                 // 模型最大长度
		TenantId          uint    `json:"tenantId"`                       // 租户ID
		TrainPublisher    string  `json:"trainPublisher"`                 // 训练发布者
		Lora              bool    `json:"lora"`                           // 是否微调lora
	}

	LogEntry struct {
		Timestamp    time.Time `json:"timestamp"`
		Loss         float64   `json:"loss"`
		LearningRate float64   `json:"learning_rate"`
		Epoch        float64   `json:"epoch"`
	}

	TrainAnalysisDetail struct {
		Timestamp time.Time `json:"timestamp"`
		Value     string    `json:"value"`
	}

	TrainAnalysisObject struct {
		List []TrainAnalysisDetail `json:"list"`
	}

	TrainAnalysis struct {
		Epoch        TrainAnalysisObject `json:"epoch"`
		Loss         TrainAnalysisObject `json:"loss"`
		LearningRate TrainAnalysisObject `json:"learningRate"`
	}

	JobResponse struct {
		Id                uint          `json:"id"`                      // ID
		JobId             string        `json:"jobId"`                   // 任务ID
		BaseModel         string        `json:"baseModel"`               // 基础模型
		TrainEpoch        int           `json:"trainEpoch"`              // 训练轮次
		TrainStatus       string        `json:"trainStatus"`             // 训练状态
		TrainDuration     string        `json:"trainDuration"`           // 训练时长 单位秒
		Process           float64       `json:"process"`                 // 训练进度
		FineTunedModel    string        `json:"fineTunedModel"`          // 微调模型
		Remark            string        `json:"remark"`                  // 备注
		FinishedAt        string        `json:"finishedAt"`              // 训练完成时间
		CreatedAt         time.Time     `json:"createdAt"`               // 创建时间
		TrainPublisher    string        `json:"trainPublisher"`          // 训练发布者
		TrainLog          string        `json:"trainLog"`                // 训练日志
		ErrorMessage      string        `json:"errorMessage"`            // 错误信息
		Lora              bool          `json:"lora"`                    // 是否微调lora
		TrainAnalysis     TrainAnalysis `json:"trainAnalysis,omitempty"` // 训练分析
		Suffix            string        `json:"suffix"`                  // 后缀
		ModelMaxLength    int           `json:"modelMaxLength"`          // 模型最大长度
		TrainBatchSize    int           `json:"trainBatchSize"`          // 训练批次大小
		LearningRate      string        `json:"learningRate"`            // 学习率
		FileUrl           string        `json:"fileUrl"`                 // 文件地址
		FileId            string        `json:"fileId"`                  // 文件ID
		StartTrainTime    string        `json:"startTrainTime"`          // 开始训练时间
		ProcPerNode       int           `json:"procPerNode"`             // 使用GPU数量
		EvalBatchSize     int           `json:"evalBatchSize"`           // 评估批次大小
		AccumulationSteps int           `json:"accumulationSteps"`       // 梯度累加步数
	}

	ListJobRequest struct {
		Page           int    `json:"page"`
		PageSize       int    `json:"pageSize"`
		FineTunedModel string `json:"fineTunedModel"` // 微调模型
		TrainStatus    string `json:"trainStatus"`    // 训练状态
	}

	ListJobResponse struct {
		List  []JobResponse `json:"list"`
		Total int64         `json:"total"`
	}

	JobIdRequest struct {
		JobId string `json:"jobId"`
	}

	DashBoardResponse struct {
		WaitingJobCount    int64  `json:"waitingJobCount"`    // 等待中任务数
		SuccessJobCount    int64  `json:"successJobCount"`    // 已完成任务数
		TotalDurationCount string `json:"totalDurationCount"` // 总训练时长
	}

	Template struct {
		Id            uint      `json:"id"`            // ID
		Name          string    `json:"name"`          // 模版名称
		BaseModel     string    `json:"baseModel"`     // 基础模型
		Content       string    `json:"content"`       // 模版内容
		MaxTokens     int       `json:"maxTokens"`     // 最大token数
		Params        string    `json:"params"`        // 模版所需要参数
		BaseModelPath string    `json:"baseModelPath"` // 基础模型路径
		TrainImage    string    `json:"trainImage"`    // 训练镜像
		ScriptFile    string    `json:"scriptFile"`    // 脚本文件
		OutputDir     string    `json:"outputDir"`     // 输出目录
		TemplateType  string    `json:"templateType"`  // 模版类型
		CreatedAt     time.Time `json:"createdAt"`     // 创建时间
		UpdatedAt     time.Time `json:"updatedAt"`     // 更新时间
		Remark        string    `json:"remark"`        // 备注
	}

	ListTemplateRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	ListTemplateResponse struct {
		List  []Template `json:"list"`
		Total int64      `json:"total"`
	}
	EstimateResponse struct {
		Time string `json:"time"`
	}

	// jobResult createJobResponse is the response struct for the CreateJob endpoint.
	jobResult struct {
		// Object is the type of object.
		Object string `json:"object"`
		// Id is the id of the job.
		Id string `json:"id"`
		// Model is the name of the model to use.
		Model string `json:"model"`
		// CreatedAt is the time the job was created.
		CreatedAt int64 `json:"created_at"`
		// FinishedAt is the time the job finished.
		FinishedAt int64 `json:"finished_at,omitempty"`
		// FineTunedModel is the name of the fine-tuned model.
		FineTunedModel string `json:"fine_tuned_model,omitempty"`
		// OrganizationId is the id of the organization.
		OrganizationId string `json:"organization_id"`
		// ResultFiles is the list of result files.
		ResultFiles []string `json:"result_files"`
		// Status is the status of the job.
		Status string `json:"status"`
		// ValidationFile is the path to the validation file.
		ValidationFile string `json:"validation_file,omitempty"`
		// TrainingFile is the path to the training file.
		TrainingFile string `json:"training_file"`
		// Error is the error message.
		Error interface{} `json:"error,omitempty"`
		// HyperParameters is the hyperparameters to use.
		HyperParameters jobHyperParameters `json:"hyperparameters"`
		// TrainedTokens is the number of trained tokens.
		TrainedTokens int `json:"trained_tokens,omitempty"`
	}

	// jobHyperParameters is the hyperparameters to use.
	jobHyperParameters struct {
		// NEpochs is the number of training epochs.
		NEpochs int `json:"n_epochs"`
	}

	updateTrainStatusRequest struct {
		Status  string `json:"status" validate:"required"`
		Message string `json:"message"`
	}
)

type Endpoints struct {
	CreateJobEndpoint               endpoint.Endpoint
	ListJobEndpoint                 endpoint.Endpoint
	CancelJobEndpoint               endpoint.Endpoint
	DashBoardEndpoint               endpoint.Endpoint
	BaseModelEndpoint               endpoint.Endpoint
	DeleteJobEndpoint               endpoint.Endpoint
	GetJobEndpoint                  endpoint.Endpoint
	EstimateEndpoint                endpoint.Endpoint
	UpdateJobFinishedStatusEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dwm map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateJobEndpoint:               makeCreateJobEndpoint(s),
		ListJobEndpoint:                 makeListJobEndpoint(s),
		CancelJobEndpoint:               makeCancelJobEndpoint(s),
		DashBoardEndpoint:               makeDashBoardEndpoint(s),
		BaseModelEndpoint:               makeBaseModelEndpoint(s),
		DeleteJobEndpoint:               makeDeleteJobEndpoint(s),
		GetJobEndpoint:                  makeGetJobEndpoint(s),
		EstimateEndpoint:                makeEstimateEndpoint(s),
		UpdateJobFinishedStatusEndpoint: makeUpdateJobFinishedStatusEndpoint(s),
	}
	for _, m := range dwm["FineTuning"] {
		eps.CreateJobEndpoint = m(eps.CreateJobEndpoint)
		eps.ListJobEndpoint = m(eps.ListJobEndpoint)
		eps.CancelJobEndpoint = m(eps.CancelJobEndpoint)
		eps.DashBoardEndpoint = m(eps.DashBoardEndpoint)
		eps.BaseModelEndpoint = m(eps.BaseModelEndpoint)
		eps.DeleteJobEndpoint = m(eps.DeleteJobEndpoint)
		eps.GetJobEndpoint = m(eps.GetJobEndpoint)
		eps.EstimateEndpoint = m(eps.EstimateEndpoint)
		//eps.UpdateJobFinishedStatusEndpoint = m(eps.UpdateJobFinishedStatusEndpoint)
	}
	return eps
}

func makeCreateJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(CreateJobRequest)
		req.TenantId, _ = middleware.GetTenantId(ctx)
		req.TrainPublisher, _ = middleware.GetEmail(ctx)
		resp, err := s.CreateJob(ctx, tenantId, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(ListJobRequest)
		resp, err := s.ListJob(ctx, tenantId, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeCancelJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(JobIdRequest)
		err = s.CancelJob(ctx, tenantId, req.JobId)
		return encode.Response{
			Error: err,
			Data:  nil,
		}, err
	}
}

func makeDashBoardEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		resp, err := s.DashBoard(ctx, tenantId)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeBaseModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		listTemplate, err := s.ListTemplate(ctx, tenantId, ListTemplateRequest{
			Page:     -1,
			PageSize: -1,
		})
		models := make([]string, 0)
		if err != nil {
			return encode.Response{
				Data:  models,
				Error: err,
			}, err
		}
		for _, v := range listTemplate.List {
			models = append(models, v.BaseModel)
		}
		data := make(map[string]interface{})
		data["list"] = models
		return encode.Response{
			Data:  data,
			Error: err,
		}, err
	}
}

func makeDeleteJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(JobIdRequest)
		err = s.DeleteJob(ctx, tenantId, req.JobId)
		return encode.Response{
			Error: err,
			Data:  nil,
		}, err
	}
}

func makeGetJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(JobIdRequest)
		resp, err := s.GetJob(ctx, tenantId, req.JobId)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeEstimateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		req := request.(CreateJobRequest)
		resp, err := s.Estimate(ctx, tenantId, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeUpdateJobFinishedStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fineTuningJobId, _ := ctx.Value(contextKeyFineTuningJobId).(string)
		req, _ := request.(updateTrainStatusRequest)
		err = s.UpdateJobFinishedStatus(ctx, fineTuningJobId, types.TrainStatus(req.Status), req.Message)
		return encode.Response{
			Error: err,
		}, err
	}
}
