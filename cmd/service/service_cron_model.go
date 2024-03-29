package service

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/model"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"sync"
)

type modelEvalLogCronJob struct {
	logger log.Logger
	Name   string
	ctx    context.Context
	store  repository.Repository
	apiSvc services.Service
	mu     sync.Mutex
}

func (s *modelEvalLogCronJob) Run() {
	// 加个锁，防止没处理完，再次执行

	models, total, jobErr := s.store.Model().ListModels(s.ctx, model.ListModelRequest{
		Page:         1,
		PageSize:     100,
		ProviderName: types.ModelProviderLocalAI.String(),
		ModelType:    types.ModelTypeTextGeneration.String(),
	})
	if jobErr != nil {
		_ = level.Warn(s.logger).Log("store.Model", "ListModels", "err", jobErr.Error())
		return
	}
	_ = level.Debug(s.logger).Log("total", total)
	for _, m := range models {
		evaluateModels, evalTotal, evalErr := s.store.ModelEvaluate().ListModelEvaluate(s.ctx, 1, 10, m.ID, string(types.EvaluateStatusRunning), "")
		if evalErr != nil {
			_ = level.Error(s.logger).Log("modelName", m.ModelName, "store.ModelEvaluate", "ListModelEvaluate", "err", evalErr.Error())
			continue
		}
		_ = level.Debug(s.logger).Log("modelName", m.ModelName, "evalTotal", evalTotal)
		for _, em := range evaluateModels {
			jobLog, jobLogErr := s.apiSvc.Runtime().GetJobLogs(s.ctx, em.JobName)
			if jobLogErr != nil {
				_ = level.Warn(s.logger).Log("apiSvc.Runtime", "GetJobLogs", "err", jobLogErr.Error())
				continue
			}
			em.EvaluateLog = jobLog
			if saveErr := s.store.ModelEvaluate().Save(s.ctx, &em); saveErr != nil {
				_ = level.Error(s.logger).Log("store.ModelEvaluate", "Save", "err", saveErr.Error())
				continue
			}
		}
	}
	_ = level.Info(s.logger).Log("msg", "model evaluate log running job log success")
}
