package nsq

func NsqTestMessageEvent(message string) {
	testDataToNsq(&TestNsqEventData{
		Message: message,
	})
}
