package chat

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/sashabaranov/go-openai"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Completion",
			"req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Completion(ctx, req)
}

func (s *logging) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletion",
			"req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletion(ctx, req)
}

func (s *logging) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletionStream",
			"req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletionStream(ctx, req)
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

func (s *logging) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Embeddings",
			"req", fmt.Sprintf("%+v", req),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Embeddings(ctx, req)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "api.chat", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}

type fschatWorker struct {
	logger  log.Logger
	next    WorkerService
	traceId string
}

func (s *fschatWorker) ListModels(ctx context.Context) (res []ModelCard, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ListModels",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ListModels(ctx)
}

func (s *fschatWorker) GetWorkerAddress(ctx context.Context, model string) (res string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetWorkerAddress",
			"model", model,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetWorkerAddress(ctx, model)
}

func (s *fschatWorker) WorkerGetConvTemplate(ctx context.Context, workerAddress, model string) (res ModelConvTemplate, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGetConvTemplate",
			"workerAddress", workerAddress,
			"model", model,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGetConvTemplate(ctx, workerAddress, model)
}

func (s *fschatWorker) WorkerGenerateStream(ctx context.Context, workerAddress string, params GenerateStreamParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	defer func(begin time.Time) {
		//b, _ := json.Marshal(params)
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGenerateStream",
			"workerAddress", workerAddress,
			"params", fmt.Sprintf("%+v", params),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGenerateStream(ctx, workerAddress, params)
}

func (s *fschatWorker) WorkerGenerate(ctx context.Context, workerAddress string, params GenerateParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGenerate",
			"workerAddress", workerAddress,
			"params", fmt.Sprintf("%+v", params),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGenerate(ctx, workerAddress, params)
}

func (s *fschatWorker) WorkerGetEmbeddings(ctx context.Context, workerAddress string, payload EmbeddingPayload) (res EmbeddingsResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGetEmbeddings",
			"workerAddress", workerAddress,
			"payload", fmt.Sprintf("%+v", payload),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGetEmbeddings(ctx, workerAddress, payload)
}

func (s *fschatWorker) WorkerCountToken(ctx context.Context, workerAddress, model string, prompt any) (res int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerCountToken",
			"workerAddress", workerAddress,
			"model", model,
			"prompt", fmt.Sprintf("%+v", prompt),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerCountToken(ctx, workerAddress, model, prompt)
}

func (s *fschatWorker) WorkerGetStatus(ctx context.Context, workerAddress string) (res WorkerStatus, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGetStatus",
			"workerAddress", workerAddress,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGetStatus(ctx, workerAddress)
}

func (s *fschatWorker) WorkerGetModelDetails(ctx context.Context, workerAddress, model string) (res ModelDetail, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerGetModelDetails",
			"workerAddress", workerAddress,
			"model", model,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerGetModelDetails(ctx, workerAddress, model)
}

func (s *fschatWorker) WorkerCheckLength(ctx context.Context, workerAddress string, model string, maxTokens int, prompt any) (res int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "WorkerCheckLength",
			"workerAddress", workerAddress,
			"model", model,
			"maxTokens", maxTokens,
			"prompt", fmt.Sprintf("%+v", prompt),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.WorkerCheckLength(ctx, workerAddress, model, maxTokens, prompt)
}

func NewFsChatWorkerLogging(logger log.Logger, traceId string) FsChatWorkerMiddleware {
	logger = log.With(logger, "api.chat", "logging")
	return func(next WorkerService) WorkerService {
		return &fschatWorker{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
