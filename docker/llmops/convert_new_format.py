import argparse
import json
import os

def convert_new_format_to_old(dataset_path, output_file):
    # 读取新格式数据，并构建候选类别列表
    with open(dataset_path, 'r', encoding='utf-8') as file:
        lines = file.readlines()
    
    candidate_list = set()
    for line in lines:
        data = json.loads(line)
        candidate_list.add(data['intent'])

    candidate_list = list(candidate_list)
    # 若存在"不符合"类别，则去除
    if "不符合" in candidate_list:
        candidate_list.remove("不符合")

    # 转换数据并写入新文件
    with open(output_file, 'w', encoding='utf-8') as file:
        for line in lines:
            data = json.loads(line)
            old_format_entry = {
                "instruction": data['instruction'],
                "input": "问题: " + data['question'] + "\n候选类别: " + "\n".join(candidate_list),
                "output": data['intent']
            }
            file.write(json.dumps(old_format_entry, ensure_ascii=False) + "\n")

parser = argparse.ArgumentParser(description='Dataset conversion from new to old format.')
parser.add_argument('--train_path', type=str, required=True,
                    help='Path to the train dataset in JSONL format.')
parser.add_argument('--test_path', type=str, default='',
                    help='Path to the test dataset in JSONL format (optional).')
parser.add_argument('--output_path', type=str, required=True,
                    help='Output path for the formatted datasets in JSONL format.')
args = parser.parse_args()

# 为训练集和测试集执行转换，如果提供了路径
if args.train_path:
    train_output_file = os.path.join(args.output_path, 'train_dataset.jsonl')
    convert_new_format_to_old(args.train_path, train_output_file)
    print(f'Train dataset converted and saved to {train_output_file}')
if args.test_path:
    test_output_file = os.path.join(args.output_path, 'test_dataset.jsonl')
    convert_new_format_to_old(args.test_path, test_output_file)
    print(f'Test dataset converted and saved to {test_output_file}')
