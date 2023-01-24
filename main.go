package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olusamimaths/go-jwt/controller"
	"github.com/olusamimaths/go-jwt/middleware"
	"github.com/olusamimaths/go-jwt/model"
)


func init() {
	model.SetDBClient()
}

func main() {
	fmt.Println("GO JWT Sample Project")
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	router.POST("/signup", controller.Signup)
	router.POST("/login",  controller.Login)
	router.GET("/api/v1", middleware.Authorize, controller.Resources)

	router.Run(":9000")
}