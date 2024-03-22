package types

import "gorm.io/gorm"

// DatasetType 数据集类型
type DatasetType string

// DatasetSampleType 数据样本类型
type DatasetSampleType string

// DatasetFormat 数据集格式
type DatasetFormat string

const (
	DatasetTypeText  DatasetType = "text"  // 文本
	DatasetTypeImage DatasetType = "image" // 图片
	DatasetTypeAudio DatasetType = "audio" // 音频

	DatasetFormatAlpaca DatasetFormat = "alpaca"

	DatasetSampleTypeAlpaca DatasetSampleType = "alpaca" // alpaca格式
)

// Dataset 数据集
type Dataset struct {
	gorm.Model
	// Name 名称
	Name string `gorm:"column:name;size:64;not null;unique;index;"`
	// Remark 描述
	Remark string `gorm:"column:remark;size:500;null"`
	// UUID UUID
	UUID string `gorm:"column:uuid;size:64;not null;unique;index;"`
	// SampleCount 样本数量
	SampleCount int `gorm:"column:sample_count;null;default:0;"`
	// CreatorEmail 创建人邮箱
	CreatorEmail string `gorm:"column:creator_email;size:64;not null;index;"`
	// TenantId 租户ID
	TenantId uint `gorm:"column:tenant_id;size:64;not null;index;"`
	// 数据集类型
	Type DatasetType `gorm:"column:type;size:24;not null;index;"`
	// 数据集格式类型
	Format DatasetFormat `gorm:"column:format;size:24;null;default:'alpaca';"`

	// 数据集样本
	Samples []DatasetSample `gorm:"foreignKey:DatasetId;references:ID"`
}

// DatasetSample 数据集样本
type DatasetSample struct {
	gorm.Model
	// UUID UUID
	UUID string `gorm:"column:uuid;size:64;not null;unique;index;"`
	// Title 样本标题
	Title string `gorm:"column:title;size:64;not null;index;"`
	// DatasetId 数据集ID
	DatasetId uint `gorm:"column:dataset_id;size:64;not null;index;"`
	// System 系统内容
	System string `gorm:"column:system;size:1000;null;"`
	// Instruction 意图
	Instruction string `gorm:"column:instruction;size:2000;null;"`
	// Input 输入
	Input string `gorm:"column:input;type:text;null;"`
	// Output 输出
	Output string `gorm:"column:output;size:5000;null;"`
	// Conversations 内容
	Conversations string `gorm:"column:conversations;type:text;null;"`
	// Label 标签
	Label string `gorm:"column:label;size:64;null;"`
	// Remark 备注
	Remark string `gorm:"column:remark;size:500;null;"`
	// CreatorEmail 创建人邮箱
	CreatorEmail string `gorm:"column:creator_email;size:64;not null;index;"`
	// Turns 对话轮次
	Turns int `gorm:"column:turns;null;default:0;"`
	// Type 数据样本类型
	//Type DatasetSampleType `gorm:"column:type;size:24;null;index;default:'alpaca';comment:数据样本类型"`
}

// TableName sets the insert table name for this struct type
func (c *Dataset) TableName() string {
	return "datasets"
}

// TableName sets the insert table name for this struct type
func (c *DatasetSample) TableName() string {
	return "dataset_samples"
}
