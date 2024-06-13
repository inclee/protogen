# 项目名称
APP_NAME = protogen

# 主文件路径
MAIN_FILE = main.go

#安装路径
INSTALL_DIR = $(HOME)/.local/bin

# 构建输出目录
OUT_DIR = bin

# 构建输出文件
OUT_FILE = $(OUT_DIR)/$(APP_NAME)

# 包含所有 Go 源文件
GOFILES := $(shell find . -type f -name '*.go')

# 默认目标
.PHONY: all
all: build

# 清理生成的文件
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(OUT_DIR)

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# 构建项目
.PHONY: build
build: $(OUT_FILE)

$(OUT_FILE): $(GOFILES)
	@echo "Building..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_FILE) $(MAIN_FILE)

# 运行项目
.PHONY: run
run: build
	@echo "Running..."
	@./$(OUT_FILE)

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 静态代码分析
.PHONY: vet
vet:
	@echo "Running vet..."
	@go vet ./...

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# 安装依赖并构建项目
.PHONY: allbuild
allbuild: deps build

# 安装依赖并构建项目
.PHONY: install
install: deps allbuild
	@mv $(OUT_FILE) $(INSTALL_DIR)

# 帮助信息
.PHONY: help
help:
	@echo "Makefile for Go project"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  all         - 默认目标（构建项目）"
	@echo "  clean       - 清理生成的文件"
	@echo "  test        - 运行测试"
	@echo "  build       - 构建项目"
	@echo "  run         - 构建并运行项目"
	@echo "  fmt         - 格式化代码"
	@echo "  vet         - 运行静态代码分析"
	@echo "  deps        - 安装依赖"
	@echo "  allbuild    - 安装依赖并构建项目"
	@echo "  help        - 显示此帮助信息"