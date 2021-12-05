package main

import (
	"errors"

	sqs "github.com/tommzn/aws-sqs"
)

type integrationMock struct {
	messages        map[string][]sqs.RawMessage
	messagesSend    map[string][]string
	acks            map[string]string
	flushed, closed bool
	chunkSize       int
}

func (mock *integrationMock) Receive(queue string) ([]sqs.RawMessage, error) {

	if queue == "error" {
		return []sqs.RawMessage{}, errors.New("Error occurs")
	}

	if messages, ok := mock.messages[queue]; ok {
		if len(messages) <= mock.chunkSize {
			delete(mock.messages, queue)
			return messages, nil
		}
		chunk, messages := messages[:mock.chunkSize], messages[mock.chunkSize:]
		mock.messages[queue] = messages
		return chunk, nil
	}
	return []sqs.RawMessage{}, nil
}

func (mock *integrationMock) Ack(queue string, receiptHandle *string) error {

	if *receiptHandle == "ack-error" {
		return errors.New("Error occurs")
	}
	mock.acks[*receiptHandle] = queue
	return nil
}

func (mock *integrationMock) send(topic string, message []byte) error {

	if topic == "error" {
		return errors.New("Error occurs")
	}

	if _, ok := mock.messagesSend[topic]; ok {
		mock.messagesSend[topic] = append(mock.messagesSend[topic], string(message))
	} else {
		mock.messagesSend[topic] = []string{string(message)}
	}
	return nil
}

func (mock *integrationMock) flush() {
	mock.flushed = true
}

func (mock *integrationMock) close() {
	mock.closed = true
}
