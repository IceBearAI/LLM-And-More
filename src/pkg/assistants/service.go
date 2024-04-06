package assistants

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/middleware"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	lmtools "github.com/tmc/langchaingo/tools"
	"strings"
	"time"
)

type Middleware func(Service) Service

// Service describes a service that adds things together.
// 助手管理服务
type Service interface {
	// Create 创建助手
	Create(ctx context.Context, tenantId uint, req createRequest) (assistantId string, err error)
	// Update 更新助手
	Update(ctx context.Context, tenantId uint, assistantId string, req updateRequest) (err error)
	// Delete 删除助手
	Delete(ctx context.Context, tenantId uint, assistantId string) (err error)
	// Get 获取助手
	Get(ctx context.Context, tenantId uint, assistantId string) (resp assistantResult, err error)
	// List 列出助手
	List(ctx context.Context, tenantId uint, name string, page, pageSize int) (resp []assistantResult, total int64, err error)
	// AddTool 给助手添加工具
	AddTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error)
	// RemoveTool 给助手移除工具
	RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error)
	// ListTool 列出助手的工具
	ListTool(ctx context.Context, tenantId uint, assistantId, name string, page, pageSize int) (resp []assistantToolResult, total int64, err error)
	// Playground 操场测试对话
	Playground(ctx context.Context, tenantId uint, assistantId string, req playgroundRequest) (resp <-chan playgroundResult, err error)
	// Publish 发布助手
	Publish(ctx context.Context, tenantId uint, assistantId string) (err error)
}

type streamLogHandler struct {
	callbacks.SimpleHandler
	stream      chan playgroundResult
	fullContent string
	currTime    time.Time
}

func (s *streamLogHandler) HandleStreamingFunc(_ context.Context, chunk []byte) {
	s.fullContent += string(chunk)
	s.stream <- playgroundResult{
		Content:     string(chunk),
		FullContent: s.fullContent,
		CreatedAt:   s.currTime,
	}
}

func newStreamLogHandler(stream chan playgroundResult) callbacks.Handler {
	return &streamLogHandler{
		stream:   stream,
		currTime: time.Now(),
	}
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	llmOptions []openai.Option
	httpOpts   []kithttp.ClientOption
}

func (s *service) Publish(ctx context.Context, tenantId uint, assistantId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	assistant, err := s.repository.Assistants().Get(ctx, tenantId, assistantId)
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	res, err := s.repository.Chat().GetChatBotByAssistantId(ctx, assistant.ID)
	if err == nil && res.ID > 0 {
		err = errors.New("助手已发布")
		_ = level.Warn(logger).Log("msg", "助手已发布", "err", err)
		return
	}
	assistants := make([]types.Assistants, 0)
	assistants = append(assistants, assistant)
	err = s.repository.Chat().CreateBot(ctx, &types.ChatBot{
		Name:                assistant.Name,
		DescriptionForHuman: assistant.Description,
		DescriptionForModel: assistant.Instructions,
		PrivateStatus:       1,
		IconUrl:             assistant.Avatar,
		BotType:             types.ChatBotTypeUser,
		Sort:                0,
		OpeningStatement:    "",
		ModelName:           assistant.ModelName,
		ModelType:           types.ChatBotModelTypeText,
		Assistants:          assistants,
	})
	if err != nil {
		err = errors.Wrap(err, "发布助手失败")
		_ = level.Warn(logger).Log("msg", "发布助手失败", "err", err)
		return
	}
	return
}

func (s *service) Playground(ctx context.Context, tenantId uint, assistantId string, req playgroundRequest) (resp <-chan playgroundResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	if len(req.Messages) == 0 {
		err = errors.New("消息不能为空")
		_ = level.Warn(logger).Log("msg", "消息不能为空", "err", err)
		return
	}
	assistant, err := s.repository.Assistants().Get(ctx, tenantId, assistantId, "Tools")
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	_ = level.Info(logger).Log("assistant.name", assistant.Name)
	// 获取历史对话记录
	var previousMessages []schema.ChatMessage
	var lastMessage message
	// 取出最后一条消息
	lastMessage = req.Messages[len(req.Messages)-1:][0]
	// 弹出 req.Messages 的最后一条消息
	req.Messages = req.Messages[:len(req.Messages)-1]
	for _, msg := range req.Messages {
		if strings.EqualFold(msg.Role, "user") {
			previousMessages = append(previousMessages, schema.HumanChatMessage{Content: msg.Content})
		} else if strings.EqualFold(msg.Role, "assistant") {
			previousMessages = append(previousMessages, schema.AIChatMessage{Content: msg.Content})
		} else if strings.EqualFold(msg.Role, "system") {
			if strings.TrimSpace(msg.Content) != "" {
				previousMessages = append(previousMessages, schema.SystemChatMessage{Content: msg.Content})
			}
		}
	}

	assistantTools, err := s.repository.Tools().GetByIds(ctx, tenantId, req.ToolIds)
	if err != nil {
		err = errors.Wrap(err, "获取工具失败")
		_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
		return
	}

	options := s.llmOptions
	options = append(options, openai.WithModel(req.ModelName))
	llm, err := openai.New(options...)
	if err != nil {
		err = errors.Wrap(err, "创建LLM失败")
		_ = level.Warn(logger).Log("msg", "创建LLM失败", "err", err)
		return
	}

	var tools []lmtools.Tool
	for _, tool := range assistantTools {
		tools = append(tools, NewDynamicTool(ToolOptions{
			Name:        tool.Name,
			Description: tool.Description,
			ActionID:    tool.UUID,
			ToolType:    tool.ToolType,
			Metadata:    tool.Metadata,
			logger:      logger,
			TraceId:     s.traceId,
			Opts:        s.httpOpts,
		}))
	}

	streamResp := make(chan playgroundResult)
	creationOptions := []agents.Option{
		agents.WithParserErrorHandler(agents.NewParserErrorHandler(func(s string) string {
			// 这里可以发告警出来
			_ = level.Warn(logger).Log("agents.WithParserErrorHandler", "agents.NewParserErrorHandler", "msg", "解析错误", "err", s)
			return s
		})),
		//agents.WithPromptFormatInstructions(assistant.Instructions),
		//agents.ZeroShotReactDescription,
	}
	if req.Stream {
		creationOptions = append(creationOptions, agents.WithCallbacksHandler(newStreamLogHandler(streamResp)))
	}
	executor := agents.NewExecutor(agents.NewConversationalAgent(llm, tools, creationOptions...), tools, []agents.Option{
		agents.WithMaxIterations(3),
		agents.WithReturnIntermediateSteps(),
		agents.WithMemory(
			memory.NewConversationBuffer(
				memory.WithChatHistory(
					memory.NewChatMessageHistory(
						memory.WithPreviousMessages(previousMessages),
					),
				),
			),
		),
	}...)
	var chainsCalls []chains.ChainCallOption
	chainsCalls = append(chainsCalls,
		chains.WithTemperature(0),
		chains.WithTopP(0),
		chains.WithMaxTokens(2048),
		//chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		//	_ = level.Info(logger).Log("msg", "chains.WithStreamingFunc", "chunk", string(chunk))
		//	return nil
		//}),
	)
	t := time.Now()
	//if req.Stream {
	go func() {
		result, err := chains.Run(ctx, executor, lastMessage.Content, chainsCalls...)
		if err != nil {
			err = errors.Wrap(err, "运行Executor失败")
			_ = level.Warn(logger).Log("msg", "运行Executor失败", "err", err)
			streamResp <- playgroundResult{
				FullContent:  err.Error(),
				FinishReason: "stop",
				Content:      err.Error(),
				CreatedAt:    t,
			}
			close(streamResp)
			return
		}
		_ = level.Info(logger).Log("msg", "运行Executor成功", "result", result)
		streamResp <- playgroundResult{
			FullContent:  result,
			FinishReason: "stop",
			CreatedAt:    t,
			Content:      result,
		}
		close(streamResp)
	}()
	return streamResp, nil
	//}
	//result, err := chains.Run(ctx, executor, lastMessage.Content, chainsCalls...)
	//if err != nil {
	//	err = errors.Wrap(err, "运行Executor失败")
	//	_ = level.Warn(logger).Log("msg", "运行Executor失败", "err", err)
	//	return nil, err
	//}
	//_ = level.Info(logger).Log("msg", "运行Executor成功", "result", result)
	//streamResp <- playgroundResult{
	//	FullContent:  result,
	//	FinishReason: "stop",
	//	CreatedAt:    time.Now(),
	//	Content:      result,
	//}
	//fmt.Println("streamResp", streamResp)
	//close(streamResp)
	//return streamResp, nil

}

func (s *service) Create(ctx context.Context, tenantId uint, req createRequest) (assistantId string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	assistant, err := s.repository.Assistants().FindByAssistantName(ctx, tenantId, req.Name)
	if err == nil {
		err = errors.New("助手已存在")
		_ = level.Warn(logger).Log("msg", "助手已存在", "assistant", assistant)
		return
	}
	assistant = types.Assistants{
		UUID:         fmt.Sprintf("assistant-%s", uuid.New().String()),
		TenantId:     tenantId,
		Name:         req.Name,
		Remark:       req.Remark,
		ModelName:    req.ModelName,
		Description:  req.Description,
		Instructions: req.Instructions,
		Metadata:     req.Metadata,
		Avatar:       req.Avatar,
		Operator:     email,
	}
	err = s.repository.Assistants().Create(ctx, &assistant)
	return
}

func (s *service) Update(ctx context.Context, tenantId uint, assistantId string, req updateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	assistant, err := s.repository.Assistants().Get(ctx, tenantId, assistantId)
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	//tools, err := s.repository.Tools().GetByIds(ctx, tenantId, req.ToolIds)
	//if err != nil {
	//	err = errors.Wrap(err, "获取工具失败")
	//	_ = level.Warn(logger).Log("msg", "获取工具失败", "err", err)
	//	return
	//}
	assistant.Name = req.Name
	assistant.Remark = req.Remark
	assistant.ModelName = req.ModelName
	assistant.Description = req.Description
	assistant.Instructions = req.Instructions
	assistant.Metadata = req.Metadata
	//assistant.Tools = tools
	if err = s.repository.Assistants().Update(ctx, &assistant); err != nil {
		err = errors.Wrap(err, "更新助手失败")
		_ = level.Warn(logger).Log("msg", "更新助手失败", "err", err)
		return
	}
	//err = s.repository.Assistants().ReplaceTools(ctx, &assistant, tools)
	//if err != nil {
	//	err = errors.Wrap(err, "更新助手工具失败")
	//	_ = level.Warn(logger).Log("msg", "更新助手工具失败", "err", err)
	//	return
	//}
	return
}

func (s *service) Delete(ctx context.Context, tenantId uint, assistantId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	err = s.repository.Assistants().Delete(ctx, tenantId, assistantId)
	if err != nil {
		err = errors.Wrap(err, "删除助手失败")
		_ = level.Warn(logger).Log("msg", "删除助手失败", "err", err)
		return
	}
	return
}

func (s *service) Get(ctx context.Context, tenantId uint, assistantId string) (resp assistantResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	assistant, err := s.repository.Assistants().Get(ctx, tenantId, assistantId, "Tools")
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	resp = assistantResult{
		AssistantId:  assistant.UUID,
		Name:         assistant.Name,
		Remark:       assistant.Remark,
		ModelName:    assistant.ModelName,
		Description:  assistant.Description,
		Instructions: assistant.Instructions,
		Metadata:     assistant.Metadata,
		CreatedAt:    assistant.CreatedAt,
		UpdatedAt:    assistant.UpdatedAt,
		Avatar:       assistant.Avatar,
		Operator:     assistant.Operator,
	}
	resp.Tools = make([]assistantToolResult, 0)
	for _, tool := range assistant.Tools {
		resp.Tools = append(resp.Tools, assistantToolResult{
			ToolId:      tool.UUID,
			Name:        tool.Name,
			ToolType:    string(tool.ToolType),
			Metadata:    tool.Metadata,
			Description: tool.Description,
			CreatedAt:   tool.CreatedAt,
			UpdatedAt:   tool.UpdatedAt,
		})
	}
	return
}

func (s *service) List(ctx context.Context, tenantId uint, name string, page, pageSize int) (resp []assistantResult, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	assistants, total, err := s.repository.Assistants().List(ctx, tenantId, name, page, pageSize, "Tools")
	if err != nil {
		err = errors.Wrap(err, "获取助手失败")
		_ = level.Warn(logger).Log("msg", "获取助手失败", "err", err)
		return
	}
	for _, assistant := range assistants {
		var tools []assistantToolResult
		for _, tool := range assistant.Tools {
			tools = append(tools, assistantToolResult{
				ToolId:      tool.UUID,
				Name:        tool.Name,
				ToolType:    string(tool.ToolType),
				Metadata:    tool.Metadata,
				Description: tool.Description,
				CreatedAt:   tool.CreatedAt,
				UpdatedAt:   tool.UpdatedAt,
			})
		}
		resp = append(resp, assistantResult{
			AssistantId:  assistant.UUID,
			Name:         assistant.Name,
			Remark:       assistant.Remark,
			ModelName:    assistant.ModelName,
			Description:  assistant.Description,
			Instructions: assistant.Instructions,
			Metadata:     assistant.Metadata,
			CreatedAt:    assistant.CreatedAt,
			UpdatedAt:    assistant.UpdatedAt,
			Tools:        tools,
			Avatar:       assistant.Avatar,
			Operator:     assistant.Operator,
		})
	}
	return
}

func (s *service) AddTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	err = s.repository.Assistants().AddTool(ctx, tenantId, assistantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "添加工具失败")
		_ = level.Warn(logger).Log("msg", "添加工具失败", "err", err)
		return
	}
	return err
}

func (s *service) RemoveTool(ctx context.Context, tenantId uint, assistantId, toolId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	err = s.repository.Assistants().RemoveTool(ctx, tenantId, assistantId, toolId)
	if err != nil {
		err = errors.Wrap(err, "移除工具失败")
		_ = level.Warn(logger).Log("msg", "移除工具失败", "err", err)
		return
	}
	return err
}

func (s *service) ListTool(ctx context.Context, tenantId uint, assistantId, name string, page, pageSize int) (resp []assistantToolResult, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	tools, total, err := s.repository.Assistants().ListTool(ctx, tenantId, assistantId, name, page, pageSize, "Tools")
	if err != nil {
		err = errors.Wrap(err, "获取助手工具失败")
		_ = level.Warn(logger).Log("msg", "获取助手工具失败", "err", err)
		return nil, 0, err
	}
	for _, tool := range tools {
		resp = append(resp, assistantToolResult{
			ToolId:      tool.UUID,
			Name:        tool.Name,
			ToolType:    string(tool.ToolType),
			Metadata:    tool.Metadata,
			Description: tool.Description,
			CreatedAt:   tool.CreatedAt,
			UpdatedAt:   tool.UpdatedAt,
		})
	}

	return
}

//// NewService returns a basic Service with all of the expected middlewares wired in.
//func NewService(logger log.Logger, traceId string, repository repository.Repository, middlewares []Middleware) Service {
//	var svc Service = NewBasicService(logger, traceId, repository)
//	for _, m := range middlewares {
//		svc = m(svc)
//	}
//	return svc
//}

// New returns a naive, stateless implementation of Service.
func New(logger log.Logger, traceId string, repository repository.Repository, httpOpts []kithttp.ClientOption, llmOptions []openai.Option) Service {
	logger = log.With(logger, "service", "assistants")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
		llmOptions: llmOptions,
		httpOpts:   httpOpts,
	}
}
