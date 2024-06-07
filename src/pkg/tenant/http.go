// Code generated . DO NOT EDIT.
package tenant

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


		CreateTenantMethodName: ems,


		DeleteTenantMethodName: ems,


		ListTenantsMethodName: ems,

		UpdateTenantMethodName: ems,
	})

	r := mux.NewRouter()




	r.Handle("/tenants", kithttp.NewServer(
		eps.CreateTenantEndpoint,
		decodeCreateTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("POST")



	r.Handle("/tenants/{id}", kithttp.NewServer(
		eps.DeleteTenantEndpoint,
		decodeDeleteTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("DELETE")



	r.Handle("/tenants", kithttp.NewServer(
		eps.ListTenantsEndpoint,
		decodeListTenantsRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")

	

	r.Handle("/tenants/{id}", kithttp.NewServer(
		eps.UpdateTenantEndpoint,
		decodeUpdateTenantRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("PUT")

	return r
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
