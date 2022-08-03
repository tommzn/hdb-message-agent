package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
)

type agent struct {
	source          messageReceiver
	target          messagePublisher
	routes          []route
	stopChan        chan bool
	logger          log.Logger
	metricPublisher metrics.Publisher
}

type route struct {
	source string
	target string
}

type kafkaClient struct {
	producer *kafka.Producer
	logger   log.Logger
}
