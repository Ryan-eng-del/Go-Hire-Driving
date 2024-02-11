package data

import (
	"valuation/internal/biz"
	"valuation/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewValuationRepo)

// Data .
type Data struct {
	mysqlClient *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}


	// 连接 mysql
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})

	if err != nil {
		data.mysqlClient = nil
		log.Error("连接 mysql 失败")
	}
	data.mysqlClient = db
	// 开发模式，每次启动记得 migrate database
	migrateTable(data)
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}

func migrateTable(d *Data) {
	if err := d.mysqlClient.AutoMigrate(biz.PriceRule{}); err != nil {
		log.Info(err)
	}
}