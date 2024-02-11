package biz

import (
	"context"
	"customer/api/valuation"
	"database/sql"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)


type Customer struct {
	CustomerWork
	gorm.Model
	CustomerToken
}

const Secret = "my-secret"


type CustomerWork struct {
	Telephone string `gorm:"type:varchar(15);uniqueIndex;" json:"telephone"`
	Name sql.NullString `gorm:"type:varchar(15);uniqueIndex;" json:"name"`
	Email sql.NullString  `gorm:"type:varchar(15);uniqueIndex;" json:"email"`
	Wechat sql.NullString  `gorm:"type:varchar(15);uniqueIndex;" json:"wechat"`
}	

type CustomerToken struct {
	Token string `gorm:"type:varchar(4095);" json:"token"`
	TokenCreated sql.NullTime `json:"token_created"`
}


type CustomerBiz struct {
}

func NewCustomerBiz() *CustomerBiz {
	return &CustomerBiz{}
}

func (b *CustomerBiz) GetEstimatePrice(ctx context.Context, origin, destination string) (int64, error) {
	// 服务间通信 grpc
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return 0, err
	}
	consulDiscovery := consul.New(consulClient)
	
	// 负载均衡 策略
	selector.SetGlobalSelector(wrr.NewBuilder())
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///Valuation"),
		grpc.WithDiscovery(consulDiscovery),
		grpc.WithMiddleware(tracing.Client()),
	)	

	if err != nil {
		return 0, err
	}
	
	valuationClient := valuation.NewValuationClient(conn)
	reply, err := valuationClient.GetEstimatePrice(ctx, &valuation.GetEstimatePriceRequest{
		Origin: origin,
		Destination: destination,
	})

	if err != nil {
		return 0, err
	}
	return reply.Price, nil
}