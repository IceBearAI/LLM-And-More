package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// MetadataHttp 动态工具元数据
type MetadataHttp struct {
	// Url 请求地址
	Url string `json:"url"`
	// Method 请求方法
	Method string `json:"method"`
	// Headers 请求头
	Headers map[string]string `json:"headers"`
	// Body 请求体
	Body interface{} `json:"body"`
	// UserAgent 请求头
	UserAgent string `json:"userAgent"`
}

// Service describes a service that adds things together.
type Service interface {
	// Create 创建工具
	Create(ctx context.Context, tenantId uint, req createRequest) (err error)
	// Update 更新工具
	Update(ctx context.Context, tenantId uint, toolId string, req createRequest) (err error)
	// Delete 删除工具
	Delete(ctx context.Context, tenantId uint, toolId string) (err error)
	// Get 获取工具
	Get(ctx context.Context, tenantId uint, toolId string) (resp toolResult, err error)
	// List 列出工具
	List(ctx context.Context, tenantId uint, name string, page, pageSize int) (resp []toolResult, total int64, err error)
	// Test 测试工具
	Test(ctx context.Context, tenantId uint, toolId string, input string) (resp string, err error)
	// Assistants 获取所关联的助手
	Assistants(ctx context.Context, tenantId uint, toolId string) (resp []assistantResult, err error)
}

// Middleware describes a service middleware.
type Middleware func(Service) Service

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) Assistants(ctx context.Context, tenantId uint, toolId string) (resp []assistantResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tool, err := s.repository.Tools().Get(ctx, tenantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	assistants, err := s.repository.Tools().GetAssistants(ctx, tenantId, tool.ID)
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	for _, assistant := range assistants {
		resp = append(resp, assistantResult{
			AssistantId: assistant.UUID,
			Name:        assistant.Name,
			Description: assistant.Description,
		})
	}
	return
}

func (s *service) Test(ctx context.Context, tenantId uint, toolId string, input string) (resp string, err error) {
	tool, err := s.repository.Tools().Get(ctx, tenantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(s.logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	if tool.ToolType != types.ToolTypeTypeFunction {
		err = errors.New("工具不是函数类型")
		_ = level.Warn(s.logger).Log("msg", "工具不是函数类型", "tool", tool)
		return
	}
	if tool.Metadata == "" {
		err = errors.New("工具没有配置")
		_ = level.Warn(s.logger).Log("msg", "工具没有配置", "tool", tool)
		return
	}
	var metaData MetadataHttp
	if err = json.Unmarshal([]byte(tool.Metadata), &metaData); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("解析元数据失败: %s", tool.Metadata))
		return "", err
	}
	//tgt, err := url.Parse(metaData.Url)
	//if err != nil {
	//	err = errors.Wrap(err, "解析URL失败")
	//	return "", err
	//}
	//
	//resp, err = tool.Test(ctx, input)
	//if err != nil {
	//	err = errors.Wrap(err, "测试工具失败")
	//	_ = level.Warn(s.logger).Log("msg", "测试工具失败", "err", err)
	//	return
	//}
	return
}

func (s *service) Create(ctx context.Context, tenantId uint, req createRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tool, err := s.repository.Tools().Get(ctx, tenantId, req.Name)
	if err == nil {
		err = errors.New("工具已存在")
		_ = level.Warn(logger).Log("msg", "工具已存在", "tool", tool)
		return
	}
	tool = types.Tools{
		UUID:        fmt.Sprintf("tool-%s", uuid.New().String()),
		TenantId:    tenantId,
		Name:        req.Name,
		Description: req.Description,
		ToolType:    types.ToolType(req.ToolType),
		Metadata:    req.Metadata,
		Operator:    req.Operator,
		Remark:      req.Remark,
	}
	err = s.repository.Tools().Create(ctx, &tool)
	return
}

func (s *service) Update(ctx context.Context, tenantId uint, toolId string, req createRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tool, err := s.repository.Tools().Get(ctx, tenantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	tool.Name = req.Name
	tool.Description = req.Description
	tool.ToolType = types.ToolType(req.ToolType)
	tool.Metadata = req.Metadata
	tool.Operator = req.Operator
	tool.Remark = req.Remark
	err = s.repository.Tools().Update(ctx, &tool)
	return
}

func (s *service) Delete(ctx context.Context, tenantId uint, toolId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tool, err := s.repository.Tools().Get(ctx, tenantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	if err = s.repository.Tools().Delete(ctx, tenantId, toolId); err != nil {
		err = errors.Wrap(err, "删除工具失败")
		_ = level.Warn(logger).Log("msg", "删除工具失败", "err", err)
		return
	}
	if err = s.repository.Tools().ClearToolRelation(ctx, []uint{tool.ID}); err != nil {
		err = errors.Wrap(err, "清除工具关联关系失败")
		_ = level.Warn(logger).Log("msg", "清除工具关联关系失败", "err", err)
		return
	}
	return
}

func (s *service) Get(ctx context.Context, tenantId uint, toolId string) (resp toolResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tool, err := s.repository.Tools().Get(ctx, tenantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	resp = toolResult{
		ToolId: tool.UUID,
	}
	return
}

func (s *service) List(ctx context.Context, tenantId uint, name string, page, pageSize int) (resp []toolResult, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tools, total, err := s.repository.Tools().List(ctx, tenantId, name, page, pageSize, "Assistants")
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}
	for _, tool := range tools {
		assistants := make([]assistantResult, 0)
		for _, assistant := range tool.Assistants {
			assistants = append(assistants, assistantResult{
				AssistantId: assistant.UUID,
				Name:        assistant.Name,
				Description: assistant.Description,
			})
		}
		resp = append(resp, toolResult{
			ToolId:      tool.UUID,
			Name:        tool.Name,
			Description: tool.Description,
			ToolType:    string(tool.ToolType),
			Metadata:    tool.Metadata,
			Operator:    tool.Operator,
			Assistants:  assistants,
			Remark:      tool.Remark,
			UpdatedAt:   tool.UpdatedAt,
		})
	}
	return
}

// New returns a naive, stateless implementation of Service.
func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "service", "tools")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
