package middleware

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/IceBearAI/aigc/src/encode"
	asdjwt "github.com/IceBearAI/aigc/src/jwt"
)

type ASDContext string

const (
	ContextUserId               ASDContext = "ctx-user-id"          // 用户ID
	ContextKeyUserEmail         ASDContext = "ctx-user-email"       // 用户邮箱
	ContextKeyUseChannel        ASDContext = "ctx-use-channel"      // 使用渠道
	ContextKeyTenantId          ASDContext = "ctx-tenant-id"        // 租户ID
	ContextKeyPublicTenantId    ASDContext = "ctx-public-tenant-id" // 租户ID
	ContextKeyServiceIp         ASDContext = "ctx-service-ip"       // ContextKeyServiceIp 当前服务IP
	ContextKeyServicePort       ASDContext = "ctx-service-port"     // ContextKeyServicePort 当前服务端口
	ContextKeyChannelId         ASDContext = "ctx-channel-id"
	ContextKeyChannelModelsList ASDContext = "ctx-channel-models-list"
)

// CheckAuthMiddleware 验证用户登录信息，并将信息定入上下文
func CheckAuthMiddleware(logger log.Logger, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//email, ok := ctx.Value(ContextKeyUserEmail).(string)
			//if ok && !strings.EqualFold(email, "") {
			//	return next(ctx, request)
			//}
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

			// 判断是否是其他渠道 获取ChannelName
			//useChannel, ok := ctx.Value(ContextKeyUseChannel).(bool)
			//if ok && useChannel {
			//	ctx = context.WithValue(ctx, ContextKeyUserEmail, token)
			//	return next(ctx, request)
			//}

			u, _ := url.Parse(ctx.Value(kithttp.ContextKeyRequestURI).(string))
			if strings.EqualFold(token, "") && uriWhitelist(u.Path, uriArray) {
				// 跳过部分接口不进行验证
				return next(ctx, request)
			}

			if token == "" {
				_ = level.Warn(logger).Log("ctx", "Value", "err", encode.ErrAuthNotLogin.Error())
				return nil, encode.ErrAuthNotLogin.Error()
			}

			var clustom asdjwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &clustom, asdjwt.JwtKeyFunc)
			if err != nil || tk == nil {
				_ = level.Error(logger).Log("jwt", "ParseWithClaims", "err", err)
				err = encode.ErrAuthNotLogin.Wrap(err)
				return
			}

			claim, ok := tk.Claims.(*asdjwt.ArithmeticCustomClaims)
			if !ok {
				_ = level.Error(logger).Log("tk", "Claims", "err", ok)
				err = encode.ErrAccountASD.Error()
				return
			}

			if claim.Email == "" {
				return nil, encode.ErrAuthNotLogin.Error()
			}

			// 区分一下token 来源
			//_ = level.Info(logger).Log("tk", "Claims", "source", claim.Source, "QwUserid", claim.QwUserid, "UserId", claim.UserId, "Email", claim.Email)

			// 查询用户是否退出
			//if rdbTK := rdb.Get(ctx, fmt.Sprintf("login:%d:token", claim.UserId)).Val(); rdbTK == "" {
			//	_ = level.Error(logger).Log("cache", "Get", "key", fmt.Sprintf("login:%d:token", claim.UserId))
			//	err = encode.ErrAuthNotLogin.Error()
			//	return nil, err
			//}

			ctx = context.WithValue(ctx, ContextKeyUserEmail, claim.Email)
			ctx = context.WithValue(ctx, ContextUserId, claim.UserId)
			ctx = context.WithValue(ctx, "Authorization", token)
			return next(ctx, request)
		}
	}
}

func keyMatch3(key1 string, key2 string) bool {
	re := regexp.MustCompile(`(.*)\{[^/]+\}(.*)`)
	for {
		if !strings.Contains(key2, "/{") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	}
	return regexMatch(key1, key2)
}

func regexMatch(key1 string, key2 string) bool {
	if !strings.Contains(key2, "[^/]") && !strings.EqualFold(key2, key1) {
		return false
	}
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}

var uriArray = []string{
	"/api/auth/login",
}

var tenantUriArray = []string{
	"/api/auth/login",
	"/api/auth/account",
}

func uriWhitelist(uri string, uriArray []string) bool {
	if strings.EqualFold(uri, "") {
		return true
	}
	if strings.EqualFold(uri, "/") {
		return true
	}
	for _, v := range uriArray {
		if strings.EqualFold(v, uri) {
			return true
		}
		if keyMatch3(uri, v) {
			return true
		}
	}
	return strings.Contains(uri, "metrics")
}

func CheckTenantMiddleware(logger log.Logger, store repository.Repository, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			publicTenantId := ctx.Value(ContextKeyPublicTenantId).(string)

			u, _ := url.Parse(ctx.Value(kithttp.ContextKeyRequestURI).(string))
			if strings.EqualFold(publicTenantId, "") && uriWhitelist(u.Path, tenantUriArray) {
				// 跳过部分接口不进行验证
				return next(ctx, request)
			}

			if publicTenantId == "" {
				_ = level.Warn(logger).Log("ctx", "Value", "err", encode.ErrAuthNotLogin.Error())
				return nil, encode.ErrTenantNotFound.Error()
			}

			res, err := store.Auth().GetTenantByUuid(ctx, publicTenantId)
			if err != nil {
				_ = level.Warn(logger).Log("repository.Chat", "GetTenantByUuid", "err", err.Error())
				return nil, encode.ErrTenantNotFound.Error()
			}
			ctx = context.WithValue(ctx, ContextKeyTenantId, res.ID)
			return next(ctx, request)
		}
	}
}

func GetTenantId(ctx context.Context) (uint, bool) {
	tenantId, ok := ctx.Value(ContextKeyTenantId).(uint)
	return tenantId, ok
}

func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyUserEmail).(string)
	return email, ok
}
