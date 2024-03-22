package datasetdocument

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
			if datasetId, ok := vars["datasetDocumentId"]; ok && !strings.EqualFold(datasetId, "") {
				ctx = context.WithValue(ctx, contextKeyDatasetDocumentId, datasetId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := MakeEndpoints(s, map[string][]endpoint.Middleware{
		"DatasetDocument": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListDocumentsEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/create", kithttp.NewServer(
		eps.CreateDocumentEndpoint,
		decodeCreateDatasetDocumentRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/{datasetDocumentId}", kithttp.NewServer(
		eps.DeleteDocumentEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse, kitopts...)).Methods(http.MethodDelete)
	return r
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.name = r.URL.Query().Get("name")
	req.page = page.GetPage(r)
	req.pageSize = page.GetPageSize(r)

	return req, nil
}

func decodeCreateDatasetDocumentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req documentCreateRequest
	// 限制上传文件大小 100M
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		return nil, encode.InvalidParams.Wrap(errors.New("文件大小超过100MB限制"))
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req.File = file
	req.FileHeader = header
	req.Name = r.FormValue("name")
	req.Remark = r.FormValue("remark")
	req.FormatType = r.FormValue("formatType")
	req.SplitType = r.FormValue("splitType")
	req.SplitMax = 0
	if strings.EqualFold(req.FormatType, "") {
		req.FormatType = "txt"
	}

	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	return req, nil
}
