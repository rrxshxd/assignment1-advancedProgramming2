package grpc

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rrxshxd/assignment1_advProg2/proto/user"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/usecase"
	"time"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	userUseCase *usecase.UserUseCase
}

func NewUserServer(userUseCase *usecase.UserUseCase) *UserServer {
	return &UserServer{userUseCase: userUseCase}
}

func (s *UserServer) RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.UserResponse, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("email, username, and password are required")
	}

	newUser, token, err := s.userUseCase.RegisterUser(req.Email, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(newUser.ID),
			Email:     newUser.Email,
			Username:  newUser.Username,
			CreatedAt: timestampPtrFromTime(newUser.CreatedAt),
		},
		Token: token,
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *user.AuthRequest) (*user.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &user.AuthResponse{
			Success:      false,
			ErrorMessage: "email and password are required",
		}, nil
	}

	u, token, err := s.userUseCase.Authenticate(req.Email, req.Password)
	if err != nil {
		return &user.AuthResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &user.AuthResponse{
		Success: true,
		Token:   token,
		UserId:  uint64(u.ID),
	}, nil
}

func (s *UserServer) GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.UserProfile, error) {
	u, addresses, err := s.userUseCase.GetUserProfile(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	protoAddresses := make([]*user.Address, len(addresses))
	for i, addr := range addresses {
		protoAddresses[i] = &user.Address{
			Id:         uint64(addr.ID),
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
			IsDefault:  addr.IsDefault,
		}
	}

	return &user.UserProfile{
		Id:        uint64(u.ID),
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: timestampPtrFromTime(u.CreatedAt),
		Addresses: protoAddresses,
	}, nil
}

func timestampPtrFromTime(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil
	}
	return ts
}
