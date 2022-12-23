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

	// get product
	product := model.GetProduct(appReq, productId)
	if product == nil {
		return framework.GetNotFoundAppResponse("product")
	}

	// author only can edit product in status 'registered'
	if appReq.User.Role == "author" {
		if product.Status == "approved" {
			return framework.GetForbiddenAppResponse("You can't modify product in status 'approved'")
		}
	} else {
		product.EditorID = appReq.User.UserID
	}

	// unmarshaling object
	var body model.DBProduct
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("error when unmarshaling object")
	}

	// update
	updatedProduct := model.UpdateProduct(appReq, &body, productId)

	return framework.GetOkAppResponse(updatedProduct)
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

	// update editor
	product.EditorID = appReq.User.UserID

	// update new
	updatedProduct := model.UpdateProduct(appReq, product, productId)

	return framework.GetOkAppResponse(updatedProduct)
}
