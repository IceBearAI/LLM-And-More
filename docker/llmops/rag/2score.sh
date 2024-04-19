python score.py \
	--device "0" \
	--test_path './data/test.josn' \
	--output_test './data/test_result.txt' \
	--model_path "./output_model" \
	--max_length 2048 \
	--do_sample false \
	--temperature 0.01\
	--mode "baichuan2_13b"
