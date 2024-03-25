# LLMOps 说明文档

## 启动说明

- 使用 Docker 启动

1. **构建 Docker 镜像**:
   在项目的根目录（包含 Dockerfile 的地方）运行以下命令来构建 Docker 镜像

```bash
docker build -t llmops-deepspeed-backend .
```

- chat_main.sh 中进行场景选择: `general`即一般场景 或 `faq`即 FAQ 场景，详见下文`场景选择和自定义`部分

2. **运行 Docker 容器**:
   使用以下命令运行 Docker 容器：

```bash
docker run --gpus all --shm-size 64gb -p 6006:6006 -p 5000:5000 \
-v /Your/Path/mnt:/app/mnt/ \
llmops-deepspeed-backend
```

- 如遇 GPU 访问的问题，可尝试安装 NVIDIA Container Toolkit
- 在运行 Docker 容器时，使用 -v 参数来挂载目录，**将`/Your/Path/mnt`替换为你的目录**
- 挂载的目录`/Your/Path/mnt`结构应如下：

```bash
mnt
   - datasets
      - train.jsonl
      - test.jsonl
   - Llama-2-13b-chat-hf
   - output_path
```

其中，datasets 是数据集，Llama-2-13b-chat-hf 是预训练模型，output_path 是训练后的模型保存路径
数据集的格式为 jsonl，例如：

```json
{
  "messages": [
    { "role": "user", "content": "ping" },
    { "role": "assistant", "content": "Pong!" }
  ]
}
```

## 场景选择和自定义

在`chat_main.sh`脚本中，可以根据训练需求和场景配置来调整相关参数，从而灵活地进行模型训练。以下是详细的说明和操作步骤：

1. **场景选择**：

   - 使用变量`SCENARIO`来指定训练场景。`"general"`用于一般场景，而`"faq"`适用于常见问题解答（FAQ）场景。例如，`SCENARIO="general"`或`SCENARIO="faq"`。

2. **GPU 设备配置**：

   - 通过环境变量`CUDA_VISIBLE_DEVICES`设置要使用的 GPU 编号，以控制训练过程中的 GPU 分配。

3. **模型配置**：

   - 通过`--model_name_or_path`参数指定基座模型的路径。建议模型路径位于挂载目录`/Your/Path/mnt`。

4. **FAQ 场景的模型模式配置**：

   - 在 FAQ 场景（`SCENARIO="faq"`）中，可通过`--mode`参数选择特定的模型模式。可选模式包括：
     - `"auto"`：自动选择模式，适用于`Llama2-13b`模型。
     - `"glm"`：ChatGLM 模式。
     - `"glm2"`：ChatGLM2 模式。
     - `"glm3"`：ChatGLM3 模式。
     - `"baichuan2_13b"`：Baichuan2-13B 模式。
     - `"qwen1.5"`：Qwen1.5 模式。

5. **FAQ 场景的额外资源**：

   - 如需 FAQ 场景的详细配置和使用说明，请参考`faq`文件夹内的`readme`文件。该文件提供了全面的操作指南和配置详情。

6. **LoRA 微调设置**：
   - 通过设置`--train_type "lora"`开启 LoRA 微调（默认开启）。若需进行全量微调，可设置为`--train_type "all"`。此外，`--lora_module_name "gate_proj,down_proj,up_proj"`允许指定模型中要应用 LoRA 的特定模块，确保所选模块与模型结构相匹配。

## `eval`目录下的系列脚本

在`eval`目录下，提供了一套Shell脚本用于执行评估模块的功能，通过命令行参数进行操作：

1. **evaluate_model_from_five_dimensions.sh - 五维度模型性能评测模块**: 用于从五个关键维度出发，对模型的性能进行全面评估

2. **diagnosis_monitoring.sh - 诊断监控模块**: 监控模型的训练过程，及时识别潜在的风险，并针对这些风险提供相应的解决建议

3. **model_performance_evaluation.sh - 模型在指定数据集的评估模块**: 评估一个特定模型在指定数据集上的性能，可以指定模型、数据集、评估指标以及评估选项

4. **analyze_similar_questions_and_intents.sh - 数据 KnowHow 分析相似模块**: 分析数据集中的问题对，以发现文本相似度高但意图标记不同的问题对，并进一步分析这些意图是否在语义上相似

### 1.evaluate_model_from_five_dimensions.sh - 五维度模型性能评测模块

- evaluate_model_from_five_dimensions.py
- **模型输入**: 可以输入模型的 checkpoint 路径（如"Llama2-13b"）
- **评测指标**: 包含五个评估维度：推理能力、阅读理解能力、中文能力、指令遵从能力和创新能力

#### 输入格式

```bash
python evaluate_model_from_five_dimensions.py "path/to/model/checkpoint" '["inference_ability", "reading_comprehension", "chinese_language_skill", "command_compliance", "innovation_capacity"]' '{"additional_parameters": {"max_seq_len": 512}}'
```

##### 字段说明

- model_name_or_path: 待评估模型的名称或路径。
- evaluation_dimensions: 需要评估的能力维度列表。
- options: 包含评估配置的选项。
  - additional_parameters: 包含评估时的额外参数，如最大序列长度。

#### 输出格式

```json
{
  "evaluation_results": {
    "inference_ability": 0.95,
    "reading_comprehension": 0.9,
    "chinese_language_skill": 0.85,
    "command_compliance": 0.88,
    "innovation_capacity": 0.8
  },
  "diagnosis": {}
}
```

##### 字段说明

- evaluation_results: 包含评估结果的字典，每个能力维度对应一个分数。
  - inference_ability: 推理能力
  - reading_comprehension: 阅读理解能力
  - chinese_language_skill: 中文能力
  - command_compliance: 指令遵从能力
  - innovation_capacity: 创新能力

#### 方法说明

1. 推理能力 (inference_ability)

   - 数据集: AFQMC 格式，该数据集用于评估模型在理解自然语言推理方面的能力。数据集中的每条数据包含两个句子以及一个标签，标签"1"表示两个句子的含义相似，标签"0"表示两个句子的含义不同。
   - 计算方法: 准确率 (Accuracy)，通过比较模型预测的标签与真实标签来计算准确率
   - 原因: AFQMC 格式数据集有利于判断模型的推理能力，包括隐含关系等。且准确率是评估分类问题中常用的指标，能够直观地反映模型在句子相似性判断任务上的性能。

2. 阅读理解能力 (reading_comprehension)

   - 数据集: CMRC 2018 格式 (cmrc2018_dev)，这是中文机器阅读理解的基准数据集，要求模型根据给定的段落回答问题。
   - 计算方法: BLEU，通过计算模型生成的答案和真实答案之间的 n-gram 重叠来评估模型的阅读理解能力。
   - 原因:CMRC 2018 数据集是经典的阅读理解判断数据集。 BLEU 评分能够衡量生成文本的质量，适合用于评估阅读理解任务，其中模型需要根据上下文生成准确的答案。

3. 中文能力 (chinese_language_skill)

   - 数据集: TNEWS (tnews_public_dev)，来自今日头条的新闻分类数据集，包含多个新闻类别，如旅游、教育、金融等。
   - 计算方法: BLEU，通过评估模型总结出的关键词与真实关键词之间的重叠程度来计算。
   - 原因: 该任务要求模型能够理解新闻内容并提取关键信息，BLEU 评分能够评估模型在此任务上的表现。

4. 指令遵从能力 (command_compliance)

   - 数据集: BIG-bench (bigbench 的子任务 subtask001_quoref_question_generation)，要求模型根据给定的段落生成问题，问题需要准确地评估对文段中提及的人物、地点或事物的理解。
   - 计算方法: BLEU，用于评估生成问题的质量和准确性。
   - 原因: 此任务考察模型生成符合上下文且逻辑连贯的问题的能力，BLEU 评分可以有效地衡量这些问题的质量。

5. 创新能力 (innovation_capacity)

   - 数据集: M3KE 格式数据集，旨在测试模型在多级多学科知识掌握方面的能力，特别是在零样本或少样本的情况下。
   - 计算方法: 准确率 (Accuracy)，通过比较模型的预测结果与真实标签来评估模型的创新能力。
   - 原因: 准确率能够直观地展示模型在处理未见过的、多学科的新情境下的表现，反映模型的学习和推广能力。

### 2.diagnosis_monitoring.sh - 诊断监控模块

- diagnosis_monitoring.py
- 通过分析模型训练过程中生成的文本格式日志文件，实现对训练过程中可能出现的风险进行准确识别和分析
- **风险识别**: 识别过拟合、欠拟合和灾难性遗忘等风险。
- **解决建议**: 针对不同风险提供具体的解决建议。
- **数据来源**: 从模型训练的 txt 格式 log 日志加载

#### 输入格式

```json
python diagnosis_monitoring.py "/mnt/output_path/log.txt"
```

##### 字段说明

- log_path: 训练日志文件 txt 的路径。

#### 输出格式

```json
{
  "code": 0,
  "msg": "Monitoring Report",
  "data": {
    "current_risks": {
      "overfitting": "High",
      "underfitting": "Low",
      "catastrophic_forgetting": "High"
    },
    "recommendations": {
      "overfitting": "建议停止训练，降低学习率，增大dropout，增大数据量，降低训练周期数，并重新训练。",
      "underfitting": "-",
      "catastrophic_forgetting": "建议回退到上个版本并调整参数。"
    }
  }
}
```

##### 字段说明

- current_risks: 当前监控到的风险等级。
- recommendations: 针对不同风险的解决建议。

### 3.model_performance_evaluation.sh - 模型在指定数据集的评估模块

- model_performance_evaluation.py.py
- 评测模块 API 用于评估特定模型在指定数据集上的性能，可以指定模型、数据集、评估指标以及评估选项
- 给出两个经验公式：memory_usage: 评估过程中的显存使用情况和 evaluation_duration: 整个评估过程的持续时间

#### 输入格式

```json
python model_performance_evaluation.py --model_name_or_path "../../mnt/models/model_folfer" --dataset_path "../../mnt/datasets/datasets.jsonl" --evaluation_metrics "Rouge" --max_seq_len 512 --per_device_batch_size 10 --gpu_id 0 --output_path "../../mnt/output_path"
```

##### 字段说明

- model_name_or_path: 待评估模型的名称或路径。
- dataset_path: 用于评估的数据集路径。
- evaluation_metrics: 指定用于评估的指标列表，包括["Acc", "F1", "BLEU", "Rouge"]。
- options: 包含评估配置的选项，包括：
  - max_seq_len: 最大序列长度。
  - per_device_batch_size: 每个设备的批处理大小。
  - gpu_id: 用于评估的 GPU 编号。
  - output_path: 评估结果的输出路径。

#### 输出格式

```json
{
  "code": 0,
  "msg": "Success",
  "data": {
    "evaluation_results": {
      "Rouge": 0.88
    }
  }
}
```

##### 字段说明

- code: 响应状态码（0 表示成功）
- msg: 响应消息（成功时通常为"Success"）
- data: 包含评估结果的对象
  - evaluation_results: 一个字典，包含指定评估指标的评分

#### 两个经验公式

- 两个经验公式在 two_empirical_formulas.ipynb
  - memory_usage: 评估过程中的显存使用情况
  - valuation_duration: 整个评估过程的持续时间

```python
def memory_usage(model_parameters, batch_size, data_type_size=1, overhead_factor=1.5, num_gpus=1):
    """
    估算训练过程中的显存使用量（以GB为单位）。

    :param model_parameters: 模型的参数总量。
    :param batch_size: 训练过程中的batch大小。
    :param data_type_size: 数据类型的大小（以字节为单位），对于32位浮点数通常为4。
    :param overhead_factor: 额外开销因子。
    :param num_gpus: 使用的 GPU 数量。
    :return: 显存使用量（GB）。
    """
    return (model_parameters * data_type_size / 5e10) * (1 + overhead_factor) * batch_size / num_gpus

def evaluation_duration(num_tokens, model_parameters, gpu_flops, gpu_utilization=0.3, num_gpus=1):
    """
    估算整个评估过程的持续时间（以天为单位）。

    :param num_tokens: 训练过程中使用的tokens总数。
    :param model_parameters: 模型的参数总量。
    :param gpu_flops: GPU的峰值性能。
    :param gpu_utilization: GPU利用率。
    :param num_gpus: 使用的 GPU 数量。
    :return: 评估持续时间（天）。
    """
    seconds_per_day = 86400
    flops_per_token = 8
    return (num_tokens * model_parameters * flops_per_token) / (gpu_flops * gpu_utilization * seconds_per_day * num_gpus)
```

### 4. analyze_similar_questions_and_intents.sh - 数据 KnowHow 分析相似模块

- analyze_similar_questions_and_intents.py
- 分析数据集中的问题对，以发现文本相似度高但意图标记不同的问题对，并进一步分析这些意图是否在语义上相似
- 输出相似度高的问题对，及它们对应的意图和回答
- 本模型使用`sbert-base-chinese-nli`模型进行相似度计算，参考`模型来源`小节

#### 输入格式

```json
python analyze_similar_questions_and_intents.py --data_file_path "../../mnt/datasets/combine_data.json" --similarity_threshold 0.91 --intent_similarity_threshold 0.86 --model_path "../../mnt/models/sbert-base-chinese-nli"
```

- 输入格式，`data_file_path`是数据集文件的路径

  - `data_file_path`格式为 json，例如：

  ```json
  {
    "intruction": "",
    "input": "['怎么了？']",
    "intent": "获取问题解决方法",
    "output": "你好！ 我能为您提供什么帮助？"
  }
  ```

- `options`是配置选项：
  - `similarity_threshold`: 相似度阈值，用于判断两个问题是否相似
  - `intent_similarity_threshold`: 意图相似度阈值，用于判断两个意图是否相似
  - `model_path`: SBERT 模型的路径

#### 输出格式

```json
{
  "code": 0,
  "msg": "Success",
  "data": {
    "mismatched_intents": [
      {
        "questionPair": ["问题1", "问题2"],
        "intent1": "意图1",
        "intent2": "意图2",
        "answer1": "答案1",
        "answer2": "答案2",
        "lineNumbers": [4, 11]
      }
    ],
    "similar_intents": [
      {
        "intentPair": ["意图1", "意图2"],
        "lineNumbers": [4, 11]
      }
    ]
  }
}
```

- `mismatched_intents` 是相似但意图不同的情况
- `similar_intents` 是相似且意图也相似的情况
- `lineNumbers` 是从 0 开始统计

## 模型来源

1. Llama-2 模型

   - <https://huggingface.co/meta-llama/Llama-2-13b-chat-hf>

2. Baichuan2 模型

   - <https://huggingface.co/baichuan-inc/Baichuan2-13B-Chat>

3. ChatGLM3 模型

   - <https://huggingface.co/THUDM/chatglm3-6b>

4. 通义千问 1.5 模型

   - <https://modelscope.cn/models/qwen/Qwen1.5-14B/summary>

5. SBERT 模型

   - <https://huggingface.co/uer/sbert-base-chinese-nli>
