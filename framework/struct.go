package framework

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AppRequest struct {
	Tx          *sql.Tx
	AccountId   int
	PathParams  *map[string]string
	QueryParams *map[string]string
	ReqBody     *string
	UseJwt      bool
}

func (r *AppRequest) GetPathParamId(idName string) int {
	id, err := strconv.Atoi((*r.PathParams)[idName])
	if err != nil {
		panic(fmt.Sprintf("invalid %s", idName))
	}

	return id
}

type AppResponse struct {
	RetMessage string
	Status     int
}

func GetOkAppResponse(retObj interface{}) *AppResponse {
	jsonData, err := json.Marshal(retObj)
	if err != nil {
		panic("fail to transform object to string")
	}

	return &AppResponse{
		Status:     http.StatusOK,
		RetMessage: string(jsonData),
	}
}

func GetNotFoundAppResponse(field string) *AppResponse {
	return &AppResponse{
		Status:     http.StatusNotFound,
		RetMessage: fmt.Sprintf("There's no %s", field),
	}
}

func GetBadRequestAppResponse(message string) *AppResponse {
	return &AppResponse{
		Status:     http.StatusBadRequest,
		RetMessage: fmt.Sprint(message),
	}
}

func GetForbiddenAppResponse(message string) *AppResponse {
	return &AppResponse{
		Status:     http.StatusForbidden,
		RetMessage: fmt.Sprint(message),
	}
}

func GetUnauthorizedAppResponse(message string) *AppResponse {
	return &AppResponse{
		Status:     http.StatusUnauthorized,
		RetMessage: fmt.Sprint(message),
	}
}
