#!/bin/bash
export AUTH=sk-001
export JOB_ID={{.JobId}}

export CUDA_DEVICE_MAX_CONNECTIONS=1

GPUS_PER_NODE={{.ProcPerNode}}

# 根据GPUS_PER_NODE数量设置CUDA_VISIBLE_DEVICES
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

MODEL="{{.BaseModelPath}}" # Set the path if you do not want to load from huggingface directly
# ATTENTION: specify the path to your training data, which should be a json file consisting of a list of conversations.
# See the section for finetuning in README for more information.
DATA="{{.DataPath}}"

mkdir -p /data/train-data/formatted_datasets
wget -O /data/train-data/formatted_datasets/{{.DataPath}} {{.FileUrl}}

# 通过设置 `--train_type "lora"`开启 LoRA 微调（默认开启）。
# 1. 若需进行全量微调，可设置为 `--train_type "all"`。
# 2.此外，`--lora_module_name "gate_proj,down_proj,up_proj"`允许指定模型中要应用 LoRA 的特定模块，确保所选模块与模型结构相匹配。
# 不同模型的LoRA模块lora_module_name选择参考:
#     'Llama2': 'gate_proj,down_proj,up_proj',
#     'Baichuan2': 'W_pack',
#     'ChatGLM3': 'query_key_value,dense_h_to_4h,dense_4h_to_h,dense',
#     'GLM': "gate_proj,down_proj,up_proj",
#     'Qwen1.5': "q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj",

if [ "$SCENARIO" == "general" ]; then
    # 一般场景
    python3 jsonl_to_arrow_format.py \
        --base_path "$MNT_PATH"

    ZERO_STAGE=2
    deepspeed llmops_deepspeed_main.py \
        --data_path "/data/train-data/formatted_datasets/${DATA}" \
        --data_output_path "$MNT_PATH"/output_path/data_output \
        --data_split 9,1,0 \
        --model_name_or_path "${MODEL}" \
        --per_device_train_batch_size {{.TrainBatchSize}} \
        --per_device_eval_batch_size {{.EvalBatchSize}} \
        --max_seq_len {{.ModelMaxLength}} \
        --learning_rate {{.LearningRate}} \
        --weight_decay 0. \
        --num_train_epochs {{.TrainEpoch}}  \
        --train_type "all" \
        --lora_dim 2 \
        --lora_module_name "q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj" \
        --gradient_accumulation_steps {{.AccumulationSteps}} \
        --lr_scheduler_type cosine \
        --num_warmup_steps 0 \
        --seed 42 \
        --gradient_checkpointing \
        --zero_stage $ZERO_STAGE \
        --deepspeed \
        --print_loss \
        --output_dir {{.OutputDir}} \
        --start_from_step -1 \
        --save_per_steps 100 \
        --tensorboard_path "{{.OutputDir}}/output_path/tensorboard" \
        --tensorboard_port 6007 \
        --enable_tensorboard

elif [ "$SCENARIO" == "faq" ]; then
    # FAQ场景参数配置
    FAQ_OUTPUT_DIR="./faq/output_model"
    FAQ_TRAIN_PATH="./faq/data/mini_finance_train.json"
    # 模式选择: "auto"、"glm"、"glm2"、"glm3"、"baichuan2_13b"、"qwen1.5"
      # "auto"：自动选择模式，适用于`Llama2-13b`模型
    FAQ_MODE="qwen1.5"

    # FAQ场景
    deepspeed ./faq/faq_train.py \
        --train_path $FAQ_TRAIN_PATH \
        --model_name_or_path ${MODEL} \
        --per_device_train_batch_size {{.TrainBatchSize}} \
        --max_len 1024 \
        --max_src_len 512 \
        --learning_rate {{.LearningRate}} \
        --weight_decay 0.1 \
        --num_train_epochs {{.TrainEpoch}} \
        --gradient_accumulation_steps {{.AccumulationSteps}} \
        --warmup_ratio 0.1 \
        --mode $FAQ_MODE \
        --train_type "all" \
        --lora_module_name "q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj" \
        --lora_dim 4 \
        --lora_alpha 64 \
        --lora_dropout 0.1 \
        --seed 1234 \
        --ds_file ./faq/ds_zero2_no_offload.json \
        --gradient_checkpointing \
        --show_loss_step 10 \
        --output_dir {{.OutputDir}}
else
  echo "Invalid scenario selection!"
fi