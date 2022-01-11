package main

import (
	"context"
	"sync"

	sqs "github.com/tommzn/aws-sqs"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
)

func newAgent(conf config.Config, logger log.Logger) (*agent, error) {

	messagePublisher, err := newKafkaClient(conf)
	return &agent{
		source:          sqs.NewConsumer(conf),
		target:          messagePublisher,
		routes:          routesFromConfig(conf),
		stopChan:        make(chan bool, 1),
		logger:          logger,
		metricPublisher: metrics.NewTimestreamPublisher(conf, logger),
	}, err
}

func routesFromConfig(conf config.Config) []route {

	routes := []route{}
	routesMap := conf.GetAsSliceOfMaps("agent.routes")
	for _, routeMap := range routesMap {
		if source, ok := routeMap["source"]; ok {
			newRoute := route{source: source}
			if target, ok := routeMap["target"]; ok {
				newRoute.target = target
			} else {
				newRoute.target = source
			}
			routes = append(routes, newRoute)
		}
	}
	return routes
}

func (agt *agent) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	defer agt.target.close()
	defer agt.logger.Flush()

	wg := &sync.WaitGroup{}
	waitCh := make(chan struct{})

	agt.logger.Infof("Run message forwarding for %d queues.", len(agt.routes))
	go func() {
		for _, route := range agt.routes {
			wg.Add(1)
			go agt.forward(route.source, route.target, wg)
		}
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
		agt.stop()
		<-waitCh
	case <-waitCh:
		agt.logger.Info("Message publishing finished.")
	}
	agt.target.flush()

	agt.metricPublisher.Flush()
	if err := agt.metricPublisher.Error(); err != nil {
		agt.logger.Error(err)
	}
	return nil
}

func (agt *agent) forward(sourceQueue, targetTopic string, wg *sync.WaitGroup) {

	defer wg.Done()
	agt.logger.Infof("Start message forwarding from qeuue %s to topic %s.", sourceQueue, targetTopic)

	messages, err := agt.source.Receive(sourceQueue)
	if err != nil {
		agt.logger.Error("Error receiving new messages: ", err)
		return
	}

	if len(messages) == 0 {
		agt.logger.Info("No new new messages at queue ", sourceQueue)
		return
	}

	agt.logger.Debugf("Process %d messages from queue %s", len(messages), sourceQueue)
	for _, message := range messages {

		if agt.shouldStop() {
			agt.logger.Info("Stop message delivery as requested.")
			return
		}

		err = agt.target.send(targetTopic, []byte(*message.Body))
		if err != nil {
			agt.logger.Errorf("Unable to publish message to topic %s, reason: %s", targetTopic, err)
			return
		}
		err = agt.source.Ack(sourceQueue, message.ReceiptHandle)
		if err != nil {
			agt.logger.Errorf("Unable to ack message processing on queue %s, reason: %s", sourceQueue, err)
			return
		}
		agt.metricPublisher.Send(createMeasurement(sourceQueue, targetTopic))
	}
}

func (agt *agent) stop() {
	agt.stopChan <- true
}

func (agt *agent) shouldStop() bool {
	return len(agt.stopChan) != 0
}

func createMeasurement(sourceQueue, targetTopic string) metrics.Measurement {
	return metrics.Measurement{
		MetricName: "hdb-message-agent",
		Tags: []metrics.MeasurementTag{
			metrics.MeasurementTag{
				Name:  "source_queue",
				Value: sourceQueue,
			},
			metrics.MeasurementTag{
				Name:  "target_queue",
				Value: targetTopic,
			},
		},
		Values: []metrics.MeasurementValue{
			metrics.MeasurementValue{
				Name:  "count",
				Value: 1,
			},
		},
	}
}
