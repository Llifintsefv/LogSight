package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	consumer,err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id": "Log-group",
		"auto.offset.reset":"earliest",
	})
	if err != nil {
		log.Fatalf("error to create consumer %s/n",err)
	}
	defer consumer.Close()

	kafkaTopic := "logs"
	consumer.SubscribeTopics([]string{kafkaTopic},nil)

	for {
		msg,err := consumer.ReadMessage(-1)
		if err == nil {
			log.Println(msg)
		} else {
			log.Println("error received %s\n",err)
		}
	}
}