// Code generated . DO NOT EDIT.
package sysauth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IceBearAI/aigc/src/encode"
	endpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var validate = validator.New()
func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware

	opts = append(opts, kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
		return ctx
	}))

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{

		AccountMethodName: ems,

		CreateAccountMethodName: ems,

		CreateTenantMethodName: ems,

		DeleteAccountMethodName: ems,

		DeleteTenantMethodName: ems,

		ListAccountMethodName: ems,

		ListTenantsMethodName: ems,

		LoginMethodName: ems,

		UpdateAccountMethodName: ems,

		UpdateTenantMethodName: ems,
	})

	r := mux.NewRouter()

	r.Handle("/account", kithttp.NewServer(
		eps.AccountEndpoint,
		decodeAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/account", kithttp.NewServer(
		eps.CreateAccountEndpoint,
		decodeCreateAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/tenants", kithttp.NewServer(
		eps.CreateTenantEndpoint,
		decodeCreateTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/account/{id}", kithttp.NewServer(
		eps.DeleteAccountEndpoint,
		decodeDeleteAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/tenant/{id}", kithttp.NewServer(
		eps.DeleteTenantEndpoint,
		decodeDeleteTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/accounts", kithttp.NewServer(
		eps.ListAccountEndpoint,
		decodeListAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/tenants", kithttp.NewServer(
		eps.ListTenantsEndpoint,
		decodeListTenantsRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/account/{id}", kithttp.NewServer(
		eps.UpdateAccountEndpoint,
		decodeUpdateAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/tenant/{id}", kithttp.NewServer(
		eps.UpdateTenantEndpoint,
		decodeUpdateTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("PUT")

	return r
}

// Account
// @Summary  获取账号信息
// @Description  获取账号信息
// @tags Account
// @Accept json
// @Produce json
// @Success 200 {object} encode.Response{data=AccountResult}
// @Router /account [GET]
func decodeAccountRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := AccountRequest{}

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// CreateAccount
// @Summary  创建账号
// @Description  创建账号
// @tags Account
// @Accept json
// @Produce json
// @Param CreateAccountRequest body CreateAccountRequest 	true "http request body"
// @Success 200 {object} encode.Response{data=Account}
// @Router /account [POST]
func decodeCreateAccountRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := CreateAccountRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		err = encode.Invalid.Wrap(errors.Wrap(err, "decode body"))
		return
	}

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// CreateTenant
// @Summary  创建租户
// @Description  创建租户
// @tags Account
// @Accept json
// @Produce json
// @Param CreateTenantRequest body CreateTenantRequest 	true "http request body"
// @Success 200 {object} encode.Response{data=TenantDetail}
// @Router /tenants [POST]
func decodeCreateTenantRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := CreateTenantRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		err = encode.Invalid.Wrap(errors.Wrap(err, "decode body"))
		return
	}

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// DeleteAccount
// @Summary  删除账号
// @Description  删除账号
// @tags Account
// @Accept json
// @Produce json
// @Param id path string true " "
// @Success 200 {object} encode.Response{}
// @Router /account/{id} [DELETE]
func decodeDeleteAccountRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := DeleteAccountRequest{}

	var _id uint

	vars := mux.Vars(r)

	_id, err = cast.ToUintE(vars["id"])

	if err != nil {
		return
	}

	req.Id = _id

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// DeleteTenant
// @Summary  删除租户
// @Description  删除租户
// @tags Account
// @Accept json
// @Produce json
// @Param id path string true " "
// @Success 200 {object} encode.Response{}
// @Router /tenant/{id} [DELETE]
func decodeDeleteTenantRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := DeleteTenantRequest{}

	var _id uint

	vars := mux.Vars(r)

	_id, err = cast.ToUintE(vars["id"])

	if err != nil {
		return
	}

	req.Id = _id

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// ListAccount
// @Summary  获取账号列表
// @Description  获取账号列表
// @tags Account
// @Accept json
// @Produce json
// @Param email query string false " "
// @Param isLdap query string false " "
// @Param nickname query string false " "
// @Param page query string false " "
// @Param pageSize query string false " "
// @Success 200 {object} encode.Response{data.list=[]Account,data.total=int64}
// @Router /accounts [GET]
func decodeListAccountRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := ListAccountRequest{}

	var _email string

	var _isLdap bool

	var _nickname string

	var _page int

	var _pageSize int

	_email = r.URL.Query().Get("email")

	_isLdapStr := r.URL.Query().Get("isLdap")

	if _isLdapStr != "" {
		_isLdap, err = cast.ToBoolE(_isLdapStr)

		if err != nil {
			return
		}
	}

	_nickname = r.URL.Query().Get("nickname")

	_pageStr := r.URL.Query().Get("page")

	if _pageStr != "" {
		_page, err = cast.ToIntE(_pageStr)

		if err != nil {
			return
		}
	}

	_pageSizeStr := r.URL.Query().Get("pageSize")

	if _pageSizeStr != "" {
		_pageSize, err = cast.ToIntE(_pageSizeStr)

		if err != nil {
			return
		}
	}

	req.Email = _email

	if _isLdapStr != "" {
		req.IsLdap = &_isLdap
	}

	req.Nickname = _nickname

	req.Page = _page

	req.PageSize = _pageSize

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// ListTenants
// @Summary  获取租户列表
// @Description  获取租户列表
// @tags Account
// @Accept json
// @Produce json
// @Param name query string false " "
// @Param page query string false " "
// @Param pageSize query string false " "
// @Success 200 {object} encode.Response{data.list=[]TenantDetail,data.total=int64}
// @Router /tenants [GET]
func decodeListTenantsRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := ListTenantRequest{}

	var _name string

	var _page int

	var _pageSize int

	_name = r.URL.Query().Get("name")

	_pageStr := r.URL.Query().Get("page")

	if _pageStr != "" {
		_page, err = cast.ToIntE(_pageStr)

		if err != nil {
			return
		}
	}

	_pageSizeStr := r.URL.Query().Get("pageSize")

	if _pageSizeStr != "" {
		_pageSize, err = cast.ToIntE(_pageSizeStr)

		if err != nil {
			return
		}
	}

	req.Name = _name

	req.Page = _page

	req.PageSize = _pageSize

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// Login
// @Summary  平台授权登陆
// @Description  平台授权登陆
// @tags Account
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest 	true "http request body"
// @Success 200 {object} encode.Response{data=LoginResult}
// @Router /login [POST]
func decodeLoginRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := LoginRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		err = encode.Invalid.Wrap(errors.Wrap(err, "decode body"))
		return
	}

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// UpdateAccount
// @Summary  更新账号
// @Description  更新账号
// @tags Account
// @Accept json
// @Produce json
// @Param id path string true " "
// @Param UpdateAccountRequest body UpdateAccountRequest 	true "http request body"
// @Success 200 {object} encode.Response{}
// @Router /account/{id} [PUT]
func decodeUpdateAccountRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := UpdateAccountRequest{}

	var _id uint

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		err = encode.Invalid.Wrap(errors.Wrap(err, "decode body"))
		return
	}

	vars := mux.Vars(r)

	_id, err = cast.ToUintE(vars["id"])

	if err != nil {
		return
	}

	req.Id = _id

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}

// UpdateTenant
// @Summary  更新租户
// @Description  更新租户
// @tags Account
// @Accept json
// @Produce json
// @Param id path string true " "
// @Param UpdateTenantRequest body UpdateTenantRequest 	true "http request body"
// @Success 200 {object} encode.Response{}
// @Router /tenant/{id} [PUT]
func decodeUpdateTenantRequest(ctx context.Context, r *http.Request) (res interface{}, err error) {

	req := UpdateTenantRequest{}

	var _id uint

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		err = encode.Invalid.Wrap(errors.Wrap(err, "decode body"))
		return
	}

	vars := mux.Vars(r)

	_id, err = cast.ToUintE(vars["id"])

	if err != nil {
		return
	}

	req.Id = _id

err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, err
}
