package main

import (
	"log/slog"
	"net"
	"os"

	urlgrpc "zipit/internal/url/grpc"
	"zipit/internal/url/repository"
	"zipit/internal/url/service"
	"zipit/pkg/config"
	"zipit/pkg/database"
	"zipit/pkg/logger"
	"zipit/pkg/shortener"

	"github.com/joho/godotenv"
	grpc "google.golang.org/grpc"

	pb "zipit/gen/url"
)

func main() {
	// 1. Initialize Infrastucture: Logger, Env vars
	logger.SetLogger()
	err := godotenv.Load()
	if err != nil {
		slog.Error("could not load env variables", "error", err)
		os.Exit(1)
	}

	// 2. Database Setup
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		slog.Error("failed to load database config file", "error", err)
		os.Exit(1)
	}
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		slog.Error("failed to establish connection to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Initialize Layers: Repository -> Service -> gRPC Handler
	repo := repository.NewPostgresRepository(db)
	shortener := shortener.NewBase62Shortener()
	urlSvc := service.NewUrlSvc(repo, shortener)
	handler := urlgrpc.NewURLHandler(urlSvc)

	// 4. Start Network Listener
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// 5. Setup and Start gRPC Server
	server := grpc.NewServer()
	pb.RegisterURLServiceServer(server, handler)

	slog.Info("url service up & ready!", "port", 5051)
	if err := server.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
