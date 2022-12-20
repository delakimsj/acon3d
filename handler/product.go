package handler

import (
	"encoding/json"
	"strconv"

	"acon3d.com/framework"
	"acon3d.com/model"
)

func GetProduct(appReq *framework.AppRequest) *framework.AppResponse {
	// check request
	productId, err := strconv.Atoi((*appReq.PathParams)["product_id"])
	if err != nil {
		return framework.GetBadRequestAppResponse("Proper product_id required")
	}

	// make a response
	product := model.GetProduct(appReq, productId)
	if product == nil {
		return framework.GetNotFoundAppResponse("product")
	} else {
		return framework.GetOkAppResponse(product)
	}
}

func GetListProduct(appReq *framework.AppRequest) *framework.AppResponse {
	status := (*appReq.QueryParams)["status"]
	condition := make(map[string]interface{})

	if status != "" {
		if model.IsValidProductStatus(status) {
			condition["status"] = status
		} else {
			return framework.GetBadRequestAppResponse("Status is not valid")
		}
	}

	ret := model.GetListProduct(appReq, &condition)

	// make a response
	return framework.GetOkAppResponse(ret)
}

func PostProduct(appReq *framework.AppRequest) *framework.AppResponse {
	// unmarshaling object
	var body model.DBProduct
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("error when unmarshaling object")
	}

	// status would be 'registered' at first
	body.Status = "registered"
	body.AuthorID = appReq.User.UserID

	// create new
	product := model.CreateProduct(appReq, &body)

	return framework.GetOkAppResponse(product)
}

func PutProduct(appReq *framework.AppRequest) *framework.AppResponse {
	// check request
	productId, err := strconv.Atoi((*appReq.PathParams)["product_id"])
	if err != nil {
		return framework.GetBadRequestAppResponse("Proper product_id required")
	}

	// unmarshaling object
	var body model.DBProduct
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("error when unmarshaling object")
	}

	// create new
	product := model.UpdateProduct(appReq, &body, productId)

	return framework.GetOkAppResponse(product)
}

func PatchProduct(appReq *framework.AppRequest) *framework.AppResponse {
	// check request
	productId, err := strconv.Atoi((*appReq.PathParams)["product_id"])
	if err != nil {
		return framework.GetBadRequestAppResponse("Proper product_id required")
	}

	var body model.DBProduct
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("Proper request body please")
	}

	// get product
	product := model.GetProduct(appReq, productId)
	if product == nil {
		return framework.GetNotFoundAppResponse("product")
	}

	// change status
	if body.Status != "" {
		if !model.IsValidProductStatus(body.Status) {
			return framework.GetBadRequestAppResponse("Proper status required")
		}
		product.Status = body.Status
	}

	// change fee
	if body.FeeType != "" {
		if !model.IsValidFeeType(body.FeeType) {
			return framework.GetBadRequestAppResponse("Proper fee_type required")
		}

		if body.FeeRate == 0 {
			return framework.GetBadRequestAppResponse("Proper fee_rate required")
		}

		product.FeeType = body.FeeType
		product.FeeRate = body.FeeRate
	}

	// update editor
	product.EditorID = appReq.User.UserID

	// update new
	updatedProduct := model.UpdateProduct(appReq, product, productId)

	return framework.GetOkAppResponse(updatedProduct)
}
