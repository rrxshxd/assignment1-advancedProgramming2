package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/repository"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	query := `	
		INSERT INTO users (email, username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(query, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE id = $1
`

	var user entity.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE email = $1
`

	var user entity.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetAddresses(userID uint) ([]entity.Address, error) {
	query := `
		SELECT id, user_id, street, city, state, postal_code, country, is_default, created_at, updated_at 
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	defer rows.Close()

	var addresses []entity.Address
	for rows.Next() {
		var address entity.Address
		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.Street,
			&address.City,
			&address.State,
			&address.PostalCode,
			&address.Country,
			&address.IsDefault,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	return addresses, nil
}
