package terminal

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Token": ems,
	})

	r := mux.NewRouter()

	r.Handle("/resource/{resourceType}/service/{serviceName}/container/{containerName}/token", kithttp.NewServer(
		eps.TokenEndpoint,
		decodeTokenRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req tokenRequest

	vars := mux.Vars(r)
	resourceType, ok := vars["resourceType"]
	if !ok {
		return nil, encode.InvalidParams.Wrap(errors.New("resourceType is required"))
	}
	serviceName, ok := vars["serviceName"]
	if !ok {
		return nil, encode.InvalidParams.Wrap(errors.New("serviceName is required"))
	}
	containerName, ok := vars["containerName"]
	if !ok {
		return nil, encode.InvalidParams.Wrap(errors.New("containerName is required"))
	}

	req.ResourceType = resourceType
	req.ServiceName = serviceName
	req.ContainerName = containerName

	return req, nil
}
