APPNAME = aigc-server
BIN = $(if $(GOPATH),$(GOPATH)/bin,./)
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/$(APPNAME)
GOPROXY = https://goproxy.cn
GOINSTALL = $(GOCMD) install
PID = .pid
VERSION = $(shell git describe --tags --always --dirty)
GO_LDFLAGS = -ldflags="-X 'github.com/IceBearAI/aigc/cmd/service.version=$(VERSION)' -X 'github.com/IceBearAI/aigc/cmd/service.buildDate=$(shell date +%FT%T%z)' -X 'github.com/IceBearAI/aigc/cmd/service.gitCommit=$(shell git rev-parse --short HEAD)' -X 'github.com/IceBearAI/aigc/cmd/service.gitVersion=$(shell git version)' -X 'github.com/IceBearAI/aigc/cmd/service.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)'"

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: account tenant cronjob-start

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=$(GOPROXY) go build -v -o $(BINARY_UNIX) $(GO_LDFLAGS) ./cmd/main.go

build:
	GO111MODULE=on GOPROXY=$(GOPROXY) go build -v -o $(BINARY_UNIX) $(GO_LDFLAGS) ./cmd/main.go

run:
	GOPROXY=$(GOPROXY) GO111MODULE=on go run ./cmd/main.go start -p :8080 -a $(APPNAME)

run-api:
	GOPROXY=$(GOPROXY) GO111MODULE=on go run ./cmd/main.go start-api -p :8081 -a $(APPNAME)-api

generate:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go generate table all

job-finetuning-run-waiting-train:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go job finetuning run-waiting-train

job-finetuning-running-log:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go job finetuning running-log

cronjob-start:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go cronjob start $(filter-out $@,$(MAKECMDGOALS))

account:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go account $(filter-out $@,$(MAKECMDGOALS))

# 定义 'tenant' 目标
tenant:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go tenant $(filter-out $@,$(MAKECMDGOALS))

# 这个规则匹配任何目标，并且什么也不做
# 它防止了 make 因为找不到目标规则而失败
%:
	@: