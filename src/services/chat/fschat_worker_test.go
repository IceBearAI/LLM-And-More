package chat

import (
	"context"
	"encoding/json"
	"testing"
)

func initSvc() WorkerService {
	return NewFastChatWorker(WithControllerAddress("http://localhost:21001"))
}

func TestWorkerService_ListModels(t *testing.T) {
	svc := initSvc()
	models, err := svc.ListModels(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(models)
}

func TestWorkerService_GetWorkerAddress(t *testing.T) {
	svc := initSvc()
	addr, err := svc.GetWorkerAddress(context.Background(), "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
}

func TestWorkerService_GetModelInfo(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	workerAddress, err := svc.GetWorkerAddress(ctx, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	info, err := svc.WorkerGetModelDetails(ctx, workerAddress, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

func TestWorkerService_GetConvTemplate(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	workerAddress, err := svc.GetWorkerAddress(ctx, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	template, err := svc.WorkerGetConvTemplate(ctx, workerAddress, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(template.Conv)
	b, _ := json.Marshal(template.Conv)
	t.Log(string(b))
}

func TestWorkerService_WorkerGetStatus(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	workerAddress, err := svc.GetWorkerAddress(ctx, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	status, err := svc.WorkerGetStatus(ctx, workerAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(status)
}

func TestWorkerService_WorkerGenerate(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	workerAddress, err := svc.GetWorkerAddress(ctx, "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	prompt := "<|im_start|>system\nYou are a helpful assistant.<|im_end|><|im_start|>user\n您好<|im_end|><|im_start|>assistant\n"
	params := GenerateParams{
		Model:        "qwen1.5-0.5b",
		Prompt:       prompt,
		Stop:         []string{"<|endoftext|>"},
		Temperature:  0.7,
		TopP:         0.6,
		MaxNewTokens: 1024,
		StopTokenIds: []int{151643,
			151644,
			151645},
	}
	resp, err := svc.WorkerGenerate(ctx, workerAddress, params)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		select {
		case r, ok := <-resp:
			if !ok {
				return
			}
			t.Log(r.Text)
		}
	}
}

func TestWorkerService_WorkerGenerateStream(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	workerAddress, err := svc.GetWorkerAddress(ctx, "qwen1.5-0.5b-chat")
	if err != nil {
		t.Error(err)
		return
	}
	prompt := "<|im_start|>system\nYou are a helpful assistant.<|im_end|><|im_start|>user\n您好!你叫什么名字？<|im_end|><|im_start|>assistant\n"
	params := GenerateStreamParams{
		Model:        "qwen1.5-0.5b-chat",
		Prompt:       prompt,
		Stop:         []string{"<|endoftext|>"},
		Temperature:  0.7,
		TopP:         0.6,
		MaxNewTokens: 1024,
		StopTokenIds: []int{151643,
			151644,
			151645},
	}
	resp, err := svc.WorkerGenerateStream(ctx, workerAddress, params)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		select {
		case r, ok := <-resp:
			if !ok {
				return
			}
			t.Log(r.Text)
		}
	}
}
