package grpc

import (
	pb "channel-page-collector/api/grpc/interface/channel-page-collector-interface/protobuf"
	"github.com/pkg/errors"
	"log"
)

type ChannelService struct {
	pb.UnimplementedChannelPageServer
}

func (*ChannelService) GetSubscriberNames(in *pb.GetSubscriberNamesRequest, stream pb.ChannelPage_GetSubscriberNamesServer) error {
	log.Printf("Received: %v %v\n", in.Keyword, in.Page)

	err := stream.Send(&pb.GetSubscriberNamesResponse{Names: []string{in.Keyword}})
	if err != nil {
		return errors.Wrap(err, "occurred an error when sending a stream message\n")
	}

	return nil
}
