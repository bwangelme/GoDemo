package main

import (
	"bwdemo/kafka_thrift/gen-go/rpclog"
	"bwdemo/kafka_thrift/logger"
	"context"
	"github.com/apache/thrift/lib/go/thrift"
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
			if err != nil {
				if !err.(kafka.Error).IsTimeout() {
					// The client will automatically try to recover from all errors.
					// Timeout is not considered an error because it is raised by
					// ReadMessage in absence of messages.
					logger.L.Warningf("Consumer error: %v (%v)\n", err, msg)
				}
				continue
			}

			de := thrift.NewTDeserializer()
			entry := rpclog.NewLogEntry()
			err = de.Read(context.Background(), entry, msg.Value)
			if err != nil {
				logger.L.Warningf("deserializer log entry failed: %v", err)
				continue
			}
			logger.L.Infof("logentry %s:%s", entry.Date, entry.Msg)
		}
	}

	logger.L.Info("Close Consumer")
	c.Close()
}
