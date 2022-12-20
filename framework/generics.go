package framework

import (
	"encoding/json"
	"fmt"
)

func DumpObject(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")

	fmt.Println("-------- DUMP OBJECT --------")
	fmt.Println(string(s))
	fmt.Println("-----------------------------")
}
