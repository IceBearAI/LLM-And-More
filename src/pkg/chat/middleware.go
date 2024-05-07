package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
)

type Middleware func(Service) Service

type contextKey string

const (
	// ContextKeyChannelId channel id
	ContextKeyChannelId contextKey = "channel-id"
)

func CheckChatMiddleware(store repository.Repository, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "CheckAuthMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}
			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			token = strings.ReplaceAll(token, "Bearer ", "")

			if token == "" {
				return nil, encode.ErrChatChannelApiKey.Error()
			}
			channelInfo, err := store.Channel().FindChannelByKey(ctx, token)
			if err != nil {
				return nil, encode.ErrChatChannelApiKey.Wrap(err)
			}
			ctx = context.WithValue(ctx, ContextKeyChannelId, channelInfo.ID)
			return next(ctx, request)
		}
	}
}
