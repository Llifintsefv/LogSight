package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Llifintsefv/LogSight/grpc-server/LogSight"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedLogCollectorServer
	kafkaProducer *kafka.Producer
}

var kafkaTopic string

func (s *Server) SendLog(ctx context.Context,req *pb.LogRequest) (*pb.LogResponse,error) {
	log.Printf("Received log: Service=%v, Level=%s, Message=%s, Timestamp=%s, IP=%s", req.Service, req.Level,req.Message,req.Timestamp,req.Ip)
	message :=  fmt.Sprintf(`{"service=%v, level=%s, message=%s, timestamp=%s, ip=%s"}`, req.Service, req.Level,req.Message,req.Timestamp,req.Ip)
	deliveryChan := make(chan kafka.Event)

	err := s.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		log.Printf("failed to produce message to kafka: %v", err)
		return &pb.LogResponse{Success: false, Message: "failed to produce to kafka"}, err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return &pb.LogResponse{Success: false, Message: "failed to deliver to kafka"}, m.TopicPartition.Error
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
	return &pb.LogResponse{Success: true, Message: "Log received"}, nil
}

func main() {
	p,err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("failed to create producer %s\n",err)
	}
	defer p.Close()
	kafkaTopic = "logs"
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			case kafka.Error:
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	lis, err := net.Listen("tcp",":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterLogCollectorServer(s,&Server{kafkaProducer: p})
	log.Printf("Server listening %v",lis.Addr())
	if err := s.Serve(lis);err != nil {
		log.Fatal(err)
	}
}