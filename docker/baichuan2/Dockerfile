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

COPY requirements.txt /app/requirements.txt
RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN pip3 install -r requirements.txt

COPY . /app/

