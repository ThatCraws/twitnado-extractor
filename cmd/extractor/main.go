package main

import (
	"log"

	"github.com/ThatCraws/twitnado-extractor/twitnado"
	"github.com/ThatCraws/twitnado-extractor/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitEnv()

	server := gin.Default()
	grp := server.Group("/v1/")
	twitnado.SetupRoutes(grp, utils.GetEnvVal("connection_string", "mongodb://localhost:27017"))

	err := server.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Couldn't start server. Error: %s", err.Error())
	}
}
