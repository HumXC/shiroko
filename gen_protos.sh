#!/usr/bin/env sh

# 在系统中安装 protoc 命令

# 安装对应 go 工具
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

PROTO_DIR="./protos"

# 遍历目录及其子目录中的所有 .proto 文件
for proto in $(find $PROTO_DIR -name "*.proto"); do
    echo "Generating Go code for $proto..."

    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        $proto
done
