package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"time"
)

func MakeHTTPHandler(s Service, mdw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware
	ems = append(ems, mdw...)
	var kitopts = []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return ctx
		}),
	}
	kitopts = append(opts, kitopts...)

	eps := MakeEndpoints(s, map[string][]endpoint.Middleware{
		"Chat": ems,
	})

	r := mux.NewRouter()

	r.Handle("/completions", kithttp.NewServer(
		eps.ChatCompletionStreamEndpoint,
		decodeChatCompletionStreamRequest,
		encodeChatCompletionStreamResponse,
		kitopts...)).Methods(http.MethodPost)
	return r
}

func decodeChatCompletionStreamRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req openai.ChatCompletionRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil

}

func encodeChatCompletionStreamResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	resp, ok := response.(encode.Response)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	if headerer, ok := response.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				writer.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(kithttp.StatusCoder); ok {
		code = sc.StatusCode()
	}
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/octet-stream")
	traceId, _ := ctx.Value("traceId").(string)
	writer.Header().Set("TraceId", traceId)
	if code == http.StatusNoContent {
		return nil
	}
	stream, ok := resp.Data.(<-chan openai.ChatCompletionStreamResponse)
	if !ok {
		return encode.InvalidParams.Wrap(errors.New("invalid response type"))
	}
	flushWriter := writer.(http.Flusher)
	for {
		select {
		case item, ok := <-stream:
			if !ok {
				return nil
			}
			streamData, _ := json.Marshal(item)
			if resp.Stream {
				_, _ = writer.Write([]byte(fmt.Sprintf("data: %s\n\n", streamData)))
				flushWriter.Flush()
			} else {
				writer.Header().Set("Content-Type", "application/json")
				_, _ = writer.Write(streamData)
			}
		case <-time.After(time.Minute * 20):
			return nil
		}
	}
}