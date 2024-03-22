package datasets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

type Middleware func(Service) Service

// Service 数据集服务逻辑模块
type Service interface {
	// List 根据租户获取数据集列表
	List(ctx context.Context, tenantId uint, page, pageSize int, query string) (datasets []datasetResult, total int64, err error)
	// Create 创建数据集
	Create(ctx context.Context, tenantId uint, name, remark string) (datasetId string, err error)
	// Update 更新数据集
	Update(ctx context.Context, tenantId uint, datasetId, name, remark string) (err error)
	// Delete 删除数据集
	Delete(ctx context.Context, tenantId uint, datasetId string) (err error)
	// Detail 获取数据集详情
	Detail(ctx context.Context, tenantId uint, datasetId string) (dataset datasetResult, err error)
	// AddSample 给数据集添加样本
	// sampleType: 样本类型，默认为text
	AddSample(ctx context.Context, tenantId uint, uuid, sampleType string, samples []message) (err error)
	// DeleteSample 删除数据集样本
	// sampleUIds: 样本唯一标识列表
	DeleteSample(ctx context.Context, tenantId uint, datasetId string, sampleUIds []string) (err error)
	// SampleList 根据数据集ID获取样本列表
	SampleList(ctx context.Context, tenantId uint, datasetId string, page, pageSize int, title string) (samples []datasetSampleResult, total int64, err error)
	// UpdateSampleMessages 更新样本对话
	UpdateSampleMessages(ctx context.Context, tenantId uint, datasetId string, sampleUId string, messages []message) (err error)
	// ExportSample 导出样本
	// format: jsonl 或 train
	ExportSample(ctx context.Context, tenantId uint, datasetId, format string) (samples []addSampleRequest, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) ExportSample(ctx context.Context, tenantId uint, datasetId, format string) (samples []addSampleRequest, err error) {
	logger := log.With(s.logger, "method", "List")
	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId, "Samples")
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return
	}
	for _, item := range dataset.Samples {
		var messages []message
		_ = json.Unmarshal([]byte(item.Conversations), &messages)
		samples = append(samples, addSampleRequest{
			Messages: messages,
		})
	}
	return
}

func (s *service) List(ctx context.Context, tenantId uint, page, pageSize int, query string) (datasets []datasetResult, total int64, err error) {
	logger := log.With(s.logger, "method", "List")
	list, total, err := s.repository.Dataset().List(ctx, tenantId, page, pageSize, query)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "List", "err", err)
		return
	}
	for _, item := range list {
		datasets = append(datasets, datasetResult{
			UUID:      item.UUID,
			Name:      item.Name,
			Remark:    item.Remark,
			Type:      string(item.Type),
			Creator:   item.CreatorEmail,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Samples:   item.SampleCount,
		})
	}

	return
}

func (s *service) Create(ctx context.Context, tenantId uint, name, remark string) (datasetId string, err error) {
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	logger := log.With(s.logger, "method", "Create")
	dataset := &types.Dataset{
		Name:         name,
		Remark:       remark,
		UUID:         uuid.New().String(),
		CreatorEmail: email,
		TenantId:     tenantId,
		Type:         types.DatasetTypeText,
	}
	err = s.repository.Dataset().Create(ctx, dataset)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "Create", "err", err)
		return
	}
	datasetId = dataset.UUID
	return datasetId, err
}

func (s *service) Update(ctx context.Context, tenantId uint, uuid, name, remark string) (err error) {
	logger := log.With(s.logger, "method", "Update")
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)

	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, uuid, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return err
	}
	dataset.Name = name
	dataset.Remark = remark
	dataset.CreatorEmail = email
	err = s.repository.Dataset().Update(ctx, &dataset)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "Update", "err", err)
		return
	}
	return err
}

func (s *service) Delete(ctx context.Context, tenantId uint, datasetId string) (err error) {
	logger := log.With(s.logger, "method", "Delete")
	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return err
	}
	err = s.repository.Dataset().Delete(ctx, dataset.ID)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "Delete", "err", err)
		return
	}
	return err
}

func (s *service) Detail(ctx context.Context, tenantId uint, datasetId string) (dataset datasetResult, err error) {
	logger := log.With(s.logger, "method", "Detail")
	datasetInfo, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return dataset, err
	}
	dataset = datasetResult{
		UUID:      datasetInfo.UUID,
		Name:      datasetInfo.Name,
		Remark:    datasetInfo.Remark,
		Type:      string(datasetInfo.Type),
		Creator:   datasetInfo.CreatorEmail,
		CreatedAt: datasetInfo.CreatedAt,
		UpdatedAt: datasetInfo.UpdatedAt,
		Samples:   datasetInfo.SampleCount,
	}
	return dataset, err
}

func (s *service) AddSample(ctx context.Context, tenantId uint, datasetId, sampleType string, samples []message) (err error) {
	logger := log.With(s.logger, "method", "AddSample")
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return err
	}
	var title string
	if len(samples) > 0 {
		title = samples[0].Content
	}
	// 根据title查询是否是有一样的数据
	sample, err := s.repository.Dataset().FindSampleByTitle(ctx, dataset.ID, title)
	if err == nil && sample.ID != 0 {
		err = errors.New(fmt.Sprintf("%s", "样本已存在"))
		_ = level.Warn(logger).Log("repository.Dataset", "FindSampleByTitle", "err", err)
		return
	}
	content, _ := json.Marshal(samples)
	err = s.repository.Dataset().CreateSample(ctx, &types.DatasetSample{
		Title:         title,
		DatasetId:     dataset.ID,
		Conversations: string(content),
		Turns:         len(samples) / 2,
		UUID:          uuid.New().String(),
		CreatorEmail:  email,
	})
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "CreateSample", "err", err)
		return err
	}
	// 更新数据集样本数量
	dataset.SampleCount = dataset.SampleCount + 1
	dataset.CreatorEmail = email
	if err = s.repository.Dataset().Update(ctx, &dataset); err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "Update", "err", err)
	}
	return nil
}

func (s *service) DeleteSample(ctx context.Context, tenantId uint, datasetId string, sampleUIds []string) (err error) {
	logger := log.With(s.logger, "method", "DeleteSample")
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return err
	}
	_ = level.Info(logger).Log("repository.Dataset", "DeleteSample", "dataset", dataset.ID, "sampleUIds", strings.Join(sampleUIds, ","))
	err = s.repository.Dataset().DeleteSampleByUUID(ctx, sampleUIds)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "DeleteSample", "err", err)
		return err
	}
	// 删除数据集样本数
	dataset.SampleCount = dataset.SampleCount - len(sampleUIds)
	dataset.CreatorEmail = email
	if err = s.repository.Dataset().Update(ctx, &dataset); err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "Update", "err", err)
	}
	return nil
}

func (s *service) SampleList(ctx context.Context, tenantId uint, datasetId string, page, pageSize int, title string) (samples []datasetSampleResult, total int64, err error) {
	logger := log.With(s.logger, "method", "SampleList")
	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return samples, total, err
	}
	list, total, err := s.repository.Dataset().SampleList(ctx, dataset.ID, page, pageSize, title)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "SampleList", "err", err)
		return samples, total, err
	}
	for _, item := range list {
		var messages []message
		_ = json.Unmarshal([]byte(item.Conversations), &messages)
		samples = append(samples, datasetSampleResult{
			UUID:          item.UUID,
			Title:         item.Title,
			Conversations: item.Conversations,
			Label:         item.Label,
			Remark:        item.Remark,
			Turns:         item.Turns,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
			CreatorEmail:  item.CreatorEmail,
			Messages:      messages,
		})
	}
	return samples, total, err
}

func (s *service) UpdateSampleMessages(ctx context.Context, tenantId uint, datasetId string, sampleUId string, messages []message) (err error) {
	logger := log.With(s.logger, "method", "UpdateSample")
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)

	dataset, err := s.repository.Dataset().FindByUUIDAndTenantId(ctx, datasetId, tenantId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "FindByUUIDAndTenantId", "err", err)
		return err
	}
	_ = level.Info(logger).Log("repository.Dataset", "UpdateSample", "dataset", dataset.ID, "sampleUId", sampleUId)
	datasetSample, err := s.repository.Dataset().FindSampleByUUID(ctx, sampleUId)
	var title string
	if len(messages) > 0 {
		title = messages[0].Content
	}
	content, _ := json.Marshal(messages)
	datasetSample.Title = title
	datasetSample.Conversations = string(content)
	datasetSample.Turns = len(messages) / 2
	datasetSample.CreatorEmail = email
	err = s.repository.Dataset().UpdateSample(ctx, &datasetSample)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Dataset", "UpdateSample", "err", err)
		return err
	}
	return err
}

func New(logger log.Logger, traceId string, store repository.Repository) Service {
	logger = log.With(logger, "service", "datasets")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: store,
	}
}
