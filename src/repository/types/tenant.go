package types

import (
	"gorm.io/gorm"
	"time"
)

const (
	SystemTenant = 1
)

type Tenants struct {
	gorm.Model
	Name           string     `gorm:"column:name;type:varchar(32);NOT NULL"`          // 租户名称
	PublicTenantID string     `gorm:"column:public_tenant_id;type:varchar(50)"`       // 租户ID
	ContactEmail   string     `gorm:"column:contact_email;type:varchar(50);NOT NULL"` // 联系人邮箱
	Accounts       []Accounts `gorm:"many2many:tenant_account_associations;foreignKey:ID;references:ID;joinForeignKey:TenantID;joinReferences:AccountID"`
	Models         []Models   `gorm:"many2many:tenant_model_associations;foreignKey:ID;references:ID;joinForeignKey:TenantID;joinReferences:ModelID"`
}

func (m *Tenants) TableName() string {
	return "tenants"
}

type Accounts struct {
	gorm.Model
	Email        string    `gorm:"column:email;type:varchar(100)"`             // 邮箱
	Nickname     string    `gorm:"column:nickname;type:varchar(100);NOT NULL"` // 昵称
	Language     string    `gorm:"column:language;type:varchar(10);NOT NULL"`  // 语言
	IsLdap       bool      `gorm:"column:is_ldap;type:tinyint(4)"`             // 是否LDAP用户
	PasswordHash string    `gorm:"column:password_hash;type:varchar(128)"`     // 密码
	Status       bool      `gorm:"column:status;type:tinyint(4);default:1"`    // 状态 0 禁用 1 启用
	Tenants      []Tenants `gorm:"many2many:tenant_account_associations;foreignKey:ID;references:ID;joinForeignKey:AccountID;joinReferences:TenantID"`
}

func (m *Accounts) TableName() string {
	return "accounts"
}

type TenantAccountAssociations struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TenantID  uint `gorm:"column:tenant_id;type:bigint(20);NOT NULL"`  // 租户ID
	AccountID uint `gorm:"column:account_id;type:bigint(20);NOT NULL"` // 账户ID
}

func (m *TenantAccountAssociations) TableName() string {
	return "tenant_account_associations"
}
