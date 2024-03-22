package fastchat

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

func (s *tracing) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("req", req, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateChatCompletionStream(ctx, req)
}

func (s *tracing) ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletion", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("model", model, "messages", messages, "temperature", temperature, "topP", topP, "presencePenalty", presencePenalty, "frequencyPenalty", frequencyPenalty, "maxToken", maxToken, "n", n, "stop", stop, "user", user, "functions", functions, "functionCall", functionCall, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletion(ctx, model, messages, temperature, topP, presencePenalty, frequencyPenalty, maxToken, n, stop, user, functions, functionCall)
}

func (s *tracing) ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature float64, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ChatCompletionStream", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("model", model, "messages", messages, "temperature", temperature, "topP", topP,
			"presencePenalty", presencePenalty, "frequencyPenalty", frequencyPenalty, "maxToken", maxToken, "n", n,
			"stop", stop, "user", user, "functions", functions, "functionCall", functionCall, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ChatCompletionStream(ctx, model, messages, temperature, topP, presencePenalty, frequencyPenalty, maxToken, n, stop, user, functions, functionCall)
}

func (s *tracing) CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("TrainingFile", req.TrainingFile, "ValidationFile", req.ValidationFile, "Model", req.Model, "Suffix", req.Suffix, "Epochs", req.Hyperparameters, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateFineTuningJob(ctx, req)
}

func (s *tracing) CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CancelFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("jobId", jobId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CancelFineTuningJob(ctx, modelName, jobId)
}

func (s *tracing) ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *tracing) RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RetrieveFineTuningJob", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("jobId", jobId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RetrieveFineTuningJob(ctx, jobId)
}

func (s *tracing) UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UploadFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("modelName", modelName, "fileName", fileName, "filePath", filePath, "purpose", purpose, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UploadFile(ctx, modelName, fileName, filePath, purpose)
}

func (s *tracing) ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ModeRations", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("model", model, "input", input, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ModeRations(ctx, model, input)
}

func (s *tracing) Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Embeddings", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("model", model, "documents", documents, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Embeddings(ctx, model, documents)
}

func (s *tracing) CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CheckLength", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("prompt", prompt, "maxToken", maxToken, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CheckLength(ctx, prompt, maxToken)
}

func (s *tracing) CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateImage", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("prompt", prompt, "size", size, "format", format, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateImage(ctx, prompt, size, format)
}

func (s *tracing) Models(ctx context.Context) (res []openai.Model, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Models", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "api.fastchat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Models(ctx)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
