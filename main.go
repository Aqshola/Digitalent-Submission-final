package main

import (
	"final-project/library"
	"final-project/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Unable to load env")
	}

	library.StartDB()
	router := gin.Default()
	routes.UserRoutes(router)
	routes.PhotoRoutes(router)
	routes.CommentRoutes(router)
	routes.SocialRoutes(router)

	router.Run(":8080")
	fmt.Println("Server Running")
}
