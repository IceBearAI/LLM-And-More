#!/bin/bash

DATA_TASK_JOB_ID="task-07093b2d-18bf-4439-a094-86acb819ab33"
API_URL="${API_HOST}/api/mgr/annotation/task/${DATA_TASK_JOB_ID}/detect/finish"
AUTH="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJuZCIsInF3VXNlcmlkIjoiIiwiZW1haWwiOiJjb25nd2FuZ0BjcmVkaXRlYXNlLmNuIiwidXNlcklkIjozLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoic3lzdGVtIiwiZXhwIjoxNzEwMjk5NjI0fQ.YgXFAw5Z8MOSWvrk-0GPfDfxxLgisoPfVXeFkvSKYVU"
TENANT_ID="9d65d67f-acd1-47ec-8858-6e57361997b4"

# 模型性能评估脚本变量
# MODEL_NAME_OR_PATH="../../mnt/models/chatglm3-6b-32k"
MODEL_NAME_OR_PATH="../../mnt/models/Qwen1.5-14B"
DATASET_PATH="../../mnt/datasets/min_test.jsonl"
EVALUATION_METRICS="Rouge"
MAX_SEQ_LEN=24
PER_DEVICE_BATCH_SIZE=1
GPU_ID=0
OUTPUT_PATH="../../mnt/output_path"

# 调用Python脚本并捕获输出和退出状态
output=$(python model_performance_evaluation.py \
  --model_name_or_path "${MODEL_NAME_OR_PATH}" \
  --dataset_path "${DATASET_PATH}" \
  --evaluation_metrics "${EVALUATION_METRICS}" \
  --max_seq_len ${MAX_SEQ_LEN} \
  --per_device_batch_size ${PER_DEVICE_BATCH_SIZE} \
  --gpu_id ${GPU_ID} \
  --output_path "${OUTPUT_PATH}" 2>&1)
status=$?

# 根据退出状态判断执行是否异常
if [ $status -eq 0 ]; then
    echo "执行成功，输出内容："
    echo "${output}"
    job_status="success"
    # 调用API并传递输出内容
    json_content=$(jq -n --arg jobStatus "$job_status" --arg output "$output" '{"status": $jobStatus, "results": $output}')
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "${json_content}"
else
    echo "执行失败，错误信息："
    echo "${output}"
    job_status="failed"
    # 调用API并传递错误信息
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"${job_status}\", \"message\": \"${output}\"}"
fi