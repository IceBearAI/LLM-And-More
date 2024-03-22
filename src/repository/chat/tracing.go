package chat

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"time"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) GetChatBotByAssistantId(ctx context.Context, assistantId uint) (res types.ChatBot, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetChatBotByAssistantId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV(
			"assistantId", assistantId,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetChatBotByAssistantId(ctx, assistantId)
}

func (s *tracing) CreateBot(ctx context.Context, data *types.ChatBot) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateBot", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV(
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
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateBot(ctx, data)
}

func (s *tracing) CreateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateAudio", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateAudio(ctx, data)
}

func (s *tracing) UpdateAudio(ctx context.Context, data *types.ChatAudio) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateAudio", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateAudio(ctx, data)
}

func (s *tracing) UpdateConversationModel(ctx context.Context, conversationIdInt uint, modelName types.ChatModel, maxToken int) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateConversationModel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("conversationIdInt", conversationIdInt, "modelName", modelName, "maxToken", maxToken, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateConversationModel(ctx, conversationIdInt, modelName, maxToken)
}

func (s *tracing) FindChannelById(ctx context.Context, id uint, preload ...string) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindChannelById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("id", id, "preload", preload, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindChannelById(ctx, id, preload...)
}

func (s *tracing) UpdateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateMessage", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateMessage(ctx, data)
}

func (s *tracing) FindModelByChannelId(ctx context.Context, channelId uint, modelId string) (res types.ChatChannelModels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindModelByChannelId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("channelId", channelId, "modelId", modelId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindModelByChannelId(ctx, channelId, modelId)
}

func (s *tracing) CreateMessage(ctx context.Context, data *types.ChatMessages) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateMessage", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateMessage(ctx, data)
}

func (s *tracing) FindChannelByApiKey(ctx context.Context, apiKey string, preload ...string) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindChannelByApiKey", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("apiKey", apiKey, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindChannelByApiKey(ctx, apiKey, preload...)
}

func (s *tracing) FindPrompts(ctx context.Context, promptType string) (res []types.ChatPrompts, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindPrompts", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("promptType", promptType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindPrompts(ctx, promptType)
}

func (s *tracing) CountChannel(ctx context.Context, channelName string, currTime time.Time) (res int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CountChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("channelName", channelName, "currTime", currTime, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CountChannel(ctx, channelName, currTime)
}

func (s *tracing) FindChannel(ctx context.Context, name string, preloads ...string) (res types.ChatChannels, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindChannel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("name", name, "preloads", preloads, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindChannel(ctx, name, preloads...)
}

func (s *tracing) FindPromptTypes(ctx context.Context) (res []types.ChatPromptTypes, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindPromptTypes", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindPromptTypes(ctx)
}

func (s *tracing) FindSystemPrompt(ctx context.Context, chatModel types.ChatModel, promptType types.ChatPromptType) (res types.ChatSystemPrompt, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindSystemPrompt", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("chatModel", chatModel, "promptType", promptType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindSystemPrompt(ctx, chatModel, promptType)
}

func (s *tracing) FindOrCreateConversation(ctx context.Context, channelName string, data *types.ChatConversation) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindOrCreateConversation", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("channelName", channelName, "email", data.Email, "conversationId", data.Uuid, "alias", data.Alias, "chatModel", data.ChatModel,
			"maxTokens", data.MaxTokens, "temperature", data.Temperature, "topP", data.TopP,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindOrCreateConversation(ctx, channelName, data)
}

func (s *tracing) FindConversationByUuid(ctx context.Context, email, uuid string) (res types.ChatConversation, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindConversationByUuid", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindConversationByUuid(ctx, email, uuid)
}

func (s *tracing) FindConversations(ctx context.Context, email string, page, pageSize int) (res []types.ChatConversation, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindConversations", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindConversations(ctx, email, page, pageSize)
}

func (s *tracing) UpdateConversation(ctx context.Context, email, id string, updates map[string]interface{}) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateConversation", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	values := make([]interface{}, 0, len(updates))
	values = append(values, "method", "UpdateConversation", "email", email, "id", id)
	for k, v := range updates {
		values = append(values, k, v)
	}
	var err error
	defer func() {
		values = append(values, "err", err)
		span.LogKV(values...)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateConversation(ctx, email, id, updates)
}

func (s *tracing) DeleteConversation(ctx context.Context, email string, conversationId uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteConversation", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "conversationId", conversationId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteConversation(ctx, email, conversationId)
}

func (s *tracing) FindByChatId(ctx context.Context, email, chatId string) (res types.Chat, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByChatId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "chatId", chatId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByChatId(ctx, email, chatId)
}

func (s *tracing) ClearHistory(ctx context.Context, email string, roleId uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ClearHistory", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "roleId", roleId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ClearHistory(ctx, email, roleId)
}

func (s *tracing) UpdateChat(ctx context.Context, data *types.Chat) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateChat", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateChat(ctx, data)
}

func (s *tracing) DeleteRole(ctx context.Context, email, roleName string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DeleteRole", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "roleName", roleName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.DeleteRole(ctx, email, roleName)
}

func (s *tracing) Create(ctx context.Context, data *types.Chat) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, data)
}

func (s *tracing) AllowUsers(ctx context.Context, email string) (res []types.ChatAllowUser, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AllowUsers", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.AllowUsers(ctx, email)
}

func (s *tracing) History(ctx context.Context, model types.ChatModel, role uint, email string, promptType types.ChatPromptType, page, pageSize int) (res []types.Chat, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "History", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("model", model, "role", role, "email", email, "promptType", promptType, "page", page, "pageSize", pageSize, "total", total, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.History(ctx, model, role, email, promptType, page, pageSize)
}

func (s *tracing) FindRoleByName(ctx context.Context, name, email string) (res types.ChatRole, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindRoleByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("name", name, "email", email, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindRoleByName(ctx, name, email)
}

func (s *tracing) FindOrCreateRole(ctx context.Context, name, email string) (res types.ChatRole, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindOrCreateRole", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("name", name, "email", email, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindOrCreateRole(ctx, name, email)
}

func (s *tracing) CreateRole(ctx context.Context, name, alias, email string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateRole", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("name", name, "alias", alias, "email", email, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.CreateRole(ctx, name, alias, email)
}

func (s *tracing) RolesByUser(ctx context.Context, email string, limit int) (res []types.ChatRole, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RolesByUser", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "limit", limit, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RolesByUser(ctx, email, limit)
}

func (s *tracing) UpdateRole(ctx context.Context, email, roleName, roleAlias string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UpdateRole", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.chat",
	})
	defer func() {
		span.LogKV("email", email, "roleName", roleName, "roleAlias", roleAlias, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.UpdateRole(ctx, email, roleName, roleAlias)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
