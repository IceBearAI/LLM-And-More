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
export BASE_MODEL_PATH="THUDM/chatglm3-6b-32k"
export BASE_MODEL_NAME="chatglm3-6b-32k"
# 场景选择: "general" 或 "faq"
export SCENARIO="faq"

#!/bin/bash
export CUDA_DEVICE_MAX_CONNECTIONS=1

NNODES=1
NODE_RANK=0
MASTER_ADDR=localhost
MASTER_PORT=54648
Q_LORA=False
export HF_HOME=/root/autodl-tmp/
export HF_ENDPOINT=https://hf-mirror.com
GPUS_PER_NODE=2

DISTRIBUTED_ARGS="
    --nproc_per_node $GPUS_PER_NODE \
    --nnodes $NNODES \
    --node_rank $NODE_RANK \
    --master_addr $MASTER_ADDR \
    --master_port $MASTER_PORT
"

LORA_NAME=''
MODENAME=$(echo "chatglm3-6b-32k" | tr '[:upper:]' '[:lower:]')

case $MODENAME in
    *'llama2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='llama2'
        echo $MODENAME
        ;;
    *'baichuan2'*)
        LORA_MODULE_NAME='W_pack'
        MODENAME='baichuan2_13b'
        echo $MODENAME
        ;;
    *'chatglm3'*)
        LORA_MODULE_NAME='query_key_value,dense_h_to_4h,dense_4h_to_h,dense'
        MODENAME='glm3'
        echo $MODENAME
        ;;
    *'glm2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm2'
        echo $MODENAME
        ;;
    *'glm'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm'
        echo $MODENAME
        ;;
    *'qwen1.5'*)
        LORA_MODULE_NAME='q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj'
        MODENAME='qwen1.5'
        echo $MODENAME
        ;;
    *)
        MODENAME='auto'
        echo "未知模型名称"
        ;;
esac


USE_LORA=true
if [ "$USE_LORA" = true ]; then
    TRAIN_TYPE="lora"
    ZERO_STAGE=2
    DS_FILE=faq/ds_zero2_no_offload.json
    echo "使用LoRA微调"
    echo "ZERO_STAGE: $ZERO_STAGE"
else
    TRAIN_TYPE="all"
    ZERO_STAGE=3
    DS_FILE=faq/ds_zero3_offload.json
    echo "使用全量微调"
    echo "ZERO_STAGE: $ZERO_STAGE"
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
JOB_ID=ssss
mkdir -p /data/train-data/
# TRAIN_LOCAL_FILE=/data/train-data/train-${JOB_ID}.jsonl
TRAIN_LOCAL_FILE="./datasets_example/faq_train_dataset.jsonl"
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

SCENARIO="faq"
echo "SCENARIO: $SCENARIO"
if [ "$SCENARIO" == "general" ]; then
  mkdir -p /data/train-data/formatted_datasets

  python3 jsonl_to_arrow_format.py \
    --train_path "./general_train_dataset.jsonl" \
    --output_path "./data/train-data/formatted_datasets"

#  output=$(torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
#  output=$(deepspeed --include localhost:$CUDA_VISIBLE_DEVICES {{.ScriptFile}} \
  output=$(deepspeed llmops_deepspeed_main.py \
      --data_path "./data/train-data/formatted_datasets" \
      --data_output_path ./data/data_output \
      --data_split 9,1,0 \
      --model_name_or_path THUDM/chatglm3-6b-32k \
      --per_device_train_batch_size 1 \
      --per_device_eval_batch_size 8 \
      --max_seq_len 512 \
      --learning_rate 0.00001 \
      --weight_decay 0. \
      --num_train_epochs 1  \
      --train_type TRAIN_TYPE \
      --lora_dim 2 \
      --lora_module_name LORA_MODULE_NAME \
      --gradient_accumulation_steps 8 \
      --lr_scheduler_type cosine \
      --num_warmup_steps 0 \
      --seed 42 \
      --gradient_checkpointing \
      --zero_stage $ZERO_STAGE \
      --deepspeed \
      --print_loss \
      --output_dir ./output \
      --start_from_step -1 \
      --save_per_steps 100 \
      --tensorboard_path ./output/tensorboard \
      --tensorboard_port 6007 \
      --enable_tensorboard)
elif [ "$SCENARIO" == "faq" ]; then
    formatted_datasets_path="./data/train-data/faq_formatted_datasets"
    mkdir -p "$formatted_datasets_path"
    TRAIN_LOCAL_FILE="./datasets_example/faq_train_dataset.jsonl"
    DS_FILE="./faq/ds_zero2_no_offload.json"

    python3 ./faq/convert_new_format.py \
        --train_path $TRAIN_LOCAL_FILE \
        --output_path "$formatted_datasets_path"

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
        --train_type "$TRAIN_TYPE" \
        --lora_module_name $LORA_MODULE_NAME \
        --lora_dim 4 \
        --lora_alpha 64 \
        --lora_dropout 0.1 \
        --seed 1234 \
        --ds_file $DS_FILE \
        --gradient_checkpointing \
        --show_loss_step 10 \
        --output_dir $OUTPUT_DIR

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