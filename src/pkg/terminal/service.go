package terminal

import (
	"context"
	"encoding/json"
	"github.com/IceBearAI/aigc/src/encode"
	aigcjwt "github.com/IceBearAI/aigc/src/jwt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/runtime"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/igm/sockjs-go/v3/sockjs"
	"net/http"
	"strings"
	"time"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts  []kithttp.ClientOption
	runtimePlatform string
	sessionTimeout  int64
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOpts returns a CreationOption that sets the http client options.
func WithHTTPClientOpts(opts ...kithttp.ClientOption) CreationOption {
	return func(co *CreationOptions) {
		co.httpClientOpts = opts
	}
}

// WithRuntimePlatform returns a CreationOption that sets the controller address.
func WithRuntimePlatform(platform string) CreationOption {
	return func(co *CreationOptions) {
		co.runtimePlatform = platform
	}
}

// WithSessionTimeout returns a CreationOption that sets the session timeout.
func WithSessionTimeout(timeout int64) CreationOption {
	return func(co *CreationOptions) {
		co.sessionTimeout = timeout
	}
}

type Middleware func(Service) Service

// Service is the interface for the terminal service.
type Service interface {
	// HandleTerminalSession handles terminal session
	HandleTerminalSession(session sockjs.Session)
	// Token 获取访问Terminal的token
	Token(ctx context.Context, tenantId, userId uint, resourceType string, serviceName, containerName string) (res tokenResult, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	apiSvc     services.Service
	options    *CreationOptions
}

func (s *service) Token(ctx context.Context, tenantId, userId uint, resourceType, serviceName, containerName string) (res tokenResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	timeout := time.Duration(s.options.sessionTimeout) * time.Second
	// 创建声明
	claims := aigcjwt.ArithmeticTerminalClaims{
		UserId:    userId,
		Container: containerName,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(timeout)),
			Issuer:    "system",
		},
	}
	e := func(ctx context.Context, i interface{}) (interface{}, error) { return ctx, nil }
	signer := kitjwt.NewSigner("terminal-auth", []byte(aigcjwt.GetJwtKey()), jwt2.SigningMethodHS256, claims)(e)
	kitCtx, err := signer(context.Background(), struct{}{})
	if err != nil {
		_ = level.Error(logger).Log("token", "SignedString", "err", err.Error())
		err = encode.ErrAuthTimeout.Wrap(err)
		return
	}
	tk, ok := kitCtx.(context.Context).Value(kitjwt.JWTContextKey).(string)
	if !ok {
		err = encode.ErrAuthTimeout.Error()
		return
	}
	var containers []string
	if resourceType == "deployment" {
		containers, err = s.apiSvc.Runtime().GetDeploymentContainerNames(ctx, serviceName)
		if err != nil {
			_ = level.Warn(logger).Log("apiSvc.Runtime", "GetDeploymentContainerNames", "err", err.Error())
		}
	} else {
		containers = []string{containerName}
	}

	res.SessionId = tk
	res.Container = containerName
	res.Containers = containers
	res.ServiceName = serviceName
	return res, nil
}

func (s *service) HandleTerminalSession(session sockjs.Session) {
	var (
		buf string
		err error
		msg runtime.Message
		//terminalSession Session
	)

	if buf, err = session.Recv(); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalSession", "can't Recv:", "err", err.Error())
		return
	}
	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalSession", "can't UnMarshal", "err", err.Error(), "buf", buf)
		return
	}

	if msg.Op != "bind" {
		_ = level.Error(s.logger).Log("handleTerminalSession: expected 'bind' message, got:", buf)
		return
	}

	var tr runtime.Result
	if err := json.Unmarshal([]byte(msg.Data), &tr); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalResult", "can't UnMarshal", "err", err.Error())
		return
	}
	ctx := context.Background()
	var userId uint
	// 验证token、权限、过期时间之类的
	userId, err = s.checkShellToken(tr.Cluster, tr.Namespace, tr.PodName, tr.Container, tr.SessionId)
	if err != nil {
		_ = level.Error(s.logger).Log("http.status", http.StatusBadRequest, "token", "not valid", "token", tr.Token, "err", err.Error())
		err = encode.ErrAuthTimeout.Wrap(err)
		return
	}

	_ = level.Debug(s.logger).Log("userId", userId, "cluster", tr.Cluster, "namespace", tr.Namespace, "podName",
		tr.PodName, "container", tr.Container, "cmd", tr.Cmd)

	ts := runtime.Session{
		Id:            tr.SessionId,
		SockJSSession: session,
		ClusterName:   tr.Cluster,
		PodName:       tr.PodName,
		Namespace:     tr.Namespace,
		Container:     tr.Container,
	}

	go s.apiSvc.Runtime().WaitForTerminal(ctx, ts, runtime.Config{
		ServiceName: tr.ServiceName,
	}, tr.Container, "bash")
	//terminalSession.sockJSSession = session
	//terminalSessions.Set(msg.SessionID, terminalSession)
	//terminalSession.bound <- nil
	return
}

func (s *service) checkShellToken(cluster, namespace, podName, container, token string) (userId uint, err error) {
	var atc aigcjwt.ArithmeticTerminalClaims
	e := func(ctx context.Context, i interface{}) (interface{}, error) { return ctx, nil }
	parser := kitjwt.NewParser(aigcjwt.JwtKeyFunc, jwt2.SigningMethodHS256, func() jwt2.Claims {
		return &atc
	})(e)
	ctx := context.WithValue(context.Background(), kitjwt.JWTContextKey, token)
	kitCtx, err := parser(ctx, struct{}{})
	if err != nil {
		_ = level.Error(s.logger).Log("jwt", "ParseWithClaims", "err", err)
		err = encode.ErrAuthTimeout.Wrap(err)
		return
	}
	claim, ok := kitCtx.(context.Context).Value(kitjwt.JWTClaimsContextKey).(*aigcjwt.ArithmeticTerminalClaims)
	if !ok {
		err = encode.ErrAuthTimeout.Error()
		return
	}

	if !strings.EqualFold(claim.Cluster, cluster) {
		err = encode.ErrAccountASD.Error()
		return
	}
	if !strings.EqualFold(claim.Namespace, namespace) {
		err = encode.ErrAccountASD.Error()
		return
	}
	if !strings.EqualFold(claim.PodName, podName) {
		err = encode.ErrAccountASD.Error()
		return
	}
	//if !strings.EqualFold(claim.Container, container) {
	//	return encode.ErrAccountASD.Error()
	//}

	// 拿到用户ID之后的操作

	return claim.UserId, nil
}

func New(logger log.Logger, traceId string, store repository.Repository, apiSvc services.Service, opts ...CreationOption) Service {
	options := &CreationOptions{
		sessionTimeout: 3600,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
		apiSvc:     apiSvc,
		options:    options,
	}
}
