package llmeval

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(Service) Service

type Service interface {
	// GetEarliestWaitingEvalTask 获取最早一个正在等待的评估任务
	GetEarliestWaitingEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error)
	// GetRunningEvalTask 获取正在运行的评估任务
	GetRunningEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error)
	// UpdateEvalStatus 更新评估任务状态
	UpdateEvalStatus(ctx context.Context, evalId uint, status types.EvalStatus, errMessage string) (err error)
	// UpdateEvalProgress 更新评估任务进度
	UpdateEvalProgress(ctx context.Context, evalId uint, score, progress float64, status types.EvalStatus) (err error)
	// UpdateEvalStartTime 更新评估任务开始时间
	UpdateEvalStartTime(ctx context.Context, evalId uint) (err error)
	// UpdateEvalDetail 更新评估任务详情
	UpdateEvalDetail(ctx context.Context, evalId uint, detail string) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) UpdateEvalDetail(ctx context.Context, evalId uint, detail string) (err error) {
	return s.db.WithContext(ctx).Model(&types.LLMEvalResults{}).
		Where("id = ?", evalId).Updates(map[string]interface{}{
		"details": detail,
	}).Error
}

func (s *service) UpdateEvalStartTime(ctx context.Context, evalId uint) (err error) {
	return s.db.WithContext(ctx).Model(&types.LLMEvalResults{}).
		Where("id = ?", evalId).Updates(map[string]interface{}{
		"started_at": gorm.Expr("now()"),
	}).Error
}

func (s *service) GetEarliestWaitingEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("status = ?", types.EvalStatusPending).Order("created_at asc").First(&evalTask).Error
	return
}

func (s *service) UpdateEvalProgress(ctx context.Context, evalId uint, score, progress float64, status types.EvalStatus) (err error) {
	return s.db.WithContext(ctx).Model(&types.LLMEvalResults{}).
		Where("id = ?", evalId).Updates(map[string]interface{}{
		"score":    score,
		"progress": progress,
		"status":   status,
	}).Error
}

func (s *service) UpdateEvalStatus(ctx context.Context, evalId uint, status types.EvalStatus, errMessage string) (err error) {
	return s.db.WithContext(ctx).Model(&types.LLMEvalResults{}).
		Where("id = ?", evalId).Updates(map[string]interface{}{
		"status":        status,
		"error_message": errMessage,
		"progress":      0,
	}).Error
}

func (s *service) GetRunningEvalTask(ctx context.Context, preloads ...string) (evalTask types.LLMEvalResults, err error) {
	db := s.db.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err = db.Where("status = ?", types.EvalStatusRunning).First(&evalTask).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
