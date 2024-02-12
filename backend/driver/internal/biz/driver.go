package biz

import (
	"context"
	"regexp"
)


type DriverBiz struct {
	driverData DriverImpl
}

func NewDriverBiz(d DriverImpl) *DriverBiz {
	return &DriverBiz{
		driverData: d,
	}
}

type DriverImpl interface {
	GetVerifyCode(ctx context.Context, telephone string) (string, error)
}

func (dbz *DriverBiz) GetVerifyCode(ctx context.Context,telephone string) (code string, err error) {
	pattern := `^1[34578]\d{9}$`

	regexPattern := regexp.MustCompile(pattern)

	if !regexPattern.MatchString(telephone) {
		return
	}
	
	code, err = dbz.driverData.GetVerifyCode(ctx, telephone)
	if err != nil {
		return
	}
	return
}