package types

import "gorm.io/gorm"

type EvaluateStatus string
type EvaluateDataType string
type EvaluateTargetType string

const (
	EvaluateStatusRunning EvaluateStatus = "running" // 运行中
	EvaluateStatusSuccess EvaluateStatus = "success" // 成功
	EvaluateStatusFailed  EvaluateStatus = "failed"  // 失败
	EvaluateStatusWaiting EvaluateStatus = "waiting" // 等待
	EvaluateStatusCancel  EvaluateStatus = "cancel"  // 取消

	EvaluateDataTypeCustom EvaluateDataType = "custom" //自定义数据集
	EvaluateDataTypeTrain  EvaluateDataType = "train"  //训练数据集

	EvaluateTargetTypeACC   EvaluateTargetType = "Acc"   //ACC
	EvaluateTargetTypeF1    EvaluateTargetType = "F1"    //F1
	EvaluateTargetTypeBleu  EvaluateTargetType = "BLEU"  //Bleu
	EvaluateTargetTypeRouge EvaluateTargetType = "Rouge" //Rouge
	EvaluateTargetTypeFive  EvaluateTargetType = "five"  //五维图
)

// ModelEvaluate 模型评测
type ModelEvaluate struct {
	gorm.Model
	ModelId        int     `gorm:"column:model_id;type:bigint(20);NOT NULL"`                  // 模型表主键 models.id
	ModelPath      string  `gorm:"column:model_path;type:varchar(255);NOT NULL"`              // 模型部署路径
	Status         string  `gorm:"column:status;type:varchar(32)"`                            // 状态
	StatusMsg      string  `gorm:"column:status_msg;type:varchar(1000)"`                      // 状态原因
	Replicas       int     `gorm:"column:replicas;type:int(11);default:0;"`                   // 实例数
	InferredType   string  `gorm:"column:inferred_type;type:varchar(50)"`                     // 推理类型GPU/CPU
	Label          string  `gorm:"column:label;type:varchar(50)"`                             // GPU调度标签
	Gpu            int     `gorm:"column:gpu;type:int(11);"`                                  // 单实例GPU数量
	Cpu            int     `gorm:"column:cpu;type:int(11);"`                                  // 单实例CPU数量
	MaxGpuMemory   int     `gorm:"column:max_gpu_memory;type:int(11);"`                       // GPU最大内存
	Memory         int     `gorm:"column:memory;type:int(11);"`                               // 内存G
	EvalTargetType string  `gorm:"column:eval_target_type;type:varchar(100);"`                // 评测指标类型
	MaxLength      int     `gorm:"column:max_length;type:int(11);"`                           // 最大输出序列长度
	BatchSize      int     `gorm:"column:batch_size;type:int(11);"`                           // 单卡Batch大小
	DataSize       int     `gorm:"column:data_size;type:int(11);"`                            // 数据量
	DataType       string  `gorm:"column:data_type;type:varchar(50);"`                        // 数据集类型
	DataFileId     string  `gorm:"column:data_file_id;type:varchar(500)"`                     // 数据集fileId
	DataUrl        string  `gorm:"column:data_url;type:varchar(500)"`                         // 数据集地址
	Remark         string  `gorm:"column:remark;type:varchar(1000)"`                          // 备注
	CompleteRate   float64 `gorm:"column:complete_rate;type:decimal(7,2);default:0;NOT NULL"` // 进度
	OperatorEmail  string  `gorm:"column:operator_email;type:varchar(50);not null;"`          // 操作人邮箱
	RiskOver       bool    `gorm:"column:risk_over;null;default:false;"`                      // 过拟合风险
	RiskUnder      bool    `gorm:"column:risk_under;null;default:false;"`                     // 欠拟合风险
	RiskDisaster   bool    `gorm:"column:risk_disaster;null;default:false;"`                  // 灾难性遗忘
	Score          float64 `gorm:"column:score;type:decimal(7,2);default:0;NOT NULL"`         // 分数
	Five1          float64 `gorm:"column:five1;type:decimal(7,2);default:0;NOT NULL"`         // 中文能力
	Five2          float64 `gorm:"column:five2;type:decimal(7,2);default:0;NOT NULL"`         // 推理能力
	Five3          float64 `gorm:"column:five3;type:decimal(7,2);default:0;NOT NULL"`         // 指令遵从能力
	Five4          float64 `gorm:"column:five4;type:decimal(7,2);default:0;NOT NULL"`         // 创新能力
	Five5          float64 `gorm:"column:five5;type:decimal(7,2);default:0;NOT NULL"`         // 阅读理解
	Models         Models  `gorm:"foreignKey:ModelId;references:ID"`                          //关联Models
	Uuid           string  `gorm:"column:uuid;type:varchar(500);"`                            // JobId uuid
	JobName        string  `gorm:"column:job_name;type:varchar(500);"`                        // JobName
	Result         string  `gorm:"column:result;type:longtext;"`                              // 回调结果
	EvaluateLog    string  `gorm:"column:evaluate_log;type:longtext;null"`                    // 回调结果
}

func (m *ModelEvaluate) TableName() string {
	return "model_evaluate"
}
