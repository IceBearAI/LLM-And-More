package modelevaluate

import (
	"context"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/tests"
	kithttp "github.com/go-kit/kit/transport/http"
	"os"
	"testing"
)

func initSvc() Service {
	_ = os.Setenv("AIGC_RUNTIME_K8S_CONFIG_PATH", "IceBearAI/aigc/k8sconfig.yaml")
	_ = os.Setenv("AIGC_RUNTIME_PLATFORM", "docker")
	_ = os.Setenv("AIGC_STORAGE_TYPE", "local")
	_ = os.Setenv("AIGC_RUNTIME_K8S_VOLUME_NAME", "aigc-data-cfs")

	//_ = os.Setenv("DOCKER_HOST", "tcp://10.143.151.50:2376")
	apiSvc, _ := tests.Init()
	if apiSvc == nil {
		panic(apiSvc)
	}
	fileOptions := []files.CreationOption{
		files.WithLocalDataPath("/data/storage"),
		files.WithServerUrl("http://localhost:8080/storage"),
	}
	fileSvc := files.NewService(tests.Logger, "", tests.Store, apiSvc, fileOptions...)
	return New(tests.Logger, "", tests.Store, apiSvc, fileSvc,
		WithCallbackHost("http://localhost:8080"),
		WithDatasetGpuTolerationValue("cpu-aigc-model"),
	)
}

func TestService_Create(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJuZCIsInF3VXNlcmlkIjoiIiwiZW1haWwiOiJjb25nd2FuZ0BjcmVkaXRlYXNlLmNuIiwidXNlcklkIjozLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoic3lzdGVtIiwiZXhwIjoxNzEwOTIyOTI4fQ.C9jox2_KozIHlTKclbeoeyiXQsYErGL49YQ7GhkL5MY")
	err := svc.Create(ctx, createRequest{
		FileId:         "e05edda4-9690-4f7f-a1c3-5d855ec8666d",
		ModelId:        95,
		EvalTargetType: "ACC",
		MaxLength:      512,
		BatchSize:      32,
		Replicas:       1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success.")
}

func TestService_Delete(t *testing.T) {
	svc := initSvc()
	ctx := context.Background()
	ctx = context.WithValue(ctx, kithttp.ContextKeyRequestAuthorization, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJuZCIsInF3VXNlcmlkIjoiIiwiZW1haWwiOiJjb25nd2FuZ0BjcmVkaXRlYXNlLmNuIiwidXNlcklkIjozLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoic3lzdGVtIiwiZXhwIjoxNzEwOTIyOTI4fQ.C9jox2_KozIHlTKclbeoeyiXQsYErGL49YQ7GhkL5MY")
	err := svc.Delete(ctx, deleteRequest{Uuid: ""})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success.")
}
