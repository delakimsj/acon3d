package main

import (
	"acon3d.com/config"
	"acon3d.com/framework"
	"acon3d.com/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// get system config
	cfg := new(config.Config)
	cfg.GetConfig("config/config.yaml")

	// set middleware
	router := gin.Default()
	router.Use(framework.TransactionMiddleware(&framework.InputTransactionMiddleware{
		// UseAuthHeader: cfg.User.UseAuthHeader,
	}))

	// admin apis
	adminAPI := router.Group("/admin")
	{
		// Product
		adminAPI.GET("/product", func(c *gin.Context) {
			appRes := handler.GetListProduct(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.POST("/product", func(c *gin.Context) {
			appRes := handler.PostProduct(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.GET("/product/:product_id", func(c *gin.Context) {
			appRes := handler.GetProduct(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.PUT("/product/:product_id", func(c *gin.Context) {
			appRes := handler.PutProduct(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.PATCH("/product/:product_id", func(c *gin.Context) { // approve product
			appRes := handler.PatchProduct(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		// Deal
		adminAPI.GET("/deal", func(c *gin.Context) {
			appRes := handler.GetListDeal(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.POST("/deal", func(c *gin.Context) {
			appRes := handler.PostDeal(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		adminAPI.PUT("/deal/:deal_id", func(c *gin.Context) {
			appRes := handler.PutDeal(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		// HTTP
		adminAPI.POST("/request_deal_modification", func(c *gin.Context) {
			appRes := handler.RequestDealModification(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})
	}

	// customer apis
	customerAPI := router.Group("/")
	{
		customerAPI.GET("/deal", func(c *gin.Context) {
			appRes := handler.GetListOnsaleDeal(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})

		customerAPI.GET("/deal/:deal_id", func(c *gin.Context) {
			appRes := handler.GetOnsaleDeal(framework.GetAppRequest(c))

			c.String(appRes.Status, appRes.RetMessage)
		})
	}
	router.Run(":8080")
}
