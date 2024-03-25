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
#RUN pip3 install -r requirements.txt
RUN pip3 install packaging ninja torch
#RUN pip3 install numpy flash-attn --no-build-isolation
RUN pip3 install fschat[model_worker,train] pydantic==1.10.13
#RUN pip3 install requests accelerate>=0.21 peft sentencepiece torch protobuf einops wandb anthropic>=0.3 ray
### 以下插件是支持lora微调
#RUN pip3 install deepspeed bitsandbytes scipy accelerator torchvision torch-optimizer torch.optim jsonlines xformers
#RUN pip3 install transformers_stream_generator==0.0.4 aiohttp
#RUN pip3 install torchaudio
#RUN pip3 install transformers==4.37.2
#RUN pip3 install optimum auto-gptq exllamav2
#RUN pip3 install vllm

# https://github.com/Dao-AILab/flash-attention/tree/main/csrc/layer_norm
# RUN git clone https://github.com/Dao-AILab/flash-attention && \
#    cd flash-attention/csrc/layer_norm && \
#    pip3 install .

# 模型挂载目录
#VOLUME /data

# 启动 controller
CMD ["python3", "-m", "fastchat.serve.controller", "--host", "0.0.0.0", "--port", "21001"]

## 启动 API
#CMD ["python3", "-m", "fastchat.serve.openai_api_server", "--host", "0.0.0.0", "--port", "8080", "--controller-address", "http://controller:21001"]
#
## 启动 worker
#CMD ["python3", "-m", "fastchat.serve.model_worker", "--host", "0.0.0.0", "--port", "22001",
#"--controller-address", "http://controller:21001", "--worker-address", "http://0.0.0.0:22001",
#"--model-path", "/data/qweb-14b-chat", "--model-name", "qweb-14b-chat", "--num-gpus", "2", "--max-gpu-memory", "6GiB"]