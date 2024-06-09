package chat

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/sashabaranov/go-openai"
	"math"
)

type Endpoints struct {
	ChatCompletionStreamEndpoint endpoint.Endpoint
	CompletionEndpoint           endpoint.Endpoint
	ModelsEndpoint               endpoint.Endpoint
	EmbeddingsEndpoint           endpoint.Endpoint
}

func MakeEndpoints(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ChatCompletionStreamEndpoint: makeChatCompletionStreamEndpoint(s),
		CompletionEndpoint:           makeCompletionEndpoint(s),
		ModelsEndpoint:               makeModelsEndpoint(s),
		EmbeddingsEndpoint:           makeEmbeddingsEndpoint(s),
	}

	for _, m := range mdw["Chat"] {
		eps.ChatCompletionStreamEndpoint = m(eps.ChatCompletionStreamEndpoint)
		eps.CompletionEndpoint = m(eps.CompletionEndpoint)
		eps.ModelsEndpoint = m(eps.ModelsEndpoint)
		eps.EmbeddingsEndpoint = m(eps.EmbeddingsEndpoint)
	}
	return eps
}

func makeEmbeddingsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}

		req := request.(openai.EmbeddingRequest)
		res, err := s.Embeddings(ctx, channelId, req)
		if req.EncodingFormat == openai.EmbeddingEncodingFormatBase64 {
			var resData []string
			for _, v := range res.Data {
				for j, f := range v.Embedding {
					if math.IsNaN(float64(f)) {
						continue
					}
					if math.IsInf(float64(f), 0) {
						continue
					}
					v.Embedding[j] = float32(math.Round(float64(f)*10000) / 10000)
				}
				resData = append(resData, embeddingEncode(v.Embedding))
			}
			return encode.Response{
				Success: true,
				Code:    200,
				Data:    resData,
				Error:   err,
			}, err
		}
		return encode.Response{
			Success: true,
			Code:    200,
			Data:    res,
			Error:   err,
		}, err
	}
}

func makeModelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}

		res, err := s.Models(ctx, channelId)
		return encode.Response{
			Success: true,
			Code:    200,
			Data:    res,
			Error:   err,
		}, err
	}
}

func makeCompletionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}
		tenantId, _ := ctx.Value(ContextKeyTenantId).(uint)
		models, ok := ctx.Value(ContextKeyChannelModelsList).([]string)
		if !ok || models == nil {
			return nil, encode.ErrChatChannelNotFound.Error()
		}
		req := request.(openai.CompletionRequest)
		if tenantId != 1 {
			if !util.StringInArray(models, req.Model) {
				return nil, encode.ErrChatChannelModelIdNotAllow.Error()
			}
		}
		res, err := s.Completion(ctx, channelId, req)
		return encode.Response{
			Success: true,
			Code:    200,
			Data:    res,
			Error:   err,
		}, err
	}

}

func makeChatCompletionStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		channelId, ok := ctx.Value(ContextKeyChannelId).(uint)
		if !ok {
			return nil, encode.ErrChatChannelApiKey.Error()
		}
		tenantId, _ := ctx.Value(ContextKeyTenantId).(uint)
		models, ok := ctx.Value(ContextKeyChannelModelsList).([]string)
		if !ok || models == nil {
			return nil, encode.ErrChatChannelNotFound.Error()
		}
		req := request.(openai.ChatCompletionRequest)
		if tenantId != 1 {
			if !util.StringInArray(models, req.Model) {
				return nil, encode.ErrChatChannelModelIdNotAllow.Error()
			}
		}
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

// Encode 方法将 float32 列表编码为 base64 编码的字符串
func embeddingEncode(floats []float32) string {
	bytes := make([]byte, len(floats)*4)
	for i, f := range floats {
		binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], math.Float32bits(f))
	}
	encodedData := base64.StdEncoding.EncodeToString(bytes)
	return encodedData
}
