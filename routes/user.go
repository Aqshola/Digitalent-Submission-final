package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	user := route.Group("/users")
	user.POST("/register", controllers.Register)
	user.POST("/login", controllers.Login)
	user.Use(middlewares.Authentication()).PUT("/:userId", controllers.UpdateUser)
	user.Use(middlewares.Authentication()).DELETE("/:userId", controllers.DeleteUser)
}
