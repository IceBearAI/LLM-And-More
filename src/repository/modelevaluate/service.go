package modelevaluate

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
	"strings"
)

type Middleware func(Service) Service

type Service interface {
	// Save 创建
	Save(ctx context.Context, data *types.ModelEvaluate) (err error)
	// ListModelEvaluate 列表
	ListModelEvaluate(ctx context.Context, page, pageSize int, modelId uint, status string, evalTargetType string) (res []types.ModelEvaluate, total int64, err error)
	// GetById 获取详情
	GetById(ctx context.Context, id uint) (res types.ModelEvaluate, err error)
	// DeleteById 删除
	DeleteById(ctx context.Context, id uint) (err error)
	// FindFiveGraphLastByModelId 获取最新的五维图信息
	FindFiveGraphLastByModelId(ctx context.Context, modelId, evaluateId uint) (res types.ModelEvaluate, err error)
	// CountEvaluate 数量统计
	CountEvaluate(ctx context.Context, evalTargetType string, status []string) (res int64, err error)
	// GetByUuid 获取详情
	GetByUuid(ctx context.Context, uuid string) (res types.ModelEvaluate, err error)
	// IsExistFiveByModelId 判断模型是否五维图数据
	IsExistFiveByModelId(ctx context.Context, modelId uint) (res bool, err error)
	// 获取微调模型损失率
}

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, data *types.ModelEvaluate) (err error) {
	return s.db.WithContext(ctx).Save(data).Error
}

func (s *service) ListModelEvaluate(ctx context.Context, page, pageSize int, modelId uint, status string, evalTargetType string) (res []types.ModelEvaluate, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.ModelEvaluate{})

	if modelId > 0 {
		query = query.Where("model_id = ?", modelId)
	}

	if !strings.EqualFold(status, "") {
		query = query.Where("status = ?", status)
	}

	if !strings.EqualFold(evalTargetType, "") {
		query = query.Where("eval_target_type = ?", evalTargetType)
	}

	err = query.Count(&total).Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error

	return
}

func (s *service) GetById(ctx context.Context, id uint) (res types.ModelEvaluate, err error) {
	err = s.db.WithContext(ctx).Model(&types.ModelEvaluate{}).Preload("Models").Where("id = ?", id).First(&res).Error
	return
}

func (s *service) DeleteById(ctx context.Context, id uint) (err error) {
	return s.db.WithContext(ctx).Model(&types.ModelEvaluate{}).Where("id = ?", id).Delete(&types.ModelEvaluate{}).Error
}

func (s *service) FindFiveGraphLastByModelId(ctx context.Context, modelId, evaluateId uint) (res types.ModelEvaluate, err error) {
	query := s.db.WithContext(ctx).Model(&types.ModelEvaluate{}).Preload("Models")

	query = query.Where("model_id = ?", modelId)

	if evaluateId > 0 {
		query = query.Where("id = ?", evaluateId)
	}

	err = query.Order("id DESC").First(&res).Error

	return
}

func (s *service) CountEvaluate(ctx context.Context, evalTargetType string, status []string) (res int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.ModelEvaluate{})

	if !strings.EqualFold(evalTargetType, "") {
		query = query.Where("eval_target_type = ?", evalTargetType)
	}

	if len(status) > 0 {
		query = query.Where("status in ?", status)
	}

	err = query.Count(&res).Error

	return
}

func (s *service) GetByUuid(ctx context.Context, uuid string) (res types.ModelEvaluate, err error) {
	err = s.db.WithContext(ctx).Model(&types.ModelEvaluate{}).Preload("Models").Where("uuid = ?", uuid).First(&res).Error
	return
}

func (s *service) IsExistFiveByModelId(ctx context.Context, modelId uint) (res bool, err error) {
	var count int64
	err = s.db.WithContext(ctx).Model(&types.ModelEvaluate{}).Where("model_id = ?", modelId).Where("eval_target_type = ?", string(types.EvaluateTargetTypeFive)).Where("status = ?", string(types.EvaluateStatusSuccess)).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
