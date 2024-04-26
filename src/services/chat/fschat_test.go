package chat

import (
	"context"
	"testing"
)

func initSvc() Service {
	return NewFsChat(WithControllerAddress("http://fschat-controller.paas.paas.test"))
}

func TestFschatService_ListModels(t *testing.T) {
	svc := initSvc()
	models, err := svc.ListModels(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(models)
}

func TestFschatService_GetWorkerAddress(t *testing.T) {
	svc := initSvc()
	addr, err := svc.GetWorkerAddress(context.Background(), "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
}

func TestFschatService_CheckModel(t *testing.T) {
	svc := initSvc()
	res, err := svc.CheckModel(context.Background(), "qwen1.5-0.5b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestFschatService_GetModelInfo(t *testing.T) {
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

func TestFschatService_GetConvTemplate(t *testing.T) {
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
	t.Log(template)
}

func TestFschatService_WorkerGetStatus(t *testing.T) {
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
