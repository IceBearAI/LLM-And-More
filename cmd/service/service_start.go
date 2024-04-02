package service

import (
	"context"
	"embed"
	"fmt"
	tiktoken2 "github.com/IceBearAI/aigc/src/helpers/tiktoken"
	"github.com/IceBearAI/aigc/src/pkg/assistants"
	"github.com/IceBearAI/aigc/src/pkg/auth"
	"github.com/IceBearAI/aigc/src/pkg/channels"
	"github.com/IceBearAI/aigc/src/pkg/datasetdocument"
	"github.com/IceBearAI/aigc/src/pkg/datasets"
	"github.com/IceBearAI/aigc/src/pkg/datasettask"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/pkg/finetuning"
	"github.com/IceBearAI/aigc/src/pkg/modelevaluate"
	"github.com/IceBearAI/aigc/src/pkg/models"
	"github.com/IceBearAI/aigc/src/pkg/sys"
	"github.com/IceBearAI/aigc/src/pkg/tools"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/pkoukk/tiktoken-go"
	"github.com/tmc/langchaingo/llms/openai"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/tracing/opentracing"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/logging"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动http服务",
		Example: `## 启动命令
aigc-server start -p :8080
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := prepare(cmd.Context()); err != nil {
				_ = level.Error(logger).Log("cmd", "start.PreRunE", "err", err.Error())
				return err
			}

			// 判断是否需要初始化数据，如果没有则初始化数据
			if !gormDB.Migrator().HasTable(types.Accounts{}) {
				_ = generateTable()
				if err = initData(); err != nil {
					_ = level.Error(logger).Log("cmd.start.PreRunE", "initData", "err", err.Error())
					return err
				}
			}
			_ = generateTable()

			//aigc-server channelID
			channelRes, err := store.Chat().FindChannelByApiKey(cmd.Context(), serverChannelKey)
			if err != nil {
				_ = level.Error(logger).Log("cmd.start.PreRunE", "FindChannelByApiKey", "err", err.Error())
				return err
			}
			channelId = int(channelRes.ID)
			return nil
		},
	}

	tracer stdopentracing.Tracer

	opts []kithttp.ServerOption

	WebFs  embed.FS
	DataFs embed.FS

	authSvc auth.Service

	fileSvc            files.Service
	channelSvc         channels.Service
	modelSvc           models.Service
	fineTuningSvc      finetuning.Service
	sysSvc             sys.Service
	datasetSvc         datasets.Service
	datasetDocumentSvc datasetdocument.Service
	datasetTaskSvc     datasettask.Service
	toolsSvc           tools.Service
	assistantsSvc      assistants.Service
	modelEvaluateSvc   modelevaluate.Service
)

func start(ctx context.Context) (err error) {

	tiktoken.SetBpeLoader(tiktoken2.NewBpeLoader(DataFs))

	authSvc = auth.New(logger, traceId, store, apiSvc)
	fileSvc = files.NewService(logger, traceId, store, apiSvc, []files.CreationOption{
		files.WithLocalDataPath(serverStoragePath),
		files.WithServerUrl(fmt.Sprintf("%s/storage", serverDomain)),
		files.WithStorageType("local"),
	}...)
	channelSvc = channels.NewService(logger, traceId, store, apiSvc)
	modelSvc = models.NewService(logger, traceId, store, apiSvc,
		models.WithGPUTolerationValue(datasetsGpuToleration),
		models.WithVolumeName(runtimeK8sVolumeName),
		models.WithControllerAddress(fsChatControllerAddress),
	)
	fineTuningSvc = finetuning.New(traceId, logger, store, fileSvc, apiSvc,
		finetuning.WithGpuTolerationValue(datasetsGpuToleration),
		finetuning.WithCallbackHost(serverDomain),
		finetuning.WithVolumeName(runtimeK8sVolumeName),
	)
	sysSvc = sys.NewService(logger, traceId, store, apiSvc)
	datasetSvc = datasets.New(logger, traceId, store)
	toolsSvc = tools.New(logger, traceId, store)
	assistantsSvc = assistants.New(logger, traceId, store, []kithttp.ClientOption{
		//kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
		//	dump, _ := httputil.DumpRequest(request, true)
		//	fmt.Println(string(dump))
		//	return ctx
		//}),
		//kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
		//	dump, _ := httputil.DumpResponse(response, true)
		//	fmt.Println(string(dump))
		//	return ctx
		//}),
	}, []openai.Option{
		openai.WithToken(serviceLocalAiToken),
		openai.WithBaseURL(serviceLocalAiHost),
	})
	datasetDocumentSvc = datasetdocument.New(traceId, logger, store)
	datasetTaskSvc = datasettask.New(traceId, logger, store, apiSvc, fileSvc,
		datasettask.WithDatasetImage(datasetsImage),
		datasettask.WithDatasetModel(datasetsModelName),
		datasettask.WithDatasetDrive(datasetsDevice),
		datasettask.WithCallbackHost(serverDomain),
		datasettask.WithVolumeName(runtimeK8sVolumeName),
	)

	modelEvaluateSvc = modelevaluate.New(logger, traceId, store, apiSvc, fileSvc,
		modelevaluate.WithDatasetGpuTolerationValue(datasetsGpuToleration),
		modelevaluate.WithCallbackHost(serverDomain),
		modelevaluate.WithVolumeName(runtimeK8sVolumeName),
	)

	if logger != nil {
		authSvc = auth.NewLogging(logger, logging.TraceId)(authSvc)
		fileSvc = files.NewLogging(logger, logging.TraceId)(fileSvc)
		channelSvc = channels.NewLogging(logger, logging.TraceId)(channelSvc)
		modelSvc = models.NewLogging(logger, logging.TraceId)(modelSvc)
		fineTuningSvc = finetuning.NewLogging(logger, logging.TraceId)(fineTuningSvc)
		sysSvc = sys.NewLogging(logger, logging.TraceId)(sysSvc)
		datasetSvc = datasets.NewLogging(logger, logging.TraceId)(datasetSvc)
		toolsSvc = tools.NewLogging(logger, logging.TraceId)(toolsSvc)
		datasetDocumentSvc = datasetdocument.NewLogging(logger, logging.TraceId)(datasetDocumentSvc)
		datasetTaskSvc = datasettask.NewLogging(logger, logging.TraceId)(datasetTaskSvc)
		modelEvaluateSvc = modelevaluate.NewLogging(logger, logging.TraceId)(modelEvaluateSvc)
	}

	if tracer != nil {
		authSvc = auth.NewTracing(tracer)(authSvc)
		fileSvc = files.NewTracing(tracer)(fileSvc)
		channelSvc = channels.NewTracing(tracer)(channelSvc)
		modelSvc = models.NewTracing(tracer)(modelSvc)
		fineTuningSvc = finetuning.NewTracing(tracer)(fineTuningSvc)
		sysSvc = sys.NewTracing(tracer)(sysSvc)
		datasetSvc = datasets.NewTracing(tracer)(datasetSvc)
		datasetDocumentSvc = datasetdocument.NewTracing(tracer)(datasetDocumentSvc)
		toolsSvc = tools.NewTracing(tracer)(toolsSvc)
		datasetTaskSvc = datasettask.NewTracing(tracer)(datasetTaskSvc)
		modelEvaluateSvc = modelevaluate.NewTracing(tracer)(modelEvaluateSvc)
	}

	g := &group.Group{}

	initHttpHandler(ctx, g)
	//initGRPCHandler(g)
	initCancelInterrupt(ctx, g)

	if cronJobAuto {
		autoCronjobHandler(ctx, g)
	}

	_ = level.Error(logger).Log("server exit", g.Run())
	return nil
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range corsHeaders {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}

func initHttpHandler(ctx context.Context, g *group.Group) {
	httpLogger := log.With(logger, "component", "http")

	opts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encode.JsonError),
		kithttp.ServerErrorHandler(logging.NewLogErrorHandler(level.Error(logger), apiSvc)),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			requestID := request.Header.Get("X-Request-Id")
			authToken := request.Header.Get("Authorization")
			xToken := request.Header.Get("X-Token")
			tenantID := request.Header.Get("X-Tenant-Id")

			// 如果 Authorization 和 X-Token 都为空，尝试从 URL 查询参数获取 X-Token
			if authToken == "" && xToken == "" {
				xToken = request.URL.Query().Get("X-Token")
			}

			// 优先使用 Authorization，如果为空，则使用 X-Token
			token := authToken
			if authToken == "" {
				token = xToken
			}

			// 如果 X-Tenant-Id 为空，尝试从 URL 查询参数获取
			if tenantID == "" {
				tenantID = request.URL.Query().Get("X-Tenant-Id")
			}

			// 更新请求头和上下文
			request.Header.Set("Authorization", token)
			ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, token)
			ctx = context.WithValue(ctx, logging.TraceId, requestID)
			ctx = context.WithValue(ctx, middleware.ContextKeyPublicTenantId, tenantID)
			// 假设 channelId 已经在之前定义
			ctx = context.WithValue(ctx, middleware.ContextKeyChannelId, channelId)
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

	authEms := []endpoint.Middleware{
		middleware.AuditMiddleware(logger, store),
		middleware.CheckTenantMiddleware(logger, store, tracer),
		middleware.CheckAuthMiddleware(logger, tracer),
	}
	authEms = append(authEms, ems...)

	r := mux.NewRouter()
	// auth模块
	r.PathPrefix("/api/auth").Handler(http.StripPrefix("/api/auth", auth.MakeHTTPHandler(authSvc, authEms, opts)))

	// file模块
	r.PathPrefix("/api/files").Handler(http.StripPrefix("/api", files.MakeHTTPHandler(fileSvc, authEms, opts)))
	// channel模块
	r.PathPrefix("/api/channels").Handler(http.StripPrefix("/api", channels.MakeHTTPHandler(channelSvc, authEms, opts)))
	// Model模块
	r.PathPrefix("/api/models").Handler(http.StripPrefix("/api", models.MakeHTTPHandler(modelSvc, authEms, opts)))
	// FineTuning模块
	r.PathPrefix("/api/finetuning").Handler(http.StripPrefix("/api", finetuning.MakeHTTPHandler(fineTuningSvc, authEms, opts)))
	// Sys模块
	r.PathPrefix("/api/sys").Handler(http.StripPrefix("/api/sys", sys.MakeHTTPHandler(sysSvc, authEms, opts)))
	// Dataset模块
	r.PathPrefix("/api/datasets").Handler(http.StripPrefix("/api/datasets", datasets.MakeHTTPHandler(datasetSvc, authEms, opts)))
	// Tools模块
	r.PathPrefix("/api/tools").Handler(http.StripPrefix("/api/tools", tools.MakeHTTPHandler(toolsSvc, authEms, opts)))
	// Assistants模块
	r.PathPrefix("/api/assistants").Handler(http.StripPrefix("/api/assistants", assistants.MakeHTTPHandler(assistantsSvc, authEms, opts)))
	// 数据集样本模块
	r.PathPrefix("/api/mgr/datasets").Handler(http.StripPrefix("/api/mgr/datasets", datasetdocument.MakeHTTPHandler(datasetDocumentSvc, authEms, opts)))
	// 数据集标注模块
	r.PathPrefix("/api/mgr/annotation/task").Handler(http.StripPrefix("/api/mgr/annotation/task", datasettask.MakeHTTPHandler(datasetTaskSvc, authEms, opts)))
	// Model Evaluate模块
	r.PathPrefix("/api/evaluate").Handler(http.StripPrefix("/api/evaluate", modelevaluate.MakeHTTPHandler(modelEvaluateSvc, authEms, opts)))
	// 对外metrics
	r.Handle("/metrics", promhttp.Handler())
	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	// 文件存储
	r.PathPrefix("/storage/").Handler(http.StripPrefix("/storage/", http.FileServer(http.Dir(serverStoragePath))))

	// web页面
	if webEmbed {
		fe, fsErr := fs.Sub(WebFs, DefaultWebPath)
		if fsErr != nil {
			_ = level.Error(logger).Log("FailedToSubPath", "web", "err", fsErr.Error())
		}
		r.PathPrefix("/").Handler(http.FileServer(http.FS(fe)))
	} else {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(webPath)))
	}

	if enableCORS {
		corsHeaders["Access-Control-Allow-Origin"] = corsAllowOrigins
		corsHeaders["Access-Control-Allow-Methods"] = corsAllowMethods
		corsHeaders["Access-Control-Allow-Headers"] = corsAllowHeaders
		corsHeaders["Access-Control-credentials"] = strconv.FormatBool(corsAllowCredentials)
	}

	http.Handle("/", accessControl(r, httpLogger))

	g.Add(func() error {
		_ = level.Debug(httpLogger).Log("transport", "HTTP", "addr", httpAddr)
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

func initCancelInterrupt(ctx context.Context, g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case sig := <-c:
			if err != nil {
				_ = level.Error(logger).Log("rocketmq", "close", "err", err)
				return err
			}
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(err error) {
		close(cancelInterrupt)
	})
}

func autoCronjobHandler(ctx context.Context, g *group.Group) {
	g.Add(func() error {
		return cronStart(ctx, cronJobNames)
	}, func(err2 error) {
		closeConnection(ctx)
		_ = level.Warn(logger).Log("")
	})
}

var localAddr string

func getLocalAddr() string {
	if localAddr != "" {
		return localAddr
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localAddr = ipNet.IP.String()
				return localAddr
			}
		}
	}

	return ""
}
