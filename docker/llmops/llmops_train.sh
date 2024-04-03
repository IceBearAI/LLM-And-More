#!/bin/bash
export CUDA_DEVICE_MAX_CONNECTIONS=1

NNODES=1
NODE_RANK=0
MASTER_ADDR=localhost
MASTER_PORT={{.MasterPort}}
Q_LORA=False

DISTRIBUTED_ARGS="
    --nproc_per_node $GPUS_PER_NODE \
    --nnodes $NNODES \
    --node_rank $NODE_RANK \
    --master_addr $MASTER_ADDR \
    --master_port $MASTER_PORT
"
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
if [ "$USE_LORA" = true ]; then
    TR_TYPE="lora"
    LORA_NAME=' '
    MODENAME=$(echo "$BASE_MODEL" | tr '[:upper:]' '[:lower:]')
    case $MODENAME in
        'llama2')
            LORA_NAME='gate_proj,down_proj,up_proj'
            ;;
        'baichuan2')
            LORA_NAME='W_pack'
            ;;
        'chatglm3')
            LORA_NAME='query_key_value,dense_h_to_4h,dense_4h_to_h,dense'
            ;;
        'glm')
            LORA_NAME='gate_proj,down_proj,up_proj'
            ;;
        'qwen1.5')
            LORA_NAME='q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj'
            ;;
        *)
            echo "未知模型名称"
            ;;
    esac
    ZERO_STAGE=2
else
    TR_TYPE="all"
    ZERO_STAGE=3
fi

if [ $GPUS_PER_NODE -eq 1 ]; then
    # 单卡
    export CUDA_VISIBLE_DEVICES=0
elif [ $GPUS_PER_NODE -eq 2 ]; then
    # 双卡
    export CUDA_VISIBLE_DEVICES=0,1
elif [ $GPUS_PER_NODE -eq 4 ]; then
    # 四卡
    export CUDA_VISIBLE_DEVICES=0,1,2,3
elif [ $GPUS_PER_NODE -eq 8 ]; then
    # 八卡
    export CUDA_VISIBLE_DEVICES=0,1,2,3,4,5,6,7
else
    echo "Invalid GPUS_PER_NODE!"
fi

URL_REGEX="^(http|https)://"

mkdir -p /data/train-data/
TRAIN_LOCAL_FILE=/data/train-data/train-${JOB_ID}.jsonl
EVAL_LOCAL_FILE=/data/train-data/eval-${JOB_ID}.jsonl

if [[ "$TRAIN_FILE" =~ $URL_REGEX ]]; then
    echo "The path is a URL. Starting download..."
    wget -O $TRAIN_LOCAL_FILE "$TRAIN_FILE"
else
    TRAIN_LOCAL_FILE="$TRAIN_FILE"
fi

if [[ "$EVAL_FILE" =~ $URL_REGEX ]]; then
    echo "The path is a URL. Starting download..."
    wget -O $EVAL_LOCAL_FILE "$EVAL_FILE"
else
    EVAL_LOCAL_FILE="$EVAL_FILE"
fi


if [ "$SCENARIO" == "general" ]; then
  mkdir -p /data/train-data/formatted_datasets

  python3 jsonl_to_arrow_format.py \
    --train_path "$TRAIN_LOCAL_FILE" \
    --test_path "$EVAL_LOCAL_FILE" \
    --output_path "/data/train-data/formatted_datasets"

#  output=$(torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
#  output=$(deepspeed --include localhost:$CUDA_VISIBLE_DEVICES {{.ScriptFile}} \
  output=$(deepspeed {{.ScriptFile}} \
      --data_path "/data/train-data/formatted_datasets" \
      --data_output_path {$OUTPUT_DIR}/data_output \
      --data_split 9,1,0 \
      --model_name_or_path "$BASE_MODEL_PATH" \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --per_device_eval_batch_size $PER_DEVICE_EVAL_BATCH_SIZE \
      --max_seq_len {{.ModelMaxLength}} \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0. \
      --num_train_epochs $NUM_TRAIN_EPOCHS  \
      --train_type TR_TYPE \
      --lora_dim 2 \
      --lora_module_name LORA_NAME \
      --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
      --lr_scheduler_type cosine \
      --num_warmup_steps 0 \
      --seed 42 \
      --gradient_checkpointing \
      --zero_stage $ZERO_STAGE \
      --deepspeed \
      --print_loss \
      --output_dir $OUTPUT_DIR \
      --start_from_step -1 \
      --save_per_steps 100 \
      --tensorboard_path "{$OUTPUT_DIR}/tensorboard" \
      --tensorboard_port 6007 \
      --enable_tensorboard)
elif [ "$SCENARIO" == "faq" ]; then


  pass

elif [ "$SCENARIO" == "rag" ]; then

  pass

else
  echo "Invalid scenario selection!"
fi
status=$?

# 根据退出状态判断执行是否异常
if [ $status -eq 0 ]; then
    # 没有发生异常，正常输出内容
    echo "执行成功!"
    echo "${API_URL}"
    # 调用API并传递输出内容
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{"status": "success"}"
else
    # 发生异常
    echo "执行失败!"
    # 调用API并传递错误信息
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": "failed", \"message\": \"${output}\"}"
fi
