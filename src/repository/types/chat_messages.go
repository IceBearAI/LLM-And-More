package types

import "gorm.io/gorm"

// ChatMessages 消息表
type ChatMessages struct {
	gorm.Model
	ModelName        string  `gorm:"column:model;size:64;notnull;index;default:'default';comment:聊天模型" json:"model"`
	ChannelId        uint    `gorm:"column:channel_id;null;index;comment:渠道ID" json:"channel_id"`
	Response         string  `gorm:"column:response;type:text;size:32768;null;comment:回复" json:"response"`
	Prompt           string  `gorm:"column:prompt;type:longtext;not null;comment:问题" json:"prompt"`
	PromptTokens     int     `gorm:"column:prompt_tokens;default:0;null;comment:问题Tokens" json:"prompt_tokens"`
	ResponseTokens   int     `gorm:"column:response_tokens;default:0;null;comment:回复Tokens" json:"response_tokens"`
	Finished         bool    `gorm:"column:finished;default:false;null;comment:是否完成" json:"finished"`
	TimeCost         string  `gorm:"column:time_cost;size:32;null;comment:耗时" json:"time_cost"`
	Temperature      float32 `orm:"column:temperature;default:0.9;null;comment:温度" json:"temperature"`
	TopP             float32 `orm:"column:top_p;default:0.9;null;comment:核心采样" json:"top_p"`
	N                int     `orm:"column:n;default:1;null;comment:聊天完成选项" json:"n"`
	User             string  `orm:"column:user;size:64;null;comment:用户" json:"user"`
	MessageId        string  `orm:"column:message_id;size:128;null;comment:消息ID" json:"message_id"`
	Object           string  `orm:"column:object;size:128;null;comment:对象" json:"object"`
	Created          int64   `gorm:"column:created;null;comment:创建时间" json:"created"`
	Messages         string  `gorm:"column:messages;type:longtext;null;comment:消息" json:"messages"`
	Error            bool    `gorm:"column:error;default:false;null;comment:是否错误" json:"error"`
	ErrorMessage     string  `gorm:"column:error_message;size:128;null;comment:错误消息" json:"error_message"`
	PresencePenalty  float32 `gorm:"column:presence_penalty;default:0.0;null;comment:存在惩罚" json:"presence_penalty"`
	FrequencyPenalty float32 `gorm:"column:frequency_penalty;default:0.0;null;comment:频率惩罚" json:"frequency_penalty"`
	MaxTokens        int     `gorm:"column:max_tokens;default:2048;null;comment:最大令牌" json:"max_tokens"`
	Stop             string  `gorm:"column:stop;size:128;null;comment:停止" json:"stop"`
}

func (*ChatMessages) TableName() string {
	return "chat_messages"
}
