package service

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type deploymentStatusCronJob struct {
	logger log.Logger
	Name   string
	ctx    context.Context
	store  repository.Repository
	apiSvc services.Service
}

func (s *deploymentStatusCronJob) Run() {
	// 1. 获取正在等待部署的模型
	models, err := s.store.Model().FindDeployPendingModels(s.ctx)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "find deploy Pending models failed", "err", err.Error())
		return
	}
	// 2. 调用fschat modelse status
	fschatModels, err := s.apiSvc.FastChat().Models(s.ctx)
	if err != nil {
		_ = level.Error(logger).Log("msg", "get fschat models failed", "err", err.Error())
		return
	}
	var modelNames []string
	for _, model := range fschatModels {
		modelNames = append(modelNames, model.ID)
	}
	// 3. 如果存在就是成功 否则就是pedding
	for _, model := range models {
		if !util.StringInArray(modelNames, model.ModelName) {
			continue
		}
		deployStatus := types.ModelDeployStatusSuccess
		// 如果20分钟还没部署完，默认它部署失败
		//if time.Now().Unix()-model.ModelDeploy.CreatedAt.Unix() > 60*20 {
		//	deployStatus = types.ModelDeployStatusFailed
		//}
		// 更新状态为成功
		err = s.store.Model().UpdateDeployStatus(s.ctx, model.ID, deployStatus)
		if err != nil {
			_ = level.Error(logger).Log("msg", "update deploy status failed", "err", err.Error())
			continue
		}
		//if deployStatus != types.ModelDeployStatusSuccess {
		//	continue
		//}
		// 更新模型状态为可用
		err = s.store.Model().SetModelEnabled(s.ctx, model.ModelName, true)
		if err != nil {
			_ = level.Error(logger).Log("msg", "set model enabled failed", "err", err.Error())
			continue
		}
	}
}
