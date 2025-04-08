package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/repository"
	"strings"
	"time"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *entity.Product) error {
	query := `INSERT INTO products (name, description, category, price, stock, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(query, product.Name, product.Description, product.Category, product.Price, product.Stock, time.Now(), time.Now()).Scan(&product.ID)
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	query := `SELECT id, name, description, category, price, stock, created_at, updated_at 
	          FROM products WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var product entity.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product can't be found: %v", err)
		}
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}

	return &product, nil
}

func (r *productRepository) Update(product *entity.Product) error {
	query := `UPDATE products 
	          SET name = $1, description = $2, category = $3, 
	              price = $4, stock = $5, updated_at = NOW() 
	          WHERE id = $6
	          RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Category,
		product.Price,
		product.Stock,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Product can't be found: %v", err)
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(id uint) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Product with such id can't be found: %d", id)
	}

	return nil
}

func (r *productRepository) FindAll(page, limit int, filters map[string]interface{}) ([]*entity.Product, error) {
	baseQuery := `SELECT id, name, description, category, price, stock, created_at, updated_at 
	              FROM products`

	var args []interface{}
	var whereClauses []string
	var argPos int = 1

	for field, value := range filters {
		switch field {
		case "category":
			whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argPos))
			args = append(args, value)
			argPos++
		case "min_price":
			whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argPos))
			args = append(args, value)
			argPos++
		case "max_price":
			whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argPos))
			args = append(args, value)
			argPos++
		case "name":
			whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argPos))
			args = append(args, "%"+value.(string)+"%")
			argPos++
		}
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		baseQuery += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during product rows iteration: %w", err)
	}

	return products, nil
}
