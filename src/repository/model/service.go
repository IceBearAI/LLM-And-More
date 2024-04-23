package model

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
	"time"
)

type Middleware func(Service) Service

type ListModelRequest struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	ModelName     string `json:"modelName"`
	Enabled       *bool  `json:"enabled"`
	IsPrivate     *bool  `json:"isPrivate"`
	IsFineTuning  *bool  `json:"isFineTuning"`
	ProviderName  string `json:"providerName"`
	ModelType     string `json:"modelType"`
	BaseModelName string `json:"baseModelName"` // null, notNull, 为空时，查全部
}

type ListEvalRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	ModelName   string `json:"modelName"`
	MetricName  string `json:"metricName"`
	Status      string `json:"status"`
	DatasetType string `json:"datasetType"`
}

type Service interface {
	// ListModels 模型分页列表
	ListModels(ctx context.Context, request ListModelRequest, preloads ...string) (res []types.Models, total int64, err error)
	// CreateModel 创建模型
	CreateModel(ctx context.Context, data *types.Models) (err error)
	// GetModel 获取模型
	GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error)
	// UpdateModel 更新模型
	UpdateModel(ctx context.Context, request UpdateModelRequest) (err error)
	// DeleteModel 删除模型
	DeleteModel(ctx context.Context, id uint) (err error)
	// FindModelsByTenantId 根据租户id查询模型
	FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error)
	// GetModelByModelName 根据模型名称查询模型
	GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error)
	// CreateEval 创建评估任务
	CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error)
	// ListEval 评估任务分页列表
	ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error)
	// UpdateEval 更新评估任务
	UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error)
	// GetEval 获取评估任务
	GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error)
	// DeleteEval 删除评估任务
	DeleteEval(ctx context.Context, id uint) (err error)
	// SaveModelDeploy 模型部署，信息入库
	SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error)
	// FindModelDeployByModeId 获取模型ID获取最近的部署
	FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error)
	// SaveModel 保存模型
	SaveModel(ctx context.Context, model *types.Models) (err error)
	// CancelModelDeploy 取消部署
	CancelModelDeploy(ctx context.Context, modelId uint) (err error)
	// FindDeployPendingModels 获取正在部署的模型
	FindDeployPendingModels(ctx context.Context) (models []types.Models, err error)
	// UpdateDeployStatus 更新部署状态
	UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error)
	// SetModelEnabled 设置模型是事可用
	SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error)
	// FindByModelId 根据id查询模型
	FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByModelId(ctx context.Context, modelId string, preloads ...string) (model types.Models, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		if preload == "ModelDeploy" {
			db = db.Preload(preload, "status = ? AND deleted_at IS NULL", types.ModelDeployStatusRunning.String())
		} else {
			db = db.Preload(preload)
		}
	}
	err = db.Where("model_name = ?", modelId).First(&model).Error
	return
}

func (s *service) SetModelEnabled(ctx context.Context, modelId string, enabled bool) (err error) {
	return s.db.WithContext(ctx).Model(&types.Models{}).Where("model_name = ?", modelId).Update("enabled", enabled).Error
}

func (s *service) FindDeployPendingModels(ctx context.Context) (models []types.Models, err error) {
	err = s.db.WithContext(ctx).
		InnerJoins("JOIN model_deploy ON model_deploy.status = ? AND model_deploy.deleted_at IS NULL", types.ModelDeployStatusPending).
		Find(&models).Error
	return
}

func (s *service) CancelModelDeploy(ctx context.Context, modelId uint) (err error) {
	model, err := s.FindModelDeployByModeId(ctx, modelId)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&types.ModelDeploy{}).Where("id = ?", model.ID).Delete(&types.ModelDeploy{}).Error; err != nil {
			return err
		}
		// 取消channel授权
		//if err = tx.Model(&types.Models{}).Where("model_id = ?", modelId).Association("Channels").Clear(); err != nil {
		//	return err
		//}
		return tx.Model(&types.Models{}).Where("id = ?", modelId).Update("enabled", false).Error
	})
}

func (s *service) SaveModel(ctx context.Context, model *types.Models) (err error) {
	err = s.db.WithContext(ctx).Save(model).Error
	return
}

func (s *service) FindModelDeployByModeId(ctx context.Context, modelId uint) (res types.ModelDeploy, err error) {
	err = s.db.WithContext(ctx).Where("model_id = ?", modelId).Order("id DESC").Limit(1).First(&res).Error
	return
}

func (s *service) DeleteEval(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&types.LLMEvalResults{}).Error
	return
}

func (s *service) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.LLMEvalResults{})
	if request.ModelName != "" {
		query = query.Where("model_name = ?", request.ModelName)
	}
	if request.MetricName != "" {
		query = query.Where("metric_name = ?", request.MetricName)
	}
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	if request.DatasetType != "" {
		query = query.Where("dataset_type = ?", request.DatasetType)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	err = s.db.WithContext(ctx).Where("model_name = ?", modelName).First(&res).Error
	return
}

func (s *service) ListModels(ctx context.Context, request ListModelRequest, preloads ...string) (res []types.Models, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Models{})
	for _, v := range preloads {
		if v == "ModelDeploy" {
			query = query.Preload(v, "status = ? AND deleted_at IS NULL", types.ModelDeployStatusRunning.String())
		} else {
			query = query.Preload(v)
		}
	}
	if request.ModelName != "" {
		query = query.Where("model_name LIKE ?", "%"+request.ModelName+"%")
	}
	if request.Enabled != nil {
		query = query.Where("enabled = ?", *request.Enabled)
	}
	if request.IsPrivate != nil {
		query = query.Where("is_private = ?", *request.IsPrivate)
	}
	if request.IsFineTuning != nil {
		query = query.Where("is_fine_tuning = ?", *request.IsFineTuning)
	}
	if request.ProviderName != "" {
		query = query.Where("provider_name = ?", request.ProviderName)
	}

	if request.ModelType != "" {
		query = query.Where("model_type = ?", request.ModelType)
	}

	if request.BaseModelName == "null" {
		query = query.Where("base_model_name is NULL or base_model_name = ''")
	} else if request.BaseModelName == "notNull" {
		query = query.Where("base_model_name is NOT NULL and base_model_name != '' ")
	}

	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("updated_at DESC").Limit(limit).Offset(offset).Preload("Tenants").Preload("ModelDeploy").Find(&res).Error
	return
}

func (s *service) CreateModel(ctx context.Context, data *types.Models) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&types.Models{}).Create(data).Error; err != nil {
			return err
		}
		if len(data.TenantId) > 0 {
			models := make([]types.TenantModelAssociations, 0)
			for _, v := range data.TenantId {
				models = append(models, types.TenantModelAssociations{
					ModelID:  data.ID,
					TenantID: v,
				})
			}
			if err = tx.Model(&types.TenantModelAssociations{}).Create(&models).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *service) GetModel(ctx context.Context, id uint, preloads ...string) (res types.Models, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		if preload == "ModelDeploy" {
			db = db.Preload(preload, "status = ? AND deleted_at IS NULL", types.ModelDeployStatusRunning.String())
		} else {
			db = db.Preload(preload)
		}
	}
	err = db.Where("id = ?", id).First(&res).Error
	return
}

type UpdateModelRequest struct {
	Id            uint
	TenantId      *[]uint
	MaxTokens     *int
	Enabled       *bool
	Remark        *string
	BaseModelName string
	Replicas      int
	Label         string
	K8sCluster    string
	InferredType  string
	Gpu           int
	Cpu           int
	Memory        int
}

func (s *service) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	data, err := s.GetModel(ctx, request.Id)
	if err != nil {
		return
	}
	if request.MaxTokens != nil {
		data.MaxTokens = *request.MaxTokens
	}
	if request.Enabled != nil {
		data.Enabled = *request.Enabled
	}
	if request.Remark != nil {
		data.Remark = *request.Remark
	}
	data.BaseModelName = request.BaseModelName
	data.Replicas = request.Replicas         //并行/实例数量
	data.Label = request.Label               //调度标签
	data.K8sCluster = request.K8sCluster     //k8s集群
	data.InferredType = request.InferredType //推理类型cpu,gpu
	data.Gpu = request.Gpu                   //GPU数
	data.Cpu = request.Cpu                   //CPU核数
	data.Memory = request.Memory             //内存G
	data.UpdatedAt = time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(data).Error; err != nil {
			return err
		}
		if request.TenantId != nil {
			// 删除原有的关联关系
			if err := tx.WithContext(ctx).Where("model_id = ?", request.Id).Delete(&types.TenantModelAssociations{}).Error; err != nil {
				return err
			}
			// 创建新的关联关系
			if len(*request.TenantId) > 0 {
				models := make([]types.TenantModelAssociations, 0)
				for _, v := range *request.TenantId {
					models = append(models, types.TenantModelAssociations{
						ModelID:  request.Id,
						TenantID: v,
					})
				}
				if err := tx.WithContext(ctx).Model(&types.TenantModelAssociations{}).Create(&models).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

func (s *service) DeleteModel(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.WithContext(ctx).Where("id = ?", id).Delete(&types.Models{}).Error; err != nil {
			return err
		}
		if err = tx.WithContext(ctx).Where("model_id = ?", id).Delete(&types.TenantModelAssociations{}).Error; err != nil {
			return err
		}
		if err = tx.WithContext(ctx).Where("model_id = ?", id).Delete(&types.ChannelModelAssociations{}).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *service) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	var tenant types.Tenants
	err = s.db.WithContext(ctx).Where("id = ?", tenantId).Preload("Models").First(&tenant).Error
	res = tenant.Models
	return
}

func (s *service) SaveModelDeploy(ctx context.Context, data *types.ModelDeploy) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) UpdateDeployStatus(ctx context.Context, modelId uint, status types.ModelDeployStatus) (err error) {
	return s.db.WithContext(ctx).Model(&types.ModelDeploy{}).Where("model_id = ?", modelId).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
