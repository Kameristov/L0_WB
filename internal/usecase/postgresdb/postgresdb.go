package postgresdb

import (
	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/internal/entity"
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
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(orders) - db.Builder: %w", err)
	}

	row := db.Pool.QueryRow(context.Background(), sql, args...)

	err = row.Scan(&e.Order_uid, &e.Track_number, &e.Entry, &e.Locale, &e.Internal_signature, &e.Customer_id, &e.Delivery_service, &e.Shardkey, &e.Sm_id, &e.Date_created, &e.Oof_shard)
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
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(delivery) - db.Builder: %w", err)
	}

	row = db.Pool.QueryRow(context.Background(), sql, args...)

	err = row.Scan(&e.Delivery.Name, &e.Delivery.Phone, &e.Delivery.Zip, &e.Delivery.City, &e.Delivery.Address, &e.Delivery.Region, &e.Delivery.Email)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(delivery) - rows.Scan: %w", err)
	}

	// Select payment
	sql, args, err = db.Builder.
		Select("transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee").
		From("payment").
		Where(squirrel.Eq{"transaction": id}).
		ToSql()
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(payment) - db.Builder: %w", err)
	}

	row = db.Pool.QueryRow(context.Background(), sql, args...)

	err = row.Scan(&e.Payment.Transaction, &e.Payment.RequestID, &e.Payment.Currency, &e.Payment.Provider, &e.Payment.Amount, &e.Payment.PaymentDto, &e.Payment.Bank, &e.Payment.DeliveryCost, &e.Payment.GoodsTotal, &e.Payment.CustomFee)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(payment) - rows.Scan: %w", err)
	}
	// Select items

	sql, args, err = db.Builder.
		Select("chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status").
		From("items").
		Where(squirrel.Eq{"order_uid": id}).
		ToSql()
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(items) - db.Builder: %w", err)
	}
	e.Items = make([]entity.Item, 0, 2)

	rows, err := db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderRepo - PutRep(items) - db.Pool.Query: %w", err)
	}
	for rows.Next() {
		item:= entity.Item{}
		err = rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgRep(items) - rows.Scan: %w", err)
		}
		e.Items = append(e.Items, item)
	}

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
		return fmt.Errorf("OrderRepo - PutRep - db.Builder: %w", err)
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
		return fmt.Errorf("OrderRepo - PutRep - db.Builder: %w", err)
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
		return fmt.Errorf("OrderRepo - PutRep - db.Builder: %w", err)
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
			return fmt.Errorf("OrderRepo - PutRep - db.Builder: %w", err)
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
		return []aggregate.Order{}, fmt.Errorf("OrderRepo - GetPgAllRepos - db.Builder: %w", err)
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

	fmt.Printf("DB oirderIdList : %v \n", oirderIdList)
	for _, orderId := range oirderIdList {
		object, err := db.GetPgRep(orderId)
		if err != nil {
			fmt.Printf("Error db.GetPgRep :%v \n", err)
			continue
		}
		orders = append(orders, object)
	}

	return orders, nil
}
