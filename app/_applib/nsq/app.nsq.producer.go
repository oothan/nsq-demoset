package nsq

import (
	gonsq "github.com/nsqio/go-nsq"
	logger "nsq-demoset/app/_applib"
	"os"
)

var producer *gonsq.Producer

func InitNSQProducer() {
	logger.Sugar.Info("Starting NSQD ... ")
	if producer == nil {
		initNsqProducer()
	}
}

func initNsqProducer() {
	var err error
	cfg := gonsq.NewConfig()
	// cfg.ReadTimeout = 60 * time.Second
	// cfg.DialTimeout = 5 * time.Second

	producer, err = gonsq.NewProducer(os.Getenv("NSQ_ADDR"), cfg)
	if err != nil {
		logger.Sugar.Error("nsq new panic")
		panic("nsq new panic")
	}

	err = producer.Ping()
	if err != nil {
		logger.Sugar.Error(err.Error(), " NSQ Ping")
		panic(err.Error())
	}

	logger.Sugar.Info("Nsqd started ... ")
}
