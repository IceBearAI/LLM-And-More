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
function set_cuda_devices {
    case $1 in
        1) export CUDA_VISIBLE_DEVICES=0 ;;
        2) export CUDA_VISIBLE_DEVICES=0,1 ;;
        4) export CUDA_VISIBLE_DEVICES=0,1,2,3 ;;
        8) export CUDA_VISIBLE_DEVICES=0,1,2,3,4,5,6,7 ;;
        *) echo "Invalid GPUS_PER_NODE!" ; exit 1 ;;
    esac
}
set_cuda_devices $GPUS_PER_NODE


LORA_MODULE_NAME=''
MODENAME=$(echo "$BASE_MODEL_NAME" | tr '[:upper:]' '[:lower:]')
case $MODENAME in
    *'llama2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='llama2'
        ;;
    *'baichuan2'*)
        LORA_MODULE_NAME='W_pack'
        MODENAME='baichuan2_13b'
        ;;
    *'glm3'*)
        LORA_MODULE_NAME='query_key_value,dense_h_to_4h,dense_4h_to_h,dense'
        MODENAME='glm3'
        ;;
    *'glm2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm2'
        ;;
    *'glm'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm'
        ;;
    *'qwen1.5'*)
        LORA_MODULE_NAME='q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj'
        MODENAME='qwen1.5'
        ;;
    *)
        MODENAME='auto'
        echo "未知模型名称"
        ;;
esac

if [ "$USE_LORA" = true ]; then
    TRAIN_TYPE="lora"
    ZERO_STAGE=2
    DS_FILE=faq/ds_zero2_no_offload.json
else
    TRAIN_TYPE="all"
    ZERO_STAGE=3
    DS_FILE=faq/ds_zero3_offload.json
fi

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
download_file "$EVAL_FILE" "$EVAL_LOCAL_FILE"

if [ "$SCENARIO" == "general" ]; then
  GENERAL_DATA_PATH=/data/train-data/formatted_datasets
  mkdir -p $GENERAL_DATA_PATH

  python3 jsonl_to_arrow_format.py \
    --train_path "$TRAIN_LOCAL_FILE" \
    --test_path "$EVAL_LOCAL_FILE" \
    --output_path "$GENERAL_DATA_PATH"

#  output=$(torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
#  output=$(deepspeed --include localhost:$CUDA_VISIBLE_DEVICES {{.ScriptFile}} \
  output=$(deepspeed {{.ScriptFile}} \
      --data_path $GENERAL_DATA_PATH \
      --data_output_path {$OUTPUT_DIR}/data_output \
      --data_split 9,1,0 \
      --model_name_or_path $BASE_MODEL_PATH \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --per_device_eval_batch_size $PER_DEVICE_EVAL_BATCH_SIZE \
      --max_seq_len $MODEL_MAX_LENGTH \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0. \
      --num_train_epochs $NUM_TRAIN_EPOCHS  \
      --train_type $TRAIN_TYPE \
      --lora_dim 2 \
      --lora_module_name $LORA_MODULE_NAME \
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
  formatted_datasets_path=/data/train-data/faq_formatted_datasets
  mkdir -p "$formatted_datasets_path"

  python3 .convert_new_format.py \
      --train_path $TRAIN_LOCAL_FILE \
      --test_path $EVAL_LOCAL_FILE \
      --output_path "$formatted_datasets_path"

#  output=$(deepspeed {{.ScriptFile}}  \
  output=$(deepspeed ./faq/faq_train.py \
      --train_path "$formatted_datasets_path/train_dataset.jsonl" \
      --model_name_or_path $BASE_MODEL_PATH \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --max_len $MODEL_MAX_LENGTH \
      --max_src_len 128 \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0.1 \
      --num_train_epochs 2 \
      --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
      --warmup_ratio 0.1 \
      --mode $MODENAME \
      --train_type $TRAIN_TYPE \
      --lora_module_name $LORA_MODULE_NAME \
      --lora_dim 4 \
      --lora_alpha 64 \
      --lora_dropout 0.1 \
      --seed 1234 \
      --ds_file $DS_FILE \
      --gradient_checkpointing \
      --show_loss_step 10 \
      --output_dir $OUTPUT_DIR)

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
