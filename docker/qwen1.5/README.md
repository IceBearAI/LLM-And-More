### lora训练

如果是使用Lora微调的则需要对模型进行合并

举例：

```
$ python3 -m fastchat.model.apply_lora --base /data/base-model/qwen-14b-chat \
--target /data/ft-model/ft-qwen-14b-chat-1-dxvk-pz-qa-lora-010414-merged \
--lora /data/ft-model/ft-qwen-14b-chat-1-dxvk-pz-qa-lora-010414
```

可能会报错

需要修改 ```/usr/local/lib/python3.10/dist-packages/fastchat/model/apply_lora.pyroot@ft-qwen-14b-chat-1-ovfn-pz-qa-lora-010413-757dd99ff-9pzjc:/app#```

```
base = AutoModelForCausalLM.from_pretrained(
    base_model_path, torch_dtype=torch.float16, trust_remote_code=True, low_cpu_mem_usage=True
)
base_tokenizer = AutoTokenizer.from_pretrained(base_model_path, trust_remote_code=True, use_fast=False)
```

把`trust_remote_code=False`改为`trust_remote_code=True`

另外一种方案是`finetune.py`后面加上以下代码:

```
if training_args.use_lora:
    model = AutoPeftModelForCausalLM.from_pretrained(
        model_args.model_name_or_path,  # path to the output directory
        device_map="auto",
        trust_remote_code=True
    ).eval()

    merged_model = model.merge_and_unload()
    # max_shard_size and safe serialization are not necessary.
    # They respectively work for sharding checkpoint and save the model to safetensors
    merged_model.save_pretrained(training_args.output_dir, max_shard_size="2048MB", safe_serialization=True)
    tokenizer.save_pretrained(training_args.output_dir)
```

### 微调数据格式

jsonl文件，每行一个json对象，包含以下字段：

```json lines
{
  "type": "chatml",
  "source": "unknown",
  "messages": [
    {
      "role": "system",
      "content": "You are a helpful assistant."
    },
    {
      "role": "user",
      "content": "Tell me something about large language models."
    },
    {
      "role": "assistant",
      "content": "Large language models are a type of language model that is trained on a large corpus of text data. They are capable of generating human-like text and are used in a variety of natural language processing tasks..."
    }
  ]
}
```

## 经验

GPU: A100-80Gx4

### SFT 全量微调

#### qwen1.5-14b-base

- **datasets**: 1k
- **train_batch_size**: 4
- **eval_batch_size**: 16
- **learning_rate**: 2e-5
- **max_seq_length**: 2048
- **max_steps**: 1000
- **epochs**: 6

### LORA 微调

#### qwen1.5-14b-base

- **datasets**: 1k
- **train_batch_size**: 16
- **eval_batch_size**: 32
- **learning_rate**: 3e-4
- **max_seq_length**: 2048
- **max_steps**: 1000
- **epochs**: 6
