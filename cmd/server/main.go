package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/axopadyani/billing-engine/internal/interface/grpc"
	postgres2 "github.com/axopadyani/billing-engine/internal/repository/adapter/db/postgres"
	"github.com/axopadyani/billing-engine/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	postgresConn, err := postgres2.InitConnection()
	if err != nil {
		log.Fatalf("error initializing postgres connection: %v", err)
	}

	loanRepo := postgres2.NewRepository(postgresConn)
	svc := service.NewService(loanRepo)

	grpcServer := grpc.NewServer(svc)
	listener, err := grpc.InitListener()
	if err != nil {
		log.Fatalf("error initializing grpc listener: %v", err)
	}

	grpcServer.Serve(listener)
}
