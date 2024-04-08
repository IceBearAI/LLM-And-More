<div align=center><img src=docs/figure/logo.png width=40% /></div>


<div align=center><h1><b>构建、训练、部署：LLM-And-More让一切变得更简单</b></h1></div>


<!-- \[ [English](README.md) | 中文 \]-->

## LLM-And-More: 从创意到服务的一站式LLM解决方案

LLM-And-More 是一个专业、开箱即用的大模型训练及应用构建一站式解决方案，包含从数据到评估、从训练到部署、从想法到服务的全流程最佳实践。

项目旨在提供一个专业、开箱即用的大模型训练及应用构建一站式解决方案。您可以轻松地通过本项目进行模型训练，并一键生成所需的产品服务。我们的项目集成了不同应用场景下的专业知识和最佳实践，以保证模型在实际生产中的表现优异。此外，我们集成了高性能模型并行框架，有效地减少了训练和推理时的算力开销，提升了整体的效率。无论您是在构建基于大模型的服务还是在进行私有化模型训练，我们的开源项目将为您提供一个强大且可靠的基础。

如有任何问题，请加入我们的微信群
<div align=center><img src=docs/figure/wechat.jpg width=50% /></div>

## 目录

- [功能模块](#功能模块)
- [支持场景](#支持场景)
- [适配模型](#适配模型)
- [安装与使用](#安装与使用)

## 功能模块

LLM-And-More致力于为专业开发者和一线业务人员提供同样专业、易用的LLM应用构建方案。为了实现这一理念，LLM-And-More将LLM应用开发过程分解为以下六个模块：

- [数据模块](#数据模块)
- [训练模块](#训练模块)
- [监控模块](#监控模块)
- [评估模块](#评估模块)
- [部署模块](#部署模块)
- [交互模块](#交互模块)

这些模块涵盖了开发一个LLM产品的一切，并注入了专业知识和性能优化组件，全流程协助您构建一个优秀的LLM应用。

### 数据模块

为协助您开始LLM应用构建的第一步，也是影响最终效果的最重要一步：数据标注，LLM-And-More提供了一个功能齐全的数据标注平台。在这里，管理人员可以新建一个数据标注任务，向标注人员分配需要标注的数据，或是亲自对数据进行标注。在完成所有标注之后，LLM-And-More将把标注完成的数据自动转换为模型可以处理的格式(
jsonl)，并存入本地数据库中，在之后的训练、评估模块中一键启用。LLM-And-More还提供了对数据质量的一键检测，用户可以通过查看数据检测报告，发现标注过程中可能产生的错误，提升模型训练的最终效果。

![数据示例](docs/figure/datasets-sample.png)

### 训练模块

为协助您完成LLM应用构建最专业，最困难的一步：模型训练，LLM-And-More提供了一个开箱即用的高性能模型训练框架，使您无需了解任何深度学习相关知识，即可轻松对齐大模型训练最佳实践。在训练模块中，用户可以自由调整所选取的基座模型，训练方式，以及batch_size，学习率等超参数。如果用户对此并不了解，LLM-And-More预置的智能默认参数将帮助用户完成参数的选取和调优。LLM-And-More将自动为用户提供DeepSpeed多卡多机加速适配，帮助用户节约训练时间，充分利用算力资源。

![微调示例](docs/figure/finetune-add.png)

### 监控模块

在训练过程中，您可能烦恼于无法清晰直观的观察模型的性能变化，LLM-And-More提供了一个智能化的模型训练监控模块，不仅可以实时、可视化的显示CPU、GPU等核心算力资源的占用情况，监控模型Loss，学习率，训练步数的变化，还可以在训练过程中智能提示您模型的潜在性能风险，并提供恰当的解决方案建议。例如，您在训练过程中可能发现系统提示“过拟合风险”，并建议您“停止训练，降低学习率或增大数据量”，您可以遵循这些建议，尝试解决问题。这节省了不必要的算力浪费，并能使您更精准的把握模型可能的表现。

![监控示例](docs/figure/finetune-graph.png)

### 评估模块

在完成训练后，您可能烦恼于无法准确了解模型的性能边界，不敢真正将模型应用于实际场景，LLM-And-More提供了一个客观准确的评估模块，让您可以全面，精准的把握模型的各项能力水准。评估模块分为两个主要评测任务，一个负责评测模型在训练过任务上的性能，主要展现模型是否能够完成用户指定的任务，例如一个客服模型是否能够正确回复用户提问；另一个负责评测模型在五个通用维度上的能力（推理能力，阅读理解能力，中文能力，指令遵从能力，创新能力），主要体现模型是否产生了灾难性遗忘，是否过于专注于特定领域而丧失了通用能力。用户可以根据两方面的反馈，调整模型的训练数据和轮次，选择表现最佳的模型应用于线上系统。

![评估示例](docs/figure/five.png)

### 部署模块(Coming Soon)

### 交互模块(Coming Soon)

## 支持场景

除了简单的提供输入输出训练您的LLM应用外，我们还提供了丰富的场景支持，帮助您更好的解决您在工作生产当中遇到的复杂问题。例如，您可以基于FAQ场景适配方案，直接构建一个客服，协助您识别用户意图，解决您淘宝店铺自动回复的问题；或是基于RAG场景方案，构建一套企业内部规章制度问答机器人。所有场景均具备独立的留个功能模块，但拥有深度定制的UI和专业Know-How。LLM-And-More支持以下场景：

- [通用场景](#通用场景)
- [FAQ场景](#FAQ场景)
- [RAG场景](#RAG场景)
- [创意写作场景](#创意写作场景)
- [Agent场景](#Agent场景)

### 通用场景

可以接受任意的输入输出，是最基础的训练场景。在该场景中，我们没有针对数据特征、应用范围等进行任何假设，在数据、训练、监控、评估、部署、交互模块中的各项参数均调整至最均衡的水平，并适配了任何场景均有收益的专业Know-How辅助模型训练。总的来说，如果您不确定您应该使用哪个场景，或认为LLM-And-More提供的任何场景均不符合您的要求，您可以选用该场景开始您的构建。

### FAQ场景

该场景适用于FAQ客服或FAQ问答机器人。FAQ(frequently asked questions)
即常见问题，通常应用于客服、快速助手、和在线论坛等场景，在这些地方，常见问题往往会反复出现，例如，用户经常以各种不同的表达方式询问发货时间。在FAQ场景中，用户的问题往往被归类，并对于每一类问题有一个统一的回答，例如针对询问发货时间的问题，统一回复“我们将尽快安排发货，请及时查看物流信息”。直接尝试使用大模型生成这些回答往往是事倍功半的，因为没有充分利用数据的特征，很难让大模型回复稳定，并常常会产生幻觉现象。为此，我们为FAQ场景设计了涵盖全部六个模块的全流程解决方案，主要引入用户意图识别，让LLM预测用户意图（在上述例子中，“查询物流信息”），而不是直接预测回复（在上述例子中，“我们将尽快安排发货，请及时查看物流信息”）。我们在大量FAQ场景中的实验表明，该解决方案可以提升30%以上的回复准确率，并显著降低模型过拟合风险与幻觉现象。

### RAG场景(Coming Soon)

### 创意写作场景(Coming Soon)

### Agent场景(Coming Soon)

## 适配模型

为了支持更广泛的应用，LLM-And-More支持多种多样的模型，并支持您选择多种不同的训练方式。

| 模型名                                                  | 模型大小                    | 支持的训练方法    |
|------------------------------------------------------|-------------------------|------------|
| [Baichuan2](https://huggingface.co/baichuan-inc)     | 7B/13B                  | 全参数训练/Lora |
| [ChatGLM3](https://huggingface.co/THUDM/chatglm3-6b) | 6B                      | 全参数训练/Lora |
| [LLaMA](https://github.com/facebookresearch/llama)   | 7B/13B/33B/65B          | 全参数训练/Lora |
| [LLaMA-2](https://huggingface.co/meta-llama)         | 7B/13B/70B              | 全参数训练/Lora |
| [Qwen](https://huggingface.co/Qwen)                  | 0.5B/1.8B/4B/7B/14B/72B | 全参数训练/Lora |

### 系统架构设计

系统是前后端分离的架构。

#### 模型推理框架

我们使用的是[FastChat](https://github.com/lm-sys/FastChat)作为模型推理框架，FastChat是一个非常优秀的开源项目。

> [FastChat](https://github.com/lm-sys/FastChat) 是一个开放平台，用于训练、服务和评估基于大型语言模型的聊天机器人。

**FastChat我们主要用其三个服务**

`controller` 用于模型的注册中心及健康检查

`worker` 服务启动模型并将当前模型注册到controller

`api` 从controller获取模型的地址代理到worker并提供标准API

我们主要通过它来实现大模型的高可用，高可扩展性。

![模型推理框架](https://github.com/lm-sys/FastChat/raw/main/assets/server_arch.png)

模型部署的操作可以参考[模型部署](docs/model/list.md)

### 模型微调

为了实现模型的微调，您可以参考我们的详细指南：[模型微调](docs/model/finetune.md)。

### 模型部署与微调

您可以将模型部署到任意配备GPU的节点上，无论是私有的K8s集群、Docker集群，还是云服务商提供的K8s集群，均能轻松对接。

### 本系统组成

本系统主要由以下几个部分组成：

- **HTTP服务**：提供Web服务接口，方便用户进行交互。
- **定时任务**：执行预定任务，如模型训练、数据预处理等。
- **训练镜像**：包含所有必要的环境和依赖，用于模型的训练和微调。

- 通过这些组件的协同工作，我们能够提供一个灵活、高效的模型微调和部署解决方案。

#### 部署流程

```mermaid
graph LR
    A[aigc] --> B[点击部署]
    B --> C[创建部署模版]
    C --> D[使用Docker或k8s进行调度]
    D --> E[挂载相应配置有模型]
    E --> F[启动模型]
    F --> G[注册到fschat-controller]
```

#### 微调训练流程

```mermaid
graph LR
    A[aigc] --> B[上传微调文件]
    B --> C[生成微调模版]
    C --> D[使用Docker或k8s进行调度]
    D --> E[挂载相应配置有模型]
    E --> F[启动训练脚本]
    F --> G[输出日志]
```

## 使用手册

[AIGC平台使用手册](docs/SUMMARY.md)

### 安装使用步骤

- 克隆项目: `git clone https://github.com/IceBearAI/aigc.git`
- 进入项目: `cd aigc-server`

该系统依赖**Docker**需要安装此服务

推理或训练节点只需要安装**Docker**和**Nvidia-Docker**
即可。[NVIDIA Container Toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)

#### 本地开发

[golang](https://github.com/golang/go)版本请安装go1.21以上版本

- 安装依赖包: `go mod tidy`
- 本地启动: `make run`
- build成x86 Linux可执行文件: `make build-linux`
- build成当前电脑可执行文件: `make build`

build完通常会保存在 `$(GOPATH)/bin/` 目录下

#### Docker-compose部署

安装docker和docker-compose可以参考官网教程：[Install Docker Engine](https://docs.docker.com/engine/install/)

需要在 当前目录下增加 `.env` 文件来设置相关的环境变量。

执行命令启动全部服务

```
$ docker-compose up
```

如果不需要执行build流程，可以进入到`docker`目录下执行`docker-compose up`即可。或把`docker-compose.yaml`的`build`注释掉。

### 项目配置

项目配置可以通过命令行传参或环境变量两种方式进行配置

#### 通过命令行传参

**需要注意的是，如果即设置了环境变量也设置了命令行参数，那么命令行参数的值会覆盖环境变量的值**

执行: `./aigc-server --help` 查看命令行参数

```bash
Usage:
  aigc-server [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cronjob     定时任务
  generate    生成命令
  help        Help about any command
  job         任务命令
  start       启动http服务

Flags:
  -c, --config.path string                配置文件路径，如果没有传入配置文件路径则默认使用环境变量
      --db.drive string                   数据库驱动 (default "sqlite")
      --db.mysql.database string          mysql数据库 (default "aigc")
      --db.mysql.host string              mysql数据库地址: mysql (default "mysql")
      --db.mysql.metrics                  是否启GORM的Metrics
      --db.mysql.password string          mysql数据库密码
      --db.mysql.port int                 mysql数据库端口 (default 3306)
      --db.mysql.user string              mysql数据库用户 (default "aigc")
  -h, --help                              help for aigc-server
  -n, --namespace string                  命名空间 (default "aigc")
      --runtime.docker.workspace string   Docker工作目录 (default "github.com/IceBearAI/LLM-And-More/storage")
      --runtime.k8s.config.path string    K8s配置文件路径
      --runtime.k8s.host string           K8s地址
      --runtime.k8s.insecure              K8s是否不安全
      --runtime.k8s.namespace string      K8s命名空间 (default "default")
      --runtime.k8s.token string          K8s Token
      --runtime.k8s.volume.name string    K8s挂载的存储名
      --runtime.platform string           运行时平台 (default "docker")
      --runtime.shm.size string           运行时共享内存大小 (default "16G")
      --server.admin.pass string          系统管理员密码 (default "admin")
      --server.admin.user string          系统管理员账号 (default "admin")
      --server.debug                      是否开启Debug模式
      --server.key string                 本系统服务密钥 (default "Aigcfj@202401")
      --server.log.drive string           本系统日志驱动, 支持syslog,term (default "term")
      --server.log.level string           本系统日志级别 (default "all")
      --server.log.name string            本系统日志名称 (default "aigc-server.log")
      --server.log.path string            本系统日志路径
  -a, --server.name string                本系统服务名称 (default "aigc-server")
      --server.storage.path string        文件存储绝对路径 (default "github.com/IceBearAI/LLM-And-More/storage")
      --service.local.ai.host string      Chat-Api 地址 (default "http://fschat-api:8000/v1")
      --service.local.ai.token string     Chat-Api Token (default "sk-001")
      --service.openai.enable             是否启用OpenAI服务
      --service.openai.host string        OpenAI服务地址 (default "https://api.openai.com/v1")
      --service.openai.model string       OpenAI模型名称 (default "gpt-3.5-turbo")
      --service.openai.org.id string      OpenAI OrgId
      --service.openai.token string       OpenAI Token

Use "aigc-server [command] --help" for more information about a command.
```

##### 启动http服务

执行: `./aigc-server start --cronjob.auto false` 启动服务

```
Usage:
  aigc-server start [flags]

Flags:
      --cors.allow.credentials           是否允许跨域访问的凭证 (default true)
      --cors.allow.headers string        允许跨域访问的头部 (default "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
      --cors.allow.methods string        允许跨域访问的方法 (default "GET,POST,PUT,DELETE,OPTIONS")
      --cors.allow.origins string        允许跨域访问的域名 (default "*")
      --cors.enable                      是否开启跨域访问
      --cors.expose.headers string       允许跨域访问的头部 (default "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
      --datasets.device string           datasets device
      --datasets.gpu.toleration string   datasets gpu toleration
      --datasets.image string            datasets image (default "dudulu/llmops:latest")
      --datasets.model.name string       datasets model name (default "uer/sbert-base-chinese-nli")
  -h, --help                             help for start
  -p, --http.port string                 服务启动的http端口 (default ":8080")
      --ldap.base.dn string              LDAP Base DN (default "OU=HABROOT,DC=ORG,DC=corp")
      --ldap.bind.pass string            LDAP Bind Password
      --ldap.bind.user string            LDAP Bind User (default "aigc_ldap")
      --ldap.group.filter string         LDAP Group Filter
      --ldap.host string                 LDAP地址 (default "ldap://ldap")
      --ldap.port int                    LDAP端口 (default 389)
      --ldap.use.ssl                     LDAP Base DN
      --ldap.user.attr strings           LDAP Attributes (default [name,mail,userPrincipalName,displayName,sAMAccountName])
      --ldap.user.filter string          LDAP User Filter (default "(userPrincipalName=%s)")
      --server.domain string             启动服务的域名 (default "http://localhost:8080")
      --storage.type string              storage type (default "local")
      --tracer.drive string              Tracer驱动 (default "jaeger")
      --tracer.enable                    是否启用Tracer
      --tracer.jaeger.host string        Tracer Jaeger Host (default "jaeger:6832")
      --tracer.jaeger.log.spans          Tracer Jaeger Log Spans
      --tracer.jaeger.param float        Tracer Jaeger Param (default 1)
      --tracer.jaeger.type string        采样器的类型 const: 固定采样, probabilistic: 随机取样, ratelimiting: 速度限制取样, remote: 基于Jaeger代理的取样 (default "const")
      --web.embed                        是否使用embed.FS (default true)
      --cronjob.auto                     是否自动执行定时任务 (default true)
```

##### 启动定时任务

改模块可以独立启动

执行: `./aigc-server cronjob start` 启动定时任务

定时任务主要有

- 获取模型部署状态
- 获取训练日志
- 获取正在等待调度训练的任务

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

##### 启动参考

**使用命令行传参**

```
$ export HF_ENDPOINT=https://hf-mirror.com DOCKER_HOST=tcp://127.0.0.1:2376
$ ./aigc-server-linux-amd64 start \
    --runtime.gpu.num 4 \
    --service.fschat.controller.host http://127.0.0.1:21001 \
    --service.fschat.api.host http://127.0.0.1:8000 \
    --server.storage.path /data/aigc/storage \
    --datasets.image dudulu/llmops:latest \
    --runtime.platform docker \
    --runtime.docker.workspace /data/aigc/storage\
```

**使用环境变量启动**

```
$ export AIGC_RUNTIME_GPU_NUM=4 
$ export AIGC_FSCHAT_CONTROLLER_ADDRESS=http://127.0.0.1:21001 
$ export AIGC_SERVICE_CHAT_API_HOST=http://127.0.0.1:8000 
$ export HF_ENDPOINT=https://hf-mirror.com 
$ export AIGC_ADMIN_SERVER_STORAGE_PATH=/data/aigc/storage
$ export AIGC_DATASETS_IMAGE=dudulu/llmops:latest
$ export AIGC_RUNTIME_PLAORM=docker DOCKER_HOST=tcp://127.0.0.1:2376 
$ export AIGC_RUNTIME_DOCKER_WORKSPACE=/data/aigc/storage
$ export AIGC_ADMIN_SERVER_DOMAIN=http://127.0.0.1:8080
$ ./aigc-server-linux-amd64 start
```

#### 系统环境变量配置

可以修改`.env`调整相关配置

##### 数据库配置

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

##### Tracer 链路追踪配置

如果想看整个调用链以下是相关配置，默认不开启。

| 变量名                            | 描述          | 值        |
|--------------------------------|-------------|----------|
| `AIGC_TRACER_ENABLE`           | 是否启用链路追踪    | `false`  |
| `AIGC_TRACER_DRIVE`            | 链路追踪驱动类型    | `jaeger` |
| `AIGC_TRACER_JAEGER_HOST`      | Jaeger 服务地址 |          |
| `AIGC_TRACER_JAEGER_PARAM`     | Jaeger 采样参数 | `1`      |
| `AIGC_TRACER_JAEGER_TYPE`      | Jaeger 采样类型 | `const`  |
| `AIGC_TRACER_JAEGER_LOG_SPANS` | 是否记录追踪日志    | `false`  |

##### 跨域配置

跨域配置，默认不开启

| 变量名                           | 描述        | 值                                                                                                   |
|-------------------------------|-----------|-----------------------------------------------------------------------------------------------------|
| `AIGC_ENABLE_CORS`            | 是否启用CORS  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_METHODS`     | 允许的HTTP方法 | `GET,POST,PUT,DELETE,OPTIONS`                                                                       |
| `AIGC_CORS_ALLOW_HEADERS`     | 允许的HTTP头  | `Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,x-tenant-id,x-token` |
| `AIGC_CORS_ALLOW_CREDENTIALS` | 是否允许携带凭证  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_ORIGINS`     | 允许的源      | `*`                                                                                                 |

##### 外部服务调用配置

chat的一些配置，假设使用的FastChat作为服务的推理框架，则配置FastChat的Api地址。

如果还有使用OpenAI的相关模型，则设置OpenAI的相关信息。

- `AIGC_SERVICE_CHAT_API_HOST`: 指的就是我们启动的fschat-api的服务地址，也就是模型推理地下
- `AIGC_FSCHAT_CONTROLLER_ADDRESS`: 指的就是我们启动的fschat-controller的地址，模型部署时会将该变理注入进去模型启动后会将信息注册到该地址

| 变量名                              | 描述                     | 值                                |
|----------------------------------|------------------------|----------------------------------|
| `AIGC_SERVICE_CHAT_API_HOST`     | 聊天API服务地址              | `http://fschat-api:8000`         |
| `AIGC_SERVICE_CHAT_API_TOKEN`    | 聊天API服务访问令牌            |                                  |
| `AIGC_SERVICE_OPENAI_ORG_ID`     | OpenAI 组织ID            |                                  |
| `AIGC_SERVICE_OPENAI_HOST`       | OpenAI 服务地址            | `https://api.openai.com/v1`      |
| `AIGC_SERVICE_OPENAI_TOKEN`      | OpenAI 服务访问令牌          |                                  |
| `AIGC_FSCHAT_CONTROLLER_ADDRESS` | FastChat Controller的地址 | `http://fschat-controller:21001` |

##### S3 存储配置

企业使用可以配置文件存储在S3上，通过设置环境变量`AIGC_STORAGE_TYPE`来配置存储类型，默认为`local`表示存在本地。

| 变量名                             | 描述         | 值 |
|---------------------------------|------------|---|
| `AIGC_SERVICE_S3_HOST`          | S3 服务地址    |   |
| `AIGC_SERVICE_S3_ACCESS_KEY`    | S3 访问密钥    |   |
| `AIGC_SERVICE_S3_SECRET_KEY`    | S3 访问密钥密码  |   |
| `AIGC_SERVICE_S3_BUCKET`        | S3 存储桶名称   |   |
| `AIGC_SERVICE_S3_BUCKET_PUBLIC` | S3 公共存储桶名称 |   |
| `AIGC_SERVICE_S3_PROJECT_NAME`  | S3 项目名称    |   |

##### LDAP 配置

如果是企业使用可以配置LDAP地址。

| 变量名                   | 描述          | 值                                                        |
|-----------------------|-------------|----------------------------------------------------------|
| `AIGC_LDAP_HOST`      | LDAP 服务器地址  | `ldap`                                                   |
| `AIGC_LDAP_BASE_DN`   | LDAP 基础DN   | `OU=HABROOT,DC=corp`                                     |
| `AIGC_LDAP_BIND_USER` | LDAP 绑定用户   |                                                          |
| `AIGC_LDAP_BIND_PASS` | LDAP 绑定用户密码 |                                                          |
| `AIGC_LDAP_USER_ATTR` | LDAP 用户属性   | `name,mail,userPrincipalName,displayName,sAMAccountName` |

##### aigc-server 环境变量配置

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

##### Runtime

可选择参数：

- `k8s`
- `docker`

运行时的平台，现支持**Kubernetes**和**Docker**作为模型运行的平台，默认为`docker`。

###### Docker 平台

当`AIGC_RUNTIME_PLATFORM`设置为`docker`时可设置Docker本身支持的变量，如：`DOCKER_`开头的相关环境变量

- `AIGC_RUNTIME_DOCKER_WORKSPACE` 是指本机的模型目录，会映射到运行模型容器里的`/data/`目录。
- `AIGC_RUNTIME_GPU_NUM` 当前主机的GPU总数量，如果不设置默认是`8`，默认会从第`0`块卡启动
-

要使用Docker API创建容器并挂载NVIDIA GPU，你需要确保你的系统上安装了NVIDIA
Docker支持（例如nvidia-docker2）并且Docker守护进程配置正确。以下是使用Docker Engine API创建容器并挂载NVIDIA GPU的基本步骤：

确保你的Docker守护进程启用了NVIDIA
GPU支持。这通常意味着你需要在Docker守护进程的配置文件中添加默认的运行时，例如`/etc/docker/daemon.json`：

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

需要注意的是`DOCKER_HOST`暂时只支持tcp连接docker，使用`unix:///run/containerd/containerd.sock`

- `DOCKER_HOST=tcp://127.0.0.1:2376`

**配置docker**

使用`systemctl status docker`查看`Unit`的位置，通常会在`/usr/lib/systemd/system/docker.service`

在`ExecStart=/usr/bin/dockerd`后面加上`-H tcp://0.0.0.0:2376`，保存后重启动docker

```
[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2376 -H fd:// --containerd=/run/containerd/containerd.sock
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutStartSec=0
RestartSec=2
Restart=always
```

**重启docker**

```
$ systemctl daemon-reload && systemctl restart docker
```

###### k8s 平台

kubernetes支持两种方式连接

**通过Host和Token**

当配置了`AIGC_RUNTIME_K8S_HOST`和`AIGC_RUNTIME_K8S_TOKEN`时，HOST为`api-server`地址，如: `https://k8s:6443`

**通过config.yaml文件连接**

只需要配置`AIGC_RUNTIME_K8S_CONFIG_PATH`所在的当前路径

- `AIGC_RUNTIME_K8S_NAMESPACE`: 最终创建的job或deployment所在的空间，默认是`default`
- `AIGC_RUNTIME_K8S_VOLUME_NAME`: 存储的PVC名称，会将它挂载到容器的`/data`目录
  > 如果`AIGC_RUNTIME_K8S_VOLUME_NAME=aigc-data`
  > 会在`AIGC_RUNTIME_DOCKER_WORKSPACE`下的`aigc-data`目录

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
| `AIGC_RUNTIME_GPU_NUM`          | 当前主机的GPU总数量       | `8`                |

##### Datasets

数据信标注内容检测的相关配置

| 变量名                            | 描述                                  | 值                            |
|--------------------------------|-------------------------------------|------------------------------|
| `AIGC_DATASETS_IMAGE`          | 检测数据集标注的相似度的镜像                      | `dudulu/llmops:latest`       |
| `AIGC_DATASETS_MODEL_NAME`     | 检测数据集的模型                            | `uer/sbert-base-chinese-nli` |
| `AIGC_DATASETS_DEVICE`         | 使用驱动如cpu,mps,cuda,npx               |                              |
| `AIGC_DATASETS_GPU_TOLERATION` | 通常是k8s设置脏节点的标签,通常我们会把带有GPU的节点设置为脏节点 |                              |

##### 其他环境变量

以下环境变量会在创建容器时注入到容器运行时的环境变量中

| 变量名           | 描述                  | 值                       |
|---------------|---------------------|-------------------------|
| `HF_ENDPOINT` | Hugging Face 终端地址   | `https://hf-mirror.com` |
| `HF_HOME`     | Hugging Face 内容缓存目录 | `/data/hf`              |
| `HTTP_PROXY`  | HTTP代理              |                         |
| `HTTPS_PROXY` | HTTPS代理             |                         |
| `NO_PROXY`    | 不使用代理的地址            |                         |

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
$ docker buildx build --push -t dudulu/aigc-server:latest --platform linux/amd64,linux/arm64 .
```
