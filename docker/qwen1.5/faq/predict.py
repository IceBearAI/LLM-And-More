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
import argparse
import torch
from model import MODE
import json
from tqdm import tqdm


def parse_args():
    parser = argparse.ArgumentParser()
    # Model
    parser.add_argument("--device", type=str, default="0", help="")
    parser.add_argument("--mode", type=str, default="glm3", help="")
    parser.add_argument("--model_path", type=str,
                        default="./output_model/", help="")
    parser.add_argument("--max_length", type=int, default=8192, help="")
    parser.add_argument("--do_sample", type=bool, default=True, help="")
    parser.add_argument("--top_p", type=float, default=0.8, help="")
    parser.add_argument("--temperature", type=float, default=0.0, help="")
    return parser.parse_args()


def predict_one_sample(instruction, input, model, tokenizer, args, mode):
    if "glm" in mode:
        result, _ = model.chat(tokenizer, instruction + input, max_length=args.max_length, do_sample=args.do_sample,
                               top_p=args.top_p, temperature=args.temperature)
    else:
        input_ids = []
        value = instruction + input
        if len(value) > args.max_length:
            # 当输入内容超长时，随向后截断
            value = value[:args.max_length]
        value_ids = tokenizer.encode(value)
        input_ids += value_ids
        input_ids.append(tokenizer.eos_token_id)
        input_ids = torch.tensor([input_ids], dtype=torch.long).to(
            torch.device("cuda:{}".format(args.device)))
        outputs = model.generate(input_ids, max_length=args.max_length, do_sample=args.do_sample,
                                 top_p=args.top_p, temperature=args.temperature)
        result = tokenizer.decode(
            outputs[0][len(input_ids[0]):], skip_special_tokens=True)

    return result


if __name__ == '__main__':
    args = parse_args()
    model = MODE[args.mode]["model"].from_pretrained(args.model_path, device_map="cuda:{}".format(
        args.device), torch_dtype=torch.float16, trust_remote_code=True)
    tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(
        args.model_path, trust_remote_code=True)
    print('finished model and tokenizer loading')

    with open('data/finance_test.json') as f:
        data = f.read().strip()
    for datum in data.split('\n'):
        datum = json.loads(datum)
        instruction = datum['instruction']
        input = datum['input']
        output = datum['output']

        answer = predict_one_sample(
            instruction, input, model, tokenizer, args, args.mode)
        answer_str = json.dumps(
            {'answer': answer, 'gold': output}, ensure_ascii=False) + '\n'
        filename = f'''result.txt'''
        with open(filename, 'a', encoding='utf-8') as f:
            f.write(answer_str)
