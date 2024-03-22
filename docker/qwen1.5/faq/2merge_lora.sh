# 定义基座模型的路径
ori_model_dir="/home/calf/ssd/models/Baichuan2-13B-Base"
# 定义生成的LoRA模型的路径
lora_dir="epoch-4-step-0"

# 复制 "./output_model/$lora_dir"到"./output_model"
cp -r "./output_model/$lora_dir/"* "./output_model/"

# 执行merge_lora.py脚本
python merge_lora.py --ori_model_dir "$ori_model_dir" --model_dir "./output_model" --mode "baichuan2_13b"

# 复制token文件和config.json到输出目录
cp  "$ori_model_dir/token"* "./$model_dir/"
cp  "$ori_model_dir/config.json" "./$model_dir/"
