package sys

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	Dict struct {
		ID        uint        `json:"id"`
		ParentID  uint        `json:"parentId"`
		Code      string      `json:"code"`
		DictValue interface{} `json:"dictValue"`
		DictLabel string      `json:"dictLabel"`
		DictType  string      `json:"dictType"`
		Sort      int         `json:"sort"`
		Remark    string      `json:"remark"`
		CreatedAt time.Time   `json:"createdAt"`
		UpdatedAt time.Time   `json:"updatedAt"`
		Children  []Dict      `json:"children"`
	}
	CreateDictRequest struct {
		ParentID  uint   `json:"parentId"`
		Code      string `json:"code"`      // 字典编码 父ID为0时，code不能为空 且唯一
		DictValue string `json:"dictValue"` // 字典值 父Id 不为0时，dictValue不能为空
		DictLabel string `json:"dictLabel" validate:"required"`
		DictType  string `json:"dictType" validate:"required"`
		Sort      int    `json:"sort"`
		Remark    string `json:"remark"`
	}
	UpdateDictRequest struct {
		ID        uint   `json:"id"`
		Code      string `json:"code"`
		DictValue string `json:"dictValue" validate:"required"`
		DictLabel string `json:"dictLabel" validate:"required"`
		DictType  string `json:"dictType" validate:"required"`
		Sort      *int   `json:"sort,omitempty"`
		Remark    string `json:"remark"`
	}
	IdRequest struct {
		Id uint `json:"id"`
	}
	ListDictRequest struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Code     string `json:"code"`
		Label    string `json:"label"`
		ParentId uint   `json:"parentId"`
	}
	ListDictResponse struct {
		List  []Dict `json:"list"`
		Total int64  `json:"total"`
	}
	GetDictTreeByCodeRequest struct {
		Code []string `json:"code"`
	}

	Audit struct {
		ID            uint            `json:"id"`
		TraceId       string          `json:"traceId"`
		Operator      string          `json:"operator"`
		CreatedAt     time.Time       `json:"createdAt"`
		IsError       bool            `json:"isError"`
		ErrorMessage  string          `json:"errorMessage"`
		Duration      string          `json:"duration"`
		RequestMethod string          `json:"requestMethod"`
		RequestUrl    string          `json:"requestUrl"`
		RequestBody   json.RawMessage `json:"requestBody"`
		ResponseBody  json.RawMessage `json:"responseBody"`
		TenantId      uint            `json:"tenantId"`
	}
	ListAuditRequest struct {
		Page      int        `json:"page"`
		PageSize  int        `json:"pageSize"`
		TraceId   string     `json:"traceId"`
		Operator  string     `json:"operator"`
		StartTime *time.Time `json:"startTime"`
		EndTime   *time.Time `json:"endTime"`
		IsError   *bool      `json:"isError"`
		Duration  float64    `json:"duration"`
	}

	ListAuditResponse struct {
		List  []Audit `json:"list"`
		Total int64   `json:"total"`
	}
)

type (
	templateListRequest struct {
		Name         string `json:"name"`
		TemplateType string `json:"templateType"`
		Page         int    `json:"page"`
		PageSize     int    `json:"pageSize"`
	}
	templateListResult struct {
		Name          string             `json:"name"`
		BaseModel     string             `json:"baseModel"`
		Content       string             `json:"content"`
		Params        string             `json:"params"`
		TrainImage    string             `json:"trainImage"`
		Remark        string             `json:"remark"`
		BaseModelPath string             `json:"baseModelPath"`
		ScriptFile    string             `json:"scriptFile"`
		OutputDir     string             `json:"outputDir"`
		MaxTokens     int                `json:"maxTokens"`
		Lora          bool               `json:"lora"`
		Enabled       bool               `json:"enabled"`
		TemplateType  types.TemplateType `json:"templateType"`
		CreatedAt     time.Time          `json:"createdAt"`
		UpdatedAt     time.Time          `json:"updatedAt"`
	}

	templateCreateRequest struct {
		Name          string             `json:"name" validate:"required"`          //名称
		BaseModel     string             `json:"baseModel" validate:"required"`     //基本模型 来源models表
		Content       string             `json:"content" validate:"required"`       //脚本模版
		Params        string             `json:"params"`                            //模版所需要参数
		TrainImage    string             `json:"trainImage" validate:"required"`    //训练镜像
		Remark        string             `json:"remark"`                            //备注
		BaseModelPath string             `json:"baseModelPath" validate:"required"` //基础模型路径
		ScriptFile    string             `json:"scriptFile"`                        //脚本文件
		OutputDir     string             `json:"outputDir"`                         //输出目录
		Enabled       bool               `json:"enabled"`                           //是否可用
		GpuLabel      string             `json:"gpuLabel"`                          //GPU 标签
		ParallelNum   int                `json:"parallelNum"`                       //并行数量
		K8sCluster    string             `json:"k8sCluster"`                        //K8S集群
		Cpu           int                `json:"cpu"`                               //CPU 核数
		Memory        int                `json:"memory"`                            //内存大小 G
		TemplateType  types.TemplateType `json:"templateType"`
	}

	templateDeleteRequest struct {
		Name string `json:"name" validate:"required"` //名称
	}
)

type Endpoints struct {
	ListDictEndpoint       endpoint.Endpoint
	CreateDictEndpoint     endpoint.Endpoint
	UpdateDictEndpoint     endpoint.Endpoint
	DeleteDictEndpoint     endpoint.Endpoint
	GetDictTreeEndpoint    endpoint.Endpoint
	ListAuditEndpoint      endpoint.Endpoint
	ListTemplateEndpoint   endpoint.Endpoint
	CreateTemplateEndpoint endpoint.Endpoint
	UpdateTemplateEndpoint endpoint.Endpoint
	DeleteTemplateEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListDictEndpoint:       makeListDictEndpoint(s),
		CreateDictEndpoint:     makeCreateDictEndpoint(s),
		UpdateDictEndpoint:     makeUpdateDictEndpoint(s),
		DeleteDictEndpoint:     makeDeleteDictEndpoint(s),
		GetDictTreeEndpoint:    makeGetDictTreeEndpoint(s),
		ListAuditEndpoint:      makeListAuditEndpoint(s),
		ListTemplateEndpoint:   makeListTemplateEndpoint(s),
		CreateTemplateEndpoint: makeCreateTemplateEndpoint(s),
		UpdateTemplateEndpoint: makeUpdateTemplateEndpoint(s),
		DeleteTemplateEndpoint: makeDeleteTemplateEndpoint(s),
	}
	for _, m := range dmw["Sys"] {
		eps.ListDictEndpoint = m(eps.ListDictEndpoint)
		eps.CreateDictEndpoint = m(eps.CreateDictEndpoint)
		eps.UpdateDictEndpoint = m(eps.UpdateDictEndpoint)
		eps.DeleteDictEndpoint = m(eps.DeleteDictEndpoint)
		eps.GetDictTreeEndpoint = m(eps.GetDictTreeEndpoint)
		eps.ListAuditEndpoint = m(eps.ListAuditEndpoint)
		eps.ListTemplateEndpoint = m(eps.ListTemplateEndpoint)
		eps.CreateTemplateEndpoint = m(eps.CreateTemplateEndpoint)
		eps.UpdateTemplateEndpoint = m(eps.UpdateTemplateEndpoint)
		eps.DeleteTemplateEndpoint = m(eps.DeleteTemplateEndpoint)

	}
	return eps
}

func makeListDictEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListDictRequest)
		resp, err := s.ListDict(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeCreateDictEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateDictRequest)
		resp, err := s.CreateDict(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeUpdateDictEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateDictRequest)
		err = s.UpdateDict(ctx, req)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteDictEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IdRequest)
		err = s.DeleteDict(ctx, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeGetDictTreeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetDictTreeByCodeRequest)
		resp, err := s.DictTreeByCode(ctx, req.Code)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListAuditEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListAuditRequest)
		resp, err := s.ListAudit(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}

func makeListTemplateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(templateListRequest)
		res, total, err := s.TemplateList(ctx, req.Page, req.PageSize, req.Name, req.TemplateType)

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

func makeCreateTemplateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(templateCreateRequest)
		err = s.TemplateCreate(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}

func makeUpdateTemplateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(templateCreateRequest)
		err = s.TemplateUpdate(ctx, req)

		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteTemplateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(templateDeleteRequest)
		err = s.TemplateDelete(ctx, req.Name)

		return encode.Response{
			Error: err,
		}, err
	}
}
