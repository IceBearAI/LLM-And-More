package datasets

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
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

var validate = validator.New()

func MakeHTTPHandler(s Service, mdw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, mdw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			vars := mux.Vars(request)
			if datasetId, ok := vars["datasetId"]; ok && !strings.EqualFold(datasetId, "") {
				ctx = context.WithValue(ctx, contextKeyDatasetId, datasetId)
			}
			if datasetSampleId, ok := vars["datasetSampleId"]; ok && !strings.EqualFold(datasetSampleId, "") {
				ctx = context.WithValue(ctx, contextKeyDatasetSampleId, datasetSampleId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := MakeEndpoints(s, map[string][]endpoint.Middleware{
		"Dataset": ems,
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
	r.Handle("/{datasetId}", kithttp.NewServer(
		eps.DetailEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetId}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPut)
	r.Handle("/{datasetId}", kithttp.NewServer(
		eps.DeleteEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{datasetId}/samples", kithttp.NewServer(
		eps.SampleListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{datasetId}/samples", kithttp.NewServer(
		eps.AddSampleEndpoint,
		decodeAddSampleRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetId}/samples/{datasetSampleId}", kithttp.NewServer(
		eps.UpdateSampleEndpoint,
		decodeAddSampleRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPut)
	r.Handle("/{datasetId}/samples/{datasetSampleId}", kithttp.NewServer(
		eps.DeleteSampleEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{datasetId}/export/samples", kithttp.NewServer(
		eps.ExportSampleEndpoint,
		decodeExportRequest,
		encodeExportSampleResponse,
		kitopts...)).Methods(http.MethodGet, http.MethodPost)
	return r
}

func decodeExportRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req exportSampleRequest
	req.Format = "jsonl"
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.Page = page.GetPage(r)
	req.PageSize = page.GetPageSize(r)
	req.Name = r.URL.Query().Get("name")
	return req, nil
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req datasetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeAddSampleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addSampleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if len(req.Messages) < 2 {
		return nil, encode.InvalidParams.Wrap(errors.New("messages length must be greater than 2"))
	}
	return req, nil
}

func encodeExportSampleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	resp := response.(encode.Response)
	if resp.Error != nil {
		return resp.Error
	}
	var writeBody []byte
	messages := resp.Data.([]addSampleRequest)
	for _, v := range messages {
		b, _ := json.Marshal(v)
		writeBody = append(writeBody, []byte(fmt.Sprintf("%s\n", b))...)
	}

	fileName := fmt.Sprintf("ft-%s.jsonl", time.Now().Format("20060102150405"))

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	_, _ = w.Write(writeBody)
	return nil
}
