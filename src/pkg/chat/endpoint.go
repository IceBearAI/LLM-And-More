package chat

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/go-kit/kit/endpoint"
	"github.com/sashabaranov/go-openai"
)

type Endpoints struct {
	ChatCompletionStreamEndpoint endpoint.Endpoint
	CompletionEndpoint           endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ChatCompletionStreamEndpoint: makeChatCompletionStreamEndpoint(s),
		CompletionEndpoint:           makeCompletionEndpoint(s),
	}

	for _, m := range mdw["Chat"] {
		eps.ChatCompletionStreamEndpoint = m(eps.ChatCompletionStreamEndpoint)
		eps.CompletionEndpoint = m(eps.CompletionEndpoint)
	}
	return eps
}

func makeCompletionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}
		req := request.(openai.CompletionRequest)

		fmt.Println(channelId)
		return encode.Response{
			Success: true,
			Code:    200,
			Data:    nil,
			Error:   err,
			Stream:  req.Stream,
		}, err
	}

}

func makeChatCompletionStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}
		req := request.(openai.ChatCompletionRequest)
		var res interface{}
		if req.Stream {
			res, err = s.ChatCompletionStream(ctx, channelId, req)
		} else {
			res, err = s.ChatCompletion(ctx, channelId, req)
		}

		return encode.Response{
			Success: true,
			Code:    200,
			Data:    res,
			Error:   err,
			Stream:  req.Stream,
		}, err
	}
}
