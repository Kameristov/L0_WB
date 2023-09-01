package aggregate

import (
	"L0_WB/internal/entity"
	"encoding/json"
)

type Order struct {
	Order_uid          string          `json:"order_uid"`
	Track_number       string          `json:"track_number"`
	Entry              string          `json:"entry"`
	Delivery           entity.Delivery `json:"delivery"`
	Payment            entity.Payment  `json:"payment"`
	Items              []entity.Item   `json:"items"`
	Locale             string          `json:"locale"`
	Internal_signature string          `json:"internal_signature"`
	Customer_id        string          `json:"customer_id"`
	Delivery_service   string          `json:"delivery_service"`
	Shardkey           string          `json:"shardkey"`
	Sm_id              int             `json:"sm_id"`
	Date_created       string          `json:"date_created"`
	Oof_shard          string          `json:"oof_shard"`
}

func NewOrder(data []byte) (Order, error) {

	order:= Order{}
	json.Unmarshal(data, &order)

	return order, nil
}

func (o *Order) GetOrderID() string {
	return o.Order_uid
}
