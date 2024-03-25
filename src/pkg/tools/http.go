package tools

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, mdw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, mdw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			vars := mux.Vars(request)
			if toolId, ok := vars["toolId"]; ok && !strings.EqualFold(toolId, "") {
				ctx = context.WithValue(ctx, contextKeyToolId, toolId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := MakeEndpoints(s, map[string][]endpoint.Middleware{
		"Tools": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/create", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/{toolId}", kithttp.NewServer(
		eps.DetailEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{toolId}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPut)
	r.Handle("/{toolId}", kithttp.NewServer(
		eps.DeleteEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{toolId}/assistants", kithttp.NewServer(
		eps.AssistantsEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	return r
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.page = page.GetPage(r)
	req.pageSize = page.GetPageSize(r)
	req.name = r.URL.Query().Get("name")
	return req, nil
}
