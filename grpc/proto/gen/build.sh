#!/bin/bash

# 创建 protoc grpc 插件

# 输出目录的绝对路径
outputRealPath=$(realpath ../../../../message/)

############################################################

# 设置工作目录为脚本所在目录
cd "$(dirname "$0")"

echo "正在构建自定义 protoc 插件..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go 未安装"
    exit 1
fi

# 构建插件 win
GOOS=windows GOARCH=amd64 go build -o protoc-gen-go-grpc-x *.go
# 构建插件 linux
GOOS=linux GOARCH=amd64 go build -o protoc-gen-go-grpc-x-linux *.go
# 构建插件 mac
GOOS=darwin GOARCH=amd64 go build -o protoc-gen-go-grpc-x-mac *.go

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "✅ 插件构建成功！"
    #echo "插件位置: $(pwd)"
    # 复制插件到输出目录
    cp -rf protoc-gen-go-grpc-x* ${outputRealPath}
    echo "✅ 插件已复制到 ${outputRealPath} 目录"
    # 删除当前目录下的插件文件
    rm -rf protoc-gen-go-grpc-x*
else
    echo "❌ 插件构建失败"
    exit 1
fi

# 任意键继续
read -n 1 -s -r -p "按任意键继续..."