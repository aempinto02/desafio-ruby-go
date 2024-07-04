package main

import (
	"ms-go/app/services/kafka"
	"ms-go/app/services/kafka/producers"
	_ "ms-go/db"
	"ms-go/router"
)

func main() {
	// Enable Kafka producer
	producers.EnableKafkaProducer = true

	// Initialize Kafka
	kafka.Initialize()

	// Start the router
	router.Run()
}
