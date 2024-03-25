FROM nvidia:cuda-12.2.2-runtime-ubuntu20.04

LABEL maintainer="solacowa@gmail.com"

COPY sources.list /etc/apt/sources.list

RUN apt-get update -y && apt-get install -y curl wget git vim
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get install -y tzdata \
	software-properties-common
RUN add-apt-repository ppa:deadsnakes/ppa && \
    apt-get -y install Python3.10

RUN apt-get install -y cuda-toolkit --fix-missing

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
RUN python3.10 get-pip.py
RUN apt-get install -y Python3.10-dev

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app/
COPY . /app/

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
#RUN pip3 install -r requirements.txt
RUN pip3 install packaging ninja torch
RUN pip3 install numpy flash-attn --no-build-isolation
RUN pip3 install fschat[model_worker,train] pydantic==1.10.13
RUN pip3 install requests accelerate>=0.21 peft sentencepiece torch protobuf einops wandb anthropic>=0.3 ray
## 以下插件是支持lora微调
RUN pip3 install deepspeed bitsandbytes scipy accelerator torchvision torch-optimizer torch.optim jsonlines xformers
RUN pip3 install transformers_stream_generator==0.0.4 aiohttp
RUN pip3 install torchaudio
RUN pip3 install transformers==4.32.0
RUN pip3 install optimum auto-gptq exllamav2
RUN pip3 install vllm

# https://github.com/Dao-AILab/flash-attention/tree/main/csrc/layer_norm
RUN git clone https://github.com/Dao-AILab/flash-attention && \
    cd flash-attention/csrc/layer_norm && \
    pip3 install .
