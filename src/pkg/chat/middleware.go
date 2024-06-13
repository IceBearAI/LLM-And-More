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
	// ContextKeyChannelModelsList is the context key that holds the channel models list.
	ContextKeyChannelModelsList contextKey = "ctx-channel-models-list"
	// ContextKeyChannelId is the context key that holds the channel id.
	ContextKeyChannelId contextKey = "ctx-channel-id"
	// ContextKeyChannelQuota is the context key that holds the channel.
	ContextKeyChannelQuota contextKey = "ctx-channel-quota"
	// ContextKeyTenantId is the context key that holds the tenant id.
	ContextKeyTenantId contextKey = "ctx-tenant-id"
)

func CheckChatMiddleware(store repository.Repository, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "CheckChatMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}
			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			if strings.HasPrefix(token, "Bearer ") {
				token = strings.TrimPrefix(token, "Bearer ")
			}

			if token == "" {
				return nil, encode.ErrChatChannelApiKey.Error()
			}
			channelInfo, err := store.Channel().FindChannelByKey(ctx, token, "ChannelModels")
			if err != nil {
				return nil, encode.ErrChatChannelApiKey.Wrap(err)
			}
			var modeNames []string
			for _, model := range channelInfo.ChannelModels {
				if model.BaseModelName != "" {
					modeNames = append(modeNames, model.BaseModelName)
				}
				modeNames = append(modeNames, model.ModelName)
			}
			ctx = context.WithValue(ctx, ContextKeyChannelId, channelInfo.ID)
			ctx = context.WithValue(ctx, ContextKeyTenantId, channelInfo.TenantId)
			ctx = context.WithValue(ctx, ContextKeyChannelModelsList, modeNames)
			ctx = context.WithValue(ctx, ContextKeyChannelQuota, channelInfo.Quota)
			return next(ctx, request)
		}
	}
}
