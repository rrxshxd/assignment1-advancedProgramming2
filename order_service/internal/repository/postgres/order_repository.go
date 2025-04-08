package postgres

import (
	"database/sql"
	"errors"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/repository"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *entity.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`INSERT INTO orders (user_id, total, status, created_at, updated_at) 
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		order.UserID, order.Total, order.Status, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity, price) 
             VALUES ($1, $2, $3, $4)`,
			order.ID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) FindByID(id uint) (*entity.Order, error) {
	query := `
        SELECT o.id, o.user_id, o.total, o.status, o.created_at, o.updated_at,
               oi.product_id, oi.quantity, oi.price
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.id = $1
    `

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order entity.Order
	var items []entity.OrderItem
	orderFound := false

	for rows.Next() {
		orderFound = true
		var item entity.OrderItem
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Total, &order.Status,
			&order.CreatedAt, &order.UpdatedAt,
			&item.ProductID, &item.Quantity, &item.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if !orderFound {
		return nil, errors.New("order not found")
	}

	order.Items = items
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint) ([]*entity.Order, error) {
	query := `
        SELECT o.id, o.user_id, o.total, o.status, o.created_at, o.updated_at,
               oi.product_id, oi.quantity, oi.price
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[uint]*entity.Order)
	for rows.Next() {
		var orderID uint
		var order entity.Order
		var item entity.OrderItem

		err := rows.Scan(
			&orderID, &order.UserID, &order.Total, &order.Status,
			&order.CreatedAt, &order.UpdatedAt,
			&item.ProductID, &item.Quantity, &item.Price,
		)
		if err != nil {
			return nil, err
		}

		if existingOrder, exists := ordersMap[orderID]; exists {
			existingOrder.Items = append(existingOrder.Items, item)
		} else {
			order.ID = orderID
			order.Items = []entity.OrderItem{item}
			ordersMap[orderID] = &order
		}
	}

	orders := make([]*entity.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(id uint, status entity.OrderStatus) error {
	_, err := r.db.Exec(
		`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}
