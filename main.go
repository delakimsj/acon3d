package main

import (
	"net/http"

	"acon3d.com/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, handler.GetListProduct())
	})

	router.Run(":8080")
}
