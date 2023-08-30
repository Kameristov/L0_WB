package http

import (
	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebInterface struct {
	t usecase.OrderUseCase
	l logger.Interface
}

func newWeb(handler *gin.Engine, t usecase.OrderUseCase, l logger.Interface) {
	r := &WebInterface{t, l}
	handler.LoadHTMLGlob("pkg/httpserver/website/*.tmpl")
	handler.GET("/", r.MainPage)

}

func (r *WebInterface) MainPage(c *gin.Context) {

	id, stat := c.GetQuery("orderId")
	if !stat {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Data": "",
		})
		return
	}

	order, err := r.t.Get(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - website - Get")
		//errorResponse(c, http.StatusInternalServerError, "database problems")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Data": "Не корректный номер заказа",
		})
		return
	}
	orderString := structToString(order)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Data": orderString,
	})

}

func structToString(st aggregate.Order) string {
	resultString := " Order Information:"
	resultString = fmt.Sprintf("%s \n Order_uid: %+v", resultString, st.Order_uid)
	resultString = fmt.Sprintf("%s \n Track_number: %+v", resultString, st.Track_number)
	resultString = fmt.Sprintf("%s \n Entry: %+v", resultString, st.Entry)
	resultString = fmt.Sprintf("%s \n Locale: %+v", resultString, st.Locale)
	resultString = fmt.Sprintf("%s \n Internal_signature: %+v", resultString, st.Internal_signature)
	resultString = fmt.Sprintf("%s \n Customer_id: %+v", resultString, st.Customer_id)
	resultString = fmt.Sprintf("%s \n Delivery_service: %+v", resultString, st.Delivery_service)
	resultString = fmt.Sprintf("%s \n Shardkey: %+v", resultString, st.Shardkey)
	resultString = fmt.Sprintf("%s \n Sm_id: %+v", resultString, st.Sm_id)
	resultString = fmt.Sprintf("%s \n Date_created: %+v", resultString, st.Date_created)
	resultString = fmt.Sprintf("%s \n Oof_shard: %+v", resultString, st.Oof_shard)
	//
	resultString = fmt.Sprintf("%s \n\n Delivery:", resultString)
	resultString = fmt.Sprintf("%s \n Name: %+v", resultString, st.Delivery.Name)
	resultString = fmt.Sprintf("%s \n Phone: %+v", resultString, st.Delivery.Phone)
	resultString = fmt.Sprintf("%s \n Zip: %+v", resultString, st.Delivery.Zip)
	resultString = fmt.Sprintf("%s \n City: %+v", resultString, st.Delivery.City)
	resultString = fmt.Sprintf("%s \n Address: %+v", resultString, st.Delivery.Address)
	resultString = fmt.Sprintf("%s \n Region: %+v", resultString, st.Delivery.Region)
	resultString = fmt.Sprintf("%s \n Email: %+v", resultString, st.Delivery.Email)
	//
	resultString = fmt.Sprintf("%s \n\n Payment:", resultString)
	resultString = fmt.Sprintf("%s \n Transaction: %+v", resultString, st.Payment.Transaction)
	resultString = fmt.Sprintf("%s \n RequestID: %+v", resultString, st.Payment.RequestID)
	resultString = fmt.Sprintf("%s \n Currency: %+v", resultString, st.Payment.Currency)
	resultString = fmt.Sprintf("%s \n Provider: %+v", resultString, st.Payment.Provider)
	resultString = fmt.Sprintf("%s \n Amount: %+v", resultString, st.Payment.Amount)
	resultString = fmt.Sprintf("%s \n PaymentDto: %+v", resultString, st.Payment.PaymentDto)
	resultString = fmt.Sprintf("%s \n Bank: %+v", resultString, st.Payment.Bank)
	resultString = fmt.Sprintf("%s \n DeliveryCost: %+v", resultString, st.Payment.DeliveryCost)
	resultString = fmt.Sprintf("%s \n GoodsTotal: %+v", resultString, st.Payment.GoodsTotal)
	resultString = fmt.Sprintf("%s \n CustomFee: %+v", resultString, st.Payment.CustomFee)
	//
	resultString = fmt.Sprintf("%s \n\n Items:", resultString)
	for num, it := range st.Items {
		resultString = fmt.Sprintf("%s \n Item № %d", resultString, num+1)
		resultString = fmt.Sprintf("%s \n ChrtId: %+v", resultString, it.ChrtId)
		resultString = fmt.Sprintf("%s \n TrackNumber: %+v", resultString, it.TrackNumber)
		resultString = fmt.Sprintf("%s \n Price: %+v", resultString, it.Price)
		resultString = fmt.Sprintf("%s \n Rid: %+v", resultString, it.Rid)
		resultString = fmt.Sprintf("%s \n Name: %+v", resultString, it.Name)
		resultString = fmt.Sprintf("%s \n Sale: %+v", resultString, it.Sale)
		resultString = fmt.Sprintf("%s \n Size: %+v", resultString, it.Size)
		resultString = fmt.Sprintf("%s \n TotalPrice: %+v", resultString, it.TotalPrice)
		resultString = fmt.Sprintf("%s \n NmId: %+v", resultString, it.NmId)
		resultString = fmt.Sprintf("%s \n Brand: %+v", resultString, it.Brand)
		resultString = fmt.Sprintf("%s \n Status: %+v", resultString, it.Status)
	}
	return resultString
}
