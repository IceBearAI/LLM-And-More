package sys

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/sys"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

type Middleware func(Service) Service

type Service interface {
	// ListDict 字典分页列表
	ListDict(ctx context.Context, request ListDictRequest) (resp ListDictResponse, err error)
	// CreateDict 创建字典
	CreateDict(ctx context.Context, data CreateDictRequest) (resp Dict, err error)
	// DictTreeByCode 针对前端的字典树
	DictTreeByCode(ctx context.Context, codes []string) (resp []Dict, err error)
	// UpdateDict 更新字典
	UpdateDict(ctx context.Context, data UpdateDictRequest) (err error)
	// DeleteDict 删除字典
	DeleteDict(ctx context.Context, id uint) (err error)
	// ListAudit 审计分页列表
	ListAudit(ctx context.Context, request ListAuditRequest) (resp ListAuditResponse, err error)
	// TemplateList 模板列表
	TemplateList(ctx context.Context, page, pageSize int, name, templateType string) (res []templateListResult, total int64, err error)
	// TemplateCreate 模板创建
	TemplateCreate(ctx context.Context, req templateCreateRequest) (err error)
	// TemplateUpdate 模板更新
	TemplateUpdate(ctx context.Context, req templateCreateRequest) (err error)
	// TemplateDelete 模板删除
	TemplateDelete(ctx context.Context, name string) (err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  services.Service
}

func (s *service) TemplateList(ctx context.Context, page, pageSize int, name, templateType string) (res []templateListResult, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "TemplateList")
	list, total, err := s.store.Sys().ListFineTuningTemplate(ctx, page, pageSize, name, templateType)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "ListFineTuningTemplate", "err", err.Error())
		return
	}

	for _, v := range list {
		dd := templateListResult{
			Name:          v.Name,
			BaseModel:     v.BaseModel,
			Content:       v.Content,
			Params:        v.Params,
			TrainImage:    v.TrainImage,
			Remark:        v.Remark,
			BaseModelPath: v.BaseModelPath,
			ScriptFile:    v.ScriptFile,
			OutputDir:     v.OutputDir,
			MaxTokens:     v.MaxTokens,
			Lora:          v.Lora,
			Enabled:       v.Enabled,
			TemplateType:  v.TemplateType,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		}

		res = append(res, dd)
	}

	return
}

func (s *service) TemplateCreate(ctx context.Context, req templateCreateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "TemplateCreate")
	isExist, err := s.store.Sys().IsExistTuningTemplate(ctx, req.Name)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "IsExistTuningTemplate", "err", err.Error())
		return
	}

	if isExist == true {
		return fmt.Errorf("该模板名称已存在，请修改模板名称")
	}

	var maxTokens int

	bashModel, err := s.store.Model().GetModelByModelName(ctx, req.BaseModel)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Model", "GetModelByModelName", "err", err.Error())
		return err
	}
	maxTokens = bashModel.MaxTokens

	data := &types.FineTuningTemplate{
		Name:          req.Name,
		BaseModel:     req.BaseModel,
		Content:       req.Content,
		Params:        req.Params,
		TrainImage:    req.TrainImage,
		Remark:        req.Remark,
		BaseModelPath: req.BaseModelPath,
		ScriptFile:    req.ScriptFile,
		OutputDir:     req.OutputDir,
		MaxTokens:     maxTokens,
		Enabled:       req.Enabled,
		TemplateType:  req.TemplateType,
		GpuLabel:      req.GpuLabel,
		ParallelNum:   req.ParallelNum,
		K8sCluster:    req.K8sCluster,
		Cpu:           req.Cpu,
		Memory:        req.Memory,
	}

	err = s.store.Sys().SaveFineTuningTemplate(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "SaveFineTuningTemplate", "err", err.Error())
		return
	}

	return nil
}

func (s *service) TemplateUpdate(ctx context.Context, req templateCreateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "TemplateUpdate")
	templateInfo, err := s.store.Sys().GetFineTuningTemplate(ctx, req.Name)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "GetFineTuningTemplate", "err", err.Error())
		return
	}

	var maxTokens int
	bashModel, err := s.store.Model().GetModelByModelName(ctx, req.BaseModel)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Model", "GetModelByModelName", "err", err.Error())
		return err
	}
	maxTokens = bashModel.MaxTokens

	templateInfo.BaseModel = req.BaseModel
	templateInfo.Content = req.Content
	templateInfo.Params = req.Params
	templateInfo.TrainImage = req.TrainImage
	templateInfo.Remark = req.Remark
	templateInfo.BaseModelPath = req.BaseModelPath
	templateInfo.ScriptFile = req.ScriptFile
	templateInfo.OutputDir = req.OutputDir
	templateInfo.MaxTokens = maxTokens
	templateInfo.Enabled = req.Enabled
	templateInfo.GpuLabel = req.GpuLabel
	templateInfo.ParallelNum = req.ParallelNum
	templateInfo.K8sCluster = req.K8sCluster
	templateInfo.Cpu = req.Cpu
	templateInfo.Memory = req.Memory
	templateInfo.TemplateType = req.TemplateType

	err = s.store.Sys().SaveFineTuningTemplate(ctx, &templateInfo)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "SaveFineTuningTemplate", "err", err.Error())
		return
	}

	return nil

}

func (s *service) TemplateDelete(ctx context.Context, name string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "TemplateUpdate")
	templateInfo, err := s.store.Sys().GetFineTuningTemplate(ctx, name)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "GetFineTuningTemplate", "err", err.Error())
		return
	}

	err = s.store.Sys().DeleteFineTuningTemplate(ctx, templateInfo.Name)
	if err != nil {
		_ = level.Error(logger).Log("s.store.Sys", "DeleteFineTuningTemplate", "err", err.Error())
		return
	}

	return nil
}

func (s *service) ListAudit(ctx context.Context, request ListAuditRequest) (resp ListAuditResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListAudit")
	var items []types.SysAudit
	req := sys.ListAuditRequest{
		Page:      request.Page,
		PageSize:  request.PageSize,
		TraceId:   request.TraceId,
		Operator:  request.Operator,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		IsError:   request.IsError,
		Duration:  request.Duration,
	}
	resp.List = make([]Audit, 0)
	items, resp.Total, err = s.store.Sys().ListAudit(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("store.Sys", "ListAudit", "err", err.Error())
		return
	}
	for _, item := range items {
		resp.List = append(resp.List, convertAudit(item))
	}
	return
}

func (s *service) ListDict(ctx context.Context, request ListDictRequest) (resp ListDictResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListDict")
	var items []types.SysDict
	req := sys.ListDictRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
		Code:     request.Code,
		Label:    request.Label,
		ParentId: request.ParentId,
	}
	resp.List = make([]Dict, 0)
	items, resp.Total, err = s.store.Sys().ListDict(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("store.Sys", "ListDict", "err", err.Error())
		return
	}
	for _, item := range items {
		resp.List = append(resp.List, convertDict(item))
	}
	return
}

func (s *service) CreateDict(ctx context.Context, data CreateDictRequest) (resp Dict, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateDict")
	req := types.SysDict{
		ParentID:  data.ParentID,
		Code:      data.Code,
		DictValue: data.DictValue,
		DictLabel: data.DictLabel,
		DictType:  data.DictType,
		Sort:      data.Sort,
		Remark:    data.Remark,
	}
	if data.ParentID == 0 {
		res, err := s.store.Sys().GetDictByCode(ctx, req.Code)
		if err == nil && res.ID > 0 {
			_ = level.Error(logger).Log("store.Sys", "GetDictByCode", "err", "字典编码已存在", "code", req.Code)
			return Dict{}, encode.InvalidParams.Wrap(errors.Errorf("字典编码%s已存在", req.Code))
		}
		req.DictValue = req.Code
	}

	var parent types.SysDict
	if data.ParentID > 0 {
		parent, err = s.store.Sys().GetDict(ctx, data.ParentID)
		if err != nil {
			_ = level.Error(logger).Log("store.Sys", "GetDict", "err", err.Error(), "id", data.ParentID)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return Dict{}, encode.InvalidParams.Wrap(errors.Errorf("父级字典不存在"))
			}
			return Dict{}, encode.ErrSystem.Wrap(errors.New("获取父级字典失败"))
		}
		req.Code = parent.Code
		res, err := s.store.Sys().GetDictByDictValue(ctx, parent.ID, req.DictValue)
		if err == nil && res.ID > 0 {
			_ = level.Error(logger).Log("store.Sys", "GetDictByDictValue", "err", "字典值已存在", "code", req.Code, "dictValue", req.DictValue)
			return Dict{}, encode.InvalidParams.Wrap(errors.Errorf("字典值%s已存在", req.DictValue))
		}
	}

	err = s.store.Sys().CreateDict(ctx, &req)
	if err != nil {
		_ = level.Error(logger).Log("store.Sys", "CreateDict", "err", err.Error())
		return
	}
	resp = convertDict(req)
	return
}

func (s *service) DictTreeByCode(ctx context.Context, code []string) (resp []Dict, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DictTreeByCode")
	resp = make([]Dict, 0)
	var items []types.SysDict
	if len(code) == 0 {
		items, err = s.store.Sys().FindDictTreeByParentId(ctx, 0)
		if err != nil {
			_ = level.Error(logger).Log("store.Sys", "FindDictTreeByParentId", "err", err.Error())
			return nil, err
		}
	} else {
		items, err = s.store.Sys().FindDictTreeByCode(ctx, code)
		if err != nil {
			_ = level.Error(logger).Log("store.Sys", "FindDictTreeByCode", "err", err.Error())
			return nil, err
		}
	}
	for _, item := range items {
		resp = append(resp, convertDict(item, true))
	}
	return
}

func (s *service) UpdateDict(ctx context.Context, req UpdateDictRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateDict")
	res, err := s.store.Sys().GetDict(ctx, req.ID)
	if err != nil {
		_ = level.Error(logger).Log("store.Sys", "GetDict", "err", err.Error(), "id", req.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return encode.InvalidParams.Wrap(errors.Errorf("字典不存在"))
		}
		return encode.ErrSystem.Wrap(errors.New("获取字典失败"))
	}
	updateChildrenCode := false
	// 只有顶级字典才能修改编码
	if req.Code != "" && res.ParentID == 0 {
		dict, err := s.store.Sys().GetDictByCode(ctx, req.Code)
		if err == nil && dict.ID > 0 && dict.ID != res.ID {
			_ = level.Error(logger).Log("store.Sys", "GetDictByCode", "err", "字典编码已存在", "code", req.Code)
			return encode.InvalidParams.Wrap(errors.Errorf("字典编码%s已存在", req.Code))
		}
		updateChildrenCode = true
		res.Code = req.Code
	}

	if req.DictLabel != "" {
		res.DictLabel = req.DictLabel
	}

	if req.DictValue != "" {
		if res.ParentID > 0 {
			dict, err := s.store.Sys().GetDictByDictValue(ctx, res.ParentID, req.DictValue)
			if err == nil && dict.ID > 0 && dict.ID != req.ID {
				_ = level.Error(logger).Log("store.Sys", "GetDictByDictValue", "err", "字典值已存在", "code", req.Code, "dictValue", req.DictValue)
				return encode.InvalidParams.Wrap(errors.Errorf("字典值%s已存在", req.DictValue))
			}
		}
		res.DictValue = req.DictValue
	}

	if req.DictType != "" {
		res.DictType = req.DictType
	}

	if req.Sort != nil {
		res.Sort = *req.Sort
	}

	if res.Remark != req.Remark {
		res.Remark = req.Remark
	}

	err = s.store.Sys().UpdateDict(ctx, &res, updateChildrenCode)
	if err != nil {
		_ = level.Error(s.logger).Log("store.Sys", "UpdateDict", "err", err.Error())
		return encode.ErrSystem.Wrap(errors.New("更新字典失败"))
	}
	return
}

func (s *service) DeleteDict(ctx context.Context, id uint) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteDict")
	err = s.store.Sys().DeleteDict(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("store.Sys", "DeleteDict", "err", err.Error())
		return
	}
	return
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc services.Service) Service {
	return &service{
		logger:  log.With(logger, "service", "models"),
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
	}
}

func convertDict(d types.SysDict, checkDictValue ...bool) Dict {
	t := Dict{
		ID:        d.ID,
		ParentID:  d.ParentID,
		Code:      d.Code,
		DictLabel: d.DictLabel,
		DictType:  d.DictType,
		DictValue: d.DictValue,
		Sort:      d.Sort,
		Remark:    d.Remark,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Children:  make([]Dict, 0),
	}
	if len(checkDictValue) > 0 && checkDictValue[0] {
		switch d.ParentDictType {
		case types.DictTypeString.String():
			t.DictValue = d.DictValue
		case types.DictTypeInt.String():
			value, _ := strconv.Atoi(d.DictValue)
			t.DictValue = value
		case types.DictTypeBool.String():
			value, _ := strconv.ParseBool(d.DictValue)
			t.DictValue = value
		default:
			t.DictValue = d.DictValue
		}
	}
	for _, child := range d.Children {
		t.Children = append(t.Children, convertDict(child, checkDictValue...))
	}
	return t
}

func convertAudit(a types.SysAudit) Audit {
	return Audit{
		ID:            a.ID,
		TraceId:       a.TraceID,
		Operator:      a.Operator,
		CreatedAt:     a.CreatedAt,
		IsError:       a.IsError,
		ErrorMessage:  a.ErrorMessage,
		RequestMethod: a.RequestMethod,
		RequestUrl:    a.RequestUrl,
		Duration:      fmt.Sprintf("%.3fs", a.Duration),
		RequestBody:   json.RawMessage(a.RequestBody),
		ResponseBody:  json.RawMessage(a.ResponseBody),
		TenantId:      a.TenantID,
	}
}
