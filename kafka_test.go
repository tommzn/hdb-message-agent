package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type KafkaTestSuite struct {
	suite.Suite
}

func TestKafkaTestSuite(t *testing.T) {
	suite.Run(t, new(KafkaTestSuite))
}

func (suite *KafkaTestSuite) TestPublishMessages() {

	conf := loadConfigForTest(nil)
	client, err := newKafkaClient(conf)
	suite.Nil(err)

	suite.Nil(client.send("test-topic", "xxx"))
	client.flush()
	client.close()
}
