package biz

import (
	"context"
	mapService "valuation/api/map"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

// GreeterUsecase is a Greeter usecase.
type ValuationBiz struct {
	log  *log.Helper
	repo ValuationRepo
}

type PriceRule struct {
	gorm.Model
	CityID uint `json:"city_id"`
	StartFee int64 `json:"start_fee"`
	DistanceFee int64 `json:"distance_fee"`
	DurationFee int64 `json:"duration_fee"`
	StartAt int `json:"start_at"`
	EndedAt int `json:"ended_at"`
}

type ValuationRepo interface {
	GetRule(cityID int, currentTime int) (*PriceRule, error)
}

// NewGreeterUsecase new a Greeter usecase.
func NewValuationBiz(repo ValuationRepo, logger log.Logger) *ValuationBiz {
	return &ValuationBiz{log: log.NewHelper(logger), repo: repo}
}


func GetMapInfo(ctx context.Context, origin string, destination string) (duration, distance string, err error) {
		// 服务间通信 grpc
		consulConfig := api.DefaultConfig()
		consulConfig.Address = "localhost:8500"
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			return
		}
		consulDiscovery := consul.New(consulClient)
		// 负载均衡 策略
		selector.SetGlobalSelector(wrr.NewBuilder())
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("discovery:///Map"),
			grpc.WithDiscovery(consulDiscovery),
		)	

		if err != nil {
			return
		}
		mapClient := mapService.NewMapClient(conn)
		resp, err := mapClient.GetDriveInfo(ctx, &mapService.GetDriveInfoRequest{
			Origin: origin,
			Destination: destination,
		})

		if err != nil {
			return
		}

		duration = resp.Duration
		distance = resp.Distance
		return
}

