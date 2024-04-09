#!/bin/bash

#API_URL="${API_HOST}/api/mgr/annotation/task/${DATA_TASK_JOB_ID}/detect/finish"
DATASET_FILE="/app/dataset.json"
DATASET_OUTPUT_FILE="/app/${DATA_TASK_JOB_ID}-result.json"

cd /app/eval/

# 判断如果 DATASET_PATH 是url，则下载文件
URL_REGEX="^(http|https)://"

if [[ $DATASET_PATH =~ $URL_REGEX ]]; then
    echo "The path is a URL. Starting download..."
    wget -O $DATASET_FILE "$DATASET_PATH"
else
    DATASET_FILE=$DATASET_PATH
fi


# 调用Python脚本并捕获输出和退出状态
output=$(python3 analyze_similar_questions_and_intents.py \
  --model_name ${DATASET_ANALYZE_MODEL} \
  --similarity_threshold 0.91 \
  --intent_similarity_threshold 0.86 \
  --dataset ${DATASET_FILE} \
  --output_file ${DATASET_OUTPUT_FILE} \
  --dataset_type ${DATASET_TYPE} 2>&1)
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
    curl -X PUT ${API_URL} -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"${job_status}\", \"message\": \"${output}\"}"
fi

rm -rf $DATASET_OUTPUT_FILE