import argparse
import json
import pandas as pd
import pyarrow as pa
import os

# 解析命令行参数
parser = argparse.ArgumentParser(description='JSONL to Arrow conversion.')
parser.add_argument('--base_path', type=str,
                    help='Mount point of the dataset.')
args = parser.parse_args()
base_path = args.base_path if args.base_path else ''
print(f'base_path: {base_path}')

def jsonl_to_dataframe(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        data = [json.loads(line) for line in file]
    return pd.DataFrame({
        'input': [d['messages'][0]['content'] for d in data],
        'output': [d['messages'][1]['content'] for d in data]
    })

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

def convert_dataset(jsonl_path, arrow_path):
    if os.path.exists(jsonl_path):
        df = jsonl_to_dataframe(jsonl_path)
        dataframe_to_arrow(df, arrow_path)
        create_state_json(os.path.dirname(arrow_path))

formatted_datasets_path = os.path.join(base_path, 'formatted_datasets')
os.makedirs(formatted_datasets_path, exist_ok=True)

# 创建dataset_info.json和dataset_dict.json
dataset_info = {
    "citation": "", "description": "Converted dataset from JSONL to Arrow format.",
    "features": {"input": {"dtype": "string", "_type": "Value"}, "output": {"dtype": "string", "_type": "Value"}},
    "homepage": "", "license": ""
}
dataset_dict = {"splits": ["train", "test"]}

for folder in dataset_dict['splits']:
    os.makedirs(f'{formatted_datasets_path}/{folder}', exist_ok=True)
    jsonl_path = f'{base_path}/datasets/{folder}.jsonl'
    arrow_path = f'{formatted_datasets_path}/{folder}/data-00000-of-00001.arrow'
    convert_dataset(jsonl_path, arrow_path)
    if not os.path.exists(jsonl_path):
        create_state_json(f'{formatted_datasets_path}/{folder}')  # 确保创建空的状态文件，即使数据文件不存在
    with open(f'{formatted_datasets_path}/{folder}/dataset_info.json', 'w', encoding='utf-8') as f:
        json.dump(dataset_info, f)

with open(f'{formatted_datasets_path}/dataset_dict.json', 'w', encoding='utf-8') as f:
    json.dump(dataset_dict, f)

print('jsonl to arrow finished!')
