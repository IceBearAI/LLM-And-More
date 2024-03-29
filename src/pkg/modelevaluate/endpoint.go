package modelevaluate

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	listRequest struct {
		ModelId        int    `json:"modelId" validate:"required,min=1"`
		Status         string `json:"status"`
		EvalTargetType string `json:"evalTargetType"`
		Page           int    `json:"page"`
		PageSize       int    `json:"pageSize"`
	}
	listResult struct {
		Id             int       `json:"id"`
		Uuid           string    `json:"uuid"`
		ModelID        int       `json:"modelId"`
		Status         string    `json:"status"`
		EvalTargetType string    `json:"evalTargetType"`
		StatusMsg      string    `json:"statusMsg"`
		Remark         string    `json:"remark"`
		DataSize       int       `json:"dataSize"`
		Complete       string    `json:"complete"`
		Score          float64   `json:"score"`
		DataType       string    `json:"dataType"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	createRequest struct {
		FileId         string `json:"fileId"`
		ModelId        int    `json:"modelId"`
		EvalTargetType string `json:"evalTargetType"`
		MaxLength      int    `json:"maxLength"`
		BatchSize      int    `json:"batchSize"`
		Replicas       int    `json:"replicas"`     //并行/实例数量
		Label          string `json:"label"`        //调度标签
		K8sCluster     string `json:"k8sCluster"`   //k8s集群
		InferredType   string `json:"inferredType"` //推理类型cpu,gpu
		Gpu            int    `json:"gpu"`          //GPU数
		Cpu            int    `json:"cpu"`          //CPU核数
		Memory         int    `json:"memory"`       //内存G
		Remark         string `json:"remark"`       //备注
		maxGpuMemory   int    `json:"maxGpuMemory"` //GPU大于2时，gpu内存
	}

	createResult struct {
	}

	cancelRequest struct {
		Uuid string `json:"uuid" validate:"required"`
	}

	deleteRequest struct {
		Uuid string `json:"uuid" validate:"required"`
	}

	fiveGraphRequest struct {
		CurrentModelId         int `json:"currentModelId" validate:"required,min=1"`         //当前模型
		CurrentModelEvaluateId int `json:"currentModelEvaluateId" validate:"required,min=1"` //当前模型评测
		Compare1ModelId        int `json:"compare1ModelId"`                                  //比对1模型
		Compare2ModelId        int `json:"compare2ModelId"`                                  //比对1模型
	}
	fiveGraphResult struct {
		RiskOver     bool      `json:"riskOver"`     // 过拟合风险
		RiskUnder    bool      `json:"riskUnder"`    // 欠拟合风险
		RiskDisaster bool      `json:"riskDisaster"` // 灾难性遗忘风险
		Remind       string    `json:"remind"`       //建议提醒
		Value        []float64 `json:"value"`        // value: [10, 9.8, 8.5, 3.9, 0] 中文能力，推理能力，指令遵从能力，创新能力，阅读理解
		Score        float64   `json:"score"`        // 总分
		ModelId      int       `json:"modelId"`      // 模型ID
		Name         string    `json:"name"`         // 模型名称
		IsFineTuning bool      `json:"isFineTuning"` // 是否为微调模型，如果是显示风险提示
	}

	finishRequest struct {
		JobId  string `json:"jobId"`  //uuid
		Result string `json:"result"` //回调json
	}

	callbackResult struct {
		Status string `json:"status,omitempty"`
		Data   struct {
			Config struct {
				ModelNameOrPath    string `json:"model_name_or_path,omitempty"`
				DatasetPath        string `json:"dataset_path,omitempty"`
				EvaluationMetrics  string `json:"evaluation_metrics,omitempty"`
				MaxSeqLen          int    `json:"max_seq_len,omitempty"`
				PerDeviceBatchSize int    `json:"per_device_batch_size,omitempty"`
				GpuID              int    `json:"gpu_id,omitempty"`
			} `json:"config,omitempty"`
			OverallEvaluationMetrics map[string]interface{} `json:"overall_evaluation_metrics,omitempty"`
			DetailedResults          []struct {
				Question    string                 `json:"question,omitempty"`
				Reference   string                 `json:"reference,omitempty"`
				ModelOutput string                 `json:"model_output,omitempty"`
				Scores      map[string]interface{} `json:"scores,omitempty"`
			} `json:"detailed_results,omitempty"`
		} `json:"data,omitempty"`
	}

	fiveResult struct {
		Status string `json:"status,omitempty"`
		Data   struct {
			InferenceAbility     float64 `json:"inference_ability,omitempty"`      // 推理能力
			ReadingComprehension float64 `json:"reading_comprehension,omitempty"`  //阅读理解能力
			ChineseLanguageSkill float64 `json:"chinese_language_skill,omitempty"` //中文能力
			CommandCompliance    float64 `json:"command_compliance,omitempty"`     //指令遵从能力
			InnovationCapacity   float64 `json:"innovation_capacity,omitempty"`    //创新能力
		} `json:"data,omitempty"`
	}
)

type Endpoints struct {
	ListEndpoint       endpoint.Endpoint
	CreateEndpoint     endpoint.Endpoint
	CancelEndpoint     endpoint.Endpoint
	DeleteEndpoint     endpoint.Endpoint
	FiveGraphEndpoint  endpoint.Endpoint
	FinishEndpoint     endpoint.Endpoint
	GetEvalLogEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:       makeListEndpoint(s),
		CreateEndpoint:     makeCreateEndpoint(s),
		CancelEndpoint:     makeCancelEndpoint(s),
		DeleteEndpoint:     makeDeleteEndpoint(s),
		FiveGraphEndpoint:  makeFiveGraphEndpoint(s),
		FinishEndpoint:     makeFinishEndpoint(s),
		GetEvalLogEndpoint: makeGetEvalLogEndpoint(s),
	}

	for _, m := range dmw["Evaluate"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.CancelEndpoint = m(eps.CancelEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.FiveGraphEndpoint = m(eps.FiveGraphEndpoint)
		eps.FinishEndpoint = m(eps.FinishEndpoint)
		eps.GetEvalLogEndpoint = m(eps.GetEvalLogEndpoint)

	}

	return eps
}

func makeGetEvalLogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenantId, _ := ctx.Value(middleware.ContextKeyTenantId).(uint)
		modelId, _ := ctx.Value(contextKeyModelName).(string)
		modelEvalId, _ := ctx.Value(contextKeyModelEvaluateId).(string)
		res, err := s.GetEvalLog(ctx, tenantId, modelId, modelEvalId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req)

		return encode.Response{
			Data: map[string]interface{}{
				"list":     res,
				"total":    total,
				"page":     req.Page,
				"pageSize": req.PageSize,
			},
			Error: err,
		}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createRequest)
		err = s.Create(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCancelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(cancelRequest)
		err = s.Cancel(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteRequest)
		err = s.Delete(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}

func makeFiveGraphEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(fiveGraphRequest)
		res1, res2, res3, err := s.FiveGraph(ctx, req)

		return encode.Response{
			Data: map[string]interface{}{
				"current":  res1,
				"compare1": res2,
				"compare2": res3,
			},
			Error: err,
		}, err
	}
}

func makeFinishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(finishRequest)
		err = s.EvalFinish(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}
