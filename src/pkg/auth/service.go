package auth

import (
	"context"
	"time"

	"github.com/IceBearAI/aigc/src/encode"
	authjwt "github.com/IceBearAI/aigc/src/jwt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/auth"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户管理
// @tags Account
type Service interface {
	// Login 平台授权登陆
	// username, password: 用户名密码
	// @kit-http /login POST
	// @kit-http-request LoginRequest true
	Login(ctx context.Context, req LoginRequest) (res LoginResult, err error)
	// Account 获取账号信息
	// @kit-http /account GET
	// @kit-http-request AccountRequest
	Account(ctx context.Context, email string) (res AccountResult, err error)
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

	err = s.store.Auth().CreateAccount(ctx, &data)
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

func (s *service) Login(ctx context.Context, req LoginRequest) (res LoginResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	username := req.Username
	password := req.Password
	if username == "" || password == "" {
		_ = level.Error(logger).Log("auth", "login, username or password is empty", "username", username, "password", password)
		return res, encode.ErrAccountLogin.Error()
	}

	//loginKey := fmt.Sprintf("aigc:auth:login:%s", username)
	//if s.rdb.Incr(ctx, loginKey).Val() > 1 {
	//	_ = s.rdb.Expire(ctx, loginKey, time.Minute).Err()
	//	_ = level.Error(logger).Log("auth", "login", "username", username, "err", encode.ErrLimiter.Error())
	//	return res, encode.ErrLimiter.Error()
	//}
	//defer func() {
	//	_ = s.rdb.Del(ctx, loginKey).Err()
	//}()
	//
	//errLoginKey := fmt.Sprintf("aigc:auth:login:err:%s", username)
	//errNum, err := s.rdb.Get(ctx, errLoginKey).Int()
	//if err != nil && !errors.Is(err, redis.Nil) {
	//	_ = level.Error(logger).Log("auth", "login, redis get errLoginKey error", "username", username, "err", err.Error())
	//	return res, encode.ErrSystem.Wrap(errors.Wrap(err,"登录失败，请联系管理员"))
	//}
	//if errNum >= 5 {
	//	_ = level.Error(logger).Log("auth", "login too many error times", "username", username, "err", encode.ErrLimiter.Error())
	//	return res, encode.ErrAccountLocked.Error()
	//}

	// 获取账号信息
	account, err := s.store.Auth().GetAccountByEmail(ctx, req.Username)
	if err != nil {
		_ = level.Error(logger).Log("auth", "login, GetAccountByEmail error", "username", username, "err", err.Error())
		return res, encode.ErrAccountNotFound.Error()
	}
	// 账号被锁定
	if !account.Status {
		_ = level.Error(logger).Log("auth", "login, account status not active", "username", username)
		return res, encode.ErrAccountLocked.Error()
	}

	if account.IsLdap {
		authenticate, err := s.apiSvc.Ldap().Authenticate(ctx, username, password)
		if err != nil {
			_ = level.Error(logger).Log("auth", "login", "username", username, "err", err.Error())
			return res, encode.ErrAccountLogin.Error()
		}
		if !authenticate {
			//_ = s.rdb.Incr(ctx, errLoginKey).Err()
			//_ = s.rdb.Expire(ctx, errLoginKey, time.Minute*30).Err()
			_ = level.Error(logger).Log("auth", "login authenticate false", "username", username)
			return res, encode.ErrAccountLogin.Error()
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
		if err != nil {
			//_ = s.rdb.Incr(ctx, errLoginKey).Err()
			//_ = s.rdb.Expire(ctx, errLoginKey, time.Minute*30).Err()
			_ = level.Error(logger).Log("auth", "login, bcrypt.CompareHashAndPassword error", "username", username, "err", err.Error())
			return res, encode.ErrAccountLogin.Error()
		}
	}
	tk, err := s.jwtToken(ctx, authjwt.TokenSourceNd, time.Duration(168)*time.Hour, username, "", account.ID)
	if err != nil {
		_ = level.Error(logger).Log("auth", "login, jwtToken error", "username", username, "err", err.Error())
		return res, encode.ErrSystem.Wrap(errors.Wrap(err, "登录失败，请联系管理员"))
	}
	res.Token = tk
	res.Username = username
	return res, nil
}

func (s *service) jwtToken(ctx context.Context, source authjwt.TokenSource, timeout time.Duration, email string, qwUserid string, accountId uint) (tk string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	expAt := jwt2.NewNumericDate(time.Now().Add(timeout))

	// 创建声明
	claims := authjwt.ArithmeticCustomClaims{
		Source: source,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
		UserId: accountId,
	}

	//不同类型，参数不一样
	claims.Email = email

	//创建token，指定加密算法为HS256
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)
	//生成token
	tk, err = token.SignedString([]byte(authjwt.GetJwtKey()))
	if err != nil {
		_ = level.Error(logger).Log("auth", "jwtToken, SignedString error", "err", err.Error(), "source", source, "email", email, "qwUserid", qwUserid, "accountId", accountId)
		return tk, nil
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
