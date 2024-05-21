# 项目配置

项目配置可以通过命令行传参或环境变量两种方式进行配置

推理或训练节点只需要安装**Docker**和**Nvidia-Docker**
即可。[NVIDIA Container Toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)

## 通过命令行传参

**需要注意的是，如果即设置了环境变量也设置了命令行参数，那么命令行参数的值会覆盖环境变量的值**

执行: `./aigc-server --help` 查看命令行参数

```bash
Usage:
  aigc-server [command]

Available Commands:
  account     用户相关操作命令
  completion  Generate the autocompletion script for the specified shell
  cronjob     定时任务
  generate    生成命令
  help        Help about any command
  job         任务命令
  start       启动http服务
  start-api   启动http api服务
  tenant      租户相关操作命令

Flags:
  -c, --config.path string                      配置文件路径，如果没有传入配置文件路径则默认使用环境变量
      --db.drive string                         数据库驱动 (default "sqlite")
      --db.mysql.database string                mysql数据库 (default "aigc")
      --db.mysql.host string                    mysql数据库地址: mysql (default "mysql")
      --db.mysql.metrics                        是否启GORM的Metrics
      --db.mysql.password string                mysql数据库密码
      --db.mysql.port int                       mysql数据库端口 (default 3306)
      --db.mysql.user string                    mysql数据库用户 (default "aigc")
  -h, --help                                    help for aigc-server
      --runtime.docker.workspace string         Docker工作目录 (default "/go/src/github.com/IceBearAI/LLM-And-More/storage")
      --runtime.gpu.num int                     GPU数量 (default 8)
      --runtime.k8s.config.path string          K8s配置文件路径
      --runtime.k8s.host string                 K8s地址
      --runtime.k8s.insecure                    K8s是否不安全
      --runtime.k8s.namespace string            K8s命名空间 (default "default")
      --runtime.k8s.token string                K8s Token
      --runtime.k8s.volume.name string          K8s挂载的存储名
      --runtime.platform string                 运行时平台 (default "docker")
      --runtime.shm.size string                 运行时共享内存大小 (default "16Gi")
      --server.admin.pass string                系统管理员密码 (default "admin")
      --server.admin.user string                系统管理员账号 (default "admin")
      --server.debug                            是否开启Debug模式
      --server.key string                       本系统服务密钥 (default "Aigcfj@202401")
      --server.log.drive string                 本系统日志驱动, 支持syslog,term (default "term")
      --server.log.level string                 本系统日志级别 (default "all")
      --server.log.name string                  本系统日志名称 (default "aigc-server.log")
      --server.log.path string                  本系统日志路径
  -a, --server.name string                      本系统服务名称 (default "aigc-server")
      --server.storage.path string              文件存储绝对路径 (default "/Users/cong/go/src/github.com/IceBearAI/LLM-And-More/storage")
      --service.fschat.controller.host string   fastchat controller address (default "http://fschat-controller:21001")
      --service.local.ai.host string            Chat-Api 地址 (default "http://localhost:8000/v1")
      --service.local.ai.token string           Chat-Api Token (default "sk-001")
      --service.openai.host string              OpenAI服务地址 (default "https://api.openai.com/v1")
      --service.openai.token string             OpenAI Token
      --tracer.drive string                     Tracer驱动 (default "jaeger")
      --tracer.enable                           是否启用Tracer
      --tracer.jaeger.host string               Tracer Jaeger Host (default "jaeger:6832")
      --tracer.jaeger.log.spans                 Tracer Jaeger Log Spans
      --tracer.jaeger.param float               Tracer Jaeger Param (default 1)
      --tracer.jaeger.type string               采样器的类型 const: 固定采样, probabilistic: 随机取样, ratelimiting: 速度限制取样, remote: 基于Jaeger代理的取样 (default "const")

Use "aigc-server [command] --help" for more information about a command.
```

## 启动后端管理http服务

执行: `./aigc-server start` 启动服务

```
Usage:
  aigc-server start [flags]

Flags:
  -c, --config.path string                      配置文件路径，如果没有传入配置文件路径则默认使用环境变量
      --db.drive string                         数据库驱动 (default "sqlite")
      --db.mysql.database string                mysql数据库 (default "aigc")
      --db.mysql.host string                    mysql数据库地址: mysql (default "mysql")
      --db.mysql.metrics                        是否启GORM的Metrics
      --db.mysql.password string                mysql数据库密码
      --db.mysql.port int                       mysql数据库端口 (default 3306)
      --db.mysql.user string                    mysql数据库用户 (default "aigc")
  -h, --help                                    help for aigc-server
      --runtime.docker.workspace string         Docker工作目录 (default "~/go/src/github.com/icowan/LLM-And-More/storage")
      --runtime.gpu.num int                     GPU数量 (default 8)
      --runtime.k8s.config.path string          K8s配置文件路径
      --runtime.k8s.host string                 K8s地址
      --runtime.k8s.insecure                    K8s是否不安全
      --runtime.k8s.namespace string            K8s命名空间 (default "default")
      --runtime.k8s.token string                K8s Token
      --runtime.k8s.volume.name string          K8s挂载的存储名
      --runtime.platform string                 运行时平台 (default "docker")
      --runtime.shm.size string                 运行时共享内存大小 (default "16Gi")
      --server.admin.pass string                系统管理员密码 (default "admin")
      --server.admin.user string                系统管理员账号 (default "admin")
      --server.debug                            是否开启Debug模式
      --server.key string                       本系统服务密钥 (default "Aigcfj@202401")
      --server.log.drive string                 本系统日志驱动, 支持syslog,term (default "term")
      --server.log.level string                 本系统日志级别 (default "all")
      --server.log.name string                  本系统日志名称 (default "aigc-server.log")
      --server.log.path string                  本系统日志路径
  -a, --server.name string                      本系统服务名称 (default "aigc-server")
      --server.storage.path string              文件存储绝对路径 (default "~/go/src/github.com/icowan/LLM-And-More/storage")
      --service.fschat.controller.host string   fastchat controller address (default "http://fschat-controller:21001")
      --service.local.ai.host string            Chat-Api 地址 (default "http://localhost:8000/v1")
      --service.local.ai.token string           Chat-Api Token (default "sk-001")
      --service.openai.host string              OpenAI服务地址 (default "https://api.openai.com/v1")
      --service.openai.only                     是否只使用OpenAI服务
      --service.openai.token string             OpenAI Token
      --tracer.drive string                     Tracer驱动 (default "jaeger")
      --tracer.enable                           是否启用Tracer
      --tracer.jaeger.host string               Tracer Jaeger Host (default "jaeger:6832")
      --tracer.jaeger.log.spans                 Tracer Jaeger Log Spans
      --tracer.jaeger.param float               Tracer Jaeger Param (default 1)
      --tracer.jaeger.type string               采样器的类型 const: 固定采样, probabilistic: 随机取样, ratelimiting: 速度限制取样, remote: 基于Jaeger代理的取样 (default "const")
```

## 启动定时任务

执行: `./aigc-server cronjob start` 启动定时任务

```
Usage:
  aigc-server cronjob start <args> [flags]

Examples:
如果 cronjob.auto 设置为 true 并且没有传入相应用的任务名称，则将自动运行所有的任务

aigc-server cronjob start -h

Flags:
      --cronjob.auto   是否自动执行定时任务 (default true)
  -h, --help           help for start
```

## 启动API服务

执行: `./aigc-server start-api` 启动API服务

```
Usage:
  aigc-server start-api [flags]

Examples:
## 启动命令
aigc-server start-api -p :8000


Flags:
  -h, --help                  help for start-api
  -p, --openapi.port string   服务启动的http api 端口 (default ":8000")
      --web.embed             是否使用embed.FS (default true)
```

## 启动FastChat-Controller

- `controller` 用于模型的注册中心及健康检查

通过docker启动:

```
$ docker run -d --network host -p 21001:21001 -it dudulu/fschat:v0.2.36 python3 -m fastchat.serve.controller --host 0.0.0.0 --port 21001
```

### 启动 aigc-server

如果`AIGC_RUNTIME_K8S_VOLUME_NAME=./aigc-data-cfs`

会在`AIGC_RUNTIME_DOCKER_WORKSPACE`下的`aigc-data-cfs`目录

**注意：需要提前将模型下载到 `AIGC_RUNTIME_DOCKER_WORKSPACE` 所指定的目录下**

可以通过环境变量或命令行传参：

#### 使用环境变量

```
$ AIGC_RUNTIME_GPU_NUM=4 AIGC_FSCHAT_CONTROLLER_ADDRESS=http://127.0.0.1:21001 HF_ENDPOINT=https://hf-mirror.com NO_PROXY=".idc,.corp,127.0.0.1,127.0.0.1" HTTP_PROXY=http://127.0.0.1:7890 HTTPS_PROXY=http://127.0.0.1:7890 AIGC_ADMIN_SERVER_STORAGE_PATH=/data/aigc/.cache/storage AIGC_DATASETS_IMAGE=reg.creditease.corp/aigc/qwen1.5-train:v0.2.36-0327 AIGC_RUNTIME_PLAORM=docker DOCKER_HOST=tcp://127.0.0.1:2376 AIGC_RUNTIME_DOCKER_WORKSPACE=/data/aigc/.cache/storage AIGC_RUNTIME_K8S_VOLUME_NAME=aigc-data-cfs AIGC_ADMIN_SERVER_DOMAIN=http://127.0.0.1:8080 ./aigc-server-linux-amd64-beta41 start
```

#### 使用命令行传参

```
$ export HF_ENDPOINT=https://hf-mirror.com DOCKER_HOST=tcp://127.0.0.1:2376
$ ./aigc-server-linux-amd64-beta4 start \
    --runtime.gpu.num 4 \
    --service.fschat.controller.host http://127.0.0.1:21001 \
    --service.local.ai.host http://127.0.0.1:8000 \
    --server.storage.path /data/aigc/.cache/storage/aigc \
    --datasets.image dudulu/llmops:latest \
    --runtime.platform docker \
    --runtime.docker.workspace /data/aigc/.cache/storage
```


## 系统公共环境变量配置

可以修改`.env`调整相关配置

## 数据库配置

目前支持两类数据库的配置，默认是使用**sqlite**,如果是使用的sqlite那么默认会存储在`AIGC_ADMIN_SERVER_STORAGE_PATH`
所配置的路径下的`storage/database/aigc.db`。

如果使用的是`mysql`驱动，则按照下面配置设置。

| 变量名                   | 描述               | 值       |
|-----------------------|------------------|---------|
| `AIGC_DB_DRIVER`      | 数据库驱动类型（可能是遗留错误） | `mysql` |
| `AIGC_MYSQL_DRIVE`    | 数据库驱动类型          | `mysql` |
| `AIGC_MYSQL_HOST`     | 数据库主机地址          | `mysql` |
| `AIGC_MYSQL_PORT`     | 数据库端口号           | `3306`  |
| `AIGC_MYSQL_USER`     | 数据库用户名           | `aigc`  |
| `AIGC_MYSQL_PASSWORD` | 数据库密码            | `admin` |
| `AIGC_MYSQL_DATABASE` | 数据库名             | `aigc`  |

## Tracer 链路追踪配置

如果想看整个调用链以下是相关配置，默认不开启。

| 变量名                            | 描述          | 值        |
|--------------------------------|-------------|----------|
| `AIGC_TRACER_ENABLE`           | 是否启用链路追踪    | `false`  |
| `AIGC_TRACER_DRIVE`            | 链路追踪驱动类型    | `jaeger` |
| `AIGC_TRACER_JAEGER_HOST`      | Jaeger 服务地址 |          |
| `AIGC_TRACER_JAEGER_PARAM`     | Jaeger 采样参数 | `1`      |
| `AIGC_TRACER_JAEGER_TYPE`      | Jaeger 采样类型 | `const`  |
| `AIGC_TRACER_JAEGER_LOG_SPANS` | 是否记录追踪日志    | `false`  |

## 跨域配置

跨域配置，默认不开启

| 变量名                           | 描述        | 值                                                                                                   |
|-------------------------------|-----------|-----------------------------------------------------------------------------------------------------|
| `AIGC_ENABLE_CORS`            | 是否启用CORS  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_METHODS`     | 允许的HTTP方法 | `GET,POST,PUT,DELETE,OPTIONS`                                                                       |
| `AIGC_CORS_ALLOW_HEADERS`     | 允许的HTTP头  | `Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,x-tenant-id,x-token` |
| `AIGC_CORS_ALLOW_CREDENTIALS` | 是否允许携带凭证  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_ORIGINS`     | 允许的源      | `*`                                                                                                 |

### 外部服务调用配置

chat的一些配置，假设使用的FastChat作为服务的推理框架，则配置FastChat的Api地址。

如果还有使用OpenAI的相关模型，则设置OpenAI的相关信息。

| 变量名                           | 描述            | 值                           |
|-------------------------------|---------------|-----------------------------|
| `AIGC_SERVICE_CHAT_API_HOST`  | 聊天API服务地址     | `http://localhost:8000/v1`  |
| `AIGC_SERVICE_CHAT_API_TOKEN` | 聊天API服务访问令牌   |                             |
| `AIGC_SERVICE_OPENAI_ORG_ID`  | OpenAI 组织ID   |                             |
| `AIGC_SERVICE_OPENAI_HOST`    | OpenAI 服务地址   | `https://api.openai.com/v1` |
| `AIGC_SERVICE_OPENAI_TOKEN`   | OpenAI 服务访问令牌 |                             |

### LDAP 配置

如果是企业使用可以配置LDAP地址。

| 变量名                   | 描述          | 值                                                        |
|-----------------------|-------------|----------------------------------------------------------|
| `AIGC_LDAP_HOST`      | LDAP 服务器地址  | `ldap`                                                   |
| `AIGC_LDAP_BASE_DN`   | LDAP 基础DN   | `OU=HABROOT,DC=corp`                                     |
| `AIGC_LDAP_BIND_USER` | LDAP 绑定用户   |                                                          |
| `AIGC_LDAP_BIND_PASS` | LDAP 绑定用户密码 |                                                          |
| `AIGC_LDAP_USER_ATTR` | LDAP 用户属性   | `name,mail,userPrincipalName,displayName,sAMAccountName` |

### aigc-server 环境变量配置

本系统基础配置，通常不需要修改。

需要注意的是`AIGC_ADMIN_SERVER_ADMIN_USER`和`AIGC_ADMIN_SERVER_ADMIN_PASS`是系统初始化的管理员账号密码，只有在系统第一次启动初始化的时候配置有效，后续调整将不会生效。

| 变量名                                     | 描述                            | 值                       |
|-----------------------------------------|-------------------------------|-------------------------|
| `AIGC_ADMIN_SERVER_HTTP_PORT`           | 服务HTTP端口                      | `:8080`                 |
| `AIGC_ADMIN_SERVER_LOG_DRIVE`           | 日志驱动类型(默认term)                | `term`                  |
| `AIGC_ADMIN_SERVER_NAME`                | 服务名称                          | `aigc-server`           |
| `AIGC_ADMIN_SERVER_DEBUG`               | 是否开启调试模式(开启后控制台显示所有Debug信息)   | `true`                  |
| `AIGC_ADMIN_SERVER_LOG_LEVEL`           | 日志级别(debug,info,warn,error)   | `all`                   |
| `AIGC_ADMIN_SERVER_LOG_PATH`            | 日志路径(设置之后会写入文件)               |                         |
| `AIGC_ADMIN_SERVER_LOG_NAME`            | 日志文件名称                        | `aigc-server.log`       |
| `AIGC_ADMIN_SERVER_DEFAULT_CHANNEL_KEY` | 默认渠道密钥                        | `sk-001`                |
| `AIGC_ADMIN_SERVER_STORAGE_PATH`        | 上传文件所存储的路径                    | `./storage/`            |
| `AIGC_ADMIN_SERVER_DOMAIN`              | 本服务的域名(容器回调传输数据，需要保证容器网络可以访问) | `http://localhost:8080` |
| `AIGC_ADMIN_SERVER_ADMIN_USER`          | 初始化默认账号                       | `admin`                 |
| `AIGC_ADMIN_SERVER_ADMIN_PASS`          | 初始化默认密码                       | `admin`                 |

## Runtime

可选择参数：

- `k8s`
- `docker`

运行时的平台，现支持**Kubernetes**和**Docker**作为模型运行的平台，默认为`docker`。

# Docker 平台

当`AIGC_RUNTIME_PLATFORM`设置为`docker`时可设置Docker本身支持的变量，如：`DOCKER_`开头的相关环境变量

- `AIGC_RUNTIME_DOCKER_WORKSPACE` 是指本机的模型目录，会映射到运行模型容器里的`/data/`目录。

要使用Docker API创建容器并挂载NVIDIA GPU，你需要确保你的系统上安装了NVIDIA Docker支持（例如nvidia-docker2）并且Docker守护进程配置正确。以下是使用Docker Engine API创建容器并挂载NVIDIA GPU的基本步骤：

确保你的Docker守护进程启用了NVIDIA GPU支持。这通常意味着你需要在Docker守护进程的配置文件中添加默认的运行时，例如`/etc/docker/daemon.json`：

```json
{
  "default-runtime": "nvidia",
  "runtimes": {
    "nvidia": {
      "path": "nvidia-container-runtime",
      "runtimeArgs": []
    }
  }
}
```

# k8s 平台

kubernetes支持两种方式连接

**通过Host和Token**

当配置了`AIGC_RUNTIME_K8S_HOST`和`AIGC_RUNTIME_K8S_TOKEN`时，HOST为`api-server`地址，如: `https://k8s:6443`

**通过config.yaml文件连接**

只需要配置`AIGC_RUNTIME_K8S_CONFIG_PATH`所在的当前路径

- `AIGC_RUNTIME_K8S_NAMESPACE`: 最终创建的job或deployment所在的空间，默认是`default`
- `AIGC_RUNTIME_K8S_VOLUME_NAME`: 存储的PVC名称，会将它挂载到容器的`/data`目录

| 变量名                             | 描述                | 值                  |
|---------------------------------|-------------------|--------------------|
| `AIGC_RUNTIME_PLATFORM`         | 运行平台              | `docker`           |
| `AIGC_RUNTIME_K8S_HOST`         | Kubernetes 服务地址   |                    |
| `AIGC_RUNTIME_K8S_TOKEN`        | Kubernetes 服务访问令牌 |                    |
| `AIGC_RUNTIME_K8S_NAMESPACE`    | Kubernetes 命名空间   | `default`          |
| `AIGC_RUNTIME_K8S_INSECURE`     | Kubernetes 不安全连接  | `false`            |
| `AIGC_RUNTIME_K8S_CONFIG_PATH`  | Kubernetes 配置路径   | `./k8sconfig.yaml` |
| `AIGC_RUNTIME_K8S_VOLUME_NAME`  | Kubernetes 卷名称    | `aigc-data`        |
| `AIGC_RUNTIME_SHM_SIZE`         | 共享内存大小            | `16G`              |
| `AIGC_RUNTIME_DOCKER_WORKSPACE` | Docker 工作空间       | `/tmp`             |

# Datasets

数据信标注内容检测的相关配置

| 变量名                            | 描述                                  | 值                            |
|--------------------------------|-------------------------------------|------------------------------|
| `AIGC_DATASETS_IMAGE`          | 检测数据集标注的相似度的镜像                      | `dudulu/llmops:latest`       |
| `AIGC_DATASETS_MODEL_NAME`     | 检测数据集的模型                            | `uer/sbert-base-chinese-nli` |
| `AIGC_DATASETS_DEVICE`         | 使用驱动如cpu,mps,cuda,npx               |                              |
| `AIGC_DATASETS_GPU_TOLERATION` | 通常是k8s设置脏节点的标签,通常我们会把带有GPU的节点设置为脏节点 |                              |

# 其他环境变量

以下环境变量会在创建容器时注入到容器运行时的环境变量中

| 变量名           | 描述                  | 值                       |
|---------------|---------------------|-------------------------|
| `HF_ENDPOINT` | Hugging Face 终端地址   | `https://hf-mirror.com` |
| `HF_HOME`     | Hugging Face 内容缓存目录 | `~/.cache/huggingface`  |
| `HTTP_PROXY`  | HTTP代理              |                         |
| `HTTPS_PROXY` | HTTPS代理             |                         |
| `NO_PROXY`    | 不使用代理的地址            |                         |


# Docker 部署

## Docker镜像

我们提供了Docker镜像，您可以直接使用我们提供的镜像，也可以自行构建。

- [LLMOps](docker/llmops/README.md)
- [百川2](docker/baichuan2/README.md)
- [FastChat](docker/fastchat/README.md)
- [Qwen](docker/qwen1.5/README.md)

### Docker交叉编译多平台

```
$ docker buildx ls
$ docker buildx rm --all-inactive
$ docker buildx create --driver-opt image=moby/buildkit:master --name builder --driver docker-container --use
$ docker buildx inspect --bootstrap
$ docker buildx create --platform linux/amd64,linux/arm64
$ docker login
$ docker buildx build --push -t dudulu/aigc-server:v0.0.0-bet03 --platform linux/amd64,linux/arm64 .
```
