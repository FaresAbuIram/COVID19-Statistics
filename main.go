package main

import (
	"os"

	"github.com/FaresAbuIram/COVID19-Statistics/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	router.Use(cors.Default())
	routes.Setup(router)
	router.Run(":" + port)
}
