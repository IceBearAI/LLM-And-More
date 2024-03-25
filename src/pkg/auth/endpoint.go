package auth

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	loginRequest struct {
		// Username 登陆用户的邮箱或邮箱前缀
		Username string `json:"username" validate:"required"`
		// Username 登陆用户的邮箱密码
		Password string `json:"password" validate:"required"`
	}

	loginResult struct {
		// Token jwt token
		Token string `json:"token"`
		// Username 登陆用户的姓名
		Username string `json:"username"`
		// Avatar 登陆用户的头像地址
		Avatar string `json:"avatar,omitempty"`
	}
	accountRequest struct {
		Email string `json:"email"`
	}
	tenant struct {
		// TenantId 租户ID
		Id string `json:"id"`
		// TenantName 租户名称
		Name string `json:"name"`
	}
	accountResult struct {
		// Tenant 租户信息
		Tenants  []tenant `json:"tenants"`
		Email    string   `json:"email"`
		Nickname string   `json:"nickname"`
		Language string   `json:"language"`
	}

	TenantDetail struct {
		Id             uint      `json:"id"`
		Name           string    `json:"name"`
		PublicTenantID string    `json:"publicTenantId"`
		ContactEmail   string    `json:"contactEmail"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	ListTenantRequest struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Name     string `json:"name"`
	}

	ListTenantResponse struct {
		Tenants []TenantDetail `json:"list"`
		Total   int64          `json:"total"`
	}

	CreateAccountRequest struct {
		Nickname string `json:"nickname" validate:"required"`
		Email    string `json:"email" validate:"required"`
		IsLdap   bool   `json:"isLdap"`
		Password string `json:"password"`
		TenantId uint   `json:"tenantId"`
		Language string `json:"language" validate:"required"`
	}

	Account struct {
		Id        uint      `json:"id"`
		Email     string    `json:"email"`
		Nickname  string    `json:"nickname"`
		Status    bool      `json:"status"`
		IsLdap    bool      `json:"isLdap"`
		Language  string    `json:"language"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	ListAccountRequest struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		IsLdap   *bool  `json:"isLdap,omitempty"`
		Status   *bool  `json:"status,omitempty"`
	}

	ListAccountResponse struct {
		Accounts []Account `json:"list"`
		Total    int64     `json:"total"`
	}

	CreateTenantRequest struct {
		Name         string `json:"name" validate:"required"`
		ContactEmail string `json:"contactEmail" validate:"required"`
	}

	UpdateAccountRequest struct {
		Id       uint   `json:"id" validate:"required"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		IsLdap   *bool  `json:"isLdap"`
		Language string `json:"language"`
		Status   *bool  `json:"status,omitempty"`
		Password string `json:"password"`
	}
)

type Endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	AccountEndpoint       endpoint.Endpoint
	CreateAccountEndpoint endpoint.Endpoint
	ListTenantsEndpoint   endpoint.Endpoint
	CreateTenantEndpoint  endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		LoginEndpoint:         makeLoginEndpoint(s),
		AccountEndpoint:       makeAccountEndpoint(s),
		CreateAccountEndpoint: makeCreateAccountEndpoint(s),
		ListTenantsEndpoint:   makeListTenantsEndpoint(s),
		CreateTenantEndpoint:  makeCreateTenantEndpoint(s),
	}

	for _, m := range dmw["Auth"] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
		eps.AccountEndpoint = m(eps.AccountEndpoint)
		eps.CreateAccountEndpoint = m(eps.CreateAccountEndpoint)
		eps.ListTenantsEndpoint = m(eps.ListTenantsEndpoint)
		eps.CreateTenantEndpoint = m(eps.CreateTenantEndpoint)
	}
	return eps
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		res, err := s.Login(ctx, req.Username, req.Password)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email, _ := middleware.GetEmail(ctx)
		res, err := s.Account(ctx, email)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateAccountRequest)
		if req.TenantId == 0 {
			req.TenantId, _ = middleware.GetTenantId(ctx)
		}
		res, err := s.CreateAccount(ctx, req)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeListTenantsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListTenantRequest)
		res, err := s.ListTenants(ctx, req)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCreateTenantEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateTenantRequest)
		resp, err := s.CreateTenant(ctx, req)
		return encode.Response{
			Data:  resp,
			Error: err,
		}, err
	}
}
