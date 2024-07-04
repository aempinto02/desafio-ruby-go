package consumers

import (
	"encoding/json"
	"log"
	"ms-go/app/models"
	"ms-go/app/services/products"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

func StartConsumer() {
	// Set up configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0 // Adjust according to your Kafka version

	// Define brokers and topics
	brokers := []string{"kafka:29092"}
	topic := "rails-to-go"

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing consumer: %v", err)
		}
	}()

	// Subscribe to the topic
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalf("Error closing partition consumer: %v", err)
		}
	}()

	// Handle signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)
			case msg := <-partitionConsumer.Messages():
				processMessage(msg.Value)
			case <-signals:
				log.Println("Received shutdown signal, shutting down consumer...")
				return
			}
		}
	}()

	log.Println("Consumer started, listening for messages...")

	// Wait for shutdown signal
	wg.Wait()
}

func processMessage(value []byte) {
	log.Printf("Processing message: %s\n", string(value))

	var product models.Product
	err := json.Unmarshal(value, &product)
	if err != nil {
		log.Printf("Error unmarshaling Kafka message: %v", err)
		return
	}

	// Upsert the product from Rails-to-Go Kafka message
	updatedProduct, err := products.UpsertProduct(product)
	if err != nil {
		log.Printf("Failed to upsert product: %s", err)
		return
	}

	log.Printf("Upserted product from rails-to-go Kafka topic: %+v", updatedProduct)

}
