package main

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type KafkaTestSuite struct {
	suite.Suite
}

func TestKafkaTestSuite(t *testing.T) {
	suite.Run(t, new(KafkaTestSuite))
}

func (suite *KafkaTestSuite) TestPublishMessages() {

	conf := loadConfigForTest(nil)
	client, err := newKafkaClient(conf, loggerForTest())
	suite.Nil(err)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go client.(*kafkaClient).forwardLogs(ctx)
	suite.Nil(client.send("test-topic", "xxx"))
	client.flush()
	client.close()

	time.Sleep(3 * time.Second)
	cancelFunc()
	time.Sleep(1 * time.Second)
}
