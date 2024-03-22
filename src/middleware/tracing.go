package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

// TracingServerBefore Http serverBefore Tracing中间件
func TracingServerBefore(tracer stdopentracing.Tracer) kithttp.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		if tracer == nil {
			return ctx
		}
		reqPath := ctx.Value(kithttp.ContextKeyRequestURI).(string)
		u, _ := url.Parse(reqPath)
		span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, u.Path, stdopentracing.Tag{
			Key:   string(ext.Component),
			Value: "TracingServerBefore",
		}, stdopentracing.Tag{
			Key:   string(ext.HTTPMethod),
			Value: ctx.Value(kithttp.ContextKeyRequestMethod),
		}, stdopentracing.Tag{
			Key:   string(ext.HTTPUrl),
			Value: ctx.Value(kithttp.ContextKeyRequestURI),
		}, stdopentracing.Tag{
			Key:   "user_agent",
			Value: ctx.Value(kithttp.ContextKeyRequestUserAgent),
		}, stdopentracing.Tag{
			Key:   "x-request-id",
			Value: ctx.Value(kithttp.ContextKeyRequestXRequestID),
		}, stdopentracing.Tag{
			Key:   "Authorization",
			Value: ctx.Value(kithttp.ContextKeyRequestAuthorization),
		}, stdopentracing.Tag{
			Key:   "referer",
			Value: ctx.Value(kithttp.ContextKeyRequestReferer),
		}, stdopentracing.Tag{
			Key:   "x-forwarded-for",
			Value: ctx.Value(kithttp.ContextKeyRequestXForwardedFor),
		}, stdopentracing.Tag{
			Key:   "AIGC-User-Id",
			Value: request.Header.Get("AIGC-User-Id"),
		}, stdopentracing.Tag{
			Key:   "AIGC-User-Email",
			Value: request.Header.Get("AIGC-User-Email"),
		}, stdopentracing.Tag{
			Key:   "AIGC-Username",
			Value: request.Header.Get("AIGC-Username"),
		}, stdopentracing.Tag{
			Key:   "AIGC-User-Roles",
			Value: request.Header.Get("AIGC-User-Roles"),
		}, stdopentracing.Tag{
			Key:   "AIGC-Sign",
			Value: request.Header.Get("AIGC-Sign"),
		}, stdopentracing.Tag{
			Key:   "AIGC-Super-Admin",
			Value: request.Header.Get("AIGC-Super-Admin"),
		}, stdopentracing.Tag{
			Key:   "X-Tenant-Id",
			Value: request.Header.Get("X-Tenant-Id"),
		})
		traceId := span.Context().(jaeger.SpanContext).TraceID().String()
		ctx = context.WithValue(ctx, "traceId", traceId)
		ctx = context.WithValue(ctx, "trace-id", traceId)
		//span = span.SetTag("TraceId", traceId)

		defer func() {
			b, _ := json.Marshal(request)
			span.LogKV("request", string(b))
			span.Finish()
		}()
		return ctx
	}
}

// TracingGrpcServerBefore grpc serverBefore Tracing中间件
func TracingGrpcServerBefore(tracer stdopentracing.Tracer) kitgrpc.ServerRequestFunc {
	return func(ctx context.Context, mds metadata.MD) context.Context {
		rpcMethod, ok := ctx.Value(kitgrpc.ContextKeyRequestMethod).(string)
		if !ok {
			rpcMethod = "/"
		}
		span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, rpcMethod, stdopentracing.Tag{
			Key:   string(ext.Component),
			Value: "ServerBefore",
		}, stdopentracing.Tag{
			Key:   string(ext.HTTPMethod),
			Value: ctx.Value(kitgrpc.ContextKeyRequestMethod),
		}, stdopentracing.Tag{
			Key:   "user_agent",
			Value: mds.Get("user-agent"),
		})

		traceId := span.Context().(jaeger.SpanContext).TraceID().String()
		ctx = context.WithValue(ctx, "traceId", traceId)
		span = span.SetTag("TraceId", span.Context().(jaeger.SpanContext).TraceID().String())

		defer func() {
			b, _ := json.Marshal(mds)
			span.LogKV("metadata", string(b))
			span.Finish()
		}()
		return ctx
	}
}

// TracingMiddleware request and response Tracing中间件
func TracingMiddleware(tracer stdopentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer == nil {
				return next(ctx, request)
			}
			span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, "TracingMiddleware", stdopentracing.Tag{
				Key:   string(ext.Component),
				Value: "Middleware",
			})
			defer func() {
				b, _ := json.Marshal(request)
				res, _ := json.Marshal(response)
				span.LogKV("request", string(b), "response", string(res), "err", err)
				span.Finish()
			}()
			return next(ctx, request)
		}
	}
}

func RecordRequestAndBody(tracer stdopentracing.Tracer, logger log.Logger, operationName string) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		// 获取并记录HTTP请求的各种入参
		method := req.Method
		u := req.URL.String()
		host := req.URL.Host
		headers := req.Header

		// 读取请求体
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			_ = level.Error(logger).Log("read body io.ReadAll", err.Error())
			return ctx
		}

		// 因为读取了整个请求体，我们需要将其写回以供后续处理
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)
		_ = level.Info(logger).Log("traceId", ctx.Value("traceId"), "method", method, "url", u, "host", host, "headers", headers, "body", string(bodyBytes))
		span.LogKV("method", method, "url", u, "host", host, "headers", headers, "body", string(bodyBytes))
		span.Finish()
		return ctx
	}
}
