# FastChat 模型推理镜像

> [FastChat](https://github.com/lm-sys/FastChat) 是一个开放平台，用于训练、服务和评估基于大型语言模型的聊天机器人。

## 简介

**FastChat我们主要用其三个服务**

`controller` 用于模型的注册中心及健康检查

`worker` 服务启动模型并将当前模型注册到controller

`api` 从controller获取模型的地址代理到worker并提供标准API

我们主要通过它来实现大模型的高可用，高可扩展性。

### fschat-controller

主要用于模型的注册中心及健康检查

### fschat-api

提供类似OpenAI的标准接口

### fschat-worker

用于部署在各个节点上，用于模型的推理

![img.png](https://github.com/lm-sys/FastChat/raw/main/assets/server_arch.png)

## 打包成docker镜像

```
$ docker build --rm -t fschat:latest .
```