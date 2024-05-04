package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/chat"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"strings"
	"time"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	workerSvc      chat.WorkerService
	openAISvc      chat.OpenAIService
}

// CreationOption is the option for the chat service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOpts is the option to set the http client options.
func WithHTTPClientOpts(opts ...kithttp.ClientOption) CreationOption {
	return func(o *CreationOptions) {
		o.httpClientOpts = opts
	}
}

// WithWorkerService is the option to set the worker service.
func WithWorkerService(svc chat.WorkerService) CreationOption {
	return func(o *CreationOptions) {
		o.workerSvc = svc
	}
}

type Service interface {
	// ChatCompletion 聊天处理
	ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error)
	// ChatCompletionStream 聊天处理流传输
	ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan openai.ChatCompletionStreamResponse, err error)
	// Models 模型列表
	Models(ctx context.Context, channelId uint) (res []openai.Model, err error)
	// Embeddings 向量化处理
	Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	options    *CreationOptions
	services   services.Service
	repository repository.Repository
}

func (s *service) Models(ctx context.Context, channelId uint) (res []openai.Model, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Embeddings(ctx context.Context, channelId uint, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) processInput(modelName string, inp any) (newInp []string) {
	fmt.Println(reflect.TypeOf(inp))
	if reflect.TypeOf(inp).Name() == "string" {
		newInp = []string{inp.(string)}
	} else if reflect.TypeOf(inp).Name() == "[]any" {
		fastInp := inp.([]any)
		if reflect.TypeOf(fastInp[0]).Name() == "int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				_ = level.Warn(s.logger).Log("msg", "model not found. Using cl100k_base encoding.")
				model := "cl100k_base"
				decoding, err = tiktoken.GetEncoding(model)
			}
			newInp = []string{decoding.Decode(inp.([]int))}
		} else if reflect.TypeOf(fastInp[0]).Name() == "[]int" {
			decoding, err := tiktoken.EncodingForModel(modelName)
			if err != nil {
				_ = level.Warn(s.logger).Log("msg", "model not found. Using cl100k_base encoding.")
				model := "cl100k_base"
				decoding, err = tiktoken.GetEncoding(model)
			}
			for _, text := range inp.([][]int) {
				newInp = append(newInp, decoding.Decode(text))
			}
		}
	}
	return
}

func (s *service) generateCompletionStreamGenerator(ctx context.Context, worker chat.WorkerService, req openai.CompletionRequest, n int, workerAddress string) {
	//modeName := req.Model
	//id := fmt.Sprintf("cmpl-%s", shortuuid.random())
	//finishStreamEvents := make([]openai.CompletionStreamResponse, 0)
	//prompts := req.Prompt.([]string)
	//for _, prompt := range prompts {
	//	for i := 0; i < n; i++ {
	//		previousText := ""
	//		genParams := chat.GenerateParams{
	//			Model:            req.Model,
	//			Prompt:           prompt,
	//			Temperature:      req.Temperature,
	//			TopP:             req.TopP,
	//			FrequencyPenalty: req.FrequencyPenalty,
	//			PresencePenalty:  req.PresencePenalty,
	//			MaxNewTokens:     req.MaxTokens,
	//			Logprobs:         req.LogProbs,
	//			Echo:             req.Echo,
	//			StopTokenIds:     nil,
	//			BestOf:           req.BestOf,
	//			Stop:             req.Stop,
	//		}
	//		content, err := worker.WorkerGenerate(ctx, workerAddress, genParams)
	//		if err != nil {
	//			_ = level.Warn(s.logger).Log("msg", "failed to generate completion", "err", err)
	//			return
	//		}
	//		decodedUnicode := content.Text
	//}
	// model_name = request.model
	//    id = f"cmpl-{shortuuid.random()}"
	//    finish_stream_events = []
	//    for text in request.prompt:
	//        for i in range(n):
	//            previous_text = ""
	//            gen_params = await get_gen_params(
	//                request.model,
	//                worker_addr,
	//                text,
	//                temperature=request.temperature,
	//                top_p=request.top_p,
	//                top_k=request.top_k,
	//                presence_penalty=request.presence_penalty,
	//                frequency_penalty=request.frequency_penalty,
	//                max_tokens=request.max_tokens,
	//                logprobs=request.logprobs,
	//                echo=request.echo,
	//                stop=request.stop,
	//            )
	//            async for content in generate_completion_stream(gen_params, worker_addr):
	//                if content["error_code"] != 0:
	//                    yield f"data: {json.dumps(content, ensure_ascii=False)}\n\n"
	//                    yield "data: [DONE]\n\n"
	//                    return
	//                decoded_unicode = content["text"].replace("\ufffd", "")
	//                delta_text = decoded_unicode[len(previous_text) :]
	//                previous_text = (
	//                    decoded_unicode
	//                    if len(decoded_unicode) > len(previous_text)
	//                    else previous_text
	//                )
	//                # todo: index is not apparent
	//                choice_data = CompletionResponseStreamChoice(
	//                    index=i,
	//                    text=delta_text,
	//                    logprobs=create_openai_logprobs(content.get("logprobs", None)),
	//                    finish_reason=content.get("finish_reason", None),
	//                )
	//                chunk = CompletionStreamResponse(
	//                    id=id,
	//                    object="text_completion",
	//                    choices=[choice_data],
	//                    model=model_name,
	//                )
	//                if len(delta_text) == 0:
	//                    if content.get("finish_reason", None) is not None:
	//                        finish_stream_events.append(chunk)
	//                    continue
	//                yield f"data: {chunk.json(exclude_unset=True, ensure_ascii=False)}\n\n"
	//    # There is not "content" field in the last delta message, so exclude_none to exclude field "content".
	//    for finish_chunk in finish_stream_events:
	//        yield f"data: {finish_chunk.json(exclude_unset=True, ensure_ascii=False)}\n\n"
	//    yield "data: [DONE]\n\n"
}

func (s *service) localAIChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (res <-chan openai.ChatCompletionStreamResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	svc := chat.NewFastChatWorker(chat.WithControllerAddress("http://localhost:21001"))
	models, err := svc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		_ = level.Warn(logger).Log("msg", "failed to list models", "err", err)
		return
	}
	var exists bool
	for _, model := range models {
		if model.ID == req.Model {
			exists = true
			break
		}
	}
	if !exists {
		err = errors.New("model not found")
		_ = level.Warn(logger).Log("msg", "model not found", "err", err)
		return
	}
	workerAddress, err := svc.GetWorkerAddress(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		_ = level.Warn(logger).Log("msg", "failed to get worker address", "err", err)
		return
	}
	_ = level.Info(logger).Log("msg", "worker address", "workerAddress", workerAddress)
	var maxTokens int
	prompts := s.processInput(req.Model, req.Messages)
	for _, prompt := range prompts {
		maxTokens, err = svc.WorkerCheckLength(ctx, workerAddress, req.Model, req.MaxTokens, prompt)
		if err != nil {
			err = errors.WithMessage(err, "failed to check length")
			_ = level.Warn(logger).Log("msg", "failed to check length", "err", err)
			return res, err
		}
	}
	_ = level.Info(logger).Log("msg", "max tokens", "maxTokens", maxTokens)
	if maxTokens != 0 && maxTokens < req.MaxTokens {
		req.MaxTokens = maxTokens
	}

	convTemplate, err := svc.WorkerGetConvTemplate(ctx, workerAddress, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to get conv template")
		_ = level.Warn(logger).Log("msg", "failed to get conv template", "err", err)
		return
	}

	var systemMessage = convTemplate.Conv.SystemMessage
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
			break
		}
	}

	if req.Stop == nil {
		// todo 应该多模型模版获取
		req.Stop = []string{convTemplate.Conv.StopStr}
	}
	stopTokenIds := convTemplate.Conv.StopTokenIds

	prompt := strings.ReplaceAll(convTemplate.Conv.SystemTemplate, "{system_message}", systemMessage)
	prompt += convTemplate.Conv.Sep
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			prompt += convTemplate.Conv.Roles[0] + "\n"
			prompt += msg.Content + convTemplate.Conv.Sep
		} else if msg.Role == "assistant" {
			prompt += convTemplate.Conv.Roles[1] + "\n"
			prompt += msg.Content + convTemplate.Conv.Sep
		}
	}
	prompt += convTemplate.Conv.Roles[1]
	_ = level.Debug(logger).Log("msg", "prompt", "prompt", prompt)

	// todo: 获取模版，并处理面prompt
	stream, err := svc.WorkerGenerateStream(ctx, workerAddress, chat.GenerateStreamParams{
		Model:            req.Model,
		Prompt:           prompt,
		Temperature:      req.Temperature,
		Logprobs:         req.LogProbs,
		TopP:             req.TopP,
		TopK:             -1,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxNewTokens:     req.MaxTokens,
		StopTokenIds:     stopTokenIds,
		Images:           nil,
		UseBeamSearch:    false,
		Stop:             req.Stop,
	})
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}

	for {
		content, ok := <-stream
		if !ok {
			fmt.Println("stream closed")
			break
		}
		if content.ErrorCode != 0 {
			fmt.Println("error")
			err = errors.New(content.Text)
			return
		}
		fmt.Println(content.Text, content.Usage)
	}

	dot := make(chan openai.ChatCompletionStreamResponse)
	go func() {
		now := time.Now().UnixMilli()
		for content := range stream {
			if content.ErrorCode != 0 {
				err = errors.New(content.Text)
				return
			}
			dot <- openai.ChatCompletionStreamResponse{
				ID:      fmt.Sprintf("cmpl-%s", shortuuid.New()),
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []openai.ChatCompletionStreamChoice{
					{
						FinishReason: openai.FinishReason(content.FinishReason),
						Delta: openai.ChatCompletionStreamChoiceDelta{
							Content: content.Text,
							Role:    "assistant",
						},
					},
				},
			}
		}
	}()
	return dot, nil
}

func (s *service) localAIChatCompletion(ctx context.Context, modelInfo types.Models, req openai.CompletionRequest) (res <-chan openai.ChatCompletionResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	svc := chat.NewFastChatWorker()
	models, err := svc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		_ = level.Warn(logger).Log("msg", "failed to list models", "err", err)
		return
	}
	var exists bool
	for _, model := range models {
		if model.ID == req.Model {
			exists = true
			break
		}
	}
	if !exists {
		err = errors.New("model not found")
		_ = level.Warn(logger).Log("msg", "model not found", "err", err)
		return
	}
	workerAddress, err := svc.GetWorkerAddress(ctx, modelInfo.ModelName)
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		_ = level.Warn(logger).Log("msg", "failed to get worker address", "err", err)
		return
	}
	var maxTokens int
	prompts := s.processInput(req.Model, req.Prompt)
	for _, prompt := range prompts {
		maxTokens, err = svc.WorkerCheckLength(ctx, workerAddress, req.Model, req.MaxTokens, prompt)
		if err != nil {
			err = errors.WithMessage(err, "failed to check length")
			_ = level.Warn(logger).Log("msg", "failed to check length", "err", err)
			return
		}
	}
	if maxTokens < req.MaxTokens {
		req.MaxTokens = maxTokens
	}
	req.Prompt = prompts
	// todo: 获取模版，并处理面prompt
	wsrStream, err := svc.WorkerGenerate(ctx, workerAddress, chat.GenerateParams{
		Model:            req.Model,
		Prompt:           "<|im_start|>system\nYou are a helpful assistant.<|im_end|><|im_start|>user\n您好!你叫什么名字？<|im_end|><|im_start|>assistant\n",
		Temperature:      req.Temperature,
		Logprobs:         req.LogProbs,
		TopP:             req.TopP,
		TopK:             -1,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxNewTokens:     req.MaxTokens,
		Echo:             req.Echo,
		StopTokenIds:     nil,
		Images:           nil,
		BestOf:           req.BestOf,
		UseBeamSearch:    false,
		Stop:             req.Stop,
	})
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}
	now := time.Now().UnixMilli()
	dot := make(chan openai.CompletionResponse)
	for content := range wsrStream {
		if content.ErrorCode != 0 {
			err = errors.New(content.Text)
			return
		}
		dot <- openai.CompletionResponse{
			ID:      fmt.Sprintf("cmpl-%s", shortuuid.New()),
			Object:  "chat.completion.chunk",
			Created: now,
			Model:   req.Model,
			Choices: []openai.CompletionChoice{
				{
					Index:        0,
					FinishReason: content.FinishReason,
					Text:         content.Text,
					LogProbs: openai.LogprobResult{
						Tokens:        content.Logprobs.Tokens,
						TokenLogprobs: content.Logprobs.TokenLogprobs,
						TopLogprobs:   content.Logprobs.TopLogprobs,
						TextOffset:    content.Logprobs.TextOffset,
					},
				},
			},
		}
	}

	// text_completions = []
	//        for text in request.prompt:
	//            gen_params = await get_gen_params(
	//                request.model,
	//                worker_addr,
	//                text,
	//                temperature=request.temperature,
	//                top_p=request.top_p,
	//                top_k=request.top_k,
	//                frequency_penalty=request.frequency_penalty,
	//                presence_penalty=request.presence_penalty,
	//                max_tokens=request.max_tokens,
	//                logprobs=request.logprobs,
	//                echo=request.echo,
	//                stop=request.stop,
	//                best_of=request.best_of,
	//                use_beam_search=request.use_beam_search,
	//            )
	//            for i in range(request.n):
	//                content = asyncio.create_task(
	//                    generate_completion(gen_params, worker_addr)
	//                )
	//                text_completions.append(content)
	//
	//        try:
	//            all_tasks = await asyncio.gather(*text_completions)
	//        except Exception as e:
	//            return create_error_response(ErrorCode.INTERNAL_ERROR, str(e))
	//
	//        choices = []
	//        usage = UsageInfo()
	//        for i, content in enumerate(all_tasks):
	//            if content["error_code"] != 0:
	//                return create_error_response(content["error_code"], content["text"])
	//            choices.append(
	//                CompletionResponseChoice(
	//                    index=i,
	//                    text=content["text"],
	//                    logprobs=create_openai_logprobs(content.get("logprobs", None)),
	//                    finish_reason=content.get("finish_reason", "stop"),
	//                )
	//            )
	//            task_usage = UsageInfo.parse_obj(content["usage"])
	//            for usage_key, usage_value in task_usage.dict().items():
	//                setattr(usage, usage_key, getattr(usage, usage_key) + usage_value)
	//
	//        return CompletionResponse(
	//            model=request.model, choices=choices, usage=UsageInfo.parse_obj(usage)
	//        )
	return
}

func (s *service) ChatCompletion(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return
	}
	svc := chat.NewFastChatWorker()
	if modelInfo.ProviderName == types.ModelProviderLocalAI {
		models, err := svc.ListModels(ctx)
		if err != nil {
			err = errors.WithMessage(err, "failed to list models")
			_ = level.Warn(logger).Log("msg", "failed to list models", "err", err)
			return res, err
		}
		var exists bool
		for _, model := range models {
			if model.ID == req.Model {
				exists = true
				break
			}
		}
		if !exists {
			err = errors.New("model not found")
			_ = level.Warn(logger).Log("msg", "model not found", "err", err)
			return res, err
		}
		workerAddress, err := svc.GetWorkerAddress(ctx, modelInfo.ModelName)
		if err != nil {
			err = errors.WithMessage(err, "failed to get worker address")
			_ = level.Warn(logger).Log("msg", "failed to get worker address", "err", err)
			return res, err
		}
		workerResult, err := svc.WorkerGenerate(ctx, workerAddress, chat.GenerateParams{
			Model:            req.Model,
			Prompt:           "",
			Temperature:      req.Temperature,
			Logprobs:         req.TopLogProbs,
			TopP:             req.TopP,
			PresencePenalty:  req.PresencePenalty,
			FrequencyPenalty: req.FrequencyPenalty,
			MaxNewTokens:     0,
			Echo:             false,
			StopTokenIds:     nil,
			Images:           nil,
			BestOf:           0,
			UseBeamSearch:    false,
			Stop:             req.Stop,
		})
		fmt.Println(workerResult)
	}
	return
}

func (s *service) ChatCompletionStream(ctx context.Context, channelId uint, req openai.ChatCompletionRequest) (stream <-chan openai.ChatCompletionStreamResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	dot := make(chan openai.ChatCompletionStreamResponse)
	defer func() {
		close(dot)
	}()

	modelInfo, err := s.repository.Model().FindByModelId(ctx, req.Model)
	if err != nil {
		err = errors.WithMessage(err, "failed to find model")
		_ = level.Warn(logger).Log("msg", "failed to find model", "err", err)
		return dot, nil
	}

	b, _ := json.Marshal(req.Messages)
	stop, _ := json.Marshal(req.Stop)
	msgData := &types.ChatMessages{
		ModelName:        req.Model,
		ChannelId:        channelId,
		Prompt:           req.Messages[len(req.Messages)-1].Content,
		Finished:         false,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		N:                req.N,
		User:             req.User,
		Messages:         string(b),
		Error:            true,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxTokens:        req.MaxTokens,
		Stop:             string(stop),
	}
	defer func() {
		if err = s.repository.Messages().Create(ctx, msgData); err != nil {
			err = errors.WithMessage(err, "failed to create message")
			_ = level.Warn(logger).Log("msg", "failed to create message", "err", err)
		}
	}()
	if modelInfo.BaseModelName != "" {
		req.Model = modelInfo.BaseModelName
	}
	completionStream, err := s.localAIChatCompletionStream(ctx, req)
	if err != nil {
		msgData.ErrorMessage = err.Error()
		_ = level.Warn(logger).Log("msg", "failed to get completion stream", "err", err)
		return dot, err
	}

	go func() {
		for content := range completionStream {
			if content.Choices[0].FinishReason == openai.FinishReasonStop {
				// 更新数据库
				msgData.Response = content.Choices[0].Delta.Content
				msgData.Error = false
				msgData.Created = content.Created
				msgData.TimeCost = time.Since(msgData.CreatedAt).String()
				msgData.Finished = true
				if err = s.repository.Messages().Update(ctx, msgData); err != nil {
					_ = level.Warn(logger).Log("msg", "failed to update message", "err", err)
				}
			}
			dot <- content
		}
	}()

	return dot, nil
}

func New(logger log.Logger, traceId string, store repository.Repository, services services.Service, opts ...CreationOption) Service {
	logger = log.With(logger, "service", "chat")
	options := &CreationOptions{}
	for _, o := range opts {
		o(options)
	}
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
		services:   services,
		options:    options,
	}
}
