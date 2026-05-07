#!/bin/bash
# entrypoint.sh - 配置同步脚本

set -e

CONFIG_FILE="/app/config/base.yaml"
ETCD_KEY="${ETCD_KEY:-/vicomova/config}"

# 安装必要工具
install_tools() {
    if command -v apk > /dev/null 2>&1; then
        apk add --no-cache inotify-tools bash > /dev/null 2>&1 || true
    fi
}

install_tools

# 上传初始配置到 etcd
upload_config() {
    if [ -f "$CONFIG_FILE" ]; then
        echo "Uploading initial config to etcd..."
        ETCDCTL_API=3 etcdctl --endpoints=http://127.0.0.1:2379 put "$ETCD_KEY" < "$CONFIG_FILE"
        echo "Config uploaded successfully"

        # 创建备份
        cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"

        # 启动配置监听
        watch_config
    else
        echo "Config file not found: $CONFIG_FILE"
    fi
}

# 监听配置变化（使用 inotifywait）
watch_config() {
    if command -v inotifywait > /dev/null 2>&1; then
        echo "Using inotifywait for real-time monitoring"
        inotifywait -m -e modify "$CONFIG_FILE" | while read path action file; do
            echo "Config changed, uploading..."
            ETCDCTL_API=3 etcdctl --endpoints=http://127.0.0.1:2379 put "$ETCD_KEY" < "$CONFIG_FILE"
            echo "Config updated at $(date +'%Y-%m-%d %H:%M:%S')"
        done
    else
        # Fallback: 轮询
        echo "Using polling (60s interval)"
        previous_hash=$(sha256sum "$CONFIG_FILE" | awk '{print $1}')
        while true; do
            sleep 60
            current_hash=$(sha256sum "$CONFIG_FILE" | awk '{print $1}')
            if [ "$current_hash" != "$previous_hash" ]; then
                echo "Config changed, uploading..."
                ETCDCTL_API=3 etcdctl --endpoints=http://127.0.0.1:2379 put "$ETCD_KEY" < "$CONFIG_FILE"
                previous_hash="$current_hash"
                echo "Config updated at $(date +'%Y-%m-%d %H:%M:%S')"
            fi
        done
    fi
}

upload_config