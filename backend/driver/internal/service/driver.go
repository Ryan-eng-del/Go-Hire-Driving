package service

import (
	"context"
	pb "driver/api/driver"
	"driver/internal/biz"
	"time"
)

type DriverService struct {
	pb.UnimplementedDriverServer
	dbz *biz.DriverBiz
}

func NewDriverService(dbz *biz.DriverBiz) *DriverService {
	return &DriverService{
		dbz: dbz,
	}
}

func (s *DriverService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeResponse, error) {
	verifyCode, err := s.dbz.GetVerifyCode(ctx, req.Telephone)
	if err != nil {
		return &pb.GetVerifyCodeResponse{
			Code: 500,
			VerifyCode: verifyCode,
			VerifyCodeTime: int32(time.Now().Unix()),
			VerifyCodeLife: 60,
			Message: "获取验证码失败",
		}, nil
	}

	return &pb.GetVerifyCodeResponse{
		Code: 200,
		VerifyCode: verifyCode,
		VerifyCodeTime: int32(time.Now().Unix()),
		VerifyCodeLife: 60,
		Message: "获取验证码成功",
	}, nil
}
