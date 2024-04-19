import json
import statistics
import warnings

import argparse
import torch
from model import MODE
import json
from tqdm import tqdm
from rouge import Rouge


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
    parser.add_argument("--output_test", type=str, default='./data/test_result.txt', help="")
    parser.add_argument("--test_path", type=str, default='./data/test.json', help="")
    return parser.parse_args()


def predict_one_sample(instruction, document,question, model, tokenizer, args, mode):
    if "glm" in mode:
        result, _ = model.chat(tokenizer, instruction +document+question, max_length=args.max_length, do_sample=args.do_sample,
                               top_p=args.top_p, temperature=args.temperature)
    else:
        input_ids = []
        max_doc_len = args.max_length - len(instruction + question)
        if len(document) > max_doc_len:
            document = document[:max_doc_len]
        value = instruction + document + question
        value_ids = tokenizer.encode(value)
        input_ids += value_ids
        input_ids.append(tokenizer.eos_token_id)
        input_ids = torch.tensor([input_ids], dtype=torch.long).to(
            torch.device("cuda:{}".format(args.device)))
        print(f"数据input_ids：{input_ids}")
        outputs = model.generate(input_ids, max_length=args.max_length, do_sample=args.do_sample,
                                 top_p=args.top_p, temperature=args.temperature,
                                 # pad_token_id=tokenizer.eos_token_id,

                                 )
        result = tokenizer.decode(
            outputs[0][len(input_ids[0]):], skip_special_tokens=True)

    return result


if __name__ == '__main__':
    args = parse_args()
    model = MODE[args.mode]["model"].from_pretrained(args.model_path, device_map="cuda:{}".format(
        args.device), torch_dtype=torch.float16, trust_remote_code=True)
    tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(
        args.model_path, trust_remote_code=True)
    # tokenizer.padding_side = 'left'
    print('finished model and tokenizer loading')
    rouge = Rouge()
    scores = {'rouge-1': [], 'rouge-2': [], 'rouge-l': []}
    
    with open(args.output_test, 'w', encoding='utf-8') as f_out:
        with open(args.test_path,'r',encoding='utf-8') as f:
            data = f.read().strip()
        for datum in data.split('\n'):
            datum = json.loads(datum)
            instruction = f'{datum["instruction"]}\n'
            document = f'背景知识：{datum["document"]}\n'
            question = f'问题：{datum["question"]}\n'
            output = datum["output"]

            answer = predict_one_sample(instruction, document, question, model, tokenizer, args, args.mode)

            score = rouge.get_scores(answer, output)
            for key in scores.keys():
                scores[key].append(score[0][key]['f'])

            result_str = json.dumps({'answer': answer, 'gold': output, 'rouge_scores': {key: score[0][key]['f'] for key in scores.keys()}}, ensure_ascii=False) + '\n'
            f_out.write(result_str)

        avg_scores = {'avg_'+key: statistics.mean(scores[key]) for key in scores.keys()}
        f_out.write(f"{avg_scores}")

    for key, value in avg_scores.items():
        print(f'{key} score: {value}')

