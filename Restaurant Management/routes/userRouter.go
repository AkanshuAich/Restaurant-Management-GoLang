package routes

import (
	"github.com/gin-gonic/gin"
	"golang-restaurant-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/users/:user_id",controller.GetUsers())
	incomingRoutes.POST("/users/signup",controller.Signup())
	incomingRoutes.POST("/users/login",controller.Login())
}