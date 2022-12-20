package handler

import (
	"acon3d.com/framework"
	"acon3d.com/model"
)

func IsValidLanguageCode(inp string) bool {
	return inp == "ko" || inp == "zh" || inp == "ja"
}

func TranslateTo(lan string, statement string) string {
	if lan == "zh" {
		return "中國語"
	} else if lan == "ja" {
		return "にほんご"
	} else if lan == "ko" {
		return statement
	} else {
		panic("Invalid language code")
	}
}

func ExchageCurrency(appReq *framework.AppRequest, lan string, price float64) float64 {
	if lan == "kr" {
		return price
	} else {
		rate := model.GetExchangeRate(appReq, lan)

		// TODO : Trim fule for Currency needs be added
		return price * rate
	}
}
