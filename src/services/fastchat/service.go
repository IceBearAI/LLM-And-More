package fastchat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/util"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Config struct {
	Debug bool
	// OpenAI地址
	OpenAiEndpoint string
	// OpenAI token
	OpenAiToken string
	// OpenAI 默认模型
	OpenAiModel string
	// OpenAI 组织ID
	OpenAiOrgId string
	// localai-api 地址
	LocalAiEndpoint string
	LocalAiToken    string
}

type CtxPlatform string
type Platform string

const (
	PlatformOpenAI Platform = "OpenAI"

	ContextKeyPlatform CtxPlatform = "ctx-platform-name"
	ContextKeyApiKey   CtxPlatform = "ctx-api-key"
)

type ApiErrResponse struct {
	Object  string `json:"object"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Middleware func(service Service) Service

type Service interface {
	// UploadFile
	// file string Required
	// Name of the JSON Lines file to be uploaded.
	// If the purpose is set to "fine-tune", the file will be used for fine-tuning.
	// purpose string Required
	// The intended purpose of the uploaded documents.
	// Use "fine-tune" for fine-tuning. This allows us to validate the format of the uploaded file.
	UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error)
	// CreateFineTuningJob 创建微调任务
	CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error)
	// RetrieveFineTuningJob 获取微调任务
	RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error)
	ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error)
	// CancelFineTuningJob 取消微调任务
	CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error)
	// ChatCompletion 聊天处理
	// model: str
	// messages: List[Dict[str, str]]
	// temperature: Optional[float] = 0.7
	// top_p: Optional[float] = 1.0
	// n: Optional[int] = 1
	// max_tokens: Optional[int] = None
	// stop: Optional[Union[str, List[str]]] = None
	// stream: Optional[bool] = False
	// presence_penalty: Optional[float] = 0.0
	// frequency_penalty: Optional[float] = 0.0
	// user: Optional[str] = None
	ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature float64, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error)
	// Models 模型列表
	Models(ctx context.Context) (res []openai.Model, err error)
	// Embeddings 创建图片
	Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error)
	// ModeRations 检测量否有不当内容
	ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error)
	// CreateImage 创建图片
	// Deprecated: use CreateSdImage
	// 暂时还无法使用
	CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error)
	// CheckLength 验证Token是否超过相应长度
	CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error)
	// CreateChatCompletionStream OpenAI Chat Completion openai 改版参数变化较大，直接使用 OpenAI入参
	CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error)
}

type chatGPTToken struct {
	Token          string `json:"token"`
	OrganizationId string `json:"organizationId"`
}

type service struct {
	localAiHost, chatGPTHost, chatGPTToken, chatGPTModel, chatGPTOrgId string
	localAiToken                                                       string
	opts                                                               []kithttp.ClientOption
	debug                                                              bool
	chatGPTTokens                                                      []chatGPTToken
}

func (s *service) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	var client *openai.Client
	client, _ = s.getClient(ctx, req.Model)
	if req.Temperature == 0 {
		req.Temperature = math.SmallestNonzeroFloat32
	}
	if req.TopP == 0 {
		req.TopP = math.SmallestNonzeroFloat32
	}
	if req.PresencePenalty == 0 {
		req.PresencePenalty = math.SmallestNonzeroFloat32
	}
	if req.FrequencyPenalty == 0 {
		req.FrequencyPenalty = math.SmallestNonzeroFloat32
	}
	stream, err = client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "CreateChatCompletionStream")
		return stream, err
	}
	return
}

func (s *service) CancelFineTuningJob(ctx context.Context, modelName, jobId string) (err error) {
	client, _ := s.getClient(ctx, modelName)
	_, err = client.CancelFineTuningJob(ctx, jobId)
	if err != nil {
		err = errors.Wrap(err, "cancel fine tune")
		return
	}
	return
}

func (s *service) ListFineTune(ctx context.Context, modelName string) (res openai.FineTuneList, err error) {
	client, _ := s.getClient(ctx, modelName)
	tunes, err := client.ListFineTunes(ctx)
	if err != nil {
		err = errors.Wrap(err, "list fine tunes")
		return
	}
	return tunes, nil
}

func (s *service) RetrieveFineTuningJob(ctx context.Context, jobId string) (res openai.FineTuningJob, err error) {
	client, _ := s.getClient(ctx, openai.GPT3Dot5Turbo)
	tune, err := client.RetrieveFineTuningJob(ctx, jobId)
	if err != nil {
		err = errors.Wrap(err, "get fine tune")
		return
	}
	return tune, nil
}

func (s *service) UploadFile(ctx context.Context, modelName, fileName, filePath, purpose string) (res openai.File, err error) {
	client, _ := s.getClient(ctx, modelName)

	file, err := client.CreateFile(ctx, openai.FileRequest{
		FileName: fileName,
		FilePath: filePath,
		Purpose:  "fine-tune",
	})
	if err != nil {
		err = errors.Wrap(err, "create file")
		return
	}
	b, _ := json.Marshal(file)
	fmt.Println(string(b))
	return file, nil
}

func (s *service) CreateFineTuningJob(ctx context.Context, req openai.FineTuningJobRequest) (res openai.FineTuningJob, err error) {
	client, _ := s.getClient(ctx, req.Model)

	tune, err := client.CreateFineTuningJob(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "create fine tune")
		return
	}
	//if tune.Status == "failed" {
	//	err = errors.New(tune.Error)
	//	return err
	//}
	b, _ := json.Marshal(tune)
	fmt.Println(string(b))
	return tune, nil
}

func (s *service) ModeRations(ctx context.Context, model, input string) (res openai.ModerationResponse, err error) {
	client, _ := s.getClient(ctx, model)
	res, err = client.Moderations(ctx, openai.ModerationRequest{
		Input: input, Model: model,
	})
	return res, err
}

func (s *service) Embeddings(ctx context.Context, model string, documents any) (res openai.EmbeddingResponse, err error) {
	if isOpenAiEmbeddingModel(model) {
		client, _ := s.getClient(ctx, model)
		res, err = client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Input: documents,
			Model: openai.AdaEmbeddingV2, // 暂时写死
		})
		if err != nil {
			err = errors.Wrap(err, "Embeddings")
			return openai.EmbeddingResponse{}, err
		}
		return res, nil
	}

	tgt, _ := url.Parse(fmt.Sprintf("%s/v1/embeddings", s.localAiHost))
	ep := kithttp.NewClient(http.MethodPost, tgt, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.opts...).Endpoint()
	_, err = ep(ctx, map[string]any{
		"input": documents,
		"model": model,
	})
	if err != nil {
		err = errors.Wrap(err, "LocalAIEmbeddings")
		return openai.EmbeddingResponse{}, err
	}
	return res, nil
}

func (s *service) CreateSdImageV1(ctx context.Context, prompt, negativePrompt, samplerIndex string, steps int) (res string, err error) {

	return
}

func (s *service) CheckLength(ctx context.Context, prompt string, maxToken int) (tokenNum int, err error) {
	tgt, _ := url.Parse(fmt.Sprintf("%s/worker/count_token", s.localAiHost))

	ep := kithttp.NewClient(http.MethodPost, tgt, kithttp.EncodeJSONRequest, func(ctx context.Context, response *http.Response) (response1 interface{}, e error) {
		if response.StatusCode != http.StatusOK {
			var res ApiErrResponse
			if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
				return nil, err
			}
			return nil, errors.New(res.Message)
		}

		var res struct {
			Count     int `json:"count"`
			ErrorCode int `json:"error_code"`
		}
		if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res.Count, nil
	}, s.opts...).Endpoint()

	type req struct {
		Prompt string `json:"prompt"`
	}

	if res, err := ep(ctx, req{
		Prompt: prompt,
	}); err == nil {
		tokenNum = res.(int)
	}

	var contextLength = 2048 // TODO 调用接口获得

	if tokenNum+maxToken > contextLength {
		//return tokenNum, errors.New(fmt.Sprintf("token num %d + max token %d > context length %d", tokenNum, maxToken, contextLength))
	}

	return
}

func (s *service) CreateImage(ctx context.Context, prompt, size, format string) (res []openai.ImageResponseDataInner, err error) {
	c := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:            fmt.Sprintf("%s/v1", s.localAiHost),
		EmptyMessagesLimit: 300,
		//HTTPClient:         http.DefaultClient,
		HTTPClient: &http.Client{
			Transport: &proxyRoundTripper{},
		},
	})

	//translation, err := c.CreateTranslation(ctx, openai.AudioRequest{
	//	Model:    openai.Whisper1,
	//	FilePath: "~/Downloads/langchain.wav",
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Println(translation.Text)

	image, err := c.CreateImage(ctx, openai.ImageRequest{
		Prompt:         prompt,
		N:              1,
		Size:           size,
		ResponseFormat: format,
	})
	if err != nil {
		err = errors.Wrap(err, "CreateImage")
		return nil, err
	}

	return image.Data, nil
}

func (s *service) Models(ctx context.Context) (res []openai.Model, err error) {
	client, _ := s.getClient(ctx, "")
	models, err := client.ListModels(ctx)
	if err != nil {
		err = errors.Wrap(err, "ListModels")
		return nil, err
	}

	return models.Models, nil
}

func (s *service) getClient(ctx context.Context, model string) (*openai.Client, string) {
	httpClient := http.DefaultClient
	if s.debug {
		httpClient = &http.Client{
			Transport: &proxyRoundTripper{},
		}
	}

	if isOpenAiModel(model) || isOpenAiEmbeddingModel(model) {
		ran := rand.Intn(len(s.chatGPTTokens))
		token := s.chatGPTTokens[ran].Token
		config := openai.DefaultConfig(token)
		config.BaseURL = s.chatGPTHost
		config.HTTPClient = httpClient
		return openai.NewClientWithConfig(config), model
	}

	fmt.Println("apiKey", s.localAiToken)

	apiKey := s.localAiToken
	platform, ok := ctx.Value(ContextKeyPlatform).(Platform)
	if !ok || platform == "" {
		if key, exists := ctx.Value(ContextKeyApiKey).(string); exists && key != "" {
			apiKey = key
		}
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = fmt.Sprintf("%s/v1", s.localAiHost)
		config.HTTPClient = httpClient
		return openai.NewClientWithConfig(config), model
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = s.chatGPTHost
	config.HTTPClient = httpClient
	return openai.NewClientWithConfig(config), model
}

func (s *service) ChatCompletionStream(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (stream *openai.ChatCompletionStream, status int, err error) {
	var client *openai.Client
	client, model = s.getClient(ctx, model)
	req := openai.ChatCompletionRequest{
		Model:            model,
		MaxTokens:        maxToken,
		Messages:         messages,
		Stream:           true,
		TopP:             float32(topP),
		Temperature:      float32(temperature),
		N:                n,
		PresencePenalty:  float32(presencePenalty),
		FrequencyPenalty: float32(frequencyPenalty),
		Stop:             stop,
		User:             user,
		Functions:        functions,
	}
	if functionCall != nil {
		if callStr, ok := functionCall.(string); ok && callStr == "" {
			req.FunctionCall = nil
		} else {
			req.FunctionCall = functionCall
		}
	}
	stream, err = client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		var er *openai.APIError
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = er
		}
		return
	}

	return stream, status, nil
}

func (s *service) ChatCompletion(ctx context.Context, model string, messages []openai.ChatCompletionMessage, temperature, topP, presencePenalty, frequencyPenalty float64, maxToken, n int, stop []string, user string, functions []openai.FunctionDefinition, functionCall any) (res openai.ChatCompletionResponse, status int, err error) {
	var client *openai.Client
	client, model = s.getClient(ctx, model)

	req := openai.ChatCompletionRequest{
		Model:            model,
		MaxTokens:        maxToken,
		Messages:         messages,
		Stream:           false,
		TopP:             float32(topP),
		Temperature:      float32(temperature),
		N:                n,
		PresencePenalty:  float32(presencePenalty),
		FrequencyPenalty: float32(frequencyPenalty),
		Stop:             stop,
		User:             user,
		Functions:        functions,
	}
	if functionCall != nil {
		if callStr, ok := functionCall.(string); ok && callStr == "" {
			req.FunctionCall = nil
		} else {
			req.FunctionCall = functionCall
		}
	}
	res, err = client.CreateChatCompletion(ctx, req)
	if err != nil {
		er := &openai.RequestError{}
		//apiErr := &openai.APIError{}
		if errors.As(err, &er) {
			status = er.HTTPStatusCode
			err = errors.New(fmt.Sprintf("%s %s", er.Error(), "可能是Token超过该模型的最大token限制了"))
		}

		return res, status, err
	}
	return
}

func decodeJsonResponse(data interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error) {
	return func(ctx context.Context, res *http.Response) (response interface{}, err error) {
		if res.StatusCode == 422 {
			body, _ := io.ReadAll(res.Body)
			return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
		}
		/*		if res.StatusCode != 200 {
				body, _ := io.ReadAll(res.Body)
				return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
			}*/
		if data == nil {
			return res, nil
		}
		if err = json.NewDecoder(res.Body).Decode(data); err != nil {
			return res, errors.Wrap(err, "json decode")
		}
		return res, nil
	}
}

func New(cfg Config, opts []kithttp.ClientOption) Service {
	chatGPTTokens := parseTokens(cfg.OpenAiToken)
	return &service{
		localAiHost:   cfg.LocalAiEndpoint,
		opts:          opts,
		chatGPTHost:   cfg.OpenAiEndpoint,
		chatGPTToken:  cfg.OpenAiToken,
		chatGPTOrgId:  cfg.OpenAiOrgId,
		debug:         cfg.Debug,
		localAiToken:  cfg.LocalAiToken,
		chatGPTTokens: chatGPTTokens,
	}
}

type proxyRoundTripper struct {
	traceId string
	before  []kithttp.RequestFunc
	after   []kithttp.ClientResponseFunc
}

func (s *proxyRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))
	defer func() {
		if res != nil {
			dump, _ = httputil.DumpResponse(res, true)
			fmt.Println(string(dump))
		}
	}()
	return http.DefaultTransport.RoundTrip(req)
}

func TokensNumFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	if strings.EqualFold(model, "") {
		model = openai.GPT3Dot5Turbo
	}
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokensPerMessage int
	var tokensPerName int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokensPerMessage = 4
		tokensPerName = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokensPerMessage = 3
		tokensPerName = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokensPerMessage = 3
		tokensPerName = 1
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3
	return numTokens
}

// TokenizerGetWord 只支持英文
func TokenizerGetWord(text string, model string) []string {
	if strings.EqualFold(model, "") {
		model = openai.GPT3Dot5Turbo
	}
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return nil
	}
	var words []string
	tokens := tkm.Encode(text, nil, nil)
	for _, v := range tokens {
		words = append(words, tkm.Decode([]int{v}))
	}
	return words
}

// InputToString 将输入转换为字符串
func InputToString(input any) (res string, err error) {
	tk, err := tiktoken.EncodingForModel(openai.GPT3Dot5Turbo)
	if err != nil {
		return "", err
	}
	switch input.(type) {
	case string:
		return input.(string), nil
	case []string:
		return tk.Decode(tk.Encode(strings.Join(input.([]string), " "), nil, nil)), nil
	case [][]int:
		var allInput []int
		for _, val := range input.([][]int) {
			for _, num := range val {
				allInput = append(allInput, num)
			}
		}
		return tk.Decode(allInput), nil
	}
	return res, nil
}

func isOpenAiModel(model string) bool {
	return util.StringContainsArray([]string{
		openai.GPT3Dot5Turbo,
		openai.GPT432K0613,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
		openai.GPT4,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5TurboInstruct,
		openai.GPT3Davinci,
		openai.GPT3Davinci002,
	}, model)
}

func isOpenAiEmbeddingModel(model string) bool {
	return util.StringInArray([]string{
		"text-similarity-ada-001",
		"text-similarity-babbage-001",
		"text-similarity-curie-001",
		"text-similarity-davinci-001",
		"text-search-ada-doc-001",
		"text-search-ada-query-001",
		"text-search-babbage-doc-001",
		"text-search-babbage-query-001",
		"text-search-curie-doc-001",
		"text-search-curie-query-001",
		"text-search-davinci-doc-001",
		"text-search-davinci-query-001",
		"code-search-ada-code-001",
		"code-search-ada-text-001",
		"code-search-babbage-code-001",
		"code-search-babbage-text-001",
		"text-embedding-ada-002",
		"text-embedding-large",
	}, model)
}

func parseTokens(cfgToken string) []chatGPTToken {
	var chatGPTTokens []chatGPTToken

	parseToken := func(token string) {
		split := strings.SplitN(token, ":", 2)
		chatToken := chatGPTToken{OrganizationId: "", Token: split[0]}

		if len(split) > 1 {
			chatToken.OrganizationId = split[0]
			chatToken.Token = split[1]
		}

		chatGPTTokens = append(chatGPTTokens, chatToken)
	}

	for _, token := range strings.Split(cfgToken, ",") {
		parseToken(token)
	}

	return chatGPTTokens
}
