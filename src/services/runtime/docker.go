package runtime

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
)

// WithWorkspace returns a CreationOption that sets the workspace.
func WithWorkspace(workspace string) CreationOption {
	return func(co *CreationOptions) {
		co.workspace = workspace
	}
}

// WithGpuNum returns a CreationOption that sets the workspace.
func WithGpuNum(gpuNum int) CreationOption {
	return func(co *CreationOptions) {
		co.gpuNum = gpuNum
	}
}

type docker struct {
	options       *CreationOptions
	dockerCli     *client.Client
	createOptions CreationOptions
}

func (s *docker) WaitForTerminal(ctx context.Context, ts Session, config Config, container, cmd string) {
	execConfig := types.ExecConfig{
		Cmd:          []string{cmd},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
	}
	cec, err := s.dockerCli.ContainerExecCreate(ctx, container, execConfig)
	if err != nil {
		log.Println("ContainerExecCreate err: ", err)
		return
	}

	if cec.ID == "" {
		log.Println("ContainerExecCreate ID is empty")
		return
	}

	hijack, err := s.dockerCli.ContainerExecAttach(ctx, cec.ID, types.ExecStartCheck{})
	if err != nil {
		log.Println("ContainerExecAttach err: ", err)
		return
	}

	defer hijack.Close()

	buf1 := make([]byte, 1024)
	go func() {
		defer hijack.CloseWrite()
		for {
			var rint int
			rint, err = ts.Read(buf1)
			if err != nil {
				log.Println("Read err: ", err)
				if err == io.EOF {
					return
				}
				log.Println("Read err: ", err)
				break
			}

			if _, err = hijack.Conn.Write(buf1[:rint]); err != nil {
				log.Println("Write err: ", err)
				break
			}
		}
	}()

	buf2 := make([]byte, 1024)
	for {
		var rint int
		rint, err = hijack.Reader.Read(buf2)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Read err: ", err)
			break
		}
		_, err = ts.Write(buf2[:rint])
		if err != nil {
			log.Println("Write err: ", err)
			break
		}
	}
}

func (s *docker) GetDeploymentContainerNames(ctx context.Context, deploymentName string) (containerNames []string, err error) {
	newName := replacerServiceName(deploymentName)
	list, err := s.dockerCli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "ContainerList err")
		return
	}

	for _, v := range list {
		if strings.HasPrefix(v.Names[0], newName+"-") {
			containerNames = append(containerNames, v.ID)
		}
	}

	return
}

func (s *docker) HandleTerminalSession(ctx context.Context, config Config) {
	// 连接到Docker container terminal
	// Connect to Docker container
	// 这里的代码是一个示例，实际上需要根据具体的业务逻辑来实现
	// The code here is an example, and the actual implementation needs to be based on specific business logic
	// 创建执行实例
	execConfig := types.ExecConfig{
		Cmd:          []string{"/bin/sh"},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
	}
	execIDResp, err := s.dockerCli.ContainerExecCreate(ctx, config.ServiceName, execConfig)
	if err != nil {
		err = errors.Wrap(err, "ContainerExecCreate err")
		return
	}

	attach, err := s.dockerCli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{})
	if err != nil {
		err = errors.Wrap(err, "ContainerExecAttach err")
		return
	}
	defer attach.Close()
	// 将容器的输出复制到标准输出
	go func() {
		if execConfig.Tty {
			_, _ = io.Copy(os.Stdout, attach.Reader)
		} else {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader)
			if err != nil {
				log.Println("stdcopy.StdCopy err: ", err)
			}
		}
	}()

	// 发送命令到容器
	_, err = attach.Conn.Write([]byte("ls\n"))
	if err != nil {
		panic(err)
	}

	// 等待命令执行完毕（这里简单地等待几秒）
	// 实际应用中应该使用更合适的同步机制
	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
	}

	// 读取执行实例的退出代码
	inspectResp, err := s.dockerCli.ContainerExecInspect(ctx, execIDResp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Exit Code: %d\n", inspectResp.ExitCode)
}

func (s *docker) CreateJob(ctx context.Context, config Config) (jobName string, err error) {
	_ = s.RemoveJob(ctx, config.ServiceName)

	exposedPorts := make(nat.PortSet)
	hostPortBindings := make(nat.PortMap)
	hostBinds := make([]string, 0)

	if len(config.ConfigData) > 0 {
		err = config.saveConfigToLocal(config.ServiceName, s.options.workspace)
		if err != nil {
			err = fmt.Errorf("config.saveConfigToLocal err: %s", err.Error())
			return
		}
	}

	for k, v := range config.Ports {
		exposedPort := fmt.Sprintf("%s/tcp", v)
		exposedPorts[nat.Port(exposedPort)] = struct{}{}
		hostPortBindings[nat.Port(exposedPort)] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: k,
			},
		}
	}

	for _, v := range config.Volumes {
		if config.HasConfigData(v.Key) {
			hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", filepath.Join(s.options.workspace, config.ServiceName, v.Key), v.Value))
		} else {
			hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", filepath.Join(s.options.workspace, v.Key), v.Value))
		}
	}
	//if s.options.k8sVolumeName != "" {
	//	hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", filepath.Join(s.options.workspace, s.options.k8sVolumeName), "/data"))
	//}

	for k, _ := range config.ConfigData {
		hostBinds = append(hostBinds, fmt.Sprintf("%s:%s", filepath.Join(s.options.workspace, config.ServiceName, k), k))
	}

	ccf := &container.Config{
		Image:        config.Image,
		Cmd:          config.Command,
		Env:          append([]string{"SERVICE_NAME=" + config.ServiceName}, config.EnvVars...),
		ExposedPorts: exposedPorts,
	}

	// s.createOptions.shmSizes.createOptions.shmSize
	defaultShmSize := s.options.shmSize
	if defaultShmSize == "" {
		defaultShmSize = "4G"
	}
	parse := resource.MustParse(defaultShmSize)
	shmSize := parse.Value()

	hcf := &container.HostConfig{
		PortBindings: hostPortBindings,
		Resources: container.Resources{
			CPUCount: config.CPU,
			Memory:   config.Memory,
		},
		NetworkMode: "host",
		//Mounts: hostMount,
		Binds:   hostBinds,
		ShmSize: shmSize,
	}
	var dr []container.DeviceRequest
	if config.GPU != 0 {
		availableGPUs, err := s.getContainerGpuNum(ctx)
		if err != nil {
			return "", errors.Wrap(err, "getContainerGpuNum err")
		}
		if len(availableGPUs) < config.GPU {
			err = errors.New("No enough GPU")
			return "", err
		}
		dr = append(dr, container.DeviceRequest{
			Driver: "nvidia",
			//Count:        -1,
			Capabilities: [][]string{{"gpu"}},
			DeviceIDs:    availableGPUs[:config.GPU],
		})

	}
	hcf.Resources.DeviceRequests = dr

	resp, err := s.dockerCli.ContainerCreate(ctx, ccf, hcf, nil, nil, config.ServiceName)
	if err != nil {
		err = fmt.Errorf("ContainerCreate err: %s", err.Error())
		return
	}

	err = s.dockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		err = fmt.Errorf("ContainerStart err: %s", err.Error())
		return
	}

	return config.ServiceName, nil
}

func (s *docker) getContainerGpuNum(ctx context.Context) (gpuNum []string, err error) {
	containers, err := s.dockerCli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "ContainerList err")
		return
	}
	usedGPUs := make(map[string]bool)
	for _, _container := range containers {
		containerJSON, err := s.dockerCli.ContainerInspect(ctx, _container.ID)
		if err != nil {
			fmt.Printf("Error inspecting container %s: %s\n", _container.ID, err)
			continue
		}

		for _, deviceRequest := range containerJSON.HostConfig.DeviceRequests {
			if deviceRequest.Driver == "nvidia" {
				for _, id := range deviceRequest.DeviceIDs {
					usedGPUs[id] = true
					fmt.Printf("Container %s is using GPU %s\n", _container.ID[:12], id)
				}
			}
		}
	}

	// Assuming you have a maximum of 8 GPUs
	totalGPUs := s.options.gpuNum
	availableGPUs := []string{}
	for i := 0; i < totalGPUs; i++ {
		if !usedGPUs[fmt.Sprint(i)] {
			availableGPUs = append(availableGPUs, fmt.Sprint(i))
		}
	}

	return availableGPUs, nil
}

func replacerServiceName(name string) string {
	replacer := strings.NewReplacer(
		"_", "-",
		".", "-",
		"::", "-", // 这个可能不需要，因为前一个已经将单个冒号替换了
		":", "-",
	)
	return strings.ToLower(replacer.Replace(name))
}

func (s *docker) CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error) {
	serviceName := config.ServiceName
	for i := 0; i < int(config.Replicas); i++ {
		config.ServiceName = fmt.Sprintf("%s-%d", replacerServiceName(config.ServiceName), i)
		if _, err = s.CreateJob(ctx, config); err != nil {
			err = fmt.Errorf("createJob err: %s, job name: %s", err.Error(), config.ServiceName)
			return serviceName, err
		}
	}
	return serviceName, nil
}

func (s *docker) GetDeploymentLogs(ctx context.Context, id, containerName string) (log string, err error) {
	out, err := s.dockerCli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		err = fmt.Errorf("ContainerLogs err: %s", err.Error())
		return
	}

	b, err := io.ReadAll(out)
	if err != nil {
		err = fmt.Errorf("io.ReadAll err: %s", err.Error())
		return
	}
	return string(b), nil
}

func (s *docker) GetJobLogs(ctx context.Context, id string) (log string, err error) {
	return s.GetDeploymentLogs(ctx, id, "")
}

func (s *docker) GetJobStatus(ctx context.Context, jobName string) (status string, err error) {
	cJson, err := s.dockerCli.ContainerInspect(ctx, jobName)
	if err != nil {
		err = fmt.Errorf("ContainerInspect err: %s", err.Error())
		return
	}

	switch cJson.State.Status {
	case "created", "restarting":
		status = "Pending"
	case "running", "paused":
		status = "Running"
	case "exited":
		status = "Failed"
	case "dead":
		status = "Failed"
	default:
		status = "Unknown"
	}
	return status, nil
}

func (s *docker) GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error) {
	return s.GetJobStatus(ctx, deploymentName)
}

func (s *docker) RemoveJob(ctx context.Context, jobName string) (err error) {
	if err = s.dockerCli.ContainerKill(ctx, jobName, "SIGKILL"); err != nil {
		err = errors.Wrap(err, "ContainerStop err")
	}

	err = s.dockerCli.ContainerRemove(ctx, jobName, container.RemoveOptions{})
	return errors.Wrap(err, "ContainerRemove err")
}

func (s *docker) RemoveDeployment(ctx context.Context, deploymentName string) (err error) {
	deploymentName = replacerServiceName(deploymentName)
	return s.RemoveJob(ctx, deploymentName)
}

func NewDocker(opts ...CreationOption) Service {
	options := &CreationOptions{
		workspace: "/tmp",
		gpuNum:    8,
	}
	for _, opt := range opts {
		opt(options)
	}
	dockerCli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	//images, _ := dockerCli.ImageSearch(context.Background(), "dudulu/fschat:latest", types.ImageSearchOptions{})
	//if len(images) == 0 {
	//	_, err := dockerCli.ImagePull(context.Background(), "dudulu/fschat:latest", types.ImagePullOptions{})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	return &docker{
		options:   options,
		dockerCli: dockerCli,
	}
}
