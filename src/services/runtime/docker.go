package runtime

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
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

func (s docker) CreateJob(ctx context.Context, config Config) (jobName string, err error) {
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

	hcf := &container.HostConfig{
		PortBindings: hostPortBindings,
		Resources: container.Resources{
			CPUCount: config.CPU,
			Memory:   config.Memory,
		},
		NetworkMode: "host",
		//Mounts: hostMount,
		Binds: hostBinds,
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

	err = s.dockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		err = fmt.Errorf("ContainerStart err: %s", err.Error())
		return
	}

	return resp.ID, nil
}

func (s docker) getContainerGpuNum(ctx context.Context) (gpuNum []string, err error) {
	containers, err := s.dockerCli.ContainerList(ctx, container.ListOptions{All: true})
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

	fmt.Println("availableGPUs", availableGPUs)

	return availableGPUs, nil
}

func (s docker) CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error) {
	return s.CreateJob(ctx, config)
}

func (s docker) GetDeploymentLogs(ctx context.Context, id string) (log string, err error) {
	out, err := s.dockerCli.ContainerLogs(ctx, id, types.ContainerLogsOptions{
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

func (s docker) GetJobLogs(ctx context.Context, id string) (log string, err error) {
	return s.GetDeploymentLogs(ctx, id)
}

func (s docker) GetJobStatus(ctx context.Context, jobName string) (status string, err error) {
	cJson, err := s.dockerCli.ContainerInspect(ctx, jobName)
	if err != nil {
		err = fmt.Errorf("ContainerInspect err: %s", err.Error())
		return
	}

	return cJson.State.Status, nil
}

func (s docker) GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error) {
	return s.GetJobStatus(ctx, deploymentName)
}

func (s docker) RemoveJob(ctx context.Context, jobName string) (err error) {
	err = s.dockerCli.ContainerStop(ctx, jobName, container.StopOptions{})
	if err != nil {
		err = errors.Wrap(err, "ContainerStop err")
	}

	err = s.dockerCli.ContainerRemove(ctx, jobName, types.ContainerRemoveOptions{})
	return errors.Wrap(err, "ContainerRemove err")
}

func (s docker) RemoveDeployment(ctx context.Context, deploymentName string) (err error) {
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
	return &docker{
		options:   options,
		dockerCli: dockerCli,
	}
}
