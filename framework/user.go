package framework

import "encoding/json"

type (
	User struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}

	DBMSUser struct {
		UserIdSequence int    `json:"user_id_sequence"`
		Data           []User `json:"data"`
	}
)

func (dbms *DBMSUser) LoadData() {
	if json.Unmarshal(*LoadDBMS("user"), &dbms) != nil {
		panic("Can't load user data")
	}
}

func GetUser(userId int) *User {
	// load deal
	var dbms DBMSUser
	dbms.LoadData()

	for _, d := range dbms.Data {
		if d.UserID == userId {
			return &d
		}
	}

	return nil
}
