import json

import deepspeed
import fire
import jieba
import torch
from nltk.translate.bleu_score import sentence_bleu, SmoothingFunction
from rouge import Rouge
from tqdm import tqdm
from transformers import AutoTokenizer, AutoModelForCausalLM, LlamaTokenizer


def calculate_bleu_chinese(reference, candidate):
    ref_words = list(jieba.cut(reference))
    cand_words = list(jieba.cut(candidate))

    # 使用平滑函数
    smoothie = SmoothingFunction().method1
    return sentence_bleu([ref_words], cand_words, smoothing_function=smoothie)


def calculate_rouge_chinese(reference, candidate):
    # 将中文文本分割为字符列表
    reference_chars = list(jieba.cut(reference))
    candidate_chars = list(jieba.cut(candidate))

    # 将字符列表转换为字符串，每个字符后加空格以符合Rouge库的处理方式
    reference_joined = ' '.join(reference_chars)
    candidate_joined = ' '.join(candidate_chars)
    if not reference_joined or not candidate_joined:
        print("Warning: one of the inputs is empty, skipping")
        return 0.0
    # 计算Rouge得分
    rouge = Rouge()
    scores = rouge.get_scores(candidate_joined, reference_joined)
    return scores[0]['rouge-l']['f']


def load_model(model_name_or_path,gpu_nums):
    """
    加载模型和分词器
    """
    # 检查是否为 Llama 模型
    if "llama" in model_name_or_path.lower() or "Llama" in model_name_or_path:
        # 加载 Llama Tokenizer
        tokenizer = LlamaTokenizer.from_pretrained(
            model_name_or_path, trust_remote_code=True, padding_side="left")
        tokenizer.pad_token = tokenizer.eos_token

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


def generate_answers_batch(model, tokenizer, batch_questions, max_length, device):
    """
    使用模型批量生成答案
    """
    inputs = tokenizer(batch_questions, return_tensors="pt", padding=True,
                       truncation=True, max_length=1024).to(device)
    # 使用模型生成答案
    outputs = model.generate(
        **inputs,
        pad_token_id=tokenizer.pad_token_id,
        num_return_sequences=1,  # 每次生成一个答案
        max_new_tokens=max_length,
        eos_token_id=tokenizer.eos_token_id
    )

    results = []
    for i in range(len(batch_questions)):
        generated_text = tokenizer.decode(outputs[i], skip_special_tokens=True)
        ChatGLM_sop = "[gMASK]sop "
        if generated_text.startswith(ChatGLM_sop):
            generated_text = generated_text[len(ChatGLM_sop):]
        generated_text = generated_text[len(batch_questions[i]):]
        generated_text = generated_text.strip()
        results.append(generated_text)

    # 解码生成的答案
    return results


def write_evaluation_results_json(eval_output_path, hyperparams, evaluation_results, scores):
    """
    将评测结果以JSON格式写入文件
    """
    results = {
        "config": hyperparams,
        "overall_evaluation_metrics": evaluation_results,
        "detailed_results": scores
    }
    with open(eval_output_path, 'w', encoding='utf-8') as file:
        json.dump(results, file, indent=4, ensure_ascii=False)

    print(f"evaluation results written to {eval_output_path}")


def evaluate_model(model_name_or_path, dataset_path, evaluation_metrics, max_seq_len, per_device_batch_size, gpu_nums,
                   output_path):
    """
    真实的评估函数
    """
    # 移动模型到指定的GPU
    device = torch.device(
        f'cuda' if torch.cuda.is_available() else 'cpu')
    # 加载模型和分词器
    model, tokenizer = load_model(model_name_or_path,gpu_nums)
    # print(f"Model loaded from {model_name_or_path}, running on {device}")
    # 加载测试集
    dataset = load_dataset(dataset_path)
    # print(f"Dataset loaded from {dataset_path}")

    # 初始化评估结果
    evaluation_results = {}

    # 如果evaluation_metrics是字符串，则转换为列表
    if isinstance(evaluation_metrics, str):
        evaluation_metrics = [evaluation_metrics]

    modified_dataset = []

    for item in dataset:
        if 'messages' in item and len(item['messages']) >= 2:
            if item['messages'][0]['role'] == 'system':
                for i in range(1,len(item['messages']) - 1):
                    if item['messages'][i]['role'] == 'user' and item['messages'][i + 1]['role'] == 'assistant':
                        questions = item['messages'][0]['content'] + item['messages'][i]['content']
                        references = item['messages'][i + 1]['content']
                        modified_dataset.append({"question": questions, "reference": references})
            else:
                for i in range(len(item['messages']) - 1):
                    if item['messages'][i]['role'] == 'user' and item['messages'][i + 1]['role'] == 'assistant':
                        questions = item['messages'][i]['content']
                        references = item['messages'][i + 1]['content']
                        modified_dataset.append({"question": questions, "reference": references})
    scores = []

    # 分批处理问题
    candidates = []
    for i in tqdm(range(0, len(questions), per_device_batch_size), desc="Evaluating", unit="batch"):
        batch_questions = questions[i:i + per_device_batch_size]
        batch_references = references[i:i + per_device_batch_size]
        batch_answers = generate_answers_batch(
            model, tokenizer, batch_questions, max_seq_len, device)

        candidates.extend(batch_answers)

        batch_scores = []
        for reference, answer in zip(batch_references, batch_answers):
            if not reference or not answer:
                print("Warning: one of the inputs for reference and answer is empty, skipping")
                continue
            single_score = {
                "question": reference,
                "reference": reference,
                "model_output": answer
            }
            if "Acc" in evaluation_metrics:
                single_score["Acc"] = int(reference == answer)
            if "F1" in evaluation_metrics:
                single_score["F1"] = int(reference == answer)
            if "BLEU" in evaluation_metrics:
                single_score["BLEU"] = calculate_bleu_chinese(
                    reference, answer)
            if "Rouge" in evaluation_metrics:
                single_score["Rouge"] = calculate_rouge_chinese(
                    reference, answer)
            batch_scores.append(single_score)
        scores.extend(batch_scores)

    # 计算平均值
    for metric in evaluation_metrics:
        evaluation_results[metric] = sum(
            [score[metric] for score in scores]) / len(scores)

    # 写入评测结果
    hyperparams = {
        "model_name_or_path": model_name_or_path,
        "dataset_path": dataset_path,
        "evaluation_metrics": evaluation_metrics,
        "max_seq_len": max_seq_len,
        "per_device_batch_size": per_device_batch_size,
        "gpu_nums": gpu_nums
    }
    eval_output_path = output_path + "/eval_results.json"

    write_evaluation_results_json(eval_output_path, hyperparams, evaluation_results, scores)

    # 清理模型以释放内存
    del model
    del tokenizer
    torch.cuda.empty_cache()

    # 构建返回的JSON格式结果
    result = {
        "code": 0,
        "msg": "Success",
        "data": {
            "evaluation_results": evaluation_results
        }
    }

    return result


# 模型评估主函数
def evaluate(model_name_or_path, dataset_path, evaluation_metrics, max_seq_len, per_device_batch_size, gpu_nums,
             output_path,local_rank):
    """
    模型评估主函数
    :param model_name_or_path: 模型名称或路径
    :param dataset_path: 数据集路径
    :param evaluation_metrics: 评估指标列表
    :param max_seq_len: 最大序列长度
    :param per_device_batch_size: 每个设备的批量大小
    :param gpu_id: GPU ID
    :param output_path: 输出路径
    """
    # 调用模型评估函数并打印结果
    result = evaluate_model(
        model_name_or_path=model_name_or_path,
        dataset_path=dataset_path,
        evaluation_metrics=evaluation_metrics,
        max_seq_len=max_seq_len,
        per_device_batch_size=per_device_batch_size,
        gpu_nums=gpu_nums,
        output_path=output_path
    )

    return result


if __name__ == '__main__':
    fire.Fire(evaluate)
