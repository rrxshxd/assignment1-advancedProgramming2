package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rrxshxd/assignment1_advProg2/proto/user"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/config"
	grpccontroller "github.com/rrxshxd/assignment1_advProg2/user_service/internal/controller/grpc"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/repository/postgres"
	"github.com/rrxshxd/assignment1_advProg2/user_service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWTSecret, time.Duration(cfg.JWTExpirationHours)*time.Hour)
	grpcServer := grpc.NewServer()
	userServer := grpccontroller.NewUserServer(userUseCase)
	user.RegisterUserServiceServer(grpcServer, userServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	log.Printf("User service is running on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve: %v", err)
	}

}
