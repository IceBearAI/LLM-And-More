package models

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
	"strings"
	"time"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, dmw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			vars := mux.Vars(request)
			if containerName, ok := vars["containerName"]; ok && !strings.EqualFold(containerName, "") {
				ctx = context.WithValue(ctx, contextKeyModelContainerName, containerName)
			}
			if modelName, ok := vars["modelName"]; ok && !strings.EqualFold(modelName, "") {
				ctx = context.WithValue(ctx, contextKeyModelName, modelName)
			}
			if id, ok := vars["id"]; ok && !strings.EqualFold(id, "") {
				ctx = context.WithValue(ctx, contextKeyModelId, id)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Model": ems,
	})

	r := mux.NewRouter()
	r.Handle("/models", kithttp.NewServer(
		eps.CreateModelEndpoint,
		decodeCreateModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models", kithttp.NewServer(
		eps.ListModelsEndpoint,
		decodeListModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.UpdateModelEndpoint,
		decodeUpdateModelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.DeleteModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/models/eval", kithttp.NewServer(
		eps.ListEvalEndpoint,
		decodeListEvalRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}", kithttp.NewServer(
		eps.GetModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{id}/deploy", kithttp.NewServer(
		eps.DeployModelEndpoint,
		decodeModelDeployRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/{id}/undeploy", kithttp.NewServer(
		eps.UndeployModelEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval", kithttp.NewServer(
		eps.CreateEvalEndpoint,
		decodeCreateEvalRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval/{id}/cancel", kithttp.NewServer(
		eps.CancelEvalEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/models/eval/{id}", kithttp.NewServer(
		eps.DeleteEvalEndpoint,
		decodeIdRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/models/{modelName}/container/{containerName}/logs", kithttp.NewServer(
		eps.GetModelLogsEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/{modelName}/info", kithttp.NewServer(
		eps.ModelInfoEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/models/chat/completions", kithttp.NewServer(
		eps.ChatCompletionStreamEndpoint,
		decodeChatCompletionStreamRequest,
		encodeChatCompletionsStreamResponse,
		kitopts...))
	return r
}

func decodeCreateModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeUpdateModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("modelId is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req UpdateModelRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req.Id = uint(mid)
	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListModelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListModelRequest{
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	req.ModelType = r.URL.Query().Get("modelType")
	req.ModelName = r.URL.Query().Get("modelName")
	req.ProviderName = r.URL.Query().Get("providerName")
	enabled := r.URL.Query().Get("enabled")
	if enabled != "" {
		b, err := strconv.ParseBool(enabled)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		req.Enabled = &b
	}
	isFineTuning := r.URL.Query().Get("isFineTuning")
	if isFineTuning != "" {
		b, err := strconv.ParseBool(isFineTuning)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		req.IsFineTuning = &b
	}
	return req, nil
}

func decodeIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("id is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req IdRequest
	req.Id = uint(mid)
	return req, nil
}

func decodeCreateEvalRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateEvalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListEvalRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListEvalRequest{
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	req.ModelName = r.URL.Query().Get("modelName")
	req.MetricName = r.URL.Query().Get("metricName")
	req.Status = r.URL.Query().Get("status")
	req.DatasetType = r.URL.Query().Get("datasetType")
	return req, nil
}

func decodeModelDeployRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("id is empty"))
	}
	mid, err := strconv.Atoi(id)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req ModelDeployRequest
	req.Id = uint(mid)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if req.Quantization == "" {
		req.Quantization = "float16"
	}
	if req.InferredType == "cpu" && req.Cpu == 0 {
		return nil, encode.InvalidParams.Wrap(errors.New("cpu is 0"))
	}
	if req.InferredType == "gpu" && req.Gpu == 0 {
		return nil, encode.InvalidParams.Wrap(errors.New("gpu is 0"))
	}

	return req, nil
}

func decodeChatCompletionStreamRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ChatCompletionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func encodeChatCompletionsStreamResponse(ctx context.Context, writer http.ResponseWriter, res interface{}) error {
	resp, ok := res.(encode.Response)
	if !ok {
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	stream, ok := resp.Data.(<-chan CompletionsStreamResult)
	if !ok {
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	flushWriter := writer.(http.Flusher)
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Paas-model", "chat/stream")
	//writer.Header().Set("Transfer-Encoding", "chunked")
	writer.WriteHeader(http.StatusOK)
	traceId, _ := ctx.Value("traceId").(string)
	writer.Header().Set("TraceId", traceId)
	for {
		select {
		case item, ok := <-stream:
			if !ok {
				return nil
			}
			if err := json.NewEncoder(writer).Encode(encode.Response{
				Success: true,
				Data:    item,
				TraceId: traceId,
			}); err != nil {
				return errors.Wrap(err, "encode error")
			}
			flushWriter.Flush()
			if item.FinishReason == "stop" {
				return nil
			}
		case <-time.After(time.Minute * 10):
			return nil
		}
	}
}
