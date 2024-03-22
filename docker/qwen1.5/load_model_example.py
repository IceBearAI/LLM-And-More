from transformers import AutoTokenizer, AutoModelForCausalLM, LlamaTokenizer
import math
import torch


def load_and_test_model(model_name_or_path):
    # 检查是否为 Llama 模型
    if "llama" in model_name_or_path.lower():
        # 加载 Llama Tokenizer
        tokenizer = LlamaTokenizer.from_pretrained(model_name_or_path)
    else:
        # 为其他模型加载 AutoTokenizer
        tokenizer = AutoTokenizer.from_pretrained(model_name_or_path)

    # 确保tokenizer具有 pad_token
    if tokenizer.pad_token is None:
        tokenizer.add_special_tokens({'pad_token': '[PAD]'})

    # 加载模型
    model = AutoModelForCausalLM.from_pretrained(model_name_or_path)

    # 打印验证信息
    print(f"Model {model_name_or_path} loaded successfully.")
    print(f"Tokenizer type: {type(tokenizer)}")
    print(f"Model type: {type(model)}")

    # 测试文本生成
    test_text_generation(model, tokenizer, prompt="how are you?")


def test_text_generation(model, tokenizer, prompt="how are you?"):
    # 准备初始文本
    inputs = tokenizer(prompt, return_tensors="pt")
    print(f"Input IDs: {inputs['input_ids']}")

    # 生成文本
    with torch.no_grad():
        outputs = model.generate(**inputs, max_length=50)
    print(f"Output IDs: {outputs}")
    generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)

    # 打印生成的文本
    print("Generated Text: ")
    print(generated_text)


# 调用函数，需要替换这里的路径为你的模型路径
model_name_or_path = "../mnt/Llama-2-13b-chat-hf"
load_and_test_model(model_name_or_path)
