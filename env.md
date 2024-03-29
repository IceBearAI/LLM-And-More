#### 系统公共环境变量配置

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

###### k8s 平台

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
