#!/bin/bash
# 自动注入以下环境变量，可以直接使用
# - TENANT_ID: 租户ID
# - JOB_ID: 任务ID
# - AUTH: 授权码
# - API_URL: 回调API地址
# - HF_ENDPOINT: huggingface地址
# - HF_HOME: huggingface家目录
# - HTTP_PROXY: HTTP代理
# - HTTPS_PROXY: HTTPS代理
# - NO_PROXY: 不代理的地址
# - GPUS_PER_NODE: 每个节点的GPU数量
# - MASTER_PORT: 主节点端口
# - USE_LORA: 是否使用LORA
# - OUTPUT_DIR: 输出路径
# - NUM_TRAIN_EPOCHS: 训练轮数
# - PER_DEVICE_TRAIN_BATCH_SIZE: 训练批次大小
# - PER_DEVICE_EVAL_BATCH_SIZE: 评估批次大小
# - GRADIENT_ACCUMULATION_STEPS: 累积步数
# - LEARNING_RATE: 学习率
# - MODEL_MAX_LENGTH: 模型最大长度
# - BASE_MODEL_PATH: 基础模型路径
# - BASE_MODEL_NAME: 基础模型名称
# - SCENARIO: 应用场景
# - TRAIN_FILE: 训练文件URL或路径
# - EVAL_FILE: 验证文件URL或路径
export CUDA_DEVICE_MAX_CONNECTIONS=1

# 使用 nvidia-smi 获取 GPU 信息
gpu_info=$(nvidia-smi --query-gpu=name --format=csv,noheader)

# 设置一个标志，初始为未找到
found=0

# 检查每个型号
for model in "RTX 4090" "RTX 3090" "RTX 4080"
do
    if echo "$gpu_info" | grep -q "$model"; then
        echo "存在 $model GPU。"
        found=1
    fi
done

# 如果没有找到任何型号
if [ $found -eq 0 ]; then
    echo "未检测到 RTX 4090、RTX 3090 或 RTX 4080 GPU。"
else
    export NCCL_P2P_DISABLE=1
    export NCCL_IB_DISABLE=1
fi

#NNODES=1
#NODE_RANK=0
#MASTER_ADDR=localhost
Q_LORA=False

DS_CONFIG_PATH="ds_config_zero3.json"

DISTRIBUTED_ARGS="
    --nproc_per_node $GPUS_PER_NODE \
    --nnodes $NNODES \
    --node_rank $NODE_RANK \
    --master_addr $MASTER_ADDR \
    --master_port $MASTER_PORT
"
if [ "$USE_LORA" == "true" ]; then
    USE_LORA=True
    DS_CONFIG_PATH="ds_config_zero2.json"
else
    USE_LORA=False
    DS_CONFIG_PATH="ds_config_zero3.json"
fi


JOB_STATUS="success"
JOB_MESSAGE=""

function callback() {
  # 根据退出状态判断执行是否异常
  if [ $JOB_STATUS -eq 0 ]; then
      # 没有发生异常，正常输出内容
      echo "执行成功!"
      echo "${API_URL}"
      # 调用API并传递输出内容
      curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"success\"}"
  else
      # 发生异常
      echo "执行失败!"
      sleep 40
      # 调用API并传递错误信息
      JOB_MESSAGE=$(jq -n --arg content "$JOB_MESSAGE" '{"status": "failed", "message": $content}')
      curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "$JOB_MESSAGE"
  fi
}

URL_REGEX="^(http|https)://"
mkdir -p /data/train-data/
TRAIN_LOCAL_FILE=/data/train-data/train-${JOB_ID}.jsonl
EVAL_LOCAL_FILE=/data/train-data/eval-${JOB_ID}.jsonl

function download_file {
    if [[ "$1" =~ $URL_REGEX ]]; then
        echo "The path is a URL. Starting download..."
        wget -O $2 "$1"
    else
        echo "File path is local. No download needed."
        cp "$1" "$2"
    fi
}

download_file "$TRAIN_FILE" "$TRAIN_LOCAL_FILE"

if [ "$EVAL_FILE" != "" ]; then
  download_file "$EVAL_FILE" "$EVAL_LOCAL_FILE"
fi

temp_file=$(mktemp)

torchrun --nproc_per_node $GPUS_PER_NODE /app/finetune.py \
    --model_name_or_path "$BASE_MODEL_PATH" \
    --data_path "${TRAIN_LOCAL_FILE}" \
    --bf16 True \
    --output_dir "$OUTPUT_DIR" \
    --num_train_epochs $NUM_TRAIN_EPOCHS \
    --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
    --per_device_eval_batch_size $PER_DEVICE_EVAL_BATCH_SIZE \
    --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
    --evaluation_strategy "no" \
    --save_strategy "epoch" \
    --save_total_limit 5 \
    --learning_rate $LEARNING_RATE \
    --weight_decay 0.1 \
    --adam_beta2 0.95 \
    --warmup_ratio 0.01 \
    --lr_scheduler_type "cosine" \
    --logging_steps 1 \
    --report_to "none" \
    --model_max_length $MODEL_MAX_LENGTH \
    --gradient_checkpointing True \
    --lazy_preprocess True \
    --use_lora ${USE_LORA} \
    --q_lora ${Q_LORA} \
    --deepspeed $DS_CONFIG_PATH  > >(tee "$temp_file") 2>&1

JOB_STATUS=$?
JOB_MESSAGE=$(<"$temp_file")
callback