#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "Generating RPC code from proto..."

cd "$PROJECT_ROOT"

kitex -module vicomova -service user -gen-path third_party/kitex_gen api/rpc/user/user.proto

# Kitex -service 会在当前目录生成 main.go、handler.go 等脚手架，清理掉
rm -f main.go handler.go build.sh script.go 2>/dev/null || true

echo "Done!"
