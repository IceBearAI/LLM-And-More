import json
import statistics
import warnings
from nltk.translate.bleu_score import corpus_bleu
from sklearn.metrics import accuracy_score, f1_score
import argparse
import torch
from model import MODE
import json
import deepspeed

from rouge import Rouge


def parse_args():
    parser = argparse.ArgumentParser()
    # Model
    parser.add_argument("--gpu_nums", type=int, default=1, help="")
    parser.add_argument("--mode", type=str, default="glm3", help="")
    parser.add_argument("--model_path", type=str,
                        default="./output_model/", help="")
    parser.add_argument("--max_length", type=int, default=8192, help="")
    parser.add_argument("--do_sample", type=bool, default=True, help="")
    parser.add_argument("--top_p", type=float, default=0.8, help="")
    parser.add_argument("--temperature", type=float, default=0.0, help="")
    parser.add_argument("--output_test", type=str, default='./data/test_result.txt', help="")
    parser.add_argument("--test_path", type=str, default='./data/test.json', help="")
    parser.add_argument("--local_rank", type=int)
    return parser.parse_args()


def predict_one_sample(instruction, document,question, model, tokenizer, args, mode, device):
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
        input_ids = torch.tensor([input_ids], dtype=torch.long).to(device)
        # print(f"数据input_ids：{input_ids}")
        outputs = model.generate(input_ids, max_length=args.max_length, do_sample=args.do_sample,
                                 top_p=args.top_p, temperature=args.temperature,
                                 # pad_token_id=tokenizer.eos_token_id,

                                 )
        result = tokenizer.decode(
            outputs[0][len(input_ids[0]):], skip_special_tokens=True)

    return result


if __name__ == '__main__':
    args = parse_args()
    device = torch.device(
        f'cuda' if torch.cuda.is_available() else 'cpu')
    model = MODE[args.mode]["model"].from_pretrained(args.model_path, torch_dtype=torch.float16, trust_remote_code=True)
    ds_model = deepspeed.init_inference(
        model=model,
        mp_size=args.gpu_nums,
        dtype=torch.float16,
        replace_method="auto",
        replace_with_kernel_inject=True,
    )
    print(f"模型加载至设备{ds_model.module.device}\n")
    tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(
        args.model_path, trust_remote_code=True)
    # tokenizer.padding_side = 'left'
    print('finished model and tokenizer loading')
    rouge = Rouge()
    rouge_scores = {'rouge-1': [], 'rouge-2': [], 'rouge-l': []}
    bleu_scores = []
    acc_scores = []
    f1_scores = []
    with open(args.output_test, 'w', encoding='utf-8') as f_out:
        with open(args.test_path,'r',encoding='utf-8') as f:
            data = f.read().strip()
        for datum in data.split('\n'):
            datum = json.loads(datum)
            instruction = f'{datum["instruction"]}\n'
            document = f'背景知识：{datum["document"]}\n'
            question = f'问题：{datum["question"]}\n'
            output = datum["output"]

            answer = predict_one_sample(instruction, document, question, ds_model, tokenizer, args, args.mode,device)

            rouge_score = rouge.get_scores(answer, output)
            for key in rouge_scores.keys():
                rouge_scores[key].append(rouge_score[0][key]['f'])

            bleu_score = corpus_bleu([[output]], [answer])
            bleu_scores.append(bleu_score)

            acc_score = accuracy_score([output], [answer])
            acc_scores.append(acc_score)

            f1_score_val = f1_score([output], [answer], average='macro')
            f1_scores.append(f1_score_val)

            result_str = json.dumps({'answer': answer, 'gold': output, 'rouge_scores': {key: rouge_score[0][key]['f'] for key in rouge_scores.keys()}}, ensure_ascii=False) + '\n'
            f_out.write(result_str)

        avg_rouge_scores = {'avg_'+key: statistics.mean(rouge_scores[key]) for key in rouge_scores.keys()}
        avg_bleu_score = statistics.mean(bleu_scores)
        avg_acc_score = statistics.mean(acc_scores)
        avg_f1_score = statistics.mean(f1_scores)
        f_out.write(f"Average Rouge Scores: {avg_rouge_scores}\n")
        f_out.write(f"Average BLEU Score: {avg_bleu_score}\n")
        f_out.write(f"Average Accuracy Score: {avg_acc_score}\n")
        f_out.write(f"Average F1 Score: {avg_f1_score}\n")

    for key, value in avg_rouge_scores.items():
        print(f'{key} score: {value}')
    print(f'Average BLEU Score: {avg_bleu_score}')
    print(f'Average Accuracy Score: {avg_acc_score}')
    print(f'Average F1 Score: {avg_f1_score}')

