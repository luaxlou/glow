# Glow Makefile
# 版本信息（可以通过 make 命令行覆盖）
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date)

# Ldflags 用于注入版本信息
LDFLAGS := -ldflags "-X github.com/luaxlou/glow/cmd/glow-server/cmd.version=$(VERSION) \
                   -X github.com/luaxlou/glow/cmd/glow-server/cmd.commit=$(COMMIT) \
                   -X github.com/luaxlou/glow/cmd/glow-server/cmd.buildDate=$(BUILD_DATE) \
                   -X github.com/luaxlou/glow/cmd/glow/cmd.glowVersion=$(VERSION) \
                   -X github.com/luaxlou/glow/cmd/glow/cmd.glowCommit=$(COMMIT) \
                   -X github.com/luaxlou/glow/cmd/glow/cmd.glowBuildDate=$(BUILD_DATE)"

# 构建目录
BUILD_DIR := bin

# 平台列表
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: all build clean test install glow-server glow

all: build

## build: 编译所有二进制文件
build: glow-server glow

## glow-server: 编译 glow-server
glow-server:
	@echo "Building glow-server..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/glow-server cmd/glow-server/main.go

## glow: 编译 glow CLI
glow:
	@echo "Building glow..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/glow cmd/glow/main.go

## build-all: 为所有平台交叉编译
build-all:
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		GOOS=$(word 1,$(subst /, ,$(platform))) \
		GOARCH=$(word 2,$(subst /, ,$(platform))) \
		go build $(LDFLAGS) -o $(BUILD_DIR)/glow-server-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform))) cmd/glow-server/main.go; \
		GOOS=$(word 1,$(subst /, ,$(platform))) \
		GOARCH=$(word 2,$(subst /, ,$(platform))) \
		go build $(LDFLAGS) -o $(BUILD_DIR)/glow-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform))) cmd/glow/main.go; \
	)

## checksums: 生成所有二进制文件的 SHA256 校验和
checksums:
	@echo "Generating SHA256 checksums..."
	@cd $(BUILD_DIR) && sha256sum * > SHA256SUMS.txt

## release: 创建 release 产物（交叉编译 + 校验和）
release: build-all checksums
	@echo "Release artifacts created in $(BUILD_DIR)/"

## install: 安装到 $GOPATH/bin 或 /usr/local/bin
install:
	@echo "Installing glow-server and glow..."
	@mkdir -p $(GOPATH)/bin
	go install $(LDFLAGS) ./cmd/glow-server
	go install $(LDFLAGS) ./cmd/glow

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

## test: 运行所有测试
test:
	go test -v ./...

## lint: 运行代码检查
lint:
	golangci-lint run ./...

## help: 显示帮助信息
help:
	@echo "Glow Makefile"
	@echo ""
	@echo "使用方法:"
	@echo "  make [target]"
	@echo ""
	@echo "可用目标:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## 默认目标
.DEFAULT_GOAL := help
