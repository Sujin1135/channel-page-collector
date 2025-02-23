package main

import (
	grpc2 "channel-page-collector/api/grpc"
	pb "channel-page-collector/api/grpc/interface/channel-page-collector-interface/protobuf"
	"channel-page-collector/internal/collector"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChannelPageServer(s, &grpc2.ChannelService{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	websiteCollector := collector.NewWebsiteCollector()

	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	subscriberNamesChan := make(chan []string, 1)

	go func() {
		defer wg.Done()
		fmt.Println("start to collect method 1")
		err := websiteCollector.CollectWithScrolling("운동", 5, subscriberNamesChan)
		if err != nil {
			log.Fatal("failed to collect subscriber names", err)
		}
	}()

	fmt.Println("channel data as follow:")
	for v := range subscriberNamesChan {
		fmt.Println(v)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)
}
