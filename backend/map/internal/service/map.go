package service

import (
	"context"
	"log"

	pb "map/api/map"
	"map/internal/biz"
)

type MapService struct {
	pb.UnimplementedMapServer
	mapBiz *biz.MapServiceBiz
}

func NewMapService(mapBiz *biz.MapServiceBiz) *MapService {
	return &MapService{
		mapBiz: mapBiz,
	}
}

func (s *MapService) GetDriveInfo(ctx context.Context, req *pb.GetDriveInfoRequest) (*pb.GetDriveInfoReply, error) {
	distance, duration, err := s.mapBiz.GetDriveInfo(req.Origin, req.Destination)

	if err != nil {
		log.Println(err)
	}

	return &pb.GetDriveInfoReply{
		Distance: distance,
		Duration: duration,
		Origin: req.Origin,
		Destination: req.Destination,
	}, nil
}
