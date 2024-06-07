// Code generated . DO NOT EDIT.
package tenant

import (
	"context"

	"github.com/IceBearAI/aigc/src/encode"
	endpoint "github.com/go-kit/kit/endpoint"
)

const (


	CreateTenantMethodName = "CreateTenant"


	DeleteTenantMethodName = "DeleteTenant"


	ListTenantsMethodName = "ListTenants"


	UpdateTenantMethodName = "UpdateTenant"
)

type Endpoints struct {


	CreateTenantEndpoint endpoint.Endpoint


	DeleteTenantEndpoint endpoint.Endpoint


	ListTenantsEndpoint endpoint.Endpoint


	UpdateTenantEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{


		CreateTenantEndpoint: makeCreateTenantEndpoint(s),


		DeleteTenantEndpoint: makeDeleteTenantEndpoint(s),


		ListTenantsEndpoint: makeListTenantsEndpoint(s),


		UpdateTenantEndpoint: makeUpdateTenantEndpoint(s),
	}


	for _, m := range dmw[CreateTenantMethodName] {
		eps.CreateTenantEndpoint = m(eps.CreateTenantEndpoint)
	}


	for _, m := range dmw[DeleteTenantMethodName] {
		eps.DeleteTenantEndpoint = m(eps.DeleteTenantEndpoint)
	}


	for _, m := range dmw[ListTenantsMethodName] {
		eps.ListTenantsEndpoint = m(eps.ListTenantsEndpoint)
	}



	for _, m := range dmw[UpdateTenantMethodName] {
		eps.UpdateTenantEndpoint = m(eps.UpdateTenantEndpoint)
	}

	return eps
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
