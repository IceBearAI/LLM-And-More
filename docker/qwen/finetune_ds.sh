#!/bin/bash
export AUTH=
export JOB_ID={{.JobId}}
export CHAT_API_FINE_TUNE_URL=https://chat-api:8080/v1/fine_tuning/jobs/${JOB_ID}/finish

export CUDA_DEVICE_MAX_CONNECTIONS=1
DIR=`pwd`

GPUS_PER_NODE={{.ProcPerNode}}
NNODES=1
NODE_RANK=0
MASTER_ADDR=localhost
MASTER_PORT={{.MasterPort}}

MODEL="{{.BaseModelPath}}" # Set the path if you do not want to load from huggingface directly
# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.
# See the section for finetuning in README for more information.
DATA="{{.DataPath}}"
DS_CONFIG_PATH="ds_config_zero3.json"

DISTRIBUTED_ARGS="
    --nproc_per_node $GPUS_PER_NODE \
    --nnodes $NNODES \
    --node_rank $NODE_RANK \
    --master_addr $MASTER_ADDR \
    --master_port $MASTER_PORT
"

mkdir -p /data/train-data/
wget -O {{.DataPath}} {{.FileUrl}}

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
    --learning_rate 1e-5 \
    --weight_decay 0.1 \
    --adam_beta2 0.95 \
    --warmup_ratio 0.01 \
    --lr_scheduler_type "cosine" \
    --logging_steps 1 \
    --report_to "none" \
    --model_max_length {{.ModelMaxLength}} \
    --gradient_checkpointing True \
    --lazy_preprocess True \
    --deepspeed $DS_CONFIG_PATH
