FROM python:3.10.14-alpine

LABEL maintainer="solacowa@gmail.com"

RUN apk add gcc python3-dev musl-dev linux-headers

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

RUN pip3 install --no-cache-dir aiohttp fastapi httpx \
    markdown2[all] nh3 numpy prompt_toolkit>=3.0.0 \
    pydantic psutil requests rich>=10.0.0 \
    shortuuid tiktoken uvicorn

RUN pip3 install fschat
RUN pip3 install pydantic==1.10.13

WORKDIR /app

COPY . /app/

CMD ["/app/endrypoint.sh"]