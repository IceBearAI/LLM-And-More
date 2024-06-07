package tenant

import (
	"context"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/auth"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// 用户管理
// @tags Account
type Service interface {
	// ListTenants 获取租户列表
	// @kit-http /tenants GET
	// @kit-http-request ListTenantRequest
	ListTenants(ctx context.Context, req ListTenantRequest) (list []TenantDetail, total int64, err error)

	// UpdateTenant 更新租户
	// @kit-http /tenant/{id} PUT
	// @kit-http-request UpdateTenantRequest true
	UpdateTenant(ctx context.Context, req UpdateTenantRequest) (err error)
	// CreateTenant 创建租户
	// @kit-http /tenants POST
	// @kit-http-request CreateTenantRequest true
	CreateTenant(ctx context.Context, req CreateTenantRequest) (res TenantDetail, err error)
	// DeleteTenant 删除租户
	// @kit-http /tenant/{id} DELETE
	// @kit-http-request DeleteTenantRequest
	DeleteTenant(ctx context.Context, id uint) (err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	//rdb     redis.UniversalClient
	apiSvc services.Service
}

// UpdateTenant implements Service.
func (s *service) UpdateTenant(ctx context.Context, req UpdateTenantRequest) (err error) {
	t, err := s.store.Auth().GetTenantById(ctx, req.Id)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.Wrap(err, "查询租户失败"))
	}

	t.Name = req.Name
	t.ContactEmail = req.ContactEmail
	t.Models, err = s.store.Model().FindByModelNames(ctx, req.ModelNames)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.Wrap(err, "查询模型失败"))
	}

	err = s.store.Auth().UpdateTenant(ctx, &t)
	if err != nil {
		return encode.ErrSystem.Wrap(errors.Wrap(err, "更新租户失败"))
	}

	return
}

// DeleteTenant implements Service.
func (s *service) DeleteTenant(ctx context.Context, id uint) (err error) {
	err = s.store.Auth().DeleteTenant(ctx, id)
	if err != nil {
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "删除租户失败"))
		return
	}
	return
}

func (s *service) CreateTenant(ctx context.Context, request CreateTenantRequest) (res TenantDetail, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateTenant")
	data := types.Tenants{
		Name:           request.Name,
		PublicTenantID: uuid.New().String(),
		ContactEmail:   request.ContactEmail,
	}

	models, err := s.store.Model().FindByModelNames(ctx, request.ModelNames)
	if err != nil {
		_ = level.Error(logger).Log("auth", "CreateTenant", "err", err.Error())
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "查询模型失败"))
		return res, err
	}
	data.Models = models
	err = s.store.Auth().CreateTenant(ctx, &data)
	if err != nil {
		_ = level.Error(logger).Log("auth", "CreateTenant", "err", err.Error())
		return res, err
	}
	res = convertTenant(&data)
	return res, nil
}

func (s *service) ListTenants(ctx context.Context, request ListTenantRequest) (list []TenantDetail, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListTenants")
	req := auth.ListTenantRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
		Name:     request.Name,
	}
	tenants, total, err := s.store.Auth().ListTenants(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("auth", "ListTenants", "err", err.Error())
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "查询租户失败"))
		return
	}
	for _, t := range tenants {
		list = append(list, convertTenant(&t))
	}
	return
}

func convertTenant(data *types.Tenants) TenantDetail {
	td := TenantDetail{
		Id:             data.ID,
		Name:           data.Name,
		PublicTenantID: data.PublicTenantID,
		ContactEmail:   data.ContactEmail,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}

	for _, v := range data.Models {
		td.ModelNames = append(td.ModelNames, v.ModelName)
	}
	return td
}

func New(logger log.Logger, traceId string,
	store repository.Repository,
	//rdb redis.UniversalClient,
	apiSvc services.Service) Service {
	return &service{
		logger:  logger,
		traceId: traceId,
		store:   store,
		//rdb:     rdb,
		apiSvc: apiSvc,
	}
}
