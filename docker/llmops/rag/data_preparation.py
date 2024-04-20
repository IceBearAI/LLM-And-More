import json
import random

def data_enhancement(input_path, output_path):
    # 对训练数据进行增强
    with open(input_path, 'r', encoding='utf-8') as f_input, open(output_path, 'w', encoding='utf-8') as f_output:
        i = 0
        mods = []
        for line in f_input:
            item = json.loads(line)
            instruction = [
                "你是一个专业的客服机器人。你需要使用提供的背景知识来回答问题，请严格根据背景知识的内容回答，对于没有背景知识的信息或与问题不匹配的背景知识，直接回答“抱歉，我不知道”：",
                f"{item['instruction']}",
                "你是一个客服机器人,对于你不知道的内容直接回答“抱歉，我不知道”：",
                "你是一个专业的客服。你要使用提供的背景知识来回答问题，没有背景知识，直接回答“抱歉，我不知道”：",
            ]

            modified_item = {
                'instruction': random.choice(instruction),
                'document': item['document'],
                'question': item['question'],
                'output': item['output']
            }
            json.dump(modified_item, f_output, ensure_ascii=False)
            f_output.write('\n')
            if i % 2 == 0:
                mods.append(modified_item)
            if i % 4 == 0:
                if mods != None:
                    mod = random.choice(mods)
                    if mod['document'] != modified_item['document']:
                        document=[
                            mod['document'],
                            ' '
                        ]
                        items = {
                            'instruction': modified_item['instruction'],
                            'document': random.choice(document),
                            'question': modified_item['question'],
                            'output': '抱歉，我不知道'
                        }
                        json.dump(items, f_output, ensure_ascii=False)
                        f_output.write('\n')
            i += 1


def prepare_rag_dataset(input_path, output_path):
    with open(input_path, 'r', encoding='utf-8') as f_input, open(output_path, 'w', encoding='utf-8') as f_output:
        for line in f_input:
            item = json.loads(line)
            instruction = "你是一个专业的客服机器人。你需要使用提供的背景知识来回答问题，请严格根据背景知识的内容回答，对于没有背景知识的信息或与问题不匹配的背景知识，直接回答“抱歉，我不知道”："
            modified_item = {
                'instruction': instruction,
                'document': item['document'],
                'question': item['question'],
                'output': item['output']
            }
            json.dump(modified_item, f_output, ensure_ascii=False)
            f_output.write('\n')