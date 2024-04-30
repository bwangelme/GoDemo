package main

import (
	"bwdemo/kafka_thrift/logger"
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})
}

func main() {
	// get kafka reader using environment variables.
	kafkaURL := "localhost:9092"
	topic := "lazycat"
	groupID := "thrift_consumer"

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	logger.L.Info("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		logger.L.WithFields(logrus.Fields{
			"topic":     m.Topic,
			"partition": m.Partition,
			"offset":    m.Offset,
			"value":     string(m.Value),
		}).Infof("message")
	}
}
