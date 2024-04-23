function set_cuda_devices {
    devices=$(printf ",%d" $(seq 0 $(($1-1)) | sed 's/ //g'))
    export CUDA_VISIBLE_DEVICES=${devices:1}
}
set_cuda_devices $GPUS_PER_NODE
deepspeed score.py \
	--gpu_nums $GPUS_PER_NODE \
	--test_path './data/test.josn' \
	--output_test './data/test_result.txt' \
	--model_path "./output_model" \
	--max_length 2048 \
	--do_sample false \
	--temperature 0.01\
	--mode "baichuan2_13b"
