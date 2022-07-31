package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	config "github.com/tommzn/go-config"
)

func newKafkaClient(conf config.Config) (messagePublisher, error) {

	server := conf.Get("kafka.servers", config.AsStringPtr("localhost"))
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   *server,
		"retries":             3,
		"delivery.timeout.ms": 3000,
	})
	return &kafkaClient{
		producer: producer,
	}, err
}

func (client *kafkaClient) send(topic string, message string) error {

	deliveryChan := make(chan kafka.Event)
	client.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	event := <-deliveryChan
	switch ev := event.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return ev.TopicPartition.Error
		}
	}
	return nil

}

func (client *kafkaClient) flush() {
	client.producer.Flush(15 * 1000)
}

func (client *kafkaClient) close() {
	client.producer.Close()
}
