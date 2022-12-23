package handler

import (
	"strconv"

	"acon3d.com/framework"
	"acon3d.com/model"
)

func GetListOnsaleDeal(appReq *framework.AppRequest) *framework.AppResponse {
	type OnsaleDealListResponse struct {
		DealID       int     `json:"deal_id"`
		Title        string  `json:"title"`
		Status       string  `json:"status"` // on-sale, closed, paused
		Price        float64 `json:"price"`
		Remark       string  `json:"remark"`
		Market       string  `json:"market"` // using if channel is various
		ProductTitle string  `json:"product_title"`
		ProductPrice float64 `json:"product_price"`
	}

	var ret []OnsaleDealListResponse

	// get query from params
	lan := (*appReq.QueryParams)["lan"]
	if !IsValidLanguageCode(lan) {
		return framework.GetBadRequestAppResponse("language code is not valid")
	}

	// get list of deal
	deals := model.GetListOnsaleDeal(appReq)
	for _, d := range *deals {
		product := model.GetProduct(appReq, d.ProductId)

		ret = append(ret, OnsaleDealListResponse{
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

func GetOnsaleDeal(appReq *framework.AppRequest) *framework.AppResponse {
	type OnsaleDealResponse struct {
		DealID          int             `json:"deal_id"`
		Title           string          `json:"title"`
		Status          string          `json:"status"` // on-sale, closed, paused
		Price           float64         `json:"price"`
		Remark          string          `json:"remark"`
		Market          string          `json:"market"` // using if channel is various
		Product         model.DBProduct `json:"product"`
		ProductReview   interface{}     `json:"product_review"`   // TODO
		ProductFileInfo interface{}     `json:"product_fileinfo"` // TODO
		ProductProperty interface{}     `json:"product_property"` // TODO
	}

	// get query from params
	dealId, err := strconv.Atoi((*appReq.PathParams)["deal_id"])
	if err != nil {
		return framework.GetBadRequestAppResponse("Proper deal_id required")
	}

	lan := (*appReq.QueryParams)["lan"]
	if !IsValidLanguageCode(lan) {
		return framework.GetBadRequestAppResponse("language code is not valid")
	}

	deal := model.GetOnsaleDeal(appReq, dealId)
	if deal == nil {
		return framework.GetNotFoundAppResponse("deal")
	}

	product := model.GetProduct(appReq, deal.ProductId)
	if product == nil {
		panic("Unexpected Server Error!!")
	}

	// make a response
	return framework.GetOkAppResponse(OnsaleDealResponse{
		DealID: deal.DealID,
		Title:  TranslateTo(lan, deal.Title),
		Status: deal.Status,
		Price:  ExchageCurrency(appReq, lan, deal.Price),
		Remark: TranslateTo(lan, deal.Remark),
		Market: deal.Market,
		Product: model.DBProduct{
			ProductID:   product.ProductID,
			Title:       TranslateTo(lan, product.Title),
			Status:      product.Status,
			Description: TranslateTo(lan, product.Description),
			Price:       ExchageCurrency(appReq, lan, product.Price),
			AuthorID:    product.AuthorID,
			EditorID:    product.EditorID,
		},
		ProductReview:   "{TODO}",
		ProductFileInfo: "{TODO}",
		ProductProperty: "{TODO}",
	})
}
