package chat

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sashabaranov/go-openai"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Completion", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Completion(ctx, req)
}

func (s *tracing) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletion", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletion(ctx, req)
}

func (s *tracing) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletionStream(ctx, req)
}

func (s *tracing) Models(ctx context.Context) (res []openai.Model, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Models", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Models(ctx)
}

func (s *tracing) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Embeddings", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Embeddings(ctx, req)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}

type fschatWorkerTracing struct {
	next   WorkerService
	tracer opentracing.Tracer
}

func (s *fschatWorkerTracing) ListModels(ctx context.Context) (res []ModelCard, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ListModels", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ListModels(ctx)
}

func (s *fschatWorkerTracing) GetWorkerAddress(ctx context.Context, model string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetWorkerAddress", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("model", model, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetWorkerAddress(ctx, model)
}

func (s *fschatWorkerTracing) WorkerGetConvTemplate(ctx context.Context, workerAddress, model string) (res ModelConvTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGetConvTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "model", model, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGetConvTemplate(ctx, workerAddress, model)
}

func (s *fschatWorkerTracing) WorkerGenerateStream(ctx context.Context, workerAddress string, params GenerateStreamParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGenerateStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "params", params, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGenerateStream(ctx, workerAddress, params)
}

func (s *fschatWorkerTracing) WorkerGenerate(ctx context.Context, workerAddress string, params GenerateParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGenerate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "params", params, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGenerate(ctx, workerAddress, params)
}

func (s *fschatWorkerTracing) WorkerGetEmbeddings(ctx context.Context, workerAddress string, payload EmbeddingPayload) (res EmbeddingsResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGetEmbeddings", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "payload", payload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGetEmbeddings(ctx, workerAddress, payload)
}

func (s *fschatWorkerTracing) WorkerCountToken(ctx context.Context, workerAddress, model string, prompt any) (res int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerCountToken", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "model", model, "prompt", prompt, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerCountToken(ctx, workerAddress, model, prompt)
}

func (s *fschatWorkerTracing) WorkerGetStatus(ctx context.Context, workerAddress string) (res WorkerStatus, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGetStatus", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGetStatus(ctx, workerAddress)
}

func (s *fschatWorkerTracing) WorkerGetModelDetails(ctx context.Context, workerAddress, model string) (res ModelDetail, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerGetModelDetails", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "model", model, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerGetModelDetails(ctx, workerAddress, model)
}

func (s *fschatWorkerTracing) WorkerCheckLength(ctx context.Context, workerAddress string, model string, maxTokens int, prompt any) (res int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WorkerCheckLength", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "services.chat",
	})
	defer func() {
		span.LogKV("workerAddress", workerAddress, "model", model, "maxTokens", maxTokens, "prompt", prompt, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.WorkerCheckLength(ctx, workerAddress, model, maxTokens, prompt)
}

func NewFsChatWorkerTracing(otTracer opentracing.Tracer) FsChatWorkerMiddleware {
	return func(next WorkerService) WorkerService {
		return &fschatWorkerTracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
