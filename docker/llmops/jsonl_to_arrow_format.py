import argparse
import json
import pandas as pd
import pyarrow as pa
import os

# 解析命令行参数
parser = argparse.ArgumentParser(description='JSONL to Arrow conversion.')
parser.add_argument('--train_path', type=str, required=True,
                    help='Path to the train dataset in JSONL format.')
parser.add_argument('--test_path', type=str, default='',
                    help='Path to the test dataset in JSONL format (optional).')
parser.add_argument('--output_path', type=str, required=True,
                    help='Output path for the formatted datasets in Arrow format.')
args = parser.parse_args()


def jsonl_to_dataframe(file_path):
    if file_path and os.path.exists(file_path):
        with open(file_path, 'r', encoding='utf-8') as file:
            data = [json.loads(line) for line in file]
        return pd.DataFrame({
            'input': [d['messages'][0]['content'] for d in data],
            'output': [d['messages'][1]['content'] for d in data]
        })
    else:
        return pd.DataFrame({'input': [], 'output': []})  # Return an empty dataframe for empty or non-existent paths


def dataframe_to_arrow(df, file_path):
    table = pa.Table.from_pandas(df)
    with pa.OSFile(file_path, 'wb') as f:
        writer = pa.ipc.new_stream(f, table.schema)
        writer.write_table(table)
        writer.close()


def create_state_json(output_path, filename="data-00000-of-00001.arrow"):
    state = {
        "_data_files": [{"filename": filename}],
        "_fingerprint": "some_unique_identifier",
        "_format_columns": None,
        "_format_kwargs": {},
        "_format_type": None,
        "_output_all_columns": False,
        "_split": None
    }
    with open(os.path.join(output_path, 'state.json'), 'w') as f:
        json.dump(state, f, indent=4)


def convert_dataset(dataset_path, output_folder):
    os.makedirs(output_folder, exist_ok=True)
    df = jsonl_to_dataframe(dataset_path)
    arrow_file_path = os.path.join(output_folder, 'data-00000-of-00001.arrow')
    dataframe_to_arrow(df, arrow_file_path)
    create_state_json(output_folder)


# 创建formatted_datasets目录
os.makedirs(args.output_path, exist_ok=True)

# 创建dataset_info.json和dataset_dict.json
dataset_info = {
    "citation": "", "description": "Converted dataset from JSONL to Arrow format.",
    "features": {"input": {"dtype": "string", "_type": "Value"}, "output": {"dtype": "string", "_type": "Value"}},
    "homepage": "", "license": ""
}
dataset_dict = {"splits": ["train", "test"]}

# 处理train数据集
convert_dataset(args.train_path, os.path.join(args.output_path, 'train'))
with open(os.path.join(args.output_path, 'train', 'dataset_info.json'), 'w', encoding='utf-8') as f:
    json.dump(dataset_info, f)

# 处理test数据集，即使test_path为空也生成空的Arrow文件
print(f'detect test_path: {args.test_path}')
convert_dataset(args.test_path, os.path.join(args.output_path, 'test'))
with open(os.path.join(args.output_path, 'test', 'dataset_info.json'), 'w', encoding='utf-8') as f:
    json.dump(dataset_info, f)

# 保存dataset_dict.json
with open(os.path.join(args.output_path, 'dataset_dict.json'), 'w', encoding='utf-8') as f:
    json.dump(dataset_dict, f)

print(f'jsonl to arrow finished, output path: {args.output_path}')