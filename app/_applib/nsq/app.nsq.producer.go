package nsq

import (
	"encoding/json"
	gonsq "github.com/nsqio/go-nsq"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/conf"
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
	//cfg.ReadTimeout = 60 * time.Second
	//cfg.DialTimeout = 5 * time.Second

	producer, err = gonsq.NewProducer(conf.NsqAddr, cfg)
	if err != nil {
		logger.Sugar.Errorf("nsq producer err: %v", err)
		panic(err)
	}

	err = producer.Ping()
	if err != nil {
		logger.Sugar.Errorf("NSQ Ping err: %v", err)
		panic(err)
	}

	logger.Sugar.Info("NSQD started ... ")
}

func testDataToNsq(data *TestNsqEventData) {
	logger.Sugar.Info("testNsqEventData")
	body, err := json.Marshal(data)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	err = producer.Publish(TestEventNsqTopic, body)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}
