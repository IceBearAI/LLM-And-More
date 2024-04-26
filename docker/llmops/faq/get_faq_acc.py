# -*- coding:utf-8 -*-
# @project: ChatGLM-Finetuning
# @filename: predict
# @author: 刘聪NLP
# @zhihu: https://www.zhihu.com/people/LiuCongNLP
# @contact: logcongcong@gmail.com
# @time: 2023/12/6 20:41
"""
    文件说明：
            
"""
import torch
import fire
import json
from tqdm import tqdm
from transformers import AutoTokenizer, AutoModelForCausalLM
from transformers import LlamaTokenizer
import deepspeed


def parse_args():
    parser = argparse.ArgumentParser()
    # Model
    parser.add_argument("--device", type=str, default="0", help="")
    parser.add_argument("--mode", type=str, default="glm3", help="")
    parser.add_argument("--model_name_or_path", type=str, default="./output_model/", help="")
    parser.add_argument("--data_path", type=str, default="./data/finance_test.json", help="")
    parser.add_argument("--max_length", type=int, default=8192, help="")
    parser.add_argument("--do_sample", type=bool, default=True, help="")
    parser.add_argument("--top_p", type=float, default=0.8, help="")
    parser.add_argument("--temperature", type=float, default=0.0, help="")
    return parser.parse_args()


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

def predict_one_sample(instruction, input, model, tokenizer, device, max_seq_len=1024, max_new_tokens=16):
    prompt = f"{instruction}\n{input}"
    # 将提示词编码为模型所需的格式
    inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                        truncation=True, max_length=max_seq_len).to(device)
    # 使用模型生成答案
    outputs = model.generate(
        **inputs,
        pad_token_id=tokenizer.pad_token_id,
        num_return_sequences=1,  # 每次生成一个答案
        max_new_tokens=max_new_tokens
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
    
    return generated_text

def main(model_name_or_path, data_path, gpu_nums=1):
    # 移动模型到GPU
    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

    # 加载模型和分词器
    model, tokenizer = load_model(model_name_or_path, gpu_nums)
    print('finished model and tokenizer loading')

    # 加载测试集
    dataset = load_dataset(data_path)

    # 初始化intent统计字典
    intent_accuracy = {}
    
    for item in dataset:
        instruction = item['instruction']
        input = item['input']
        expected_output = item['output']
        intent = expected_output

        # 预测
        actual_output = predict_one_sample(instruction, input, model, tokenizer, device)
        
        # 初始化intent统计数据
        if intent not in intent_accuracy:
            intent_accuracy[intent] = {'correct': 0, 'total': 0}
        
        # 更新统计数据
        intent_accuracy[intent]['total'] += 1
        if actual_output.strip() == expected_output.strip():
            intent_accuracy[intent]['correct'] += 1
        
        # 输出当前预测结果和正确答案
        # answer_str = json.dumps({'answer': actual_output, 'gold': expected_output}, ensure_ascii=False) + '\n'
        # print(answer_str)
    
    # 计算每个intent的准确率
    for intent in intent_accuracy:
        intent_accuracy[intent]['accuracy'] = intent_accuracy[intent]['correct'] / intent_accuracy[intent]['total']

    print(intent_accuracy)
    # 返回intent准确率字典
    return intent_accuracy

if __name__ == '__main__':
    fire.Fire(main)