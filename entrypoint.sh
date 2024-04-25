#!/bin/sh

nohup python3 -m fastchat.serve.controller > /tmp/controller.log 2>&1 &
nohup python3 -m fastchat.serve.model_worker --model-path $MODEL_PATH --model-name $MODEL_NAME > /tmp/worker.sh 2>&1 &

# shellcheck disable=SC2046
# shellcheck disable=SC2006
while [ `grep -c "Uvicorn running on" /tmp/worker.sh` -eq '0' ];do
        sleep 1s;
        echo "wait worker running"
done

python3 -m fastchat.serve.gradio_web_server