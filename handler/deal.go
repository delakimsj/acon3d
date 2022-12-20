package handler

import (
	"encoding/json"

	"acon3d.com/framework"
	"acon3d.com/model"
)

type DealListResponse struct {
	DealID       int     `json:"deal_id"`
	Title        string  `json:"title"`
	Status       string  `json:"status"` // on-sale, closed, paused
	Price        float64 `json:"price"`
	Remark       string  `json:"remark"`
	Market       string  `json:"market"` // using if channel is various
	ProductTitle string  `json:"product_title"`
	ProductPrice float64 `json:"product_price"`
}

func GetListDeal(appReq *framework.AppRequest) *framework.AppResponse {
	var ret []DealListResponse
	condition := make(map[string]interface{})

	// get query from params
	lan := (*appReq.QueryParams)["lan"]
	if !IsValidLanguageCode(lan) {
		return framework.GetBadRequestAppResponse("language code is not valid")
	}

	status := (*appReq.QueryParams)["status"]
	if status != "" {
		if !model.IsValidDealStatus(status) {
			return framework.GetBadRequestAppResponse("Proper status code required")
		}

		condition["status"] = status
	}

	// get list of deal
	deals := model.GetListDeal(appReq, &condition)
	for _, d := range *deals {
		product := model.GetProduct(appReq, d.ProductId)

		ret = append(ret, DealListResponse{
			DealID:       d.DealID,
			Title:        TranslateTo(lan, d.Title),
			Status:       d.Status,
			Price:        ExchageCurrency(appReq, lan, d.Price),
			Remark:       TranslateTo(lan, d.Remark),
			Market:       d.Market,
			ProductTitle: TranslateTo(lan, product.Title),
			ProductPrice: ExchageCurrency(appReq, lan, product.Price),
		})
	}

	// make a response
	return framework.GetOkAppResponse(ret)
}

type DealResponse struct {
	DealID  int               `json:"deal_id"`
	Title   string            `json:"title"`
	Status  string            `json:"status"` // on-sale, closed, paused
	Price   float64           `json:"price"`
	Remark  string            `json:"remark"`
	Market  string            `json:"market"` // using if channel is various
	Product model.DBMSProduct `json:"product"`
	// ProductReview model.DBMSProductReview `json:"product_review"`
}

func GetDeal(appReq *framework.AppRequest) *framework.AppResponse {
	return nil
}

func PostDeal(appReq *framework.AppRequest) *framework.AppResponse {
	// unmarshaling object
	var body model.DBDeal
	if json.Unmarshal([]byte(*appReq.ReqBody), &body) != nil {
		return framework.GetBadRequestAppResponse("error when unmarshaling object")
	}

	if !model.IsValidDealStatus(body.Status) {
		return framework.GetBadRequestAppResponse("Proper deal status required")
	}

	product := model.GetProduct(appReq, body.ProductId)
	if product == nil {
		return framework.GetBadRequestAppResponse("Proper product id required")
	}

	// create new
	deal := model.CreateDeal(appReq, &body)

	return framework.GetOkAppResponse(deal)
}

func PutDeal(appReq *framework.AppRequest) *framework.AppResponse {
	return nil
}

func PatchDeal(appReq *framework.AppRequest) *framework.AppResponse {
	return nil
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
