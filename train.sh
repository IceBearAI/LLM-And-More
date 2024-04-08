#!/bin/bash

BASE_MODEL_PATH=/data/base-model/qwen1-5-0-5b
BASE_MODEL_NAME=qwen1.5-0.5b

if [ -n "$BASE_MODEL_PATH" ]; then
	BASE_MODEL=$BASE_MODEL_PATH
else
	BASE_MODEL=$BASE_MODEL_NAME
fi

echo $BASE_MODEL