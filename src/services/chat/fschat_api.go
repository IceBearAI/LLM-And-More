package chat

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"math"
	"strings"
	"time"
)

type fsChatApiClient struct {
	options   *CreationOptions
	template  Templates
	openaiSvc Service
}

func (s *fsChatApiClient) Completion(ctx context.Context, req openai.CompletionRequest) (res openai.CompletionResponse, err error) {
	streamResp, err := s.ChatCompletionStream(ctx, openai.ChatCompletionRequest{})
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}
	var content CompletionStreamResponse
	for {
		rs, ok := <-streamResp
		if !ok {
			break
		}
		content = rs
	}
	res = openai.CompletionResponse{
		Usage: content.Usage,
		Choices: []openai.CompletionChoice{
			{
				FinishReason: "stop",
				Text:         content.Choices[0].Delta.Content,
			},
		},
	}
	return
}

func (s *fsChatApiClient) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (res CompletionResponse, err error) {
	streamResp, err := s.ChatCompletionStream(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}
	var content CompletionStreamResponse
	var fullContent string
	for {
		rs, ok := <-streamResp
		if !ok {
			break
		}
		if len(rs.Choices) > 0 {
			content = rs
			fullContent += rs.Choices[0].Delta.Content
		}
	}
	res = CompletionResponse{
		Usage: openai.Usage{
			PromptTokens:     content.Usage.PromptTokens,
			CompletionTokens: content.Usage.CompletionTokens,
			TotalTokens:      content.Usage.TotalTokens,
		},
		ChatCompletionResponse: openai.ChatCompletionResponse{
			ID:      fmt.Sprintf("cmpl-%s", shortuuid.New()),
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []openai.ChatCompletionChoice{
				{
					FinishReason: "stop",
					Message: openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: fullContent,
					},
				},
			},
		},
	}
	return
}

func (s *fsChatApiClient) ChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (stream <-chan CompletionStreamResponse, err error) {
	workerAddress, err := s.getAddress(ctx, string(req.Model))
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		return
	}

	dot := make(chan CompletionStreamResponse)

	genParams, err := s.genParams(ctx, req, workerAddress)
	if err != nil {
		err = errors.WithMessage(err, "failed to get gen params")
		return
	}
	var newMaxTokens = genParams.MaxNewTokens
	if newMaxTokens, err = s.options.workerSvc.WorkerCheckLength(ctx, workerAddress, req.Model, req.MaxTokens, genParams.Prompt); err != nil {
		newMaxTokens = 1024
	}
	genParams.MaxNewTokens = newMaxTokens

	streamResp, err := s.options.workerSvc.WorkerGenerateStream(ctx, workerAddress, genParams)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate stream")
		return
	}

	go func() {
		now := time.Now().Unix()
		defer close(dot)
		streamId := fmt.Sprintf("cmpl-%s", shortuuid.New())
		var previousText string

		for {
			content, ok := <-streamResp
			if !ok {
				break
			}
			if content.ErrorCode != 0 {
				err = errors.New(content.Text)
				return
			}

			text := content.Text
			// 替换所有的Unicode替代字符\ufffd为空字符串
			decodedUnicode := strings.Replace(text, "\ufffd", "", -1)

			// 获取新的字符串，它是当前文本去掉与之前文本相同部分后的结果
			deltaText := decodedUnicode
			if len(previousText) < len(decodedUnicode) {
				deltaText = decodedUnicode[len(previousText):]
			}
			if len(decodedUnicode) > len(previousText) {
				previousText = decodedUnicode
			} else {
				deltaText = ""
			}

			dot <- CompletionStreamResponse{
				Usage: struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				}{PromptTokens: content.Usage.PromptTokens, CompletionTokens: content.Usage.CompletionTokens, TotalTokens: content.Usage.TotalTokens},
				ChatCompletionStreamResponse: openai.ChatCompletionStreamResponse{
					ID:      streamId,
					Object:  "chat.completion.chunk",
					Created: now,
					Model:   req.Model,
					Choices: []openai.ChatCompletionStreamChoice{
						{
							FinishReason: openai.FinishReason(content.FinishReason),
							Delta: openai.ChatCompletionStreamChoiceDelta{
								Content: deltaText,
								Role:    "assistant",
							},
						},
					},
				},
			}
			if content.FinishReason == "stop" {
				return
			}
		}
	}()

	return dot, nil
}

func (s *fsChatApiClient) Models(ctx context.Context) (res []openai.Model, err error) {
	models, err := s.options.workerSvc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return nil, err
	}
	for _, model := range models {
		res = append(res, openai.Model{
			ID:   model.ID,
			Root: model.Root,
		})
	}
	return
}

func (s *fsChatApiClient) getAddress(ctx context.Context, modelName string) (workerAddress string, err error) {
	models, err := s.options.workerSvc.ListModels(ctx)
	if err != nil {
		err = errors.WithMessage(err, "failed to list models")
		return
	}
	var exists bool
	for _, model := range models {
		if strings.ToLower(model.ID) == strings.ToLower(modelName) ||
			strings.ToLower(model.Root) == strings.ToLower(modelName) {
			exists = true
			break
		}
	}
	if !exists {
		err = errors.New("model not found")
		return
	}
	workerAddress, err = s.options.workerSvc.GetWorkerAddress(ctx, modelName)
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		return
	}
	return

}

func processEmbedding(emb interface{}, n int) (openai.Embedding, error) {
	if embd, ok := emb.([]interface{}); ok {
		var embeddings []float32
		for _, v := range embd {
			if vv, vvOk := v.([]interface{}); vvOk {
				for _, vvv := range vv {
					embeddings = append(embeddings, float32(vvv.(float64)))
				}
			} else if vv, vvOk := v.(string); vvOk {
				respf, err := embeddingDecode(vv)
				if err != nil {
					return openai.Embedding{}, err
				}
				embeddings = append(embeddings, respf...)
			}
		}
		return openai.Embedding{
			Index:     n,
			Object:    "embedding",
			Embedding: embeddings,
		}, nil
	} else if embStr, ok := emb.(string); ok {
		respf, err := embeddingDecode(embStr)
		if err != nil {
			return openai.Embedding{}, err
		}
		return openai.Embedding{
			Index:     n,
			Object:    "embedding",
			Embedding: respf,
		}, nil
	}
	return openai.Embedding{}, nil
}

func (s *fsChatApiClient) Embeddings(ctx context.Context, req openai.EmbeddingRequest) (res openai.EmbeddingResponse, err error) {
	workerAddress, err := s.getAddress(ctx, string(req.Model))
	if err != nil {
		err = errors.WithMessage(err, "failed to get worker address")
		return
	}

	var input []string
	input = processInput(string(req.Model), req.Input)
	var data []openai.Embedding
	var tokenNum int
	var batch []string
	batchSize := 64
	var n = 0
	for i := 0; i < len(input); i += batchSize {
		batch = input[i:min(i+batchSize, len(input))]
		payload := EmbeddingPayload{
			Model:          string(req.Model),
			Input:          batch,
			EncodingFormat: string(req.EncodingFormat),
			User:           req.User,
		}
		resp, err := s.options.workerSvc.WorkerGetEmbeddings(ctx, workerAddress, payload)
		if err != nil {
			err = errors.WithMessage(err, "failed to get embeddings")
			return res, err
		}
		if resp.ErrorCode != 0 {
			err = errors.New(resp.Text)
			return res, err
		}
		embedding, err := processEmbedding(resp.Embedding, n)
		if err != nil {
			err = errors.WithMessage(err, "failed to process embedding")
			return res, err
		}
		data = append(data, embedding)
	}
	res = openai.EmbeddingResponse{
		Object: "list",
		Model:  req.Model,
		Usage:  openai.Usage{PromptTokens: tokenNum, TotalTokens: tokenNum},
		Data:   data,
	}

	return
}

// embeddingDecode 方法将 base64 编码的字符串解码为 float32 列表
func embeddingDecode(encodedData string) ([]float32, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}
	floats := make([]float32, len(bytes)/4)
	for i := 0; i < len(floats); i++ {
		floats[i] = math.Float32frombits(binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4]))
	}
	return floats, nil
}

// Encode 方法将 float32 列表编码为 base64 编码的字符串
func embeddingEncode(floats []float32) string {
	bytes := make([]byte, len(floats)*4)
	for i, f := range floats {
		binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], math.Float32bits(f))
	}
	encodedData := base64.StdEncoding.EncodeToString(bytes)
	return encodedData
}

func NewFsChatApi(opts ...CreationOption) Service {
	options := &CreationOptions{
		endpoints: []Endpoint{
			{
				Host:     "http://localhost:8000/v1",
				Token:    "",
				Platform: "localai",
			},
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	tp := NewTemplates()
	return &fsChatApiClient{
		options:  options,
		template: Register(tp),
	}
}

func (s *fsChatApiClient) genParams(ctx context.Context, req openai.ChatCompletionRequest, workerAddress string) (params GenerateStreamParams, err error) {

	conv, ok := s.template.GetByModelName(ctx, req.Model)
	if !ok {
		err = errors.New("failed to get conv template")
		return
	}
	var imageList []string
	for _, v := range req.Messages {
		if v.Role == "system" {
			conv.SetSystemMessage(strings.TrimSpace(v.Content))
		} else if v.Role == "user" {
			if v.MultiContent != nil {
				var textList []string
				for _, item := range v.MultiContent {
					if item.Type == openai.ChatMessagePartTypeImageURL {
						if item.ImageURL != nil {
							imageList = append(imageList, item.ImageURL.URL)
						}
					} else {
						textList = append(textList, item.Text)
					}
				}
				text := strings.Repeat("<image>\n", len(imageList))
				text += strings.Join(textList, "\n")
				conv.AppendMessage(conv.Roles[0], text)
			} else {
				conv.AppendMessage(conv.Roles[0], strings.TrimSpace(v.Content))
			}
		} else if v.Role == "assistant" {
			conv.AppendMessage(conv.Roles[1], strings.TrimSpace(v.Content))
		}
	}
	conv.AppendMessage(conv.Roles[1], "")

	//convTemplate, err := s.options.workerSvc.WorkerGetConvTemplate(ctx, workerAddress, req.Model)
	//if err != nil {
	//	err = errors.WithMessage(err, "failed to get conv template")
	//	return
	//}
	prompt := conv.GetPrompt()

	if req.Stop == nil && conv.StopStr != nil {
		req.Stop = append(req.Stop, conv.StopStr...)
	}

	genParams := GenerateStreamParams{
		Model:            req.Model,
		Prompt:           prompt,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		TopK:             -1,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxNewTokens:     req.MaxTokens,
		StopTokenIds:     conv.StopTokenIds,
		Stop:             req.Stop,
		Echo:             false,
		Images:           imageList,
	}
	if req.N > 0 {
		genParams.N = &req.N
	}
	if req.LogProbs {
		logProbs := true
		genParams.Logprobs = &logProbs
	}

	return genParams, nil
}

func processInput(modelName string, inp any) (newInp []string) {
	switch v := inp.(type) {
	case string:
		newInp = append(newInp, v)
	case []string:
		newInp = append(newInp, v...)
	case []interface{}:
		prompts, _ := ConvertToSliceOfStrings(v)
		newInp = append(newInp, prompts...)
	case []int:
		decoding, err := getEncoding(modelName)
		if err != nil {
			return
		}
		newInp = append(newInp, decoding.Decode(v))
	case [][]int:
		decoding, err := getEncoding(modelName)
		if err != nil {
			return
		}
		for _, text := range v {
			newInp = append(newInp, decoding.Decode(text))
		}
	}

	return newInp
}

func getEncoding(modelName string) (*tiktoken.Tiktoken, error) {
	decoding, err := tiktoken.EncodingForModel(modelName)
	if err != nil {
		modelName = "cl100k_base"
		decoding, err = tiktoken.GetEncoding(modelName)
	}
	return decoding, err
}

func ConvertToSliceOfStrings(input []interface{}) ([]string, error) {
	var result []string
	for _, v := range input {
		if str, ok := v.(string); ok {
			result = append(result, str)
		} else {
			return nil, fmt.Errorf("element is not a string")
		}
	}
	return result, nil
}
