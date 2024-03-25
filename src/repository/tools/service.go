package tools

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(Service) Service

// Service 工具
type Service interface {
	// Create 创建工具
	Create(ctx context.Context, tool *types.Tools) (err error)
	// Update 更新工具
	Update(ctx context.Context, tool *types.Tools) (err error)
	// Delete 删除工具
	Delete(ctx context.Context, tenantId uint, toolId string) (err error)
	// Get 获取工具
	Get(ctx context.Context, tenantId uint, toolId string) (tool types.Tools, err error)
	// List 列出工具
	List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error)
	// GetByIds 根据Ids获取工具
	GetByIds(ctx context.Context, tenantId uint, toolIds []string) (tools []types.Tools, err error)
	// ClearToolRelation 清除关联关系
	ClearToolRelation(ctx context.Context, toolIds []uint) (err error)
	// GetAssistants 获取所关联的助手
	GetAssistants(ctx context.Context, tenantId, toolId uint) (resp []types.Assistants, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) GetAssistants(ctx context.Context, tenantId, toolId uint) (resp []types.Assistants, err error) {
	err = s.db.WithContext(ctx).Model(&types.Assistants{}).
		Joins("join assistant_tool_associations on assistant_tool_associations.assistant_id = assistants.id").
		Where("assistant_tool_associations.tool_id = ?", toolId).Find(&resp).Error
	return
}

func (s *service) ClearToolRelation(ctx context.Context, toolIds []uint) (err error) {
	return s.db.WithContext(ctx).
		Where("tool_id in (?)", toolIds).Delete(&types.AssistantToolAssociations{}).Error
}

func (s *service) GetByIds(ctx context.Context, tenantId uint, toolIds []string) (tools []types.Tools, err error) {
	err = s.db.WithContext(ctx).Where("tenant_id = ? and uuid in ?", tenantId, toolIds).Find(&tools).Error
	return
}

func (s *service) Create(ctx context.Context, tool *types.Tools) (err error) {
	return s.db.WithContext(ctx).Create(tool).Error
}

func (s *service) Update(ctx context.Context, tool *types.Tools) (err error) {
	return s.db.WithContext(ctx).Updates(tool).Error
}

func (s *service) Delete(ctx context.Context, tenantId uint, toolId string) (err error) {
	return s.db.WithContext(ctx).Where("tenant_id = ? and uuid = ?", tenantId, toolId).Delete(&types.Tools{}).Error
}

func (s *service) Get(ctx context.Context, tenantId uint, toolId string) (tool types.Tools, err error) {
	err = s.db.WithContext(ctx).Where("tenant_id = ? and uuid = ?", tenantId, toolId).First(&tool).Error
	return
}

func (s *service) List(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (tools []types.Tools, total int64, err error) {
	db := s.db.WithContext(ctx).Where("tenant_id = ?", tenantId)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	err = db.Model(&types.Tools{}).Count(&total).Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tools).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
