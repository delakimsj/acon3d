package model

import (
	"encoding/json"

	"acon3d.com/framework"
)

type (
	DBExchangeRate struct {
		ExchangeRateID int `json:"exchange_rate_id"`
		// BaseDate     string `json:"base_date"`
		Currency string  `json:"currency"` // ja, zh
		Rate     float64 `json:"rate"`
	}

	DBMSExchangeRate struct {
		ExchangeRateIdSequence int              `json:"exchange_rate_id_sequence"`
		Data                   []DBExchangeRate `json:"data"`
	}
)

func (dbms *DBMSExchangeRate) LoadData() {
	if json.Unmarshal(*framework.LoadDBMS("exchange_rate"), &dbms) != nil {
		panic("Can't load exchange_rate data")
	}
}

func GetExchangeRate(appReq *framework.AppRequest, currency string) float64 {
	var dbms DBMSExchangeRate
	dbms.LoadData()

	for _, v := range dbms.Data {
		if v.Currency == currency {
			return v.Rate
		}
	}

	panic("No currency data")
}
