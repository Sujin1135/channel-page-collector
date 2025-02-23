package main

import (
	grpc2 "channel-page-collector/api/grpc"
	pb "channel-page-collector/api/grpc/interface/channel-page-collector-interface/protobuf"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
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
	channelService := grpc2.NewChannelService()
	s := grpc.NewServer()
	pb.RegisterChannelPageServer(s, channelService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
