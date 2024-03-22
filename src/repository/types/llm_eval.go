package types

import "gorm.io/gorm"

type EvalStatus string
type EvalDataSetType string

const (
	// EvalStatusPending 评估状态：等待评估
	EvalStatusPending EvalStatus = "pending"
	// EvalStatusRunning 评估状态：正在评估
	EvalStatusRunning EvalStatus = "running"
	// EvalStatusSuccess 评估状态：评估成功
	EvalStatusSuccess EvalStatus = "success"
	// EvalStatusFailed 评估状态：评估失败
	EvalStatusFailed EvalStatus = "failed"
	// EvalStatusCancel 评估状态：评估取消
	EvalStatusCancel EvalStatus = "cancel"
	// EvalDataSetTypeTrain 评估数据集类型：训练集
	EvalDataSetTypeTrain EvalDataSetType = "train"
	// EvalDataSetTypeCustom 评估数据集类型：自定义
	EvalDataSetTypeCustom EvalDataSetType = "custom"
)

func (c EvalStatus) String() string {
	return string(c)
}

func (c EvalDataSetType) String() string {
	return string(c)
}

// LLMEvalDatasetSamples 数据集样本的表
type LLMEvalDatasetSamples struct {
	gorm.Model
	// 数据集文件ID
	DatasetFileId uint `gorm:"column:dataset_file_id;NOT NULL;"` // 数据集文件ID

}

// LLMEvalResults 存储评估结果的表
type LLMEvalResults struct {
	gorm.Model
	// 评估状态
	Status EvalStatus `gorm:"column:status;size:20;NOT NULL;DEFAULT:'pending';index;"` // 评估状态
	// 模型名称
	ModelName string `gorm:"column:model_name;type:varchar(50);NOT NULL;index;"` // 模型名称
	// 评估指标名称
	MetricName string `gorm:"column:metric_name;type:varchar(50);NOT NULL;"` // 评估指标名称
	// 评估指标分数
	Score float64 `gorm:"column:score;type:decimal(5,2);NOT NULL;"` // 评估指标分数
	// 评估结果详情
	Details *string `gorm:"column:details;type:json;NULL;"` // 评估结果详情
	// 评估数据集文件ID
	DatasetFileId uint `gorm:"column:dataset_file_id;NULL;"` // 评估数据集文件ID
	// Progress 进度
	Progress float64 `gorm:"column:progress;null;"`
	// EvalTotal 评估总数
	EvalTotal int `gorm:"column:eval_total;null;"`
	// ErrorMessage 错误信息
	ErrorMessage string `gorm:"column:error_message;size:1000;null;"`
	// CompletedAt 评估完成时间
	CompletedAt *gorm.DeletedAt `gorm:"column:completed_at;NULL;"` // 评估完成时间
	// StartedAt 评估开始时间
	StartedAt *gorm.DeletedAt `gorm:"column:started_at;NULL;"` // 评估开始时间
	// EvalDataSetType 评估数据集类型
	DatasetType string `gorm:"column:dataset_type;size:50;NOT NULL;DEFAULT:'train';"` // 评估数据集类型
	// Remark
	Remark      string `gorm:"column:remark;size:255;null;"`
	DatasetFile Files  `gorm:"foreignKey:DatasetFileId;references:ID"`
	UUid        string `gorm:"column:uuid;size:50;null;comment:uuid"`
}

// TableName sets the insert table name for this struct type
func (c *LLMEvalResults) TableName() string {
	return "llm_eval_results"
}
