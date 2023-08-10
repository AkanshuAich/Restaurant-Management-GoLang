package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

func TableRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/Tables",controller.GetTables())
	incomingRoutes.GET("/Tables/:Table_id",controller.GetTable())
	incomingRoutes.POST("/Tables",controller.CreateTable())
	incomingRoutes.PATCH("/Tables/:Table_id", controller.UpdateTable())
}