package kafka

import (
	"fmt"
	"log"
	"ms-go/app/services/kafka/consumers"
	"time"

	"github.com/IBM/sarama"
)

func Initialize() {
	go consumers.StartConsumer()
}

var (
	brokerAddress = "kafka:29092"
	topic         = "go_to_rails"
)

func SendMessage(message []byte) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionNone
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Failed to close Kafka producer: %v", err)
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
