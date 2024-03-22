package datasets

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(Service) Service

// Service 数据集服务
type Service interface {
	// List 根据租户获取数据集列表
	List(ctx context.Context, tenantId uint, page, pageSize int, name string, preloads ...string) (datasets []types.Dataset, total int64, err error)
	// Create 创建数据集
	Create(ctx context.Context, dataset *types.Dataset) (err error)
	// Update 更新数据集
	Update(ctx context.Context, dataset *types.Dataset) (err error)
	// Delete 删除数据集
	Delete(ctx context.Context, id uint) (err error)
	// FindById 根据ID查找数据集
	FindById(ctx context.Context, id uint, preloads ...string) (dataset types.Dataset, err error)
	// FindByUUID 根据UUID查找数据集
	FindByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.Dataset, err error)
	// FindByUUIDAndTenantId 根据UUID和租户ID查找数据集
	FindByUUIDAndTenantId(ctx context.Context, uuid string, tenantId uint, preloads ...string) (dataset types.Dataset, err error)
	// CreateSample 创建数据集样本
	CreateSample(ctx context.Context, sample *types.DatasetSample) (err error)
	// UpdateSample 更新数据集样本
	UpdateSample(ctx context.Context, sample *types.DatasetSample) (err error)
	// DeleteSample 删除数据集样本
	DeleteSample(ctx context.Context, id uint) (err error)
	// SampleList 根据数据集ID获取样本列表
	SampleList(ctx context.Context, datasetId uint, page, pageSize int, title string, preloads ...string) (samples []types.DatasetSample, total int64, err error)
	// DeleteSampleByUUID 根据样本UUID删除样本
	DeleteSampleByUUID(ctx context.Context, uuid []string) (err error)
	// FindSampleByTitle 根据标题查找样本
	FindSampleByTitle(ctx context.Context, datasetId uint, title string, preloads ...string) (dataset types.DatasetSample, err error)
	// FindSampleByUUID 根据UUID查找样本
	FindSampleByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.DatasetSample, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindSampleByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.DatasetSample, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("uuid = ?", uuid).First(&dataset).Error
	return
}

func (s *service) FindSampleByTitle(ctx context.Context, datasetId uint, title string, preloads ...string) (dataset types.DatasetSample, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("dataset_id = ? AND title = ?", datasetId, title).First(&dataset).Error
	return
}

func (s *service) DeleteSampleByUUID(ctx context.Context, uuid []string) (err error) {
	return s.db.WithContext(ctx).Where("uuid in (?)", uuid).Delete(&types.DatasetSample{}).Error
}

func (s *service) SampleList(ctx context.Context, datasetId uint, page, pageSize int, title string, preloads ...string) (samples []types.DatasetSample, total int64, err error) {
	offset := (page - 1) * pageSize
	db := s.db.WithContext(ctx).Model(&types.DatasetSample{}).Where("dataset_id = ?", datasetId)
	if title != "" {
		db = db.Where("title like ?", "%"+title+"%")
	}
	err = db.Count(&total).Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&samples).Error
	return
}

func (s *service) List(ctx context.Context, tenantId uint, page, pageSize int, name string, preloads ...string) (datasets []types.Dataset, total int64, err error) {
	offset := (page - 1) * pageSize
	db := s.db.WithContext(ctx).Model(types.Dataset{}).Where("tenant_id = ?", tenantId)
	if name != "" {
		db = db.Where("`name` like ?", "%"+name+"%")
	}
	err = db.Count(&total).Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&datasets).Error
	return
}

func (s *service) Create(ctx context.Context, dataset *types.Dataset) (err error) {
	return s.db.WithContext(ctx).Create(dataset).Error
}

func (s *service) Update(ctx context.Context, dataset *types.Dataset) (err error) {
	return s.db.WithContext(ctx).Model(dataset).Where("id = ?", dataset.ID).Updates(dataset).Error
}

func (s *service) Delete(ctx context.Context, id uint) (err error) {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&types.Dataset{}).Error
}

func (s *service) FindById(ctx context.Context, id uint, preloads ...string) (dataset types.Dataset, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("id = ?", id).First(&dataset).Error
	return
}

func (s *service) FindByUUID(ctx context.Context, uuid string, preloads ...string) (dataset types.Dataset, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("uuid = ?", uuid).First(&dataset).Error
	return
}

func (s *service) FindByUUIDAndTenantId(ctx context.Context, uuid string, tenantId uint, preloads ...string) (dataset types.Dataset, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("uuid = ? and tenant_id = ?", uuid, tenantId).First(&dataset).Error
	return
}

func (s *service) CreateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	return s.db.WithContext(ctx).Create(sample).Error
}

func (s *service) UpdateSample(ctx context.Context, sample *types.DatasetSample) (err error) {
	return s.db.WithContext(ctx).Model(sample).Where("id = ?", sample.ID).Updates(sample).Error
}

func (s *service) DeleteSample(ctx context.Context, id uint) (err error) {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&types.DatasetSample{}).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
