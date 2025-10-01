package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func main() {
	kafkaBrokers := getEnv("KAFKA_BROKERS", "kafka:9092")
	topic := getEnv("KAFKA_TOPIC_NOTIFICATION_EVENTS", "notification-events")
	groupID := getEnv("SMS_SERVICE_GROUP", "sms-service-group")
	env := getEnv("ENVIRONMENT", "development")

	fmt.Printf("Starting SMS Service in %s environment\n", env)
	fmt.Printf("Listening for SMS notifications on %s with group %s\n", topic, groupID)

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBrokers},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer consumer.Close()

	fmt.Println("⏳ Waiting for Kafka Group Coordinator (this can take 2-3 minutes)...")
	time.Sleep(120 * time.Second)

	fmt.Println("Kafka consumer created.")

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Failed to read message from Kafka topic %s with error: %s\n", topic, err)
			time.Sleep(10 * time.Second)

			continue
		}

		var event map[string]interface{}

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Printf("Failed to unmarshal Kafka message %s with error: %s\n", string(msg.Value), err)
			continue
		}

		fmt.Printf("Parsed event: %+v\n", event)

		if eventType, ok := event["event_type"].(string); ok && eventType == "user_registered" {
			phone, _ := event["phone"].(string)
			userID, _ := event["user_id"].(string)

			fmt.Printf(`
----------------------------------
[SMS SERVICE] Sending welcome SMS:
To: %s
User ID: %s
Time: %s
Environment: %s
----------------------------------
`, phone, userID, time.Now().Format(time.RFC3339), env)
		}

	}
}
