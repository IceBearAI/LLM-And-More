package datasettask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strings"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, mdw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, mdw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			vars := mux.Vars(request)
			if datasetId, ok := vars["datasetTaskId"]; ok && !strings.EqualFold(datasetId, "") {
				ctx = context.WithValue(ctx, contextKeyDatasetTaskId, datasetId)
			}
			if taskSegmentId, ok := vars["datasetTaskSegmentId"]; ok && !strings.EqualFold(taskSegmentId, "") {
				ctx = context.WithValue(ctx, contextKeyDatasetTaskSegmentId, taskSegmentId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := NewEndpoints(s, map[string][]endpoint.Middleware{
		"DatasetTask": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListTasksEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/create", kithttp.NewServer(
		eps.CreateTaskEndpoint,
		decodeCreateDatasetTaskRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetTaskId}/delete", kithttp.NewServer(
		eps.DeleteTaskEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{datasetTaskId}/segment/next", kithttp.NewServer(
		eps.GetTaskSegmentNextEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetTaskId}/clean", kithttp.NewServer(
		eps.CleanAnnotationTaskEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPut, http.MethodPost)
	r.Handle("/{datasetTaskId}/segment/{datasetTaskSegmentId}/mark", kithttp.NewServer(
		eps.AnnotationTaskSegmentEndpoint,
		decodeTaskSegmentMarkRequest,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetTaskId}/split", kithttp.NewServer(
		eps.SplitAnnotationDataSegmentEndpoint,
		decodeTaskSplitRequest,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetTaskId}/segment/{datasetTaskSegmentId}/abandoned", kithttp.NewServer(
		eps.AbandonTaskSegmentEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPut)
	r.Handle("/{datasetTaskId}/export", kithttp.NewServer(
		eps.ExportAnnotationDataEndpoint,
		decodeTaskExportRequest,
		encodeExportResponse, kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetTaskId}/detect/finish", kithttp.NewServer(
		eps.TaskDetectFinishEndpoint,
		decodeTaskDetectFinishRequest,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPut)
	r.Handle("/{datasetTaskId}/info", kithttp.NewServer(
		eps.TaskInfoEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetTaskId}/detect/cancel", kithttp.NewServer(
		eps.CancelCheckTaskDatasetSimilarEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPut, http.MethodPost)
	r.Handle("/{datasetTaskId}/detect/annotation/async", kithttp.NewServer(
		eps.AsyncCheckTaskDatasetSimilarEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetTaskId}/detect/annotation/log", kithttp.NewServer(
		eps.GetCheckTaskDatasetSimilarLogEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetTaskId}/segment/{datasetTaskSegmentId}/generation/annotation", kithttp.NewServer(
		eps.GenerationAnnotationContentEndpoint,
		decodeGenerationAnnotationRequest,
		encode.JsonResponse, kitopts...)).Methods(http.MethodPost)
	return r
}

func decodeGenerationAnnotationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req generationAnnotationContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskListRequest
	req.name = r.URL.Query().Get("name")
	req.page = page.GetPage(r)
	req.pageSize = page.GetPageSize(r)

	return req, nil
}

func decodeCreateDatasetTaskRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskCreateRequest
	// 限制上传文件大小 100M
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeTaskSegmentMarkRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskSegmentAnnotationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeTaskSplitRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskSplitAnnotationDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}

func decodeTaskExportRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskExportAnnotationDataRequest
	req.FormatType = r.URL.Query().Get("formatType")
	if req.FormatType == "" {
		req.FormatType = "default"
	}

	return req, nil
}

func encodeExportResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(encode.Response)
	if resp.Error != nil {
		encode.JsonError(ctx, resp.Error, w)
		return nil
	}
	filePath := resp.Data.(string)
	body, _ := os.ReadFile(filePath)
	defer func() {
		_ = os.Remove(filePath)
	}()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(filePath)))
	_, err := w.Write(body)
	return err
}

func decodeTaskDetectFinishRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req taskDetectFinishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}
