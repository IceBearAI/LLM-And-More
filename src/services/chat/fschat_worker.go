package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

const (
	// ImagePlaceholderStr 图片占位符
	ImagePlaceholderStr = "$$<image>$$"
)

// WithControllerAddress is the option to set the controller address.
func WithControllerAddress(addr string) WorkerCreationOption {
	return func(o *WorkerCreationOptions) {
		o.controllerAddress = addr
	}
}

type ErrorResponse struct {
	Text      string `json:"text"`
	ErrorCode int    `json:"error_code"`
}

type worker struct {
	options *WorkerCreationOptions
}

func (s *worker) WorkerCheckLength(ctx context.Context, workerAddress string, model string, maxTokens int, prompt any) (res int, err error) {
	if maxTokens <= 0 {
		maxTokens = 1024
	}
	mc, err := s.WorkerGetModelDetails(ctx, workerAddress, model)
	if err != nil {
		err = errors.Wrap(err, "failed to get model details")
		return
	}
	contextLen := mc.ContextLength
	tokenNum, err := s.WorkerCountToken(ctx, workerAddress, model, prompt)
	if err != nil {
		err = errors.Wrap(err, "failed to count token")
		return
	}
	res = min(maxTokens, contextLen-tokenNum)
	if res <= 0 {
		err = errors.Errorf("This model's maximum context length is %d tokens. However, your messages resulted in %d tokens. Please reduce the length of the messages.", contextLen, tokenNum)
		return 0, err
	}
	return
}

func (s *worker) ListModels(ctx context.Context) (modelList []ModelCard, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/list_models", s.options.controllerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}

	var res struct {
		Models []string `json:"models"`
	}
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	for _, v := range res.Models {
		modelList = append(modelList, ModelCard{
			ID:   v,
			Root: v,
		})
	}
	return
}

func (s *worker) GetWorkerAddress(ctx context.Context, model string) (res string, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/get_worker_address", s.options.controllerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	type workerAddressResponse struct {
		Address string `json:"address"`
	}
	var war workerAddressResponse
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&war), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, map[string]string{"model": model})
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return war.Address, nil
}

func (s *worker) WorkerGetConvTemplate(ctx context.Context, workerAddress string, model string) (res ModelConvTemplate, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/worker_get_conv_template", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, map[string]any{"model": model})
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return
}
func (s *worker) WorkerGenerateStream(ctx context.Context, workerAddress string, params GenerateStreamParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/worker_generate_stream", workerAddress))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	opts := s.options.httpClientOpts
	opts = append(opts, kithttp.BufferedStream(true))

	ep := kithttp.NewClient(http.MethodPost, u, encodeRequest, decodeResponse, opts...).Endpoint()
	resStream, err := ep(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call endpoint")
	}

	dot := make(chan WorkerGenerateStreamResponse)
	go func() {
		defer close(dot)
		rc := resStream.(io.ReadCloser)
		defer rc.Close()

		scanner := bufio.NewScanner(rc)
		scanner.Split(splitAtDelimiter([]byte("\x00")))

		for scanner.Scan() {
			var resp WorkerGenerateStreamResponse
			if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
				dot <- WorkerGenerateStreamResponse{
					ErrorCode:    1,
					Text:         err.Error(),
					FinishReason: "stop",
				}
				return
			}
			dot <- resp
		}

		if err := scanner.Err(); err != nil {
			dot <- WorkerGenerateStreamResponse{
				ErrorCode:    1,
				Text:         err.Error(),
				FinishReason: "stop",
			}
		}
	}()

	return dot, nil
}

func (s *worker) WorkerGenerate(ctx context.Context, workerAddress string, params GenerateParams) (res <-chan WorkerGenerateStreamResponse, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/worker_generate", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	opts := s.options.httpClientOpts
	opts = append(opts, kithttp.BufferedStream(true))

	ep := kithttp.NewClient(http.MethodPost, u, encodeRequest, decodeResponse, opts...).Endpoint()
	resStream, err := ep(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call endpoint")
	}

	dot := make(chan WorkerGenerateStreamResponse)
	go func() {
		defer close(dot)
		rc := resStream.(io.ReadCloser)
		defer rc.Close()

		scanner := bufio.NewScanner(rc)
		scanner.Split(splitAtDelimiter([]byte("\x00")))

		for scanner.Scan() {
			var resp WorkerGenerateStreamResponse
			if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
				dot <- WorkerGenerateStreamResponse{
					ErrorCode:    1,
					Text:         err.Error(),
					FinishReason: "stop",
				}
				return
			}
			dot <- resp
		}

		if err := scanner.Err(); err != nil {
			dot <- WorkerGenerateStreamResponse{
				ErrorCode:    1,
				Text:         err.Error(),
				FinishReason: "stop",
			}
		}
	}()

	return dot, nil
}

func (s *worker) WorkerGetEmbeddings(ctx context.Context, workerAddress string, payload EmbeddingPayload) (res EmbeddingsResponse, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/worker_get_embeddings", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}

	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, payload)
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return

	// request.input = process_input(request.model, request.input)
	//
	//    data = []
	//    token_num = 0
	//    batch_size = WORKER_API_EMBEDDING_BATCH_SIZE
	//    batches = [
	//        request.input[i : min(i + batch_size, len(request.input))]
	//        for i in range(0, len(request.input), batch_size)
	//    ]
	//    for num_batch, batch in enumerate(batches):
	//        payload = {
	//            "model": request.model,
	//            "input": batch,
	//            "encoding_format": request.encoding_format,
	//        }
	//        embedding = await get_embedding(payload)
	//        if "error_code" in embedding and embedding["error_code"] != 0:
	//            return create_error_response(embedding["error_code"], embedding["text"])
	//        data += [
	//            {
	//                "object": "embedding",
	//                "embedding": emb,
	//                "index": num_batch * batch_size + i,
	//            }
	//            for i, emb in enumerate(embedding["embedding"])
	//        ]
	//        token_num += embedding["token_num"]
	//    return EmbeddingsResponse(
	//        data=data,
	//        model=request.model,
	//        usage=UsageInfo(
	//            prompt_tokens=token_num,
	//            total_tokens=token_num,
	//            completion_tokens=None,
	//        ),
	//    ).dict(exclude_none=True)
	return
}

func (s *worker) WorkerCountToken(ctx context.Context, workerAddress, model string, prompt any) (res int, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/count_token", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	var resp struct {
		Count     int `json:"count"`
		ErrorCode int `json:"error_code"`
	}
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&resp), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, map[string]any{"model": model, "prompt": prompt})
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return resp.Count, nil
}

func (s *worker) WorkerGetStatus(ctx context.Context, workerAddress string) (res WorkerStatus, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/worker_get_status", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return
}

func (s *worker) WorkerGetModelDetails(ctx context.Context, workerAddress, model string) (res ModelDetail, err error) {
	u, err := url.Parse(fmt.Sprintf("%s/model_details", workerAddress))
	if err != nil {
		err = errors.Wrap(err, "failed to parse url")
		return
	}
	ep := kithttp.NewClient(http.MethodPost, u, kithttp.EncodeJSONRequest, decodeJsonResponse(&res), s.options.httpClientOpts...).Endpoint()
	_, err = ep(ctx, map[string]string{"model": model})
	if err != nil {
		err = errors.Wrap(err, "failed to call endpoint")
		return
	}
	return
}

func NewFastChatWorker(opts ...WorkerCreationOption) WorkerService {
	options := &WorkerCreationOptions{
		controllerAddress: "http://fschat-controller:21001",
	}
	for _, opt := range opts {
		opt(options)
	}
	return &worker{
		options: options,
	}
}

func decodeJsonResponse(data interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error) {
	return func(ctx context.Context, res *http.Response) (response interface{}, err error) {
		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return res, fmt.Errorf("http status code is %d, body %s", res.StatusCode, string(body))
		}
		if data == nil {
			return res, nil
		}
		if err = json.NewDecoder(res.Body).Decode(data); err != nil {
			return res, errors.Wrap(err, "json decode")
		}
		return res, nil
	}
}

// splitAtDelimiter returns a bufio.SplitFunc closure that splits at the given delimiter.
func splitAtDelimiter(delimiter []byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, delimiter); i >= 0 {
			// We have a full delimiter-terminated chunk.
			return i + len(delimiter), data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated chunk. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}

func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Accept", "text/event-stream")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Connection", "keep-alive")
	var b bytes.Buffer
	r.Body = io.NopCloser(&b)
	return json.NewEncoder(&b).Encode(request)
}

func decodeResponse(ctx context.Context, response *http.Response) (interface{}, error) {
	return response.Body, nil
}
