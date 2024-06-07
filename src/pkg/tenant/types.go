package tenant

import "time"

type Tenant struct {
	// TenantId 租户ID
	Id string `json:"id"`
	// TenantName 租户名称
	Name string `json:"name"`
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

type ListTenantRequest struct {
	Page     int    `json:"page" param:"query,page"`
	PageSize int    `json:"pageSize" param:"query,pageSize"`
	Name     string `json:"name" param:"query,name"`
}

type CreateTenantRequest struct {
	Name         string   `json:"name" validate:"required"`
	ContactEmail string   `json:"contactEmail" validate:"required"`
	ModelNames   []string `json:"modelNames"`
}

type DeleteTenantRequest struct {
	Id uint `json:"publicTenantId" validate:"required" param:"path,id"`
}

type UpdateTenantRequest struct {
	Id           uint     `json:"id" param:"path,id" validate:"required"`
	Name         string   `json:"name" validate:"required"`
	ContactEmail string   `json:"contactEmail" validate:"required"`
	ModelNames   []string `json:"modelNames"`
}
