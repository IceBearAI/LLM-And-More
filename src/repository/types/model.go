package types

import (
	"gorm.io/gorm"
	"time"
)

// Models 模型表
type Models struct {
	gorm.Model
	ProviderName ModelProvider `gorm:"column:provider_name;type:varchar(50);default:localai;NOT NULL"`      // 模型供应商 openai、localai
	ModelType    ModelType     `gorm:"column:model_type;type:varchar(30);default:text-generation;NOT NULL"` // 模型类型 text-generation、embeddings、whisper
	ModelName    string        `gorm:"column:model_name;type:varchar(50);NOT NULL"`                         // 模型名称
	// 映射模型,如果有值则表示是映射模型,值为映射的模型名称
	BaseModelName      string             `gorm:"column:base_model_name;type:varchar(50);null;comment:映射模型名称"` // 映射模型名称
	MaxTokens          int                `gorm:"column:max_tokens;type:int(11);default:2048;NOT NULL"`        // 最长上下文
	IsPrivate          bool               `gorm:"column:is_private;type:tinyint(1);default:0;NOT NULL"`        // 是否是私有模型
	IsFineTuning       bool               `gorm:"column:is_fine_tuning;type:tinyint(1);default:0;NOT NULL"`    // 是否是微调模型
	Enabled            bool               `gorm:"column:enabled;type:tinyint(1);default:0"`                    // 是否启用
	Remark             string             `gorm:"column:remark;size:255;null;comment:备注"`
	ModelDeploy        ModelDeploy        `gorm:"foreignKey:ModelID;references:ID"`
	Tenants            []Tenants          `gorm:"many2many:tenant_model_associations;foreignKey:ID;references:ID;joinForeignKey:ModelID;joinReferences:TenantID"`
	TenantId           []uint             `gorm:"-"`
	FineTuningTrainJob FineTuningTrainJob `gorm:"foreignKey:FineTunedModel;references:ModelName"`
	Parameters         float64            `gorm:"column:parameters;type:decimal(7,2);default:0;NOT NULL"` // 模型参数量
	LastOperator       string             `gorm:"column:last_operator;size:100;null;comment:最后操作人"`
	Replicas           int                `gorm:"column:replicas;default:1;null;comment:并行/实例数量"`
	Label              string             `gorm:"column:label;size:500;null;comment:调度标签"`
	K8sCluster         string             `gorm:"column:k8s_cluster;size:500;null;comment:k8s集群"`
	InferredType       string             `gorm:"column:inferred_type;size:500;null;comment:推理类型cpu,gpu"`
	Gpu                int                `gorm:"column:gpu;default:0;null;comment:GPU数"`
	Cpu                int                `gorm:"column:cpu;default:0;null;comment:CPU核数"`
	Memory             int                `gorm:"column:memory;default:1;null;comment:内存G"`

	Channels []ChatChannels `gorm:"many2many:channel_model_associations;foreignKey:id;joinForeignKey:model_id;References:id;joinReferences:channel_id"`
}

func (m *Models) TableName() string {
	return "models"
}

func (m *Models) CanDelete() bool {
	switch m.ModelType {
	case ModelTypeTextGeneration:
		return m.ProviderName == ModelProviderLocalAI && (m.ModelDeploy.Status == "" || m.ModelDeploy.Status == ModelDeployStatusFailed.String())
	case ModelTypeDigitalhuman:
		return true
	case ModelTypeVoice:
		return true
	}
	return false

}

func (m *Models) CanDeploy() bool {
	switch m.ModelType {
	case ModelTypeTextGeneration:
		return m.ProviderName == ModelProviderLocalAI && (m.ModelDeploy.Status == "" || m.ModelDeploy.Status == ModelDeployStatusFailed.String()) && (m.BaseModelName == "")
	case ModelTypeDigitalhuman:
		return false
	case ModelTypeVoice:
		return false
	}
	return false
}

func (m *Models) CanUndeploy() bool {
	switch m.ModelType {
	case ModelTypeTextGeneration:
		return m.ProviderName == ModelProviderLocalAI && (m.ModelDeploy.Status == ModelDeployStatusPending.String() ||
			m.ModelDeploy.Status == ModelDeployStatusRunning.String() ||
			m.ModelDeploy.Status == ModelDeployStatusSuccess.String())
	case ModelTypeDigitalhuman:
		return false
	case ModelTypeVoice:
		return false
	}
	return false
}

// ModelDeploy 模型部署
type ModelDeploy struct {
	gorm.Model
	ModelID      uint64 `gorm:"column:model_id;type:bigint(20) unsigned;NOT NULL"` // 模型表主键 models.id
	ModelPath    string `gorm:"column:model_path;type:varchar(255);NOT NULL"`      // 模型部署路径
	Status       string `gorm:"column:status;type:varchar(32)"`                    // 部署状态
	Replicas     int    `gorm:"column:replicas;type:int(11);default:0;"`           // 实例数
	InferredType string `gorm:"column:inferred_type;type:varchar(50)"`             // 推理类型GPU/CPU
	Label        string `gorm:"column:label;type:varchar(50)"`                     // GPU调度标签
	Gpu          int    `gorm:"column:gpu;type:int(11);"`                          // 单实例GPU数量
	Cpu          int    `gorm:"column:cpu;type:int(11);"`                          // 单实例CPU数量
	Quantization string `gorm:"column:quantization;type:varchar(50)"`              // 量化
	Vllm         bool   `gorm:"column:vllm;default:false;"`                        // 是否开启VLLM
}

func (m *ModelDeploy) TableName() string {
	return "model_deploy"
}

type TenantModelAssociations struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TenantID  uint `gorm:"column:tenant_id;type:bigint(20) unsigned;NOT NULL"` // 租户ID
	ModelID   uint `gorm:"column:model_id;type:bigint(20) unsigned;NOT NULL"`  // 模型ID
}

func (m *TenantModelAssociations) TableName() string {
	return "tenant_model_associations"
}

type ModelProvider string
type ModelType string
type ModelDeployStatus string

const (
	ModelProviderOpenAI  ModelProvider = "OpenAI"
	ModelProviderLocalAI ModelProvider = "LocalAI"

	ModelTypeTextGeneration ModelType = "text-generation"
	ModelTypeEmbeddings     ModelType = "embeddings"
	ModelTypeWhisper        ModelType = "whisper"
	ModelTypeDigitalhuman   ModelType = "digitalhuman"
	ModelTypeVoice          ModelType = "voice"

	ModelDeployStatusPending ModelDeployStatus = "pending"
	ModelDeployStatusRunning ModelDeployStatus = "running"
	ModelDeployStatusSuccess ModelDeployStatus = "success"
	ModelDeployStatusFailed  ModelDeployStatus = "failed"
)

func (m ModelProvider) String() string {
	return string(m)
}

func (m ModelType) String() string {
	return string(m)
}

func (m ModelDeployStatus) String() string {
	return string(m)
}
