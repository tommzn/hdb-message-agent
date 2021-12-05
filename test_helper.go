package main

import (
	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
)

func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func integrationMockForTest(queue string, numberOfMessages *int) *integrationMock {
	agent := &integrationMock{
		messages:     make(map[string][]sqs.RawMessage),
		messagesSend: make(map[string][]string),
		acks:         make(map[string]string),
		flushed:      false,
		closed:       false,
		chunkSize:    10,
	}
	agent.messages[queue] = messagesForTest(numberOfMessages)
	return agent
}

func messagesForTest(numberOfMessages *int) []sqs.RawMessage {

	numOfMessages := 3
	if numberOfMessages != nil {
		numOfMessages = *numberOfMessages
	}
	messages := []sqs.RawMessage{}
	for i := 0; i < numOfMessages; i++ {
		messages = append(messages, messageForTest(nil, nil, nil))
	}
	return messages
}

func newRandomValue() *string {
	val := utils.NewId()
	return &val
}

func messageForTest(messageId, receiptHandle, body *string) sqs.RawMessage {

	if messageId == nil {
		messageId = newRandomValue()
	}
	if receiptHandle == nil {
		receiptHandle = newRandomValue()
	}
	if body == nil {
		body = newRandomValue()
	}
	return sqs.RawMessage{
		MessageId:     messageId,
		ReceiptHandle: receiptHandle,
		Body:          body,
	}
}
