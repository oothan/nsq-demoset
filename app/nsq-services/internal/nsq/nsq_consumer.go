package nsq

import (
	"encoding/json"
	"fmt"
	gonsq "github.com/nsqio/go-nsq"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/_applib/nsq"
	"nsq-demoset/app/conf"
	"nsq-demoset/app/nsq-services/internal/services"
)

func NewNsqTestMessageEventConsumer() {
	nsq.NewNsqConsumer(conf.Nsq().Addr, nsq.TestEventNsqTopic, "1", handelNoteDynamicEvent, 1)
}

func handelNoteDynamicEvent(message *gonsq.Message) error {
	defer func() {}()
	logger.Sugar.Info("testNsqEventData", string(message.Body))

	var testMsg nsq.TestNsqEventData
	if err := json.Unmarshal(message.Body, &testMsg); err != nil {
		logger.Sugar.Error(err)
		return err
	}

	if testMsg.Message != "" {

		go func() {
			fmt.Println("send to telegram")

			msg := fmt.Sprintf("NSQ Message: %v", testMsg.Message)
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
