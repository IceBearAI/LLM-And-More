package finetuning

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

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, dmw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			vars := mux.Vars(request)
			if fineTuningJobId, ok := vars["fineTuningJobId"]; ok && !strings.EqualFold(fineTuningJobId, "") {
				ctx = context.WithValue(ctx, contextKeyFineTuningJobId, fineTuningJobId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"FineTuning": ems,
	})

	r := mux.NewRouter()
	r.Handle("/finetuning", kithttp.NewServer(
		eps.CreateJobEndpoint,
		decodeCreateJobRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/finetuning", kithttp.NewServer(
		eps.ListJobEndpoint,
		decodeListJobRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/finetuning/{jobId}/cancel", kithttp.NewServer(
		eps.CancelJobEndpoint,
		decodeJobIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/finetuning/dashboard", kithttp.NewServer(
		eps.DashBoardEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/finetuning/base/model", kithttp.NewServer(
		eps.BaseModelEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/finetuning/{jobId}", kithttp.NewServer(
		eps.DeleteJobEndpoint,
		decodeJobIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/finetuning/{jobId}", kithttp.NewServer(
		eps.GetJobEndpoint,
		decodeJobIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/finetuning/estimate", kithttp.NewServer(
		eps.EstimateEndpoint,
		decodeCreateJobRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/finetuning/{fineTuningJobId}/finish", kithttp.NewServer(
		eps.UpdateJobFinishedStatusEndpoint,
		decodeUpdateJobFinishedStatusRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	return r
}

func decodeCreateJobRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//if err := validate.Struct(req); err != nil {
	//	return nil, encode.InvalidParams.Wrap(err)
	//}
	return req, nil
}

func decodeListJobRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ListJobRequest
	req = ListJobRequest{
		Page:           page.GetPage(r),
		PageSize:       page.GetPageSize(r),
		FineTunedModel: r.URL.Query().Get("fineTunedModel"),
		TrainStatus:    r.URL.Query().Get("trainStatus"),
	}
	return req, nil
}

func decodeJobIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]
	req := JobIdRequest{
		JobId: jobId,
	}
	return req, nil
}

func decodeUpdateJobFinishedStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateTrainStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}
