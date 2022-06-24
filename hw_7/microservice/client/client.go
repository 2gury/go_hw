package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"lectures-2022-1/08_microservices/99_hw/microservice/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func firstClient() {
	conn, err := grpc.Dial("127.0.0.1:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ctx := context.Background()
	md := metadata.Pairs(
		"consumer", "first",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	biz := service.NewBizClient(conn)
	biz.Check(ctx, &service.Nothing{})
}

func secondClient() {
	conn, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer conn.Close()

	admin := service.NewAdminClient(conn)

	ctx := context.Background()
	md := metadata.Pairs(
		"consumer", "second",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, _ := admin.Statistics(ctx, &service.StatInterval{IntervalSeconds: 2})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happed", err)
				return
			}
			fmt.Println("second: " + event.String())
		}
	}(wg)

	stream.CloseSend()
	wg.Wait()
}

func thirdClient() {
	conn, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer conn.Close()

	admin := service.NewAdminClient(conn)

	ctx := context.Background()
	md := metadata.Pairs(
		"consumer", "third",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, _ := admin.Statistics(ctx, &service.StatInterval{IntervalSeconds: 3})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happed", err)
				return
			}
			fmt.Println("third: " + event.String())
		}
	}(wg)

	stream.CloseSend()
	wg.Wait()
}

func main() {
	// go secondClient()
	thirdClient()
}
