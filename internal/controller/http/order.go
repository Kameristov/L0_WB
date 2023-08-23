package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/pkg/logger"
)

type OrderRoutes struct {
	t usecase.OrderUseCase
	l logger.Interface
}

func newOrderRoutes(handler *gin.RouterGroup, t usecase.OrderUseCase, l logger.Interface) {
	r := &OrderRoutes{t, l}

	h := handler.Group("/Order")
	{
		h.GET("/order", r.getOrderById)
	}
}

type doIdRequest struct {
	Id string `json:"id"`
}

type OrderByIdResponse struct {
	Order aggregate.Order
}

// @Summary     Show history
// @Description Show all Order history
// @ID          history
// @Tags  	    Order
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /Order/history [get]
func (r *OrderRoutes) getOrderById(c *gin.Context) {

	var request doIdRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	Orders, err := r.t.Get(c.Request.Context(),request.Id)
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, OrderByIdResponse{Orders})
}
