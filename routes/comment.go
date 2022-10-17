package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	comment := route.Group("/comments").Use(middlewares.Authentication())
	comment.POST("/", controllers.CreateComment)
	comment.GET("/", controllers.GetComments)
	comment.PUT("/:commentId", controllers.UpdateComment)
	comment.DELETE("/:commentId", controllers.DeleteComment)

}
