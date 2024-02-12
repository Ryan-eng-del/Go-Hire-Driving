package data

import (
	"context"
	"driver/api/verifyCode"
	"driver/internal/biz"
	"driver/internal/conf"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
)

type DriverRepo struct {
	d *Data
	conf *conf.Service
}


func NewDriverRepo(d *Data,conf *conf.Service) biz.DriverImpl {
	return &DriverRepo{
		d: d,
		conf: conf,
	}
}


func (d *DriverRepo) GetVerifyCode(ctx context.Context,telephone string) (code string, err error) {
		// 服务间通信 grpc
		consulConfig := api.DefaultConfig()
		consulConfig.Address = d.conf.Consul.Addr
		consulClient, err := api.NewClient(consulConfig)
		consulDiscovery := consul.New(consulClient)
	
		// 负载均衡 策略
		selector.SetGlobalSelector(wrr.NewBuilder())
	
		if err != nil {
			return
		}
	
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("discovery:///VerifyCode"),
			grpc.WithDiscovery(consulDiscovery),
			grpc.WithMiddleware(tracing.Client()),
		)	
	
		if err != nil  {
			return
		}
		
		defer conn.Close()

		client := verifyCode.NewVerifyCodeClient(conn)
	
		reply, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
			Length: 8,
			Type: 1,
		})
	
		if err != nil {
			return 
		}

		const LIFE = 60
		redisSetResult := d.d.redisClient.Set(context.Background(), "DVC:" + telephone, code, LIFE * time.Second)
	
		if _, err = redisSetResult.Result(); err != nil {
			log.Println(err, "error")
			return
		}
	return reply.Code, nil
}