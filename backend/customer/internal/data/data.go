package data

import (
	"customer/internal/biz"
	"customer/internal/conf"
	"fmt"
	logg "log"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCustomData)

// Data .
type Data struct {
	redisClient *redis.Client
	mysqlClient *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}

	// 连接 redis
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s/1?dial_timeout=2", c.Redis.User, c.Redis.Password, c.Redis.Addr))
	if err != nil {
		data.redisClient = nil
		log.Error("连接 redis 失败")
	}
	data.redisClient = redis.NewClient(opt)

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
		data.redisClient.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}

func migrateTable(d *Data) {
	if err := d.mysqlClient.AutoMigrate(biz.Customer{}); err != nil {
		logg.Println(err)
	}
}

func (d *CustomData) GetToken(id interface{}) (string, error) {
	c := biz.Customer{}
	result := d.data.mysqlClient.First(&c, id)
	if result.Error != nil {
		return "", result.Error
	}
	return c.Token, nil
}