package sys

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Middleware defines a middleware for service
type Middleware func(Service) Service

type ListDictRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Code     string `json:"code"`
	Label    string `json:"label"`
	ParentId uint   `json:"parentId"`
}

type ListAuditRequest struct {
	Page      int        `json:"page"`
	PageSize  int        `json:"pageSize"`
	TraceId   string     `json:"traceId"`
	Operator  string     `json:"operator"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	IsError   *bool      `json:"isError"`
	Duration  float64    `json:"duration"`
}

type Service interface {
	// ListDict 字典分页列表
	ListDict(ctx context.Context, request ListDictRequest) (items []types.SysDict, total int64, err error)
	// CreateDict 创建字典
	CreateDict(ctx context.Context, data *types.SysDict) (err error)
	// GetDict 获取字典
	GetDict(ctx context.Context, id uint) (res types.SysDict, err error)
	// UpdateDict 更新字典
	UpdateDict(ctx context.Context, data *types.SysDict, checkChildrenCode bool) (err error)
	// DeleteDict 删除字典
	DeleteDict(ctx context.Context, id uint) (err error)
	// GetDictByCode 根据code获取字典
	GetDictByCode(ctx context.Context, code string) (res types.SysDict, err error)
	// FindDictTreeByParentId 根据parentId获取字典树
	FindDictTreeByParentId(ctx context.Context, parentId uint, parentDictType ...string) (res []types.SysDict, err error)
	// GetDictByDictValue 根据字典值获取字典
	GetDictByDictValue(ctx context.Context, parentId uint, dictValue string) (res types.SysDict, err error)
	// FindDictTreeByCode 根据code获取字典树
	FindDictTreeByCode(ctx context.Context, code []string) (res []types.SysDict, err error)
	// CreateAudit 创建审计
	CreateAudit(ctx context.Context, data *types.SysAudit) (err error)
	// ListAudit 审计分页列表
	ListAudit(ctx context.Context, request ListAuditRequest) (items []types.SysAudit, total int64, err error)
	//ListFineTuningTemplate 模板列表
	ListFineTuningTemplate(ctx context.Context, page, pageSize int, name, templateType string) (res []types.FineTuningTemplate, total int64, err error)
	//SaveFineTuningTemplate 模板创建/更新
	SaveFineTuningTemplate(ctx context.Context, data *types.FineTuningTemplate) (err error)
	//DeleteFineTuningTemplate 模板删除
	DeleteFineTuningTemplate(ctx context.Context, name string) (err error)
	//IsExistTuningTemplate 模板是否存在
	IsExistTuningTemplate(ctx context.Context, name string) (isExist bool, err error)
	//GetFineTuningTemplate 模板详情
	GetFineTuningTemplate(ctx context.Context, name string) (res types.FineTuningTemplate, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) ListFineTuningTemplate(ctx context.Context, pageNum, pageSize int, name, templateType string) (res []types.FineTuningTemplate, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.FineTuningTemplate{})

	if !strings.EqualFold(templateType, "") {
		query = query.Where("template_type = ?", templateType)
	}

	if !strings.EqualFold(name, "") {
		query = query.Where("name like ?", "%"+name+"%")
	}
	limit, offset := page.Limit(pageNum, pageSize)
	err = query.Count(&total).Order("updated_at DESC").Offset(offset).Limit(limit).Find(&res).Error

	return
}

func (s *service) SaveFineTuningTemplate(ctx context.Context, data *types.FineTuningTemplate) (err error) {
	return s.db.WithContext(ctx).Save(data).Error
}

func (s *service) DeleteFineTuningTemplate(ctx context.Context, name string) (err error) {
	return s.db.WithContext(ctx).Model(&types.FineTuningTemplate{}).Where("name = ?", name).Delete(&types.FineTuningTemplate{}).Error
}

func (s *service) IsExistTuningTemplate(ctx context.Context, name string) (isExist bool, err error) {
	var c int64
	if err := s.db.WithContext(ctx).Model(&types.FineTuningTemplate{}).Where("name = ?", name).Count(&c).Error; err != nil {
		return false, err
	} else {
		if c > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) GetFineTuningTemplate(ctx context.Context, name string) (res types.FineTuningTemplate, err error) {
	err = s.db.WithContext(ctx).Model(&types.FineTuningTemplate{}).Where("name = ?", name).First(&res).Error
	return
}

func (s *service) CreateAudit(ctx context.Context, data *types.SysAudit) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) ListAudit(ctx context.Context, request ListAuditRequest) (items []types.SysAudit, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.SysAudit{})
	if request.TraceId != "" {
		query = query.Where("trace_id = ?", request.TraceId)
	}
	if request.Operator != "" {
		query = query.Where("operator = ?", request.Operator)
	}
	if request.StartTime != nil {
		query = query.Where("created_at >= ?", request.StartTime)
	}
	if request.EndTime != nil {
		query = query.Where("created_at <= ?", request.EndTime)
	}
	if request.IsError != nil {
		query = query.Where("is_error = ?", request.IsError)
	}
	if request.Duration != 0 {
		query = query.Where("duration >= ?", request.Duration)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	return
}

func (s *service) FindDictTreeByCode(ctx context.Context, code []string) (res []types.SysDict, err error) {
	var settings []types.SysDict
	err = s.db.WithContext(ctx).Where("code in ? and parent_id = ?", code, 0).Find(&settings).Error
	if err != nil {
		return
	}
	for i, setting := range settings {
		children, err := s.FindDictTreeByParentId(ctx, setting.ID, setting.DictType)
		if err != nil {
			return nil, err
		}
		settings[i].Children = children
	}
	return settings, nil
}

func (s *service) GetDictByDictValue(ctx context.Context, parentId uint, dictValue string) (res types.SysDict, err error) {
	err = s.db.WithContext(ctx).Where("parent_id = ? and dict_value = ?", parentId, dictValue).First(&res).Error
	return
}

func (s *service) ListDict(ctx context.Context, request ListDictRequest) (items []types.SysDict, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.SysDict{})
	if request.Code != "" {
		query = query.Where("code LIKE ?", "%"+request.Code+"%")
	}
	if request.Label != "" {
		query = query.Where("dict_label LIKE ?", "%"+request.Label+"%")
	}
	if request.ParentId > 0 {
		query = query.Where("id = ?", request.ParentId)
	} else {
		query = query.Where("parent_id = ?", 0)
	}

	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("sort DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	if request.ParentId > 0 {
		for i, item := range items {
			children, err := s.FindDictTreeByParentId(ctx, item.ID, item.DictType)
			if err != nil {
				return nil, 0, err
			}
			items[i].Children = children
		}
	}
	return
}

func (s *service) CreateDict(ctx context.Context, data *types.SysDict) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) GetDict(ctx context.Context, id uint) (res types.SysDict, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) UpdateDict(ctx context.Context, data *types.SysDict, checkChildrenCode bool) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		if checkChildrenCode {
			ids, err2 := s.getDictChildrenIds(ctx, data.ID)
			if err2 != nil {
				return err2
			}
			if err := tx.Model(&types.SysDict{}).Where("id in ?", ids).Update("code", data.Code).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *service) DeleteDict(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&types.SysDict{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&types.SysDict{}).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *service) GetDictByCode(ctx context.Context, code string) (res types.SysDict, err error) {
	err = s.db.WithContext(ctx).Where("code = ? and parent_id = ?", code, 0).First(&res).Error
	return
}

func (s *service) FindDictTreeByParentId(ctx context.Context, parentId uint, parentDictType ...string) (res []types.SysDict, err error) {
	var settings []types.SysDict
	err = s.db.WithContext(ctx).Where("parent_id = ?", parentId).Order("sort desc").Find(&settings).Error
	if err != nil {
		return
	}
	for i, setting := range settings {
		// 当 parentId 不为 0 时，设置 ParentDictType
		if parentId != 0 && len(parentDictType) > 0 {
			settings[i].ParentDictType = parentDictType[0]
		}

		children, err := s.FindDictTreeByParentId(ctx, setting.ID, setting.DictType)
		if err != nil {
			return nil, err
		}
		settings[i].Children = children
	}
	return settings, nil
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

func (s *service) getDictChildrenIds(ctx context.Context, parentId uint) (ids []uint, err error) {
	var settings []types.SysDict
	err = s.db.WithContext(ctx).Select("id").Where("parent_id = ?", parentId).Find(&settings).Error
	if err != nil {
		return
	}
	for _, setting := range settings {
		ids = append(ids, setting.ID)
		childrenIds, err := s.getDictChildrenIds(ctx, setting.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, childrenIds...)
	}
	return ids, nil
}
