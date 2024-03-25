# Golang 打包基础镜像
FROM golang:1.21.5 AS build-env

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV BUILDPATH=github.com/IceBearAI/aigc
ENV GOINSECURE=github.com/IceBearAI/aigc
ENV CGO_ENABLED=0
#ENV GOOS=linux
#ENV GOARCH=amd64
RUN mkdir -p /go/src/${BUILDPATH}
COPY . /go/src/${BUILDPATH}

WORKDIR /go/src/${BUILDPATH}/

RUN go build -o /go/bin/aigc-server -ldflags="-X 'github.com/IceBearAI/aigc/cmd/service.version=$(git describe --tags --always --dirty)' \
                                           -X 'github.com/IceBearAI/aigc/cmd/service.buildDate=$(date +%FT%T%z)' \
                                           -X 'github.com/IceBearAI/aigc/cmd/service.gitCommit=$(git rev-parse --short HEAD)' \
                                           -X 'github.com/IceBearAI/aigc/cmd/service.gitBranch=$(git rev-parse --abbrev-ref HEAD)'" ./cmd/main.go

# 运行镜像
FROM alpine:latest

COPY --from=build-env /go/bin/aigc-server /usr/local/aigc-server/bin/aigc-server

WORKDIR /usr/local/aigc-server/
ENV PATH=$PATH:/usr/local/aigc-server/bin/

CMD ["aigc-server", "start"]
