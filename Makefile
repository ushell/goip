# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=goip
BINARY_UNIX=$(BINARY_NAME)_unix

# 目录
CMD_DIR=./cmd/server
CONFIG_DIR=./configs
DATA_DIR=./data

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建参数
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
.PHONY: all
all: clean deps build

# 清理
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf ./dist

# 依赖管理
.PHONY: deps
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# 构建
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(CMD_DIR)

# 交叉编译
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) $(CMD_DIR)

# 测试
.PHONY: test
test:
	$(GOTEST) -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 运行
.PHONY: run
run:
	$(GOCMD) run $(CMD_DIR)

# 格式化代码
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# 静态检查
.PHONY: lint
lint:
	@if [ -z "$(shell which golangci-lint)" ]; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# 生成proto文件
.PHONY: proto
proto:
	@if [ -z "$(shell which protoc)" ]; then \
		echo "Please install protoc first"; \
		exit 1; \
	fi
	@if [ -z "$(shell which protoc-gen-go)" ]; then \
		echo "Installing protoc-gen-go..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	fi
	@if [ -z "$(shell which protoc-gen-go-grpc)" ]; then \
		echo "Installing protoc-gen-go-grpc..."; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	fi
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

# Docker构建
.PHONY: docker-build
docker-build:
	docker build -t goip:latest .

# Docker运行
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -p 50051:50051 goip:latest

# Docker Compose运行
.PHONY: docker-compose
docker-compose:
	docker-compose -f deployments/docker-compose.yml up -d

# Docker Compose停止
.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f deployments/docker-compose.yml down

# 帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - 清理、安装依赖并构建"
	@echo "  clean        - 清理构建产物"
	@echo "  deps         - 安装依赖"
	@echo "  build        - 构建应用"
	@echo "  build-linux  - 交叉编译Linux版本"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试并生成覆盖率报告"
	@echo "  run          - 运行应用"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 运行静态检查"
	@echo "  proto        - 生成proto文件"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"
	@echo "  docker-compose - 使用Docker Compose