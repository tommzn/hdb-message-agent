package main

import (
	sqs "github.com/tommzn/aws-sqs"
)

type messageReceiver interface {
	Receive(string) ([]sqs.RawMessage, error)
	Ack(string, *string) error
}

type messagePublisher interface {
	send(string, string) error
	flush()
	close()
}
