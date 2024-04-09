#!/bin/bash

DATASET_FILE="/app/dataset-${JOB_ID}-eval.jsonl"
OUTPUT_PATH="/app/eval"
DATASET_OUTPUT_FILE="${OUTPUT_PATH}/eval_results.json"

cd /app/eval/

# 判断如果 DATASET_PATH 是url，则下载文件
URL_REGEX="^(http|https)://"

if [[ $DATASET_PATH =~ $URL_REGEX ]]; then
    echo "The path is a URL. Starting download..."
    wget -O $DATASET_FILE "$DATASET_PATH"
else
    DATASET_FILE=$DATASET_PATH
fi

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
    CUDA_VISIBLE_DEVICES=${devices:1}
#     export CUDA_VISIBLE_DEVICES=${devices:1}
}
set_cuda_devices $GPUS_PER_NODE

# 调用Python脚本并捕获输出和退出状态
output=$(python3 model_performance_evaluation.py \
  --model_name_or_path "${MODEL_PATH}" \
  --dataset_path "${DATASET_FILE}" \
  --evaluation_metrics "${EVALUATION_METRICS}" \
  --max_seq_len ${MAX_SEQ_LEN} \
  --per_device_batch_size ${PER_DEVICE_BATCH_SIZE} \
  --gpu_id $CUDA_VISIBLE_DEVICES \
  --output_path "${OUTPUT_PATH}" 2>&1)
status=$?

# 根据退出状态判断执行是否异常
if [ $status -eq 0 ]; then
    # 没有发生异常，正常输出内容
    echo "执行成功，输出内容："
    echo "${output}"
    job_status="success"
    # 调用API并传递输出内容
#    content=$(<"${DATASET_OUTPUT_FILE}")
    json_content=$(jq -c '.' "$DATASET_OUTPUT_FILE")
    new_json=$(jq -n --argjson content "$json_content" '{"status": "success", "data": $content}')
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "${new_json}"
else
    # 发生异常
    echo "执行失败，错误信息："
    echo "${output}"
    job_status="failed"
    # 调用API并传递错误信息
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": "${job_status}", \"message\": \"${output}\"}"
fi

rm -rf $DATASET_OUTPUT_FILE