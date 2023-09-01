// Package http implements routing paths. Each services in own file.
package http

import (
	"L0_WB/internal/usecase"
	"L0_WB/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.OrderUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/api")
	{
		newOrderRoutes(h, t, l)
	}

	newWeb(handler, t, l)

}
