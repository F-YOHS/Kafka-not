package main

import (
	"context"
	"fmt"
	"os"
	"sync"
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
	AnalyticTopic := getEnv("KAFKA_TOPIC_ANALYTICS_EVENTS", "analytics-events")
	emailTopic := getEnv("KAFKA_TOPIC_NOTIFICATION_EVENTS", "notification-events")
	smsTopic := getEnv("KAFKA_TOPIC_NOTIFICATION_EVENTS", "notification-events")
	groupId := getEnv("ANALYTICS_SERVICE_GROUP", "analytics-service-group")
	env := getEnv("ENVIRONMENT", "development")

	fmt.Printf("Starting Analytics Service in %s environment\n", env)
	fmt.Printf("Listening to topics: %s, %s, %s\n", AnalyticTopic, emailTopic, smsTopic)

	var wg sync.WaitGroup

	allTopic := []string{AnalyticTopic, emailTopic, smsTopic}

	fmt.Println("⏳ Waiting for Kafka Group Coordinator (this can take 2-3 minutes)...")
	time.Sleep(120 * time.Second)

	for _, topic := range allTopic {
		wg.Add(1)
		go func(nowTopic string) {
			defer wg.Done()

			consumer := kafka.NewReader(kafka.ReaderConfig{
				Brokers:  []string{kafkaBrokers},
				Topic:    nowTopic,
				GroupID:  groupId,
				MaxBytes: 10e6,
				MinBytes: 10e3,
			})
			defer consumer.Close()

			fmt.Printf("Analytics: listening to %s\n", nowTopic)

			for {
				msg, err := consumer.ReadMessage(context.Background())
				if err != nil {
					fmt.Printf("Error reading message: %s\n", err)
					time.Sleep(3 * time.Second)
					continue
				}

				fmt.Printf("Message received: %s\n", string(msg.Value))
			}

		}(topic)

		wg.Wait()
	}
}
