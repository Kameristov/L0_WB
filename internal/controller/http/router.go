// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.OrderUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newOrderRoutes(h, t, l)
	}
}
