package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"strings"
	"time"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) GetChatBotByAssistantId(ctx context.Context, assistantId uint) (res types.ChatBot, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "GetChatBotByAssistantId", "assistantId", assistantId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.GetChatBotByAssistantId(ctx, assistantId)
}

func (s *logging) CreateBot(ctx context.Context, data *types.ChatBot) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId), "method", "CreateBot",
			"name", data.Name,
			"email", data.Email,
			"assistants", data.Assistants,
			"descriptionForHuman", data.DescriptionForHuman,
			"descriptionForModel", data.DescriptionForModel,
			"privateStatus", data.PrivateStatus,
			"iconUrl", data.IconUrl,
			"botType", data.BotType,
			"sort", data.Sort,
			"openingStatement", data.OpeningStatement,
			"modelName", data.ModelName,
			"modelType", data.ModelType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateBot(ctx, data)
}

func (s *logging) CreateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateAudio", "channelId", data.ChannelId, "filename", data.FileName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateAudio(ctx, data)
}

func (s *logging) UpdateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateAudio", "channelId", data.ChannelId, "filename", data.FileName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateAudio(ctx, data)
}

func (s *logging) UpdateConversationModel(ctx context.Context, conversationIdInt uint, modelName types.ChatModel, maxToken int) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateConversationModel", "conversationIdInt", conversationIdInt, "modelName", modelName.String(), "maxToken", maxToken,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateConversationModel(ctx, conversationIdInt, modelName, maxToken)
}

func (s *logging) FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindChannelById", "id", id, "preload", strings.Join(preload, ","),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindChannelById(ctx, id, preload...)
}

func (s *logging) UpdateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateMessage",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateMessage(ctx, data)
}

func (s *logging) FindModelByChannelId(ctx context.Context, channelId uint, modelId string) (res types.ChatChannelModels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindModelByChannelId",
			"channelId", channelId,
			"modelId", modelId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindModelByChannelId(ctx, channelId, modelId)
}

func (s *logging) CreateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CreateMessage",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CreateMessage(ctx, data)
}

func (s *logging) FindChannelByApiKey(ctx context.Context, apiKey string, preload ...string) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindChannelByApiKey", "apiKey", apiKey,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindChannelByApiKey(ctx, apiKey, preload...)
}

func (s *logging) FindPrompts(ctx context.Context, promptType string) (res []types.ChatPrompts, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindPrompts", "promptType", promptType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindPrompts(ctx, promptType)
}

func (s *logging) CountChannel(ctx context.Context, channelName string, currTime time.Time) (res int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CountChannel", "channelName", channelName, "currTime", currTime,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.CountChannel(ctx, channelName, currTime)
}

func (s *logging) FindChannel(ctx context.Context, name string, preloads ...string) (res types.ChatChannels, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindChannel", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindChannel(ctx, name, preloads...)
}

func (s *logging) FindPromptTypes(ctx context.Context) (res []types.ChatPromptTypes, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindPromptTypes",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindPromptTypes(ctx)
}

func (s *logging) FindSystemPrompt(ctx context.Context, chatModel types.ChatModel, promptType types.ChatPromptType) (res types.ChatSystemPrompt, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindSystemPrompt", "chatModel", chatModel, "promptType", promptType,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindSystemPrompt(ctx, chatModel, promptType)
}

func (s *logging) FindOrCreateConversation(ctx context.Context, channelName string, data *types.ChatConversation) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindOrCreateConversation", "channelName", channelName, "email", data.Email, "conversationId", data.Uuid, "alias", data.Alias, "chatModel", data.ChatModel.String(),
			"sysPrompt", data.SysPrompt, "model", data.Model, "maxTokens", data.MaxTokens, "temperature", data.Temperature, "topP", data.TopP,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindOrCreateConversation(ctx, channelName, data)
}

func (s *logging) FindConversationByUuid(ctx context.Context, email, uuid string) (res types.ChatConversation, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindConversationByUuid", "email", email, "uuid", uuid,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindConversationByUuid(ctx, email, uuid)
}

func (s *logging) FindConversations(ctx context.Context, email string, page, pageSize int) (res []types.ChatConversation, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindConversations", "email", email, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindConversations(ctx, email, page, pageSize)
}

func (s *logging) UpdateConversation(ctx context.Context, email, id string, updates map[string]interface{}) error {
	values := make([]interface{}, 0, len(updates))
	values = append(values, "method", "UpdateConversation", "email", email, "id", id)
	for k, v := range updates {
		values = append(values, k, v)
	}
	var err error
	defer func(begin time.Time) {
		additionalValues := []interface{}{s.traceId, ctx.Value(s.traceId), "took", time.Since(begin), "err", err}
		values = append(values, additionalValues...)
		_ = s.logger.Log(values...)
	}(time.Now())
	err = s.next.UpdateConversation(ctx, email, id, updates)
	return err
}

func (s *logging) DeleteConversation(ctx context.Context, email string, conversationId uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "DeleteConversation", "email", email, "conversationId", conversationId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.DeleteConversation(ctx, email, conversationId)
}

func (s *logging) FindByChatId(ctx context.Context, email, chatId string) (res types.Chat, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByChatId", "email", email, "chatId", chatId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByChatId(ctx, email, chatId)
}

func (s *logging) ClearHistory(ctx context.Context, email string, roleId uint) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ClearHistory", "email", email, "roleId", roleId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ClearHistory(ctx, email, roleId)
}

func (s *logging) UpdateChat(ctx context.Context, data *types.Chat) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UpdateChat",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UpdateChat(ctx, data)
}

func (s *logging) Create(ctx context.Context, data *types.Chat) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, data)
}

func (s *logging) AllowUsers(ctx context.Context, email string) (res []types.ChatAllowUser, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AllowUsers", "email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.AllowUsers(ctx, email)
}

func (s *logging) History(ctx context.Context, model types.ChatModel, role uint, email string, promptType types.ChatPromptType, page, pageSize int) (res []types.Chat, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "History", "model", model, "role", role, "email", email, "promptType", promptType, "page", page, "pageSize", pageSize, "total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.History(ctx, model, role, email, promptType, page, pageSize)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.chat", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
