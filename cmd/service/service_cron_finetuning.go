package service

import (
	"context"
	"github.com/IceBearAI/aigc/src/pkg/finetuning"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type fineTuningRunWaitingTrainCronJob struct {
	fineTuning finetuning.Service
	logger     log.Logger
	Name       string
	ctx        context.Context
}

func (s *fineTuningRunWaitingTrainCronJob) Run() {
	if err := s.fineTuning.RunWaitingTrain(s.ctx); err != nil {
		_ = level.Warn(s.logger).Log("msg", "fine tuning run waiting train failed", "err", err.Error())
		return
	}
	_ = level.Debug(s.logger).Log("msg", "fine tuning run waiting train success")
}

type fineTuningRunningLogCronJob struct {
	fineTuning finetuning.Service
	logger     log.Logger
	Name       string
	ctx        context.Context
	store      repository.Repository
	runFun     func(ctx context.Context, runningJobs []types.FineTuningTrainJob) error
}

func (s *fineTuningRunningLogCronJob) Run() {
	runningJobs, err := s.store.FineTuning().FindFineTuningJobRunning(s.ctx)
	if err != nil {
		_ = level.Warn(s.logger).Log("msg", "find running job failed", "err", err.Error())
		return
	}
	if funErr := s.runFun(s.ctx, runningJobs); funErr != nil {
		_ = level.Warn(s.logger).Log("msg", "fine tuning running job log failed", "err", funErr.Error())
		return
	}
	_ = level.Info(s.logger).Log("msg", "fine tuning running job log success")
}
