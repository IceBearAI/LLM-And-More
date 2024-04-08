package types

import "gorm.io/gorm"

const (
	ChatBotTypeDefault    ChatBotType      = 1 // 默认智能体
	ChatBotTypeSystem     ChatBotType      = 2 // 系统内置智能体
	ChatBotTypeUser       ChatBotType      = 3 // 用户自定义智能体
	ChatBotModelTypeText  ChatBotModelType = 1 // 智能体模型类型 文本
	ChatBotModelTypeImage ChatBotModelType = 2 // 智能体模型类型 图片
	ChatBotModelTypeFile  ChatBotModelType = 3 // 智能体模型类型 文件
)

type ChatBotType int
type ChatBotModelType int

func (c ChatBotModelType) Int() int {
	return int(c)
}
func (c ChatBotType) Int() int {
	return int(c)
}

// ChatBot 会话智能体
type ChatBot struct {
	gorm.Model
	Name                string           `gorm:"column:name;type:varchar(100);NOT NULL"`                   // 名称
	Email               string           `gorm:"column:email;size:128;not null;index;" json:"email"`       // comment:邮箱
	DescriptionForHuman string           `gorm:"column:description_for_human;type:varchar(2000);NOT NULL"` // 智能体描述
	DescriptionForModel string           `gorm:"column:description_for_model;type:varchar(2000)"`          // 智能体针对模型描述(作为系统提示词使用)
	PrivateStatus       int              `gorm:"column:private_status;type:tinyint(4);default:1;NOT NULL"` // 私有状态 1.公开 2.不公开
	IconUrl             string           `gorm:"column:icon_url;type:varchar(255);NOT NULL"`               // 图标地址
	BotType             ChatBotType      `gorm:"column:bot_type;type:tinyint(4);default:2;NOT NULL"`       // 智能体类型 1.默认 2.系统内置 2.用户自定义
	Sort                int              `gorm:"column:sort;type:int(11);default:0;NOT NULL"`              // 排序
	OpeningStatement    string           `gorm:"column:opening_statement;type:varchar(500)"`               // 开场白
	ModelName           string           `gorm:"column:model_name;type:varchar(100);NOT NULL"`             // 模型名称
	ModelType           ChatBotModelType `gorm:"column:model_type;type:tinyint(4);default:1;NOT NULL"`     // 模型类型 1.文本 2.图片
	Assistants          []Assistants     `gorm:"many2many:chat_bot_assistant_associations;foreignKey:id;joinForeignKey:chat_bot_id;References:id;joinReferences:assistant_id"`
}

func (m *ChatBot) TableName() string {
	return "chat_bot"
}

type ChatBotAssistantAssociations struct {
	gorm.Model
	ChatBotId   uint `gorm:"column:chat_bot_id;type:bigint(20);NOT NULL"`  // 渠道表主键ID channels.id
	AssistantId uint `gorm:"column:assistant_id;type:bigint(20);NOT NULL"` // 模型表主键ID models.id
}

func (m *ChatBotAssistantAssociations) TableName() string {
	return "chat_bot_assistant_associations"
}
