# -*- coding:utf-8 -*-
# @project: ChatGPT
# @filename: model
# @author: 刘聪NLP
# @zhihu: https://www.zhihu.com/people/LiuCongNLP
# @contact: logcongcong@gmail.com
# @time: 2023/8/6 16:13
"""
    文件说明：
            
"""
import sys,os
BASE = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, BASE)
from faq.glm3_32k.modeling_chatglm import ChatGLMForConditionalGeneration as ChatGLM3_32kForConditionalGeneration
from faq.glm3_32k.tokenization_chatglm import ChatGLMTokenizer as ChatGLM3_32kTokenizer
from faq.glm3_32k.configuration_chatglm import ChatGLMConfig as ChatGLM3_32kConfig
from faq.glm3.modeling_chatglm import ChatGLMForConditionalGeneration as ChatGLM3ForConditionalGeneration
from faq.glm3.tokenization_chatglm import ChatGLMTokenizer as ChatGLM3Tokenizer
from faq.glm3.configuration_chatglm import ChatGLMConfig as ChatGLM3Config
from faq.glm2.modeling_chatglm import ChatGLMForConditionalGeneration as ChatGLM2ForConditionalGeneration
from faq.glm2.tokenization_chatglm import ChatGLMTokenizer as ChatGLM2Tokenizer
from faq.glm2.configuration_chatglm import ChatGLMConfig as ChatGLM2Config
from faq.glm1.modeling_chatglm import ChatGLMForConditionalGeneration
from faq.glm1.tokenization_chatglm import ChatGLMTokenizer
from faq.glm1.configuration_chatglm import ChatGLMConfig
# 自动选择模型
from transformers import AutoModelForCausalLM, AutoTokenizer, AutoConfig


from utils import GLMPromptDataSet, GLM2PromptDataSet, GLM3PromptDataSet, Baichuan2For13bSupervisedDataset, SupervisedDataset

MODE = {"glm": {"model": ChatGLMForConditionalGeneration, "tokenizer": ChatGLMTokenizer, "config": ChatGLMConfig, "dataset": GLMPromptDataSet},
        "glm2": {"model": ChatGLM2ForConditionalGeneration, "tokenizer": ChatGLM2Tokenizer, "config": ChatGLM2Config, "dataset": GLM2PromptDataSet},
        "glm3": {"model": ChatGLM3ForConditionalGeneration, "tokenizer": ChatGLM3Tokenizer, "config": ChatGLM3Config, "dataset": GLM3PromptDataSet},
        "glm3_32k": {"model": ChatGLM3_32kForConditionalGeneration, "tokenizer": ChatGLM3_32kTokenizer, "config": ChatGLM3_32kConfig, "dataset": GLM3PromptDataSet},
        "baichuan2_13b": {"model": AutoModelForCausalLM, "tokenizer": AutoTokenizer, "config": AutoConfig, "dataset": Baichuan2For13bSupervisedDataset},
        "qwen1.5": {"model": AutoModelForCausalLM, "tokenizer": AutoTokenizer, "config": AutoConfig, "dataset": SupervisedDataset},
        "auto": {"model": AutoModelForCausalLM, "tokenizer": AutoTokenizer, "config": AutoConfig, "dataset": SupervisedDataset}
        }
