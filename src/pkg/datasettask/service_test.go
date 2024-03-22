package datasettask

import (
	"context"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/tests"
	kithttp "github.com/go-kit/kit/transport/http"
	"os"
	"testing"
)

func initSvc() Service {
	_ = os.Setenv("AIGC_RUNTIME_K8S_CONFIG_PATH", "./k8sconfig.yaml")
	//_ = os.Setenv("AIGC_RUNTIME_PLATFORM", "docker")
	_ = os.Setenv("AIGC_RUNTIME_PLATFORM", "k8s")
	_ = os.Setenv("AIGC_STORAGE_TYPE", "local")
	_ = os.Setenv("AIGC_RUNTIME_K8S_VOLUME_NAME", "aigc-data-cfs")

	//_ = os.Setenv("DOCKER_HOST", "tcp://127.0.0.1:2376")
	apiSvc, err := tests.Init()
	if apiSvc == nil || err != nil {
		panic(err)
	}
	fileOptions := []files.CreationOption{
		files.WithLocalDataPath("/data/storage"),
		files.WithServerUrl("http://localhost:8080/storage"),
	}
	fileSvc := files.NewService(tests.Logger, "", tests.Store, apiSvc, fileOptions...)
	return New("", tests.Logger, tests.Store, apiSvc, fileSvc,
		WithDatasetImage("dudulu/llmops:v0.8-0314"),
		WithDatasetModel("uer/sbert-base-chinese-nli"),
		WithCallbackHost("http://localhost:8080"),
		WithDatasetGpuTolerationValue("cpu-aigc-model"),
	)
}

func TestService_AsyncCheckTaskDatasetSimilar(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJuZCIsInF3VXNlcmlkIjoiIiwiZW1haWwiOiJjb25nd2FuZ0BjcmVkaXRlYXNlLmNuIiwidXNlcklkIjozLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoic3lzdGVtIiwiZXhwIjoxNzEwOTIyOTI4fQ.C9jox2_KozIHlTKclbeoeyiXQsYErGL49YQ7GhkL5MY")
	err := svc.AsyncCheckTaskDatasetSimilar(ctx, 1, "task-e510e4f6-d755-4b2e-a32e-4d8614d3a71f")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success.")
}

func TestService_CancelCheckTaskDatasetSimilar(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	err := svc.CancelCheckTaskDatasetSimilar(ctx, 1, "task-e510e4f6-d755-4b2e-a32e-4d8614d3a71f")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success.")
}
