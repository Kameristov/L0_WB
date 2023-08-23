package aggregate

import (
	"L0_EVRONE/internal/entity"
	"encoding/json"
)

type Order struct {
	orderInfo *entity.OrderInfo
	delivery  *entity.Delivery
	payment   *entity.Payment
	items     []*entity.Item
}

func NewOrder(jsonData []byte) (Order, error) {

	order := Order{
		orderInfo: &entity.OrderInfo{},
		delivery:  &entity.Delivery{},
		payment:   &entity.Payment{},
		items:     make([]*entity.Item, 0),
	}

	json.Unmarshal(jsonData, &order)

	return order, nil
}

func (o *Order) GetOrderID() string {
	return o.orderInfo.OrderUid
}
