package fastchat

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateChatCompletionStream", "req", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateChatCompletionStream(ctx, req)
}

func (s *logging) ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletion", "model", model, "messages", messages, "temperature", temperature, "topP", topP, "presencePenalty", presencePenalty, "frequencyPenalty", frequencyPenalty, "maxToken", maxToken, "n", n, "stop", stop, "user", user, "functionCall", functionCall,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletion(ctx, model, messages, temperature, topP, presencePenalty, frequencyPenalty, maxToken, n, stop, user, functions, functionCall)
}

func (s *logging) ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature float64, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletionStream", "model", model, "temperature", temperature, "topP", topP, "presencePenalty", presencePenalty, "frequencyPenalty", frequencyPenalty, "maxToken", maxToken, "n", n, "stop", strings.Join(stop, ","), "user", user,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletionStream(ctx, model, messages, temperature, topP, presencePenalty, frequencyPenalty, maxToken, n, stop, user, functions, functionCall)
}

func (s *logging) CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateFineTuningJob", "TrainingFile", req.TrainingFile, "ValidationFile",
			req.ValidationFile, "Model", req.Model, "Suffix", req.Suffix, "Epochs", req.Hyperparameters,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateFineTuningJob(ctx, req)
}

func (s *logging) CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CancelFineTuningJob", "modelName", modelName, "jobId", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CancelFineTuningJob(ctx, modelName, jobId)
}

func (s *logging) ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListFineTune", "modelName", modelName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListFineTune(ctx, modelName)
}

func (s *logging) RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "RetrieveFineTuningJob", "jobId", jobId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.RetrieveFineTuningJob(ctx, jobId)
}

func (s *logging) UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UploadFile", "modelName", modelName, "fileName", fileName, "filePath", filePath, "purpose", purpose,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UploadFile(ctx, modelName, fileName, filePath, purpose)
}

func (s *logging) ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ModeRations", "model", model, "input", input,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ModeRations(ctx, model, input)
}

func (s *logging) Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Embeddings", "model", model, "documents", documents,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Embeddings(ctx, model, documents)
}

func (s *logging) CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CheckLength", "prompt", prompt, "maxToken", maxToken,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CheckLength(ctx, prompt, maxToken)
}

func (s *logging) CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateImage", "prompt", prompt, "size", size, "format", format,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateImage(ctx, prompt, size, format)
}

func (s *logging) Models(ctx context.Context) (res []openai.Model, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Models",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Models(ctx)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.fastchat", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
