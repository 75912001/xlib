#!/bin/bash

# 输出运行时间
start_time=$(date +%s)
echo "脚本开始运行时间: $(date)"

# 设置工作目录为脚本所在目录
cd "$(dirname "$0")"

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请先安装 protoc: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# 检查必要的 Go 插件是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装 protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "安装 protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# 创建输出目录（如果不存在）
OutputDir=./
mkdir -p ${OutputDir}

# 生成 protobuf 代码
echo "正在生成 protobuf 代码..."
protoc \
    --go_out=${OutputDir} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${OutputDir} \
    --go-grpc_opt=paths=source_relative \
    ./*.proto

# 检查生成是否成功
if [ $? -eq 0 ]; then
    echo "✅ 代码生成成功！"
    echo "生成的文件位于: ${OutputDir}"
else
    echo "❌ 代码生成失败"
    exit 1
fi

# 任意键继续
read -n 1 -s -r -p "按任意键继续..."