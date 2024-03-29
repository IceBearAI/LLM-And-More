package types

import (
	"gorm.io/gorm"
)

// ChatChannels 渠道表
type ChatChannels struct {
	gorm.Model
	Name          string   `gorm:"column:name;size:64;not null;unique;index;" json:"name"`        // comment:名称
	Alias         string   `gorm:"column:alias;size:64;null;default:'New Channel';" json:"alias"` // comment:渠道名称
	Remark        string   `gorm:"column:remark;size:128;null;" json:"remark"`                    // comment:备注
	Quota         int      `gorm:"column:quota;null;default:10;" json:"quota"`                    // comment:配额
	Models        string   `gorm:"column:models;size:255;null;" json:"models"`                    // comment:支持模型
	OnlyOpenAI    bool     `gorm:"column:only_openai;null;default:false;" json:"only_openai"`     // comment:仅使用openai
	ApiKey        string   `gorm:"column:api_key;index;unique;size:128;" json:"api_key"`          // comment:ApiKey
	Email         string   `gorm:"column:email;size:128;null;" json:"email"`                      // comment:邮箱
	LastOperator  string   `gorm:"column:last_operator;size:100;null;" json:"last_operator"`      // comment:最后操作人
	TenantId      uint     `gorm:"column:tenant_id;type:bigint(20);NOT NULL"`                     // 租户ID
	ChannelModels []Models `gorm:"many2many:channel_model_associations;foreignKey:id;joinForeignKey:channel_id;References:id;joinReferences:model_id"`
	ModelId       []uint   `gorm:"-" json:"modelId"`
	Tenant        Tenants  `gorm:"foreignKey:TenantId;references:ID"`
}

func (*ChatChannels) TableName() string {
	return "chat_channels"
}

// ChannelModelAssociations 渠道和模型中间表
type ChannelModelAssociations struct {
	//gorm.Model
	ChannelID uint `gorm:"column:channel_id;type:bigint(20);NOT NULL"` // 渠道表主键ID channels.id
	ModelID   uint `gorm:"column:model_id;type:bigint(20)"`            // 模型表主键ID models.id
}

func (m *ChannelModelAssociations) TableName() string {
	return "channel_model_associations"
}
