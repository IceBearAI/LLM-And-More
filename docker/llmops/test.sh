#!/bin/bash
clear
ps aux | grep llmops_deepspeed_main.py | awk '{print $2}' | xargs kill -9
ps aux | grep faq_train.py | awk '{print $2}' | xargs kill -9
MNT_PATH=../mnt

# 使用的GPU编号
INCLUDE_GPUS=0,1,2,3,4,5,6,7
export CUDA_VISIBLE_DEVICES=$INCLUDE_GPUS

# 场景选择: "general" 或 "faq"
SCENARIO="general"

# 选择基座模型
MODEL_NAME_OR_PATH="$MNT_PATH/models/Qwen1.5-14B"

echo "INCLUDE_GPUS: $INCLUDE_GPUS"

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
    python jsonl_to_arrow_format.py \
        --base_path "$MNT_PATH"

    ZERO_STAGE=2
    deepspeed llmops_deepspeed_main.py \
        --data_path "$MNT_PATH"/formatted_datasets \
        --data_output_path "$MNT_PATH"/output_path/data_output \
        --data_split 9,1,0 \
        --model_name_or_path "$MODEL_NAME_OR_PATH" \
        --per_device_train_batch_size 2 \
        --per_device_eval_batch_size 10 \
        --max_seq_len 512 \
        --learning_rate 1e-5 \
        --weight_decay 0. \
        --num_train_epochs 2  \
		--train_type "all" \
		--lora_dim 2 \
		--lora_module_name "q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj" \
        --gradient_accumulation_steps 10 \
        --lr_scheduler_type cosine \
        --num_warmup_steps 0 \
        --seed 42 \
        --gradient_checkpointing \
        --zero_stage $ZERO_STAGE \
        --deepspeed \
        --print_loss \
        --output_dir "$MNT_PATH"/output_path \
        --start_from_step -1 \
        --save_per_steps 100 \
        --tensorboard_path "$MNT_PATH"/output_path/tensorboard \
        --tensorboard_port 6007 \
		--enable_tensorboard| tee $MNT_PATH/output_path/log.txt \

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
        --model_name_or_path $MODEL_NAME_OR_PATH \
        --per_device_train_batch_size 1 \
        --max_len 1024 \
        --max_src_len 512 \
        --learning_rate 1e-4 \
        --weight_decay 0.1 \
        --num_train_epochs 2 \
        --gradient_accumulation_steps 16 \
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
        --output_dir $FAQ_OUTPUT_DIR
else
    echo "Invalid scenario selection!"
fi