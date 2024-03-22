package types

import (
	"gorm.io/gorm"
	"time"
)

// DatasetAnnotationType is a string enum type
type DatasetAnnotationType string

// DatasetAnnotationStatus is a string enum type
type DatasetAnnotationStatus string

// DatasetAnnotationDetectionStatus is a string enum type
type DatasetAnnotationDetectionStatus string

// DatasetAnnotationSegmentType is a string enum type
type DatasetAnnotationSegmentType string

const (
	// DatasetAnnotationTypeRAG RAG类型
	DatasetAnnotationTypeRAG DatasetAnnotationType = "rag"
	// DatasetAnnotationTypeFAQ 问答类型
	DatasetAnnotationTypeFAQ DatasetAnnotationType = "faq"
	// DatasetAnnotationTypeGeneral 通用类型
	DatasetAnnotationTypeGeneral DatasetAnnotationType = "general"

	// DatasetAnnotationStatusPending 等待中
	DatasetAnnotationStatusPending DatasetAnnotationStatus = "pending"
	// DatasetAnnotationStatusProcessing 处理中
	DatasetAnnotationStatusProcessing DatasetAnnotationStatus = "processing"
	// DatasetAnnotationStatusCompleted 已完成
	DatasetAnnotationStatusCompleted DatasetAnnotationStatus = "completed"
	// DatasetAnnotationStatusAbandoned 已废弃
	DatasetAnnotationStatusAbandoned DatasetAnnotationStatus = "abandoned"
	// DatasetAnnotationStatusCleaned 已清除
	DatasetAnnotationStatusCleaned DatasetAnnotationStatus = "cleaned"

	// DatasetAnnotationSegmentTypeTrain 训练集
	DatasetAnnotationSegmentTypeTrain DatasetAnnotationSegmentType = "train"
	// DatasetAnnotationSegmentTypeTest 测试集
	DatasetAnnotationSegmentTypeTest DatasetAnnotationSegmentType = "test"

	// DatasetAnnotationDetectionStatusPending 等待中
	DatasetAnnotationDetectionStatusPending DatasetAnnotationDetectionStatus = "pending"
	// DatasetAnnotationDetectionStatusProcessing 处理中
	DatasetAnnotationDetectionStatusProcessing DatasetAnnotationDetectionStatus = "processing"
	// DatasetAnnotationDetectionStatusCompleted 已完成
	DatasetAnnotationDetectionStatusCompleted DatasetAnnotationDetectionStatus = "completed"
	// DatasetAnnotationDetectionStatusCanceled 已取消
	DatasetAnnotationDetectionStatusCanceled DatasetAnnotationDetectionStatus = "canceled"
)

// DatasetAnnotationTaskSegment is a struct type
type DatasetAnnotationTaskSegment struct {
	gorm.Model
	// UUID 样本ID
	UUID string `gorm:"column:uuid;type:string;size:64;uniqueIndex;not null;"`
	// DataAnnotationID 标注任务ID
	DataAnnotationID uint `gorm:"column:data_annotation_id;index;"` // comment:标注任务ID
	// SegmentID 样本ID
	SegmentID      uint                         `gorm:"column:segment_id;type:int;index;"`             // comment:样本ID
	AnnotationType DatasetAnnotationType        `gorm:"column:annotation_type;type:varchar(12);index"` // comment:标注类型
	SegmentContent string                       `gorm:"column:segment_content;type:text;"`             // comment:样本内容
	Document       string                       `gorm:"column:document;size:2000;null;"`               // comment:标注文本
	Instruction    string                       `gorm:"column:instruction;size:2000;null;"`            // comment:标注说明
	Input          string                       `gorm:"column:input;size:2000;null;"`                  // comment:标注输入
	Question       string                       `gorm:"column:question;size:2000;null;"`               // comment:标注问题
	Intent         string                       `gorm:"column:intent;size:32;null;"`                   // comment:标注意图
	Output         string                       `gorm:"column:output;size:2000;null;"`                 // comment:输出结果
	Status         DatasetAnnotationStatus      `gorm:"column:status;size:12;index;"`                  // comment:标注状态
	SegmentType    DatasetAnnotationSegmentType `gorm:"column:segment_type;size:12;index;"`            // comment:样本类型
	CreatorEmail   string                       `gorm:"column:creator_email;size:32;"`                 // comment:创建人邮箱

	DataAnnotationTask DatasetAnnotationTask `gorm:"foreignKey:DataAnnotationID;references:ID"`
}

// DatasetAnnotationTask is a struct type
type DatasetAnnotationTask struct {
	gorm.Model
	DatasetDocumentId uint                             `gorm:"column:dataset_document_id;type:int;index;"`       // comment:数据集ID
	UUID              string                           `gorm:"column:uuid;type:string;size:64;uniqueIndex;"`     // comment:标注任务ID
	Name              string                           `gorm:"column:name;type:string;size:64;index;"`           // comment:任务名称
	TenantID          uint                             `gorm:"column:tenant_id;index;"`                          // comment:租户ID
	Principal         string                           `gorm:"column:principal;size:32;index;"`                  // comment:负责人
	AnnotationType    string                           `gorm:"column:annotation_type;size:12;index;"`            // comment:标注类型
	Status            DatasetAnnotationStatus          `gorm:"column:status;size:12;index;default:pending;"`     // comment:标注状态
	CompletedAt       *time.Time                       `gorm:"column:completed_at;null;"`                        // comment:完成时间
	DataSequence      string                           `gorm:"column:data_sequence;size:12;"`                    // comment:数据序列
	Total             int                              `gorm:"column:total;type:int;default:0;"`                 // comment:需要标的数据总量
	Completed         int                              `gorm:"column:completed;default:0;"`                      // comment:完成标注量
	Abandoned         int                              `gorm:"column:abandoned;default:0;"`                      // comment:废弃标注量
	TrainTotal        int                              `gorm:"column:train_total;default:0;"`                    // comment:训练数据总量
	TestTotal         int                              `gorm:"column:test_total;default:0;"`                     // comment:测试数据总量
	Remark            string                           `gorm:"column:remark;size:1000;"`                         // comment:备注
	TestReport        string                           `gorm:"column:test_report;type:text;"`                    // comment:测试数据仓库
	JobName           string                           `gorm:"column:job_name;size:64;null;"`                    // comment:任务名称
	DetectionStatus   DatasetAnnotationDetectionStatus `gorm:"column:detection_status;size:12;default:pending;"` // comment:检测状态
	Segments          []DatasetAnnotationTaskSegment   `gorm:"foreignKey:DataAnnotationID;references:ID"`

	DatasetDocument DatasetDocument `gorm:"foreignKey:DatasetDocumentId;references:ID"`
}

// DatasetDocument is a struct type
type DatasetDocument struct {
	gorm.Model
	UUID         string `gorm:"column:uuid;type:string;size:64;uniqueIndex;not null;"` // comment:文档ID
	Name         string `gorm:"column:name;type:string;size:64;index;"`                // comment:文档名称
	Remark       string `gorm:"column:remark;size:500;null;"`                          // comment:描述
	SegmentCount int    `gorm:"column:segment_count;type:int;default:0;"`              // comment:样本数量
	CreatorEmail string `gorm:"column:creator_email;size:64;index;"`                   // comment:创建人邮箱
	TenantID     uint   `gorm:"column:tenant_id;index;"`                               // comment:租户ID
	FormatType   string `gorm:"column:format;size:12;null;default:alpaca;"`            // comment:数据集类型
	SplitType    string `gorm:"column:split_type;size:12;null;default:'\n\n';"`        // comment:切割方式
	SplitMax     int    `gorm:"column:split_max;type:int;default:1000;"`               // comment:切割的最大数据块
	FileName     string `gorm:"column:file_name;size:64;null;"`                        // comment:文件名

	DatasetDocumentSegments []DatasetDocumentSegment `gorm:"foreignKey:DatasetDocumentId;references:ID"`
}

// DatasetDocumentSegment is a struct type
type DatasetDocumentSegment struct {
	gorm.Model
	UUID              string `gorm:"column:uuid;type:string;size:64;uniqueIndex;not null;"` // comment:样本ID
	DatasetDocumentId uint   `gorm:"column:dataset_document_id;type:int;index;"`            // comment:文档ID
	SegmentContent    string `gorm:"column:segment_content;type:longtext;"`                 // comment:样本内容
	WordCount         int    `gorm:"column:word_count;type:int;"`                           // comment:字数
	SerialNumber      int    `gorm:"column:serial_number;type:int;index;"`                  // comment:序号

	DatasetDocument DatasetDocument `gorm:"foreignKey:DatasetDocument;references:ID"`
}

// TableName sets the insert table name for this struct type
func (c *DatasetAnnotationTask) TableName() string {
	return "dataset_annotations"
}

// TableName sets the insert table name for this struct type
func (c *DatasetAnnotationTaskSegment) TableName() string {
	return "dataset_annotation_segments"
}

// TableName sets the insert table name for this struct type
func (c *DatasetDocument) TableName() string {
	return "dataset_documents"
}
