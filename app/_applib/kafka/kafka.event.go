package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestKafkaMessageEvent(message string) {
	deliveryChan := make(chan kafka.Event)
	go DeliveryReport(deliveryChan)

	Publish(TestKafkaData{Message: message}, Test_Kafka_Topic, Producer, nil, deliveryChan)
}
