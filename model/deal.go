package model

import (
	"encoding/json"

	"acon3d.com/framework"
)

type (
	DBDeal struct {
		DealID    int     `json:"deal_id"`
		Title     string  `json:"title"`
		Status    string  `json:"status"` // on-sale, closed, paused, mod_requested
		ProductId int     `json:"product_id"`
		Price     float64 `json:"price"`
		Remark    string  `json:"remark"`
		Market    string  `json:"market"`             // using if channel is various
		FeeType   string  `json:"fee_type,omitempty"` // ratio | fixed
		FeeRate   float64 `json:"fee_rate,omitempty"`
	}

	DBMSDeal struct {
		DealIdSequence int      `json:"deal_id_sequence"`
		Data           []DBDeal `json:"data"`
	}
)

func (dbms *DBMSDeal) IncreaseSequence() {
	dbms.DealIdSequence = dbms.DealIdSequence + 1
}

func (dbms *DBMSDeal) LoadData() {
	if json.Unmarshal(*framework.LoadDBMS("deal"), &dbms) != nil {
		panic("Can't load deal data")
	}
}

func IsValidDealStatus(inp string) bool {
	return inp == "on-sale" || inp == "closed" || inp == "paused"
}

func findDealIndex(dbms *DBMSDeal, dealId int) int {
	for i, v := range dbms.Data {
		if v.DealID == dealId {
			return i
		}
	}
	return -1 // not found
}

func GetDeal(appReq *framework.AppRequest, dealId int) *DBDeal {
	// load deal
	var dbms DBMSDeal
	dbms.LoadData()

	for _, d := range dbms.Data {
		if d.DealID == dealId {
			return &d
		}
	}

	return nil
}

func GetListOnsaleDeal(appReq *framework.AppRequest) *[]DBDeal {
	var dbms DBMSDeal
	dbms.LoadData()

	var ret []DBDeal
	for _, v := range dbms.Data {
		if v.Status == "on-sale" {
			ret = append(ret, v)
		}
	}

	return &ret
}

func GetOnsaleDeal(appReq *framework.AppRequest, dealId int) *DBDeal {
	var dbms DBMSDeal
	dbms.LoadData()

	for _, v := range dbms.Data {
		if v.Status == "on-sale" && v.DealID == dealId {
			return &v
		}
	}

	return nil
}

func GetListDeal(appReq *framework.AppRequest, condition *map[string]interface{}) *[]DBDeal {
	var dbms DBMSDeal

	dbms.LoadData()

	return &dbms.Data
}

func CreateDeal(appReq *framework.AppRequest, inp *DBDeal) *DBDeal {
	var dbms DBMSDeal
	dbms.LoadData()

	dbms.IncreaseSequence()
	(*inp).DealID = dbms.DealIdSequence

	dbms.Data = append(dbms.Data, *inp)

	framework.WriteToDBFile("deal", dbms)

	return inp
}

func UpdateDeal(appReq *framework.AppRequest, inp *DBDeal, dealId int) *DBDeal {
	var dbms DBMSDeal
	dbms.LoadData()

	index := findDealIndex(&dbms, dealId)
	if index == -1 {
		return nil
	} else {
		dbms.Data[index] = *inp
		dbms.Data[index].DealID = dealId
	}

	framework.WriteToDBFile("deal", dbms)

	return &dbms.Data[index]
}
