package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	// Create 创建记录
	Create(ctx context.Context, data *types.Chat) (err error)
	// UpdateChat 更新聊天记录
	UpdateChat(ctx context.Context, data *types.Chat) (err error)
	// AllowUsers 获取允许聊天的所有用户
	AllowUsers(ctx context.Context, email string) (res []types.ChatAllowUser, err error)
	// History 获取聊天记录
	History(ctx context.Context, model types.ChatModel, id uint, email string, promptType types.ChatPromptType, page, pageSize int) (res []types.Chat, total int64, err error)
	// FindRoleByName 获取角色
	FindRoleByName(ctx context.Context, name, email string) (res types.ChatRole, err error)
	// FindOrCreateRole 获取或创建角色
	FindOrCreateRole(ctx context.Context, name, email string) (res types.ChatRole, err error)
	// ClearHistory 清空历史记录
	ClearHistory(ctx context.Context, email string, id uint) (err error)
	// FindByChatId 根据chatId查找聊天记录
	FindByChatId(ctx context.Context, email, chatId string) (res types.Chat, err error)
	// FindOrCreateConversation 获取或创建会话
	FindOrCreateConversation(ctx context.Context, channelName string, data *types.ChatConversation) (err error)
	// FindConversationByUuid 根据uuid查找会话
	FindConversationByUuid(ctx context.Context, email, uuid string) (res types.ChatConversation, err error)
	// FindConversations 获取会话列表
	FindConversations(ctx context.Context, email string, page, pageSize int) (res []types.ChatConversation, err error)
	// UpdateConversation 更新会话信息
	UpdateConversation(ctx context.Context, email, id string, updates map[string]interface{}) (err error)
	// UpdateConversationModel 更新会话模型
	UpdateConversationModel(ctx context.Context, conversationIdInt uint, modelName types.ChatModel, maxToken int) (err error)
	// DeleteConversation 删除会话
	DeleteConversation(ctx context.Context, email string, conversationId uint) (err error)
	// FindSystemPrompt 获取系统提示
	FindSystemPrompt(ctx context.Context, chatModel types.ChatModel, promptType types.ChatPromptType) (res types.ChatSystemPrompt, err error)
	// FindPromptTypes 创建系统提示
	FindPromptTypes(ctx context.Context) (res []types.ChatPromptTypes, err error)
	// FindChannel 获取渠道
	FindChannel(ctx context.Context, name string, preload ...string) (res types.ChatChannels, err error)
	// CountChannel 统计渠道对话数
	CountChannel(ctx context.Context, channelName string, currTime time.Time) (res int, err error)
	// FindPrompts 获取所有提词
	FindPrompts(ctx context.Context, promptType string) (res []types.ChatPrompts, err error)
	// FindChannelByApiKey 根据ApiKey获取渠道
	FindChannelByApiKey(ctx context.Context, apiKey string, preload ...string) (res types.ChatChannels, err error)
	// FindChannelById 根据id获取渠道
	FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error)
	// CreateMessage 创建Message
	CreateMessage(ctx context.Context, data *types.ChatMessages) (err error)
	// UpdateMessage 更新Message
	UpdateMessage(ctx context.Context, data *types.ChatMessages) (err error)
	// FindModelByChannelId 根据chatId获取Message
	FindModelByChannelId(ctx context.Context, channelId uint, modelId string) (res types.ChatChannelModels, err error)

	// CreateRole 创建角色
	// Deprecated: 废弃
	CreateRole(ctx context.Context, name, alias, email string) (err error)
	// RolesByUser 获取用户角色
	// Deprecated: 废弃
	RolesByUser(ctx context.Context, email string, limit int) (res []types.ChatRole, err error)
	// UpdateRole 更新聊天角色信息
	// Deprecated: 废弃
	UpdateRole(ctx context.Context, email, roleName, roleAlias string) (err error)
	// DeleteRole 删除角色
	// Deprecated: 废弃
	DeleteRole(ctx context.Context, email, roleName string) (err error)
	// CreateAudio 创建音频记录
	CreateAudio(ctx context.Context, data *types.ChatAudio) (err error)
	// UpdateAudio 更新音频记录
	UpdateAudio(ctx context.Context, data *types.ChatAudio) (err error)
	// CreateBot 创建智能体
	CreateBot(ctx context.Context, data *types.ChatBot) (err error)
	// GetChatBotByAssistantId 根据assistantId获取智能体
	GetChatBotByAssistantId(ctx context.Context, assistantId uint) (res types.ChatBot, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) GetChatBotByAssistantId(ctx context.Context, assistantId uint) (res types.ChatBot, err error) {
	db := s.db.WithContext(ctx)
	err = db.Joins("JOIN chat_bot_assistant_associations ON chat_bot.id = chat_bot_assistant_associations.chat_bot_id").
		Where("chat_bot_assistant_associations.assistant_id = ?", assistantId).
		First(&res).Error
	return
}

func (s *service) CreateBot(ctx context.Context, data *types.ChatBot) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) UpdateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) CreateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	err = s.db.WithContext(ctx).Model(&types.ChatAudio{}).Create(data).Error
	return
}

func (s *service) UpdateConversationModel(ctx context.Context, conversationIdInt uint, modelName types.ChatModel, maxToken int) (err error) {
	err = s.db.WithContext(ctx).Model(&types.ChatConversation{}).Where("id = ?", conversationIdInt).Updates(types.ChatConversation{
		ChatModel: modelName,
		MaxTokens: maxToken,
	}).Error
	return
}

func (s *service) FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v)
	}
	err = db.Where("id = ?", id).First(&res).Error
	return
}

func (s *service) UpdateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) FindModelByChannelId(ctx context.Context, channelId uint, modelId string) (res types.ChatChannelModels, err error) {
	err = s.db.WithContext(ctx).Where("channel_id = ? AND model = ?", channelId, modelId).First(&res).Error
	return
}

func (s *service) CreateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) FindChannelByApiKey(ctx context.Context, apiKey string, preload ...string) (res types.ChatChannels, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v)
	}
	err = db.Where("api_key = ?", apiKey).First(&res).Error
	return
}

func (s *service) FindPrompts(ctx context.Context, promptType string) (res []types.ChatPrompts, err error) {
	err = s.db.WithContext(ctx).Where("prompt_type = ?", promptType).Order("id DESC").Find(&res).Error
	return
}

func (s *service) CountChannel(ctx context.Context, channelName string, currTime time.Time) (res int, err error) {
	//err = s.db.WithContext(ctx).Where("name = ? AND created_at between ? AND ?",
	//		channelName, fmt.Sprintf("%s 00:00:00", currTime.Format("2006-01-02")),
	//		fmt.Sprintf("%s 23:59:59", currTime.Format("2006-01-02"))).Count(&res).Error
	return
}

func (s *service) FindChannel(ctx context.Context, name string, preload ...string) (res types.ChatChannels, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v, func(db *gorm.DB) *gorm.DB {
			return db.Order("chat_channel_models.updated_at DESC")
		})
	}
	err = db.Where("name = ?", name).First(&res).Error
	return
}

func (s *service) FindPromptTypes(ctx context.Context) (res []types.ChatPromptTypes, err error) {
	err = s.db.WithContext(ctx).Order("id DESC").Find(&res).Error
	return
}

func (s *service) FindSystemPrompt(ctx context.Context, chatModel types.ChatModel, promptType types.ChatPromptType) (res types.ChatSystemPrompt, err error) {
	err = s.db.WithContext(ctx).Where("chat_model = ? AND prompt_type = ?", chatModel, promptType).First(&res).Error
	return
}

func (s *service) DeleteConversation(ctx context.Context, email string, id uint) (err error) {
	return s.db.WithContext(ctx).Where("email = ? AND id = ?", email, id).Delete(&types.ChatConversation{}).Error
}

func (s *service) UpdateConversation(ctx context.Context, email, uuid string, updates map[string]interface{}) (err error) {
	return s.db.WithContext(ctx).Model(&types.ChatConversation{}).Where("email = ? AND uuid = ?", email, uuid).Updates(updates).Error
}

func (s *service) FindConversations(ctx context.Context, email string, page, pageSize int) (res []types.ChatConversation, err error) {
	err = s.db.WithContext(ctx).Where("email = ?", email).Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error
	return
}

func (s *service) FindConversationByUuid(ctx context.Context, email, uuid string) (res types.ChatConversation, err error) {
	err = s.db.WithContext(ctx).Where("email = ? AND uuid = ?", email, uuid).First(&res).Error
	return
}

func (s *service) FindOrCreateConversation(ctx context.Context, channelName string, data *types.ChatConversation) (err error) {
	if strings.EqualFold(data.Uuid, "") {
		data.Uuid = uuid.New().String()
	}
	if len(data.Alias) > 32 {
		data.Alias = string([]rune(data.Alias)[0:10])
	}

	if channelName != "" {
		channel, err := s.FindChannel(ctx, channelName)
		if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
			return err
		}
		data.ChannelId = channel.ID
	}
	return s.db.WithContext(ctx).Where("email = ? AND uuid = ?", data.Email, data.Uuid).FirstOrCreate(&data).Error
}

func (s *service) FindByChatId(ctx context.Context, email, chatId string) (res types.Chat, err error) {
	err = s.db.WithContext(ctx).Where("email = ? AND chat_id = ?", email, chatId).First(&res).Error
	return
}

func (s *service) ClearHistory(ctx context.Context, email string, conversationId uint) (err error) {
	return s.db.WithContext(ctx).Where("email = ? AND conversation_id = ?", email, conversationId).Delete(&types.Chat{}).Error
}

func (s *service) UpdateChat(ctx context.Context, data *types.Chat) (err error) {
	return s.db.WithContext(ctx).Save(data).Error
}

func (s *service) DeleteRole(ctx context.Context, email, roleName string) (err error) {
	return s.db.WithContext(ctx).Where("email = ? AND name = ?", email, roleName).Delete(&types.ChatRole{}).Error
}

func (s *service) UpdateRole(ctx context.Context, email, roleName, roleAlias string) (err error) {
	roleInfo, err := s.FindRoleByName(ctx, roleName, email)
	if err != nil {
		return
	}
	roleInfo.Alias = roleAlias
	return s.db.WithContext(ctx).Save(&roleInfo).Error
}

func (s *service) RolesByUser(ctx context.Context, email string, limit int) (res []types.ChatRole, err error) {
	err = s.db.WithContext(ctx).Where("email = ?", email).Order("id DESC").Limit(limit).Find(&res).Error
	return
}

func (s *service) CreateRole(ctx context.Context, name, alias, email string) (err error) {
	return s.db.WithContext(ctx).Create(&types.ChatRole{
		Name:  name,
		Alias: alias,
		Email: email,
	}).Error
}

func (s *service) FindOrCreateRole(ctx context.Context, name, email string) (res types.ChatRole, err error) {
	res.Name = name
	res.Email = email
	err = s.db.WithContext(ctx).Where("name = ? AND email = ?", name, email).FirstOrCreate(&res).Error
	return
}

func (s *service) FindRoleByName(ctx context.Context, name, email string) (res types.ChatRole, err error) {
	err = s.db.WithContext(ctx).Where("name = ? AND email = ?", name, email).First(&res).Error
	return
}

func (s *service) History(ctx context.Context, model types.ChatModel, conversationId uint, email string, promptType types.ChatPromptType, page, pageSize int) (res []types.Chat, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Chat{}) /*.Where("chat_model = ?", model)*/
	if email != "" {
		query = query.Where("email = ?", email)
	}
	if conversationId > 0 {
		query = query.Where("conversation_id = ?", conversationId)
	}
	if promptType != "" {
		query = query.Where("prompt_type = ?", promptType)
	}
	err = query.Count(&total).Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&res).Error
	return
}

func (s *service) AllowUsers(ctx context.Context, email string) (res []types.ChatAllowUser, err error) {
	query := s.db.WithContext(ctx)
	if email != "" {
		query = query.Where("email = ?", email)
	}
	err = query.Find(&res).Error
	return
}

func (s *service) Create(ctx context.Context, data *types.Chat) (err error) {
	return s.db.WithContext(ctx).Create(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
