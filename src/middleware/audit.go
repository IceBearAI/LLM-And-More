package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func AuditMiddleware(logger log.Logger, store repository.Repository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			l := log.With(logger, "middleware", "audit")
			begin := time.Now()
			method := ctx.Value(kithttp.ContextKeyRequestMethod).(string)

			if !strings.EqualFold(method, http.MethodPost) &&
				!strings.EqualFold(method, http.MethodDelete) &&
				!strings.EqualFold(method, http.MethodPatch) &&
				!strings.EqualFold(method, http.MethodPut) {
				return next(ctx, request)
			}

			defer func() {
				// 错误处理
				if err := recover(); err != nil {
					_ = level.Error(l).Log("err", err, "msg", "panic occurred in audit middleware")
				}

				uri := ctx.Value(kithttp.ContextKeyRequestURI).(string)
				u, err := url.Parse(uri)
				if err != nil {
					_ = level.Error(l).Log("msg", "parse url error", "err", err.Error())
					return
				}
				traceId := ctx.Value("traceId").(string)
				email, _ := GetEmail(ctx)
				tenantId, _ := GetTenantId(ctx)

				reqBody, err := json.Marshal(request)
				if err != nil {
					_ = level.Error(l).Log("msg", "marshal request error", "err", err.Error(), "request", fmt.Sprintf("%+v", request))
					return
				}
				reqBodyStr := string(reqBody)
				re := regexp.MustCompile(`"password":"[^"]*"`)
				reqBodyStr = re.ReplaceAllString(reqBodyStr, `"password":"******"`)
				var respBody []byte
				if err == nil {
					respBody, err = json.Marshal(response)
					if err != nil {
						_ = level.Error(l).Log("msg", "marshal response error", "err", err.Error(), "response", fmt.Sprintf("%+v", response))
						return
					}
				}

				errorMsg := ""
				if err != nil {
					errorMsg = err.Error()
				}
				repoErr := store.Sys().CreateAudit(ctx, &types.SysAudit{
					Operator:      email,
					TenantID:      tenantId,
					RequestMethod: method,
					RequestUrl:    u.String(),
					RequestBody:   reqBodyStr,
					ResponseBody:  string(respBody),
					IsError:       err != nil,
					ErrorMessage:  errorMsg,
					TraceID:       traceId,
					Duration:      time.Since(begin).Seconds(),
				})
				if repoErr != nil {
					_ = level.Error(l).Log("msg", "audit error", "err", repoErr.Error())
				}
			}()
			return next(ctx, request)
		}
	}
}
