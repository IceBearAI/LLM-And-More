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

# Set proxy if provided
[ -n "$HTTP_PROXY" ] && export http_proxy=$HTTP_PROXY
[ -n "$HTTPS_PROXY" ] && export https_proxy=$HTTPS_PROXY
[ -n "$NO_PROXY" ] && export no_proxy=$NO_PROXY

MODEL_WORKER="fastchat.serve.model_worker"
OS_TYPE=$(uname)

# Set model worker based on the type
case "$MODEL_WORKER_TYPE" in
    sglang)
        MODEL_WORKER="fastchat.serve.sglang_worker"
        ;;
    vllm)
        MODEL_WORKER="fastchat.serve.vllm_worker"
        ;;
    tensorrt)
        MODEL_PATH_ENGINE="${MODEL_PATH}/trtllm/engines/fp16/${NUM_GPUS}-gpu"
        MODEL_PATH_CKPT="${MODEL_PATH}/trtllm/ckpt/fp16/${NUM_GPUS}-gpu"
        BASE_MODEL_NAMES=("qwen" "llama" "phi" "chatglm" "baichuan")

        if [ ! -d "$MODEL_PATH_ENGINE" ]; then
           selected_base_model=""
           model_name_lower=$(echo "$MODEL_NAME" | tr '[:upper:]' '[:lower:]')
           for BASE_MODEL_NAME in "${BASE_MODEL_NAMES[@]}"; do
               if echo "$model_name_lower" | grep -q "$BASE_MODEL_NAME"; then
                   selected_base_model="$BASE_MODEL_NAME"
                   break
               fi
           done

            # Generate TensorRT engine if the directory does not exist
            python3 /app/TensorRT-LLM/examples/"$selected_base_model"/convert_checkpoint.py \
              --model_dir "$MODEL_PATH" \
              --output_dir "$MODEL_PATH_CKPT" \
              --dtype float16 --tp_size "$NUM_GPUS"

            trtllm-build --checkpoint_dir "$MODEL_PATH_CKPT" \
              --output_dir "$MODEL_PATH_ENGINE" \
              --gemm_plugin float16

            MODEL_PATH=$MODEL_PATH_ENGINE
        else
          MODEL_PATH=$MODEL_PATH_ENGINE
        fi

        MODEL_WORKER="fastchat.serve.trt_worker"
        ;;
esac

echo "Using model worker: $MODEL_WORKER"
echo "Model path: $MODEL_PATH"


# Set quantization
QUANTIZATION=$([ "$QUANTIZATION" == "8bit" ] && echo "--load-8bit" || echo "")

# Set number of GPUs
NUM_GPUS=$([ "$NUM_GPUS" -gt 0 ] && echo "--num-gpus $NUM_GPUS" || echo "")

# Set device option
if [ "$INFERRED_TYPE" == "cpu" ]; then
    DEVICE_OPTION=$([ "$OS_TYPE" == "Darwin" ] && echo "--device mps" || echo "--device cpu")
else
    DEVICE_OPTION=""
fi

# Set max GPU memory
MAX_GPU_MEMORY=$([ "$MAX_GPU_MEMORY" -gt 0 ] && echo "--max-gpu-memory ${MAX_GPU_MEMORY}GiB" || echo "")

# Set pod IP
MY_POD_IP=${MY_POD_IP:-$(hostname -I | awk '{print $1}')}

# Set model path
MODEL_PATH=${MODEL_PATH:+"--model-path $MODEL_PATH"}

python3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \
    --controller-address $CONTROLLER_ADDRESS \
    --worker-address http://$MY_POD_IP:$HTTP_PORT \
    --model-name $MODEL_NAME \
    $MODEL_PATH $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION