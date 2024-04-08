package runtime

import (
	"context"
	"testing"
)

// var s Service = NewDocker()
// var id string = ""

// var s, err = NewK8s(WithK8sConfigPath("./k8sconfig.yaml"), WithNamespace("dev"))
// var token string = `rQ-a28Qz1HxPZYMxYUjB51zqAK8M3hs8hvIKbh30r3Z2FpRjIBboFZaoRTClp9UirJj_SoYs4XOKBcxoDmJAxPuvQXD4QPyR8TrCswIAHeP5DOtovFg_9HgN_0wGRROzmSg6VKR096PljPB0YqOMulZPMyS52qKE8PHy8IA6ggf_CSzzwEesv4Zs9002zf8TOJAH6ZmJyyVut2i1zgg8mnb6eSN1Oe8nHl210bukpaIz2N1l1b5vEzHb3jE-NjKw5Q9EoacL0t-_pFEMsBjNMyMQuXlshb9KIDeRqcEmke9SCqY1I1EKyHXsaTHs4qMD87QviT9v_Ffz8X-DM7xDNw`
// var s, err = NewK8s(WithK8sToken("https://127.0.0.1:6443", token, true), WithNamespace("default"))
// var id string = "ddwded2"
var id string = "74ef6a371ef98d3c8bac5a54b88031db9d792f5f3a205dd8d8c6ed29b3b127a1"
var image = "nginx:latest"

var s = NewDocker(WithWorkspace("/Users/cong/go/src/github.com/icowan/LLM-And-More/storage"))

func TestService_CreateDeployment(t *testing.T) {

	//_ = os.Setenv("DOCKER_HOST", "tcp://127.0.0.1:2375")

	ctx := context.Background()

	jobName, err := s.CreateDeployment(ctx, Config{
		ServiceName: id,
		Image:       image,
		Cpu:         1,
		GPU:         0,
		// Command: []string{
		// "/bin/bash",
		// "/app/dataset_analyze_similar.sh",
		// },
		Volumes: []Volume{
			{
				Key:   "/Users/cong/go/src/github.com/icowan/LLM-And-More/storage",
				Value: "/data/",
			},
		},
		Command: []string{
			"/bin/bash",
			"-c",
			"echo start sleep ;sleep 10000000;",
		},
		EnvVars: []string{
			"TZ=Asia/Shanghai",
		},
		ConfigData: map[string]string{
			"/app/dataset.json": "{\"instruction\":\"asdfasdfasdf\",\"input\":\"\",\"output\":\"ffff\",\"intent\":\"ffff\",\"document\":\"\",\"question\":\"fadgawevansdf\"}\n",
			"/app/hello":        "say hello",
		},
		Replicas: 1,
		ShmSize:  "24MB",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(jobName)
	id = jobName
}

func TestService_CreateJob(t *testing.T) {
	jobName, err := s.CreateJob(context.Background(), Config{
		ServiceName: id,
		Image:       image,
		Cpu:         1,
		Memory:      1,
		GPU:         0,
		// Command: []string{
		// "/bin/bash",
		// "/app/dataset_analyze_similar.sh",
		// },
		Command: []string{
			"/bin/bash",
			"-c",
			"echo start sleep ;sleep 10000000;",
		},
		EnvVars: []string{
			"TZ=Asia/Shanghai",
		},
		ConfigData: map[string]string{
			"/app/dataset.json": "{\"instruction\":\"asdfasdfasdf\",\"input\":\"\",\"output\":\"ffff\",\"intent\":\"ffff\",\"document\":\"\",\"question\":\"fadgawevansdf\"}\n",
			"/app/hello":        "say hello",
		},
		Replicas: 1,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(jobName)
	id = jobName
}

func TestService_GetDeploymentStatus(t *testing.T) {
	status, err := s.GetDeploymentStatus(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(status)
}

func TestService_GetDeploymentLog(t *testing.T) {
	log, err := s.GetDeploymentLogs(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(log)
}

func TestService_GetJobStatus(t *testing.T) {
	status, err := s.GetJobStatus(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(status)
}

func TestService_GetJobLog(t *testing.T) {
	log, err := s.GetJobLogs(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(log)
}

func TestService_RemoveDeployment(t *testing.T) {
	err := s.RemoveDeployment(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestService_RemoveJob(t *testing.T) {
	err := s.RemoveJob(context.Background(), id)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
