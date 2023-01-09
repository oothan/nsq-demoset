package nsq

import (
	gonsq "github.com/nsqio/go-nsq"
	logger "nsq-demoset/app/_applib"
	"time"
)

func NewNsqConsumer(addr, topic, channel string, handel func(message *gonsq.Message) error, concurrency int) error {
	config := gonsq.NewConfig()
	config.LookupdPollInterval = 1 * time.Second

	consumer, err := gonsq.NewConsumer(topic, channel, config)
	if err != nil {
		logger.Sugar.Errorf("NSQ Consumer err: %v", err)
		panic(err)
	}
	consumer.AddHandler(gonsq.HandlerFunc(handel))
	//consumer.AddConcurrentHandlers(gonsq.HandlerFunc(handel), concurrency)

	err = consumer.ConnectToNSQD(addr)
	if err != nil {
		logger.Sugar.Errorf("NSQD Connect err: %v", err)
		panic(err)
	}

	return nil
}
