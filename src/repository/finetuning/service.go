package finetuning

import (
	"bytes"
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
	"strings"
	"text/template"
)

type Middleware func(Service) Service

type ListFineTuningTemplateRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type Service interface {
	// CreateFineTuningJob 创建训练任务
	CreateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error)
	// FindFineTuningTemplateByModel 根据模型名称查找微调模版
	FindFineTuningTemplateByModel(ctx context.Context, modelName string, preloads ...string) (template types.FineTuningTemplate, err error)
	// FindFineTuningTemplateByModelType 根据模型名称和模版类型查找微调模版
	FindFineTuningTemplateByModelType(ctx context.Context, modelName string, templateType types.TemplateType, preloads ...string) (template types.FineTuningTemplate, err error)
	// FindFineTuningJobByJobId 根据jobId查找微调任务
	FindFineTuningJobByJobId(ctx context.Context, jobId string, preloads ...string) (job types.FineTuningTrainJob, err error)
	// FindFineTuningJob 根据id查找微调任务
	FindFineTuningJob(ctx context.Context, id uint, preloads ...string) (job types.FineTuningTrainJob, err error)
	// EncodeFineTuningJobTemplate 生成微调模版
	EncodeFineTuningJobTemplate(ctx context.Context, tpl string, job *types.FineTuningTrainJob) (re string, err error)
	// FindFineTuningJobLastByStatus 查找最后一个微调任务
	FindFineTuningJobLastByStatus(ctx context.Context, status types.TrainStatus, orderBy string, preloads ...string) (jobs types.FineTuningTrainJob, err error)
	// HasRunningJob 是否有正在Running的任务
	HasRunningJob(ctx context.Context) (has bool, err error)
	// UpdateFineTuningJob 更新微调任务
	UpdateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error)
	// ListFindTuningJob 分页列出微调任务
	ListFindTuningJob(ctx context.Context, request ListFindTuningJobRequest) (jobs []types.FineTuningTrainJob, total int64, err error)
	// CountFineTuningJobByStatus 根据状态统计任务数量
	CountFineTuningJobByStatus(ctx context.Context) (res map[string]int64, err error)
	// CountFineTuningJobDuration 统计训练任务时长
	CountFineTuningJobDuration(ctx context.Context) (res int64, err error)
	// DeleteFineTuningJob 删除微调任务
	DeleteFineTuningJob(ctx context.Context, id uint) (err error)
	// ListFineTuningTemplate 分页列出微调模版
	ListFineTuningTemplate(ctx context.Context, request ListFineTuningTemplateRequest) (res []types.FineTuningTemplate, total int64, err error)
	// GetFineTuningJobByModelName 根据模型名称查找微调任务
	GetFineTuningJobByModelName(ctx context.Context, modelName string, preloads ...string) (job types.FineTuningTrainJob, err error)
	// FindFineTuningTemplateByType 根据类型查找微调模版
	FindFineTuningTemplateByType(ctx context.Context, modelName string, templateType types.TemplateType) (tpl types.FineTuningTemplate, err error)
	// FindFineTuningJobRunning 查找正在运行的任务
	FindFineTuningJobRunning(ctx context.Context, preloads ...string) (jobs []types.FineTuningTrainJob, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindFineTuningJobRunning(ctx context.Context, preloads ...string) (jobs []types.FineTuningTrainJob, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("train_status = ?", types.TrainStatusRunning).Find(&jobs).Error
	return
}

func (s *service) FindFineTuningTemplateByType(ctx context.Context, modelName string, templateType types.TemplateType) (tpl types.FineTuningTemplate, err error) {
	err = s.db.WithContext(ctx).Where("base_model = ? AND template_type = ?", modelName, templateType).First(&tpl).Error
	return
}

func (s *service) FindFineTuningTemplateByModelType(ctx context.Context, modelName string, templateType types.TemplateType, preloads ...string) (template types.FineTuningTemplate, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("base_model = ? and template_type = ?", modelName, templateType).First(&template).Error
	return
}

func (s *service) GetFineTuningJobByModelName(ctx context.Context, modelName string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("fine_tuned_model = ?", modelName).First(&job).Error
	return
}

func (s *service) ListFineTuningTemplate(ctx context.Context, request ListFineTuningTemplateRequest) (res []types.FineTuningTemplate, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.FineTuningTemplate{})
	limit, offset := page.Limit(request.Page, request.PageSize)
	query = query.Where("template_type = ?", "train")
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) DeleteFineTuningJob(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Delete(&types.FineTuningTrainJob{}, id).Error
	return
}

func (s *service) CountFineTuningJobDuration(ctx context.Context) (res int64, err error) {
	err = s.db.WithContext(ctx).Model(&types.FineTuningTrainJob{}).Where("train_status in ?", []types.TrainStatus{
		types.TrainStatusSuccess, types.TrainStatusFailed,
	}).Select("sum(train_duration) as total").Scan(&res).Error
	return res, nil
}

func (s *service) CountFineTuningJobByStatus(ctx context.Context) (res map[string]int64, err error) {
	res = make(map[string]int64)
	type statusCount struct {
		TrainStatus string `gorm:"column:train_status"`
		Total       int64  `gorm:"column:total"`
	}
	var statusCounts []statusCount
	err = s.db.WithContext(ctx).Model(&types.FineTuningTrainJob{}).Select("train_status, count(*) as total").Group("train_status").Scan(&statusCounts).Error
	if err != nil {
		return
	}
	for _, v := range statusCounts {
		res[v.TrainStatus] = v.Total
	}
	return
}

type ListFindTuningJobRequest struct {
	Page           int    `json:"page"`
	PageSize       int    `json:"pageSize"`
	FineTunedModel string `json:"fineTunedModel"` // 微调模型
	TrainStatus    string `json:"trainStatus"`    // 训练状态
}

func (s *service) ListFindTuningJob(ctx context.Context, request ListFindTuningJobRequest) (res []types.FineTuningTrainJob, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.FineTuningTrainJob{})
	if request.FineTunedModel != "" {
		query = query.Where("fine_tuned_model = ?", request.FineTunedModel)
	}
	if request.TrainStatus != "" {
		query = query.Where("train_status = ?", request.TrainStatus)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Omit("train_log").Find(&res).Error
	return
}

func (s *service) UpdateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	job.Template = types.FineTuningTemplate{}
	job.FineTuningFile = types.Files{}
	err = s.db.WithContext(ctx).Model(&job).Save(job).Error
	return
}

func (s *service) HasRunningJob(ctx context.Context) (has bool, err error) {
	var count int64
	err = s.db.WithContext(ctx).Model(&types.FineTuningTrainJob{}).Where("train_status = ?", types.TrainStatusRunning).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		has = true
	}
	return
}

func (s *service) FindFineTuningJobLastByStatus(ctx context.Context, status types.TrainStatus, orderBy string, preloads ...string) (jobs types.FineTuningTrainJob, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if strings.EqualFold(orderBy, "") {
		orderBy = "id asc"
	}
	err = db.Where("train_status = ?", status).Order(orderBy).First(&jobs).Error
	return
}

func (s *service) FindFineTuningJobByJobId(ctx context.Context, jobId string, preloads ...string) (job types.FineTuningTrainJob, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("job_id = ?", jobId).First(&job).Error
	return
}

func (s *service) FindFineTuningJob(ctx context.Context, id uint, preloads ...string) (job types.FineTuningTrainJob, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("id = ?", id).First(&job).Error
	return
}

func (s *service) EncodeFineTuningJobTemplate(ctx context.Context, tpl string, job *types.FineTuningTrainJob) (re string, err error) {
	tmpl, err := template.New(job.JobId).Parse(tpl)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, job)
	if err != nil {
		return "", err
	}
	re = buffer.String()
	return
}

func (s *service) CreateFineTuningJob(ctx context.Context, job *types.FineTuningTrainJob) (err error) {
	err = s.db.WithContext(ctx).Create(&job).Error
	return
}

func (s *service) FindFineTuningTemplateByModel(ctx context.Context, modelName string, preloads ...string) (template types.FineTuningTemplate, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("base_model = ? and template_type = ?", modelName, "train").First(&template).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
