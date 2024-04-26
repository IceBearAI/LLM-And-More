package types

import (
	"gorm.io/gorm"
	"time"
)

type TemplateType string
type TrainStatus string
type ScenarioType string

const (
	TemplateTypeTrain     TemplateType = "train"     // 训练模版
	TemplateTypeInference TemplateType = "inference" // 推理模版

	TrainStatusRunning TrainStatus = "running" // 运行中
	TrainStatusSuccess TrainStatus = "success" // 成功
	TrainStatusFailed  TrainStatus = "failed"  // 失败
	TrainStatusWaiting TrainStatus = "waiting" // 等待
	TrainStatusCancel  TrainStatus = "cancel"  // 取消

	// ScenarioTypeGeneral 通用应用场景
	ScenarioTypeGeneral ScenarioType = "general"
	// ScenarioTypeFAQ 知识问答应用场景
	ScenarioTypeFAQ ScenarioType = "faq"
	// ScenarioTypeRAG RAG应用场景
	ScenarioTypeRAG ScenarioType = "rag"
)

func (t TrainStatus) String() string {
	return string(t)
}

// FineTuningTemplate 微调模版
type FineTuningTemplate struct {
	gorm.Model
	Name          string       `gorm:"column:name;size:64;not null;unique;index;"` // comment:名称
	BaseModel     string       `gorm:"column:base_model;size:64;not null;index;"`  // comment:模型
	Content       string       `gorm:"column:content;type:longtext;not null;"`     // comment:脚本模版
	Params        string       `gorm:"column:params;text;null;"`                   // comment:模版所需要参数
	TrainImage    string       `gorm:"column:train_image;size:500;not null;"`      // comment:训练镜像
	Remark        string       `gorm:"column:remark;size:500;null;"`               // comment:备注
	BaseModelPath string       `gorm:"column:base_model_path;size:500;null;"`      // comment:基础模型路径
	ScriptFile    string       `gorm:"column:script_file;size:500;null;"`          // comment:脚本文件
	OutputDir     string       `gorm:"column:output_dir;size:500;null;"`           // comment:输出目录
	MaxTokens     int          `gorm:"column:max_tokens;default:2048;null;"`       // comment:最大token数
	Lora          bool         `gorm:"column:lora;null;default:false;"`            // comment:是否使用lora微调
	Enabled       bool         `gorm:"column:enabled;default:false;"`              // comment:可用状态
	ParallelNum   int          `gorm:"column:parallel_num;default:1;null;"`        // comment:并行数量
	GpuLabel      string       `gorm:"column:gpu_label;size:500;null;"`            // comment:gpu标签
	K8sCluster    string       `gorm:"column:k8s_cluster;size:500;null;"`          // comment:k8s集群
	Cpu           int          `gorm:"column:cpu;default:1;null;"`                 // comment:CPU核数
	Memory        int          `gorm:"column:memory;default:1;null;"`              // comment:内存G
	TemplateType  TemplateType `gorm:"column:template_type;size:24;null;"`         // comment:模版类型
}

// FineTuningTrainJob 微调训练任务
type FineTuningTrainJob struct {
	// FineTuningTrainJob 微调训练任务
	gorm.Model
	// Name 名称
	JobId string `gorm:"column:job_id;size:64;not null;unique;index;"`
	// FineTunedModel 微调模型
	FineTunedModel string `gorm:"column:fine_tuned_model;size:64;not null;unique;index;"`
	// ChannelId 渠道ID
	ChannelId uint `gorm:"column:channel_id;size:64;not null;index;"`
	// TemplateId 模版ID
	TemplateId uint `gorm:"column:template_id;size:64;not null;index;"`
	// 文件ID
	FileId string `gorm:"column:file_id;size:64;not null;index;"`
	// BaseModel 模型名称
	BaseModel string `gorm:"column:base_model;size:128;not null;index;"`
	// TrainBatchSize 训练批次 default 1
	TrainBatchSize int `gorm:"column:train_batch_size;null;default:1;"`
	// EvalBatchSize 评估批次 default 1
	EvalBatchSize int `gorm:"column:eval_batch_size;null;default:1;"`
	// AccumulationSteps 梯度累加步数 default 1
	AccumulationSteps int `gorm:"column:accumulation_steps;null;default:1;"`
	// TrainEpoch 训练轮次 default 1
	TrainEpoch int `gorm:"column:train_epoch;null;default:1;"`
	// ProcPerNode 每个节点使用GPU数量 default 1
	ProcPerNode int `gorm:"column:proc_per_node;null;default:1;"`
	// EvalSteps 1500
	EvalSteps int `gorm:"column:eval_steps;null;default:1500;"`
	// SaveSteps 1500
	SaveSteps int `gorm:"column:save_steps;null;default:1500;"`
	// SaveTotalLimit default 8
	SaveTotalLimit int `gorm:"column:save_total_limit;null;default:8;"`
	// LearningRate default 2e-5
	LearningRate float64 `gorm:"column:learning_rate;null;default:2e-5;"`
	// WeightDecay default 0.
	WeightDecay float32 `gorm:"column:weight_decay;null;default:0.;"`
	// WarmupRatio default 0.04
	WarmupRatio float32 `gorm:"column:warmup_ratio;null;default:0.04;"`
	// LoggingSteps default 1
	LoggingSteps int `gorm:"column:logging_steps;null;default:1;"`
	// ModelMaxLength default 2048
	ModelMaxLength int `gorm:"column:model_max_length;null;default:2048;"`
	// 是否微调lora
	Lora bool `gorm:"column:lora;null;default:false;"`
	// BaseModelPath 基础模型路径
	BaseModelPath string `gorm:"column:base_model_path;size:500;null;"`
	// DataPath 数据集文件地址
	DataPath string `gorm:"column:data_path;size:500;null;"`
	// OutputDir 输出目录
	OutputDir string `gorm:"column:output_dir;size:500;null;"`
	// ScriptFile 脚本文件
	ScriptFile string `gorm:"column:script_file;size:500;null;"`
	// MasterPort master端口
	MasterPort int `gorm:"column:master_port;null;"`
	// FileUrl 文件地址
	FileUrl string `gorm:"column:file_url;null;"`
	// Suffix 后缀
	Suffix string `gorm:"column:suffix;null;"`
	// ValidationFile 验证文件
	ValidationFile string `gorm:"column:validation_file;null;"`
	// TrainStatus 训练状态
	TrainStatus TrainStatus `gorm:"column:train_status;null;"`
	// TrainDuration 训练时长 单位秒
	TrainDuration int `gorm:"column:train_duration;null;"`
	// Progress 训练进度
	Progress float64 `gorm:"column:progress;null;"`
	// ProgressEpochs 训练轮次
	ProgressEpochs float64 `gorm:"column:progress_epochs;null;"`
	// ProgressLoss 训练损失
	ProgressLoss float64 `gorm:"column:progress_loss;null;"`
	// ProgressLearningRate 学习率
	ProgressLearningRate float64 `gorm:"column:progress_learning_rate;null;"`
	// TrainPublisher 训练发布者
	TrainPublisher string `gorm:"column:train_publisher;null;"`
	// TrainScript 训练脚本内容
	TrainScript string `gorm:"column:train_script;type:text;null;"`
	// PaasJobName paas job name
	PaasJobName string `gorm:"column:paas_job_name;size:64;null;"`
	// FinishedAt 完成时间
	FinishedAt *time.Time `gorm:"column:finished_at;null;"`
	// ErrorMessage 错误信息
	ErrorMessage string `gorm:"column:error_message;type:text;null;"`
	// Remark 备注
	Remark string `gorm:"column:remark;size:128;null;"`
	// TrainLog 训练日志
	TrainLog string `gorm:"column:train_log;type:longtext;null;"`
	// TenantID 租户ID
	TenantID uint `gorm:"column:tenant_id;type:bigint(20);NOT NULL"`
	// StartTrainTime 开始训练时间
	StartTrainTime *time.Time `gorm:"column:start_train_time;null;"`
	// Scenario 应用场景
	Scenario ScenarioType `gorm:"column:scenario;null;default:'general'"`
	// Diagnosis 诊断
	Diagnosis string `gorm:"column:diagnosis;null;"`

	// Template 微调模版
	Template FineTuningTemplate `gorm:"foreignKey:TemplateId;references:ID"`
	// Files 文件
	FineTuningFile Files `gorm:"foreignKey:FileId;references:file_id"`
	// BaseModelInfo 基础模型
	BaseModelInfo *Models `gorm:"foreignKey:BaseModel;references:ModelName"`
}

// TableName sets the insert table name for this struct type
func (c *FineTuningTemplate) TableName() string {
	return "fine_tuning_template"
}

// TableName sets the insert table name for this struct type
func (c *FineTuningTrainJob) TableName() string {
	return "fine_tuning_train_job"
}

func (c *FineTuningTrainJob) CanCancel() bool {
	return c.TrainStatus == TrainStatusRunning || c.TrainStatus == TrainStatusWaiting
}

func (c *FineTuningTrainJob) CanDelete() bool {
	return c.TrainStatus == TrainStatusFailed || c.TrainStatus == TrainStatusCancel
}
