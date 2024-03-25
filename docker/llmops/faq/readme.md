# FAQ场景训练指南

## 概览
本指南涵盖了FAQ场景训练、模型合并（LoRA checkpoint）、评估和打分的完整流程。请根据实际需求选择相应的模型，并按照以下步骤进行操作。

## 环境准备和参数设置
1. **参数调整**: 在`1train.sh`和`2merge_lora.sh`脚本中根据需要调整以下参数：
   - `--model_name_or_path`: 基座Base模型存储路径。
   - `--mode`: 基座Base模型类型。可选项包括：
     - `glm`: ChatGLM
     - `glm2`: ChatGLM2
     - `glm3`: ChatGLM3
     - `baichuan2_13b`: Baichuan2-13B
     - `auto`: Llama2

## 训练步骤
### 1. 训练模型 (`1train.sh`)
- **显卡选择**: 根据实际情况调整`--include localhost`后的数字，用逗号隔开表示使用的显卡编号。
- **批次大小调整**: 如果显存不足，请适当降低`--per_device_train_batch_size`的值。
- **开始训练**: 运行命令`bash 1train.sh`。训练完成后，结果将保存在`output_model`文件夹中。

### 2. 合并LoRA checkpoint (`2merge_lora.sh`)
- **参数替换**: 确保替换`--model_name_or_path`为模型的实际路径。
- **指定LoRA模型名称**: 输入`model_dir`以定义生成的LoRA模型的名称。
- **执行合并**: 运行命令`bash 2merge_lora.sh`以合并LoRA checkpoint到原始模型。

## 评估和打分
### 3. 模型评测 (`3predict.sh`)
- 在`3predict.sh`中替换`--mode`为对应模型类别，然后执行命令`bash 3predict.sh`进行模型评测。

### 4. 打分 (`4score.sh`)
- 执行命令`bash 4score.sh`进行模型性能打分。

确保在每一步中根据实际情况调整相应的参数。上述步骤将指导您完成从训练到评估的整个流程。