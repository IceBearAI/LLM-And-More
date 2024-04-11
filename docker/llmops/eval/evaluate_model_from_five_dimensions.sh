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
# - DATASET_PATH: 数据集路径
# - EVALUATION_METRICS: 评估指标
# - MAX_SEQ_LEN: 最大长度
# - PER_DEVICE_BATCH_SIZE: 评估批次

# Python脚本执行的参数
EVAL_DIMENSIONS='["inference_ability", "reading_comprehension", "chinese_language_skill", "command_compliance", "innovation_capacity"]'
OPTIONS='{"additional_parameters": {"max_seq_len": 512}}'
DATASET_OUTPUT_FILE="/app/eval/eval-model-result.json"

cd /app/eval/

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

temp_file=$(mktemp)

# 执行 Python 脚本并捕获输出和退出状态
deepspeed  evaluate_model_from_five_dimensions.py \
  --model_name_or_path "${MODEL_PATH}" \
  --gpu_nums "$GPUS_PER_NODE" \
  --evaluation_dimensions "${EVAL_DIMENSIONS}" \
  --output_file ${DATASET_OUTPUT_FILE} \
  --options "${OPTIONS}" > >(tee "$temp_file") 2>&1
status=$?

output=$(<"$temp_file")
output=$(echo "$output" | sed 's/"/\\"/g')


# 根据退出状态判断执行是否异常
if [ $status -eq 0 ]; then
    # 没有发生异常，正常输出内容
    echo "执行成功，输出内容："
    job_status="success"
    # 调用API并传递输出内容
#    content=$(<"${DATASET_OUTPUT_FILE}")
    json_content=$(jq -c '.' "$DATASET_OUTPUT_FILE")
    new_json=$(jq -n --argjson content "$json_content" '{"status": "success", "data": $content}')
    curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "${new_json}"
else
    # 发生异常
    echo "执行失败，错误信息："
    job_status="failed"
    # 调用API并传递错误信息
    curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"${job_status}\", \"message\": \"${output}\"}"
fi

rm -rf $DATASET_OUTPUT_FILE