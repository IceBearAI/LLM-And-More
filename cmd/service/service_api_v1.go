package service

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/encode"
	tiktoken2 "github.com/IceBearAI/aigc/src/helpers/tiktoken"
	"github.com/IceBearAI/aigc/src/logging"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/pkg/chat"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pkoukk/tiktoken-go"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	servicesChat "github.com/IceBearAI/aigc/src/services/chat"
)

var (
	chatApi chat.Service

	apiV1StartCmd = &cobra.Command{
		Use:   "start-api",
		Short: "启动http api服务",
		Example: `## 启动命令
aigc-server start-api -p :8081
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startApiHttpServer(cmd.Context())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := prepare(cmd.Context()); err != nil {
				_ = level.Error(logger).Log("cmd", "start.PreRunE", "err", err.Error())
				return err
			}

			// 判断是否需要初始化数据，如果没有则初始化数据
			if !gormDB.Migrator().HasTable(types.ChatMessages{}) {
				_ = generateTable()
				if err = initData(); err != nil {
					_ = level.Error(logger).Log("cmd.start.PreRunE", "initData", "err", err.Error())
					return err
				}
			}
			return nil
		},
	}

	chatWorkerSvc servicesChat.WorkerService
)

func init() {
	apiV1StartCmd.PersistentFlags().StringVarP(&openApiAddr, "openapi.port", "p", ":8081", "服务启动的http api 端口")
	apiV1StartCmd.PersistentFlags().BoolVar(&webEmbed, "web.embed", true, "是否使用embed.FS")
}

func startApiHttpServer(ctx context.Context) error {
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
	chatWorkerSvc = servicesChat.NewFastChatWorker(
		servicesChat.WithControllerAddress(fsChatControllerAddress),
		//servicesChat.WithWorkerCreationOptionHTTPClientOpts(clientOptions...),
	)
	chatApi = chat.New(logger, "traceId", store, apiSvc,
		chat.WithWorkerService(chatWorkerSvc),
		chat.WithHTTPClientOpts(clientOptions...),
	)
	fieldKeys := []string{"method"}
	chatApi = chat.NewInstrumentingService(
		prometheus.NewCounterFrom(prometheus2.CounterOpts{
			Namespace: "chat",
			Subsystem: "service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(prometheus2.SummaryOpts{
			Namespace: "chat",
			Subsystem: "service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys))(chatApi)

	if logger != nil {
	}

	if tracer != nil {
	}

	apiGroup := &group.Group{}

	initApiHttpHandler(ctx, apiGroup)
	initCancelInterrupt(ctx, apiGroup)

	return apiGroup.Run()
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
		chat.CheckChatMiddleware(store, tracer),
		middleware.TracingMiddleware(tracer),                                                      // 2
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 1
	}

	r := mux.NewRouter()

	// auth模块
	r.PathPrefix("/v1/chat").Handler(http.StripPrefix("/v1/chat", chat.MakeHTTPHandler(chatApi, ems, opts)))
	// 对外metrics
	prometheus2.MustRegister(
		chat.NewChatQueueGaugeService(logger, chatWorkerSvc),
		chat.NewChatAvgSpeedGaugeService(logger, chatWorkerSvc),
	)
	r.Handle("/metrics", promhttp.Handler())
	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})

	http.Handle("/", accessControl(r, httpLogger))

	g.Add(func() error {
		_ = level.Debug(httpLogger).Log("transport", "HTTP", "openapi.addr", openApiAddr)
		return http.ListenAndServe(openApiAddr, nil)
	}, func(e error) {
		closeConnection(ctx)
		_ = level.Error(httpLogger).Log("transport", "HTTP", "httpListener.Close", "http", "err", e)
		//if rdb != nil {
		//	_ = level.Debug(logger).Log("rdb", "close", "err", rdb.Close())
		//}
		os.Exit(1)
	})
}
