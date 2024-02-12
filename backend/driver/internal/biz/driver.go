package biz

import "context"


type DriverBiz struct {
	driverData DriverImpl
}

func NewDriverBiz(d DriverImpl) *DriverBiz {
	return &DriverBiz{
		driverData: d,
	}
}

type DriverImpl interface {
	SetVerifyCode(ctx context.Context, telephone string) (error)
}

func (dbz *DriverBiz) GetVerifyCode(ctx context.Context,telephone string) (string, error) {
	dbz.driverData.SetVerifyCode(ctx, telephone)
	return "", nil
}