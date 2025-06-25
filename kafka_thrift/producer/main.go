package main

import (
	"context"
	"time"

	"bwdemo/kafka_thrift/gen-go/rpclog"
	"bwdemo/kafka_thrift/logger"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.L.Warningf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					logger.L.Infof("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "lazycat"
	entry := rpclog.LogEntry{
		Msg:  "hello, the server start",
		Date: time.Now().Format("2006-01-02 15:04:05 -0700"),
	}
	se := thrift.NewTSerializer()
	val, err := se.Write(context.Background(), &entry)
	if err != nil {
		logger.L.Warningf("Serialize failed: %v", err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          val,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
