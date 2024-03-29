package data

import (
	"log"
	"valuation/internal/biz"
)

type ValuationRepo struct {
	data *Data
}

func NewValuationRepo(data *Data) biz.ValuationRepo {
	return &ValuationRepo{data: data}
}

func (r *ValuationRepo) GetRule(cityID int, currentTime int) (*biz.PriceRule, error) {
	pr := biz.PriceRule{}
	result := r.data.mysqlClient.Where("city_id=? AND start_at >= ?", cityID, currentTime).First(&pr)
	log.Println(pr)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pr, nil
}