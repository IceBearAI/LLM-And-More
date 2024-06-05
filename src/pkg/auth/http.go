package auth

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, dmw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return ctx
		}),
		kithttp.ServerAfter(func(ctx context.Context, response http.ResponseWriter) context.Context {
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Auth": ems,
	})

	r := mux.NewRouter()

	r.Handle("/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/account", kithttp.NewServer(
		eps.AccountEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/account", kithttp.NewServer(
		eps.CreateAccountEndpoint,
		decodeCreateAccountRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/tenants", kithttp.NewServer(
		eps.ListTenantsEndpoint,
		decodeListTenantsRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/tenants", kithttp.NewServer(
		eps.CreateTenantEndpoint,
		decodeCreateTenantRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	email, _ := middleware.GetEmail(ctx)
	return accountRequest{
		Email: email,
	}, nil
}

func decodeCreateAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !req.IsLdap && req.Password == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("创建非LDAP账户，密码不能为空"))
	}
	return req, nil
}

func decodeListTenantsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListTenantRequest{
		Name:     r.URL.Query().Get("name"),
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	return req, nil
}

func decodeCreateTenantRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}
