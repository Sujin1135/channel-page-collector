package grpc

import (
	pb "channel-page-collector/api/grpc/interface/channel-page-collector-interface/protobuf"
	"channel-page-collector/internal/collector"
	"log"
	"sync"
)

type ChannelService struct {
	pb.UnimplementedChannelPageServer
	collector *collector.Collector
}

func NewChannelService() *ChannelService {
	return &ChannelService{
		collector: collector.NewWebsiteCollector(),
	}
}

func (c *ChannelService) GetSubscriberNames(in *pb.GetSubscriberNamesRequest, stream pb.ChannelPage_GetSubscriberNamesServer) error {
	var wg sync.WaitGroup
	wg.Add(2)

	subscriberNamesChan := make(chan []string, 1)

	go func() {
		defer wg.Done()
		err := c.collector.CollectWithScrolling(in.Keyword, int(in.Page), subscriberNamesChan)
		if err != nil {
			log.Printf("failed to collect subscriber names: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		for subscriberNames := range subscriberNamesChan {
			err := stream.Send(&pb.GetSubscriberNamesResponse{Names: subscriberNames})
			if err != nil {
				log.Printf("failed to send subscriber names to grpc client: %v\n", err)
			}
		}
	}()

	wg.Wait()

	log.Println("Ended to collect youtube handles")

	return nil
}
