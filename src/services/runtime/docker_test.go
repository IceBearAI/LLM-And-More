package runtime

import (
	"context"
	"os"
	"testing"
)

var (
	_ = os.Setenv("DOCKER_HOST", "tcp://10.170.32.94:2375")
)

func TestService_GetContainers(t *testing.T) {
	dockerCli := NewDocker(
		WithWorkspace("/Users/cong/go/src/github.com/icowan/LLM-And-More/storage"),
	)
	containers, err := dockerCli.GetContainers(context.Background(), "qwen1-5-0-5")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(containers)
}
