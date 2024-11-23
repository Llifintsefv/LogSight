package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Llifintsefv/LogSight/grpc-server/LogSight"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051",grpc.WithInsecure())
    if err != nil {
        log.Fatalf("fail to create grpc client %s\n",err)
    }
    defer conn.Close()

    client := pb.NewLogCollectorClient(conn)
    req := &pb.LogRequest{Service: "test", Level: "error", Message: "test", Timestamp: "test", Ip: "test"}
    req2 := &pb.LogRequest{Service: "test2", Level: "error2", Message: "test2", Timestamp: "test2", Ip: "test3"}
    req3 := &pb.LogRequest{Service: "test3", Level: "error3", Message: "test3", Timestamp: "test3", Ip: "test3"}
    req4 := &pb.LogRequest{Service: "test4", Level: "error4", Message: "test4", Timestamp: "test4", Ip: "test4"}
    req5 := &pb.LogRequest{Service: "test5", Level: "error5", Message: "test5", Timestamp: "test5", Ip: "test5"}
    req6 := &pb.LogRequest{Service: "test6", Level: "error6", Message: "test6", Timestamp: "test6", Ip: "test6"}
    request := []*pb.LogRequest{req,req2,req3,req4,req5,req6}
    for {
        for _, req := range request {
            res, err := client.SendLog(context.Background(),req)
            if err != nil {
                log.Fatalf("failed to call grpc server %s\n",err)
            }
            log.Printf("Response from grpc server: %s\n",res.Message)
        }
        time.Sleep(1 * time.Second)
    }
}