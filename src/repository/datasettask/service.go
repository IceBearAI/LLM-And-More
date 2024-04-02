package datasettask

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Middleware func(Service) Service

// Service is the interface that provides datasettask methods.
type Service interface {
	// GetDatasetDocumentByUUID returns a dataset document by UUID.
	GetDatasetDocumentByUUID(ctx context.Context, tenantId uint, uuid string, preload ...string) (res *types.DatasetDocument, err error)
	// CreateTask creates a new task.
	CreateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error)
	// ListTasks returns all tasks.
	ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (res []types.DatasetAnnotationTask, total int64, err error)
	// DeleteTask deletes a task.
	DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error)
	// UpdateTask updates a task.
	UpdateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error)
	// GetTask returns a task.
	GetTask(ctx context.Context, tenantId uint, uuid string, preloads ...string) (res *types.DatasetAnnotationTask, err error)
	// AddTaskSegments adds task segments.
	AddTaskSegments(ctx context.Context, data []types.DatasetAnnotationTaskSegment) (err error)
	// GetTaskOneSegment 获取一条待标注任务样本
	GetTaskOneSegment(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error)
	// GetTaskSegmentByUUID 获取一条任务样本
	GetTaskSegmentByUUID(ctx context.Context, taskId uint, uuid string, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error)
	// UpdateTaskSegment 更新任务样本
	UpdateTaskSegment(ctx context.Context, data *types.DatasetAnnotationTaskSegment) (err error)
	// GetTaskSegments 获取任务样本
	GetTaskSegments(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, page, pageSize int, preload ...string) (res []types.DatasetAnnotationTaskSegment, total int64, err error)
	// GetTaskSegmentByRand 获取任务样本
	GetTaskSegmentByRand(ctx context.Context, taskId uint, testPercent float64, status types.DatasetAnnotationStatus, segmentType types.DatasetAnnotationSegmentType, preload ...string) (res []types.DatasetAnnotationTaskSegment, err error)
	// UpdateTaskSegmentType 更新任务样本类型
	UpdateTaskSegmentType(ctx context.Context, segmentIds []uint, segmentType types.DatasetAnnotationSegmentType) (err error)
	// GetDatasetDocumentSegmentByRange returns a dataset document segment by range.
	GetDatasetDocumentSegmentByRange(ctx context.Context, datasetDocumentId uint, start, end int, preload ...string) (res []types.DatasetDocumentSegment, err error)
	// DeleteDatasetDocument 删除数据集文档
	DeleteDatasetDocument(ctx context.Context, tenantId uint, uuid string) (err error)
	// DeleteDatasetDocumentById 删除数据集文档
	DeleteDatasetDocumentById(ctx context.Context, id uint, unscoped bool) (err error)
	// ListDatasetDocuments returns all dataset documents.
	ListDatasetDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []types.DatasetDocument, total int64, err error)
	// CreateDatasetDocument creates a new dataset document.
	CreateDatasetDocument(ctx context.Context, data *types.DatasetDocument) (err error)
	// AddDatasetDocumentSegments adds dataset document segments.
	AddDatasetDocumentSegments(ctx context.Context, data []types.DatasetDocumentSegment) (err error)
	// UpdateDatasetDocumentSegmentCount 更新数据集文档样本数量
	UpdateDatasetDocumentSegmentCount(ctx context.Context, datasetDocumentId uint, count int) (err error)
	// GetTaskSegmentPrev 获取上一条已标注的内容
	GetTaskSegmentPrev(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus) (res types.DatasetAnnotationTaskSegment, err error)
	// GetTaskByDetection 获取正在评估检测的任务
	GetTaskByDetection(ctx context.Context, status types.DatasetAnnotationStatus, detectionStatus types.DatasetAnnotationDetectionStatus, preload ...string) (res []types.DatasetAnnotationTask, err error)
	// GetSegmentFaqIntentInSegmentId 获取相应用segment的faq的意图数据
	GetSegmentFaqIntentInSegmentId(ctx context.Context, segmentIds []uint, annotationStatus types.DatasetAnnotationStatus, annotationType types.DatasetAnnotationType) (res []types.DatasetAnnotationTaskSegment, err error)
}

type service struct {
	db         *gorm.DB
	randomFunc string
}

func (s *service) GetSegmentFaqIntentInSegmentId(ctx context.Context, segmentIds []uint, annotationStatus types.DatasetAnnotationStatus, annotationType types.DatasetAnnotationType) (res []types.DatasetAnnotationTaskSegment, err error) {
	err = s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("segment_id in (?) and status = ? AND annotation_type = ?", segmentIds, annotationStatus, annotationType).
		Group("intent").Find(&res).Error
	return
}

func (s *service) GetTaskByDetection(ctx context.Context, status types.DatasetAnnotationStatus, detectionStatus types.DatasetAnnotationDetectionStatus, preload ...string) (res []types.DatasetAnnotationTask, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTask{}).Where("status = ? and detection_status = ?", status, detectionStatus)
	for _, p := range preload {
		query = query.Preload(p)
	}
	err = query.Find(&res).Error
	return
}

func (s *service) GetTaskSegmentPrev(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus) (res types.DatasetAnnotationTaskSegment, err error) {
	err = s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("data_annotation_id = ? and status = ?", taskId, status).
		Order("updated_at desc").
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return
}

func (s *service) GetTaskSegmentByRand(ctx context.Context, datasetId uint, percent float64, status types.DatasetAnnotationStatus, segmentType types.DatasetAnnotationSegmentType, preload ...string) (res []types.DatasetAnnotationTaskSegment, err error) {
	var total int64
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("data_annotation_id = ? and status = ? and segment_type = ?", datasetId, status, segmentType)
	if err = query.Count(&total).Error; err != nil {
		return nil, err
	}
	for _, p := range preload {
		query = query.Preload(p)
	}
	//limit := math.Round(percent * float64(total))

	err = query.Order(s.randomFunc).Limit(int(percent * float64(total))).Find(&res).Error
	return
}

func (s *service) DeleteDatasetDocumentById(ctx context.Context, id uint, unscoped bool) (err error) {
	if unscoped {
		return s.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&types.DatasetDocument{}).Error
	}
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&types.DatasetDocument{}).Error
}

func (s *service) AddDatasetDocumentSegments(ctx context.Context, data []types.DatasetDocumentSegment) (err error) {
	return s.db.WithContext(ctx).Model(types.DatasetDocumentSegment{}).Create(data).Error
}

func (s *service) UpdateDatasetDocumentSegmentCount(ctx context.Context, datasetDocumentId uint, count int) (err error) {
	return s.db.WithContext(ctx).Model(types.DatasetDocument{}).Where("id = ?", datasetDocumentId).
		Update("segment_count", count).Error
}

func (s *service) ListDatasetDocuments(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []types.DatasetDocument, total int64, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetDocument{}).Where("tenant_id = ?", tenantId)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Count(&total).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error
	return
}

func (s *service) CreateDatasetDocument(ctx context.Context, data *types.DatasetDocument) (err error) {
	return s.db.WithContext(ctx).Model(data).Create(data).Error
}

func (s *service) DeleteDatasetDocument(ctx context.Context, tenantId uint, uuid string) (err error) {
	return s.db.WithContext(ctx).Where("tenant_id = ? and uuid = ?", tenantId, uuid).Delete(&types.DatasetDocument{}).Error
}

func (s *service) GetDatasetDocumentSegmentByRange(ctx context.Context, datasetDocumentId uint, start, end int, preload ...string) (res []types.DatasetDocumentSegment, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetDocumentSegment{}).Where("dataset_document_id = ?", datasetDocumentId)
	for _, p := range preload {
		query = query.Preload(p)
	}
	err = query.Order("serial_number ASC").Offset(start).Limit(end - start).Find(&res).Error
	return
}

func (s *service) GetDatasetDocumentByUUID(ctx context.Context, tenantId uint, uuid string, preload ...string) (res *types.DatasetDocument, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetDocument{}).Where("tenant_id = ? and uuid = ?", tenantId, uuid)
	for _, p := range preload {
		query = query.Preload(p)
	}
	err = query.First(&res).Error
	return
}

func (s *service) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int, preloads ...string) (res []types.DatasetAnnotationTask, total int64, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTask{}).Where("tenant_id = ?", tenantId)
	for _, p := range preloads {
		query = query.Preload(p)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Count(&total).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Omit("detection_log").Find(&res).Error
	return
}

func (s *service) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	return s.db.WithContext(ctx).Where("tenant_id = ? and uuid = ?", tenantId, uuid).
		Delete(&types.DatasetAnnotationTask{}).Error
}

func (s *service) UpdateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	return s.db.WithContext(ctx).Model(data).Updates(data).Error
}

func (s *service) GetTask(ctx context.Context, tenantId uint, uuid string, preloads ...string) (res *types.DatasetAnnotationTask, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTask{}).Where("tenant_id = ? and uuid = ?", tenantId, uuid)
	for _, p := range preloads {
		query = query.Preload(p)
	}
	err = query.First(&res).Error
	return
}

func (s *service) AddTaskSegments(ctx context.Context, data []types.DatasetAnnotationTaskSegment) (err error) {
	return s.db.WithContext(ctx).Create(data).Error
}

func (s *service) GetTaskOneSegment(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	err = s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("data_annotation_id = ? and status = ?", taskId, status).
		Order(s.randomFunc).
		First(&res).Error
	return
}

func (s *service) GetTaskSegmentByUUID(ctx context.Context, taskId uint, uuid string, preload ...string) (res *types.DatasetAnnotationTaskSegment, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("data_annotation_id = ? and uuid = ?", taskId, uuid)
	for _, p := range preload {
		query = query.Preload(p)
	}
	err = query.First(&res).Error
	return
}

func (s *service) UpdateTaskSegment(ctx context.Context, data *types.DatasetAnnotationTaskSegment) (err error) {
	return s.db.WithContext(ctx).Model(data).Updates(data).Error
}

func (s *service) GetTaskSegments(ctx context.Context, taskId uint, status types.DatasetAnnotationStatus, page, pageSize int, preload ...string) (res []types.DatasetAnnotationTaskSegment, total int64, err error) {
	query := s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).
		Where("data_annotation_id = ? and status = ?", taskId, status)
	for _, p := range preload {
		query = query.Preload(p)
	}
	err = query.Count(&total).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error
	return
}

func (s *service) UpdateTaskSegmentType(ctx context.Context, segmentIds []uint, segmentType types.DatasetAnnotationSegmentType) (err error) {
	return s.db.WithContext(ctx).Model(types.DatasetAnnotationTaskSegment{}).Where("id in (?)", segmentIds).
		Update("segment_type", segmentType).Error
}

func (s *service) CreateTask(ctx context.Context, data *types.DatasetAnnotationTask) (err error) {
	return s.db.WithContext(ctx).Create(data).Error
}

func New(db *gorm.DB) Service {
	var randomFunc = "rand()"
	if db.Dialector.Name() == "sqlite" {
		randomFunc = "random()"
	}
	return &service{
		db:         db,
		randomFunc: randomFunc,
	}
}
