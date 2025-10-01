package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	groupID := getEnv("EMAIL_SERVICE_GROUP", "email-service-group")
	env := getEnv("ENVIRONMENT", "development")

	fmt.Printf("Starting Email Service in %s environment\n", env)
	fmt.Printf("Listening to topic: %s, group: %s\n", topic, groupID)

	fmt.Println("⏳ Waiting for Kafka Group Coordinator (this can take 2-3 minutes)...")
	time.Sleep(120 * time.Second)

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBrokers},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer consumer.Close()

	fmt.Println("Kafka consumer created.")

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v", err)

			time.Sleep(10 * time.Second)
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		fmt.Printf("Parsed event: %+v\n", event)

		if eventType, ok := event["event_type"].(string); ok && eventType == "user_registered" {
			email, _ := event["email"].(string)
			userID, _ := event["user_id"].(string)

			fmt.Printf(`
--------------------------------------
[EMAIL SERVICE] Sending welcome email:
To: %s
User ID: %s
Time: %s
Environment: %s
--------------------------------------
`, email, userID, time.Now().Format(time.RFC3339), env)
		}
	}
}
