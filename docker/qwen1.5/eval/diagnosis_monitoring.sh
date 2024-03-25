#!/bin/bash

DATA_TASK_JOB_ID="task-07093b2d-18bf-4439-a094-86acb819ab33"
API_URL="${API_HOST}/api/mgr/annotation/task/${DATA_TASK_JOB_ID}/detect/finish"
AUTH="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJuZCIsInF3VXNlcmlkIjoiIiwiZW1haWwiOiJjb25nd2FuZ0BjcmVkaXRlYXNlLmNuIiwidXNlcklkIjozLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoic3lzdGVtIiwiZXhwIjoxNzEwMjk5NjI0fQ.YgXFAw5Z8MOSWvrk-0GPfDfxxLgisoPfVXeFkvSKYVU"
TENANT_ID="9d65d67f-acd1-47ec-8858-6e57361997b4"

# 日志监控脚本变量
LOG_PATH="../../mnt/output_path/log.txt"

# 调用Python脚本并捕获输出和退出状态
output=$(python diagnosis_monitoring.py --log_path="${LOG_PATH}" 2>&1)
status=$?

# 根据退出状态判断执行是否异常
if [ $status -eq 0 ]; then
    echo "执行成功，输出内容："
    echo "${output}"
    job_status="success"
    # 生成要发送的 JSON 数据
    json_content=$(jq -n --arg output "$output" --arg jobStatus "$job_status" '{"status": $jobStatus, "data": $output}')
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "${json_content}"
else
    echo "执行失败，错误信息："
    echo "${output}"
    job_status="failed"
    # 发送失败状态和错误信息
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"${job_status}\", \"message\": \"${output}\"}"
fi