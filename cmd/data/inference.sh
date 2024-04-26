#!/bin/bash

# 自动注入以下环境变量，可以直接使用
# - MODEL_PATH: 模型路径
# - MODEL_NAME: 模型名称
# - HTTP_PORT: HTTP端口
# - QUANTIZATION: 量化 8bit
# - NUM_GPUS: GPU数量
# - MAX_GPU_MEMORY: GPU内存
# - USE_VLLM: 是否是VLLM
# - INFERRED_TYPE: 推理类型
# - CONTROLLER_ADDRESS: fschat控制器地址
# - HF_ENDPOINT: huggingface地址
# - HF_HOME: huggingface家目录
# - HTTP_PROXY: HTTP代理
# - HTTPS_PROXY: HTTPS代理
# - NO_PROXY: 不代理的地址
# - MODEL_WORKER_TYPE: 模型工作类型

# shellcheck disable=SC2153
if [ "$HTTP_PROXY" != "" ]; then
    export http_proxy=$HTTP_PROXY
fi

if [ "$HTTPS_PROXY" != "" ]; then
    export https_proxy=$HTTPS_PROXY
fi

if [ "$NO_PROXY" != "" ]; then
    export no_proxy=$NO_PROXY
fi

MODEL_WORKER=fastchat.serve.model_worker
OS_TYPE=$(uname)

if [ "$MODEL_WORKER_TYPE" == "sglang" ]; then
    MODEL_WORKER="fastchat.serve.sglang_worker"
elif [ "$MODEL_WORKER_TYPE" == "vllm" ]; then
    MODEL_WORKER="fastchat.serve.vllm_worker"
fi

# 量化配置
if [ "$QUANTIZATION" == "8bit" ]; then
    QUANTIZATION="--load-8bit"
else
    QUANTIZATION=""
fi

# NUM_GPUS
if [ "$NUM_GPUS" -gt 0 ]; then
    NUM_GPUS="--num-gpus $NUM_GPUS"
else
    NUM_GPUS=""
fi

# CPU推理CPU，mps
if [ "$INFERRED_TYPE" == "cpu" ] && [ "$OS_TYPE" == "Darwin" ]; then
    DEVICE_OPTION="--device mps"
elif [ "$INFERRED_TYPE" == "cpu" ]; then
    DEVICE_OPTION="--device cpu"
else
    DEVICE_OPTION=""
fi

# MAX_GPU_MEMORY
if [ "$MAX_GPU_MEMORY" -gt 0 ]; then
    MAX_GPU_MEMORY="--max-gpu-memory ${MAX_GPU_MEMORY}GiB"
else
    MAX_GPU_MEMORY=""
fi

# 当前Pod的IP
if [ ! -n "$MY_POD_IP" ]; then
	MY_POD_IP=$(hostname -I | awk '{print $1}')
fi

# 模型路径
if [ -n "$MODEL_PATH" ]; then
	MODEL_PATH="--model-path $MODEL_PATH"
fi

python3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \
    --controller-address $CONTROLLER_ADDRESS \
    --worker-address http://$MY_POD_IP:$HTTP_PORT \
    --model-name $MODEL_NAME \
    $MODEL_PATH $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION