package types

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type ChatModel string
type ChatPromptType string
type ChatSendType string
type FileExt string
type AudioType string

func (c ChatPromptType) String() string {
	return string(c)
}
func (c ChatModel) String() string {
	return string(c)
}

func (c ChatSendType) String() string {
	return string(c)
}

func (c FileExt) String() string {
	return string(c)
}

func (c AudioType) String() string {
	return string(c)
}

const (
	ChatPromptTypeImagine       ChatPromptType = "imagine"
	ChatPromptTypeText          ChatPromptType = "text"
	ChatPromptTypeQA            ChatPromptType = "qa"
	ChatSendTypeText            ChatSendType   = "text"
	ChatSendTypeAudio           ChatSendType   = "audio"
	ChatSendTypeImagine         ChatSendType   = "imagine"
	ChatSendTypeFile            ChatSendType   = "file"
	FileExtMp3                  FileExt        = "mp3"
	FileExtPng                  FileExt        = "png"
	FileExtPdf                  FileExt        = "pdf"
	AudioTypeTextToAudio        AudioType      = "text_to_audio"
	AudioTypeAudioTranslation   AudioType      = "audio_translation"
	AudioTypeAudioTranscription AudioType      = "audio_transcription"
)

// Chat 对话消息表
type Chat struct {
	gorm.Model
	ChatModel      ChatModel      `gorm:"column:chat_model;size:64;null;index;default:'default'" json:"chat_model"`
	Email          string         `gorm:"column:email;size:128;not null;index;" json:"email"`
	Prompt         string         `gorm:"column:prompt;type:text;size:8192;not null;" json:"prompt"`
	Response       string         `gorm:"column:response;type:longtext;null;" json:"response"`
	BeginTime      time.Time      `gorm:"column:begin_time;null;" json:"begin_time"`
	EndTime        time.Time      `gorm:"column:end_time;null;" json:"end_time"`
	Temperature    float64        `gorm:"column:temperature;null;" json:"temperature"`
	TopP           float64        `gorm:"column:top_p;null;" json:"top_p"`
	Status         int            `gorm:"column:status;null;" json:"status"`
	RoleId         uint           `gorm:"column:role_id;null;index;" json:"role_id"`
	TimeCost       string         `gorm:"column:time_cost;size:16;null;" json:"time_cost"`
	Error          bool           `gorm:"column:error;null;default:false;" json:"error"`
	ChatId         string         `gorm:"column:chat_id;size:64;null;index;" json:"chat_id"`
	MaxLength      int            `gorm:"column:max_length;null;" json:"max_length"`
	ConversationId uint           `gorm:"column:conversation_id;null;index;" json:"conversation_id"`
	PromptType     ChatPromptType `gorm:"column:prompt_type;size:32;null;default:'text';" json:"prompt_type"`
	PanUrl         string         `gorm:"column:pan_url;size:50;null;" json:"pan_url"`
	SendType       ChatSendType   `gorm:"column:send_type;size:10;null;default:'text';" json:"send_type"`
	Ext            string         `gorm:"column:ext;size:128;null;" json:"ext"`

	Conversation ChatConversation `gorm:"foreignKey:ConversationId;references:ID" json:"-"`
}

// ChatConversation 对话表
type ChatConversation struct {
	gorm.Model
	Uuid        string    `gorm:"column:uuid;size:40;not null;unique;" json:"uuid"`                          // comment:UUID
	Alias       string    `gorm:"column:alias;size:64;null;default:'New Chat';" json:"alias"`                // comment:别名
	Email       string    `gorm:"column:email;size:128;not null;index;" json:"email"`                        // comment:邮箱
	ChatModel   ChatModel `gorm:"column:chat_model;size:64;null;index;default:'default';" json:"chat_model"` // comment:聊天模型
	ChannelId   uint      `gorm:"column:channel_id;null;index;" json:"channel_id"`                           // comment:渠道ID
	Temperature float64   `gorm:"column:temperature;null;" json:"temperature"`                               // comment:温度
	TopP        float64   `gorm:"column:top_p;null;" json:"top_p"`                                           // comment:采样性
	MaxTokens   int       `gorm:"column:max_tokens;default:2048;null;" json:"max_tokens"`                    // comment:最大支持Tokens
	SysPrompt   string    `gorm:"column:sys_prompt;size:255;null;" json:"sys_prompt"`                        // comment:系统提示

	Channel ChatChannels `gorm:"foreignKey:channel_id;references:id" json:"-"`
}

// ChatSystemPrompt 系统提示表
type ChatSystemPrompt struct {
	gorm.Model
	ChatModel  ChatModel      `gorm:"column:chat_model;size:64;notnull;index;default:'default';" json:"chat_model"`
	Content    string         `gorm:"column:content;type:text;size:8192;not null;" json:"content"`
	PromptType ChatPromptType `gorm:"column:prompt_type;size:32;notnull;index;default:'text';" json:"prompt_type"`
}

// ChatPromptTypes 提示类型表
type ChatPromptTypes struct {
	gorm.Model
	Name   string `gorm:"column:name;size:32;not null;index;" json:"name"`
	Alias  string `gorm:"column:alias;size:64;null;default:'New Chat';" json:"alias"`
	Remark string `gorm:"column:remark;size:128;null;" json:"remark"`
}

// ChatPrompts 提示表
type ChatPrompts struct {
	gorm.Model
	Title      string         `gorm:"column:title;size:64;not null;index;" json:"title"`                           // comment:标题
	Content    string         `gorm:"column:content;type:text;size:8192;not null;" json:"content"`                 // comment:内容
	PromptType ChatPromptType `gorm:"column:prompt_type;size:32;notnull;index;default:'text';" json:"prompt_type"` // comment:提示类型
}

// ChatChannelModels 渠道模型表
type ChatChannelModels struct {
	gorm.Model
	ChannelId    uint         `gorm:"column:channel_id;null;index;" json:"channel_id"`                    // comment:渠道ID
	ModelName    ChatModel    `gorm:"column:model;size:64;notnull;index;default:'default';" json:"model"` // comment:聊天模型
	MaxTokens    int          `gorm:"column:max_tokens;default:2048;null;" json:"max_tokens"`             // comment:最大支持Tokens
	IsPrivate    bool         `gorm:"column:is_private;null;default:false;" json:"is_private"`            // comment:是否为本地私有模型
	ChatChannels ChatChannels `gorm:"foreignKey:id;references:channel_id" json:"-"`
}

// ChatMessages 消息表
type ChatMessages struct {
	gorm.Model
	ModelName      ChatModel `gorm:"column:model;size:64;notnull;index;default:'default';" json:"model"` // comment:聊天模型
	ChannelId      uint      `gorm:"column:channel_id;null;index;" json:"channel_id"`                    // comment:渠道ID
	Response       string    `gorm:"column:response;type:longtext;size:65536;null;" json:"response"`     // comment:回复
	Prompt         string    `gorm:"column:prompt;type:text;size:32768;not null;" json:"prompt"`         // comment:问题
	PromptTokens   int       `gorm:"column:prompt_tokens;default:0;null;" json:"prompt_tokens"`          // comment:问题Tokens
	ResponseTokens int       `gorm:"column:response_tokens;default:0;null;" json:"response_tokens"`      // comment:回复Tokens
	Finished       bool      `gorm:"column:finished;default:false;null;" json:"finished"`                // comment:是否完成
	TimeCost       string    `gorm:"column:time_cost;size:32;null;" json:"time_cost"`                    // comment:耗时
	Temperature    float64   `orm:"column:temperature;default:0.9;null;" json:"temperature"`             // comment:温度
	TopP           float64   `orm:"column:top_p;default:0.9;null;" json:"top_p"`                         // comment:核心采样
	N              int       `orm:"column:n;default:1;null;" json:"n"`                                   // comment:聊天完成选项
	User           string    `orm:"column:user;size:64;null;" json:"user"`                               // comment:用户
	MessageId      string    `orm:"column:message_id;size:128;null;" json:"message_id"`                  // comment:消息ID
	Object         string    `orm:"column:object;size:128;null;" json:"object"`                          // comment:对象
	Created        int64     `gorm:"column:created;null;" json:"created"`                                // comment:创建时间
	Messages       string    `gorm:"column:messages;type:text;size:65536;null;" json:"messages"`         // comment:消息
}

// 正向标签 masterpiece, best quality, top quality, ultra highres, 8k hdr, 8k wallpaper, RAW, huge file size, intricate details, sharp focus, natural lighting, realistic, professional, delicate, amazing, CG, finely detailed, beautiful detailed, colourful
// 反向标签 paintings, sketches, lowres, normal quality, worst quality, low quality, cropped, dot, mole, ugly, grayscale, monochrome, duplicate, morbid, mutilated, missing fingers, extra fingers, too many fingers, fused fingers, mutated hands, bad hands, poorly drawn hands, poorly drawn face, poorly drawn eyebrows, bad anatomy, cloned face, long neck, extra legs, extra arms, missing arms missing legs, malformed limbs, deformed, simple background, bad proportions, disfigured, skin spots, skin blemishes, age spot, bad feet, error, text, extra digit, fewer digits, jpeg artifacts, signature, username, blurry, watermark, mask, logo

type ChatRole struct {
	gorm.Model
	Name  string `gorm:"column:name;size:32;not null;index;" json:"name"`            // comment:角色名称
	Alias string `gorm:"column:alias;size:32;null;default:'New Chat';" json:"alias"` // comment:角色别名
	Email string `gorm:"column:email;size:128;not null;index;" json:"email"`         // comment:邮箱
}

type ChatAllowUser struct {
	gorm.Model
	Email string `gorm:"column:email;size:128;not null;index;" json:"email"` // comment:邮箱
}

func (c *Chat) TableName() string {
	return "chat"
}

func (*ChatRole) TableName() string {
	return "chat_role"
}

func (*ChatConversation) TableName() string {
	return "chat_conversation"
}

func (*ChatAllowUser) TableName() string {
	return "chat_allow_user"
}

func (*ChatSystemPrompt) TableName() string {
	return "chat_system_prompt"
}

func (*ChatPromptTypes) TableName() string {
	return "chat_prompt_types"
}

func (*ChatPrompts) TableName() string {
	return "chat_prompts"
}

func (*ChatChannelModels) TableName() string {
	return "chat_channel_models"
}

// RemoveGPT4Models 移除含有GPT4的模型
func RemoveGPT4Models(models []ChatChannelModels) []ChatChannelModels {
	if len(models) == 0 {
		return models
	}
	var newModels []ChatChannelModels
	for _, model := range models {
		if !strings.Contains(strings.ToLower(model.ModelName.String()), "gpt-4") {
			newModels = append(newModels, model)
		}
	}
	return newModels
}

// ChatAudio 音频文件
type ChatAudio struct {
	gorm.Model
	ChannelId     uint      `json:"channel_id" gorm:"channel_id"`         // 渠道ID
	FileName      string    `json:"file_name" gorm:"file_name"`           // 上传文件名
	TargetPath    string    `json:"target_path" gorm:"target_path"`       // 网盘路径
	PanUrl        string    `json:"pan_url" gorm:"pan_url"`               // 生成网盘地址
	AudioText     string    `json:"audio_text" gorm:"audio_text"`         // 音频对应的文本
	AudioDuration float64   `json:"audio_duration"`                       // 音频时长
	AudioType     AudioType `json:"audio_type"`                           // 音频处理类型  text_to_audio audio_translation audio_transcription
	TranslateText string    `json:"translate_text" gorm:"translate_text"` // 音频文本翻译
}

// TableName 表名称
func (*ChatAudio) TableName() string {
	return "chat_audio"
}
