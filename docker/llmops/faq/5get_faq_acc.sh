ps aux | grep get_faq_acc.py | awk '{print $2}' | xargs kill -9
clear
#!/bin/bash

# 进行测试的数据集
TEST_LOCAL_FILE="../datasets_example/faq_train_dataset.jsonl"

# 首先将数据集转换为新格式
formatted_datasets_path=../data/train-data/faq_formatted_datasets
mkdir -p "$formatted_datasets_path"

python3 ../convert_new_format.py \
	--test_path $TEST_LOCAL_FILE \
	--output_path "$formatted_datasets_path"

python get_faq_acc.py \
	--model_name_or_path "./output_model" \
	--data_path "$formatted_datasets_path/test_dataset.jsonl" \