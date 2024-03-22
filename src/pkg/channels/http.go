package channels

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
		"Channel": ems,
	})

	r := mux.NewRouter()
	r.Handle("/channels", kithttp.NewServer(
		eps.CreateChannelEndpoint,
		decodeCreateChannelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPost)
	r.Handle("/channels", kithttp.NewServer(
		eps.ListChannelsEndpoint,
		decodeListChannelsRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/channels/{id}", kithttp.NewServer(
		eps.UpdateChannelEndpoint,
		decodeUpdateChannelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodPut)
	r.Handle("/channels/{id}", kithttp.NewServer(
		eps.DeleteChannelEndpoint,
		decodeDeleteChannelRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodDelete)
	r.Handle("/channels/models", kithttp.NewServer(
		eps.ListChannelModelsEndpoint,
		decodeListChannelModelsRequest,
		encode.JsonResponse,
		kitopts...,
	)).Methods(http.MethodGet)
	r.Handle("/channels/chat/completions", kithttp.NewServer(
		eps.ChatCompletionStreamEndpoint,
		decodeChatCompletionStreamRequest,
		encodeChatCompletionsStreamResponse,
		kitopts...))
	return r
}

func decodeCreateChannelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err := validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}

func decodeUpdateChannelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	channelId := vars["id"]
	if channelId == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("channelId is empty"))
	}
	id, err := strconv.Atoi(channelId)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	var req UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	err = validate.Struct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req.Id = uint(id)
	return req, nil
}

func decodeListChannelsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListChannelRequest{
		Page:     page.GetPage(r),
		PageSize: page.GetPageSize(r),
	}
	name := r.URL.Query().Get("name")
	if name != "" {
		req.Name = &name
	}
	email := r.URL.Query().Get("email")
	if email != "" {
		req.Email = &email
	}
	projectName := r.URL.Query().Get("projectName")
	if projectName != "" {
		req.ProjectName = &projectName
	}
	serviceName := r.URL.Query().Get("serviceName")
	if serviceName != "" {
		req.ServiceName = &serviceName
	}
	return req, nil
}

func decodeDeleteChannelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	channelId := vars["id"]
	if channelId == "" {
		return nil, encode.InvalidParams.Wrap(errors.New("channelId is empty"))
	}
	id, err := strconv.Atoi(channelId)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req := IdRequest{
		Id: uint(id),
	}
	return req, nil
}

func decodeListChannelModelsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := ListChannelModelsRequest{}
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
