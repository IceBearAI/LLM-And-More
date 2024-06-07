package account

import (
	"context"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/auth"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户管理
type Service interface {
	// CreateAccount 创建账号
	// @kit-http /account POST
	// @kit-http-request CreateAccountRequest true
	CreateAccount(ctx context.Context, req CreateAccountRequest) (res Account, err error)
	// ListAccount 获取账号列表
	// @kit-http /accounts GET
	// @kit-http-request ListAccountRequest
	ListAccount(ctx context.Context, req ListAccountRequest) (list []Account, total int64, err error)
	// UpdateAccount 更新账号
	// @kit-http /account/{id} PUT
	// @kit-http-request UpdateAccountRequest true
	UpdateAccount(ctx context.Context, req UpdateAccountRequest) (err error)

	// DeleteAccount 删除账号
	// @kit-http /account/{id} DELETE
	// @kit-http-request DeleteAccountRequest
	DeleteAccount(ctx context.Context, id uint) (err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	//rdb     redis.UniversalClient
	apiSvc services.Service
}

// DeleteAccount implements Service.
func (s *service) DeleteAccount(ctx context.Context, id uint) (err error) {
	err = s.store.Auth().DeleteAccount(ctx, id)
	if err != nil {
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "删除账号失败"))
		return
	}
	return
}

func (s *service) UpdateAccount(ctx context.Context, request UpdateAccountRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UpdateAccount")
	data, err := s.store.Auth().GetAccountById(ctx, request.Id)
	if err != nil {
		_ = level.Error(logger).Log("auth", "UpdateAccount", "err", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return encode.ErrAccountNotFound.Error()
		}
		return encode.ErrSystem.Wrap(errors.Wrap(err, "查询账号失败"))
	}
	if request.Email != "" {
		res, err := s.store.Auth().GetAccountByEmail(ctx, request.Email)
		if err == nil && res.ID != data.ID {
			_ = level.Error(logger).Log("auth", "UpdateAccount", "err", "邮箱账号已存在")
			return encode.InvalidParams.Wrap(errors.Wrap(err, "要更新邮箱账号已存在"))
		}
		data.Email = request.Email
	}
	if request.Nickname != "" {
		data.Nickname = request.Nickname
	}

	if request.Language != "" {
		data.Language = request.Language
	}

	if request.Status != nil {
		data.Status = *request.Status
	}
	if request.IsLdap != nil {
		data.IsLdap = *request.IsLdap
	}

	if request.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			_ = level.Error(logger).Log("auth", "UpdateAccount", "err", err.Error())
			return encode.ErrSystem.Wrap(errors.Wrap(err, "更新账号失败，请联系管理员"))
		}
		data.PasswordHash = string(hash)
	}

	data.Tenants, err = s.store.Tenants().FindByPublicTenantIdItems(ctx, request.TenantPublicTenantIdItems)
	if err != nil {
		_ = level.Error(logger).Log("auth", "UpdateAccount", "err", err.Error())
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "查询租户失败"))
		return
	}

	err = s.store.Auth().UpdateAccount(ctx, &data)
	if err != nil {
		_ = level.Error(logger).Log("auth", "UpdateAccount", "err", err.Error())
		return err
	}
	return nil
}

func (s *service) ListAccount(ctx context.Context, request ListAccountRequest) (list []Account, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListAccount")
	req := auth.ListAccountRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
		Email:    request.Email,
		IsLdap:   request.IsLdap,
	}
	accounts, total, err := s.store.Auth().ListAccount(ctx, req)
	if err != nil {
		_ = level.Error(logger).Log("auth", "ListAccount", "err", err.Error())
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "查询账号失败"))
		return
	}
	for _, a := range accounts {
		list = append(list, convertAccount(&a))
	}
	return
}

func (s *service) CreateAccount(ctx context.Context, request CreateAccountRequest) (res Account, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateAccount")
	data := types.Accounts{
		Email:    request.Email,
		Nickname: request.Nickname,
		Language: request.Language,
		IsLdap:   request.IsLdap,
		Status:   true,
	}
	if !data.IsLdap {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			_ = level.Error(logger).Log("auth", "CreateAccount", "err", err.Error())
			return res, encode.ErrSystem.Wrap(errors.Wrap(err, "创建账号失败，请联系管理员"))
		}
		data.PasswordHash = string(hash)
	}
	data.Tenants, err = s.store.Tenants().FindByPublicTenantIdItems(ctx, request.TenantPublicTenantIdItems)
	if err != nil {
		_ = level.Error(logger).Log("auth", "CreateAccount", "err", err.Error())
		err = encode.ErrSystem.Wrap(errors.Wrap(err, "查询租户失败"))
		return
	}

	err = s.store.Auth().CreateAccountV2(ctx, &data)
	if err != nil {
		_ = level.Error(logger).Log("auth", "CreateAccount", "err", err.Error())
		return res, err
	}
	res = convertAccount(&data)
	return res, nil
}

func (s *service) Account(ctx context.Context, email string) (res AccountResult, err error) {
	email, ok := middleware.GetEmail(ctx)
	if !ok {
		_ = level.Error(s.logger).Log("auth", "Account, GetEmail error", "email", email)
		return res, encode.ErrAccountNotFound.Error()
	}
	account, err := s.store.Auth().GetAccountByEmail(ctx, email, "Tenants")
	if err != nil {
		_ = level.Error(s.logger).Log("auth", "account, GetAccountByEmail error", "email", email, "err", err.Error())
		return res, encode.ErrAccountNotFound.Error()
	}
	res.Email = account.Email
	res.Nickname = account.Nickname
	res.Language = account.Language
	res.Tenants = make([]Tenant, 0)
	for _, t := range account.Tenants {
		res.Tenants = append(res.Tenants, Tenant{
			Id:   t.PublicTenantID,
			Name: t.Name,
		})
	}
	return res, nil
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

func convertAccount(data *types.Accounts) Account {
	res := Account{
		Id:        data.ID,
		Email:     data.Email,
		IsLdap:    data.IsLdap,
		Nickname:  data.Nickname,
		Language:  data.Language,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	for _, v := range data.Tenants {
		res.Tenants = append(res.Tenants, convertTenant(&v))
	}
	return res

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
