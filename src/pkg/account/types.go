package account

import "time"

type AccountRequest struct {
	Email string `json:"email" param:",email"`
}
type Tenant struct {
	// TenantId 租户ID
	Id string `json:"id"`
	// TenantName 租户名称
	Name string `json:"name"`
}
type AccountResult struct {
	// Tenant 租户信息
	Tenants  []Tenant `json:"tenants"`
	Email    string   `json:"email"`
	Nickname string   `json:"nickname"`
	Language string   `json:"language"`
}

type TenantDetail struct {
	Id             uint      `json:"id"`
	Name           string    `json:"name"`
	PublicTenantID string    `json:"publicTenantId"`
	ContactEmail   string    `json:"contactEmail"`
	ModelNames     []string  `json:"modelNames"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateAccountRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	IsLdap   bool   `json:"isLdap"`
	Password string `json:"password"`
	Language string `json:"language" validate:"required"`

	TenantPublicTenantIdItems []string `json:"tenantPublicTenantIdItems"`
}

type Account struct {
	Id        uint      `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Status    bool      `json:"status"`
	IsLdap    bool      `json:"isLdap"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Tenants []TenantDetail `json:"tenants"`
}

type ListAccountRequest struct {
	Page     int    `json:"page" param:"query,page"`
	PageSize int    `json:"pageSize" param:"query,pageSize"`
	Email    string `json:"email" param:"query,email"`
	Nickname string `json:"nickname" param:"query,nickname"`
	IsLdap   *bool  `json:"isLdap,omitempty" param:"query,isLdap"`
	Status   *bool  `json:"status,omitempty"`
}

type UpdateAccountRequest struct {
	Id                        uint     `json:"id" param:"path,id" validate:"required"`
	Nickname                  string   `json:"nickname"`
	Email                     string   `json:"email"`
	IsLdap                    *bool    `json:"isLdap"`
	Language                  string   `json:"language"`
	Status                    *bool    `json:"status,omitempty"`
	Password                  string   `json:"password"`
	TenantPublicTenantIdItems []string `json:"tenantPublicTenantIdItems"`
}

type DeleteAccountRequest struct {
	Id uint `json:"id" validate:"required" param:"path,id"`
}
