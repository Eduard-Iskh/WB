package order

import (
	"context"
	"errors"
	"fmt"
	domain "wildberies/L0/backend/internal/entify"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db *pgxpool.Pool
}

// создаем ссылку на структуру, возвращаем интерфейс
func NewOrderRepository(db *pgxpool.Pool) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {

	var idToReference = order.OrderUID

	//TODOO одна транзакция

	query_orders := `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	query_deliveries := `
        INSERT INTO deliveries (deliveries_id, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	query_payments := `
        INSERT INTO payments (transaction_id, transaction, request_id, currency, provider, amount, 
            payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	query_items := `
            INSERT INTO items (items_id, chrt_id, track_number, price, rid, name, sale, size,
                total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `

	// Вставка основного заказа
	_, err := r.db.Exec(ctx,
		query_orders,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		// Проверяем, является ли ошибка ошибкой уникальности PostgreSQL
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Код ошибки нарушения уникальности
				return fmt.Errorf("order with UID %s already exists: \n%w", order.OrderUID, err)
			}
		}
		return fmt.Errorf("creating order: %w", err)
	}
	// Вставка данных о доставке

	_, err = r.db.Exec(
		ctx,
		query_deliveries,
		idToReference, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email)

	if err != nil {
		return fmt.Errorf("inserting delivery: %w", err)
	}

	// Вставка данных о платеже

	_, err = r.db.Exec(
		ctx,
		query_payments,
		idToReference, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("inserting payment: %w", err)
	}

	//Вставка элементов заказа
	for _, item := range order.Items {
		_, err = r.db.Exec(
			ctx,
			query_items,
			idToReference, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
			item.Status)

		if err != nil {
			return fmt.Errorf("inserting item: %w", err)
		}
	}

	return nil
}

func (r *orderRepository) GetById(ctx context.Context, id string) (*domain.Order, error) {
	order := &domain.Order{
		Delivery: domain.Delivery{},
		Payment:  domain.Payment{},
		Items:    []domain.Item{},
	}

	// Получение основного заказа и связанных данных
	query := `
        SELECT 
            o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
            o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
            d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
            p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        LEFT JOIN deliveries d ON o.order_uid = d.deliveries_id
        LEFT JOIN payments p ON o.order_uid = p.transaction_id
        WHERE o.order_uid = $1
    `

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no Order with this id: %w", err)
		}
		return nil, fmt.Errorf("getting Order by id: %w", err)
	}

	// Получение элементов заказа
	itemsQuery := `
        SELECT chrt_id, track_number, price, rid, name, sale, size,
               total_price, nm_id, brand, status
        FROM items
        WHERE items_id = $1
    `

	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("getting Items for order: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning Item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating items: %w", err)
	}

	return order, nil
}
