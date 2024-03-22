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
		_ = level.Error(s.logger).Log("msg", "fine tuning run waiting train failed", "err", err.Error(), "name", s.Name)
		return
	}
	_ = level.Info(s.logger).Log("msg", "fine tuning run waiting train success", "name", s.Name)
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
		_ = level.Error(logger).Log("msg", "find running job failed", "err", err.Error(), "name", s.Name)
		return
	}
	if funErr := s.runFun(s.ctx, runningJobs); funErr != nil {
		_ = level.Error(logger).Log("msg", "fine tuning running job log failed", "err", funErr.Error(), "name", s.Name)
		return
	}
	_ = level.Info(s.logger).Log("msg", "fine tuning running job log success", "name", s.Name)
}
