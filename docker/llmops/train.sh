#!/bin/bash
export AUTH=sk-001
export JOB_ID={{.JobId}}

export CUDA_DEVICE_MAX_CONNECTIONS=1
DIR=`pwd`

GPUS_PER_NODE={{.ProcPerNode}}
NNODES=1
NODE_RANK=0
MASTER_ADDR=localhost
MASTER_PORT=6001
USE_LORA={{.Lora}}
#USE_LORA=1

MODEL="{{.BaseModelPath}}" # Set the path if you do not want to load from huggingface directly
# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.
# See the section for finetuning in README for more information.
DATA="{{.DataPath}}"

DISTRIBUTED_ARGS="
    --nproc_per_node $GPUS_PER_NODE \
    --nnodes $NNODES \
    --node_rank $NODE_RANK \
    --master_addr $MASTER_ADDR \
    --master_port $MASTER_PORT
"

mkdir -p /data/train-data/
wget -O {{.DataPath}} {{.FileUrl}}

if [ "$USE_LORA" -eq true ]; then
    USE_LORA=True
    DS_CONFIG_PATH="ds_config_zero2.json"
else
    USE_LORA=False
    DS_CONFIG_PATH="ds_config_zero3.json"
fi

torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
    --model_name_or_path $MODEL \
    --data_path $DATA \
    --bf16 True \
    --output_dir {{.OutputDir}} \
    --num_train_epochs {{.TrainEpoch}} \
    --per_device_train_batch_size {{.TrainBatchSize}} \
    --per_device_eval_batch_size {{.EvalBatchSize}} \
    --gradient_accumulation_steps {{.AccumulationSteps}} \
    --evaluation_strategy "no" \
    --save_strategy "steps" \
    --save_steps 1000 \
    --save_total_limit 10 \
    --learning_rate {{.LearningRate}} \
    --weight_decay 0.1 \
    --adam_beta2 0.95 \
    --warmup_ratio 0.01 \
    --lr_scheduler_type "cosine" \
    --logging_steps 1 \
    --report_to "none" \
    --model_max_length {{.ModelMaxLength}} \
    --lazy_preprocess True \
    --use_lora ${USE_LORA} \
    --gradient_checkpointing \
    --deepspeed ${DS_CONFIG_PATH}
