package service

import (
	"context"
	"log"
	pb "valuation/api/valuation"
	"valuation/internal/biz"
)

type ValuationService struct {
	pb.UnimplementedValuationServer
	valuationBiz *biz.ValuationBiz
}

func NewValuationService(valuationBiz *biz.ValuationBiz) *ValuationService {
	return &ValuationService{
		valuationBiz: valuationBiz,
	}
}

func (s *ValuationService) GetEstimatePrice(ctx context.Context, req *pb.GetEstimatePriceRequest) (*pb.GetEstimatePriceReply, error) {
	duration, distance, err := s.valuationBiz.GetMapInfo(ctx, req.Origin, req.Destination)
	if err != nil {
		log.Println(err)
	}
	total, err := s.valuationBiz.GetPrice(ctx, duration, distance, 1, 6)

	if err != nil {
		log.Println(err)
	}
	return &pb.GetEstimatePriceReply{
		Price: total,
		Origin: req.Origin,
		Destination: req.Destination,
	}, nil
}
