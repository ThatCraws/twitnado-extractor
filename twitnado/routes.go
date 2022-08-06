package twitnado

import "github.com/gin-gonic/gin"

// instantiate handler
var handler nadoHandlerInterface = NewNadoHandler()

// define handler methods for routes
type nadoHandlerInterface interface {
	searchQuery(ctx *gin.Context)
	store(ctx *gin.Context)
}

func SetupRoutes(grp *gin.RouterGroup) {
	grp.GET("/search", handler.searchQuery)
	grp.POST("/store", handler.store)
}
