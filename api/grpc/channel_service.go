package grpc

import (
	pb "channel-page-collector/api/grpc/interface/channel-page-collector-interface/protobuf"
	"channel-page-collector/internal/collector"
	"github.com/pkg/errors"
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
	wg.Add(1)

	subscriberNamesChan := make(chan []string, 1)

	go func() {
		defer wg.Done()
		err := c.collector.CollectWithScrolling(in.Keyword, int(in.Page), subscriberNamesChan)
		if err != nil {
			log.Fatal("failed to collect subscriber names", err)
		}
	}()

	for subscriberNames := range subscriberNamesChan {
		err := stream.Send(&pb.GetSubscriberNamesResponse{Names: subscriberNames})
		if err != nil {
			return errors.Wrap(err, "failed to send subscriber names to grpc client")
		}
	}

	return nil
}
