package chat

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

// EmbeddingPayload 嵌入载荷
type EmbeddingPayload struct {
	// Model 模型
	Model string `json:"model"`
	// Input 输入
	Input []string `json:"input"`
	// Output 输出
	EncodingFormat string `json:"encoding_format"`
	// User 用户
	User string `json:"user"`
	// Engine 引擎
	Engine string `json:"engine"`
}

// EmbeddingsResponse 嵌入响应
type EmbeddingsResponse struct {
	//Object string           `json:"object"`
	//Data   []map[string]any `json:"data"`
	//Model  string           `json:"model"`
	//Usage  UsageInfo        `json:"usage"`
	Embedding []interface {
	} `json:"embedding"`
	TokenNum  int    `json:"token_num"`
	Text      string `json:"text"`
	ErrorCode int    `json:"error_code"`
}

// UsageInfo 使用信息
type UsageInfo struct {
	// PromptTokens 提示令牌
	PromptTokens int `json:"prompt_tokens"`
	// TotalTokens 总令牌
	TotalTokens int `json:"total_tokens"`
	// CompletionTokens 完成令牌
	CompletionTokens int `json:"completion_tokens"`
}

// GenerateStreamParams 生成流参数
type GenerateStreamParams struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float32 `json:"temperature"`
	Logprobs         *bool   `json:"logprobs,omitempty"`
	TopP             float32 `json:"top_p"`
	TopK             int     `json:"top_k"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	MaxNewTokens     int     `json:"max_new_tokens"`
	Echo             bool    `json:"echo"`
	StopTokenIds     []int   `json:"stop_token_ids"`
	Images           []any   `json:"images,omitempty"`
	BestOf           *int    `json:"best_of,omitempty"`
	UseBeamSearch    *bool   `json:"use_beam_search,omitempty"`
	Stop             any     `json:"stop"`
	N                *int    `json:"n,omitempty"`
}

// GenerateParams 生成流参数
type GenerateParams struct {
	Model            string      `json:"model"`
	Prompt           string      `json:"prompt"`
	Temperature      float32     `json:"temperature"`
	Logprobs         int         `json:"logprobs"`
	TopP             float32     `json:"top_p"`
	TopK             int         `json:"top_k"`
	PresencePenalty  float32     `json:"presence_penalty"`
	FrequencyPenalty float32     `json:"frequency_penalty"`
	MaxNewTokens     int         `json:"max_new_tokens"`
	Echo             bool        `json:"echo"`
	StopTokenIds     interface{} `json:"stop_token_ids"`
	Images           []any       `json:"images"`
	BestOf           int         `json:"best_of"`
	UseBeamSearch    bool        `json:"use_beam_search"`
	Stop             any         `json:"stop"`
}

// WorkerCreationOptions is the options for the faceswap service.
type WorkerCreationOptions struct {
	httpClientOpts    []kithttp.ClientOption
	controllerAddress string
	logger            log.Logger
}

// WorkerCreationOption is the option for the chat service.
type WorkerCreationOption func(*WorkerCreationOptions)

// WithWorkerCreationOptionHTTPClientOpts is the option to set the http client options.
func WithWorkerCreationOptionHTTPClientOpts(opts ...kithttp.ClientOption) WorkerCreationOption {
	return func(o *WorkerCreationOptions) {
		o.httpClientOpts = opts
	}
}

// WithWorkerCreationOptionControllerAddress is the option to set the controller address.
func WithWorkerCreationOptionControllerAddress(controllerAddress string) WorkerCreationOption {
	return func(o *WorkerCreationOptions) {
		o.controllerAddress = controllerAddress
	}
}

// WithWorkerCreationOptionLogger is the option to set the logger.
func WithWorkerCreationOptionLogger(logger log.Logger) WorkerCreationOption {
	return func(o *WorkerCreationOptions) {
		o.logger = logger
	}
}

// ModelPermission 模型权限
type ModelPermission struct {
	ID                 string `json:"id"`
	Object             string `json:"object"`
	Created            int    `json:"created"`
	AllowCreateEngine  bool   `json:"allow_create_engine"`
	AllowSampling      bool   `json:"allow_sampling"`
	AllowLogprobs      bool   `json:"allow_logprobs"`
	AllowSearchIndices bool   `json:"allow_search_indices"`
	AllowView          bool   `json:"allow_view"`
	AllowFineTuning    bool   `json:"allow_fine_tuning"`
	Organization       string `json:"organization"`
	Group              string `json:"group"`
	IsBlocking         bool   `json:"is_blocking"`
}

// ModelCard 模型卡片
type ModelCard struct {
	ID string `json:"id"`
	// Object 对象
	Object string `json:"object"`
	// Created 创建时间
	Created int `json:"created"`
	// OwnedBy 拥有者
	OwnedBy string `json:"owned_by"`
	// Root 根
	Root string `json:"root"`
	// Parent 父级
	Parent string `json:"parent"`
	// Permission 权限
	Permission []ModelPermission `json:"permission"`
}

// ModelDetail 模型详情
type ModelDetail struct {
	// ContextLength 上下文长度
	ContextLength int `json:"context_length"`
}

// Conv 对话
type Conv struct {
	Name           string     `json:"name"`
	SystemTemplate string     `json:"system_template"`
	SystemMessage  string     `json:"system_message"`
	Roles          []string   `json:"roles"`
	Messages       [][]string `json:"messages"`
	Offset         int        `json:"offset"`
	SepStyle       int        `json:"sep_style"`
	Sep            string     `json:"sep"`
	Sep2           string     `json:"sep2"`
	StopStr        string     `json:"stop_str"`
	StopTokenIds   []int      `json:"stop_token_ids"`
	Template       string     `json:"template,omitempty"`
	Stop           []string   `json:"stop,omitempty"`
}

// ModelConvTemplate 对话模板
type ModelConvTemplate struct {
	Conv Conv `json:"conv"`
}

// WorkerStatus 工作状态
type WorkerStatus struct {
	// ModelNames 模型名称
	ModelNames []string `json:"model_names"`
	// Speed 速度
	Speed int `json:"speed"`
	// QueueLength 队列长度
	QueueLength int `json:"queue_length"`
}

type Logprobs struct {
	TextOffset    []int                `json:"text_offset"`
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
}

// WorkerGenerateStreamResponse 工作生成流响应
type WorkerGenerateStreamResponse struct {
	Text         string    `json:"text"`
	ErrorCode    int       `json:"error_code"`
	Usage        UsageInfo `json:"usage"`
	FinishReason string    `json:"finish_reason"`
	Logprobs     Logprobs  `json:"logprobs"`
}

type FsChatWorkerMiddleware func(WorkerService) WorkerService

// WorkerService 服务接口
type WorkerService interface {
	// ListModels 列出模型 调用 controller
	ListModels(ctx context.Context) (res []ModelCard, err error)
	// GetWorkerAddress 获取worker地址 调用 controller
	GetWorkerAddress(ctx context.Context, model string) (res string, err error)
	// WorkerGetConvTemplate 获取worker对话模板
	WorkerGetConvTemplate(ctx context.Context, workerAddress, model string) (res ModelConvTemplate, err error)
	// WorkerGenerateStream 生成worker对话流
	WorkerGenerateStream(ctx context.Context, workerAddress string, params GenerateStreamParams) (res <-chan WorkerGenerateStreamResponse, err error)
	// WorkerGenerate 生成worker对话
	WorkerGenerate(ctx context.Context, workerAddress string, params GenerateParams) (res <-chan WorkerGenerateStreamResponse, err error)
	// WorkerGetEmbeddings 获取worker嵌入
	WorkerGetEmbeddings(ctx context.Context, workerAddress string, payload EmbeddingPayload) (res EmbeddingsResponse, err error)
	// WorkerCountToken 计数令牌
	WorkerCountToken(ctx context.Context, workerAddress, model string, prompt any) (res int, err error)
	// WorkerGetStatus 获取worker状态
	WorkerGetStatus(ctx context.Context, workerAddress string) (res WorkerStatus, err error)
	// WorkerGetModelDetails 获取worker模型详情
	WorkerGetModelDetails(ctx context.Context, workerAddress, model string) (res ModelDetail, err error)
	// WorkerCheckLength 检查worker长度
	WorkerCheckLength(ctx context.Context, workerAddress string, model string, maxTokens int, prompt any) (res int, err error)
}
