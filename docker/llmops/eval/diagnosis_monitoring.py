#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import json
import re
import fire

def extract_losses_from_log_file(log_file_path):
    # 正则表达式模式
    eval_loss_pattern = r"Evaluation finished, Average Loss: (\d+\.\d+)"
    epoch_loss_pattern = r"Epoch: (\d+), Step: (\d+), Rank: \d+, loss = (\d+\.\d+)"

    with open(log_file_path, 'r') as file:
        log_content = file.read()

    # 将内容按eval_loss_pattern切割，元素末尾包含切割本身
    split_logs = re.split(eval_loss_pattern, log_content)

    # 将切割标记和其后的内容重新组合
    split_logs_with_marker = [f"{content}{split_logs[i + 1]}" for i,
                              content in enumerate(split_logs[:-1]) if i % 2 == 0]

    results = []

    for i, log in enumerate(split_logs_with_marker):
        log_lines = log.split('\n')
        eval_loss = float(log_lines[-1].strip())

        # 获得log里最后一个匹配eval_loss_pattern
        loss_match = re.findall(epoch_loss_pattern, log)[-1]
        epoch, step, epoch_loss = int(loss_match[0]), int(
            loss_match[1]), float(loss_match[2])
        results.append({
            "eval_loss": eval_loss,
            "epoch": epoch,
            "epoch_loss": epoch_loss
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
    results = extract_losses_from_log_file(log_path)
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
    total_epoch_loss_change = 0
    total_eval_loss_change = 0
    for i in range(len(results) - 1):
        total_epoch_loss_change += results[i +
                                           1]["epoch_loss"] - results[i]["epoch_loss"]
        total_eval_loss_change += results[i +
                                          1]["eval_loss"] - results[i]["eval_loss"]
    avg_epoch_loss_change = total_epoch_loss_change / (len(results) - 1)
    avg_eval_loss_change = total_eval_loss_change / (len(results) - 1)
    # 检查过拟合
    if float(results[-1]["eval_loss"]) > float(results[-2]["eval_loss"]):
        current_risks['overfitting'] = "High"
        recommendations['overfitting'] = "建议停止训练，降低学习率，增大dropout，增大数据量，降低训练周期数，并重新训练。"
    # 检查欠拟合
    if float(results[-1]["epoch_loss"]) > _thresholds['underfitting']:
        current_risks['underfitting'] = "High"
        recommendations['underfitting'] = "建议增大学习率，减小dropout，增大训练周期数，并重新训练。"
    # 检查灾难性遗忘
    if avg_epoch_loss_change > _thresholds['catastrophic_forgetting']:
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