python rag_demo.py \
	--device "0" \
	--doc_path './data/output.jsonl' \
	--model_path "./output_model" \
	--rag_history_path './data/rag_history.txt' \
	--max_length 8192 \
	--do_sample false \
	--temperature 0.01
	--retrieval_method 'sentence_transformers' \
	--sentence_asymmetrical_path \
	--threshold 0.68 \
	--mode "baichuan2_13b"

# 	--sentence_unsymmetrical_path \
