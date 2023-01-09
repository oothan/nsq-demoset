package main

import (
	"log"
	logger "nsq-demoset/app/_applib"
	libnsq "nsq-demoset/app/_applib/nsq"
	"nsq-demoset/app/nsq-services/conf"
	_ "nsq-demoset/app/nsq-services/conf"
	"nsq-demoset/app/nsq-services/internal/nsq"
	"nsq-demoset/app/nsq-services/internal/services"
)

func main() {
	logger.Sugar.Info("Initialized NSQ ... ")

	conf.InitYaml()
	services.NewTelegramBot(conf.Telegram().TokenID, conf.Telegram().GroupID)

	libnsq.InitNSQProducer()
	nsq.NewNsqTestMessageEventConsumer()

	log.Println("bkNsq started...")
	select {}
}
