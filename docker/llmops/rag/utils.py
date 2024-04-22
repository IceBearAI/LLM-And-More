import json
import os
import random
from typing import Dict

import numpy as np
import torch
from torch.utils.data import Dataset
from transformers import set_seed

from retrieval_side import retrieval_json, retrieve_documents, is_supported_format


def glm_load_data(data_path, tokenizer, max_len, max_src_len, is_skip, model_type, retrieval_method, st, top_k):
    all_data = []
    skip_data_number = 0
    corpus, meta, ret_model, embeddings = retrieval_json(data_path, retrieval_method, st)

    def process_sample(sample, max_len, max_src_len):
        skip_flag = False

        instruction = sample["instruction"]
        question = sample["question"]
        doc = retrieve_documents(question, corpus, meta, embeddings, ret_model, retrieval_method, top_k)
        if doc[0]['score'] < 0.5:
            document = sample["document"]
            print("the score of doc is too low")

        else:
            if top_k > 1:
                docs = [f"片段{i}:{s['doc']}" for i, s in enumerate(doc)]
                document = "\n".join(docs)
            else:
                document = doc[0]['doc']
                # meta=doc[0]['meta']
                # document=f"{document}\n数据来源：{meta}"

        encoded_doc = tokenizer.encode(document, add_special_tokens=False)
        src = ""
        if max_len < max_src_len:
            max_src_len = max_len

        if model_type == "GLM3PromptDataSet":
            src = [tokenizer.get_command("<|user|>")] + tokenizer.encode("\n", add_special_tokens=False)
        max_doc_len = max_src_len - len(
            tokenizer.tokenize("{}\n背景知识：\n问题：{}\n".format(instruction, question))) - len(src)
        if len(encoded_doc) > max_doc_len:
            skip_flag = True
            encoded_doc = encoded_doc[:max_doc_len]
            document = tokenizer.decode(encoded_doc)

        # Common processing steps
        if model_type == "GLMPromptDataSet":
            src_tokens = tokenizer.tokenize(
                "{}\n背景知识：{}\n问题：{}\n".format(instruction, document, question))
            tgt_tokens = tokenizer.tokenize(sample["output"])
            max_tgt_len = max_len - 3 - len(src_tokens)
        elif model_type == "GLM2PromptDataSet":
            src_tokens = tokenizer.tokenize(
                "{}\n背景知识：{}\n问题：{}\n".format(instruction, document, question))
            tgt_tokens = tokenizer.tokenize(sample["output"])
            max_tgt_len = max_len - 3 - len(src_tokens)
        elif model_type == "GLM3PromptDataSet":
            src_tokens = src + tokenizer.encode(f'{instruction}\n背景知识：{document}\n问题：{question}\n',
                                                add_special_tokens=False)
            tgt_tokens = [tokenizer.get_command("<|assistant|>")] + tokenizer.encode("\n",
                                                                                     add_special_tokens=False) + \
                         tokenizer.encode(sample["output"], add_special_tokens=False)
            max_tgt_len = max_len - 6 - len(src_tokens)

        if len(tgt_tokens) > max_tgt_len:
            tgt_tokens = tgt_tokens[:max_tgt_len]
            skip_flag = True

        # Specific processing steps
        if model_type == "GLMPromptDataSet":
            tokens = src_tokens + ["[gMASK]", "<sop>"] + tgt_tokens + ["<eop>"]
            input_ids = tokenizer.convert_tokens_to_ids(tokens)
            context_length = input_ids.index(tokenizer.bos_token_id)
            mask_position = context_length - 1
            labels = [-100] * context_length + input_ids[mask_position + 1:]
        elif model_type == "GLM2PromptDataSet":
            tokens = src_tokens + tgt_tokens + ["</s>"]
            assert len(tokens) <= max_len
            input_ids = [tokenizer.get_command("[gMASK]"),
                         tokenizer.get_command("sop")] + tokenizer.convert_tokens_to_ids(tokens)
            context_length = len(src_tokens) + 2
            labels = [-100] * context_length + input_ids[context_length:]
        elif model_type == "GLM3PromptDataSet":
            input_ids = [tokenizer.get_command("[gMASK]"),
                         tokenizer.get_command("sop")] + src_tokens + tgt_tokens + [tokenizer.eos_token_id]
            context_length = len(src_tokens) + 2
            labels = [-100] * context_length + input_ids[context_length:]

        assert len(input_ids) == len(labels)
        assert len(input_ids) <= max_len
        if is_skip and skip_flag:
            nonlocal skip_data_number
            skip_data_number += 1
            return None
        return {"input_ids": input_ids, "labels": labels}

    if os.path.isdir(data_path):
        for filename in os.listdir(data_path):
            if is_supported_format(filename):
                with open(os.path.join(data_path, filename), "r", encoding="utf-8") as fh:
                    for line in fh:
                        sample = json.loads(line.strip())
                        processed_sample = process_sample(sample, max_len, max_src_len)
                        if processed_sample:
                            all_data.append(processed_sample)
        if len(all_data) == 0:
            print("目录中不包含json或jsonl文件")
            return None
    else:
        if is_supported_format(data_path):
            with open(data_path, "r", encoding="utf-8") as fh:
                for line in fh:
                    sample = json.loads(line.strip())
                    processed_sample = process_sample(sample, max_len, max_src_len)
                    if processed_sample:
                        all_data.append(processed_sample)
        else:
            print("文件格式错误，请传入json或jsonl文件")
            return None
    print("the number of skipping data is {}".format(skip_data_number))
    return all_data


class GLMPromptDataSet(Dataset):
    def __init__(self, data_path, tokenizer, max_len, max_src_len, is_skip, retrieval_method, st, top_k):
        self.all_data = glm_load_data(data_path, tokenizer, max_len, max_src_len, is_skip, "GLMPromptDataSet",
                                      retrieval_method, st, top_k)

    def __len__(self):
        return len(self.all_data)

    def __getitem__(self, item):
        return self.all_data[item]


class GLM2PromptDataSet(Dataset):
    def __init__(self, data_path, tokenizer, max_len, max_src_len, is_skip, retrieval_method, st, top_k):
        self.all_data = glm_load_data(data_path, tokenizer, max_len, max_src_len, is_skip, "GLM2PromptDataSet",
                                      retrieval_method, st, top_k)

    def __len__(self):
        return len(self.all_data)

    def __getitem__(self, item):
        return self.all_data[item]


class GLM3PromptDataSet(Dataset):
    def __init__(self, data_path, tokenizer, max_len, max_src_len, is_skip, retrieval_method, st, top_k):
        self.all_data = glm_load_data(data_path, tokenizer, max_len, max_src_len, is_skip, "GLM3PromptDataSet",
                                      retrieval_method, st, top_k)

    def __len__(self):
        return len(self.all_data)

    def __getitem__(self, item):
        return self.all_data[item]


class Baichuan2For13bSupervisedDataset(Dataset):
    """Dataset for supervised fine-tuning."""

    def __init__(
            self,
            data_path,
            tokenizer,
            max_len,
            max_src_len,
            is_skip,
            retrieval_method,
            st,
            top_k=1,
            user_tokens=[195],
            assistant_tokens=[196],
    ):
        super(Baichuan2For13bSupervisedDataset, self).__init__()
        self.data = []

        if os.path.isdir(data_path):
            for filename in os.listdir(data_path):
                if is_supported_format(filename):
                    with open(os.path.join(data_path, filename), "r", encoding="utf-8") as f:
                        self.data.extend([json.loads(line) for line in f])
            if len(self.data) == 0:
                print("目录中不包含json或jsonl文件")
        else:
            if is_supported_format(data_path):
                with open(data_path, "r", encoding="utf-8") as f:
                    self.data.extend([json.loads(line) for line in f])
            else:
                print("文件格式错误，请传入json或jsonl文件")

        self.corpus, self.meta, self.ret_model, self.embeddings = retrieval_json(data_path, retrieval_method, st)
        self.top_k = top_k
        self.retrieval_method = retrieval_method
        self.tokenizer = tokenizer
        self.model_max_length = max_len
        self.max_src_len = max_src_len
        self.is_skip = is_skip
        self.user_tokens = user_tokens
        self.assistant_tokens = assistant_tokens
        self.ignore_index = -100
        # item = self.preprocessing(self.data[0])
        # print(f"item['input_ids']: {item['input_ids']}")

    def __len__(self):
        return len(self.data)

    def preprocessing(self, example):
        input_ids = []
        labels = []
        instruction = f'{example["instruction"]}\n'
        question = f'问题：{example["question"]}\n'
        output = example["output"]

        doc = retrieve_documents(question, self.corpus, self.meta, self.embeddings, self.ret_model,
                                 self.retrieval_method, top_k=self.top_k)
        if doc[0]['score'] < 0.5:
            document = f'背景知识：{example["document"]}\n'

        else:
            if self.top_k > 1:
                docs = [f"片段{i}:{s['doc']}" for i, s in enumerate(doc)]
                docs = "\n".join(docs)
                document = f'背景知识：{docs}\n'
            else:
                document = doc[0]['doc']
                # meta=doc[0]['meta']
                # document=f"{document}\n数据来源：{meta}"

        encoded_ins = self.tokenizer.encode(instruction, add_special_tokens=False)
        encoded_doc = self.tokenizer.encode(document, add_special_tokens=False)
        encoded_question = self.tokenizer.encode(question, add_special_tokens=False)
        max_doc_len = self.max_src_len - len(encoded_ins) - len(encoded_question) - len(self.user_tokens) - len(
            self.assistant_tokens) - 2
        if len(encoded_doc) > max_doc_len:
            encoded_doc = encoded_doc[:max_doc_len]

        value_ids = self.user_tokens + encoded_ins + encoded_doc + encoded_question
        output_ids = self.assistant_tokens + self.tokenizer.encode(output, add_special_tokens=False)

        input_ids += value_ids + output_ids
        labels += [self.ignore_index] * len(value_ids) + output_ids

        input_ids = input_ids[:self.model_max_length - 1] + [self.tokenizer.eos_token_id]
        labels = labels[:self.model_max_length - 1] + [self.tokenizer.eos_token_id]

        input_ids += [self.tokenizer.pad_token_id] * (self.model_max_length - len(input_ids))
        labels += [self.ignore_index] * (self.model_max_length - len(labels))

        return {
            "input_ids": input_ids,
            "labels": labels,
        }

    def __getitem__(self, idx) -> Dict[str, torch.Tensor]:
        return self.preprocessing(self.data[idx])


class SupervisedDataset(Dataset):
    """Dataset for supervised fine-tuning."""

    def __init__(
            self,
            data_path,
            tokenizer,
            max_len,
            max_src_len,
            is_skip,
            retrieval_method,
            st,
            top_k=1
    ):
        super(SupervisedDataset, self).__init__()
        # self.data = json.load(open(data_path))
        self.data = []

        if os.path.isdir(data_path):
            for filename in os.listdir(data_path):
                if is_supported_format(filename):
                    with open(os.path.join(data_path, filename), "r", encoding="utf-8") as f:
                        self.data.extend([json.loads(line) for line in f])
            if len(self.data) == 0:
                print("目录中不包含json或jsonl文件")
        else:
            if is_supported_format(data_path):
                with open(data_path, "r", encoding="utf-8") as f:
                    self.data.extend([json.loads(line) for line in f])
            else:
                print("文件格式错误，请传入json或jsonl文件")

        self.corpus, self.meta, self.ret_model, self.embeddings = retrieval_json(data_path, retrieval_method, st)
        self.retrieval_method = retrieval_method
        self.top_k = top_k
        self.tokenizer = tokenizer
        self.model_max_length = max_len
        self.max_src_len = max_src_len
        self.is_skip = is_skip
        self.ignore_index = -100
        # item = self.preprocessing(self.data[0])
        # print(f"item['input_ids']: {item['input_ids']}")
        # print("input:", self.tokenizer.decode(item["input_ids"]))
        # labels = []
        # for id_ in item["labels"]:
        #     if id_ == -100:
        #         continue
        #
        #     labels.append(id_)
        # print("label:", self.tokenizer.decode(labels))

    def __len__(self):
        return len(self.data)

    def preprocessing(self, example):
        input_ids = []
        labels = []
        instruction = f'{example["instruction"]}\n'
        question = f'问题：{example["question"]}\n'
        output = example["output"]
        doc = retrieve_documents(question, self.corpus, self.meta, self.embeddings, self.ret_model,
                                 self.retrieval_method, top_k=self.top_k)
        if doc[0]['score'] < 0.5:
            document = f'背景知识：{example["document"]}\n'

        else:
            if self.top_k > 1:
                docs = [f"片段{i}:{s['doc']}" for i, s in enumerate(doc)]
                docs = "\n".join(docs)
                document = f'背景知识：{docs}\n'
            else:
                document = doc[0]['doc']
                # meta=doc[0]['meta']
                # document=f"{document}\n数据来源：{meta}"

        encoded_ins = self.tokenizer.encode(instruction, add_special_tokens=False)
        encoded_doc = self.tokenizer.encode(document, add_special_tokens=False)
        encoded_question = self.tokenizer.encode(question, add_special_tokens=False)
        max_doc_len = self.max_src_len - len(encoded_ins) - len(encoded_question) - 2

        if len(encoded_doc) > max_doc_len:
            encoded_doc = encoded_doc[:max_doc_len]

        value_ids = encoded_ins + encoded_doc + encoded_question
        output_ids = self.tokenizer.encode(output, add_special_tokens=False)

        input_ids += value_ids + [self.tokenizer.eos_token_id] + output_ids
        labels += [self.ignore_index] * len(value_ids) + [self.ignore_index] + output_ids

        input_ids = input_ids[:self.model_max_length - 1] + [self.tokenizer.eos_token_id]
        labels = labels[:self.model_max_length - 1] + [self.tokenizer.eos_token_id]

        # if hasattr(self.tokenizer, "pad_token_id") and self.tokenizer.pad_token_id is not None:
        #     pad_token_id = self.tokenizer.pad_token_id
        # else:
        #     pad_token_id = 2
        # input_ids += [pad_token_id] * (self.model_max_length - len(input_ids))
        # labels += [self.ignore_index] * (self.model_max_length - len(labels))

        return {
            "input_ids": input_ids,
            "labels": labels,
        }

    def __getitem__(self, idx) -> Dict[str, torch.Tensor]:
        return self.preprocessing(self.data[idx])


class DataCollator(object):
    def __init__(self, tokenizer):
        self.tokenizer = tokenizer
        if hasattr(tokenizer, "pad_token_id") and tokenizer.pad_token_id is not None:
            self.pad_token_id = tokenizer.pad_token_id
        else:
            self.pad_token_id = 2
        print(f"pad_token_id: {self.pad_token_id}")

    def __call__(self, batch):
        lengths = [len(instance["input_ids"]) for instance in batch]
        batch_max_len = max(lengths)

        input_ids_batch, labels_batch = [], []
        for instance in batch:
            input_ids = instance["input_ids"]
            labels = instance["labels"]

            padding_len = batch_max_len - len(input_ids)
            input_ids = input_ids + [self.pad_token_id] * padding_len
            labels = labels + [-100] * padding_len

            input_ids_batch.append(input_ids)
            labels_batch.append(labels)

        return {"input_ids": torch.tensor(input_ids_batch, dtype=torch.long),
                "labels": torch.tensor(labels_batch, dtype=torch.long)}


def print_trainable_parameters(model):
    trainable_params = 0
    all_param = 0
    for _, param in model.named_parameters():
        num_params = param.numel()
        if num_params == 0 and hasattr(param, "ds_numel"):
            num_params = param.ds_numel

        all_param += num_params
        if param.requires_grad:
            trainable_params += num_params
    print("trainable params: {} || all params: {} || trainable%: {}".format(trainable_params, all_param,
                                                                            100 * trainable_params / all_param))


def print_rank_0(msg, rank=0):
    if rank <= 0:
        print(msg)


def to_device(batch, device):
    output = {}
    for k, v in batch.items():
        try:
            output[k] = v.to(device)
        except:
            output[k] = v
    return output


def set_random_seed(seed):
    if seed is not None:
        set_seed(seed)
        random.seed(seed)
        np.random.seed(seed)
        torch.manual_seed(seed)
        torch.cuda.manual_seed_all(seed)


def save_model(model, tokenizer, output_dir, model_name, state_dict=None):
    save_dir = os.path.join(output_dir, model_name)
    if state_dict == None:
        model.save_pretrained(save_dir, torch_dtype=torch.float16)
    else:
        model.save_pretrained(
            save_dir, state_dict=state_dict, torch_dtype=torch.float16)
    tokenizer.save_pretrained(save_dir)
