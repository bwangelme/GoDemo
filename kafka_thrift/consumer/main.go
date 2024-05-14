package main

import (
	"bwdemo/kafka_thrift/logger"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "thrift_consumer",
		"auto.offset.reset": "earliest",
	})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		logger.L.Panicf("new consumer err: %v", err)
	}

	err = c.SubscribeTopics([]string{"lazycat"}, nil)
	if err != nil {
		logger.L.Panicf("subscribe topic error: %v", err)
	}

	logger.L.Info("Start Consumer")
OuterLoop:
	for {
		select {
		case sig := <-sigChan:
			logger.L.WithFields(logrus.Fields{
				"signal": sig.String(),
			}).Infof("get signal")
			break OuterLoop
		default:
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}

	logger.L.Info("Close Consumer")
	c.Close()
}
