package main

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func newKafkaClient(conf config.Config, logger log.Logger) (messagePublisher, error) {

	producer, err := kafka.NewProducer(kafkaConfig(conf))
	return &kafkaClient{
		producer: producer,
		logger:   logger,
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
		client.logger.Debug(ev.String())
	}
	return nil

}

func (client *kafkaClient) flush() {
	client.producer.Flush(15 * 1000)
}

func (client *kafkaClient) close() {
	client.producer.Close()
}

func (client *kafkaClient) forwardLogs(ctx context.Context) {

	for {
		select {
		case logEvent, ok := <-client.producer.Logs():
			if ok {
				switch {
				case logEvent.Level <= 3:
					client.logger.Error("Kafka Logs: ", logEvent.Message)
				case logEvent.Level <= 6:
					client.logger.Info("Kafka Logs: ", logEvent.Message)
				default:
					client.logger.Debug("Kafka Logs: ", logEvent.Message)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func kafkaConfig(conf config.Config) *kafka.ConfigMap {

	server := conf.Get("kafka.servers", config.AsStringPtr("localhost"))
	logLevel := log.LogLevelFromConfig(conf)
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":      *server,
		"retries":                3,
		"delivery.timeout.ms":    3000,
		"go.logs.channel.enable": true,
		"log_level":              logLevel.SyslogLevel(),
	}
	if logLevel == log.Debug {
		configMap.SetKey("debug", "all")
	}
	return configMap
}
