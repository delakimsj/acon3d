package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"acon3d.com/framework"
	"acon3d.com/model"
)

func isValidDealInput(appReq *framework.AppRequest, inp *model.DBDeal) (bool, string) {
	if !model.IsValidDealStatus(inp.Status) {
		return false, "deal status"
	}

	product := model.GetProduct(appReq, inp.ProductId)
	if product == nil {
		return false, "product id"
	}

	// change fee
	if inp.FeeType != "" {
		if !model.IsValidFeeType(inp.FeeType) {
			return false, "fee_type"
		}

		if inp.FeeRate == 0 {
			return false, "fee_rate"
		}
	} else {
		return false, "fee_type"
	}

	return true, ""
}

func GetListDeal(appReq *framework.AppRequest) *framework.AppResponse {
	condition := make(map[string]interface{})

	// get list of deal
	deals := model.GetListDeal(appReq, &condition)

	// make a response
	return framework.GetOkAppResponse(deals)
}

func PostDeal(appReq *framework.AppRequest) *framework.AppResponse {
	// unmarshaling object
	var body model.DBDeal
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("Proper request required")
	}

	if ok, field := isValidDealInput(appReq, &body); !ok {
		return framework.GetBadRequestAppResponse(
			fmt.Sprintf("Proper %s required", field),
		)
	}

	// create new
	deal := model.CreateDeal(appReq, &body)

	return framework.GetOkAppResponse(deal)
}

func PutDeal(appReq *framework.AppRequest) *framework.AppResponse {
	// unmarshaling object
	var body model.DBDeal
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("error when unmarshaling object")
	}

	if ok, field := isValidDealInput(appReq, &body); !ok {
		return framework.GetBadRequestAppResponse(
			fmt.Sprintf("Proper %s required", field),
		)
	}

	dealId, err := strconv.Atoi((*appReq.PathParams)["deal_id"])
	if err != nil {
		return framework.GetBadRequestAppResponse("Proper deal_id required")
	}

	// update new
	deal := model.UpdateDeal(appReq, &body, dealId)

	return framework.GetOkAppResponse(deal)
}

func RequestDealModification(appReq *framework.AppRequest) *framework.AppResponse {
	type RequestModificationRequest struct {
		DealID int `json:"deal_id"`
	}

	// unmarshaling request
	var body RequestModificationRequest
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("Proper body required")
	}

	// set the deal status 'mod_requested'
	deal := model.GetDeal(appReq, body.DealID)
	if deal == nil {
		return framework.GetNotFoundAppResponse("deal")
	}

	deal.Status = "mod_requested"

	updatedDeal := model.UpdateDeal(appReq, deal, body.DealID)
	if updatedDeal == nil {
		panic("Unexpected server error!!")
	}

	// set the product status 'registered'
	product := model.GetProduct(appReq, deal.ProductId)
	if product == nil {
		panic("Unexpected server error!!")
	}

	product.Status = "registered"

	updatedProduct := model.UpdateProduct(appReq, product, deal.ProductId)
	if updatedProduct == nil {
		panic("Unexpected server error!!")
	}

	return framework.GetOkAppResponse(updatedDeal)
}
