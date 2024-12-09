package routes

import (
	"productManagmentBackend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/products", controllers.CreateProduct)
	router.GET("/products/:id", controllers.GetProduct)
	router.GET("/products", controllers.GetProducts)
}
