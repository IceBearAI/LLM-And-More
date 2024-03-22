package tools

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/tests"
	"testing"
)

func initSvc() Service {
	_, err := tests.Init()
	if err != nil {
		panic(err)
	}
	return New(tests.Logger, "traceId", tests.Store)
}

func TestService_Create(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	err := svc.Create(ctx, 1, createRequest{
		Name: "get_current_time2",
		Description: `对于获取当前时间很有帮助。
	无需任何参数，直接返回当前时间。`,
		ToolType: string(types.ToolTypeTypeFunction),
		Metadata: `{"url": "https://f.m.suning.com/api/ct.do", "method": "GET", "headers": }, "body": "", "userAgent": "aigc-server"}`,
	})
	if err != nil {
		t.Error(err)
		return
	}
	return
}

func TestService_Delete(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	err := svc.Delete(ctx, 1, "tool-31e7d2b6-8668-4b29-866b-cf39f96346b4")
	if err != nil {
		t.Error(err)
		return
	}
	return
}
