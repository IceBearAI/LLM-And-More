FROM nvidia/cuda:12.2.2-runtime-ubuntu20.04

LABEL maintainer="solacowa@gmail.com"

#RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list.d/debian.sources
COPY sources.list /etc/apt/sources.list

RUN apt-get update -y && apt-get install -y curl wget git vim
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get install -y cuda-toolkit --fix-missing

RUN apt-get install -y --no-install-recommends \
    software-properties-common

# 添加Deadsnakes PPA
RUN add-apt-repository ppa:deadsnakes/ppa
# 安装Python 3.10
RUN apt-get install -y --no-install-recommends \
    python3.10 \
    python3.10-venv \
    python3.10-distutils \
    python3.10-dev

RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
RUN python3 get-pip.py

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

COPY pip.conf /root/.config/pip/

#RUN pip3 install -r requirements.txt
RUN pip3 install packaging ninja torch numpy \
    deepspeed transformers datasets fsspec pyarrow \
    rouge tqdm jieba nltk scikit-learn peft sentencepiece \
    sentence-transformers fire xformers protobuf \
    einops wandb anthropic>=0.3 ray accelerate requests \
    bitsandbytes scipy accelerator torchvision torch-optimizer torch.optim jsonlines \
    transformers_stream_generator torchaudio aiohttp

RUN pip3 install flash-attn --no-build-isolation
RUN pip3 install fschat[model_worker,train] pydantic==1.10.13
RUN pip3 install vllm

WORKDIR /app/

#COPY ./requirements.txt /app/requirements.txt
#RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY . /app/
