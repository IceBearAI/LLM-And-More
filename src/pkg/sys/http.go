package sys

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, dmw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sys": ems,
	})

	r := mux.NewRouter()
	r.Handle("/dict", kithttp.NewServer(
		eps.CreateDictEndpoint,
		decodeCreateDictRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/dict", kithttp.NewServer(
		eps.ListDictEndpoint,
		decodeListDictRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/dict/{id}", kithttp.NewServer(
		eps.UpdateDictEndpoint,
		decodeUpdateDictRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	r.Handle("/dict/{id}", kithttp.NewServer(
		eps.DeleteDictEndpoint,
		decodeDictIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/dict/tree", kithttp.NewServer(
		eps.GetDictTreeEndpoint,
		decodeDictCodeRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/audit", kithttp.NewServer(
		eps.ListAuditEndpoint,
		decodeListAuditRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/template", kithttp.NewServer(
		eps.ListTemplateEndpoint,
		decodeListTemplateRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/template", kithttp.NewServer(
		eps.CreateTemplateEndpoint,
		decodeCreateTemplateRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/template/{name}", kithttp.NewServer(
		eps.UpdateTemplateEndpoint,
		decodeUpdateTemplateRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	r.Handle("/template/{name}", kithttp.NewServer(
		eps.DeleteTemplateEndpoint,
		decodeDeleteTemplateRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	return r
}

func decodeListTemplateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := templateListRequest{
		Page:         page.GetPage(r),
		PageSize:     page.GetPageSize(r),
		Name:         r.URL.Query().Get("name"),
		TemplateType: r.URL.Query().Get("templateType"),
	}
	return req, nil
}

func decodeCreateTemplateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req templateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	if req.OutputDir == "" {
		req.OutputDir = "/data/ft-model"
	}

	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeUpdateTemplateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req templateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	vars := mux.Vars(r)
	req.Name = vars["name"]

	if req.OutputDir == "" {
		req.OutputDir = "/data/ft-model"
	}

	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeDeleteTemplateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req templateDeleteRequest
	vars := mux.Vars(r)
	req.Name = vars["name"]

	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeCreateDictRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateDictRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	// 父ID大于0，字典值不能为空
	if req.ParentID > 0 && req.DictValue == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("字典值不能为空"))
	}
	if req.ParentID == 0 && req.Code == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("字典编码不能为空"))
	}

	err = validate.Struct(req)
	return req, err
}

func decodeUpdateDictRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("dictId is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req UpdateDictRequest
	req.ID = uint(mid)
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err = validate.Struct(req)
	return req, err
}

func decodeDictIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("dictId is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req IdRequest
	req.Id = uint(mid)
	return req, nil
}

func decodeDictCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetDictTreeByCodeRequest
	query := r.URL.Query()
	req.Code = query["code"]
	return req, nil
}

func decodeListDictRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ListDictRequest
	req.Page = page.GetPage(r)
	req.PageSize = page.GetPageSize(r)
	req.Code = r.URL.Query().Get("code")
	req.Label = r.URL.Query().Get("label")
	parentId := r.URL.Query().Get("parentId")
	if parentId != "" {
		parse, err := strconv.Atoi(parentId)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(errors.New("parentId格式错误"))
		}
		req.ParentId = uint(parse)
	}
	return req, nil
}

func decodeListAuditRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ListAuditRequest
	req.Page = page.GetPage(r)
	req.PageSize = page.GetPageSize(r)
	req.TraceId = r.URL.Query().Get("traceId")
	req.Operator = r.URL.Query().Get("operator")
	startTime := r.URL.Query().Get("startTime")
	if startTime != "" {
		parse, err := time.Parse(time.DateTime, startTime)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(errors.New("起始时间格式错误"))
		}
		req.StartTime = &parse
	}
	endTime := r.URL.Query().Get("endTime")
	if endTime != "" {
		parse, err := time.Parse(time.DateTime, endTime)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(errors.New("结束时间格式错误"))
		}
		req.EndTime = &parse
	}

	isError := r.URL.Query().Get("isError")
	if isError != "" {
		parse, err := strconv.ParseBool(isError)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(errors.New("isError格式错误"))
		}
		req.IsError = &parse
	}

	duration := r.URL.Query().Get("duration")
	if duration != "" {
		parse, err := strconv.ParseFloat(duration, 64)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(errors.New("duration格式错误"))
		}
		req.Duration = parse
	}
	return req, nil
}
