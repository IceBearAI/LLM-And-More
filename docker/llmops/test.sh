ps aux | grep llmops_deepspeed_main.py | awk '{print $2}' | xargs kill -9
ps aux | grep faq_train.py | awk '{print $2}' | xargs kill -9
clear
export OUTPUT_DIR="./faq/output_model"
export NUM_TRAIN_EPOCHS=2
export PER_DEVICE_TRAIN_BATCH_SIZE=1
export PER_DEVICE_EVAL_BATCH_SIZE=4
export GRADIENT_ACCUMULATION_STEPS=2
export LEARNING_RATE=0.0001
export MODEL_MAX_LENGTH=256
# export BASE_MODEL_PATH="/home/ubuntu/model/THUDM/chatglm3-6b-32k"
# export BASE_MODEL_PATH="/home/ubuntu/model/THUDM/chatglm3-6b"
export BASE_MODEL_PATH="/home/ubuntu/model/Qwen/Qwen1.5-4B"
export BASE_MODEL_NAME="Qwen1.5-4B"
GPUS_PER_NODE=4
SCENARIO="faq"
USE_LORA=true
# 场景选择: "general" 或 "faq"

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
echo "TRAIN_TYPE: $TRAIN_TYPE"
echo "ZERO_STAGE: $ZERO_STAGE"
echo "DS_FILE: $DS_FILE"

URL_REGEX="^(http|https)://"
TRAIN_LOCAL_FILE="./datasets_example/faq_train_dataset.jsonl"
EVAL_LOCAL_FILE="./datasets_example/faq_train_dataset.jsonl"

echo "TRAIN_LOCAL_FILE: $TRAIN_LOCAL_FILE"
echo "EVAL_LOCAL_FILE: $EVAL_LOCAL_FILE"

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


#  output=$(torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
#  output=$(deepspeed --include localhost:$CUDA_VISIBLE_DEVICES {{.ScriptFile}} \
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
      --seed 42 \
      --gradient_checkpointing \
      --zero_stage $ZERO_STAGE \
      --offload \
      --only_optimize_lora \
      --deepspeed \
      --output_dir $OUTPUT_DIR \
      --start_from_step -1 \
      --save_per_steps 100  > >(tee "$temp_file") 2>&1
elif [ "$SCENARIO" == "faq" ]; then
  formatted_datasets_path=./data/train-data/faq_formatted_datasets
  mkdir -p "$formatted_datasets_path"

  python3 ./convert_new_format.py \
      --train_path $TRAIN_LOCAL_FILE \
      --test_path $EVAL_LOCAL_FILE \
      --output_path "$formatted_datasets_path"

#  output=$(deepspeed {{.ScriptFile}}  \
  deepspeed ./faq/faq_train.py \
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
      --output_dir $OUTPUT_DIR  > >(tee "$temp_file") 2>&1

    # 设置是否自动合并LoRA模型变量
    IS_MERGE_LORA=true
    # 需要同时满足USE_LORA和IS_MERGE_LORA为true
    if [ "$USE_LORA" = true ] && [ "$IS_MERGE_LORA" = true ]; then
        # 执行merge_lora.py脚本
        python ./faq/merge_lora.py --ori_model_dir "$BASE_MODEL_PATH" --model_dir "$OUTPUT_DIR" --mode "$MODENAME"

        # 复制token文件和config.json到输出目录
        cp  "$BASE_MODEL_PATH/config"* "./$OUTPUT_DIR/"
        cp  "$BASE_MODEL_PATH/config.json" "./$OUTPUT_DIR/"
        cp  "$BASE_MODEL_PATH/token"* "./$OUTPUT_DIR/"

        # 如果$MODENAME中包含"glm3"
        if [[ $MODENAME == *"glm3"* ]]; then
            cp  "$BASE_MODEL_PATH/modeling_chatglm.py" "./$OUTPUT_DIR/"
            cp  "$BASE_MODEL_PATH/quantization.py" "./$OUTPUT_DIR/"
        fi
    fi



elif [ "$SCENARIO" == "rag" ]; then

  pass

else
  echo "Invalid scenario selection!"
fi

JOB_STATUS=$?
JOB_MESSAGE=$(<"$temp_file")
callback