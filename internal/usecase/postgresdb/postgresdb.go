package postgresdb

import (
	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/pkg/postgres"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// OrderRepo -.
type OrderRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

// GetRep -.
func (db *OrderRepo) GetPgRep(id string) (aggregate.Order, error) {
	e := aggregate.Order{}
	// Select orders
	sql, args, err := db.Builder.
		Select("order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard").
		From("orders").
		Where(squirrel.Eq{"order_uid": id}).
		ToSql()
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(orders) - r.Builder: %w", err)
	}

	rows := db.Pool.QueryRow(context.Background(), sql, args...)

	err = rows.Scan(&e.Order_uid, &e.Track_number, &e.Entry, &e.Locale, &e.Internal_signature, &e.Customer_id, &e.Delivery_service, &e.Shardkey, &e.Sm_id, &e.Date_created, &e.Oof_shard)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(orders) - rows.Scan: %w", err)
	}

	// Select delivery
	sql, args, err = db.Builder.
		Select("name, phone, zip, city, address, region, email").
		From("delivery").
		Where(squirrel.Eq{"order_uid": id}).
		ToSql()
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(delivery) - r.Builder: %w", err)
	}

	rows = db.Pool.QueryRow(context.Background(), sql, args...)

	err = rows.Scan( e.Delivery.Name, e.Delivery.Phone, e.Delivery.Zip, e.Delivery.City, e.Delivery.Address, e.Delivery.Region, e.Delivery.Email)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(delivery) - rows.Scan: %w", err)
	}

	// Select payment

	// Select items
	return e, nil
}

// PutRep -.
func (db *OrderRepo) PutPgRep(t aggregate.Order) error {
	// Insert orders
	sql, args, err := db.Builder.
		Insert("orders").
		Columns("order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard").
		Values(t.Order_uid, t.Track_number, t.Entry, t.Locale, t.Internal_signature, t.Customer_id, t.Delivery_service, t.Shardkey, t.Sm_id, t.Date_created, t.Oof_shard).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Builder: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Pool.Exec: %w", err)
	}

	// Insert delivery
	sql, args, err = db.Builder.
		Insert("delivery").
		Columns("order_uid, name, phone, zip, city, address, region, email").
		Values(t.Order_uid, t.Delivery.Name, t.Delivery.Phone, t.Delivery.Zip, t.Delivery.City, t.Delivery.Address, t.Delivery.Region, t.Delivery.Email).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Builder: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Pool.Exec: %w", err)
	}

	// Insert payment
	sql, args, err = db.Builder.
		Insert("payment").
		Columns("transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee").
		Values(t.Payment.Transaction, t.Payment.RequestID, t.Payment.Currency, t.Payment.Provider, t.Payment.Amount, t.Payment.PaymentDto, t.Payment.Bank, t.Payment.DeliveryCost, t.Payment.GoodsTotal, t.Payment.CustomFee).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Builder: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepo - PutRep - r.Pool.Exec: %w", err)
	}

	// Insert items
	for _, item := range t.Items {
		sql, args, err = db.Builder.
			Insert("items").
			Columns("order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status").
			Values(t.Order_uid, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status).
			ToSql()
		if err != nil {
			return fmt.Errorf("OrderRepo - PutRep - r.Builder: %w", err)
		}

		_, err = db.Pool.Exec(context.Background(), sql, args...)
		if err != nil {
			return fmt.Errorf("OrderRepo - PutRep - r.Pool.Exec: %w", err)
		}
	}

	return nil
}

func (db *OrderRepo) GetPgAllRepos() ([]aggregate.Order, error) {
	orders := make([]aggregate.Order, 0, 10)

	sql, _, err := db.Builder.
		Select("order_uid").
		From("orders").
		ToSql()
	if err != nil {
		return []aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgAllRepos - r.Builder: %w", err)
	}

	rows, err := db.Pool.Query(context.Background(), sql)
	if err != nil {
		return []aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgAllRepos - r.Pool.Query: %w", err)
	}
	oirderIdList := make([]string, 0, 10)
	for rows.Next() {
		var oderId string
		err = rows.Scan(&oderId)
		if err != nil {
			return []aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgAllRepos - rows.Scan: %w", err)
		}
		oirderIdList = append(oirderIdList, oderId)
	}
	rows.Close()

	fmt.Printf("oirderIdList : %v \n", oirderIdList)
	for _, orderId := range oirderIdList {
		object, err := db.GetPgRep(orderId)
		if err != nil {
			fmt.Printf("Error db.GetPgRep :%v \n", err)
			continue
		}
		orders = append(orders, object)
	}

	fmt.Printf("orders : %v \n", orders)
	return orders, nil
}
