package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func PhotoRoutes(route *gin.Engine) {
	photo := route.Group("/photos").Use(middlewares.Authentication())
	photo.POST("/", controllers.CreatePhoto)
	photo.GET("/", controllers.GetPhoto)
	photo.PUT("/:photoId", controllers.UpdatePhoto)
	photo.DELETE("/:photoId", controllers.DeletePhoto)
}
