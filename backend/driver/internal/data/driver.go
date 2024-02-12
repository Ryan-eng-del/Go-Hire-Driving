package data

import (
	"context"
	"driver/internal/biz"
)

type DriverRepo struct {
	d *Data
}


func NewDriverRepo(d *Data) biz.DriverImpl {
	return &DriverRepo{
		d: d,
	}
}


func (d *DriverRepo) SetVerifyCode(ctx context.Context,telephone string) error {
	return  nil
}