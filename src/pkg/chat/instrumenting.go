package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/services/chat"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sashabaranov/go-openai"
	"math"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (s *instrumentingService) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ChatCompletion").Add(1)
		s.requestLatency.With("method", "ChatCompletion").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.ChatCompletion(ctx, channelId, req)
}

func (s *instrumentingService) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ChatCompletionStream").Add(1)
		s.requestLatency.With("method", "ChatCompletionStream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.ChatCompletionStream(ctx, channelId, req)
}

func (s *instrumentingService) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Models").Add(1)
		s.requestLatency.With("method", "Models").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Models(ctx, channelId)
}

func (s *instrumentingService) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Embeddings").Add(1)
		s.requestLatency.With("method", "Embeddings").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Embeddings(ctx, channelId, req)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(s Service) Service {
		return &instrumentingService{
			requestCount:   counter,
			requestLatency: latency,
			next:           s,
		}
	}
}

func NewChatQueueGaugeService(logger log.Logger, workerSvc chat.WorkerService) prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "chat",
		Name:      "queue_size",
		Help:      "Current size of the chat queue.",
	}, func() float64 {
		ctx := context.Background()
		models, err := workerSvc.ListModels(ctx)
		if err != nil {
			_ = level.Warn(logger).Log("msg", "failed to get models", "err", err)
			return 0
		}
		var queueSize = 0
		for _, v := range models {
			_ = level.Debug(logger).Log("msg", "model", "model", v)
			workerAddress, err := workerSvc.GetWorkerAddress(ctx, v.ID)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "failed to get worker address", "model", v.ID, "err", err)
				continue
			}
			_ = level.Debug(logger).Log("msg", "worker address", "model", v.ID, "address", workerAddress)
			workerStatus, err := workerSvc.WorkerGetStatus(context.Background(), workerAddress)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "failed to get worker status", "model", v.ID, "err", err)
				continue
			}
			queueSize += workerStatus.QueueLength
		}
		return float64(queueSize)
	})
}

func NewChatAvgSpeedGaugeService(logger log.Logger, workerSvc chat.WorkerService) prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "chat",
		Name:      "avg_speed",
		Help:      "Current model avg speed of the chat.",
	}, func() float64 {
		ctx := context.Background()
		models, err := workerSvc.ListModels(ctx)
		if err != nil {
			_ = level.Warn(logger).Log("msg", "failed to get models", "err", err)
			return 0
		}
		var speed = 0
		totalModel := len(models)
		for _, v := range models {
			_ = level.Debug(logger).Log("msg", "model", "model", v)
			workerAddress, err := workerSvc.GetWorkerAddress(ctx, v.ID)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "failed to get worker address", "model", v.ID, "err", err)
				continue
			}
			_ = level.Debug(logger).Log("msg", "worker address", "model", v.ID, "address", workerAddress)
			workerStatus, err := workerSvc.WorkerGetStatus(context.Background(), workerAddress)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "failed to get worker status", "model", v.ID, "err", err)
				continue
			}
			speed += workerStatus.Speed
		}
		avgSpeed := float64(speed) / float64(totalModel)
		// 返回保留两位小数
		return math.Round(avgSpeed*100) / 100
	})
}
