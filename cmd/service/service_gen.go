package service

import (
	"fmt"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

var (
	shellStart = `#!/bin/bash

# 自动注入以下环境变量，可以直接使用
# - MODEL_PATH: 模型路径
# - MODEL_NAME: 模型名称
# - HTTP_PORT: HTTP端口
# - QUANTIZATION: 量化 8bit
# - NUM_GPUS: GPU数量
# - MAX_GPU_MEMORY: GPU内存
# - USE_VLLM: 是否是VLLM
# - INFERRED_TYPE: 推理类型
# - CONTROLLER_ADDRESS: fschat控制器地址
# - HF_ENDPOINT: huggingface地址
# - HF_HOME: huggingface家目录
# - HTTP_PROXY: HTTP代理
# - HTTPS_PROXY: HTTPS代理
# - NO_PROXY: 不代理的地址

# shellcheck disable=SC2153
if [ "$HTTP_PROXY" != "" ]; then
    export http_proxy=$HTTP_PROXY
fi

if [ "$HTTPS_PROXY" != "" ]; then
    export https_proxy=$HTTPS_PROXY
fi

if [ "$NO_PROXY" != "" ]; then
    export no_proxy=$NO_PROXY
fi

MODEL_WORKER=fastchat.serve.model_worker
OS_TYPE=$(uname)

# MODEL_WORKER 
if [ "$USE_VLLM" == "true" ]; then
    MODEL_WORKER="fastchat.serve.vllm_worker"
fi

# 量化配置
if [ "$QUANTIZATION" == "8bit" ]; then
    QUANTIZATION="--load-8bit"
else
    QUANTIZATION=""
fi

# NUM_GPUS
if [ "$NUM_GPUS" -gt 0 ]; then
    NUM_GPUS="--num-gpus $NUM_GPUS"
else
    NUM_GPUS=""
fi

# CPU推理CPU，mps
if [ "$INFERRED_TYPE" == "cpu" ] && [ "$OS_TYPE" == "Darwin" ]; then
    DEVICE_OPTION="--device mps"
elif [ "$INFERRED_TYPE" == "cpu" ]; then
    DEVICE_OPTION="--device cpu"
else
    DEVICE_OPTION=""
fi

# MAX_GPU_MEMORY
if [ "$MAX_GPU_MEMORY" -gt 0 ]; then
    MAX_GPU_MEMORY="--max-gpu-memory ${MAX_GPU_MEMORY}GiB"
else
    MAX_GPU_MEMORY=""
fi

# 当前Pod的IP
if [ ! -n "$MY_POD_IP" ]; then
	MY_POD_IP=$(hostname -I | awk '{print $1}')
fi

# 模型路径
if [ -n "$MODEL_PATH" ]; then
	MODEL_PATH="--model-path $MODEL_PATH"
fi

python3 -m $MODEL_WORKER --host 0.0.0.0 --port $HTTP_PORT \
    --controller-address $CONTROLLER_ADDRESS \
    --worker-address http://$MY_POD_IP:$HTTP_PORT \
    --model-name $MODEL_NAME \
    $MODEL_PATH $QUANTIZATION $NUM_GPUS $MAX_GPU_MEMORY $DEVICE_OPTION
`
	shellTrain = `#!/bin/bash
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
# - MASTER_PORT: 主节点端口
# - USE_LORA: 是否使用LORA
# - OUTPUT_DIR: 输出路径
# - NUM_TRAIN_EPOCHS: 训练轮数
# - PER_DEVICE_TRAIN_BATCH_SIZE: 训练批次大小
# - PER_DEVICE_EVAL_BATCH_SIZE: 评估批次大小
# - GRADIENT_ACCUMULATION_STEPS: 累积步数
# - LEARNING_RATE: 学习率
# - MODEL_MAX_LENGTH: 模型最大长度
# - BASE_MODEL_PATH: 基础模型路径
# - BASE_MODEL_NAME: 基础模型名称
# - SCENARIO: 应用场景
# - TRAIN_FILE: 训练文件URL或路径
# - EVAL_FILE: 验证文件URL或路径
export CUDA_DEVICE_MAX_CONNECTIONS=1

JOB_STATUS="success"
JOB_MESSAGE=""

if [ "$MASTER_PORT" == "" ]; then
    MASTER_PORT=29500
fi

# shellcheck disable=SC2153
if [ "$HTTP_PROXY" != "" ]; then
    export http_proxy=$HTTP_PROXY
fi

if [ "$HTTPS_PROXY" != "" ]; then
    export https_proxy=$HTTPS_PROXY
fi

if [ "$NO_PROXY" != "" ]; then
    export no_proxy=$NO_PROXY
fi

function callback() {
  # 根据退出状态判断执行是否异常
  if [ $JOB_STATUS -eq 0 ]; then
      # 没有发生异常，正常输出内容
      echo "执行成功!"
      echo "${API_URL}"
      # 调用API并传递输出内容
      curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "{\"status\": \"success\"}"
  else
      # 发生异常
      echo "执行失败!"
      # 调用API并传递错误信息
      sleep 40
      JOB_MESSAGE=$(jq -n --arg content "$JOB_MESSAGE" '{"status": "failed", "message": $content}')
      curl -X PUT "${API_URL}" -H "Authorization: ${AUTH}" -H "X-Tenant-Id: ${TENANT_ID}" -H "Content-Type: application/json" -d "$JOB_MESSAGE"
  fi
}

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


LORA_MODULE_NAME=''
MODENAME=$(echo "$BASE_MODEL_NAME" | tr '[:upper:]' '[:lower:]')
case $MODENAME in
    *'llama2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='llama2'
        ;;
    *'baichuan2'*)
        LORA_MODULE_NAME='W_pack'
        MODENAME='baichuan2_13b'
        ;;
    *'glm3'*)
        LORA_MODULE_NAME='query_key_value,dense_h_to_4h,dense_4h_to_h,dense'
        MODENAME='glm3'
        ;;
    *'glm2'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm2'
        ;;
    *'glm'*)
        LORA_MODULE_NAME='gate_proj,down_proj,up_proj'
        MODENAME='glm'
        ;;
    *'qwen1.5'*)
        LORA_MODULE_NAME='q_proj,k_proj,v_proj,o_proj,up_proj,gate_proj,down_proj'
        MODENAME='qwen1.5'
        ;;
    *)
        MODENAME='auto'
        echo "未知模型名称"
        ;;
esac

if [ "$USE_LORA" = true ]; then
    TRAIN_TYPE="lora"
    ZERO_STAGE=2
    DS_FILE=faq/ds_zero2_no_offload.json
else
    TRAIN_TYPE="all"
    ZERO_STAGE=3
    DS_FILE=faq/ds_zero3_offload.json
fi

URL_REGEX="^(http|https)://"
mkdir -p /data/train-data/
TRAIN_LOCAL_FILE=/data/train-data/train-${JOB_ID}.jsonl
EVAL_LOCAL_FILE=/data/train-data/eval-${JOB_ID}.jsonl

function download_file {
    if [[ "$1" =~ $URL_REGEX ]]; then
        echo "The path is a URL. Starting download..."
        wget -O $2 "$1"
    else
        echo "File path is local. No download needed."
        cp "$1" "$2"
    fi
}

download_file "$TRAIN_FILE" "$TRAIN_LOCAL_FILE"

if [ "$EVAL_FILE" != "" ]; then
  download_file "$EVAL_FILE" "$EVAL_LOCAL_FILE"
fi

temp_file=$(mktemp)

if [ "$SCENARIO" == "general" ]; then
  GENERAL_DATA_PATH=/data/train-data/formatted_datasets/${JOB_ID}
  mkdir -p $GENERAL_DATA_PATH
  if [ -n "$EVAL_FILE" ]; then
    python3 jsonl_to_arrow_format.py \
        --train_path "$TRAIN_LOCAL_FILE" \
        --test_path "$EVAL_LOCAL_FILE" \
        --output_path "$GENERAL_DATA_PATH"
  else
    python3 jsonl_to_arrow_format.py \
        --train_path "$TRAIN_LOCAL_FILE" \
        --output_path "$GENERAL_DATA_PATH"
  fi


#  output=$(torchrun $DISTRIBUTED_ARGS {{.ScriptFile}} \
#  output=$(deepspeed --include localhost:$CUDA_VISIBLE_DEVICES {{.ScriptFile}} \
  deepspeed /app/llmops_deepspeed_main.py \
      --data_path $GENERAL_DATA_PATH \
      --data_output_path $OUTPUT_DIR/data_output \
      --data_split 9,1,0 \
      --model_name_or_path $BASE_MODEL_PATH \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --per_device_eval_batch_size $PER_DEVICE_EVAL_BATCH_SIZE \
      --max_seq_len $MODEL_MAX_LENGTH \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0. \
      --num_train_epochs $NUM_TRAIN_EPOCHS  \
      --train_type $TRAIN_TYPE \
      --lora_dim 2 \
      --lora_module_name $LORA_MODULE_NAME \
      --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
      --lr_scheduler_type cosine \
      --num_warmup_steps 0 \
      --seed 42 \
      --gradient_checkpointing \
      --zero_stage $ZERO_STAGE \
      --offload \
      --only_optimize_lora \
      --deepspeed \
      --output_dir $OUTPUT_DIR \
      --start_from_step -1 \
      --save_per_steps 100  > >(tee "$temp_file") 2>&1
elif [ "$SCENARIO" == "faq" ]; then
  formatted_datasets_path=/data/train-data/faq_formatted_datasets
  mkdir -p "$formatted_datasets_path"

  python3 /app/convert_new_format.py \
      --train_path $TRAIN_LOCAL_FILE \
      --test_path $EVAL_LOCAL_FILE \
      --output_path "$formatted_datasets_path"

#  output=$(deepspeed {{.ScriptFile}}  \
  deepspeed /app/faq/faq_train.py \
      --train_path "$formatted_datasets_path/train_dataset.jsonl" \
      --model_name_or_path $BASE_MODEL_PATH \
      --per_device_train_batch_size $PER_DEVICE_TRAIN_BATCH_SIZE \
      --max_len $MODEL_MAX_LENGTH \
      --max_src_len 128 \
      --learning_rate $LEARNING_RATE \
      --weight_decay 0.1 \
      --num_train_epochs 2 \
      --gradient_accumulation_steps $GRADIENT_ACCUMULATION_STEPS \
      --warmup_ratio 0.1 \
      --mode $MODENAME \
      --train_type $TRAIN_TYPE \
      --lora_module_name $LORA_MODULE_NAME \
      --lora_dim 4 \
      --lora_alpha 64 \
      --lora_dropout 0.1 \
      --seed 1234 \
      --ds_file $DS_FILE \
      --gradient_checkpointing \
      --show_loss_step 10 \
      --output_dir $OUTPUT_DIR  > >(tee "$temp_file") 2>&1

elif [ "$SCENARIO" == "rag" ]; then

  pass

else
  echo "Invalid scenario selection!"
fi

JOB_STATUS=$?
JOB_MESSAGE=$(<"$temp_file")
callback
`

	generateCmd = &cobra.Command{
		Use:               "generate command <args> [flags]",
		Short:             "生成命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 生成命令
可用的配置类型：
[table, init-data]

aigc-server generate -h
`,
	}

	genTableCmd = &cobra.Command{
		Use:               `table <args> [flags]`,
		Short:             "生成数据库表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
aigc-server generate table all
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
			}()

			if len(args) > 0 && args[0] == "all" {
				return generateTable()
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = log.NewLogfmtLogger(os.Stdout)
			return prepare(cmd.Context())
		},
	}
)

func generateTable() (err error) {
	_ = logger.Log("migrate", "table", "Chat", gormDB.AutoMigrate(types.Chat{}))
	_ = logger.Log("migrate", "table", "ChatAllowUser", gormDB.AutoMigrate(types.ChatAllowUser{}))
	_ = logger.Log("migrate", "table", "ChatConversation", gormDB.AutoMigrate(types.ChatConversation{}))
	_ = logger.Log("migrate", "table", "ChatSystemPrompt", gormDB.AutoMigrate(types.ChatSystemPrompt{}))
	_ = logger.Log("migrate", "table", "ChatPromptTypes", gormDB.AutoMigrate(types.ChatPromptTypes{}))
	_ = logger.Log("migrate", "table", "ChatChannels", gormDB.AutoMigrate(types.ChatChannels{}))
	_ = logger.Log("migrate", "table", "ChatPrompts", gormDB.AutoMigrate(types.ChatPrompts{}))
	_ = logger.Log("migrate", "table", "ChatChannelModels", gormDB.AutoMigrate(types.ChatChannelModels{}))
	_ = logger.Log("migrate", "table", "ChatMessages", gormDB.AutoMigrate(types.ChatMessages{}))
	_ = logger.Log("migrate", "table", "Dataset", gormDB.AutoMigrate(types.Dataset{}))
	_ = logger.Log("migrate", "table", "DatasetSample", gormDB.AutoMigrate(types.DatasetSample{}))
	_ = logger.Log("migrate", "table", "Assistants", gormDB.AutoMigrate(types.Assistants{}))
	_ = logger.Log("migrate", "table", "Tools", gormDB.AutoMigrate(types.Tools{}))
	_ = logger.Log("migrate", "table", "AssistantToolAssociations", gormDB.AutoMigrate(types.AssistantToolAssociations{}))
	_ = logger.Log("migrate", "table", "Files", gormDB.AutoMigrate(types.Files{}))
	_ = logger.Log("migrate", "table", "SysAudit", gormDB.AutoMigrate(types.SysAudit{}))
	_ = logger.Log("migrate", "table", "FineTuningTrainJob", gormDB.AutoMigrate(types.FineTuningTrainJob{}))
	_ = logger.Log("migrate", "table", "FineTuningTemplate", gormDB.AutoMigrate(types.FineTuningTemplate{}))
	_ = logger.Log("migrate", "table", "Tenants", gormDB.AutoMigrate(types.Tenants{}))
	_ = logger.Log("migrate", "table", "Models", gormDB.AutoMigrate(types.Models{}))
	_ = logger.Log("migrate", "table", "SysDict", gormDB.AutoMigrate(types.SysDict{}))
	_ = logger.Log("migrate", "table", "ModelDeploy", gormDB.AutoMigrate(types.ModelDeploy{}))
	_ = logger.Log("migrate", "table", "LLMEvalResults", gormDB.AutoMigrate(types.LLMEvalResults{}))
	_ = logger.Log("migrate", "table", "DatasetDocument", gormDB.AutoMigrate(types.DatasetDocument{}))
	_ = logger.Log("migrate", "table", "DatasetDocumentSegment", gormDB.AutoMigrate(types.DatasetDocumentSegment{}))
	_ = logger.Log("migrate", "table", "DatasetAnnotationTask", gormDB.AutoMigrate(types.DatasetAnnotationTask{}))
	_ = logger.Log("migrate", "table", "DatasetAnnotationTaskSegment", gormDB.AutoMigrate(types.DatasetAnnotationTaskSegment{}))
	_ = logger.Log("migrate", "table", "ModelEvaluate", gormDB.AutoMigrate(types.ModelEvaluate{}))
	//err = initData()
	//if err != nil {
	//	return err
	//}
	return
}

// 初始化数据
func initData() (err error) {
	tenant := types.Tenants{
		Name:           "系统租户",
		PublicTenantID: uuid.New().String(),
		ContactEmail:   serverAdminUser,
	}
	_ = logger.Log("init", "data", "SysDict", gormDB.Create(&tenant).Error)
	password, err := bcrypt.GenerateFromPassword([]byte(serverAdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_ = logger.Log("init", "data", "account", gormDB.Save(&types.Accounts{
		Email:        serverAdminUser,
		Nickname:     "系统管理员",
		Language:     "zh",
		IsLdap:       false,
		PasswordHash: string(password),
		Status:       true,
		Tenants:      []types.Tenants{tenant},
	}).Error)
	if serverChannelKey == "" {
		serverChannelKey = "sk-" + string(util.Krand(48, util.KC_RAND_KIND_ALL))
	}
	_ = logger.Log("init", "data", "ChatChannels", gormDB.Create(&types.ChatChannels{
		Name:       "default",
		Alias:      "默认渠道",
		Remark:     "默认渠道",
		Quota:      10000,
		Models:     "default",
		OnlyOpenAI: false,
		ApiKey:     serverChannelKey,
		Email:      serverAdminUser,
		TenantId:   tenant.ID,
	}).Error)

	var templateModels []string
	templateModels = append(templateModels, "qwen1.5-0.5b", "qwen1.5-1.8b", "qwen1.5-1.8b-chat", "qwen1.5-4b", "qwen1.5-4b-chat",
		"qwen1.5-7b", "qwen1.5-7b-chat", "qwen1.5-14b", "qwen1.5-14b-chat", "qwen1.5-32b", "qwen1.5-32b-chat", "qwen1.5-72b", "qwen1.5-72b-chat",
		"chatglm3-6b", "chatglm3-6b-32k",
		"llama-2-13b-chat", "llama-2-13b", "llama-2-7b-chat", "llama-2-7b-chat",
		"baichuan2-7b-base", "baichuan2-7b-chat", "baichuan2-13b-base", "baichuan2-13b-chat",
	)
	replacer := strings.NewReplacer(
		"::", "-", // 这个可能不需要，因为前一个已经将单个冒号替换了
		":", "-",
	)
	for _, model := range templateModels {
		_ = logger.Log("init", "data", "models-template-inference", gormDB.Create(&types.FineTuningTemplate{
			Name:          fmt.Sprintf("%s", model),
			BaseModel:     model,
			MaxTokens:     32768,
			Content:       shellStart,
			TrainImage:    "dudulu/llmops:latest",
			BaseModelPath: "/data/base-model/" + strings.ToLower(replacer.Replace(model)),
			TemplateType:  "inference",
			ScriptFile:    "/app/start.sh",
			Enabled:       true,
			OutputDir:     "/data/ft-model",
		}).Error)
		_ = logger.Log("init", "data", "models-template-train", gormDB.Create(&types.FineTuningTemplate{
			Name:          fmt.Sprintf("%s-train", model),
			BaseModel:     model,
			MaxTokens:     32768,
			Content:       shellTrain,
			TrainImage:    "dudulu/llmops:latest",
			BaseModelPath: "/data/base-model/" + strings.ToLower(replacer.Replace(model)),
			TemplateType:  "train",
			ScriptFile:    "/app/train.sh",
			Enabled:       true,
			OutputDir:     "/data/ft-model",
		}).Error)

	}

	_ = logger.Log("init", "data", "sys_dict", gormDB.Exec(initSysDictSql).Error)
	//_ = logger.Log("init", "data", "finetuning_template", gormDB.Exec(ftTemplateSql).Error)
	_ = logger.Log("init", "data", "models", gormDB.Exec(modelSql).Error)
	return err
}

var (
	modelSql = `INSERT INTO models (created_at, updated_at, deleted_at, provider_name, model_type, model_name, max_tokens, is_private, is_fine_tuning, enabled, remark, parameters, last_operator, base_model_name, replicas, label, k8s_cluster, inferred_type, gpu, cpu, memory)
VALUES
	('2024-02-04 13:02:48.112', '2024-03-19 14:08:35.667', NULL, 'OpenAI', 'text-generation', 'gpt-3.5-turbo', 4096, 0, 0, 1, 'OpenAI GPT-3.5-turbo', 20.00, '', NULL, 1, NULL, NULL, NULL, 0, 0, 1),
	('2024-03-18 17:34:59.542', '2024-03-19 14:51:35.630', NULL, 'LocalAI', 'text-generation', 'qwen1.5-0.5b', 32768, 1, 0, 0, '', 0.50, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:03:11.073', '2024-03-19 14:41:08.194', NULL, 'LocalAI', 'text-generation', 'qwen1.5-1.8b', 32768, 0, 0, 0, '', 1.80, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:03:34.619', '2024-03-19 14:41:41.709', NULL, 'LocalAI', 'text-generation', 'qwen1.5-1.8b-chat', 32768, 0, 0, 0, '', 1.80, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:03:51.375', '2024-03-19 14:41:36.354', NULL, 'LocalAI', 'text-generation', 'qwen1.5-4b', 32768, 0, 0, 0, '', 3.98, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:04:11.425', '2024-03-19 14:41:11.423', NULL, 'LocalAI', 'text-generation', 'qwen1.5-4b-chat', 32768, 0, 0, 0, '', 3.98, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:04:29.257', '2024-03-19 14:41:18.790', NULL, 'LocalAI', 'text-generation', 'qwen1.5-7b', 32768, 0, 0, 0, '', 7.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:04:45.241', '2024-03-19 14:41:24.050', NULL, 'LocalAI', 'text-generation', 'qwen1.5-7b-chat', 32768, 0, 0, 0, '', 7.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:05:04.519', '2024-03-19 14:41:27.394', NULL, 'LocalAI', 'text-generation', 'qwen1.5-14b', 32768, 0, 0, 0, '', 14.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:05:27.624', '2024-03-19 14:41:47.741', NULL, 'LocalAI', 'text-generation', 'qwen1.5-14b-chat', 32768, 0, 0, 0, '', 14.20, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:06:26.666', '2024-03-19 14:41:33.633', NULL, 'LocalAI', 'text-generation', 'qwen1.5-72b', 32768, 0, 0, 0, '', 72.30, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:06:43.121', '2024-03-19 14:41:30.391', NULL, 'LocalAI', 'text-generation', 'qwen1.5-72b-chat', 32768, 0, 0, 0, '', 72.30, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'qwen1.5-32b', 32768, 0, 0, 0, '', 32.40, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'qwen1.5-32b-chat', 32768, 0, 0, 0, '', 32.40, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'chatglm3-6b-32k', 32768, 0, 0, 0, '', 6.40, 'admin', '', 1, '', '', '', 0, 0, 1),
	'2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'chatglm3-6b', 8192, 0, 0, 0, '', 6.40, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'baichuan2-7b-base', 4096, 0, 0, 0, '', 6.8, 'admin', '', 1, '', '', '', 0, 0, 1),
	('2024-03-19 14:08:27.352', '2024-03-19 14:41:01.068', NULL, 'LocalAI', 'text-generation', 'baichuan2-7b-chat', 4096, 0, 0, 0, '', 6.8, 'admin', '', 1, '', '', '', 0, 0, 1);
`

	initSysDictSql = `INSERT INTO sys_dict (id, created_at, updated_at, deleted_at, parent_id, code, dict_value, dict_label, dict_type, sort, remark)
VALUES
	(2, '2023-11-22 16:19:52.000', '2024-01-29 10:32:18.000', NULL, 0, 'speak_gender', 'gender', '性别', 'int', 1, '性别'),
	(3, '2023-11-22 16:23:19.000', '2024-01-29 10:32:18.000', NULL, 2, 'speak_gender', '1', '男', 'int', 1, '性别:男'),
	(4, '2023-11-22 16:24:27.000', '2024-01-29 10:32:18.000', NULL, 2, 'speak_gender', '2', '女', 'int', 0, '性别:女'),
	(5, '2023-11-23 10:17:31.000', '2023-11-30 10:42:43.000', NULL, 0, 'speak_age_group', 'speak_age_group', '年龄段', 'int', 0, ''),
	(6, '2023-11-23 10:18:31.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '1', '少年', 'int', 5, ''),
	(7, '2023-11-23 10:18:46.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '2', '青年', 'int', 4, ''),
	(8, '2023-11-23 10:18:56.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '3', '中年', 'int', 4, ''),
	(9, '2023-11-23 10:19:21.000', '2023-11-23 10:20:51.000', NULL, 5, 'speak_age_group', '4', '老年', 'int', 2, ''),
	(10, '2023-11-23 10:25:07.000', '2023-11-30 10:42:43.000', NULL, 0, 'speak_style', 'speak_style', '风格', 'int', 0, ''),
	(11, '2023-11-23 10:25:53.000', '2023-11-23 10:25:53.000', NULL, 10, 'speak_style', '1', '温柔', 'int', 5, ''),
	(12, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '2', '阳光', 'int', 4, ''),
	(13, '2023-11-23 10:28:24.000', '2023-12-08 14:04:05.000', NULL, 0, 'speak_area', 'speak_area', '适应范围', 'int', 0, ''),
	(14, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '1', '客服', 'int', 5, ''),
	(15, '2023-11-23 10:29:14.000', '2023-11-23 10:29:14.000', NULL, 13, 'speak_area', '2', '小说', 'int', 4, ''),
	(16, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'speak_lang', 'speak_lang', '语言', 'string', 0, ''),
	(17, '2023-11-23 10:33:28.000', '2023-11-23 10:33:28.000', NULL, 16, 'speak_lang', 'zh-CN', '中文（普通话，简体）', 'string', 100, ''),
	(18, '2023-11-23 10:34:08.000', '2023-11-23 10:34:08.000', NULL, 16, 'speak_lang', 'zh-HK', '中文（粤语，繁体）', 'string', 99, ''),
	(19, '2023-11-23 10:34:30.000', '2023-11-23 10:34:30.000', NULL, 16, 'speak_lang', 'en-US', '英语（美国）', 'string', 98, ''),
	(20, '2023-11-23 10:35:07.000', '2023-11-23 10:35:07.000', NULL, 16, 'speak_lang', 'en-GB', '英语（英国）', 'string', 97, ''),
	(21, '2023-11-23 10:44:23.000', '2023-11-23 10:44:23.000', NULL, 0, 'speak_provider', 'speak_provider', '供应商', 'string', 0, ''),
	(22, '2023-11-23 10:44:50.000', '2023-11-23 10:44:50.000', NULL, 21, 'speak_provider', 'azure', '微软', 'string', 0, ''),
	(23, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '3', '自然流畅', 'int', 0, ''),
	(24, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '4', '亲切温和', 'int', 0, ''),
	(25, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '5', '温柔甜美', 'int', 0, ''),
	(26, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '6', '成熟知性', 'int', 0, ''),
	(27, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '7', '大气浑厚', 'int', 0, ''),
	(28, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '8', '稳重磁性', 'int', 0, ''),
	(29, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '9', '年轻时尚', 'int', 0, ''),
	(30, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '10', '轻声耳语', 'int', 0, ''),
	(31, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '11', '可爱甜美', 'int', 0, ''),
	(32, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '12', '呆萌可爱', 'int', 0, ''),
	(33, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '13', '激情力度', 'int', 0, ''),
	(34, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '14', '饱满活泼', 'int', 0, ''),
	(35, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '15', '诙谐幽默', 'int', 0, ''),
	(36, '2023-11-23 10:26:04.000', '2023-11-23 10:26:04.000', NULL, 10, 'speak_style', '16', '淳朴方言', 'int', 0, ''),
	(37, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '3', '新闻', 'int', 0, ''),
	(38, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '4', '纪录片', 'int', 0, ''),
	(39, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '5', '解说', 'int', 0, ''),
	(40, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '6', '教育', 'int', 0, ''),
	(41, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '7', '广告', 'int', 0, ''),
	(42, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '8', '直播', 'int', 0, ''),
	(43, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '9', '助理', 'int', 0, ''),
	(44, '2023-11-23 10:28:57.000', '2023-11-23 10:28:57.000', NULL, 13, 'speak_area', '10', '特色', 'int', 0, ''),
	(45, '2023-11-23 10:34:08.000', '2023-11-24 17:54:17.000', NULL, 16, 'speak_lang', 'zh-CN-henan', '中文（中原官话河南，简体）', 'string', 99, ''),
	(46, '2023-11-23 10:34:08.000', '2023-11-24 17:54:19.000', NULL, 16, 'speak_lang', 'zh-CN-liaoning', '中文（东北官话，简体）', 'string', 99, ''),
	(47, '2023-11-23 10:34:08.000', '2023-11-24 17:54:20.000', NULL, 16, 'speak_lang', 'zh-TW', '中文（台湾普通话，繁体）', 'string', 99, ''),
	(48, '2023-11-23 10:34:08.000', '2023-11-24 17:54:22.000', NULL, 16, 'speak_lang', 'zh-CN-GUANGXI', '中文（广西口音普通话，简体）', 'string', 99, ''),
	(49, '2023-11-23 10:35:07.000', '2023-11-23 10:35:07.000', NULL, 16, 'speak_lang', 'ko-KR', '韩语(韩国)', 'string', 97, ''),
	(50, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'ja-JP', '日语（日本）', 'string', 97, ''),
	(51, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'fil-PH', '菲律宾语（菲律宾）', 'string', 97, ''),
	(52, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'es-MX', '西班牙语(墨西哥)', 'string', 97, ''),
	(53, '2023-11-23 10:35:07.000', '2023-11-24 19:45:54.000', NULL, 16, 'speak_lang', 'ru-RU', '俄语（俄罗斯）', 'string', 97, ''),
	(54, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'audio_tagged_lang', 'audio_tagged_lang', '标记语言', 'string', 0, ''),
	(55, '2023-11-23 10:32:39.000', '2023-12-06 14:55:30.000', NULL, 54, 'audio_tagged_lang', 'zh', '中文', 'string', 1001, ''),
	(56, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'en', '英文', 'string', 99, ''),
	(57, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'es', '西班牙语', 'string', 98, ''),
	(58, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'de', '德语', 'string', 97, ''),
	(59, '2023-11-23 10:32:39.000', '2023-11-28 17:32:38.000', NULL, 54, 'audio_tagged_lang', 'tl', '他加禄语', 'string', 96, ''),
	(60, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 54, 'audio_tagged_lang', 'fil', '菲律宾语', 'string', 95, ''),
	(61, '2023-11-23 10:32:39.000', '2023-11-23 10:32:39.000', NULL, 0, 'sys_dict_type', 'sys_dict_type', '字典类型', 'string', 0, ''),
	(62, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'string', '字符串类型', 'string', 100, ''),
	(63, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'int', '数字类型', 'int', 99, ''),
	(64, '2023-11-23 10:32:39.000', '2023-12-08 13:55:43.000', NULL, 61, 'sys_dict_type', 'bool', '布尔类型', 'bool', 98, ''),
	(87, '2023-12-07 17:18:32.000', '2023-12-07 17:18:32.000', NULL, 0, 'language', 'language', '国际化', 'string', 0, ''),
	(88, '2023-12-14 15:53:15.000', '2023-12-14 16:37:19.000', NULL, 0, 'model_eval_dataset_type', 'model_eval_dataset_type', '模型评估数据集类型', 'string', 110, ''),
	(89, '2023-12-14 15:54:28.000', '2023-12-14 16:37:19.000', NULL, 88, 'model_eval_dataset_type', 'train', '训练集', 'string', 99, ''),
	(90, '2023-12-14 15:54:57.000', '2023-12-14 16:37:19.000', NULL, 88, 'model_eval_dataset_type', 'custom', '自定义', 'string', 98, ''),
	(91, '2023-12-14 15:55:56.000', '2023-12-14 15:59:39.000', NULL, 0, 'model_eval_metric', 'model_eval_metric', '模型评估指标', 'string', 99, ''),
	(92, '2023-12-14 15:58:48.000', '2023-12-14 15:58:48.000', NULL, 0, 'model_eval_status', 'model_eval_status', '模型评估状态', 'string', 99, ''),
	(93, '2023-12-14 16:00:31.000', '2023-12-14 16:00:31.000', NULL, 92, 'model_eval_status', 'pending', '等待评估', 'string', 99, ''),
	(94, '2023-12-14 16:00:44.000', '2023-12-14 16:00:44.000', NULL, 92, 'model_eval_status', 'running', '正在评估', 'string', 98, ''),
	(95, '2023-12-14 16:00:56.000', '2023-12-14 16:00:56.000', NULL, 92, 'model_eval_status', 'success', '评估成功', 'string', 97, ''),
	(96, '2023-12-14 16:01:09.000', '2023-12-14 16:01:09.000', NULL, 92, 'model_eval_status', 'failed', '评估失败', 'string', 97, ''),
	(97, '2023-12-14 16:01:23.000', '2023-12-14 16:01:23.000', NULL, 92, 'model_eval_status', 'cancel', '评估取消', 'string', 96, ''),
	(98, '2023-12-14 16:11:47.000', '2023-12-14 16:14:05.000', NULL, 91, 'model_eval_metric', 'equal', '完全匹配', 'string', 98, ''),
	(99, '2023-12-14 17:20:04.000', '2023-12-14 17:20:04.000', NULL, 0, 'model_deploy_status', 'model_deploy_status', '模型部署状态', 'string', 0, ''),
	(100, '2023-12-14 17:20:34.000', '2023-12-15 11:39:30.000', NULL, 99, 'model_deploy_status', 'pending', '部署中', 'string', 0, ''),
	(101, '2023-12-14 17:20:45.000', '2023-12-15 11:39:35.000', NULL, 99, 'model_deploy_status', 'running', '运行中', 'string', 0, ''),
	(102, '2023-12-14 17:21:06.000', '2023-12-14 17:22:34.000', NULL, 99, 'model_deploy_status', 'success', '完成', 'string', 0, ''),
	(103, '2023-12-14 17:24:54.000', '2023-12-14 17:24:54.000', NULL, 99, 'model_deploy_status', 'failed', '失败', 'string', 0, ''),
	(104, '2023-12-14 20:44:04.000', '2024-01-30 11:19:36.000', NULL, 0, 'model_provider_name', 'model_provider_name', '模型供应商', 'string', 0, ''),
	(105, '2023-12-14 20:45:32.000', '2024-01-30 11:19:36.000', NULL, 104, 'model_provider_name', 'LocalAI', 'LocalAI', 'string', 0, ''),
	(106, '2023-12-14 20:45:44.000', '2024-01-30 11:19:36.000', NULL, 104, 'model_provider_name', 'OpenAI', 'OpenAI', 'string', 0, ''),
	(107, '2023-12-15 15:09:42.000', '2023-12-15 15:09:42.000', NULL, 0, 'digitalhuman_synthesis_status', 'digitalhuman_synthesis_status', '数字人合成状态', 'string', 0, ''),
	(108, '2023-12-15 15:10:48.000', '2023-12-15 15:10:48.000', NULL, 107, 'digitalhuman_synthesis_status', 'running', '合成中', 'string', 0, ''),
	(109, '2023-12-15 15:11:12.000', '2023-12-15 15:11:12.000', NULL, 107, 'digitalhuman_synthesis_status', 'success', '已完成', 'string', 0, ''),
	(110, '2023-12-15 15:11:30.000', '2023-12-15 15:11:30.000', NULL, 107, 'digitalhuman_synthesis_status', 'failed', '失败', 'string', 0, ''),
	(111, '2023-12-15 15:11:48.000', '2023-12-15 15:11:48.000', NULL, 107, 'digitalhuman_synthesis_status', 'waiting', '等待中', 'string', 0, ''),
	(112, '2023-12-15 15:12:02.000', '2023-12-15 15:12:02.000', NULL, 107, 'digitalhuman_synthesis_status', 'cancel', '已取消', 'string', 0, ''),
	(113, '2023-12-20 19:07:56.000', '2023-12-20 19:07:56.000', NULL, 0, 'digitalhuman_posture', 'digitalhuman_posture', '数字人姿势', 'int', 0, ''),
	(114, '2023-12-20 19:09:10.000', '2023-12-21 10:11:35.000', NULL, 113, 'digitalhuman_posture', '1', '全身', 'int', 0, ''),
	(115, '2023-12-20 19:09:39.000', '2023-12-21 10:11:44.000', NULL, 113, 'digitalhuman_posture', '2', '半身', 'int', 0, ''),
	(116, '2023-12-20 19:10:22.000', '2023-12-21 10:11:53.000', NULL, 113, 'digitalhuman_posture', '3', '大半身', 'int', 0, ''),
	(117, '2023-12-20 19:10:34.000', '2023-12-21 10:11:58.000', NULL, 113, 'digitalhuman_posture', '4', '坐姿', 'int', 0, ''),
	(118, '2023-12-20 19:16:05.000', '2023-12-20 19:16:05.000', NULL, 0, 'digitalhuman_resolution', 'digitalhuman_resolution', '数字人分辨率', 'int', 0, ''),
	(119, '2023-12-20 19:20:03.000', '2023-12-20 19:20:03.000', NULL, 118, 'digitalhuman_resolution', '1', '480P', 'int', 0, ''),
	(120, '2023-12-20 19:20:22.000', '2023-12-20 19:20:22.000', NULL, 118, 'digitalhuman_resolution', '2', '720P', 'int', 0, ''),
	(121, '2023-12-20 19:20:43.000', '2023-12-20 19:20:43.000', NULL, 118, 'digitalhuman_resolution', '3', '1080P', 'int', 0, ''),
	(122, '2023-12-20 19:20:51.000', '2023-12-20 19:20:51.000', NULL, 118, 'digitalhuman_resolution', '4', '2K', 'int', 0, ''),
	(123, '2023-12-20 19:21:13.000', '2023-12-20 19:21:13.000', NULL, 118, 'digitalhuman_resolution', '5', '4K', 'int', 0, ''),
	(124, '2023-12-20 19:21:31.000', '2023-12-20 19:21:31.000', NULL, 118, 'digitalhuman_resolution', '6', '8K', 'int', 0, ''),
	(125, '2023-12-22 11:13:26.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 0, 'model_type', 'model_type', '模型类型', 'string', 0, ''),
	(126, '2023-12-22 11:17:02.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 125, 'model_type', 'train', '微调训练', 'string', 0, ''),
	(127, '2023-12-22 11:17:42.000', '2023-12-22 11:22:53.000', '2023-12-22 11:23:07.000', 125, 'model_type', 'inference', '模型推理', 'string', 0, ''),
	(128, '2023-12-22 11:24:15.000', '2023-12-22 11:24:15.000', NULL, 0, 'template_type', 'template_type', '模板类型', 'string', 0, ''),
	(129, '2023-12-22 11:25:41.000', '2023-12-22 11:30:11.000', NULL, 128, 'template_type', 'train', '微调训练', 'string', 0, ''),
	(130, '2023-12-22 11:26:50.000', '2023-12-22 11:29:04.000', NULL, 128, 'template_type', 'inference', '模型推理', 'string', 0, ''),
	(131, '2024-01-08 16:44:16.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 0, 'model_quantify', 'model_quantify', '模型量化', 'string', 0, '模特部署量化'),
	(132, '2024-01-08 16:45:30.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 131, 'model_quantify', 'bf16', '半精度', 'int', 0, ''),
	(133, '2024-01-08 16:47:23.000', '2024-01-09 16:57:15.000', '2024-01-09 16:57:29.000', 131, 'model_quantify', '8bit', '1/4精度', 'int', 1, '四分之一精度'),
	(134, '2024-01-09 10:50:12.000', '2024-01-09 10:50:12.000', NULL, 0, 'model_deploy_label', 'model_deploy_label', '模型部署标签', 'string', 0, ''),
	(135, '2024-01-09 10:51:24.000', '2024-03-19 17:17:07.017', NULL, 134, 'model_deploy_label', 'cpu-aigc-model', 'cpu-aigc-model', 'string', 0, ''),
	(136, '2024-01-09 10:52:19.000', '2024-01-09 10:52:19.000', NULL, 0, 'model_deploy_quantization', 'model_deploy_quantization', '模型部署量化', 'string', 0, '模型部署量化'),
	(137, '2024-01-09 10:52:40.000', '2024-01-09 10:52:40.000', NULL, 136, 'model_deploy_quantization', 'float16', 'float16', 'string', 0, ''),
	(138, '2024-01-09 10:52:46.000', '2024-01-09 10:52:46.000', NULL, 136, 'model_deploy_quantization', '8bit', '8bit', 'string', 0, ''),
	(156, '2024-01-23 10:10:54.000', '2024-01-23 10:10:54.000', NULL, 0, 'assistant_tool_type', 'assistant_tool_type', 'AI助手工具类型', 'string', 0, 'AI助手工具类型'),
	(157, '2024-01-23 10:12:23.000', '2024-01-25 15:06:41.000', NULL, 156, 'assistant_tool_type', 'function', 'API接口', 'string', 3, ''),
	(158, '2024-01-23 10:12:47.000', '2024-01-23 10:13:45.000', NULL, 156, 'assistant_tool_type', 'retrieval', '知识库', 'string', 2, ''),
	(159, '2024-01-23 10:13:01.000', '2024-01-23 10:13:25.000', NULL, 156, 'assistant_tool_type', 'code_interpreter', '代码执行', 'string', 1, ''),
	(160, '2024-01-24 11:07:37.000', '2024-01-24 11:07:37.000', NULL, 0, 'http_method', 'http_method', '请求方法', 'string', 111, 'http请求方法'),
	(161, '2024-01-24 11:08:10.000', '2024-01-24 11:09:02.000', NULL, 160, 'http_method', 'get', 'GET', 'string', 4, ''),
	(162, '2024-01-24 11:08:21.000', '2024-01-24 11:09:08.000', NULL, 160, 'http_method', 'post', 'POST', 'string', 3, ''),
	(163, '2024-01-24 11:08:36.000', '2024-01-24 11:09:12.000', NULL, 160, 'http_method', 'put', 'PUT', 'string', 2, ''),
	(164, '2024-01-24 11:09:54.000', '2024-01-24 11:09:54.000', NULL, 160, 'http_method', 'delete', 'DEL', 'string', 1, ''),
	(165, '2024-01-24 11:14:47.000', '2024-01-24 11:14:47.000', NULL, 0, 'programming_language', 'programming_language', '编程语言', 'string', 112, ''),
	(166, '2024-01-24 11:15:12.000', '2024-01-24 11:15:12.000', NULL, 165, 'programming_language', 'python', 'Python', 'string', 1, ''),
	(172, '2024-03-19 11:25:22.770', '2024-03-19 11:25:22.770', NULL, 0, 'textannotation_type', 'textannotation_type', '文本标注类型', 'string', 0, ''),
	(173, '2024-03-19 11:25:47.575', '2024-03-19 11:25:47.575', NULL, 172, 'textannotation_type', 'rag', '检索增强生成', 'string', 0, ''),
	(174, '2024-03-19 11:26:00.272', '2024-03-19 11:26:00.272', NULL, 172, 'textannotation_type', 'faq', '知识问答', 'string', 0, ''),
	(175, '2024-03-19 11:26:12.417', '2024-03-19 11:26:12.417', NULL, 172, 'textannotation_type', 'general', '通用', 'string', 0, ''),
	(176, '2024-03-19 14:01:56.036', '2024-03-19 14:01:56.036', NULL, 0, 'model_type', 'model_type', '模型类型', 'string', 0, '模型类型：文本模型，语音模型，数字人模型等'),
	(177, '2024-03-19 14:02:19.712', '2024-03-19 14:02:19.712', NULL, 176, 'model_type', 'text-generation', '文本', 'string', 0, ''),
	(178, '2024-03-19 14:02:28.164', '2024-03-19 14:02:28.164', NULL, 176, 'model_type', 'embedding', 'embedding', 'string', 0, ''),
	(179, '2024-03-19 18:04:29.830', '2024-03-19 18:04:29.830', NULL, 0, 'model_evaluate_target_type', 'model_evaluate_target_type', '模型评测指标', 'string', 0, ''),
	(180, '2024-03-19 18:05:04.555', '2024-03-19 18:05:04.555', NULL, 179, 'model_evaluate_target_type', 'Acc', 'ACC', 'string', 0, ''),
	(181, '2024-03-19 18:05:14.515', '2024-03-19 18:05:14.515', NULL, 179, 'model_evaluate_target_type', 'F1', 'F1', 'string', 0, ''),
	(182, '2024-03-19 18:05:22.487', '2024-03-19 18:05:22.487', NULL, 179, 'model_evaluate_target_type', 'BLEU', 'BLEU', 'string', 0, ''),
	(183, '2024-03-19 18:05:30.619', '2024-03-19 18:05:30.619', NULL, 179, 'model_evaluate_target_type', 'Rouge', 'Rouge', 'string', 0, ''),
	(184, '2024-03-19 18:05:38.596', '2024-03-19 18:05:38.596', NULL, 179, 'model_evaluate_target_type', 'five', '五维图', 'string', 0, '');`
)
