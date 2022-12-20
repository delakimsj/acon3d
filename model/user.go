package model

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"acon3d.com/framework"
)

const (
	UserDBFile = "data/user.json"
)

type (
	DBUser struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}

	DBMSUser struct {
		UserIdSequence int      `json:"product_id_sequence"`
		Data           []DBUser `json:"data"`
	}
)

func loadDBMSUser() *DBMSUser {
	jsonFile, err := os.Open(ProductDBFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// unmarshaling object
	var body DBMSUser
	if json.Unmarshal([]byte(byteValue), &body) != nil {
		panic("error when unmarshaling object")
	}

	return &body
}

func GetUser(appReq *framework.AppRequest, userId int) *DBUser {
	dbms := loadDBMSUser()

	for _, d := range dbms.Data {
		if d.UserID == userId {
			return &d
		}
	}

	return nil
}
