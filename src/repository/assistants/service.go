package assistants

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Middleware func(service Service) Service

// Service 助手管理数据操作接口
type Service interface {
	// Create 创建助手
	Create(ctx context.Context, assistant *types.Assistants) (err error)
	// Update 更新助手
	Update(ctx context.Context, assistant *types.Assistants) (err error)
	// Delete 删除助手
	Delete(ctx context.Context, tenantId uint, assistantId string) (err error)
	// Get 获取助手
	Get(ctx context.Context, tenantId uint, assistantId string, preloads ...string) (assistant types.Assistants, err error)
	// List 列出助手
	List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (assistants []types.Assistants, total int64, err error)
	// AddTool 给助手添加工具
	AddTool(ctx context.Context, tenantId uint, assistantId string, toolId string) (err error)
	// RemoveTool 给助手移除工具
	RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error)
	// ListTool 列出助手的工具
	ListTool(ctx context.Context, tenantId uint, assistantId, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error)
	// FindByAssistantName 根据助手名称获取助手
	FindByAssistantName(ctx context.Context, tenantId uint, assistantName string) (assistant types.Assistants, err error)
	// ReplaceTools 替换助手的工具
	ReplaceTools(ctx context.Context, assistant *types.Assistants, tools []types.Tools) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) ReplaceTools(ctx context.Context, assistant *types.Assistants, tools []types.Tools) (err error) {
	err = s.db.WithContext(ctx).Model(assistant).Transaction(func(tx *gorm.DB) error {
		if err = tx.Association("Tools").Unscoped().Clear(); err != nil {
			return err
		}

		if err = tx.Model(&types.Assistants{}).Association("Tools").Append(tools); err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *service) FindByAssistantName(ctx context.Context, tenantId uint, assistantName string) (assistant types.Assistants, err error) {
	err = s.db.WithContext(ctx).Where("tenant_id = ? and name = ?", tenantId, assistantName).First(&assistant).Error
	return
}

func (s *service) Create(ctx context.Context, assistant *types.Assistants) (err error) {
	return s.db.WithContext(ctx).Create(assistant).Error
}

func (s *service) Update(ctx context.Context, assistant *types.Assistants) (err error) {
	return s.db.WithContext(ctx).Where("id = ?", assistant.ID).Updates(assistant).Error
}

func (s *service) Delete(ctx context.Context, tenantId uint, assistantId string) (err error) {
	db := s.db.WithContext(ctx)
	assistant := types.Assistants{}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, assistantId).First(&assistant).Error
	if err != nil {
		return
	}
	err = db.Select("Tools").Delete(&assistant).Error
	return
}

func (s *service) Get(ctx context.Context, tenantId uint, assistantId string, preloads ...string) (assistant types.Assistants, err error) {
	db := s.db.WithContext(ctx)
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, assistantId).First(&assistant).Error
	return
}

func (s *service) List(ctx context.Context, tenantId uint, name string, pageNum, pageSize int, preloads ...string) (assistants []types.Assistants, total int64, err error) {
	db := s.db.WithContext(ctx).Model(&types.Assistants{})
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	limit, offset := page.Limit(pageNum, pageSize)
	db = db.Where("tenant_id = ?", tenantId)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if err = db.Order("id DESC").Count(&total).Offset(offset).Limit(limit).Find(&assistants).Error; err != nil {
		return
	}
	return
}

func (s *service) AddTool(ctx context.Context, tenantId uint, assistantId string, toolId string) (err error) {
	db := s.db.WithContext(ctx)
	assistant := types.Assistants{}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, assistantId).First(&assistant).Error
	if err != nil {
		return
	}
	tool := types.Tools{}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, toolId).First(&tool).Error
	if err != nil {
		return
	}
	// 判断是否已存在
	var total int64
	if err = db.Model(&types.AssistantToolAssociations{}).
		Where("assistant_id = ? and tool_id = ?", assistant.ID, tool.ID).Count(&total).Error; err != nil {
		return
	}
	if total > 0 {
		return errors.New("工具已存在，无需重复添加")
	}
	assistantTool := types.AssistantToolAssociations{
		AssistantId: assistant.ID,
		ToolId:      tool.ID,
	}
	err = db.Create(&assistantTool).Error
	return
}

func (s *service) RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error) {
	db := s.db.WithContext(ctx)
	assistant := types.Assistants{}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, assistantId).First(&assistant).Error
	if err != nil {
		return
	}
	tool := types.Tools{}
	err = db.Where("tenant_id = ? and uuid = ?", tenantId, toolId).First(&tool).Error
	if err != nil {
		return
	}
	err = db.Where("assistant_id = ? and tool_id = ?", assistant.ID, tool.ID).Delete(&types.AssistantToolAssociations{}).Error
	return
}

func (s *service) ListTool(ctx context.Context, tenantId uint, assistantId, name string, pageNum, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	// TODO: 这里会有bug 如果有name条件分页会有问题，应该把assistant_tool_associations所有的tool_id全部查出来
	db := s.db.WithContext(ctx)
	assistant := types.Assistants{}
	if err = db.Where("tenant_id = ? and uuid = ?", tenantId, assistantId).First(&assistant).Error; err != nil {
		return
	}
	limit, offset := page.Limit(pageNum, pageSize)
	var toolIDs []uint
	if err = db.Table("assistant_tool_associations").Where("assistant_id = ?", assistant.ID).
		Order("tool_id DESC").
		Count(&total).Offset(offset).Limit(limit).Pluck("tool_id", &toolIDs).Error; err != nil {
		return
	}
	db = db.Where("id in ?", toolIDs)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if err = db.Order("id DESC").Find(&tools).Error; err != nil {
		return nil, 0, err
	}
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
