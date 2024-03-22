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
