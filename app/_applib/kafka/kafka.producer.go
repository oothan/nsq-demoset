package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	logger "nsq-demoset/app/_applib"
)

var Producer *kafka.Producer

func InitKafkaProducer() {
	logger.Sugar.Info("Starting Kafka Producer ... ")
	if Producer == nil {
		initKafkaProducer()
	}
}

func initKafkaProducer() {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "localhost:9092",
		"delivery.timeout.ms": "1",
		"acks":                "all",
		"enable.idempotence":  "true",
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	Producer = p
}

func Publish(msg any, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	message := &kafka.Message{
		Value: bytes,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}
	err = producer.Produce(message, deliveryChan)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch e.(type) {
		case *kafka.Message:
			e := <-deliveryChan
			msg := e.(*kafka.Message)

			if msg.TopicPartition.Error != nil {
				logger.Sugar.Error("error in topic partition.")
			} else {
				logger.Sugar.Debug("Message event : ", msg.TopicPartition)
				// @TODO something what you want over msg
			}
		}
	}
}
