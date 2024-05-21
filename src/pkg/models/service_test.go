package models

import (
	"context"
	"github.com/IceBearAI/aigc/tests"
	"github.com/go-kit/log"
	"os"
	"testing"
)

func initSvc() Service {
	_ = os.Setenv("AIGC_RUNTIME_K8S_CONFIG_PATH", "github.com/IceBearAI/aigc/k8sconfig.yaml")
	_ = os.Setenv("AIGC_RUNTIME_PLATFORM", "docker")
	_ = os.Setenv("AIGC_STORAGE_TYPE", "local")
	_ = os.Setenv("AIGC_RUNTIME_K8S_VOLUME_NAME", "aigc-data-cfs")
	_ = os.Setenv("AIGC_ADMIN_SERVER_STORAGE_PATH", "/Users/dudu/go/src/github.com/icowan/LLM-And-More/storage")

	apiSvc, err := tests.Init()
	if err != nil {
		panic(err)
	}
	return NewService(log.NewLogfmtLogger(os.Stdout), "", tests.Store, apiSvc)
}

func TestService_Deploy(t *testing.T) {
	ctx := context.Background()
	err := initSvc().Deploy(ctx, ModelDeployRequest{
		Id:           95,
		Replicas:     1,
		Label:        "aigc-cpu-model",
		InferredType: "cpu",
		Gpu:          1,
		Cpu:          0,
		Quantization: "",
		Vllm:         false,
		MaxGpuMemory: 0,
		K8sCluster:   "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success.")
}

func TestService_ModelTree(t *testing.T) {
	ctx := context.Background()
	res, err := initSvc().ModelTree(ctx, "chatglm3-6b", "config.json")
	if err != nil {
		t.Error(err)
		return
	}
	if res.Object == "tree" {
		for _, v := range res.Tree {
			t.Log(v.Name, v.IsDir, v.Size, v.UpdatedAt)
		}
	}

	if res.Object == "file" {
		t.Log(res.FileInfo)
		t.Log(res.FileContent)
	}
}

func TestService_ModelCard(t *testing.T) {
	ctx := context.Background()
	res, err := initSvc().ModelCard(ctx, "chatglm3-6b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.ReadmeContent)
}
