// Code generated . DO NOT EDIT.
package sysauth

import (
	"context"

	"github.com/IceBearAI/aigc/src/encode"
	endpoint "github.com/go-kit/kit/endpoint"
)

const (
	AccountMethodName = "Account"

	CreateAccountMethodName = "CreateAccount"

	CreateTenantMethodName = "CreateTenant"

	DeleteAccountMethodName = "DeleteAccount"

	DeleteTenantMethodName = "DeleteTenant"

	ListAccountMethodName = "ListAccount"

	ListTenantsMethodName = "ListTenants"

	LoginMethodName = "Login"

	UpdateAccountMethodName = "UpdateAccount"

	UpdateTenantMethodName = "UpdateTenant"
)

type Endpoints struct {
	AccountEndpoint endpoint.Endpoint

	CreateAccountEndpoint endpoint.Endpoint

	CreateTenantEndpoint endpoint.Endpoint

	DeleteAccountEndpoint endpoint.Endpoint

	DeleteTenantEndpoint endpoint.Endpoint

	ListAccountEndpoint endpoint.Endpoint

	ListTenantsEndpoint endpoint.Endpoint

	LoginEndpoint endpoint.Endpoint

	UpdateAccountEndpoint endpoint.Endpoint

	UpdateTenantEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{

		CreateAccountEndpoint: makeCreateAccountEndpoint(s),

		CreateTenantEndpoint: makeCreateTenantEndpoint(s),

		DeleteAccountEndpoint: makeDeleteAccountEndpoint(s),

		DeleteTenantEndpoint: makeDeleteTenantEndpoint(s),

		ListAccountEndpoint: makeListAccountEndpoint(s),

		ListTenantsEndpoint: makeListTenantsEndpoint(s),

		UpdateAccountEndpoint: makeUpdateAccountEndpoint(s),

		UpdateTenantEndpoint: makeUpdateTenantEndpoint(s),
	}

	for _, m := range dmw[AccountMethodName] {
		eps.AccountEndpoint = m(eps.AccountEndpoint)
	}

	for _, m := range dmw[CreateAccountMethodName] {
		eps.CreateAccountEndpoint = m(eps.CreateAccountEndpoint)
	}

	for _, m := range dmw[CreateTenantMethodName] {
		eps.CreateTenantEndpoint = m(eps.CreateTenantEndpoint)
	}

	for _, m := range dmw[DeleteAccountMethodName] {
		eps.DeleteAccountEndpoint = m(eps.DeleteAccountEndpoint)
	}

	for _, m := range dmw[DeleteTenantMethodName] {
		eps.DeleteTenantEndpoint = m(eps.DeleteTenantEndpoint)
	}

	for _, m := range dmw[ListAccountMethodName] {
		eps.ListAccountEndpoint = m(eps.ListAccountEndpoint)
	}

	for _, m := range dmw[ListTenantsMethodName] {
		eps.ListTenantsEndpoint = m(eps.ListTenantsEndpoint)
	}

	for _, m := range dmw[LoginMethodName] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
	}

	for _, m := range dmw[UpdateAccountMethodName] {
		eps.UpdateAccountEndpoint = m(eps.UpdateAccountEndpoint)
	}

	for _, m := range dmw[UpdateTenantMethodName] {
		eps.UpdateTenantEndpoint = m(eps.UpdateTenantEndpoint)
	}

	return eps
}

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(CreateAccountRequest)

		var res Account

		res, err = s.CreateAccount(
			ctx,

			req,
		)

		return encode.Response{
			Data:  res,
			Error: err,
		}, nil

	}
}

func makeCreateTenantEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(CreateTenantRequest)

		var res TenantDetail

		res, err = s.CreateTenant(
			ctx,

			req,
		)

		return encode.Response{
			Data:  res,
			Error: err,
		}, nil

	}
}

func makeDeleteAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(DeleteAccountRequest)

		err = s.DeleteAccount(
			ctx,

			req.Id,
		)

		return encode.Response{
			Data:  map[string]interface{}{},
			Error: err,
		}, nil

	}
}

func makeDeleteTenantEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(DeleteTenantRequest)

		err = s.DeleteTenant(
			ctx,

			req.Id,
		)

		return encode.Response{
			Data:  map[string]interface{}{},
			Error: err,
		}, nil

	}
}

func makeListAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(ListAccountRequest)

		var list []Account

		var total int64

		list, total, err = s.ListAccount(
			ctx,

			req,
		)

		return encode.Response{
			Data: map[string]interface{}{
				"list":  list,
				"total": total},
			Error: err,
		}, nil

	}
}

func makeListTenantsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(ListTenantRequest)

		var list []TenantDetail

		var total int64

		list, total, err = s.ListTenants(
			ctx,

			req,
		)

		return encode.Response{
			Data: map[string]interface{}{
				"list":  list,
				"total": total},
			Error: err,
		}, nil

	}
}



func makeUpdateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(UpdateAccountRequest)

		err = s.UpdateAccount(
			ctx,

			req,
		)

		return encode.Response{
			Data:  map[string]interface{}{},
			Error: err,
		}, nil

	}
}

func makeUpdateTenantEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(UpdateTenantRequest)

		err = s.UpdateTenant(
			ctx,

			req,
		)

		return encode.Response{
			Data:  map[string]interface{}{},
			Error: err,
		}, nil

	}
}
