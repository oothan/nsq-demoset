package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	logger "nsq-demoset/app/_applib"
	libkafka "nsq-demoset/app/_applib/kafka"
	"nsq-demoset/app/kafka_services/internal/services"
)

func NewKafkaEventConsumer() {
	NewKafkaConsumer(libkafka.Test_Kafka_Topic, handleTestMessageKafkaEvent)
}

func NewKafkaConsumer(topic string, handle func(message *kafka.Message) error) error {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "goapp_consumer",
		"group.id":          "goapp_group",
	}

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	topics := []string{topic}
	err = consumer.SubscribeTopics(topics, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			logger.Sugar.Error(err)
			return err
		}

		go handle(msg)
	}
}

func handleTestMessageKafkaEvent(message *kafka.Message) error {
	defer func() {}()
	logger.Sugar.Debug(libkafka.Test_Kafka_Topic)

	msg := &libkafka.TestKafkaData{}
	if err := json.Unmarshal(message.Value, msg); err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if msg.Message != "" {
		go func() {
			fmt.Println("send to telegram")

			msg := fmt.Sprintf("Kafka Message: %v", msg.Message)
			_, err := services.TeleBot.SendMessage(msg)
			if err != nil {
				logger.Sugar.Error(err)
				return
			}
			logger.Sugar.Debug("send to telegram success")
		}()
	}

	return nil
}
