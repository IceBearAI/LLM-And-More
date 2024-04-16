package runtime

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	endpoint       string
	workspace      string
	gpuNum         int
	shmSize        string
	namespace      string
	k8sConfigPath  string
	k8sTokenModel  k8sTokenModel
	k8sVolumeName  string
	labelName      string
}

type k8sTokenModel struct {
	host     string
	token    string
	insecure bool
}

// CreationOption is a creation option for the runtime service.
type CreationOption func(*CreationOptions)

type Volume struct {
	// 宿主路径
	Key string
	// 容器路径
	Value string
}

type Env struct {
	Name  string
	Value string
}

type Config struct {
	namespace          string
	ServiceName        string
	Image              string
	Cpu                int64
	Memory             int64
	Command            []string
	EnvVars            []string
	Volumes            []Volume
	GpuTolerationValue string
	// k: 宿主端口, v: 容器端口
	Ports map[string]string
	CPU   int64
	GPU   int
	// k: 文件路径, v: 文件内容
	ConfigData map[string]string

	// 支队k8s的deployment有效
	Replicas int32

	ShmSize string

	restartPolicy v1.RestartPolicy
	LabelName     string
}

func (c Config) FilePath2Key(filePath string) string {
	return strings.Trim(strings.ReplaceAll(filePath, "/", "-"), "-")
}

func (c Config) GenConfigMap() (res v1.ConfigMap) {
	data := make(map[string]string, 0)
	for k, v := range c.ConfigData {
		data[c.FilePath2Key(k)] = v
	}
	res = v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      c.ServiceName,
		},
		Data: data,
	}
	return
}

func (c Config) GenVolumeAndVolumeMount() (volumes []v1.Volume, volumeMounts []v1.VolumeMount) {
	volumeCurrent := map[string]bool{}
	for _, v := range c.Volumes {
		// volumeName := c.FilePath2Key(v.Key)
		if _, ok := volumeCurrent[v.Key]; !ok {
			volumes = append(volumes, v1.Volume{
				Name: v.Key,
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: v.Key,
					},
				},
			})
		}
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      v.Key,
			MountPath: v.Value,
		})
		volumeCurrent[v.Key] = true
	}

	if len(c.ConfigData) > 0 {
		var items []v1.KeyToPath
		for k, _ := range c.ConfigData {
			_, fileName := filepath.Split(k)
			items = append(items, v1.KeyToPath{
				Key:  c.FilePath2Key(k),
				Path: fileName,
			})

			volumeMounts = append(volumeMounts, v1.VolumeMount{
				Name:      c.ServiceName,
				MountPath: k,
				SubPath:   fileName,
			})
		}

		volume := v1.Volume{
			Name: c.ServiceName,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: c.ServiceName,
					},
					Items: items,
				},
			},
		}

		volumes = append(volumes, volume)
	}

	if c.ShmSize != "" {
		sizeList, _ := resource.ParseQuantity(c.ShmSize)

		volumes = append(volumes, v1.Volume{
			Name: "dshm",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium:    "Memory",
					SizeLimit: &sizeList,
				},
			},
		})
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      "dshm",
			MountPath: "/dev/shm",
		})
	}

	return

}

func (c Config) GenDeploymentLabels() map[string]string {
	return map[string]string{
		c.LabelName: c.ServiceName,
	}
}

func (c Config) GenJobLabels() map[string]string {
	return map[string]string{
		c.LabelName: c.ServiceName,
	}
}

func (c Config) GenJob() (res batchv1.Job, err error) {
	spec, err := c.GenJobSpec()
	if err != nil {
		err = errors.Wrap(err, "c.GenJobSpec")
		return
	}

	res = batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      c.ServiceName,
			Labels:    c.GenJobLabels(),
		},
		Spec: spec,
	}

	return
}

func (c Config) GenJobSpec() (res batchv1.JobSpec, err error) {
	spec, err := c.GenPodTemplateSpec()
	if err != nil {
		err = errors.Wrap(err, "pod template spec error")
		return
	}
	return batchv1.JobSpec{
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: c.GenJobLabels(),
			},
			Spec: spec,
		},
	}, nil
}

func (c Config) GenDeployment() (res appsv1.Deployment, err error) {
	spec, err := c.GenDeploymentSpec()
	if err != nil {
		err = errors.Wrap(err, "c.GenDeploymentSpec")
		return
	}

	res = appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      c.ServiceName,
			Labels:    c.GenDeploymentLabels(),
		},
		Spec: spec,
	}

	return
}

func (c Config) GenDeploymentSpec() (res appsv1.DeploymentSpec, err error) {
	spec, err := c.GenPodTemplateSpec()
	if err != nil {
		err = errors.Wrap(err, "c.GenPodTemplateSpec")
		return
	}

	res = appsv1.DeploymentSpec{
		Replicas: &c.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: c.GenDeploymentLabels(),
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: c.GenDeploymentLabels(),
			},
			Spec: spec,
		},
	}

	return
}

func (c Config) GenPodTemplateSpec() (res v1.PodSpec, err error) {
	var volumes []v1.Volume
	var containers []v1.Container
	var tolerations []v1.Toleration

	volumes, _ = c.GenVolumeAndVolumeMount()
	containers, err = c.GenContainers()
	if err != nil {
		err = errors.Wrap(err, "c.GenContainers")
	}

	if c.GpuTolerationValue != "" {
		tolerations = append(tolerations, v1.Toleration{
			Key:      "nvidia.com/gpu",
			Operator: v1.TolerationOpEqual,
			Value:    c.GpuTolerationValue,
			Effect:   v1.TaintEffectNoSchedule,
		})
	}

	res = v1.PodSpec{
		Tolerations:   tolerations,
		Volumes:       volumes,
		Containers:    containers,
		RestartPolicy: c.restartPolicy,
		DNSPolicy:     v1.DNSClusterFirst,
	}
	return
}

func (c Config) GenContainers() (res []v1.Container, err error) {
	var volumeMounts []v1.VolumeMount
	var envs []v1.EnvVar
	var ports []v1.ContainerPort
	var resources = make(v1.ResourceList, 0)
	_, volumeMounts = c.GenVolumeAndVolumeMount()

	for k, v := range c.Ports {
		var containerPort, hostPort int64
		containerPort, err = strconv.ParseInt(v, 10, 32)
		if err != nil {
			err = fmt.Errorf("port %s strconv.ParseInt err: %s", k, err.Error())
			return
		}

		hostPort, err = strconv.ParseInt(k, 10, 32)
		if err != nil {
			err = fmt.Errorf("port %s strconv.ParseInt err: %s", v, err.Error())
			return
		}

		ports = append(ports, v1.ContainerPort{
			Name:          fmt.Sprintf("port-%s", k),
			ContainerPort: int32(containerPort),
			HostPort:      int32(hostPort),
			Protocol:      v1.ProtocolTCP,
		})
	}

	for _, v := range c.EnvVars {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			err = fmt.Errorf("env %s split err", v)
			return
		}
		envs = append(envs, v1.EnvVar{
			Name:  kv[0],
			Value: kv[1],
		})
	}

	if c.Cpu != 0 {
		resources[v1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%d", c.Cpu))
	}

	if c.Memory != 0 {
		resources[v1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%dGi", c.Memory))
	}

	if c.GPU != 0 {
		resources["nvidia.com/gpu"] = resource.MustParse(strconv.Itoa(c.GPU))
	}

	res = append(res, v1.Container{
		Name:         c.ServiceName,
		Image:        c.Image,
		Command:      c.Command,
		Env:          envs,
		VolumeMounts: volumeMounts,
		Ports:        ports,
		Resources:    v1.ResourceRequirements{Limits: resources, Requests: resources},
	})
	return
}

func (c Config) HasConfigData(s string) bool {
	_, ok := c.ConfigData[s]
	return ok
}

func (c Config) Tar(name, workspace string) (*bytes.Buffer, error) {
	err := c.saveConfigToLocal(name, workspace)
	if err != nil {
		return nil, err
	}

	return c.tarDirectory(name, workspace)
}

func (c Config) saveConfigToLocal(name, workspace string) (err error) {
	dir := filepath.Join(workspace, name)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = errors.Wrap(err, "os.MkdirAll")
		return
	}

	for k, v := range c.ConfigData {
		var f *os.File
		dataFilePath := filepath.Join(dir, k)
		mkdir := filepath.Dir(dataFilePath)
		err = os.MkdirAll(mkdir, 0755)
		if err != nil {
			err = errors.Wrap(err, "os.MkdirAll")
			return
		}
		f, err = os.OpenFile(dataFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			err = fmt.Errorf("%s openFile err: %s", dataFilePath, err.Error())
			return
		}
		_, err = f.Write([]byte(v))
		if err != nil {
			err = fmt.Errorf("%s write err: %s", dataFilePath, err.Error())
			return
		}

		err = f.Close()
		if err != nil {
			err = fmt.Errorf("%s close err: %s", dataFilePath, err.Error())
			return
		}
	}

	return
}

func (c Config) tarDirectory(name, workspace string) (*bytes.Buffer, error) {
	source := filepath.Join(workspace, name)
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	defer func() {
		_ = tw.Close()
	}()

	// 递归地添加目录中的文件和子目录到 tar 归档
	err := filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建 tar 头部
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// 必要时调整头部信息
		header.Name = strings.TrimPrefix(file, source)

		// 写入头部信息
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// 如果不是普通文件则不继续
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// 将文件内容复制到 tar 归档
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 确保 tar 归档写入完成
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buffer, nil
}

// Service is a service interface
//
//go:generate gowrap gen -g -p ./ -i Service -bt "ce_log:logging.go ce_trace:tracing.go"
type Service interface {
	// CreateJob 创建job
	CreateJob(ctx context.Context, config Config) (jobName string, err error)
	// CreateDeployment 创建deployment
	CreateDeployment(ctx context.Context, config Config) (deploymentName string, err error)
	// GetDeploymentLogs 获取部署的日志
	GetDeploymentLogs(ctx context.Context, deploymentName, containerName string) (log string, err error)
	// GetJobLogs 获取job的日志
	GetJobLogs(ctx context.Context, jobName string) (log string, err error)
	// GetJobStatus 获取job的状态
	GetJobStatus(ctx context.Context, jobName string) (status string, err error)
	// GetDeploymentStatus 获取部署的状态
	GetDeploymentStatus(ctx context.Context, deploymentName string) (status string, err error)
	// RemoveJob 删除job
	RemoveJob(ctx context.Context, jobName string) (err error)
	// RemoveDeployment 删除部署
	RemoveDeployment(ctx context.Context, deploymentName string) (err error)
	// WaitForTerminal 处理客户端发来的ws建立请求
	WaitForTerminal(ctx context.Context, ts Session, config Config, container, cmd string)
	// GetDeploymentContainerNames 获取部署的容器名
	GetDeploymentContainerNames(ctx context.Context, deploymentName string) (containerNames []string, err error)
}

// Middleware is a service middleware
type Middleware func(service Service) Service

func New(platform string, opts ...CreationOption) (Service, error) {
	switch platform {
	case "k8s":
		return NewK8s(opts...)
	case "docker":
		return NewDocker(opts...), nil
	default:
		return NewDocker(opts...), nil
	}
}
