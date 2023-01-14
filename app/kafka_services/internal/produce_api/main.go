package main

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"log"
	"net/http"
	logger "nsq-demoset/app/_applib"
)

func producerHandler(kafkaWriter *kafka.Writer) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln(err)
		}
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
			Value: body,
		}
		err = kafkaWriter.WriteMessages(req.Context(), msg)

		if err != nil {
			wrt.Write([]byte(err.Error()))
			log.Fatalln(err)
		}
	})
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	kafkaWriter := getKafkaWriter("kafka:9092", "topic1")
	defer kafkaWriter.Close()

	http.HandleFunc("/", producerHandler(kafkaWriter))

	logger.Sugar.Info("start producer-api ... ")
	logger.Sugar.Fatal(http.ListenAndServe(":8080", nil))
}
