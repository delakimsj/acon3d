package model

import (
	"encoding/json"

	"acon3d.com/framework"
)

const (
	ProductDBFile = "data/product.json"
)

type (
	DBProduct struct {
		ProductID   int     `json:"product_id"`
		Title       string  `json:"title"`
		Status      string  `json:"status"` // registered | approved
		Description string  `json:"description"`
		Price       float64 `json:"price"`              // price is in KRW
		FeeType     string  `json:"fee_type,omitempty"` // ratio | fixed
		FeeRate     float64 `json:"fee_rate,omitempty"`
		AuthorID    int     `json:"author_id"`
		EditorID    int     `json:"editor_id"`
	}

	DBMSProduct struct {
		ProductIdSequence int         `json:"product_id_sequence"`
		Data              []DBProduct `json:"data"`
	}
)

func IsValidProductStatus(inp string) bool {
	return inp == "registered" || inp == "approved"
}

func IsValidFeeType(inp string) bool {
	return inp == "ratio" || inp == "fixed"
}

func (dmbs *DBMSProduct) IncreaseSequence() {
	dmbs.ProductIdSequence = dmbs.ProductIdSequence + 1
}

func (dbms *DBMSProduct) LoadData() {
	if json.Unmarshal(*framework.LoadDBMS("product"), &dbms) != nil {
		panic("Can't load product data")
	}
}

func findProductIndex(dbms *DBMSProduct, productId int) int {
	for i, v := range dbms.Data {
		if v.ProductID == productId {
			return i
		}
	}
	return -1 // not found
}

func GetListProduct(appReq *framework.AppRequest, condition *map[string]interface{}) *[]DBProduct {
	var dbms DBMSProduct
	dbms.LoadData()

	ret := make([]DBProduct, 0)
	val, ok := (*condition)["status"]
	if ok {
		for _, v := range dbms.Data {
			if v.Status == val {
				ret = append(ret, v)
			}
		}
	} else {
		return &dbms.Data
	}

	return &ret
}

func GetProduct(appReq *framework.AppRequest, productId int) *DBProduct {
	var dbms DBMSProduct
	dbms.LoadData()

	for _, d := range dbms.Data {
		if d.ProductID == productId {
			return &d
		}
	}

	return nil
}

func CreateProduct(appReq *framework.AppRequest, inp *DBProduct) *DBProduct {
	var dbms DBMSProduct
	dbms.LoadData()

	dbms.IncreaseSequence()

	(*inp).ProductID = dbms.ProductIdSequence

	dbms.Data = append(dbms.Data, *inp)

	framework.WriteToDBFile("product", dbms)

	return inp
}

func UpdateProduct(appReq *framework.AppRequest, inp *DBProduct, productId int) *DBProduct {
	var dbms DBMSProduct
	dbms.LoadData()

	index := findProductIndex(&dbms, productId)
	if index == -1 {
		return nil
	} else {
		dbms.Data[index] = *inp
	}

	framework.WriteToDBFile("product", dbms)

	return &dbms.Data[index]
}
