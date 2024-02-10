package biz

import (
	"database/sql"

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