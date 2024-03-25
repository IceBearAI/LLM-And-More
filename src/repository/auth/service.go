package auth

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(Service) Service

type Service interface {
	// GetTenantByUuid 根据uuid获取租户信息
	GetTenantByUuid(ctx context.Context, uuid string, preload ...string) (res types.Tenants, err error)
	// GetAccountByEmail 根据email获取账号信息
	GetAccountByEmail(ctx context.Context, email string, preload ...string) (res types.Accounts, err error)
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, data *types.Tenants) (err error)
	// UpdateTenant 更新租户
	UpdateTenant(ctx context.Context, data *types.Tenants) (err error)
	// DeleteTenant 删除租户
	DeleteTenant(ctx context.Context, id uint) (err error)
	// ListTenants 租户列表
	ListTenants(ctx context.Context, request ListTenantRequest) (res []types.Tenants, total int64, err error)
	// CreateAccount 创建账号
	CreateAccount(ctx context.Context, data *types.Accounts, tenantId uint) (err error)
	// ListAccount 获取账号列表
	ListAccount(ctx context.Context, request ListAccountRequest) (res []types.Accounts, total int64, err error)
	// UpdateAccount 更新账号
	UpdateAccount(ctx context.Context, data *types.Accounts) (err error)
	// GetAccountById 根据id获取账号信息
	GetAccountById(ctx context.Context, id uint) (res types.Accounts, err error)
	// DeleteAccount 删除账号
	DeleteAccount(ctx context.Context, id uint) (err error)
}

type ListTenantRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`
}

type ListAccountRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	IsLdap   *bool  `json:"isLdap"`
	Status   *bool  `json:"status"`
}
type service struct {
	db *gorm.DB
}

func (s *service) UpdateTenant(ctx context.Context, data *types.Tenants) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) DeleteTenant(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&types.Tenants{}, id).Error; err != nil {
			return err
		}
		return tx.Delete(&types.TenantAccountAssociations{}, "tenant_id = ?", id).Error
	})
	return
}

func (s *service) ListAccount(ctx context.Context, request ListAccountRequest) (res []types.Accounts, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Accounts{})
	if request.Email != "" {
		query = query.Where("email LIKE ?", "%"+request.Email+"%")
	}
	if request.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+request.Nickname+"%")
	}
	if request.IsLdap != nil {
		query = query.Where("is_ldap = ?", request.IsLdap)
	}
	if request.Status != nil {
		query = query.Where("status = ?", request.Status)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) UpdateAccount(ctx context.Context, data *types.Accounts) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) GetAccountById(ctx context.Context, id uint) (res types.Accounts, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) DeleteAccount(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&types.Accounts{}, id).Error; err != nil {
			return err
		}
		return tx.Delete(&types.TenantAccountAssociations{}, "account_id = ?", id).Error
	})
	return
}

func (s *service) CreateTenant(ctx context.Context, data *types.Tenants) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) ListTenants(ctx context.Context, request ListTenantRequest) (res []types.Tenants, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Tenants{})
	if request.Name != "" {
		query = query.Where("name LIKE ?", "%"+request.Name+"%")
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) CreateAccount(ctx context.Context, data *types.Accounts, tenantId uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(data).Error; err != nil {
			return err
		}
		return tx.Create(&types.TenantAccountAssociations{
			TenantID:  tenantId,
			AccountID: data.ID,
		}).Error
	})
	return
}

func (s *service) GetAccountByEmail(ctx context.Context, email string, preload ...string) (res types.Accounts, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v)
	}
	err = db.Where("email = ?", email).First(&res).Error
	return
}

func (s *service) GetTenantByUuid(ctx context.Context, uuid string, preload ...string) (res types.Tenants, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v)
	}
	err = db.Where("public_tenant_id = ?", uuid).First(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
