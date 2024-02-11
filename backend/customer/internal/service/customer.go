package service

import (
	"context"
	pb "customer/api/customer"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"customer/internal/data"
	"log"
	"regexp"
	"time"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/consul/api"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	CustomData *data.CustomData
}

func NewCustomerService(d *data.CustomData) *CustomerService {
	return &CustomerService{
		CustomData: d,
	}
}

func (s *CustomerService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequst) (*pb.GetVerifyCodeResponse, error) {
	pattern := `^1[34578]\d{9}$`

	regexPattern := regexp.MustCompile(pattern)
	if !regexPattern.MatchString(req.Telephone) {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			Message: "电话号码格式错误",
		}, nil
	}
	
	// 服务间通信 grpc
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "localhost:8500"
	consulClient, err := api.NewClient(consulConfig)
	consulDiscovery := consul.New(consulClient)

	// 负载均衡 策略
	selector.SetGlobalSelector(wrr.NewBuilder())
	
	if err != nil {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			Message: "连接服务发现失败",
		}, nil
	}

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///VerifyCode"),
		grpc.WithDiscovery(consulDiscovery),
	)	

	if err != nil  {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			Message: "连接验证码服务失败",
		}, nil
	}
	
	defer conn.Close()

	client := verifyCode.NewVerifyCodeClient(conn)

	reply, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
		Length: 8,
		Type: 1,
	})

	if err != nil {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			Message: "获取验证码失败",
		}, nil
	}

	if err := s.CustomData.SetVerifyCode(req.Telephone, reply.Code); err != nil {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			Message: "Redis Set 操作失败",
		}, nil
	}

	return &pb.GetVerifyCodeResponse{
		VerifyCode: reply.Code,
		VerifyCodeTime: int32(time.Now().Unix()),
		Message: "获取验证码成功",
		VerifyCodeLife: 60,
	}, nil
}


func (s *CustomerService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	code, err := s.CustomData.GetVerifyCode(req.Telephone)
	if code == "" || code != req.VerifyCode || err != nil {
		log.Println(err)
		return &pb.LoginResponse{
			Code: 500,
			Message: "验证码不匹配",
			}, nil
	}

	customer, err := s.CustomData.GetCustomerByTelephone(req.Telephone)

	if err != nil {
		log.Println(err)
		return &pb.LoginResponse{
			Code: 500,
			Message: "顾客信息获取错误",
			}, nil
	}

	const durations = 24 * time.Hour
	
	token, err := s.CustomData.GenerateTokenAndSave(customer, durations, biz.Secret)

	if token == "" || err != nil {
		log.Println(err)
		return &pb.LoginResponse{
			Code: 500,
			Message: "生成Token失败",
			}, nil
	}

	return &pb.LoginResponse{
		Code: 500,
		Message: "获取 Token 成功",
		Token: token,
		TokenCreateAt: time.Now().Unix(),
		TokenLife: 24 * 3600,
	}, nil
}

func (s *CustomerService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	claims, isExists := jwt.FromContext(ctx)
	claimsMap := claims.(jwtv4.MapClaims)
	userId := claimsMap["user_id"]


	if !isExists {
		return &pb.LogoutResponse{
			Code: 500,
			Message: "Token 未传递",
		}, nil
	}

	err := s.CustomData.DelToken(userId)
	
	if err != nil {
		log.Println(err)
		return &pb.LogoutResponse{
			Code: 500,
			Message: "Token 删除失败",
		}, nil
	}

	return &pb.LogoutResponse{
		Code: 200,
		Message: "退出登录成功",
	}, nil
}
