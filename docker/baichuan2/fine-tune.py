import os
import math
import pathlib
import time
from enum import IntEnum, auto
from typing import Optional, Dict
from dataclasses import dataclass, field
import json

import jsonlines
import numpy as np
import requests
import torch
from torch.utils.data import Dataset
import transformers
from transformers.training_args import TrainingArguments


class SeparatorStyle(IntEnum):
    """Separator styles."""

    ADD_COLON_SINGLE = auto()
    ADD_COLON_TWO = auto()
    ADD_COLON_SPACE_SINGLE = auto()
    NO_COLON_SINGLE = auto()
    NO_COLON_TWO = auto()
    ADD_NEW_LINE_SINGLE = auto()
    LLAMA2 = auto()
    CHATGLM = auto()
    CHATML = auto()
    CHATINTERN = auto()
    DOLLY = auto()
    RWKV = auto()
    PHOENIX = auto()
    ROBIN = auto()
    FALCON_CHAT = auto()
    CHATGLM3 = auto()


@dataclass
class ModelArguments:
    model_name_or_path: Optional[str] = field(default="baichuan-inc/Baichuan2-7B-Base")


@dataclass
class DataArguments:
    data_path: str = field(
        default=None, metadata={"help": "Path to the training data."}
    )


@dataclass
class TrainingArguments(transformers.TrainingArguments):
    cache_dir: Optional[str] = field(default=None)
    optim: str = field(default="adamw_torch")
    model_max_length: int = field(
        default=512,
        metadata={
            "help": "Maximum sequence length. Sequences will be right padded (and possibly truncated)."
        },
    )
    use_lora: bool = field(default=False)


class SupervisedDataset(Dataset):
    """Dataset for supervised fine-tuning."""

    def __init__(
            self,
            data_path,
            tokenizer,
            model_max_length,
            user_tokens=[195],
            assistant_tokens=[196],
    ):
        super(SupervisedDataset, self).__init__()
        self.data = json.load(open(data_path))
        self.tokenizer = tokenizer
        self.model_max_length = model_max_length
        self.user_tokens = user_tokens
        self.assistant_tokens = assistant_tokens
        self.ignore_index = -100
        item = self.preprocessing(self.data[0])
        print("input:", self.tokenizer.decode(item["input_ids"]))
        labels = []
        for id_ in item["labels"]:
            if id_ == -100:
                continue

            labels.append(id_)
        print("label:", self.tokenizer.decode(labels))

    def __len__(self):
        return len(self.data)

    def preprocessing(self, example):
        input_ids = []
        labels = []

        for message in example["conversations"]:
            from_ = message["from"]
            value = message["value"]
            value_ids = self.tokenizer.encode(value)

            if from_ == "human":
                input_ids += self.user_tokens + value_ids
                labels += [self.tokenizer.eos_token_id] + [self.ignore_index] * len(
                    value_ids
                )
            else:
                input_ids += self.assistant_tokens + value_ids
                labels += [self.ignore_index] + value_ids
        input_ids.append(self.tokenizer.eos_token_id)
        labels.append(self.tokenizer.eos_token_id)
        input_ids = input_ids[: self.model_max_length]
        labels = labels[: self.model_max_length]
        input_ids += [self.tokenizer.pad_token_id] * (
                self.model_max_length - len(input_ids)
        )
        labels += [self.ignore_index] * (self.model_max_length - len(labels))
        input_ids = torch.LongTensor(input_ids)
        labels = torch.LongTensor(labels)
        attention_mask = input_ids.ne(self.tokenizer.pad_token_id)
        return {
            "input_ids": input_ids,
            "labels": labels,
            "attention_mask": attention_mask,
        }

    def __getitem__(self, idx) -> Dict[str, torch.Tensor]:
        return self.preprocessing(self.data[idx])


def safe_save_model_for_hf_trainer(trainer: transformers.Trainer, output_dir: str):
    """Collects the state dict and dump to disk."""
    state_dict = trainer.model.state_dict()
    if trainer.args.should_save:
        cpu_state_dict = {key: value.cpu() for key, value in state_dict.items()}
        del state_dict
        trainer._save(output_dir, state_dict=cpu_state_dict)  # noqa


# def make_supervised_data_module(
#         tokenizer: transformers.PreTrainedTokenizer, data_args, train_ratio=0.98
# ) -> Dict:
#     """Make dataset and collator for supervised fine-tuning."""
#     train_ratio = min(train_ratio, 1.0)
#     dataset_cls = (
#         SupervisedDataset if data_args.lazy_preprocess else SupervisedDataset
#     )
#     data_path = data_args.data_path
#     if data_path.endswith(".json"):
#         raw_data = json.load(open(data_path, "r"))
#     elif data_path.endswith(".jsonl"):
#         with jsonlines.open(data_path, mode="r") as reader:
#             raw_data = [item for item in reader]
#
#     # Split train/test
#     np.random.seed(0)
#     perm = np.random.permutation(len(raw_data))
#     split = int(len(perm) * train_ratio)
#     train_indices = perm[:split]
#     if train_ratio < 1:
#         eval_indices = perm[split:]
#     else:
#         # if train_ratio==1, we use 5% of data as eval data, make sure trainer will not throw error when eval data is empty
#         eval_indices = perm[-int(len(perm) * 0.05):]
#     train_raw_data = [raw_data[i] for i in train_indices]
#     eval_raw_data = [raw_data[i] for i in eval_indices]
#     # rank0_print(f"#train {len(train_raw_data)}, #eval {len(eval_raw_data)}")
#
#     train_dataset = dataset_cls(train_raw_data, tokenizer=tokenizer)
#     eval_dataset = dataset_cls(eval_raw_data, tokenizer=tokenizer)
#     return dict(train_dataset=train_dataset, eval_dataset=eval_dataset)


def train():
    parser = transformers.HfArgumentParser(
        (ModelArguments, DataArguments, TrainingArguments)
    )
    model_args, data_args, training_args = parser.parse_args_into_dataclasses()

    # 获取模型的配置
    config = transformers.AutoConfig.from_pretrained(
        model_args.model_name_or_path,
        trust_remote_code=True,
        cache_dir=training_args.cache_dir,
    )
    # Set RoPE scaling factor
    orig_ctx_len = getattr(config, "max_position_embeddings", None)
    if orig_ctx_len and training_args.model_max_length > orig_ctx_len:
        scaling_factor = float(math.ceil(training_args.model_max_length / orig_ctx_len))
        config.rope_scaling = {"type": "linear", "factor": scaling_factor}
    config.use_cache = False

    model = transformers.AutoModelForCausalLM.from_pretrained(
        model_args.model_name_or_path,
        trust_remote_code=True,
        config=config,
        cache_dir=training_args.cache_dir,
    )

    tokenizer = transformers.AutoTokenizer.from_pretrained(
        model_args.model_name_or_path,
        use_fast=False,
        config=config,
        trust_remote_code=True,
        model_max_length=training_args.model_max_length,
        cache_dir=training_args.cache_dir,
        padding_side="right",
    )
    tokenizer.pad_token = tokenizer.unk_token
    print(f"tokens len: {len(tokenizer)}")
    # model.resize_token_embeddings(len(tokenizer))

    # 是否使用lora训练
    if training_args.use_lora:
        from peft import LoraConfig, TaskType, get_peft_model

        peft_config = LoraConfig(
            task_type=TaskType.CAUSAL_LM,
            target_modules=["W_pack"],
            inference_mode=False,
            r=1,
            lora_alpha=32,
            lora_dropout=0.1,
        )
        model.enable_input_require_grads()
        model = get_peft_model(model, peft_config)
        model.print_trainable_parameters()

    # SFT 监督学习
    dataset = SupervisedDataset(
        data_args.data_path, tokenizer, training_args.model_max_length
    )
    # make_supervised_data_module(tokenizer, data_args, train_ratio=0.98)
    trainer = transformers.Trainer(
        model=model, args=training_args, train_dataset=dataset, tokenizer=tokenizer
    )
    # 断点续训
    if list(pathlib.Path(training_args.output_dir).glob("checkpoint-*")):
        trainer.train(resume_from_checkpoint=True)
    else:
        trainer.train()
    trainer.save_state()
    safe_save_model_for_hf_trainer(trainer=trainer, output_dir=training_args.output_dir)


def job_finished(status: str = "success", message: str = ""):
    job_id = os.getenv("JOB_ID")
    authorization = os.getenv("AUTH")
    url = os.getenv("CHAT_API_FINE_TUNE_URL", "http://chat-api:8000/api/v1/fine_tuning/jobs/" + job_id + "/finish")
    headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + authorization
    }
    data = {"status": status, "message": message}

    response = requests.put(url, headers=headers, json=data)

    print(response.status_code)
    print(response.json())


if __name__ == "__main__":
    status = "success"
    message = ""
    try:
        train()
    except Exception as e:
        status = "failed"
        message = str(e)
        print("发生了一个异常:", str(e))

    if local_rank <= 0:
        job_finished(status, message)