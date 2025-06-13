package main

import (
	_ "me-pague/docs"
	"me-pague/internal/db"
	"me-pague/internal/controller"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Me Pague API
// @version 1.0
// @description API simples para registro de pagamentos entre usuários.
// @host localhost:8080
// @BasePath /
func main() {
	db.Init()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/user", controller.CreateUser)
	r.GET("/user/:id", controller.GetUser)

	r.GET("/billing", controller.GetBilling)

	r.POST("/payment", controller.CreatePayment)

	r.Run(":8080")
}
