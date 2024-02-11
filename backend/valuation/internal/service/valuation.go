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
	duration, distance, err := biz.GetMapInfo(ctx, req.Origin, req.Destination)
	if  err != nil {
		log.Println(err)
	}
	
	log.Println(duration, distance)
	return &pb.GetEstimatePriceReply{}, nil
}
