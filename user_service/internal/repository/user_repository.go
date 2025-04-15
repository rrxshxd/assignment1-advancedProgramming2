package repository

import "github.com/rrxshxd/assignment1_advProg2/user_service/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	GetAddresses(userID uint) ([]entity.Address, error)
}
