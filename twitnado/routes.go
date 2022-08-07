package twitnado

import (
	"github.com/gin-gonic/gin"
)

// define handler methods for routes
type nadoHandlerInterface interface {
	searchQuery(ctx *gin.Context)
	store(ctx *gin.Context)
}

func SetupRoutes(grp *gin.RouterGroup, connString string) {
	var handler nadoHandlerInterface = NewNadoHandler(connString)

	grp.GET("/search", handler.searchQuery)
	grp.POST("/store", handler.store)
}
