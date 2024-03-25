FROM nvidia/cuda-12.2.2-cudnn8-devel-ubuntu22.04

LABEL maintainer="solacowa@gmail.com"

#RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list.d/debian.sources
COPY sources.list /etc/apt/sources.list

RUN apt-get update -y && apt-get install -y curl wget git vim
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get install -y tzdata \
	software-properties-common \
    git wget unzip bzip2 sudo build-essential
RUN add-apt-repository ppa:deadsnakes/ppa && \
    apt-get -y install Python3.10 Python3.10-dev

RUN apt-get install -y cuda-toolkit --fix-missing

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
RUN python3.10 get-pip.py

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
#RUN pip3 install -r requirements.txt
RUN pip3 install packaging ninja torch
RUN pip3 install numpy flash-attn --no-build-isolation

WORKDIR /app

#ENV CUDA_HOME=/usr/local/cuda
#ENV CUDA_ROOT=$CUDA_HOME
#ENV PATH=$PATH:$CUDA_ROOT/bin:$HOME/bin

COPY requirements.txt /app/requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY . /app

