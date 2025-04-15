package usecase

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserUseCase struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	jwtExpires time.Duration
}

func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpires time.Duration) *UserUseCase {
	return &UserUseCase{userRepo, jwtSecret, jwtExpires}
}

func (uc *UserUseCase) RegisterUser(email, username, password string) (*entity.User, string, error) {
	existingUser, err := uc.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (uc *UserUseCase) Authenticate(email, password string) (*entity.User, string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password: %w", err)
	}

	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	user.Password = "" // For security purposes(not exposing the password)

	return user, token, nil
}

func (uc *UserUseCase) GetUserProfile(userID uint) (*entity.User, []entity.Address, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = "" // Same as in Aunthenticate function

	addresses, err := uc.userRepo.GetAddresses(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	return user, addresses, nil
}

func (uc *UserUseCase) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(uc.jwtExpires).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
