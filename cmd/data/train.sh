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

JOB_STATUS="success"
JOB_MESSAGE=""

if [ "$MASTER_PORT" == "" ]; then
    MASTER_PORT=29500
fi

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
      # 调用API并传递错误信息
      sleep 40
      JOB_MESSAGE=$(jq -n --arg content "$JOB_MESSAGE" '{"status": "failed", "message": $content}')
      curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "$JOB_MESSAGE"
  fi
}

function set_cuda_devices {
    # 验证输入是否为空
    if [ -z "$1" ]; then
        echo "Error: No argument provided."
        echo "Usage: set_cuda_devices <number_of_gpus>"
        exit 1
    fi

    if ! [[ "$1" =~ ^[0-9]+$ ]]; then
        echo "Error: Invalid argument. Please provide a positive integer."
        exit 1
    fi

    # 使用循环动态构建
    devices=$(printf ",%d" $(seq 0 $(($1-1)) | sed 's/ //g'))
    export CUDA_VISIBLE_DEVICES=${devices:1}
}
set_cuda_devices $GPUS_PER_NODE


IS_MERGE_LORA=true
merge_lora_models() {
    # 参数检查
    if [ "$#" -ne 4 ]; then
        echo "Usage: merge_lora_models <merge_lora_script> <base_model_path> <output_dir> <modename>"
        return 1
    fi

    merge_lora_script="$1"
    base_model_path="$2"
    output_dir="$3"
    modename="$4"
    # 执行 merge_lora.py 脚本
    python "$merge_lora_script" --ori_model_dir "$base_model_path" --model_dir "$output_dir" --mode "$modename"
    # 复制 token 文件和 config.json 到输出目录
    cp  "$base_model_path/config"* "./$output_dir/"
    cp  "$base_model_path/config.json" "./$output_dir/"
    cp  "$base_model_path/token"* "./$output_dir/"
    set +e
    if [[ $modename == *"glm3"* ]]; then
        cp  "$base_model_path/modeling_chatglm.py" "./$output_dir/"
        cp  "$base_model_path/quantization.py" "./$output_dir/"
    fi
    set -e
}

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
    *'glm3_32k'*)
        LORA_MODULE_NAME='query_key_value,dense_h_to_4h,dense_4h_to_h,dense'
        MODENAME='glm3_32k'
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
DATA_DIR="/data/train-data/"
mkdir -p "$DATA_DIR"
if [[ -f "$TRAIN_FILE" || "$TRAIN_FILE" =~ $URL_REGEX ]]; then
    TRAIN_LOCAL_FILE="${DATA_DIR}train-${JOB_ID}.jsonl"
else
    TRAIN_LOCAL_FILE="${DATA_DIR}train_dir_${JOB_ID}"
    mkdir -p "$TRAIN_LOCAL_FILE"
fi
EVAL_LOCAL_FILE="${DATA_DIR}eval-${JOB_ID}.jsonl"

function download_file {
    local SOURCE="$1"
    local DEST="$2"
    if [[ -d "$SOURCE" ]]; then
        echo "Source is a directory. Copying files..."
        cp -r "$SOURCE"/* "$DEST"
    elif [[ -f "$SOURCE" ]]; then
        echo "Source is a file. Copying file..."
        cp "$SOURCE" "$DEST"
    elif [[ "$SOURCE" =~ $URL_REGEX ]]; then
        echo "Source is a URL. Starting download..."
        wget -O "$DEST" "$SOURCE"
    else
        echo "Invalid source provided."
        return 1
    fi
}

download_file "$TRAIN_FILE" "$TRAIN_LOCAL_FILE"

if [ "$EVAL_FILE" != "" ]; then
  download_file "$EVAL_FILE" "$EVAL_LOCAL_FILE"
fi

temp_file=$(mktemp)

if [ "$SCENARIO" == "general" ]; then
  GENERAL_DATA_PATH=/data/train-data/formatted_datasets/${JOB_ID}
  mkdir -p $GENERAL_DATA_PATH
  if [ -n "$EVAL_FILE" ]; then
    python3 jsonl_to_arrow_format.py \
        --train_path "$TRAIN_LOCAL_FILE" \
        --test_path "$EVAL_LOCAL_FILE" \
        --output_path "$GENERAL_DATA_PATH"
  else
    python3 jsonl_to_arrow_format.py \
        --train_path "$TRAIN_LOCAL_FILE" \
        --output_path "$GENERAL_DATA_PATH"
  fi

  deepspeed /app/llmops_deepspeed_main.py \
      --data_path $GENERAL_DATA_PATH \
      --data_output_path $OUTPUT_DIR/data_output \
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
      --seed 1234 \
      --gradient_checkpointing \
      --zero_stage $ZERO_STAGE \
      --offload \
      --only_optimize_lora \
      --deepspeed \
      --output_dir $OUTPUT_DIR \
      --start_from_step -1 \
      --save_per_steps 100  > >(tee "$temp_file") 2>&1
elif [ "$SCENARIO" == "faq" ]; then
  formatted_datasets_path=/data/train-data/faq_formatted_datasets
  mkdir -p "$formatted_datasets_path"

  python3 /app/convert_new_format.py \
      --train_path $TRAIN_LOCAL_FILE \
      --test_path $EVAL_LOCAL_FILE \
      --output_path "$formatted_datasets_path"

  deepspeed /app/faq/faq_train.py \
      --train_path "$formatted_datasets_path/train_dataset.jsonl" \
      --model_name_or_path $BASE_MODEL_PATH \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --max_len $MODEL_MAX_LENGTH \
      --max_src_len 128 \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0.1 \
      --num_train_epochs $NUM_TRAIN_EPOCHS \
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
      --output_dir $OUTPUT_DIR  > >(tee "$temp_file") 2>&1

  if [ "$USE_LORA" = true ] && [ "$IS_MERGE_LORA" = true ]; then
    merge_lora_models "./faq/merge_lora.py" "$BASE_MODEL_PATH" "$OUTPUT_DIR" "$MODENAME"
  fi

elif [ "$SCENARIO" == "rag" ]; then
    RAG_ENHANCEMENT=False
#     RETRIEVAL_METHOD="bm25"
    RETRIEVAL_METHOD="st"
    ST="BAAI/bge-large-zh-v1.5"

    deepspeed /app/rag/rag_train.py \
        --train_path $TRAIN_LOCAL_FILE \
        --enhancement $RAG_ENHANCEMENT \
        --merge_lora $IS_MERGE_LORA \
        --retrieval_method $RETRIEVAL_METHOD \
        --top_k 1 \
        --st $ST \
        --model_name_or_path $BASE_MODEL_PATH \
        --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
        --per_device_eval_batch_size $PER_DEVICE_EVAL_BATCH_SIZE \
        --max_len $MODEL_MAX_LENGTH \
        --max_src_len 1024 \
        --learning_rate $LEARNING_RATE \
        --weight_decay 0.1 \
        --num_train_epochs $NUM_TRAIN_EPOCHS \
        --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
        --warmup_ratio 0.1 \
        --mode $MODENAME \
        --train_type $TRAIN_TYPE \
        --load_in_kbits 16 \
        --fp16 \
        --lora_module_name  $LORA_MODULE_NAME \
        --lora_dim 4 \
        --lora_alpha 64 \
        --lora_dropout 0.1 \
        --seed 1234 \
        --ds_file $DS_FILE \
        --gradient_checkpointing \
        --show_loss_step 10 \
        --output_dir $OUTPUT_DIR  > >(tee "$temp_file") 2>&1
    set +e
    if [[ $MODENAME == *"glm3"* ]]; then
        cp "$BASE_MODEL_PATH/modeling_chatglm.py" "./$OUTPUT_DIR" || cp "$BASE_MODEL_PATH/quantization.py" "./$OUTPUT_DIR"
    fi
    set -e
else
  echo "Invalid scenario selection!"
fi

JOB_STATUS=$?
JOB_MESSAGE=$(<"$temp_file")
callback