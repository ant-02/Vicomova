#!/bin/bash
# config-sync.sh - 监听 config/base.yaml 变化并同步到 etcd

set -e

CONFIG_FILE="${1:-/app/config/base.yaml}"
ETCD_ENDPOINTS="${ETCD_ENDPOINTS:-127.0.0.1:2379}"
CONFIG_PREFIX="${CONFIG_PREFIX:-/vicomova/config}"

# 确保 etcdctl 可用
if ! command -v etcdctl &> /dev/null; then
    echo "Installing etcdctl..."
    go install go.etcd.io/etcd/client/v3@latest
fi

export ETCDCTL_API=3
export ETCDCTL_ENDPOINTS=$ETCD_ENDPOINTS

echo "Watching $CONFIG_FILE for changes..."

# 初始同步
sync_to_etcd() {
    local key="$1"
    local value="$2"
    local full_key="${CONFIG_PREFIX}/${key}"
    echo "Syncing: $full_key = $value"
    etcdctl put "$full_key" "$value"
}

# 读取 yaml 并同步（简单实现，仅支持顶层 key）
sync_yaml_to_etcd() {
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Config file not found: $CONFIG_FILE"
        return
    fi

    # 读取 jwt.secret
    local jwt_secret=$(grep "secret:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
    if [ -n "$jwt_secret" ]; then
        sync_to_etcd "jwt/secret" "$jwt_secret"
    fi

    # 读取 redis.host
    local redis_host=$(grep -A2 "redis:" "$CONFIG_FILE" | grep "host:" | awk '{print $2}')
    if [ -n "$redis_host" ]; then
        sync_to_etcd "redis/host" "$redis_host"
    fi
}

# 监听文件变化 (使用 inotifywait 或 poll)
if command -v inotifywait &> /dev/null; then
    # Linux: 使用 inotify
    inotifywait -m -e modify "$CONFIG_FILE" | while read path action file; do
        echo "File changed: $path$file"
        sync_yaml_to_etcd
    done
else
    # macOS / fallback: 轮询
    echo "Using polling (install inotifywait for real-time on Linux)"
    last_mtime=""
    while true; do
        current_mtime=$(stat -f "%m" "$CONFIG_FILE" 2>/dev/null || stat -c "%Y" "$CONFIG_FILE" 2>/dev/null)
        if [ "$current_mtime" != "$last_mtime" ]; then
            last_mtime="$current_mtime"
            sync_yaml_to_etcd
        fi
        sleep 2
    done
fi