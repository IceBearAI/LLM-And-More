package models

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	Model struct {
		Id            uint   `json:"id"`
		BaseModelName string `json:"baseModelName"`
		ProviderName  string `json:"providerName"`
		ModelName     string `json:"modelName"`
		ModelType     string `json:"modelType"`
		MaxTokens     int    `json:"maxTokens"`
		//IsPrivate    bool      `json:"isPrivate"`
		IsFineTuning   bool      `json:"isFineTuning"`
		Enabled        bool      `json:"enabled"`
		Remark         string    `json:"remark"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		Tenants        []Tenant  `json:"tenants"`
		DeployStatus   string    `json:"deployStatus"`
		Operation      []string  `json:"operation"`
		JobId          string    `json:"jobId"`
		LastOperator   string    `json:"lastOperator"`
		Parameters     float64   `json:"parameters"`
		Replicas       int       `json:"replicas"`     //并行/实例数量
		Label          string    `json:"label"`        //调度标签
		K8sCluster     string    `json:"k8sCluster"`   //k8s集群
		InferredType   string    `json:"inferredType"` //推理类型cpu,gpu
		Gpu            int       `json:"gpu"`          //GPU数
		Cpu            int       `json:"cpu"`          //CPU核数
		Memory         int       `json:"memory"`       //内存G
		ServiceName    string    `json:"serviceName"`
		ContainerNames []string  `json:"containerNames"`
	}

	Tenant struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	}

	CreateModelRequest struct {
		ModelName     string `json:"modelName" validate:"required"`
		ModelType     string `json:"modelType" validate:"required"`
		BaseModelName string `json:"baseModelName"`
		MaxTokens     int    `json:"maxTokens" validate:"required"`
		TenantId      []uint `json:"tenantId"`
		//IsPrivate    bool    `json:"isPrivate"`
		IsFineTuning bool    `json:"isFineTuning"`
		Enabled      bool    `json:"enabled"`
		Remark       string  `json:"remark"`
		ProviderName string  `json:"providerName"`
		Email        string  `json:"email"`
		Parameters   float64 `json:"parameters"`
		Replicas     int     `json:"replicas"`     //并行/实例数量
		Label        string  `json:"label"`        //调度标签
		K8sCluster   string  `json:"k8sCluster"`   //k8s集群
		InferredType string  `json:"inferredType"` //推理类型cpu,gpu
		Gpu          int     `json:"gpu"`          //GPU数
		Cpu          int     `json:"cpu"`          //CPU核数
		Memory       int     `json:"memory"`       //内存G
	}

	UpdateModelRequest struct {
		Id            uint    `json:"id" validate:"required"`
		TenantId      *[]uint `json:"tenantId"`
		MaxTokens     *int    `json:"maxTokens"`
		Enabled       *bool   `json:"enabled"`
		Remark        *string `json:"remark"`
		BaseModelName string  `json:"baseModelName"`
		Replicas      int     `json:"replicas"`     //并行/实例数量
		Label         string  `json:"label"`        //调度标签
		K8sCluster    string  `json:"k8sCluster"`   //k8s集群
		InferredType  string  `json:"inferredType"` //推理类型cpu,gpu
		Gpu           int     `json:"gpu"`          //GPU数
		Cpu           int     `json:"cpu"`          //CPU核数
		Memory        int     `json:"memory"`       //内存G
	}
	ListModelRequest struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		ModelName string `json:"modelName,omitempty"`
		Enabled   *bool  `json:"enabled,omitempty"`
		//IsPrivate    *bool  `json:"isPrivate,omitempty"`
		IsFineTuning *bool  `json:"isFineTuning,omitempty"`
		ProviderName string `json:"providerName"`
		ModelType    string `json:"modelType"`
	}
	ListModelResponse struct {
		Models []Model `json:"list"`
		Total  int64   `json:"total"`
	}
	IdRequest struct {
		Id uint `json:"id"`
	}

	Eval struct {
		Id          uint      `json:"id"`
		ModelName   string    `json:"modelName"`
		DatasetType string    `json:"datasetType"`
		Progress    float64   `json:"progress"`
		Score       float64   `json:"score"`
		CreatedAt   time.Time `json:"createdAt"`
		Status      string    `json:"status"`
		Duration    string    `json:"duration"`
		MetricName  string    `json:"metricName"`
		Remark      string    `json:"remark"`
		EvalTotal   int       `json:"evalTotal"`
		StartedAt   string    `json:"startedAt"`
	}

	CreateEvalRequest struct {
		ModelName     string  `json:"modelName" validate:"required"`
		DatasetType   string  `json:"datasetType" validate:"required"`
		MetricName    string  `json:"metricName" validate:"required"`
		Remark        string  `json:"remark"`
		EvalPercent   float64 `json:"evalPercent"`
		DatasetFileId string  `json:"fileId"`
	}

	ListEvalRequest struct {
		Page        int    `json:"page"`
		PageSize    int    `json:"pageSize"`
		ModelName   string `json:"modelName"`
		DatasetType string `json:"datasetType"`
		MetricName  string `json:"metricName"`
		Status      string `json:"status"`
	}

	ListEvalResponse struct {
		Total int64  `json:"total"`
		List  []Eval `json:"list"`
	}

	ModelDeployRequest struct {
		Id           uint   `json:"id"`
		Replicas     int    `json:"replicas" validate:"gte=1"`
		Label        string `json:"label"`
		InferredType string `json:"inferredType"` //推理类型
		Gpu          int    `json:"gpu"`          //gpu 数量
		Cpu          int    `json:"cpu"`          //cpu 数量
		Quantization string `json:"quantization"`
		Vllm         bool   `json:"vllm"`
		// 指定每个 GPU 用于存储模型权重的最大内存。这允许它为激活分配更多内存，因此您可以使用更长的上下文长度或更大的批量大小。
		MaxGpuMemory int    `json:"maxGpuMemory"`
		K8sCluster   string `json:"k8sCluster"` //k8s集群
	}

	// modelLogRequest 模型日志请求
	modelLogRequest struct {
		ModelName     string
		ContainerName string
	}
)

type Endpoints struct {
	ListModelsEndpoint    endpoint.Endpoint
	CreateModelEndpoint   endpoint.Endpoint
	UpdateModelEndpoint   endpoint.Endpoint
	DeleteModelEndpoint   endpoint.Endpoint
	GetModelEndpoint      endpoint.Endpoint
	DeployModelEndpoint   endpoint.Endpoint
	UndeployModelEndpoint endpoint.Endpoint
	CreateEvalEndpoint    endpoint.Endpoint
	ListEvalEndpoint      endpoint.Endpoint
	CancelEvalEndpoint    endpoint.Endpoint
	DeleteEvalEndpoint    endpoint.Endpoint
	GetModelLogsEndpoint  endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListModelsEndpoint:    makeListModelsEndpoint(s),
		CreateModelEndpoint:   makeCreateModelEndpoint(s),
		UpdateModelEndpoint:   makeUpdateModelEndpoint(s),
		DeleteModelEndpoint:   makeDeleteModelEndpoint(s),
		GetModelEndpoint:      makeGetModelEndpoint(s),
		DeployModelEndpoint:   makeDeployModelEndpoint(s),
		UndeployModelEndpoint: makeUndeployModelEndpoint(s),
		CreateEvalEndpoint:    makeCreateEvalEndpoint(s),
		ListEvalEndpoint:      makeListEvalEndpoint(s),
		CancelEvalEndpoint:    makeCancelEvalEndpoint(s),
		DeleteEvalEndpoint:    makeDeleteEvalEndpoint(s),
		GetModelLogsEndpoint:  makeGetModelLogsEndpoint(s),
	}
	for _, m := range dmw["Model"] {
		eps.ListModelsEndpoint = m(eps.ListModelsEndpoint)
		eps.CreateModelEndpoint = m(eps.CreateModelEndpoint)
		eps.UpdateModelEndpoint = m(eps.UpdateModelEndpoint)
		eps.DeleteModelEndpoint = m(eps.DeleteModelEndpoint)
		eps.GetModelEndpoint = m(eps.GetModelEndpoint)
		eps.DeployModelEndpoint = m(eps.DeployModelEndpoint)
		eps.UndeployModelEndpoint = m(eps.UndeployModelEndpoint)
		eps.CreateEvalEndpoint = m(eps.CreateEvalEndpoint)
		eps.ListEvalEndpoint = m(eps.ListEvalEndpoint)
		eps.CancelEvalEndpoint = m(eps.CancelEvalEndpoint)
		eps.GetModelLogsEndpoint = m(eps.GetModelLogsEndpoint)
	}
	return eps
}

func makeGetModelLogsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		modelName, _ := ctx.Value(contextKeyModelName).(string)
		containerName, _ := ctx.Value(contextKeyModelContainerName).(string)
		resp, err := s.GetModelLogs(ctx, modelName, containerName)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListModelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListModelRequest)
		resp, err := s.ListModels(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeCreateModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email, _ := middleware.GetEmail(ctx)
		req := request.(CreateModelRequest)
		req.Email = email
		resp, err := s.CreateModel(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeUpdateModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateModelRequest)
		err = s.UpdateModel(ctx, req)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeDeleteModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.DeleteModel(ctx, req.Id)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeGetModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		resp, err := s.GetModel(ctx, req.Id)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeDeployModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ModelDeployRequest)
		err = s.Deploy(ctx, req)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeUndeployModelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.Undeploy(ctx, req.Id)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeCreateEvalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateEvalRequest)
		resp, err := s.CreateEval(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListEvalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListEvalRequest)
		resp, err := s.ListEval(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeCancelEvalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.CancelEval(ctx, req.Id)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeDeleteEvalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.DeleteEval(ctx, req.Id)
		resp := struct{}{}
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}
