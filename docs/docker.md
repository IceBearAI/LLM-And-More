## fschat

### 启动fschat-controller

设置环境变量

HOST_IP=127.0.0.1


```
$ docker run -d --network host -p 21001:21001 -it dudulu/fschat:v0.2.36 python3 -m fastchat.serve.controller --host 0.0.0.0 --port 21001
```

### 启动fschat-api

```
$ docker run -d --network host -p 8000:8000 -it dudulu/fschat:v0.2.36 python3 -m fastchat.serve.openai_api_server --host 0.0.0.0 --port 8000 --controller-address http://$(hostname -I | awk '{print $1}'):21001
```


## aigc-server

### 启动 aigc-server

如果`AIGC_RUNTIME_K8S_VOLUME_NAME=./aigc-data-cfs`

会在`AIGC_RUNTIME_DOCKER_WORKSPACE`下的`aigc-data-cfs`目录

#### 使用环境变量

```
$ AIGC_RUNTIME_GPU_NUM=4 AIGC_FSCHAT_CONTROLLER_ADDRESS=http://127.0.0.1:21001 AIGC_SERVICE_CHAT_API_HOST=http://127.0.0.1:8000 HF_ENDPOINT=https://hf-mirror.com NO_PROXY=".idc,.corp,127.0.0.1,127.0.0.1" HTTP_PROXY=http://127.0.0.1:7890 HTTPS_PROXY=http://127.0.0.1:7890 AIGC_ADMIN_SERVER_STORAGE_PATH=/data/aigc/.cache/storage AIGC_DATASETS_IMAGE=reg.creditease.corp/aigc/qwen1.5-train:v0.2.36-0327 AIGC_RUNTIME_PLAORM=docker DOCKER_HOST=tcp://127.0.0.1:2376 AIGC_RUNTIME_DOCKER_WORKSPACE=/data/aigc/.cache/storage AIGC_RUNTIME_K8S_VOLUME_NAME=aigc-data-cfs AIGC_ADMIN_SERVER_DOMAIN=http://127.0.0.1:8080 ./aigc-server-linux-amd64-beta41 start
```

#### 使用命令行传参

```
$ export HF_ENDPOINT=https://hf-mirror.com DOCKER_HOST=tcp://127.0.0.1:2376
$ ./aigc-server-linux-amd64-beta4 start \
    --runtime.gpu.num 4 \
    --service.fschat.controller.host http://127.0.0.1:21001 \
    --service.fschat.api.host http://127.0.0.1:8000 \
    --service.local.ai.host http://127.0.0.1:8000 \
    --server.storage.path /data/aigc/.cache/storage/aigc \
    --datasets.image dudulu/llmops:latest \
    --runtime.platform docker \
    --runtime.docker.workspace /data/aigc/.cache/storage
```
