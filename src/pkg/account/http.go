// Code generated . DO NOT EDIT.
package account

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


		CreateAccountMethodName: ems,


		DeleteAccountMethodName: ems,


		ListAccountMethodName: ems,


		UpdateAccountMethodName: ems,

	})

	r := mux.NewRouter()


	r.Handle("/accounts", kithttp.NewServer(
		eps.CreateAccountEndpoint,
		decodeCreateAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("POST")


	r.Handle("/accounts/{id}", kithttp.NewServer(
		eps.DeleteAccountEndpoint,
		decodeDeleteAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("DELETE")



	r.Handle("/accounts", kithttp.NewServer(
		eps.ListAccountEndpoint,
		decodeListAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")



	r.Handle("/accounts/{id}", kithttp.NewServer(
		eps.UpdateAccountEndpoint,
		decodeUpdateAccountRequest,
		encode.JsonResponse,
		opts...,
	)).Methods("PUT")

	return r
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


