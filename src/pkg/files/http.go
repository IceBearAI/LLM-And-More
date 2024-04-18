package files

import (
	"bufio"
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

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
		"File": ems,
	})

	r := mux.NewRouter()
	r.Handle("/files", kithttp.NewServer(
		eps.CreateFileEndpoint,
		decodeCreateFileRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/files", kithttp.NewServer(
		eps.ListFilesEndpoint,
		decodeListFilesRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/files/{fileId}", kithttp.NewServer(
		eps.GetFileEndpoint,
		decodeGetFileRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/files/{fileId}", kithttp.NewServer(
		eps.DeleteFileEndpoint,
		decodeGetFileRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	return r
}

func decodeCreateFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// 限制上传文件大小 1000M
	if err := r.ParseMultipartForm(4000 << 20); err != nil {
		return nil, encode.InvalidParams.Wrap(errors.New("文件大小超过4000MB限制"))
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	purpose := r.PostFormValue("purpose")
	line := 0
	tokens := 0
	if purpose == types.FilePurposeFineTune.String() || purpose == types.FilePurposeFineTuneEval.String() {
		// 检查是否为JSONL格式
		maxCapacity := 1024 * 1024          // 1MB
		buf := make([]byte, 0, maxCapacity) // maxCapacity 是你希望设置的新的缓冲区大小
		scanner := bufio.NewScanner(file)
		scanner.Buffer(buf, maxCapacity)
		enc, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
		if err != nil {
			return nil, encode.ErrSystem.Wrap(errors.New("统计微调数据tokens错误"))
		}
		for scanner.Scan() {
			//var data MessagesWrapper
			//if err = json.Unmarshal(scanner.Bytes(), &data); err != nil {
			//	return nil, encode.InvalidParams.Wrap(errors.New("微调数据格式错误，仅支持jsonl格式"))
			//}
			tokens += len(enc.Encode(string(scanner.Bytes()), nil, nil))
			line++
		}
		if err = scanner.Err(); err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, encode.ErrSystem.Wrap(err)
		}
	}

	req := FileRequest{
		Purpose:    purpose,
		Header:     header,
		File:       file,
		FileType:   strings.TrimLeft(filepath.Ext(header.Filename), "."), // 去掉文件后缀的点
		LineCount:  line,
		TokenCount: tokens,
	}
	return req, nil
}

func decodeListFilesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListFileRequest{
		Purpose:  r.URL.Query().Get("purpose"),
		FileName: r.URL.Query().Get("fileName"),
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	tenantId, _ := middleware.GetTenantId(ctx)
	req.TenantId = tenantId
	return req, nil
}

func decodeGetFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	if fileId == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("fileId is required"))
	}
	req := GetFileRequest{
		FileId: fileId,
	}
	return req, nil
}
