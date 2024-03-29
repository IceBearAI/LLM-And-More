package service

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type datasetCheckTaskSimilarCronJob struct {
	logger log.Logger
	Name   string
	ctx    context.Context
	store  repository.Repository
	apiSvc services.Service
}

func (s *datasetCheckTaskSimilarCronJob) Run() {
	tasks, err := s.store.DatasetTask().GetTaskByDetection(s.ctx, types.DatasetAnnotationStatusCompleted, types.DatasetAnnotationDetectionStatusProcessing)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "get dataset task failed", "err", err.Error())
		return
	}
	for _, task := range tasks {
		jobLog, err := s.apiSvc.Runtime().GetJobLogs(s.ctx, task.JobName)
		if err != nil {
			_ = level.Warn(s.logger).Log("msg", "get job logs failed", "err", err.Error())
			continue
		}
		task.DetectionLog = jobLog
		if err := s.store.DatasetTask().UpdateTask(s.ctx, &task); err != nil {
			_ = level.Warn(s.logger).Log("msg", "update task failed", "err", err.Error())
		}
	}
	_ = level.Debug(s.logger).Log("msg", "check task similar job done")
}
