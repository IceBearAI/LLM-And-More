package datasets

import (
	"context"
	"github.com/IceBearAI/aigc/tests"
	"testing"
)

func initSvc() Service {
	_, err := tests.Init()
	if err != nil {
		return nil
	}
	return New(tests.Logger, "", tests.Store)
}

func TestService_Create(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	datasetId, err := svc.Create(ctx, 1, "test111", "test222")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("datasetId:", datasetId)
}

func TestService_Update(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	err := svc.Update(ctx, 1, "2fa7f17e-8b58-454f-9063-fd8a919ff7ff", "test-aaa", "testasdhfasdfasd")
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Delete(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	err := svc.Delete(ctx, 1, "2fa7f17e-8b58-454f-9063-fd8a919ff7ff")
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Detail(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	dataset, err := svc.Detail(ctx, 1, "88c7349b-1b31-480b-a17e-c19760fe836a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dataset)
}

func TestService_List(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	datasets, total, err := svc.List(ctx, 1, 1, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("total:", total)
	t.Log(datasets)
}

func TestService_CreateSample(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	err := svc.AddSample(ctx, 1, "74a32b55-e08e-4834-bf2e-89e1c1f7da6d", "text", []message{
		{
			Role:    "user",
			Content: "您好",
		}, {
			Role:    "user",
			Content: "您好！我是智能机器人！",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_UpdateSample(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	err := svc.UpdateSampleMessages(ctx, 1, "74a32b55-e08e-4834-bf2e-89e1c1f7da6d", "102f7885-3f72-409b-93b8-e8ec625fa90f", []message{
		{
			Role:    "user",
			Content: "您好",
		}, {
			Role:    "assistant",
			Content: "您好！我是智能机器人！",
		},
		{
			Role:    "user",
			Content: "你是谁训练的？",
		},
		{
			Role:    "user",
			Content: "我是你训练的？",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_DeleteSample(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	err := svc.DeleteSample(ctx, 1, "74a32b55-e08e-4834-bf2e-89e1c1f7da6d", []string{"13af6446-77a3-4654-a606-1237ea47c8d0"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_SampleList(t *testing.T) {
	svc := initSvc()
	if svc == nil {
		t.Fatal("init service failed")
	}
	ctx := context.Background()
	samples, total, err := svc.SampleList(ctx, 1, "74a32b55-e08e-4834-bf2e-89e1c1f7da6d", 1, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("total:", total)
	t.Log(samples)
}
