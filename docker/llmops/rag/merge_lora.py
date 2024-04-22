import torch
from model import MODE
import argparse
from peft import PeftModel


def set_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--model_name_or_path', default="ChatGLM2-6B", type=str, help='')
    parser.add_argument('--output_dir', default="output-glm2", type=str, help='')
    parser.add_argument('--lora_dir', default="output-glm2/epoch-2-step-3900/", type=str, help='')
    parser.add_argument('--mode', default="glm2", type=str, help='')

    return parser.parse_args()

def merge_lora(args,lora_dir):
    base_model = MODE[args.mode]["model"].from_pretrained(args.model_name_or_path, torch_dtype=torch.float16, trust_remote_code=True)
    base_tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(args.model_name_or_path, trust_remote_code=True,
                                                   use_fast=False)
    lora_model = PeftModel.from_pretrained(base_model, lora_dir, torch_dtype=torch.float16, trust_remote_code=True)
    lora_model.to("cpu")
    model = lora_model.merge_and_unload()
    model.save_pretrained(args.output_dir, max_shard_size="2GB")
    model.generation_config.save_pretrained(args.output_dir)
    model.config.save_pretrained(args.output_dir)
    base_tokenizer.save_vocabulary(args.output_dir)
    base_tokenizer.save_pretrained(args.output_dir)
    print(f"Merged model saved to {args.output_dir}")

if __name__ == '__main__':
    args = set_args()
    merge_lora(args,args.lora_dir)

