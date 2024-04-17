#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
evaluate_model_from_five_dimensions.py
模型性能评估脚本
该脚本提供了一个命令行接口，用于从五个维度评估预训练语言模型的性能
"""

import json

import deepspeed
import fire
import jieba
import torch
from nltk.translate.bleu_score import sentence_bleu, SmoothingFunction
from transformers import AutoTokenizer, AutoModelForCausalLM
from transformers import LlamaTokenizer


def load_model(model_name_or_path, gpu_nums):
    """
    加载模型和分词器
    """
    # 检查是否为 Llama 模型
    if "llama" in model_name_or_path.lower():
        # 加载 Llama Tokenizer
        tokenizer = LlamaTokenizer.from_pretrained(
            model_name_or_path, trust_remote_code=True)
    else:
        # 为其他模型加载 AutoTokenizer
        tokenizer = AutoTokenizer.from_pretrained(
            model_name_or_path, trust_remote_code=True)

    # 确保tokenizer具有 pad_token
    if tokenizer.pad_token is None:
        tokenizer.add_special_tokens({'pad_token': '[PAD]'})

    # 加载模型
    model = AutoModelForCausalLM.from_pretrained(
        model_name_or_path, trust_remote_code=True)
    ds_model = deepspeed.init_inference(
        model=model,
        mp_size=gpu_nums,
        dtype=torch.float16,
        replace_method="auto",
        replace_with_kernel_inject=True,
    )
    print(f"模型加载至设备{ds_model.module.device}\n")

    return ds_model, tokenizer


def load_dataset(dataset_path):
    """
    加载测试集
    """
    if dataset_path.endswith('.jsonl'):
        with open(dataset_path, 'r', encoding='utf-8') as file:
            dataset = [json.loads(line) for line in file.readlines()]
    else:
        with open(dataset_path, 'r', encoding='utf-8') as file:
            dataset = json.load(file)
    return dataset


def calculate_bleu_chinese(reference, candidate):
    ref_words = list(jieba.cut(reference))
    cand_words = list(jieba.cut(candidate))

    # 使用平滑函数
    smoothie = SmoothingFunction().method1
    return sentence_bleu([ref_words], cand_words, smoothing_function=smoothie)


def calculate_simply_similarity(reference, candidate):
    """计算两个简短字符串之间的相似度"""

    def levenshtein_distance(s1, s2):
        """计算两个字符串之间的编辑距离"""
        if len(s1) > len(s2):
            s1, s2 = s2, s1

        distances = range(len(s1) + 1)
        for index2, char2 in enumerate(s2):
            new_distances = [index2 + 1]
            for index1, char1 in enumerate(s1):
                if char1 == char2:
                    new_distances.append(distances[index1])
                else:
                    new_distances.append(1 + min((distances[index1], distances[index1 + 1], new_distances[-1])))
            distances = new_distances

        return distances[-1]

    def similarity(s1, s2):
        """计算两个字符串之间的相似度"""
        max_len = max(len(s1), len(s2))
        return 1 - levenshtein_distance(s1, s2) / max_len

    return similarity(reference, candidate)


def evaluate_inference_ability(model, tokenizer, device, dataset_path="./eval_datasets/inference_ability.jsonl",
                               max_seq_len=512):
    # 实现推理能力评估逻辑
    dataset = load_dataset(dataset_path)
    correct_predictions = 0
    total_predictions = len(dataset)

    for item in dataset:
        sentence1, sentence2, true_label = item['sentence1'], item['sentence2'], int(
            item['label'])
        # 将句子编码为模型所需的格式
        # 构造提示词，提示模型这是一个句子相似性判断任务
        prompt = f"判断下面两个句子是否相似：句子1：'{sentence1}' 句子2：'{sentence2}'。相似请回答'相似'，不相似请回答'不相似'。请直接回答，禁止回答思考或无关的话。"

        # 将提示词编码为模型所需的格式
        inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                           truncation=True, max_length=max_seq_len).to(device)

        # 使用模型生成答案
        outputs = model.generate(
            **inputs,
            pad_token_id=tokenizer.pad_token_id,
            num_return_sequences=1,  # 每次生成一个答案
            max_new_tokens=16,
        )

        # 解码生成的文本
        generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)
        generated_text = generated_text.strip()
        ChatGLM_sop = ["[gMASK]sop ","[gMASK] sop "]
        for sop in ChatGLM_sop:
            if generated_text.startswith(sop):
                generated_text = generated_text[len(sop):]
                break
        generated_text = generated_text[len(prompt):]
        generated_text = generated_text.strip()
        # 基于生成的文本判断句子是否相似
        prediction = 0 if "不相似" in generated_text else 1

        # 计算准确率
        if prediction == true_label:
            correct_predictions += 1

    accuracy = correct_predictions / total_predictions
    return accuracy


def evaluate_reading_comprehension(model, tokenizer, device, dataset_path="./eval_datasets/reading_comprehension.json",
                                   max_seq_len=512):
    # 实现阅读理解能力评估逻辑
    dataset = load_dataset(dataset_path)
    total_bleu_score = 0
    total_questions = 0

    for article in dataset["data"]:
        for paragraph in article["paragraphs"]:
            context = paragraph["context"]
            for qa in paragraph["qas"]:
                question = qa["question"]
                true_answers = [answer["text"] for answer in qa["answers"]]

                # 构造输入文本
                input_text = f"文章：'{context}' 问题：'{question}' 。请直接回答，禁止回答思考或无关的话。"

                # 将输入文本编码为模型所需的格式
                inputs = tokenizer(input_text, return_tensors="pt", padding=True,
                                   truncation=True).to(device)

                # 使用模型生成答案
                outputs = model.generate(
                    **inputs,
                    pad_token_id=tokenizer.pad_token_id,
                    num_return_sequences=1,  # 每次生成一个答案
                    max_new_tokens=30,
                )

                # 解码生成的文本
                generated_text = tokenizer.decode(
                    outputs[0], skip_special_tokens=True)
                generated_text = generated_text.strip()
                ChatGLM_sop = ["[gMASK]sop ","[gMASK] sop "]
                for sop in ChatGLM_sop:
                    if generated_text.startswith(sop):
                        generated_text = generated_text[len(sop):]
                        break
                generated_text = generated_text[len(input_text):]
                generated_text = generated_text.strip()

                # 计算BLEU得分
                bleu_scores = [calculate_bleu_chinese(
                    true_answer, generated_text) for true_answer in true_answers]
                avg_bleu_score = sum(bleu_scores) / len(bleu_scores)
                total_bleu_score += avg_bleu_score
                total_questions += 1
    # 计算平均BLEU得分
    average_bleu_score = total_bleu_score / \
                         total_questions if total_questions > 0 else 0
    return average_bleu_score


def evaluate_chinese_language_skill(model, tokenizer, device,
                                    dataset_path="./eval_datasets/chinese_language_skill.jsonl", max_seq_len=512):
    # 实现中文语言技能评估逻辑
    dataset = load_dataset(dataset_path)
    total_bleu_score = 0
    total_examples = len(dataset)

    for item in dataset:
        sentence, true_label_desc = item['sentence'], item['label_desc']
        keywords = item.get('keywords', '')

        # 构造提示词，提示模型这是一个新闻分类任务
        prompt = f"新闻标题：'{sentence}' 。按顺序总结新闻标题的关键词，并按英文逗号隔开。请直接中文回答，禁止回答思考或无关的话。"

        # 将提示词编码为模型所需的格式
        inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                           truncation=True, max_length=max_seq_len).to(device)

        # 使用模型生成答案
        outputs = model.generate(
            **inputs,
            pad_token_id=tokenizer.pad_token_id,
            num_return_sequences=1,  # 每次生成一个答案
            max_new_tokens=16
        )

        # 解码生成的文本
        generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)
        generated_text = generated_text.strip()
        ChatGLM_sop = ["[gMASK]sop ","[gMASK] sop "]
        for sop in ChatGLM_sop:
            if generated_text.startswith(sop):
                generated_text = generated_text[len(sop):]
                break
        generated_text = generated_text[len(prompt):]
        generated_text = generated_text.strip()
        
        # 计算BLEU得分
        bleu_score = calculate_bleu_chinese(keywords, generated_text)
        total_bleu_score += bleu_score
    # 计算平均BLEU得分
    average_bleu_score = total_bleu_score / \
                         total_examples if total_examples > 0 else 0
    return average_bleu_score


def evaluate_command_compliance(model, tokenizer, device, dataset_path="./eval_datasets/command_compliance.json",
                                max_seq_len=512):
    # 实现指令遵从能力评估逻辑
    dataset = load_dataset(dataset_path)
    total_bleu_score = 0
    total_examples = 0

    task_prefix = "Answer the given question. Your answer must be a single span in the passage."
    for item in dataset['examples']:
        # 打印进度
        passage = item['input']
        target = item['target']
        target_str = " ".join(target)

        # 构造提示词，提示模型这是一个问题生成任务
        prompt = f"{task_prefix}\n{passage}"

        # 将提示词编码为模型所需的格式
        inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                           truncation=True, max_length=128+max_seq_len).to(device)

        # 使用模型生成答案
        outputs = model.generate(
            **inputs,
            pad_token_id=tokenizer.pad_token_id,
            num_return_sequences=1,  # 每次生成一个答案
            max_new_tokens=16,
        )

        # 解码生成的文本
        generated_text = tokenizer.decode(
            outputs[0], skip_special_tokens=True)
        generated_text = generated_text.strip()
        ChatGLM_sop = ["[gMASK]sop ","[gMASK] sop "]
        for sop in ChatGLM_sop:
            if generated_text.startswith(sop):
                generated_text = generated_text[len(sop):]
                break
        generated_text = generated_text[len(prompt):]
        generated_text = generated_text.strip()
        # 对于每个真实问题，计算BLEU得分
        # 如果target_str以.结尾，去掉
        if target_str.endswith("."):
            target_str = target_str[:-1]
        bleu_score = calculate_simply_similarity(target_str, generated_text)
        total_bleu_score += bleu_score
        total_examples += 1

    # 计算平均BLEU得分
    average_bleu_score = total_bleu_score / total_examples if total_examples > 0 else 0
    return average_bleu_score


def evaluate_innovation_capacity(model, tokenizer, device, dataset_path="./eval_datasets/innovation_capacity.jsonl",
                                 max_seq_len=512):
    # 实现创新能力评估逻辑
    dataset = load_dataset(dataset_path)
    correct_predictions = 0
    total_questions = len(dataset)
    for item in dataset:
        question = item['question']
        options = {key: value for key, value in item.items() if key in [
            'A', 'B', 'C', 'D']}
        true_answer = item['answer']

        # 构造提示词，提示模型这是一个选择题答案生成任务
        prompt = f"问题是：'{question}' 选项是：{options}。请直接回答ABCD选项中的一个字母，禁止回答思考或无关的话。"

        # 将提示词编码为模型所需的格式
        inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                           truncation=True, max_length=max_seq_len).to(device)

        # 使用模型生成答案
        outputs = model.generate(
            **inputs,
            pad_token_id=tokenizer.pad_token_id,
            num_return_sequences=1,  # 每次生成一个答案
            max_new_tokens=8,
        )

        # 解码生成的文本   
        generated_text = tokenizer.decode(
            outputs[0], skip_special_tokens=True).strip()
        generated_text = generated_text.strip()
        ChatGLM_sop = ["[gMASK]sop ","[gMASK] sop "]
        for sop in ChatGLM_sop:
            if generated_text.startswith(sop):
                generated_text = generated_text[len(sop):]
                break
        generated_text = generated_text[len(prompt):]
        generated_text = generated_text.strip()
        # 基于生成的文本判断选择题答案
        prediction = ''
        if 'A' in generated_text:
            prediction = 'A'
        elif 'B' in generated_text:
            prediction = 'B'
        elif 'C' in generated_text:
            prediction = 'C'
        elif 'D' in generated_text:
            prediction = 'D'

        # 计算准确率
        if prediction == true_answer:
            correct_predictions += 1

    accuracy = correct_predictions / total_questions
    return accuracy


def main(model_name_or_path, evaluation_dimensions, output_file, options, gpu_nums,local_rank):
    """
    模型评估主函数

    :param model_name_or_path: 模型名称或路径
    :param evaluation_dimensions: 评估维度列表
    :param output_file: 结果输出到文件
    :param options: 额外选项，如最大序列长度
    :param gpu_nums: GPU数量
    """
    # 移动模型到GPU
    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
    # 加载模型和分词器
    model, tokenizer = load_model(model_name_or_path, gpu_nums)

    # 加载额外选项
    additional_parameters = options.get("additional_parameters", {})
    max_seq_len = additional_parameters.get("max_seq_len", 512)
    result_dict = {}

    if "inference_ability" in evaluation_dimensions:
        # 评估推理能力
        acc_inference_ability = evaluate_inference_ability(model, tokenizer, device, max_seq_len=max_seq_len)
        result_dict["inference_ability"] = acc_inference_ability

    if "reading_comprehension" in evaluation_dimensions:
        # 评估阅读理解能力
        bleu_reading_comprehension = evaluate_reading_comprehension(model, tokenizer, device, max_seq_len=max_seq_len)
        result_dict["reading_comprehension"] = bleu_reading_comprehension

    if "chinese_language_skill" in evaluation_dimensions:
        # 评估中文语言技能
        bleu_chinese_language_skill = evaluate_chinese_language_skill(model, tokenizer, device, max_seq_len=max_seq_len)
        result_dict["chinese_language_skill"] = bleu_chinese_language_skill

    if "command_compliance" in evaluation_dimensions:
        # 评估指令遵从能力
        bleu_command_compliance = evaluate_command_compliance(model, tokenizer, device, max_seq_len=max_seq_len)
        result_dict["command_compliance"] = bleu_command_compliance

    if "innovation_capacity" in evaluation_dimensions:
        # 评估创新能力
        accuracy_innovation_capacity = evaluate_innovation_capacity(
            model, tokenizer, device, max_seq_len=max_seq_len)
        result_dict["innovation_capacity"] = accuracy_innovation_capacity

    with open(output_file, "w", encoding="utf-8") as f:
        f.write(json.dumps(result_dict, ensure_ascii=False))
    # return json.dumps(result_dict, ensure_ascii=False, indent=4)
    # print(result_dict)


if __name__ == '__main__':
    fire.Fire(main)