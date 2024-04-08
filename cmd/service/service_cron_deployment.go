package service

import (
	"context"
	"fmt"
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
		_ = level.Warn(s.logger).Log("msg", "find deploy Pending models failed", "err", err.Error())
		return
	}
	// 2. 调用fschat modelse status
	fschatModels, err := s.apiSvc.FastChat().Models(s.ctx)
	if err != nil {
		_ = level.Warn(s.logger).Log("msg", "get fschat models failed", "err", err.Error())
		return
	}
	var modelNames []string
	for _, m := range fschatModels {
		modelNames = append(modelNames, m.ID)
	}
	// 3. 如果存在就是成功 否则就是pedding
	for _, m := range models {
		serviceName := util.ReplacerServiceName(m.ModelName)
		deploymentStatus, err := s.apiSvc.Runtime().GetDeploymentStatus(s.ctx, fmt.Sprintf("%s-%d", serviceName, m.ID))
		if err != nil {
			_ = level.Warn(logger).Log("msg", "get deployment status failed", "err", err.Error())
		}
		deployStatus := types.ModelDeployStatusSuccess

		if deploymentStatus == "Failed" {
			deployStatus = types.ModelDeployStatusFailed
		} else {
			if !util.StringInArray(modelNames, m.ModelName) {
				continue
			}
		}

		// 如果20分钟还没部署完，默认它部署失败
		//if time.Now().Unix()-model.ModelDeploy.CreatedAt.Unix() > 60*20 {
		//	deployStatus = types.ModelDeployStatusFailed
		//}
		// 更新状态为成功
		err = s.store.Model().UpdateDeployStatus(s.ctx, m.ID, deployStatus)
		if err != nil {
			_ = level.Warn(s.logger).Log("msg", "update deploy status failed", "err", err.Error())
			continue
		}
		//if deployStatus != types.ModelDeployStatusSuccess {
		//	continue
		//}
		// 更新模型状态为可用
		err = s.store.Model().SetModelEnabled(s.ctx, m.ModelName, true)
		if err != nil {
			_ = level.Warn(s.logger).Log("msg", "set model enabled failed", "err", err.Error())
			continue
		}
	}
	_ = level.Debug(s.logger).Log("msg", "deployment status job done")
}
