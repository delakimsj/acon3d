package framework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const ()

func LoadDBMS(dataType string) *[]byte {
	dataMap := map[string]string{
		"product":       "data/product.json",
		"deal":          "data/deal.json",
		"exchange_rate": "data/exchange_rate.json",
		"user":          "data/user.json",
	}

	dbFile, ok := dataMap[dataType]
	if !ok {
		panic(fmt.Sprintf("Can't load data of %s", dataType))
	}

	jsonFile, err := os.Open(dbFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return &byteValue
}

func WriteToDBFile(dataType string, data interface{}) {
	dataMap := map[string]string{
		"product":       "data/product.json",
		"deal":          "data/deal.json",
		"exchange_rate": "data/exchange_rate.json",
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(dataMap[dataType], file, 0644)
}
