package main

import (
	"log"

	"github.com/ThatCraws/twitnado-extractor/twitnado"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	grp := server.Group("/v1/")
	twitnado.SetupRoutes(grp)

	err := server.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Couldn't start server. Error: %s", err.Error())
	}
}
