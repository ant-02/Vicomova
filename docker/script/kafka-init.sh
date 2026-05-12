#!/bin/bash
set -e

# Kafka 初始化脚本
# 用于手动创建必要的 topic

KAFKA_BROKER="${KAFKA_BROKER:-localhost:9092}"
TOPIC_VIDEO_VIEW="${TOPIC_VIDEO_VIEW:-video-view}"

echo "Waiting for Kafka to be ready..."
until kafka-topics.sh --bootstrap-server "$KAFKA_BROKER" --list > /dev/null 2>&1; do
    echo "Kafka not ready, waiting..."
    sleep 2
done

echo "Kafka is ready!"

# 创建 video-view topic
echo "Creating topic: $TOPIC_VIDEO_VIEW"
kafka-topics.sh --bootstrap-server "$KAFKA_BROKER" \
    --create \
    --topic "$TOPIC_VIDEO_VIEW" \
    --partitions 10 \
    --replication-factor 1 \
    --if-not-exists

echo "Topic $TOPIC_VIDEO_VIEW created or already exists"

# 列出所有 topic
echo "Current topics:"
kafka-topics.sh --bootstrap-server "$KAFKA_BROKER" --list