package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func SocialRoutes(route *gin.Engine) {
	social := route.Group("/socialmedias").Use(middlewares.Authentication())
	social.POST("/", controllers.CreateSocial)
	social.GET("/", controllers.GetSocials)
	social.PUT("/:socialMediaId", controllers.UpdateSocial)
	social.DELETE("/:socialMediaId", controllers.DeleteSocial)

}
