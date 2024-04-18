package terminal

import (
	"context"
	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	// tokenRequest 获取token请求
	tokenRequest struct {
		ResourceType  string `json:"resourceType"`
		ContainerName string `json:"containerName"`
		ServiceName   string `json:"serviceName"`
	}

	// tokenResult 获取token结果
	tokenResult struct {
		Namespace   string    `json:"namespace,omitempty"`
		PodName     string    `json:"podName,omitempty"`
		Container   string    `json:"container"`
		ErrMsg      string    `json:"errMsg"`
		SessionId   string    `json:"sessionId"`
		Token       string    `json:"token"`
		Cluster     string    `json:"cluster,omitempty"`
		Containers  []string  `json:"containers"`
		Phase       string    `json:"phase,omitempty"`
		HostIp      string    `json:"hostIp,omitempty"`
		PodIp       string    `json:"podIp,omitempty"`
		StartTime   time.Time `json:"startTime"`
		ServiceName string    `json:"serviceName,omitempty"`
		Pods        []string  `json:"pods,omitempty"`
	}
)

type Endpoints struct {
	TokenEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		TokenEndpoint: makeTokenEndpoint(s),
	}

	for _, m := range dmw["Token"] {
		eps.TokenEndpoint = m(eps.TokenEndpoint)
	}
	return eps
}

func makeTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tenantId, _ := middleware.GetTenantId(ctx)
		userId, _ := ctx.Value(middleware.ContextUserId).(uint)
		req := request.(tokenRequest)
		res, err := s.Token(ctx, tenantId, userId, req.ResourceType, req.ServiceName)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
