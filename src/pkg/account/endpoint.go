// Code generated . DO NOT EDIT.
package account

import (
	"context"

	"github.com/IceBearAI/aigc/src/encode"
	endpoint "github.com/go-kit/kit/endpoint"
)

const (

	CreateAccountMethodName = "CreateAccount"


	DeleteAccountMethodName = "DeleteAccount"


	ListAccountMethodName = "ListAccount"


	UpdateAccountMethodName = "UpdateAccount"

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


		DeleteAccountEndpoint: makeDeleteAccountEndpoint(s),


		ListAccountEndpoint: makeListAccountEndpoint(s),


		UpdateAccountEndpoint: makeUpdateAccountEndpoint(s),

	}



	for _, m := range dmw[CreateAccountMethodName] {
		eps.CreateAccountEndpoint = m(eps.CreateAccountEndpoint)
	}


	for _, m := range dmw[DeleteAccountMethodName] {
		eps.DeleteAccountEndpoint = m(eps.DeleteAccountEndpoint)
	}


	for _, m := range dmw[ListAccountMethodName] {
		eps.ListAccountEndpoint = m(eps.ListAccountEndpoint)
	}



	for _, m := range dmw[UpdateAccountMethodName] {
		eps.UpdateAccountEndpoint = m(eps.UpdateAccountEndpoint)
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

