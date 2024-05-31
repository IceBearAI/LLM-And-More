FROM nvidia/cuda:12.2.2-runtime-ubuntu20.04

LABEL maintainer="solacowa@gmail.com"

#RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list.d/debian.sources
#COPY sources.list /etc/apt/sources.list

RUN apt-get update -y && apt-get install -y curl wget git vim jq
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

RUN apt-get install openmpi-bin libopenmpi-dev

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

WORKDIR /app/

RUN pip3 install packaging ninja torch install numpy bitsandbytes

# RUN pip3 install fschat
# 如果需要使用sglang vllm版本0.3.3
# RUN pip3 sglang zmq rpyc

RUN pip3 install flash-attn --no-build-isolation

COPY ./fschat/pyproject.toml /app/fschat/pyproject.toml
WORKDIR /app/fschat/
RUN pip3 install -e .

COPY ./requirements.txt /app/requirements.txt

RUN pip3 install --no-cache-dir -r requirements.txt

RUN pip3 install tensorrt
RUN pip3 install tensorrt_llm -U --pre --extra-index-url https://pypi.nvidia.com
RUN pip3 install transformers==4.41.2 --no-cache-dir
RUN pip3 install vllm==0.4.2

RUN git clone https://github.com/NVIDIA/TensorRT-LLM.git

ADD fschat /app/fschat

COPY . /app/
WORKDIR /app