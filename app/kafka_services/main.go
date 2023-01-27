package main

import (
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/conf"
	"nsq-demoset/app/kafka_services/internal/kafka"
	"nsq-demoset/app/kafka_services/internal/services"
)

func main() {
	conf.InitYaml()
	services.NewTelegramBot(conf.Telegram().TokenID, conf.Telegram().GroupID)

	//libkafka.InitKafkaProducer()
	kafka.NewKafkaEventConsumer()

	logger.Sugar.Info("Kafka Consumer started ... ")
	select {}
}
