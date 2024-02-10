package data

import (
	"context"
	"customer/internal/biz"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type CustomData struct {
	data *Data
}


func NewCustomData (data *Data) *CustomData {
	return &CustomData{
		data: data,
	}
}

func (d *CustomData) SetVerifyCode(telephone string, code string) error {
	const LIFE = 60
	redisSetResult := d.data.redisClient.Set(context.Background(), "CVC:" + telephone, code, LIFE * time.Second)

	if _, err := redisSetResult.Result(); err != nil {
		log.Println(err, "error")
		return err
	}
	

	return nil
}


func (d *CustomData) GetVerifyCode (telephone string) (string, error) {
	 str, err := d.data.redisClient.Get(context.Background(), "CVC:" + telephone).Result()
	 return str, err
}

func (d *CustomData) GetCustomerByTelephone(telephone string) (*biz.Customer, error) {
	customer := &biz.Customer{}

	result := d.data.mysqlClient.Where("telephone=?", telephone).First(customer)

	if result.Error == nil && customer.ID > 0 {
		return customer, nil
	}

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		customer.Telephone = telephone
		customer.Name = sql.NullString{Valid: false}
		customer.Wechat = sql.NullString{Valid: false}
		customer.Email = sql.NullString{Valid: false}

		result = d.data.mysqlClient.Create(customer)

		if result.Error != nil {
			return nil, result.Error
		} else {
			return customer, nil
		}
	}

	return nil, result.Error
}

func (d *CustomData) GenerateTokenAndSave(customer *biz.Customer, duration time.Duration, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "Customer",
		Subject: "customer-authentication",
		Audience: []string{"customer"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		NotBefore: nil,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID: strconv.Itoa(int(customer.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	
	customer.Token = signedToken
	customer.TokenCreated = sql.NullTime{
		Valid: true,
		Time: time.Now(),
	}
	d.data.mysqlClient.Save(customer)
	return signedToken, nil
}