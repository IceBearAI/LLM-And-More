#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import json
import re
import fire

def extract_data_from_logs(log_file_path):
    # 定义正则表达式以匹配日志中的数值
    # 例如:{'rank': 0, 'epoch': 1, 'step': 1, 'loss': 17.162593841552734, 'learning_rate': 2e-05, 'eval_loss': 17.33646011352539}
    pattern = re.compile(r"'rank': (\d+), 'epoch': ([\d\.]+), 'step': (\d+), 'loss': ([\d\.]+), 'learning_rate': ([\d\.e\-]+), 'eval_loss': ([\d\.]+)")

    results = []

    with open(log_file_path, 'r') as file:
        for line in file:
            match = pattern.search(line)
            if match:
                rank, epoch, step, loss, learning_rate, eval_loss = match.groups()
                results.append({
                    'rank': int(rank),
                    'epoch': float(epoch),
                    'step': int(step),
                    'train_loss': float(loss),
                    'learning_rate': float(learning_rate),
                    'eval_loss': float(eval_loss)
                })

    return results

def diagnosis_monitoring(log_path):
    # 设置默认警报阈值
    _thresholds = {
        'overfitting': 0.02,  # 过拟合阈值：验证损失的增加率
        'underfitting': 0.5,  # 欠拟合阈值：训练损失高于此值
        'catastrophic_forgetting': 1.5  # 灾难性遗忘阈值：训练损失的突然增加率
    }
    # 从日志文件中提取损失
    results = extract_data_from_logs(log_path)
    print(f"results: {results}")
    # 初始化风险等级和建议
    current_risks = {
        'overfitting': "Low",
        'underfitting': "Low",
        'catastrophic_forgetting': "Low"
    }
    recommendations = {
        'overfitting': "-",
        'underfitting': "-",
        'catastrophic_forgetting': "-"
    }
    # 计算平均变化率
    total_train_loss_change = 0
    total_eval_loss_change = 0
    for i in range(len(results) - 1):
        total_train_loss_change += results[i +
                                           1]["train_loss"] - results[i]["train_loss"]
        total_eval_loss_change += results[i +
                                          1]["eval_loss"] - results[i]["eval_loss"]
    avg_train_loss_change = total_train_loss_change / (len(results) - 1)
    avg_eval_loss_change = total_eval_loss_change / (len(results) - 1)
    # 检查过拟合
    if float(results[-1]["eval_loss"]) > float(results[-2]["eval_loss"]):
        current_risks['overfitting'] = "High"
        recommendations['overfitting'] = "建议停止训练，降低学习率，增大dropout，增大数据量，降低训练周期数，并重新训练。"
    # 检查欠拟合
    if float(results[-1]["train_loss"]) > _thresholds['underfitting']:
        current_risks['underfitting'] = "High"
        recommendations['underfitting'] = "建议增大学习率，减小dropout，增大训练周期数，并重新训练。"
    # 检查灾难性遗忘
    if avg_train_loss_change > _thresholds['catastrophic_forgetting']:
        current_risks['catastrophic_forgetting'] = "High"
        recommendations['catastrophic_forgetting'] = "建议回退到上个版本并调整参数。"
    # 准备返回的数据
    response_data = {
        "code": 0,
        "msg": "Monitoring Report",
        "data": {
            'current_risks': current_risks,
            'recommendations': recommendations
        }
    }
    return response_data

def main(log_path):
    result = diagnosis_monitoring(log_path)
    print(json.dumps(result, ensure_ascii=False, indent=4))

if __name__ == '__main__':
    fire.Fire(main)
    # main("./eval_datasets/test_log_info.txt")