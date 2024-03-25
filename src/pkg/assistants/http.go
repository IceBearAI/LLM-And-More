package assistants

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
			if toolId, ok := vars["toolId"]; ok && !strings.EqualFold(toolId, "") {
				ctx = context.WithValue(ctx, contextKeyToolId, toolId)
			}
			if assistantId, ok := vars["assistantId"]; ok && !strings.EqualFold(assistantId, "") {
				ctx = context.WithValue(ctx, contextKeyAssistantId, assistantId)
			}
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := MakeEndpoints(s, map[string][]endpoint.Middleware{
		"Assistants": ems,
	})

	r := mux.NewRouter()

	r.Handle("/{assistantId}/playground", kithttp.NewServer(
		eps.PlaygroundEndpoint,
		decodePlaygroundRequest,
		encodePlaygroundStreamResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/create", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{assistantId}", kithttp.NewServer(
		eps.GetEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{assistantId}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPut)
	r.Handle("/{assistantId}", kithttp.NewServer(
		eps.DeleteEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{assistantId}/tools", kithttp.NewServer(
		eps.ListToolEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodGet)
	r.Handle("/{assistantId}/tools", kithttp.NewServer(
		eps.AddToolEndpoint,
		decodeAddToolRequest,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	r.Handle("/{assistantId}/tools/{toolId}", kithttp.NewServer(
		eps.RemoveToolEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodDelete)
	r.Handle("/{assistantId}/publish", kithttp.NewServer(
		eps.PublishEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		kitopts...)).Methods(http.MethodPost)
	return r
}

func decodeAddToolRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req addToolRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req listRequest
	req.Name = r.URL.Query().Get("name")
	req.Page = page.GetPage(r)
	req.PageSize = page.GetPageSize(r)
	return req, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req updateRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req createRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err = validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodePlaygroundRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req playgroundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func encodePlaygroundStreamResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) (err error) {
	resp, ok := response.(encode.Response)
	if !ok {
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	stream, ok := resp.Data.(<-chan playgroundResult)
	if !ok {
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	if !resp.Stream {
		for {
			select {
			case item, ok := <-stream:
				if !ok {
					return nil
				}
				if item.FinishReason == "stop" {
					writer.Header().Set("TraceId", resp.TraceId)
					resp.Data = item
					return kithttp.EncodeJSONResponse(ctx, writer, resp)
				}
			case <-time.After(time.Minute * 10):
				return nil
			}
		}
	}
	flushWriter := writer.(http.Flusher)
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("aigc-model", "chat/stream")
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
				Code:    200,
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
	return
}
