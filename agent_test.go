package main

import (
	"context"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type AgentTestSuite struct {
	suite.Suite
	mock       *integrationMock
	agent      *agent
	ctx        context.Context
	cancelFunc func()
}

func TestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}

func (suite *AgentTestSuite) init(queue string, numberOfMessages *int) {

	suite.mock = integrationMockForTest(queue, numberOfMessages)
	var err error
	suite.agent, err = newAgent(loadConfigForTest(nil), loggerForTest())
	suite.Nil(err)
	suite.agent.source = suite.mock
	suite.agent.target = suite.mock
	suite.ctx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *AgentTestSuite) TestProcessMessages() {

	suite.init("test-queue", nil)

	go suite.agent.Run(suite.ctx, suite.waitGroupForTest())
	time.Sleep(500 * time.Millisecond)
	suite.cancelFunc()

	suite.True(suite.mock.flushed)
	suite.True(suite.mock.closed)
	messagesSend, ok := suite.mock.messagesSend["test-topic"]
	suite.True(ok)
	suite.Len(messagesSend, 3)
}

func (suite *AgentTestSuite) TestStopProcessing() {

	queue := "test-queue"
	numberOfMessages := 300000
	suite.init(queue, &numberOfMessages)
	suite.mock.chunkSize = 30000

	go suite.agent.Run(suite.ctx, suite.waitGroupForTest())
	time.Sleep(100 * time.Millisecond)
	suite.cancelFunc()

	time.Sleep(3 * time.Second)
	suite.True(suite.mock.flushed)
	suite.True(suite.mock.closed)
	messagesSend, ok := suite.mock.messagesSend["test-topic"]
	suite.True(ok)
	suite.True(len(messagesSend) > 10 && len(messagesSend) < numberOfMessages)
	suite.True(len(suite.mock.acks) > 10)
}

func (suite *AgentTestSuite) TestMessageSendFailed() {

	suite.init("test-queue-2", nil)

	go suite.agent.Run(suite.ctx, suite.waitGroupForTest())
	time.Sleep(500 * time.Millisecond)
	suite.cancelFunc()

	suite.True(suite.mock.flushed)
	suite.True(suite.mock.closed)
	suite.Len(suite.mock.messagesSend, 0)
}

func (suite *AgentTestSuite) TestMessageAckFailed() {

	queue := "test-queue"
	receiptHandle := "ack-error"
	suite.init(queue, nil)
	suite.mock.messages[queue][1] = messageForTest(nil, &receiptHandle, nil)

	go suite.agent.Run(suite.ctx, suite.waitGroupForTest())
	time.Sleep(500 * time.Millisecond)
	suite.cancelFunc()

	suite.True(suite.mock.flushed)
	suite.True(suite.mock.closed)
	suite.Len(suite.mock.messagesSend, 1)
	suite.Len(suite.mock.acks, 1)
}

func (suite *AgentTestSuite) waitGroupForTest() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	return wg
}
