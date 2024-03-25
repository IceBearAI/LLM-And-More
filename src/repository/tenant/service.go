package tenant

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Middleware func(service Service) Service

type Service interface {
	// FindTenant 通过ID获取租户信息
	FindTenant(ctx context.Context, id uint, preloads ...string) (res types.Tenants, err error)
	// AddModel 添加模型
	AddModel(ctx context.Context, id uint, models ...*types.Models) (err error)
}

type service struct {
	db *gorm.DB
}

func (s service) AddModel(ctx context.Context, id uint, models ...*types.Models) (err error) {
	tenant, err := s.FindTenant(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "find tenant error")
		return
	}
	err = s.db.WithContext(ctx).Model(&tenant).Association("Models").Append(models)
	return
}

func (s service) FindTenant(ctx context.Context, id uint, preloads ...string) (res types.Tenants, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("id = ?", id).First(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
