package service

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	tiktoken2 "github.com/IceBearAI/aigc/src/helpers/tiktoken"
	"github.com/IceBearAI/aigc/src/logging"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/pkg/auth"
	"github.com/IceBearAI/aigc/src/pkg/chat"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pkoukk/tiktoken-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var (
	chatApi chat.Service
)

func startApiHttpServer(ctx context.Context, g *group.Group) {
	var clientOptions []kithttp.ClientOption
	if serverDebug {
		clientOptions = append(clientOptions, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			return ctx
		}),
			kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
				dump, _ := httputil.DumpResponse(response, true)
				fmt.Println(string(dump))
				return ctx
			}))
	}

	tiktoken.SetBpeLoader(tiktoken2.NewBpeLoader(DataFs))

	if logger != nil {
	}

	if tracer != nil {
	}

	apiGroup := &group.Group{}

	initApiHttpHandler(ctx, apiGroup)
	initCancelInterrupt(ctx, apiGroup)

	_ = level.Error(logger).Log("server exit", g.Run())
}

func initApiHttpHandler(ctx context.Context, g *group.Group) {
	httpLogger := log.With(logger, "openapi", "http")
	opts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encode.OpenAIErrorEncoder),
		kithttp.ServerErrorHandler(logging.NewLogErrorHandler(level.Error(logger), apiSvc)),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			authToken := request.Header.Get("Authorization")

			// 如果 Authorization 和 X-Token 都为空，尝试从 URL 查询参数获取 X-Token
			if authToken == "" {
				authToken = request.URL.Query().Get("Authorization")
			}

			// 更新请求头和上下文
			request.Header.Set("Authorization", authToken)
			ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, authToken)
			return ctx
		}),
	}

	if tracer != nil {
		opts = append(opts,
			kithttp.ServerBefore(
				opentracing.HTTPToContext(tracer, "HTTPToContext", logger),
				middleware.TracingServerBefore(tracer),
			))
	}

	if serverDebug {
		opts = append(opts, kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			return ctx
		}))
	}

	ems := []endpoint.Middleware{
		middleware.TracingMiddleware(tracer),                                                      // 2
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 1
	}

	r := mux.NewRouter()

	// auth模块
	r.PathPrefix("/v1/chat").Handler(http.StripPrefix("/v1/chat", auth.MakeHTTPHandler(authSvc, ems, opts)))
	// 对外metrics
	r.Handle("/metrics", promhttp.Handler())
	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})

	g.Add(func() error {
		_ = level.Debug(httpLogger).Log("transport", "HTTP", "openapi.addr", openApiAddr)
		return http.ListenAndServe(httpAddr, nil)
	}, func(e error) {
		closeConnection(ctx)
		_ = level.Error(httpLogger).Log("transport", "HTTP", "httpListener.Close", "http", "err", e)
		//if rdb != nil {
		//	_ = level.Debug(logger).Log("rdb", "close", "err", rdb.Close())
		//}
		os.Exit(1)
	})
}
