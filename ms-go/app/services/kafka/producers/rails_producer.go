package producers

import (
	"encoding/json"
	"fmt"
	"log"
	"ms-go/app/models"
	"ms-go/app/services/kafka"
)

var (
	EnableKafkaProducer = false
)

func SendToRailsKafkaTopic(product models.Product) {
	if !EnableKafkaProducer {
		fmt.Println("Kafka producer is disabled in test mode.")
		return
	}

	// Serialize struct to JSON
	jsonBytes, err := json.Marshal(product)
	if err != nil {
		log.Fatalf("Failed to serialize product to JSON: %v", err)
	}

	// Produce message to Rails Kafka topic
	kafka.SendMessage(jsonBytes)

	fmt.Println("Message produced successfully to Kafka topic: go_to_rails")
}
