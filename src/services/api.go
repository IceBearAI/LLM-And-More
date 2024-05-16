package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/services/chat"
	"github.com/IceBearAI/aigc/src/services/ldapcli"
	"github.com/IceBearAI/aigc/src/services/runtime"
	"github.com/IceBearAI/aigc/src/services/s3"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"net/http/httputil"
)

type Config struct {
	Namespace, ServiceName string
	Ldap                   ldapcli.Config
	StorageType            string
	Runtime                []runtime.CreationOption
	RuntimePlatform        string
	ChatOptions            []chat.CreationOption
	OnlyOpenAI             bool
}

type ProviderName string

const (
	ProviderOpenAI  ProviderName = "OpenAI"
	ProviderFsChat  ProviderName = "FastChat"
	ProviderLocalAI ProviderName = "LocalAI"
)

type ContextKey string

// Service 所有调用外部服务在API聚合
type Service interface {
	// Ldap ldap客户端
	Ldap() ldapcli.Service
	// Runtime runtime服务
	Runtime() runtime.Service
	// Chat chat服务
	Chat(providerName ProviderName) chat.Service
}

type api struct {
	logger     log.Logger
	s3Client   s3.Service
	traceId    string
	ldapSvc    ldapcli.Service
	runtimeSvc runtime.Service
	chatSvc    map[ProviderName]chat.Service
	onlyOpenAI bool
}

func (s *api) Chat(providerName ProviderName) chat.Service {
	if s.onlyOpenAI {
		providerName = ProviderOpenAI
	}
	if svc, ok := s.chatSvc[providerName]; ok {
		return svc
	}
	return s.chatSvc[ProviderOpenAI]
}

func (s *api) Runtime() runtime.Service {
	return s.runtimeSvc
}

func (s *api) Ldap() ldapcli.Service {
	return s.ldapSvc
}

func (s *api) S3Client(ctx context.Context) s3.Service {
	return s.s3Client
}

// NewApi 中间件有顺序,在后面的会最先执行
func NewApi(_ context.Context, logger log.Logger, traceId string, debug bool, tracer opentracing.Tracer, cfg *Config, opts []kithttp.ClientOption) Service {
	logger = log.With(logger, "api", "Api")
	if debug {
		opts = append(opts, kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			return ctx
		}),
			kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
				dump, _ := httputil.DumpResponse(response, true)
				fmt.Println(string(dump))
				return ctx
			}),
		)
	}

	ldapSvc := ldapcli.New(cfg.Ldap)
	runtimeSvc, err := runtime.New(cfg.RuntimePlatform, cfg.Runtime...)
	_ = level.Info(logger).Log("runtimePlatform", cfg.RuntimePlatform)
	if err != nil {
		_ = level.Error(logger).Log("runtime.New", "err", err.Error())
	}

	// 初始化chat服务
	chatSvc := make(map[ProviderName]chat.Service)
	chatSvc[ProviderOpenAI] = chat.NewOpenAI(cfg.ChatOptions...)
	chatSvc[ProviderFsChat] = chat.NewFsChatApi(cfg.ChatOptions...)

	if logger != nil {
		ldapSvc = ldapcli.NewLogging(logger, traceId)(ldapSvc)
		//s3Cli = s3.NewLogging(logger, traceId)(s3Cli)
		runtimeSvc = runtime.NewLogging(logger, traceId)(runtimeSvc)

		if debug {
			b, _ := json.Marshal(cfg.Ldap)
			_ = level.Debug(logger).Log("ldap.config", string(b))
			//b, _ = json.Marshal(cfg.S3)
			//_ = level.Debug(logger).Log("s3.config", string(b))
		}
		chatSvc[ProviderFsChat] = chat.NewLogging(logger, traceId, string(ProviderFsChat))(chatSvc[ProviderFsChat])
		chatSvc[ProviderOpenAI] = chat.NewLogging(logger, traceId, string(ProviderOpenAI))(chatSvc[ProviderOpenAI])
	}

	// 如果tracer有的话
	if tracer != nil {
		//s3Cli = s3.NewTracing(tracer)(s3Cli)
		ldapSvc = ldapcli.NewTracing(tracer)(ldapSvc)
		runtimeSvc = runtime.NewTracing(tracer)(runtimeSvc)
		chatSvc[ProviderFsChat] = chat.NewTracing(tracer)(chatSvc[ProviderFsChat])
		chatSvc[ProviderOpenAI] = chat.NewTracing(tracer)(chatSvc[ProviderOpenAI])
	}

	return &api{
		ldapSvc: ldapSvc,
		//s3Client:    s3Cli,
		runtimeSvc: runtimeSvc,
		chatSvc:    chatSvc,
		onlyOpenAI: cfg.OnlyOpenAI,
	}
}
