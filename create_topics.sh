#!/bin/bash
set -e

KAFKA_BROKERS=${KAFKA_BROKERS:-kafka:9092}
USER_TOPIC=${KAFKA_TOPIC_USER_EVENTS:-user-events}
PAYMENT_TOPIC=${KAFKA_TOPIC_PAYMENT_EVENTS:-payment-events}
NOTIFICATION_TOPIC=${KAFKA_TOPIC_NOTIFICATION_EVENTS:-notification-events}
ANALYTICS_TOPIC=${KAFKA_TOPIC_ANALYTICS_EVENTS:-analytics-events}
REPLICATION_FACTOR=${REPLICATION_FACTOR:-1}
PARTITIONS=${PARTITIONS:-3}

echo "Waiting for Kafka to be ready on $KAFKA_BROKERS..."

max_attempts=30
attempt=1
while [[ $attempt -le $max_attempts ]]; do
    if kafka-topics --bootstrap-server "$KAFKA_BROKERS" --list >/dev/null 2>&1; then
        echo "Kafka is ready!"
        break
    fi

    echo "Attempt $attempt/$max_attempts: Kafka not ready, waiting..."
    sleep 5
    ((attempt++))

    if [[ $attempt -gt $max_attempts ]]; then
        echo "Error: Kafka not ready after $max_attempts attempts"
        exit 1
    fi
done

echo "Creating topics on broker: $KAFKA_BROKERS"

for topic in "$USER_TOPIC" "$PAYMENT_TOPIC" "$NOTIFICATION_TOPIC" "$ANALYTICS_TOPIC"; do
    echo "Creating topic: $topic"
    if kafka-topics --bootstrap-server "$KAFKA_BROKERS" --list | grep -q "^$topic$"; then
        echo "Topic $topic already exists, skipping..."
    else
        kafka-topics --bootstrap-server "$KAFKA_BROKERS" --create \
            --topic "$topic" \
            --partitions "$PARTITIONS" \
            --replication-factor "$REPLICATION_FACTOR"
        echo "Topic $topic created successfully"
    fi
done

echo "All topics processed successfully"