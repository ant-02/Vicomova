#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "Generating RPC code from proto..."

cd "$PROJECT_ROOT"

# 禁止生成 kitex_info.yaml
export KITEX_GENERATE_INFO=false

# User service
kitex -module vicomova -service user -gen-path third_party/kitex_gen api/rpc/user/user.proto

# Video service
kitex -module vicomova -service video -gen-path third_party/kitex_gen api/rpc/video/video.proto

# Interaction service
kitex -module vicomova -service interaction -gen-path third_party/kitex_gen api/rpc/interaction/interaction.proto

# Chat service
kitex -module vicomova -service chat -gen-path third_party/kitex_gen api/rpc/chat/chat.proto

# Commerce service
kitex -module vicomova -service commerce -gen-path third_party/kitex_gen api/rpc/commerce/commerce.proto

# Kitex -service 会在当前目录生成 main.go、handler.go 等脚手架，清理掉
rm -f main.go handler.go build.sh script.go 2>/dev/null || true

echo "Done!"
