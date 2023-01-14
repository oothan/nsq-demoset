package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	logger "nsq-demoset/app/_applib"
	"strings"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func main() {
	kafkaReader := getKafkaReader("kafka:9092", "topic1", "logger-group")
	defer kafkaReader.Close()

	logger.Sugar.Info("start consuming ... ")
	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			logger.Sugar.Fatal(err)
		}
		logger.Sugar.Debugf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
