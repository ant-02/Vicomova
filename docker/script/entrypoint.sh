#!/bin/bash
set -e

CONFIG_FILE="/app/config/base.yaml"
ETCD_KEY="${ETCD_KEY:-/vicomova/config}"
ETCD_ENDPOINTS="${ETCD_ENDPOINTS:-etcd:2379}"

echo "Using ETCD_ENDPOINTS: $ETCD_ENDPOINTS"
echo "Using ETCD_KEY: $ETCD_KEY"

upload_config() {
    if [ -f "$CONFIG_FILE" ]; then
        echo "Uploading initial config to etcd..."
        ETCDCTL_API=3 etcdctl --endpoints=http://$ETCD_ENDPOINTS put "$ETCD_KEY" < "$CONFIG_FILE"
        echo "Config uploaded successfully"

        watch_config
    else
        echo "Config file not found: $CONFIG_FILE"
    fi
}

watch_config() {
    if command -v inotifywait > /dev/null 2>&1; then
        echo "Using inotifywait for real-time monitoring"
        inotifywait -m -e modify "$CONFIG_FILE" | while read path action file; do
            echo "Config changed, uploading..."
            ETCDCTL_API=3 etcdctl --endpoints=http://$ETCD_ENDPOINTS put "$ETCD_KEY" < "$CONFIG_FILE"
            echo "Config updated at $(date +'%Y-%m-%d %H:%M:%S')"
        done
    else
        echo "Using polling (60s interval)"
        previous_hash=$(sha256sum "$CONFIG_FILE" | awk '{print $1}')
        while true; do
            sleep 60
            current_hash=$(sha256sum "$CONFIG_FILE" | awk '{print $1}')
            if [ "$current_hash" != "$previous_hash" ]; then
                echo "Config changed, uploading..."
                ETCDCTL_API=3 etcdctl --endpoints=http://$ETCD_ENDPOINTS put "$ETCD_KEY" < "$CONFIG_FILE"
                previous_hash="$current_hash"
                echo "Config updated at $(date +'%Y-%m-%d %H:%M:%S')"
            fi
        done
    fi
}

upload_config
